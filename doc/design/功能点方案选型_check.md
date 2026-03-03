# 功能点方案选型检查（功能点方案选型_check.md）

> 基于需求文档分析，识别需要进行方案选型的功能点

## 功能点方案选型一览

| 功能点 | 功能描述 | 是否进行方案选型 | 选型维度 |
|-------|---------|----------------|---------|
| F001 IP策略配置接口 | PUT /api/token/:id/ip_policy 接口实现 | 否 | 已明确：REST接口，与现有接口风格一致 |
| F002 Token.IpPolicy数据模型扩展 | ip_policy TEXT字段+JSON序列化 | **是** | 序列化方案选型（JSON vs 单独表 vs 扩展字段） |
| F003 CIDR工具函数库 | ValidateCIDRList / ParseCIDRList / IPMatchesCIDRList | **是** | CIDR匹配实现方案选型（标准库 vs 外部库 vs 缓存策略） |
| F004 可信代理感知IP提取 | GetClientIP + InitTrustedProxies | **是** | 可信代理配置方式选型（环境变量 vs 配置文件 vs DB） |
| F005 IP策略执行（鉴权中间件扩展） | TokenAuth扩展：白/黑名单校验逻辑 | 否 | 已明确：在现有AllowIps校验后插入新块 |

## 需进行方案选型的功能点

### F002 数据序列化方案选型（待4.2章节分析）

**选型决策点**：IpPolicy 数据（mode + ips 列表）如何持久化到数据库？

候选方案：
- **方案A**：JSON 序列化到单一 TEXT 列（`tokens.ip_policy`）
- **方案B**：新建独立表 `token_ip_policies`（id, token_id, mode, ips JSON）
- **方案C**：将 mode 和 ips 分别存储为两个列（`ip_policy_mode TEXT`, `ip_policy_ips TEXT`）

### F003 CIDR匹配实现选型（待4.2章节分析）

**选型决策点**：每次请求的 CIDR 匹配如何实现？

候选方案：
- **方案A**：每次请求调用 `ParseCIDRList` 即时解析
- **方案B**：Token 级内存缓存（TTL 60s）
- **方案C**：`sync.Map` 全局缓存（key=token.IpPolicy JSON hash）

### F004 可信代理配置方式选型（待4.2章节分析）

**选型决策点**：`TRUSTED_PROXIES` 如何配置和初始化？

候选方案：
- **方案A**：环境变量（`TRUSTED_PROXIES`）+ 服务启动时初始化一次
- **方案B**：数据库 `options` 表存储，支持运行时动态变更
- **方案C**：配置文件（YAML/TOML）静态配置

## 现有架构组件清单

| 组件 | 类型 | 现有实现位置 | 备注 |
|-----|------|-------------|------|
| GORM ORM | ORM框架 | `model/main.go:254` | AutoMigrate 已包含 &Token{} |
| Gin Web框架 | HTTP框架 | `router/api-router.go` | 所有接口基于 Gin |
| Redis | 缓存 | `model/token_cache.go` | Token 对象已有 Redis 缓存 |
| MySQL/PostgreSQL/SQLite | 数据库 | `model/main.go` | 三库兼容是强制要求 |
| 现有 IsIpInCIDRList | CIDR工具 | `common/ip.go:33` | 参考现有实现，新函数风格保持一致 |

## 现有技术栈清单

| 层级 | 技术 | 版本 |
|-----|-----|-----|
| 后端语言 | Go | 1.22+ |
| HTTP 框架 | Gin | latest |
| ORM | GORM v2 | latest |
| JSON | encoding/json（通过 common/json.go 封装） | 标准库 |
| IP/CIDR | net 包 | 标准库 |
