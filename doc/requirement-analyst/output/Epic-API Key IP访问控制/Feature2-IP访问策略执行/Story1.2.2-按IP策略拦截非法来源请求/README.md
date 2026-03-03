# 【Story 1.2.2】按 IP 策略拦截非法来源请求

## 1、用户需求

**用户故事**：作为**平台网关**，我需要在每次 API Key 鉴权时，根据该 Key 绑定的 IP 策略（白名单或黑名单）对请求来源 IP 进行校验，对不符合策略的请求立即拒绝并返回标准化错误，确保 Key 的访问范围严格受控，防止 Key 被从非授权来源使用。

**依赖关系**：本 Story 依赖 Story 1.2.1 提供的 `GetClientIP()` 真实客户端 IP 提取能力，在其基础上执行策略匹配和拦截。

**涉及用户角色**：

| 角色 | 影响 |
|------|------|
| API Key 使用者 | 来源 IP 须符合 Key 的 IP 策略，否则收到 HTTP 403 |
| 平台管理员 | 通过 IP 策略（Feature 1）控制每个 Key 的可用范围 |
| 平台运营者 | 事故中通过黑名单快速封堵攻击来源 |

---

## 2、功能性需求

### 2.1 IP 策略校验（白名单 / 黑名单）

#### 2.1.1 功能描述

在 `middleware/auth.go` Token 有效性校验（存在、启用、有效期、额度、模型白名单）通过后，读取 Token 的 `IpPolicy` 字段，结合 `GetClientIP()` 获取的真实客户端 IP，执行以下逻辑：

```
若 IpPolicy == nil 或 mode == ""：
    跳过校验，继续处理请求

若 mode == "whitelist"：
    客户端 IP 命中 CIDR 列表 → 通过
    客户端 IP 未命中 CIDR 列表 → 拒绝（HTTP 403 + IP_NOT_ALLOWED）

若 mode == "blacklist"：
    客户端 IP 命中 CIDR 列表 → 拒绝（HTTP 403 + IP_NOT_ALLOWED）
    客户端 IP 未命中 CIDR 列表 → 通过
```

#### 2.1.2 正常场景

**场景 1：白名单模式 — 来源 IP 在白名单内，请求通过**
- **Given**: Token IpPolicy `{mode:"whitelist", ips:["1.2.3.4/32"]}`；客户端 IP `1.2.3.4`
- **When**: 使用该 Token 调用任意业务接口
- **Then**: IP 校验通过，请求继续处理，返回正常业务响应

**场景 2：白名单模式 — 来源 IP 在 CIDR 网段内，请求通过**
- **Given**: Token IpPolicy `{mode:"whitelist", ips:["10.0.0.0/8"]}`；客户端 IP `10.255.255.1`
- **When**: 调用业务接口
- **Then**: IP 校验通过（`10.255.255.1` 属于 `10.0.0.0/8`）

**场景 3：黑名单模式 — 来源 IP 不在黑名单内，请求通过**
- **Given**: Token IpPolicy `{mode:"blacklist", ips:["5.5.5.0/24"]}`；客户端 IP `8.8.8.8`
- **When**: 调用业务接口
- **Then**: IP 校验通过（`8.8.8.8` 不在 `5.5.5.0/24` 内），请求继续处理

**场景 4：无 IP 策略 — 跳过校验**
- **Given**: Token IpPolicy 为 nil（未配置）
- **When**: 调用业务接口
- **Then**: 跳过 IP 策略校验，请求正常处理（与无此功能时行为一致）

#### 2.1.3 异常场景

**场景 5：白名单模式 — 来源 IP 不在白名单内，请求被拒绝**
- **Given**: Token IpPolicy `{mode:"whitelist", ips:["1.2.3.4/32"]}`；客户端 IP `8.8.8.8`
- **When**: 使用该 Token 调用业务接口
- **Then**: HTTP 403，响应体：`{"success":false,"message":"IP not allowed","error_code":"IP_NOT_ALLOWED"}`；请求不被转发到上游

**场景 6：黑名单模式 — 来源 IP 命中黑名单，请求被拒绝**
- **Given**: Token IpPolicy `{mode:"blacklist", ips:["5.5.5.0/24"]}`；客户端 IP `5.5.5.100`
- **When**: 使用该 Token 调用业务接口
- **Then**: HTTP 403 + `IP_NOT_ALLOWED`；请求不被转发到上游

**场景 7：CIDR 解析失败时降级不拦截**
- **Given**: Token `ip_policy` 字段存储了非法 CIDR（如数据库被手工篡改为 `{"mode":"whitelist","ips":["bad"]}`）
- **When**: 调用业务接口，中间件执行 `ParseCIDRList`
- **Then**: CIDR 解析失败，记录 WARN 日志；**降级处理为不拦截**，请求继续处理（保证服务可用性优先）

#### 2.1.4 安全场景

**场景 8：白名单配置 + XFF 伪造攻击**
- **Given**: Token 白名单 `["1.1.1.1/32"]`；攻击者 IP `8.8.8.8`，伪造 `X-Forwarded-For: 1.1.1.1`；`TRUSTED_PROXIES` 未配置
- **When**: 请求到达鉴权中间件
- **Then**: `GetClientIP()` 返回 `"8.8.8.8"`（XFF 被忽略）；IP 校验失败，返回 HTTP 403；攻击无效

**场景 9：IP 策略在 Token 有效性校验之后执行**
- **Given**: Token 已过期（或额度为 0）；Token 有白名单策略 `["1.2.3.4/32"]`；客户端 IP `1.2.3.4`
- **When**: 请求到达鉴权中间件
- **Then**: 先返回"Token 已过期"错误（HTTP 401/403），不执行 IP 策略校验；IP 策略仅对有效 Token 生效

