# 【Story 1.1.1】配置 API Key 的 IP 访问策略

## 1、用户需求

**用户故事**：作为**平台用户（或管理员）**，我想要通过接口为指定 API Key 设置 IP 访问策略（白名单或黑名单模式 + IP/CIDR 列表），以便在 Key 泄露或遭受异常来源请求时，能精准限制该 Key 的可用 IP 范围。

**涉及用户角色**：

| 角色 | 操作权限 |
|------|---------|
| 平台管理员（is_admin=true） | 可操作任意用户的 Key，可设置白名单和黑名单 |
| 普通用户 | 只能操作自己 userId 下的 Key；是否可设置黑名单模式待产品确认（Q4），暂定仅限白名单 |

---

## 2、功能性需求

### 2.1 IP 策略配置接口

#### 2.1.1 功能描述

提供 `PUT /api/keys/{id}/ip_policy` REST 接口，供前端管理页面和 API 调用方配置 API Key 的 IP 访问策略。策略由两个字段组成：
- `mode`：策略模式，枚举值 `whitelist`（白名单） / `blacklist`（黑名单） / `""`（无策略/清空）
- `ips`：IP/CIDR 列表，支持精确 IP（如 `1.2.3.4`）和 CIDR 网段（如 `10.0.0.0/8`）

策略持久化到 `tokens` 表的 `ip_policy` 字段（TEXT 类型，JSON 序列化），配置后立即对该 Key 的后续请求生效。

#### 2.1.2 正常场景

**场景 1：管理员设置白名单策略**
- **Given**：请求方为管理员（is_admin=true），目标 Key id=1 存在
- **When**：`PUT /api/keys/1/ip_policy`，body `{"mode":"whitelist","ips":["1.2.3.4/32","10.0.0.0/8"]}`
- **Then**：HTTP 200，`{"success":true,"message":"IP 策略更新成功"}`；数据库 `tokens.ip_policy` 字段更新为对应 JSON

**场景 2：用户为自己的 Key 设置白名单策略**
- **Given**：请求方为普通用户（userId=100），目标 Key 属于该用户（owner userId=100）
- **When**：`PUT /api/keys/1/ip_policy`，body `{"mode":"whitelist","ips":["192.168.1.0/24"]}`
- **Then**：HTTP 200，策略更新成功

**场景 3：清空 IP 策略（恢复无限制）**
- **Given**：目标 Key 已设置过 IP 策略（mode=whitelist）
- **When**：`PUT /api/keys/1/ip_policy`，body `{"mode":"","ips":[]}`
- **Then**：HTTP 200，`tokens.ip_policy` 更新为 null 或空策略；该 Key 后续请求不再进行 IP 校验

**场景 4：设置黑名单策略（管理员）**
- **Given**：请求方为管理员，目标 Key 存在
- **When**：`PUT /api/keys/1/ip_policy`，body `{"mode":"blacklist","ips":["5.5.5.0/24","6.6.6.6"]}`
- **Then**：HTTP 200，策略更新成功

#### 2.1.3 异常场景

**场景 5：CIDR 格式错误**
- **Given**：IP 列表中包含非法格式
- **When**：body `{"mode":"whitelist","ips":["999.0.0.0/8","not-an-ip"]}`
- **Then**：HTTP 400，错误信息明确指出非法条目（如 `"invalid IP/CIDR: 999.0.0.0/8"`）；策略不被保存

**场景 6：mode 取值非法**
- **Given**：body 中 mode 为非枚举值
- **When**：body `{"mode":"allow","ips":[]}`
- **Then**：HTTP 400，错误信息说明 mode 只允许 `whitelist`/`blacklist`/空值

**场景 7：目标 Key 不存在**
- **Given**：请求的 id 在数据库中不存在
- **When**：`PUT /api/keys/99999/ip_policy`
- **Then**：HTTP 404

**场景 8：ips 为空列表但 mode 不为空（待确认 Q1）**
- **Given**：body `{"mode":"whitelist","ips":[]}`（mode 非空但 ips 为空列表）
- **When**：发起请求
- **Then**：【暂定】HTTP 200，等同于清空策略（mode 被忽略），策略不生效；待与产品确认最终行为

#### 2.1.4 安全场景

**场景 9：普通用户操作他人 Key**
- **Given**：请求方为普通用户（userId=100），目标 Key 属于其他用户（owner userId=200）
- **When**：`PUT /api/keys/5/ip_policy`
- **Then**：HTTP 403，`{"success":false,"message":"forbidden"}`；策略不被修改

**场景 10：未登录用户调用接口**
- **Given**：请求未携带有效 Bearer Token
- **When**：`PUT /api/keys/1/ip_policy`
- **Then**：HTTP 401（由现有鉴权中间件返回，不进入 Handler）

