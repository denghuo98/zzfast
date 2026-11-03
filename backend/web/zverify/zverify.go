// Package zverify 提供验证码的生成、存储、限流与校验能力。
//
// 该包只负责验证码状态机，不关心短信厂商：发送动作通过 Sender 注入，
// 业务场景策略由调用方传入 Policy。验证码默认存储在 zcache 中，
// 校验成功后立即失效（一次性使用）。
package zverify

import (
	"context"
	"crypto/hmac"
	crand "crypto/rand"
	"crypto/sha256"
	"crypto/subtle"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"math/big"
	"regexp"
	"sort"
	"strings"
	"time"

	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gcache"
	"github.com/gogf/gf/v2/util/grand"

	"github.com/denghuo98/zzfast/web/zcache"
	"github.com/denghuo98/zzfast/zconsts"
)

// Sender 验证码发送器。
type Sender interface {
	Send(ctx context.Context, mobile, template, code string) error
}

// SendError 表示短信供应商发送失败，调用方可据此对外返回通用错误并记录内部详情。
type SendError struct {
	Err error
}

func (e *SendError) Error() string {
	return "验证码发送失败"
}

func (e *SendError) Unwrap() error {
	return e.Err
}

// RejectError 表示可直接展示给用户的业务拒绝（冷却、限流、验证码错误等）。
type RejectError struct {
	Message string
}

func (e *RejectError) Error() string {
	return e.Message
}

func reject(message string) error {
	return &RejectError{Message: message}
}

func rejectf(format string, args ...any) error {
	return &RejectError{Message: fmt.Sprintf(format, args...)}
}

var mobilePattern = regexp.MustCompile(`^\+?[1-9]\d{5,14}$`)

// Policy 单次发送/校验使用的策略。
type Policy struct {
	Scene          string
	Template       string
	HashKey        string
	Length         int
	TTL            time.Duration
	Cooldown       time.Duration
	MaxAttempts    int
	MaxIpDaily     int
	MaxMobileDaily int
}

// Service 验证码服务。
type Service struct {
	cache  *gcache.Cache
	sender Sender
	now    func() time.Time
}

// New 创建验证码服务。
func New(cache *gcache.Cache, sender Sender) *Service {
	return &Service{
		cache:  cache,
		sender: sender,
		now:    time.Now,
	}
}

type record struct {
	Hash     string `json:"hash"`
	Attempts int    `json:"attempts"`
	SentAt   int64  `json:"sentAt"`
	ExpireAt int64  `json:"expireAt"`
}

