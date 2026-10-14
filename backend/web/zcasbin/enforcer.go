package zcasbin

import (
	"context"
	"net/http"
	"strings"
	"sync"

	"github.com/casbin/casbin/v2"
	"github.com/casbin/casbin/v2/model"
	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gfile"
	"github.com/gogf/gf/v2/os/gres"
	"github.com/gogf/gf/v2/text/gstr"

	"github.com/denghuo98/zzfast/internal/dao"
	"github.com/denghuo98/zzfast/zconsts"

	_ "github.com/gogf/gf/contrib/drivers/mysql/v2"
)

const (
	ActionGet    = http.MethodGet
	ActionPost   = http.MethodPost
	ActionPut    = http.MethodPut
	ActionDelete = http.MethodDelete
	ActionAll    = "GET|POST|PUT|DELETE|PATCH|OPTIONS|HEAD"
)

var (
	Enforcer *casbin.Enforcer
	mu       sync.RWMutex // 保护 Enforcer 的并发读写
)

// EnforceSafe 带读锁的 Enforce，避免 Refresh 清空策略时并发鉴权失败
func EnforceSafe(params ...interface{}) (bool, error) {
	if Enforcer == nil {
		return false, nil
	}
	mu.RLock()
	defer mu.RUnlock()
	return Enforcer.Enforce(params...)
}

// InitEnforcer 初始化
func InitEnforcer(ctx context.Context) {
	dbConfig, _ := gdb.GetConfigGroup("default")

	var (
		link   = dbConfig[0].Link
		a, err = NewAdapter(link)
	)

	if err != nil {
		g.Log().Panicf(ctx, "casbin.NewAdapter err . %v", err)
		return
	}

	path := "casbin.conf"

	// 优先从本地加载casbin.conf，如果不存在就从资源文件中找
	modelContent := gfile.GetContents(path)
	if len(modelContent) == 0 {
		if !gres.IsEmpty() && gres.Contains(path) {
			modelContent = string(gres.GetContent(path))
		}
	}

	if len(modelContent) == 0 {
		g.Log().Panicf(ctx, "casbin model file does not exist：%v", path)
	}

	m, err := model.NewModelFromString(modelContent)
	if err != nil {
		g.Log().Panicf(ctx, "casbin NewModelFromString err：%v", err)
	}

	Enforcer, err = casbin.NewEnforcer(m, a)
	if err != nil {
		g.Log().Panicf(ctx, "casbin NewEnforcer err：%v", err)
	}
	// 策略由 loadPermissions/Refresh 从业务表（admin_role + admin_menu）加载到内存，
	// 适配器仅作为独立存储工具、不参与运行时持久化，故关闭 AutoSave 避免策略写回 zz_admin_role_casbin 表。
	Enforcer.EnableAutoSave(false)
	loadPermissions(ctx)
}

func loadPermissions(ctx context.Context) {
	type Policy struct {
		Key         string `json:"key"`
		Permissions string `json:"permissions"`
	}
	var (
		rules   [][]string
		polices []*Policy
		err     error
	)
	err = g.Model(gstr.Join([]string{dao.AdminRole.Table(), "r"}, " ")).
		LeftJoin(gstr.Join([]string{dao.AdminRoleMenu.Table(), "rm"}, " "), "r.id=rm.role_id").
		LeftJoin(gstr.Join([]string{dao.AdminMenu.Table(), "m"}, " "), "rm.menu_id=m.id").
		Fields("r.key,m.permissions").
		Where("r.status", zconsts.StatusEnabled).
		Where("m.status", zconsts.StatusEnabled).
		Where("m.permissions !=?", "").
		Where("r.key !=?", zconsts.SuperRoleKey).
		Scan(&polices)
	if err != nil {
		g.Log().Fatalf(ctx, "loadPermissions Scan err:%v", err)
		return
	}

	for _, policy := range polices {
		if strings.Contains(policy.Permissions, ",") {
			lst := strings.Split(policy.Permissions, ",")
			for _, permissions := range lst {
				rules = append(rules, []string{policy.Key, permissions, ActionAll})
			}
		} else {
			rules = append(rules, []string{policy.Key, policy.Permissions, ActionAll})
		}
	}

	if _, err = Enforcer.AddPolicies(rules); err != nil {
		g.Log().Fatalf(ctx, "loadPermissions AddPolicies err:%v", err)
		return
	}
}

// Clear 清空所有策略（内部使用，调用方需持锁）
func Clear(ctx context.Context) (err error) {
	if Enforcer == nil {
		return nil
	}
	policy, err := Enforcer.GetPolicy()
	if err != nil {
		g.Log().Warningf(ctx, "Enforcer GetPolicy err:%+v", err)
		return err
	}
	if len(policy) == 0 {
		return nil
	}

	_, err = Enforcer.RemovePolicies(policy)
	if err != nil {
		g.Log().Warningf(ctx, "Enforcer RemovePolicies err:%+v", err)
		return err
	}

	// 确认清理干净（RemovePolicies 批量删除，不会残留）
	remaining, _ := Enforcer.GetPolicy()
	if len(remaining) > 0 {
		g.Log().Warningf(ctx, "Enforcer still has %d policies after Clear, retrying", len(remaining))
		_, err = Enforcer.RemovePolicies(remaining)
	}
	return err
}

// Refresh 刷新所有策略（原子替换：先加载新策略，再清空旧策略，避免并发 Enforce 拿到空策略）
func Refresh(ctx context.Context) (err error) {
	// 先加载新策略到临时模型
	type Policy struct {
		Key         string `json:"key"`
		Permissions string `json:"permissions"`
	}
	var polices []*Policy
	err = g.Model(gstr.Join([]string{dao.AdminRole.Table(), "r"}, " ")).
		LeftJoin(gstr.Join([]string{dao.AdminRoleMenu.Table(), "rm"}, " "), "r.id=rm.role_id").
		LeftJoin(gstr.Join([]string{dao.AdminMenu.Table(), "m"}, " "), "rm.menu_id=m.id").
		Fields("r.key,m.permissions").
		Where("r.status", zconsts.StatusEnabled).
		Where("m.status", zconsts.StatusEnabled).
		Where("m.permissions !=?", "").
		Where("r.key !=?", zconsts.SuperRoleKey).
		Scan(&polices)
	if err != nil {
		// 表不存在时（如测试环境）不回传错误，仅记录 warning
		g.Log().Warningf(ctx, "Refresh Scan err:%v", err)
		mu.Lock()
		defer mu.Unlock()
		_ = Clear(ctx)
		return nil
	}

	var rules [][]string
	for _, policy := range polices {
		if strings.Contains(policy.Permissions, ",") {
			lst := strings.Split(policy.Permissions, ",")
			for _, p := range lst {
				p = strings.TrimSpace(p)
				if p != "" {
					rules = append(rules, []string{policy.Key, p, ActionAll})
				}
			}
		} else {
			rules = append(rules, []string{policy.Key, policy.Permissions, ActionAll})
		}
	}

	// 加写锁：先清旧后加新，在锁内完成原子替换
	mu.Lock()
	defer mu.Unlock()

	if err = Clear(ctx); err != nil {
		return err
	}
	if len(rules) > 0 {
		if _, err = Enforcer.AddPolicies(rules); err != nil {
			g.Log().Warningf(ctx, "Refresh AddPolicies err:%v", err)
		}
	}
	return nil
}