**场景 11：批量提交大量 CIDR 条目（防滥用）**
- **Given**：请求 body 中 ips 包含超过上限条目数（上限待确认 Q2，暂定 100）
- **When**：`PUT /api/keys/1/ip_policy`，`ips` 包含 101 条
- **Then**：HTTP 400，错误信息说明超出最大条目数限制

#### 2.1.5 边界场景

**场景 12：单条精确 IP（无 CIDR 前缀）**
- **Given**：body `{"mode":"whitelist","ips":["1.2.3.4"]}`（无 /32）
- **When**：发起请求
- **Then**：HTTP 200，系统自动将 `1.2.3.4` 视为 `1.2.3.4/32` 处理，等价于精确匹配

**场景 13：ips 和 mode 均为空**
- **Given**：body `{"mode":"","ips":null}` 或 `{}`
- **When**：发起请求
- **Then**：HTTP 200，等同于清空策略（无需校验）

**场景 14：IPv6 地址格式**
- **Given**：body 中包含 IPv6 CIDR，如 `"::1/128"`
- **When**：发起请求
- **Then**：【当前版本】若 IPv6 不在支持范围内，返回 HTTP 400 并注明当前版本仅支持 IPv4

#### 2.1.6 约束条件

| 约束 | 规格 |
|------|------|
| 接口权限 | 管理员：可操作任意 Key；普通用户：只能操作自己的 Key |
| mode 枚举值 | 仅限 `whitelist`、`blacklist`、`""`（空字符串） |
| ips 每条格式 | 精确 IP（如 `1.2.3.4`）或 CIDR（如 `1.2.3.0/24`），当前版本仅支持 IPv4 |
| ips 最大条目数 | 100 条（待确认 Q2，暂定） |
| 存储格式 | `ip_policy` 字段为 TEXT 类型，JSON 序列化；兼容 SQLite/MySQL/PostgreSQL |
| 策略生效时延 | 依赖 Token 缓存 TTL（通常 60s），非实时生效；可配置更低 TTL 优化 |

#### 2.1.7 依赖与风险

| 项 | 说明 |
|----|------|
| 依赖 Feature 2 | IP 策略执行由 Feature 2（鉴权中间件）负责读取并执行，本 Feature 仅负责配置和持久化 |
| 旧 IP 白名单字段共存 | `tokens` 表原有 `subnet` 字段（精确 IP 白名单）与新 `ip_policy` 并存；执行优先级在 Story 1.2.2 中明确 |
| Q1 待确认 | `ips=[]` + `mode` 非空时语义（当前暂定等同清空） |
| Q4 待确认 | 普通用户是否可设置黑名单模式（当前暂定不允许） |

---

## 3、非功能性需求

### 3.1 性能分析

| 指标 | 目标值 | 说明 |
|------|--------|------|
| 接口响应时间（P99） | < 200ms | 含数据库写入，策略配置为低频操作 |
| 接口响应时间（P50） | < 50ms | 正常场景目标 |
| 并发支持 | 与现有 Token 配置接口一致，不单独加压 | 策略配置属于管理面，非高频场景 |

### 3.2 可靠性

- 数据库写入失败时返回 HTTP 500，入参校验失败不写库（写前校验）
- `AutoMigrate` 新增 `ip_policy` 列为 NULL 类型，不影响存量 Token 行，存量 Token 默认无策略（不拦截）
- 三数据库（SQLite/MySQL/PostgreSQL）均需通过自动化 CI 回归测试

### 3.3 可维护性

- 策略更新成功时写入操作日志（type=3 管理操作日志），记录操作方 userId、目标 Token Id、新策略内容
- 接口错误记录到日志，关键字 `update_ip_policy_failed`，便于排障
- 策略字段迁移由 GORM `AutoMigrate` 自动完成，无需手动 DDL

### 3.4 开放性

- `PUT /api/keys/{id}/ip_policy` 接口需纳入 OpenAPI 文档（如有）
- 响应格式与现有接口保持一致：`{"success":true/false,"message":"..."}`

### 3.5 安全性

| 安全需求 | 措施 |
|---------|------|
| 接口权限控制 | 校验请求方身份（管理员 or Key 归属用户），防止越权 |
| 输入校验 | CIDR 格式在写入前严格校验，防止非法数据入库 |
| 错误信息安全 | 错误信息不暴露内部数据库结构，仅返回"非法条目"描述 |
| 日志安全 | 操作日志记录 token_id 和操作者，不记录完整 IP 列表（防日志过大） |

### 3.6 可测试性

- `ValidateCIDRList` 函数有独立单元测试（`common/ip_matcher_test.go`）
- Handler 层有接口测试覆盖正常场景（白名单/黑名单设置、清空）和异常场景（格式错误/权限不足/Key不存在）
- 数据库迁移测试在三种数据库 CI 环境中运行