**场景 10：被禁用的 Token 不执行 IP 策略**
- **Given**: Token 状态为 disabled；Token 有白名单策略
- **When**: 请求到达鉴权中间件
- **Then**: 返回 Token 禁用错误，不执行 IP 策略校验

#### 2.1.5 边界场景

**场景 11：IP 策略 ips 为空列表**
- **Given**: Token IpPolicy `{mode:"whitelist", ips:[]}`（mode 非空但 ips 为空，待确认 Q1）
- **When**: 调用业务接口
- **Then**：【暂定】等同于无策略，跳过校验（空白名单不拦截任何 IP）；待与产品最终确认

**场景 12：多个 CIDR 条目命中边界**
- **Given**: Token 白名单 `["10.0.0.0/8","172.16.0.0/12"]`；客户端 IP `172.31.255.255`（边界 IP）
- **When**: 调用业务接口
- **Then**: `172.31.255.255` 属于 `172.16.0.0/12`（`172.16.0.0`–`172.31.255.255`），校验通过

**场景 13：精确 IP 策略（无 CIDR 前缀）**
- **Given**: Token 白名单 `["1.2.3.4"]`（无 `/32` 前缀）
- **When**: 客户端 IP `1.2.3.4` 调用接口
- **Then**: 系统自动将 `1.2.3.4` 视为 `1.2.3.4/32`，精确匹配通过

**场景 14：同一请求的 relay 路由（转发请求）**
- **Given**: Token 有白名单策略；请求路径为 `/v1/chat/completions`（relay 路由）
- **When**: 请求通过 IP 校验
- **Then**: 请求正常转发到上游模型接口，无额外 IP 校验

#### 2.1.6 约束条件

| 约束 | 规格 |
|------|------|
| 校验执行时机 | 在 Token 存在性/状态/有效期/额度/模型白名单校验**之后**执行 |
| 响应格式 | `{"success":false,"message":"IP not allowed","error_code":"IP_NOT_ALLOWED"}`，与现有鉴权错误格式一致 |
| HTTP 状态码 | 403（无论白名单/黑名单，统一使用 403） |
| 降级策略 | CIDR 解析失败时不拦截（宽松降级），记录 WARN 日志 |
| 无策略时行为 | `IpPolicy == nil` 或 `mode == ""` 时，跳过校验（向后兼容） |
| 与旧 IP 白名单字段关系 | 旧 `subnet` 字段（精确 IP 白名单）执行顺序在新 `IpPolicy` 之前；两者均通过后请求才能继续 |
| 错误码 | 固定为 `IP_NOT_ALLOWED`，防止泄露策略模式（白名单/黑名单）信息 |

#### 2.1.7 依赖与风险

| 项 | 说明 |
|----|------|
| 依赖 Story 1.2.1 | `common.GetClientIP(c)` 提供真实客户端 IP |
| 依赖 `common/ip_matcher.go` | `ParseCIDRList`、`IPMatchesCIDRList` 工具函数 |
| 依赖 Feature 1 (Story 1.1.1) | `tokens.ip_policy` 字段已持久化 |
| 风险：Q1 待确认 | `ips=[]` + `mode` 非空时行为（场景 11） |
| 风险：Q6 待确认 | IP 拒绝事件是否写入数据库 logs（type=5）便于后台审计 |

---

## 3、非功能性需求

### 3.1 性能分析

| 指标 | 目标值 | 说明 |
|------|--------|------|
| IP 策略校验增加的鉴权延迟（P99） | < 1ms | CIDR 匹配为纯内存操作，O(n)，n ≤ 100 |
| 高并发下内存稳定性 | 无内存泄漏 | `ParseCIDRList` 每次调用分配临时内存，GC 管理 |
| CIDR 解析频次 | 每次请求均调用 `ParseCIDRList` | 如性能压测发现瓶颈，可增加 Token 级 CIDR 缓存（TTL=60s） |

### 3.2 可靠性

- CIDR 解析失败时降级不拦截，保证服务可用性；降级须记录 WARN 日志
- 鉴权链执行 `c.Abort()` 后必须立即 `return`，防止后续 handler 仍被执行
- IP 策略校验失败不影响其他 Token 的请求处理（无全局状态污染）

### 3.3 可维护性

- IP 拦截事件记录日志（INFO），包含 token_id、client_ip、mode，便于运维审计
- 可选：IP 拒绝事件写入数据库 logs（type=5），供管理后台查询（待确认 Q6）
- 错误响应格式与现有鉴权错误保持一致，调用方无需新增错误处理逻辑

### 3.4 安全性

| 安全需求 | 措施 |
|---------|------|
| 防止策略信息泄露 | 错误码固定为 `IP_NOT_ALLOWED`，不区分白名单/黑名单被拒，防止攻击者推断策略类型 |
| 防止降级被滥用 | CIDR 解析失败只在入库时（Feature 1）产生，运行时降级为不拦截；配合 WARN 日志监控 |
| 鉴权链完整性 | IP 策略校验在 Token 有效性校验之后，无效 Token 不进入 IP 校验逻辑 |

### 3.5 可测试性

- 单元测试：`IPMatchesCIDRList` 覆盖白名单命中/未命中、黑名单命中/未命中、边界 IP、空列表等场景
- 集成测试：端到端覆盖白名单拒绝（HTTP 403 + 响应体格式）、黑名单拒绝、无策略通过
- 安全测试：XFF 伪造攻击在白名单场景下被正确拒绝
- 测试文件：`new-api/common/ip_matcher_test.go`（工具函数）+ `new-api/middleware/auth_test.go`（中间件集成）
