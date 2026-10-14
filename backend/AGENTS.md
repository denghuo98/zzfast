# CLAUDE.md

本文件为 Claude Code（claude.ai/code）在处理此仓库代码时提供指导。

## 项目概述

ZZFast 是一个基于 [GoFrame v2](https://goframe.org) 构建的 Go Web 后端框架。它提供了完整的管理后台，包含 RBAC 权限控制（Casbin）、JWT 认证、消息队列、缓存和结构化日志等功能。

## 常用命令

- **运行所有服务**：`go run examples/basic/main.go`
- **仅运行 HTTP 服务**：`go run examples/basic/main.go http`
- **仅运行队列服务**：`go run examples/basic/main.go queue`
- **热重载（需要 gf CLI）**：`gf run examples/basic/main.go`
- **运行所有测试**：`go test ./...`
- **运行指定包的测试**：`go test ./zservice/logic/admin/...`
- **运行单个测试**：`go test -run TestAdminMember_Edit ./zservice/logic/admin/...`
- **文档开发服务器**：`cd docs && pnpm dev`
- **文档构建**：`cd docs && pnpm build`

## 架构

### 分层结构

代码库遵循严格的分层架构。数据仅沿单一方向流动：

```
zapi/          -> API 请求/响应结构体定义（含 g.Meta 路由标签）
zcontroller/   -> HTTP 处理器；薄层，负责调用服务并格式化响应
zservice/      -> 服务接口（IAdminMember、IAdminRole 等）和单例访问器
zservice/logic/-> 服务实现；包含所有业务逻辑
internal/dao/  -> 数据访问对象（GoFrame 生成）；包装 `internal/dao/internal/`
internal/model/-> `entity/`（ORM 结构体）和 `do/`（数据对象）
web/           -> 框架工具：zresp、zcontext、ztoken、zcache、zqueue、zcasbin 等
zdb/           -> 数据库工具和自动迁移
zqueues/       -> 队列消费者主题注册
zconsts/       -> 应用常量和错误定义
zschema/       -> 在 service 和 controller 层之间共享的 DTO
zcmd/          -> CLI 命令定义（http、queue、all）
```

### 服务注册模式

服务采用单例注册模式。每个服务实现：

1. 在 `zservice/` 中定义接口（如 `IAdminMember`）
2. 在 `zservice/logic/<domain>/` 中实现（如 `sAdminMember`）
3. 在 `init()` 函数中注册自身：

```go
func init() {
    zservice.RegisterAdminMember(NewAdminMember())
}
```

在其他地方通过 `zservice.AdminMember().List(ctx, ...)` 访问服务。如果单例未注册，会触发 panic。

### 请求路由

控制器是基于结构体的。GoFrame 使用结构体标签和 `g.Meta` 来定义路由：

```go
type MemberListReq struct {
    g.Meta `path:"/member/list" method:"get" tags:"SYS-02-用户管理" summary:"获取用户列表"`
    adminSchema.MemberListInput
}
```

控制器在 `zcontroller/admin.go` 中绑定到路由组，并按组应用中间件。

### 中间件栈

全局中间件（在 `zcmd/http.go` 中应用于所有路由）：
- `zservice.Middleware().Ctx` — 初始化请求上下文
- `zservice.Middleware().CORS` — CORS 处理
- `zservice.Middleware().ResponseHandler` — 标准化 JSON 响应

认证中间件（在 `zcontroller/admin.go` 中应用于 `/admin/*` 路由）：
- `zservice.Middleware().AdminAuth` — JWT 验证 + Casbin 权限检查

### 响应标准化

所有 JSON 响应通过 `web/zresp/response.go` 处理。控制器返回 `(res *XxxRes, err error)`。`ResponseHandler` 中间件拦截处理器响应，将错误包装成标准的 `{code, message, data/error, timestamp, traceId}` JSON 结构，并通过 `zresp.RJson()` 写入响应。

### 数据库与迁移

`zdb/zmigrate` 包在启动时自动运行（`init()`）。它会：
- 从 `config.yaml` 读取数据库配置
- 如果没有配置，回退到本地 SQLite 文件（`zzfast.db`）
- 如果表不存在则创建表，使用针对不同数据库类型（MySQL vs SQLite）的 SQL

`internal/dao/` 中的 DAO 是对 GoFrame 生成的 `internal/dao/internal/` 结构体的薄包装。建议使用 `dao.Xxx.Ctx(ctx).Where(...).Scan(&dest)` 进行查询。

### 认证与授权

- **认证**：从 `Authorization` 请求头解析 JWT 令牌。令牌配置（密钥、过期时间）位于 `config.yaml` 的 `system.token` 下。
- **授权**：Casbin RBAC。`web/zcasbin` 包从 `casbin.conf` 和数据库加载策略。角色-权限检查在 `AdminAuth` 中间件中执行。
- **超级管理员**：硬编码的用户名（`superAdmin`）和默认角色标识（`super`），可绕过某些限制。

### 队列系统

`web/zqueue/` 中的自定义队列抽象当前实现了磁盘队列（disk），架构上预留了 Redis、RocketMQ、Kafka 等多后端扩展点。生产者和消费者按 `groupName` 实例化。消费者主题在 `zqueues/` 包中注册。运行 `queue` 或 `all` 命令时，队列作为独立 goroutine 启动。

### 测试约定

测试位于 `zservice/logic/<domain>/zz_*_test.go`，使用内存 SQLite。每个测试：
- 通过 `gdb.SetConfig()` 配置，使用 `sqlite::@file(:memory:)?cache=shared`
- 通过原始 SQL 手动创建所需的表
- 插入必要的参考数据（如在测试成员之前先创建角色）
- 在 `defer` 中清理表

测试使用 `github.com/gogf/gf/v2/test/gtest` 断言库。

## 代码规范

来自 `.cursor/rules/project.mdc` 和 `docs/development-manual/development-standards.md`：

- **包命名**：小写，不含下划线或大小写混合（如 `admin`，而非 `admin_service`）
- **接口命名**：以 `I` 为前缀（如 `IAdminMember`）
- **错误处理**：使用 `gerror.Wrap(err, zconsts.ErrorORM)` 或 `gerror.New("message")` 包装业务错误。不要返回没有上下文的原始错误。
- **注释**：导出的函数/类型必须有以其名称开头的注释。
- **Git 提交**：`<type>(<scope>): <subject>`，其中 type 为 `feat`、`fix`、`docs`、`style`、`refactor`、`test`、`chore` 之一。

## 重要文件

- `examples/basic/main.go` — 应用入口点
- `config.yaml` — 运行时配置（数据库、缓存、令牌、队列、服务器端口）
- `casbin.conf` — Casbin RBAC 模型配置
- `zcmd/cmd.go` — CLI 命令树（all、http、queue、help）
- `zcmd/http.go` — HTTP 服务器设置和中间件绑定
- `zcontroller/admin.go` — Admin 路由组注册
- `zservice/logic/logic.go` — 空白导入，触发所有服务 `init()` 注册
