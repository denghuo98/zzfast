# ZZFrame

基于 [GoFrame v2](https://goframe.org) 的 Go Web 后端框架，提供完整的管理后台：RBAC 权限控制（Casbin）、JWT 认证、消息队列、缓存和结构化日志。

## 项目结构

```
zzframe/
├── examples/basic/main.go    # 应用入口
├── config.example.yaml       # 配置文件模板（复制为 config.yaml 使用）
├── casbin.conf               # Casbin RBAC 模型配置
├── CLAUDE.md                 # Claude Code 开发指导
│
├── zapi/                     # API 请求/响应结构体（含 g.Meta 路由标签）
├── zcontroller/              # HTTP 控制器（薄层，调用 service）
├── zservice/                 # 服务接口定义 + 单例访问器
│   └── logic/                # 服务实现 + 中间件实现 + 测试
├── internal/                 # 私有代码
│   ├── dao/                  # 数据访问对象（GoFrame 生成）
│   └── model/                # entity/（ORM 结构体）+ do/（数据对象）
├── web/                      # 框架工具库（zresp/ztoken/zcache/zqueue/zcasbin 等）
├── zdb/                      # 数据库工具（迁移 + 通用查询）
├── zqueues/                  # 队列消费者主题注册
├── zconsts/                  # 应用常量和错误定义
├── zschema/                  # 跨层共享 DTO
├── zcmd/                     # CLI 命令定义
│
├── docs/                     # VitePress 文档站点
└── hack/                     # 初始化 SQL 脚本
```

数据流单向：`zapi → zcontroller → zservice → zservice/logic → internal/dao → DB`。

## 快速开始

### 环境要求

- Go 1.24+

### 安装与运行

```bash
git clone https://github.com/denghuo98/zfast.git
cd zzframe

# 创建配置文件
cp config.example.yaml config.yaml
# 编辑 config.yaml，至少修改 token.secretKey 和 superAdmin.password

# 安装依赖
go mod tidy

# 运行（默认启动 HTTP + 队列）
go run examples/basic/main.go

# 仅 HTTP 服务
go run examples/basic/main.go http

# 仅队列服务
go run examples/basic/main.go queue
```

### 配置

通过 `config.yaml` 配置（使用 `g.Cfg()` 读取），模板见 `config.example.yaml`。主要包括：

| 配置项 | 说明 |
|--------|------|
| system.token.secretKey | JWT 签名密钥（生产环境务必修改） |
| system.superAdmin.password | 超管初始密码 |
| system.cache.adapter | 缓存适配器（file/memory/redis） |
| system.upload.drive | 上传驱动（local/cos/oss/minio） |
| queue.driver | 队列驱动（当前仅 disk） |
| database.default.link | 数据库连接串（不配则用 SQLite） |

## 架构概览

- **分层架构**：严格分层单向流动，controller 薄层 → service 业务逻辑 → dao 数据访问
- **服务注册**：单例模式 + init() 自动注册
- **认证**：JWT（从 Authorization 头解析）
- **授权**：Casbin RBAC（策略从业务表加载到内存）
- **队列**：磁盘队列（disk），预留多后端扩展点
- **缓存**：支持 file/memory/redis 三种适配器
- **数据库**：启动自动建表（MySQL/SQLite），无配置时回退 SQLite

## 开发规范

详见 `.cursor/rules/project.mdc`、`docs/development-manual/development-standards.md` 和 `CLAUDE.md`。

## License

MIT