// Send 生成并发送验证码，返回有效期（秒）。
func (s *Service) Send(ctx context.Context, p Policy, mobile, ip string) (expiresIn int, err error) {
	if s.sender == nil {
		return 0, gerror.New("未配置验证码发送器")
	}
	p = p.withDefaults()
	if p.HashKey == "" {
		return 0, gerror.New("验证码哈希密钥未配置")
	}
	mobile = NormalizeMobile(mobile)

	if mobile == "" {
		return 0, reject("手机号不能为空")
	}
	if !mobilePattern.MatchString(mobile) {
		return 0, reject("手机号格式不正确")
	}
	if p.Template == "" {
		return 0, reject("当前场景未配置短信模板")
	}

	lockKeys := []string{lockKey(p.Scene, mobile)}
	if ip != "" {
		lockKeys = append(lockKeys, ipLockKey(p.Scene, ip, s.now().Format("20060102")))
	}

	var (
		code     string
		codeHash string
		date     string
	)
	err = s.withLocks(ctx, lockKeys, 10*time.Second, func() error {
		now := s.now()

		if v, err := s.cache.Get(ctx, cooldownKey(p.Scene, mobile)); err != nil && !zcache.IsExpired(err) {
			return gerror.Wrap(err, "验证码冷却检查失败")
		} else if v != nil && !v.IsEmpty() {
			next := time.Unix(v.Int64(), 0).Add(p.Cooldown)
			if now.Before(next) {
				wait := int(next.Sub(now).Seconds()) + 1
				return rejectf("发送过于频繁，请 %d 秒后重试", wait)
			}
		}

		date = now.Format("20060102")
		if err := s.checkDaily(ctx, mobileDailyKey(p.Scene, mobile, date), p.MaxMobileDaily, "该手机号今日发送次数过多，请明日再试"); err != nil {
			return err
		}
		if ip != "" {
			if err := s.checkDaily(ctx, ipDailyKey(p.Scene, ip, date), p.MaxIpDaily, "当前 IP 今日发送次数过多，请明日再试"); err != nil {
				return err
			}
		}

		generated, codeErr := secureDigits(p.Length)
		if codeErr != nil {
			return gerror.Wrap(codeErr, "验证码生成失败")
		}
		code = generated
		codeHash = hashCode(p.HashKey, p.Scene, mobile, code)
		rec := &record{
			Hash:     codeHash,
			Attempts: 0,
			SentAt:   now.Unix(),
			ExpireAt: now.Add(p.TTL).Unix(),
		}
		raw, _ := json.Marshal(rec)
		if err := s.cache.Set(ctx, codeKey(p.Scene, mobile), string(raw), p.TTL); err != nil {
			return gerror.Wrap(err, "验证码保存失败")
		}
		if err := s.cache.Set(ctx, cooldownKey(p.Scene, mobile), now.Unix(), p.Cooldown); err != nil {
			_, _ = s.cache.Remove(ctx, codeKey(p.Scene, mobile))
			return gerror.Wrap(err, "验证码冷却保存失败")
		}

		if err := s.incrDaily(ctx, mobileDailyKey(p.Scene, mobile, date)); err != nil {
			g.Log().Warningf(ctx, "验证码手机号每日计数写入失败 scene=%s err=%v", p.Scene, err)
		}
		if ip != "" {
			if err := s.incrDaily(ctx, ipDailyKey(p.Scene, ip, date)); err != nil {
				g.Log().Warningf(ctx, "验证码 IP 每日计数写入失败 scene=%s err=%v", p.Scene, err)
			}
		}

		expiresIn = int(p.TTL.Seconds())
		return nil
	})
	if err != nil {
		return 0, err
	}

	if err = s.sender.Send(ctx, mobile, p.Template, code); err != nil {
		if decErr := s.withLocks(ctx, lockKeys, 5*time.Second, func() error {
			rec, _ := s.getRecord(ctx, p.Scene, mobile)
			if rec == nil || rec.Hash == codeHash {
				_, _ = s.cache.Remove(ctx, codeKey(p.Scene, mobile))
				_, _ = s.cache.Remove(ctx, cooldownKey(p.Scene, mobile))
				if date != "" {
					s.decrDaily(ctx, mobileDailyKey(p.Scene, mobile, date))
					if ip != "" {
						s.decrDaily(ctx, ipDailyKey(p.Scene, ip, date))
					}
				}
			}
			return nil
		}); decErr != nil {
			g.Log().Warningf(ctx, "验证码发送失败后清理状态失败 scene=%s err=%v", p.Scene, decErr)
		}
		return 0, &SendError{Err: err}
	}
	return expiresIn, nil
}

