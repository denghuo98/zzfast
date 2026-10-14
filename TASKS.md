# zzfast 后端优化任务清单

基于对 `backend/`（zzfast）逐模块 review 产出的缺点清单，全部优化项汇总如下。
每完成一项勾选 `[x]`。优先级：P0 阻塞/风险 > P1 质量 > P2 工程化。

---

## P0 — 阻塞 / 安全风险

- [x] **1. 重命名 Go module 路径**
  `go.mod` 的 `module github.com/denghuo98/zzframe` → 统一改为 zzfast 的新 module 路径，并全局替换所有 `import` 引用。

- [x] **2. 密码哈希 MD5 → bcrypt**
  `member.go`、`common/site.go`、`InitSuperAdmin` 等处 `gmd5.MustEncryptString(pwd+salt)` 改为 bcrypt（`golang.org/x/crypto/bcrypt`），保留/改造盐逻辑，迁移已有数据兼容方案。

- [x] **3. 登录密码加密 AES-ECB → AES-GCM（或 CBC+随机 IV）**
  `web/zencrypt` 的 ECB 改安全模式，前后端密钥/算法需同步；`verifyPassword` 解密逻辑同步调整。

- [x] **4. 事务 defer 提交错误处理**
  `member.Edit`、`role.Edit` 的 `defer func(){ Commit/Rollback }()` 改为显式处理 `Commit` 返回错误，避免静默吞错。

- [x] **5. 菜单删除后刷新 Casbin 策略**
  `menu.Delete` 删除 DB 后调用 `zcasbin.Refresh(ctx)`，与 `Edit` 行为一致，避免权限缓存不同步。

- [x] **6. 密钥/默认密码解耦 + 强校验**
  `zconsts.RequestEncryptKey`、`DefaultSuperAdminPassword`、`DefaultUserPassword` 与 `system.requestEncryptKey`、`system.superAdmin.password` 解耦；`secretKey`/`requestEncryptKey` 增加长度强度校验，缺失时报错而非静默用默认值。

- [x] **7. 修复 `RequestEncryptKey` Deprecated 但仍被使用的问题**
  `site.go:verifyPassword` 从配置读取密钥（`GetRequestEncryptConfig`），删除 `zconsts.RequestEncryptKey` 的默认值依赖。

---

## P1 — 代码质量 / 健壮性

- [x] **8. 角色/成员解绑空数组逻辑**
  `role.Edit` 与 `member.Edit` 中 `len(MenuIds/RoleIds)>0` 才更新关系，改为支持传空数组清空所有绑定（用 `== nil` 区分"未传"与"传空"）。

- [x] **9. 列表查询 N+1 优化**
  `member.embedRoles`、`role.embedMenus` 内层 `lo.Filter` O(n*m) 改为 Map 一次映射；`menu.treeList` 递归改为一次遍历构建索引树。

- [x] **10. `zcontext.GetData(ctx, key)` 签名修复**
  忽略的 `key` 形参改为按 key 取值（或调整命名/签名），名实相符。

- [x] **11. `zgorm.IsUnique` 表字段缓存**
  `GetPkField` 每次反射查 `TableFields`，改为按表名缓存主键字段，避免重复查询。

- [x] **12. 补齐日志/审计中间件**
  `zqueue` 已定义 `request_log`、`serve_log` 主题但无消费者；补访问日志中间件（记录请求耗时/响应）并接线消费者。

- [x] **13. 统一软删除/物理删除语义**
  `member.Delete` 软删除(status=2)，`role.Delete` 物理删除，统一为软删除或显式区分。

- [x] **14. 清理 `zservice` 机械重复注释**
  精简单例访问器/注册器的重复注释，保留关键信息。

---

## P2 — 工程化 / 可维护性

- [x] **15. 打破 `web/` 对 `internal/` 的依赖**
  `zstorager` 直接 import `internal/model/entity`，将通用模型下沉或改用 `zschema` 类型，保证 `web/` 可独立为 SDK。

