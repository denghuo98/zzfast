package common

import (
	"github.com/gogf/gf/v2/encoding/gbase64"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"

	"github.com/denghuo98/zzfast/internal/dao"
	"github.com/denghuo98/zzfast/internal/model/entity"
	"github.com/denghuo98/zzfast/web/zcontext"
	"github.com/denghuo98/zzfast/web/zencrypt"
	"github.com/denghuo98/zzfast/web/ztoken"
	"github.com/denghuo98/zzfast/zconsts"
	adminSchema "github.com/denghuo98/zzfast/zschema/admin"
	commonSchema "github.com/denghuo98/zzfast/zschema/common"
	systemSchema "github.com/denghuo98/zzfast/zschema/system"
	"github.com/denghuo98/zzfast/zschema/zweb"
	"github.com/denghuo98/zzfast/zservice"
)

type sCommonSite struct {
}

func init() {
	zservice.RegisterCommonSite(NewCommonSite())
}

func NewCommonSite() *sCommonSite {
	return &sCommonSite{}
}

// Ping 心跳检测
func (s *sCommonSite) Ping(ctx g.Ctx) (status string) {
	return "Pong"
}

func (s *sCommonSite) InitSuperAdmin(ctx g.Ctx) (err error) {
	conf, err := zservice.SystemConfig().GetSuperAdmin(ctx)
	if err != nil {
		return err
	}
	// 先判断是否存在角色
	var roleEntity *entity.AdminRole
	if err = dao.AdminRole.Ctx(ctx).Where(dao.AdminRole.Columns().Key, zconsts.SuperRoleKey).Scan(&roleEntity); err != nil {
		return gerror.Wrap(err, zconsts.ErrorORM)
	}
	if roleEntity == nil {
		roleEntity, err = zservice.AdminRole().Edit(ctx, adminSchema.RoleEditInput{
			AdminRole: entity.AdminRole{
				Name:   zconsts.SuperRoleKey,
				Key:    zconsts.SuperRoleKey,
				Remark: "系统超级管理员角色，禁止修改删除",
				Sort:   1,
				Status: 1,
			},
		})
		if err != nil {
			return err
		}
		g.Log().Info(ctx, "已自动创建超级管理员角色")
	}

	// 如果配置文件中没有设置密码，使用默认值
	if conf.Password == "" {
		conf.Password = zconsts.DefaultSuperAdminPassword
		g.Log().Warningf(ctx, "未配置超级管理员密码，使用默认密码，生产环境请务必通过 system.superAdmin.password 覆盖")
	}
	// 判断是否存在，不存在则创建
	var memberEntity *entity.AdminMember
	if err = dao.AdminMember.Ctx(ctx).Where(dao.AdminMember.Columns().Username, zconsts.DefaultSuperAdminUsername).Scan(&memberEntity); err != nil {
		return gerror.Wrap(err, zconsts.ErrorORM)
	}
	if memberEntity == nil {
		payload := adminSchema.MemberEditInput{
			Username: zconsts.DefaultSuperAdminUsername,
			Password: conf.Password,
			RealName: "超级管理员",
			Email:    "",
			Sex:      3,
			RoleIds:  []int64{roleEntity.Id},
			Remark:   "系统超级管理员，禁止编辑或删除",
			Status:   1,
		}
		err = zservice.AdminMember().Edit(ctx, &payload)
		if err != nil {
			return err
		}
		g.Log().Info(ctx, "已自动创建超级管理员用户")
	} else {
		// 仅当配置密码与现有哈希不匹配时更新（bcrypt 自带盐）
		if !zencrypt.VerifyPassword(conf.Password, memberEntity.PasswordHash) {
			hash, hashErr := zencrypt.HashPassword(conf.Password)
			if hashErr != nil {
				return hashErr
			}
			if _, err = dao.AdminMember.Ctx(ctx).WherePri(memberEntity.Id).Data(g.Map{
				dao.AdminMember.Columns().PasswordHash: hash,
			}).Update(); err != nil {
				return gerror.Wrap(err, zconsts.ErrorORM)
			}
			g.Log().Info(ctx, "已自动更新超级管理员密码")
		}
	}

	return nil
}

func (s *sCommonSite) AccountLogin(ctx g.Ctx, in *commonSchema.SiteAccountLoginInput) (out *commonSchema.SiteLoginOutput, err error) {
	defer func() {
		// 推送登录事件
		zservice.SysLoginLog().Push(ctx, &systemSchema.SysLoginLogPushInput{
			Response: out,
			Error:    err,
		})
	}()

	var mb *entity.AdminMember
	if err = dao.AdminMember.Ctx(ctx).Where(dao.AdminMember.Columns().Username, in.Username).Scan(&mb); err != nil {
		err = gerror.Wrap(err, zconsts.ErrorORM)
		return
	}

	if mb == nil {
		err = gerror.New("账号不存在")
		return
	}

	if mb.Status != zconsts.StatusEnabled {
		err = gerror.New("账号已禁用")
		return
	}

	out = new(commonSchema.SiteLoginOutput)
	out.Id = mb.Id
	out.Username = mb.Username

	key, err := zservice.SystemConfig().GetRequestEncryptKey(ctx)
	if err != nil {
		return
	}

	if err = s.verifyPassword(in.Password, key, mb.PasswordHash); err != nil {
		return
	}

	return s.handleLogin(ctx, mb)
}