// Verify 校验验证码，成功后立即失效。
func (s *Service) Verify(ctx context.Context, p Policy, mobile, code string) error {
	p = p.withDefaults()
	if p.HashKey == "" {
		return gerror.New("验证码哈希密钥未配置")
	}
	mobile = NormalizeMobile(mobile)
	code = strings.TrimSpace(code)

	if mobile == "" {
		return reject("手机号不能为空")
	}
	if !mobilePattern.MatchString(mobile) {
		return reject("手机号格式不正确")
	}
	if code == "" {
		return reject("验证码不能为空")
	}

	return s.withLocks(ctx, []string{lockKey(p.Scene, mobile)}, 10*time.Second, func() error {
		now := s.now()
		key := codeKey(p.Scene, mobile)

		rec, err := s.getRecord(ctx, p.Scene, mobile)
		if err != nil {
			return err
		}
		if rec == nil {
			return reject("验证码错误或已过期")
		}

		if now.Unix() > rec.ExpireAt {
			_, _ = s.cache.Remove(ctx, key)
			return reject("验证码已过期，请重新发送")
		}

		if rec.Attempts >= p.MaxAttempts {
			_, _ = s.cache.Remove(ctx, key)
			return reject("验证码错误次数过多，请重新发送")
		}

		want := []byte(rec.Hash)
		got := []byte(hashCode(p.HashKey, p.Scene, mobile, code))
		if subtle.ConstantTimeCompare(want, got) != 1 {
			rec.Attempts++
			remaining := time.Unix(rec.ExpireAt, 0).Sub(now)
			if remaining < time.Second {
				remaining = time.Second
			}
			raw, _ := json.Marshal(rec)
			if err = s.cache.Set(ctx, key, string(raw), remaining); err != nil {
				return gerror.Wrap(err, "验证码状态更新失败")
			}
			return reject("验证码错误")
		}

		if _, err = s.cache.Remove(ctx, key); err != nil {
			return gerror.Wrap(err, "验证码状态更新失败")
		}
		return nil
	})
}

func (p Policy) withDefaults() Policy {
	if p.Length <= 0 {
		p.Length = zconsts.DefaultVerifyCodeLength
	}
	if p.TTL <= 0 {
		p.TTL = zconsts.DefaultVerifyCodeTTL * time.Second
	}
	if p.Cooldown <= 0 {
		p.Cooldown = zconsts.DefaultVerifyCodeCooldown * time.Second
	}
	if p.MaxAttempts <= 0 {
		p.MaxAttempts = zconsts.DefaultVerifyCodeMaxAttempts
	}
	if p.MaxIpDaily == 0 {
		p.MaxIpDaily = zconsts.DefaultVerifyCodeMaxIpDaily
	}
	if p.MaxMobileDaily == 0 {
		p.MaxMobileDaily = zconsts.DefaultVerifyCodeMaxMobDaily
	}
	return p
}

func (s *Service) withLocks(ctx context.Context, keys []string, ttl time.Duration, fn func() error) error {
	sort.Strings(keys)
	tokens := make(map[string]string, len(keys))
	releaseCtx := context.WithoutCancel(ctx)
	release := func() {
		for key, token := range tokens {
			v, err := s.cache.Get(releaseCtx, key)
			if err == nil && v != nil && v.String() == token {
				_, _ = s.cache.Remove(releaseCtx, key)
			}
		}
	}

	for _, key := range keys {
		token := grand.S(16)
		acquired := false
		for attempt := 0; attempt < 3; attempt++ {
			ok, err := s.cache.SetIfNotExist(ctx, key, token, ttl)
			if err != nil {
				release()
				return gerror.Wrap(err, "验证码操作失败")
			}
			if ok {
				acquired = true
				break
			}
			if attempt < 2 {
				time.Sleep(20 * time.Millisecond)
			}
		}
		if !acquired {
			release()
			return reject("操作过于频繁，请稍后再试")
		}
		tokens[key] = token
	}
	defer release()

	return fn()
}

func (s *Service) getRecord(ctx context.Context, scene, mobile string) (*record, error) {
	v, err := s.cache.Get(ctx, codeKey(scene, mobile))
	if err != nil {
		if zcache.IsExpired(err) {
			return nil, nil
		}
		return nil, gerror.Wrap(err, "验证码读取失败")
	}
	if v == nil || v.IsEmpty() {
		return nil, nil
	}
	rec := &record{}
	if err = json.Unmarshal(v.Bytes(), rec); err != nil {
		g.Log().Warningf(ctx, "验证码缓存解析失败 scene=%s err=%v", scene, err)
		_, _ = s.cache.Remove(ctx, codeKey(scene, mobile))
		return nil, nil
	}
	return rec, nil
}

