# DFX设计输入检查（DFX设计输入_check.md）

> 基于需求文档 doc/requirement/requirement.md 和代码仓库分析收集的 DFX 设计所需信息

## 1. 非功能需求清单

| 维度 | 具体要求 | 来源 |
|------|---------|------|
| **性能** | F001 P99 < 200ms；F004 GetClientIP P99 < 0.1ms；F005 IP策略校验 P99 < 1ms（n≤100条CIDR） | requirement.md §5.1 |
| **安全** | 防XFF伪造（仅可信代理时采信）；权限双重校验（拥有者或管理员）；CIDR入口校验；统一IP_NOT_ALLOWED错误码 | requirement.md §5.2 |
| **可靠性** | CIDR解析失败宽松降级（不拦截，记WARN日志）；请求隔离（无全局状态污染）；DB写失败返回500无中间状态 | requirement.md §5.3 |
| **可观测性** | INFO级拦截日志（token_id/client_ip/mode）；WARN级降级日志（token_id/error）；频繁拦截可配置告警 | requirement.md §5.4 |
| **可维护性** | 新建 ip_matcher_test.go；扩展 auth_test.go；并发安全测试（race detector） | requirement.md §5.5 |
| **可扩展性** | TRUSTED_PROXIES 环境变量配置；Token级CIDR缓存扩展点；条目数上限可配置 | requirement.md §5.6 |

## 2. 业务关键流程

| 流程 | 描述 | 关键节点 |
|------|------|---------|
| IP策略配置 | PUT /api/token/:id/ip_policy → 校验 → 持久化 | ValidateCIDRList → GetTokenByIds → UpdateIpPolicy |
| 请求鉴权执行 | API请求 → 提取ClientIP → 策略校验 → 放行/拦截 | GetClientIP → ParseCIDRList → IPMatchesCIDRList |

## 3. 外部依赖和集成点

| 依赖 | 类型 | 集成方式 |
|------|------|---------|
| Go标准库 `net` 包 | 工具库 | 直接调用（net.ParseCIDR、net.ParseIP、net.IPNet.Contains） |
| GORM AutoMigrate | ORM框架 | Token struct 添加字段后自动 ADD COLUMN |
| gin Framework | HTTP框架 | GetClientIP依赖 *gin.Context；IpPolicy块在TokenAuth中调用 |
| 环境变量 `TRUSTED_PROXIES` | 配置 | 服务启动时在 common/init.go:InitEnv() 读取一次 |
| Redis Token缓存 | 缓存 | Token整体序列化（含IpPolicy），与现有缓存兼容 |

## 4. 敏感数据和权限要求

| 数据/操作 | 敏感性 | 权限要求 |
|---------|-------|---------|
| Token.IpPolicy 策略配置 | 中（影响访问控制） | Token拥有者或管理员（双重校验） |
| tokens.ip_policy 字段 | 中（含IP白/黑名单信息） | 通过GORM ORM访问，不直接暴露SQL |
| TRUSTED_PROXIES 环境变量 | 高（影响IP提取安全性） | 部署级配置，运维人员管理 |
| 客户端真实IP | 中 | 仅用于内部策略比对，不持久化 |

## 5. 现有安全机制

| 机制 | 代码位置 | 说明 |
|------|---------|------|
| JWT身份认证 | middleware/auth.go:UserAuth | 所有管理接口均经过JWT校验 |
| 用户权限校验 | controller/token.go:225 (UpdateToken) | Token拥有者或管理员才可修改 |
| 旧AllowIps精确IP白名单 | middleware/auth.go:316-330 | 现有机制，不变 |
| GORM参数化查询 | model/*.go | ORM自动防SQL注入 |
| 统一错误码 | common/api_error.go | 防止错误信息泄露 |

## 6. 现有可靠性机制

| 机制 | 代码位置 | 说明 |
|------|---------|------|
| Token Redis缓存 | model/token_cache.go | Token数据缓存，减少DB读取 |
| GORM数据库连接池 | model/main.go | 连接池管理，断线重连 |
| Gin Recovery中间件 | router/ | panic自动恢复，防止进程崩溃 |
| 多数据库支持 | model/main.go | SQLite/MySQL/PostgreSQL三库兼容 |

## 7. 现有监控和日志实现

| 实现 | 代码位置 | 说明 |
|------|---------|------|
| 日志框架 | common/log.go（推测） | INFO/WARN/ERROR级别日志 |
| 请求日志中间件 | middleware/ | HTTP请求记录 |
| 无Prometheus/指标暴露 | - | 当前项目未内置指标暴露 |

## 8. 公用代码库和工具函数

| 组件 | 代码位置 | 本模块复用方式 |
|------|---------|--------------|
| `common/ip.go:IsIpInCIDRList` | common/ip.go:33 | 参考现有CIDR匹配模式，扩展实现 |
| `common/json.go` | common/json.go | JSON序列化/反序列化（Valuer/Scanner实现） |
| `common/init.go:InitEnv` | common/init.go:31 | 在末尾追加InitTrustedProxies调用 |
| `controller/token.go:UpdateToken` | controller/token.go:225 | 参考权限校验和响应格式 |