// verifyPassword 验证密码
func (s *sCommonSite) verifyPassword(input string, key []byte, hash string) (err error) {
	// 先解密文本
	unBase64, err := gbase64.Decode([]byte(input))
	if err != nil {
		return err
	}

	unAes, err := zencrypt.AesGCMDecrypt(unBase64, key)
	if err != nil {
		return err
	}

	if !zencrypt.VerifyPassword(string(unAes), hash) {
		return gerror.New("用户密码错误")
	}

	return nil
}

func (s *sCommonSite) handleLogin(ctx g.Ctx, mb *entity.AdminMember) (out *commonSchema.SiteLoginOutput, err error) {
	// 获取用户绑定的角色
	roles, err := s.getMemberRoles(ctx, mb.Id)
	if err != nil {
		return nil, err
	}

	identity := &zweb.Identity{
		Id:       mb.Id,
		Roles:    roles,
		Username: mb.Username,
		RealName: mb.RealName,
		Avatar:   mb.Avatar,
		Email:    mb.Email,
		Mobile:   mb.Mobile,
		App:      "app",
		LoginAt:  gtime.Now(),
	}

	tokenStr, expires, err := ztoken.Login(ctx, identity)
	if err != nil {
		return nil, err
	}
	out = new(commonSchema.SiteLoginOutput)
	out.Id = mb.Id
	out.Username = mb.Username
	out.Token = tokenStr
	out.Expires = expires

	return out, nil
}

func (s *sCommonSite) getMemberRoles(ctx g.Ctx, id int64) (roles []*zweb.IdentityRole, err error) {
	roleIds, err := dao.AdminMemberRole.Ctx(ctx).Where(dao.AdminMemberRole.Columns().MemberId, id).Array(dao.AdminMemberRole.Columns().RoleId)
	if err != nil {
		return nil, gerror.Wrap(err, zconsts.ErrorORM)
	}

	var roleEntities []*entity.AdminRole
	if err = dao.AdminRole.Ctx(ctx).WhereIn(dao.AdminRole.Columns().Id, roleIds).Scan(&roleEntities); err != nil {
		return nil, gerror.Wrap(err, zconsts.ErrorORM)
	}
	for _, role := range roleEntities {
		roles = append(roles, &zweb.IdentityRole{
			Id:   role.Id,
			Name: role.Name,
			Key:  role.Key,
		})
	}
	return roles, nil
}

// BindUserContext 将用户信息绑定到上下文中
// 主要是重新查询数据库获取用户信息，因为用户信息可能发生变化
func (s *sCommonSite) BindUserContext(ctx g.Ctx, claims *zweb.Identity) (err error) {
	var mb *entity.AdminMember
	if err = dao.AdminMember.Ctx(ctx).Where(dao.AdminMember.Columns().Id, claims.Id).Scan(&mb); err != nil {
		return gerror.Wrap(err, zconsts.ErrorORM)
	}
	if mb == nil {
		return gerror.New("用户不存在")
	}
	if mb.Status != zconsts.StatusEnabled {
		return gerror.New("用户已禁用，如有疑问请联系管理员")
	}

	roles, err := s.getMemberRoles(ctx, mb.Id)
	if err != nil {
		return err
	}

	user := &zweb.Identity{
		Id:       mb.Id,
		Username: mb.Username,
		Roles:    roles,
		RealName: mb.RealName,
		Avatar:   mb.Avatar,
		Email:    mb.Email,
		Mobile:   mb.Mobile,
		App:      claims.App,
		LoginAt:  claims.LoginAt,
	}

	zcontext.SetUser(ctx, user)

	return nil
}

// GetRoutes 获取所有路由信息
func (s *sCommonSite) GetRoutes(ctx g.Ctx) (out *commonSchema.SiteRoutesOutput, err error) {
	routes := g.Server().GetRoutes()
	out = &commonSchema.SiteRoutesOutput{
		List: make([]*commonSchema.SiteRouteInfo, 0, len(routes)),
	}

	for _, route := range routes {
		out.List = append(out.List, &commonSchema.SiteRouteInfo{
			Method: route.Method,
			Route:  route.Route,
			Domain: route.Domain,
		})
	}

	return out, nil
}

// Logout 注销登录
func (s *sCommonSite) Logout(ctx g.Ctx) error {
	r := g.RequestFromCtx(ctx)
	if r == nil {
		return gerror.New("无法获取HTTP请求上下文")
	}
	return ztoken.Logout(r)
}