func (s *Service) checkDaily(ctx context.Context, key string, limit int, message string) error {
	if limit <= 0 {
		return nil
	}
	v, err := s.cache.Get(ctx, key)
	if err != nil {
		if zcache.IsExpired(err) {
			return nil
		}
		return gerror.Wrap(err, "验证码限流检查失败")
	}
	if v != nil && v.Int() >= limit {
		return reject(message)
	}
	return nil
}

func (s *Service) incrDaily(ctx context.Context, key string) error {
	v, err := s.cache.Get(ctx, key)
	if err != nil && !zcache.IsExpired(err) {
		return err
	}
	count := 0
	if v != nil {
		count = v.Int()
	}
	return s.cache.Set(ctx, key, count+1, untilEndOfDay(s.now()))
}

func (s *Service) decrDaily(ctx context.Context, key string) {
	v, err := s.cache.Get(ctx, key)
	if err != nil && !zcache.IsExpired(err) {
		return
	}
	count := 0
	if v != nil {
		count = v.Int()
	}
	if count <= 1 {
		_, _ = s.cache.Remove(ctx, key)
		return
	}
	_ = s.cache.Set(ctx, key, count-1, untilEndOfDay(s.now()))
}

func secureDigits(length int) (string, error) {
	if length <= 0 {
		return "", gerror.New("验证码长度必须大于 0")
	}
	buf := make([]byte, length)
	for i := 0; i < length; i++ {
		n, err := crand.Int(crand.Reader, big.NewInt(10))
		if err != nil {
			return "", err
		}
		buf[i] = byte('0' + n.Int64())
	}
	return string(buf), nil
}

func hashCode(key, scene, mobile, code string) string {
	payload := []byte(scene + "\x00" + mobile + "\x00" + code)
	mac := hmac.New(sha256.New, []byte(key))
	_, _ = mac.Write(payload)
	return hex.EncodeToString(mac.Sum(nil))
}

func codeKey(scene, mobile string) string {
	return "verify:code:" + scene + ":" + mobile
}

func cooldownKey(scene, mobile string) string {
	return "verify:cooldown:" + scene + ":" + mobile
}

func lockKey(scene, mobile string) string {
	return "verify:lock:" + scene + ":" + mobile
}

// normalizeMobile 统一手机号表示：去掉空格/横线，并把中国区号 +86/0086 转成 11 位本地号码。
func NormalizeMobile(mobile string) string {
	mobile = strings.TrimSpace(mobile)
	mobile = strings.ReplaceAll(mobile, " ", "")
	mobile = strings.ReplaceAll(mobile, "-", "")
	switch {
	case strings.HasPrefix(mobile, "+86"):
		return strings.TrimPrefix(mobile, "+86")
	case strings.HasPrefix(mobile, "0086"):
		return strings.TrimPrefix(mobile, "0086")
	case strings.HasPrefix(mobile, "86") && len(mobile) == 13:
		return strings.TrimPrefix(mobile, "86")
	default:
		return mobile
	}
}

func ipLockKey(scene, ip, date string) string {
	return "verify:lock:ip:" + scene + ":" + ip + ":" + date
}

func mobileDailyKey(scene, mobile, date string) string {
	return "verify:daily:mobile:" + scene + ":" + mobile + ":" + date
}

func ipDailyKey(scene, ip, date string) string {
	return "verify:daily:ip:" + scene + ":" + ip + ":" + date
}

func untilEndOfDay(now time.Time) time.Duration {
	next := time.Date(now.Year(), now.Month(), now.Day()+1, 0, 0, 0, 0, now.Location())
	return next.Sub(now)
}