- [x] **16. 配置字段命名对齐**
  `config.example.yaml` 的 `drive/fileSize` 等与 `webSchema.UploadConfig` 的 `uploadDrive/uploadFileSize` 命名不一致，统一 JSON tag 映射，去除宽松匹配的坑。

- [x] **17. 消除通用命名重名风险**
  `zapi/admin` 与 `zapi/common` 的 `DeleteReq/DeleteRes` 等通用名改为带领域前缀（如 `MemberDeleteReq`），降低误用。

- [x] **18. `z` 前缀包重命名（取消）**
  经评估决定**保留 `z` 前缀**：① 明确标识框架内部包，避免与标准库/第三方包名（db/api/cmd/service/consts 等）冲突；② 收益仅去掉一个字母；③ 全仓库 176+ 文件 import 替换风险高。

- [x] **19. 测试前缀规范化**
  `zz_*_test.go` 改为 `*_test.go`，评估对构建工具/CI 的影响。

---

## 集成测试发现并已修复的额外 Bug（不在原 review 清单内）

- [x] **A. Casbin Enforcer 误接 DB adapter 导致策略写回失败**
  `InitEnforcer` 用 `casbin.NewEnforcer(m, a)` 接上了 adapter，`Refresh`/`loadPermissions` 调 `AddPolicies` 时会把策略持久化到 `zz_admin_role_casbin` 表（id=0 触发 `UNIQUE constraint failed`）。已通过 `Enforcer.EnableAutoSave(false)` 关闭运行时持久化，策略仅加载到内存（`web/zcasbin/enforcer.go`）。

- [x] **B. 登录日志列表空结果导致空指针 panic**
  `SysLoginLog.List` 在 `totalCount==0` 时返回 `out=nil`，控制器 `*out` 解引用 panic。已改为始终返回非 nil 的 `out` 与空切片（`zservice/logic/system/login_log.go`）。

## 集成测试发现的遗留问题（未修复，需跟进）

- [x] **C. `exceptAuth` 白名单路径与真实路由不一致**
  实际路由带 `/admin` 前缀（如 `/admin/site/ping`、`/admin/menu/dynamic`），但 `config.example.yaml` 的 `exceptAuth` 写的是 `/site/ping`、`/menu/dynamic`。已修复：清空 `exceptAuth`——`/admin/site/*` 公共接口本就在 `AdminAuth` 中间件之前绑定、天然免鉴权；动态菜单 `/admin/menu/dynamic` 依赖登录态（按角色过滤菜单），不可匿名，否则会因 `GetUser` 为 nil 报错且泄露菜单结构。

- [x] **D. 成员列表接口泄露密码敏感字段**
  `MemberListOutputItem` 内嵌 `*entity.AdminMember`，`/admin/member/list` 直接返回 `passwordHash` 与 `salt`。已改为显式脱敏字段结构（含 `orm` 标签以保持 `Scan` 映射），不再暴露 `passwordHash`/`salt`/`passwordResetToken`。

- [ ] **E. `serverRoot: "public"` 目录缺失**
  无 `public/` 目录时服务启动即 `FATA`。已补 `public/.gitkeep`；建议后续接入前端静态资源或改为可选。

- [ ] **F. `/admin/site/routes` 匿名路由枚举泄露（新发现）**
  `Site.Routes` 在 `AdminAuth` 之前绑定，匿名可访问 `GET /admin/site/routes`，返回全部 28 条 API 路由（method + path），泄露系统攻击面。建议：移入鉴权组，或仅在 dev 模式可用，或按需裁剪返回内容。

---

## 说明

- P0 必须先行，涉及安全与正确性。
- P1 为质量项，建议紧随其后。
- P2 多为大范围重构，建议独立分次提交，避免与 P0/P1 混淆。
- 每项改动需同步跑 `go build ./...` 与 `go test ./...` 验证，并配合真实 HTTP 集成测试（见 `itest.py`）。
