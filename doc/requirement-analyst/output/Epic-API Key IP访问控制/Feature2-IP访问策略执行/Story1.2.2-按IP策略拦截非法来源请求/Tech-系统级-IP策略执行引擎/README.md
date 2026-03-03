# 【Tech-系统级】IP 策略执行引擎

## 【技术故事】

作为系统，需要在 API Key 鉴权流程中，根据 Key 绑定的 IP 策略（whitelist/blacklist + CIDR 列表）对请求来源 IP 进行匹配判断，对违规请求立即拒绝并返回标准化的 HTTP 403 响应（错误码 `IP_NOT_ALLOWED`），同时确保 CIDR 解析失败时降级不拦截，保证服务可用性优先。

**所属 Story**：Story 1.2.2 — 按 IP 策略拦截非法来源请求

**子系统职责**：
1. CIDR 解析与 IP 命中判断工具函数（`common/ip_matcher.go`）
2. 在鉴权中间件中读取 Token 的 `IpPolicy` 并调用工具函数执行白名单/黑名单校验
3. 违规时返回 HTTP 403 + `IP_NOT_ALLOWED`，并调用 `c.Abort()` 终止请求链

---

## 【验收条件】

### 一、功能性验收条件

**TC-1: 白名单 — 命中允许通过**
- **Given**: Token IpPolicy `{mode:"whitelist", ips:["1.2.3.4/32"]}`；客户端 IP `1.2.3.4`
- **When**: 中间件执行 IP 策略校验
- **Then**: 调用 `c.Next()`，请求继续处理，不产生拒绝响应

**TC-2: 白名单 — 未命中拒绝**
- **Given**: Token IpPolicy `{mode:"whitelist", ips:["1.2.3.4/32"]}`；客户端 IP `8.8.8.8`
- **When**: 中间件执行 IP 策略校验
- **Then**: HTTP 403；响应体 `{"success":false,"message":"IP not allowed","error_code":"IP_NOT_ALLOWED"}`；调用 `c.Abort()`

**TC-3: 黑名单 — 命中拒绝**
- **Given**: Token IpPolicy `{mode:"blacklist", ips:["5.5.5.0/24"]}`；客户端 IP `5.5.5.100`
- **When**: 中间件执行 IP 策略校验
- **Then**: HTTP 403 + `IP_NOT_ALLOWED`；`c.Abort()` 调用，不转发上游

**TC-4: 黑名单 — 未命中允许通过**
- **Given**: Token IpPolicy `{mode:"blacklist", ips:["5.5.5.0/24"]}`；客户端 IP `8.8.8.8`
- **When**: 中间件执行 IP 策略校验
- **Then**: 调用 `c.Next()`，请求正常处理

**TC-5: 无策略 — 跳过校验**
- **Given**: Token IpPolicy 为 nil 或 mode 为空字符串
- **When**: 中间件执行 IP 策略校验
- **Then**: 跳过校验，直接调用 `c.Next()`

**TC-6: CIDR 解析失败 — 降级不拦截**
- **Given**: Token `ip_policy` 包含无法解析的 CIDR（如数据库被篡改）
- **When**: 中间件调用 `ParseCIDRList` 时返回 error
- **Then**: 记录 WARN 日志（含 token_id、error）；跳过 IP 策略校验，调用 `c.Next()`（保证服务可用性）

**TC-7: 精确 IP 匹配（无 /32 前缀）**
- **Given**: Token 白名单 `["1.2.3.4"]`（无 `/32`）；客户端 IP `1.2.3.4`
- **When**: 执行 IP 策略校验
- **Then**: `ParseCIDRList` 自动补全为 `1.2.3.4/32`；匹配成功，请求通过

**TC-8: CIDR 边界精度校验**
- **Given**: Token 白名单 `["10.0.0.0/8"]`；分别测试 `10.255.255.255`（应命中）和 `11.0.0.1`（不应命中）
- **When**: 执行 CIDR 匹配
- **Then**: `10.255.255.255` 通过；`11.0.0.1` 拒绝

**TC-9: IP 策略在 Token 有效性校验之后执行**
- **Given**: Token 已过期；Token 有白名单策略；客户端 IP 在白名单内
- **When**: 请求到达鉴权中间件
- **Then**: 中间件先返回 Token 过期错误（HTTP 401），不进入 IP 策略校验块

**TC-10: `c.Abort()` 后不继续执行**
- **Given**: IP 策略校验失败，调用 `c.JSON(403, ...) + c.Abort() + return`
- **When**: 后续中间件链和 handler 逻辑
- **Then**: 均不再执行（relay 转发、计费、日志等后置 handler 均跳过）

### 二、非功能性验收条件

**TC-N1: 性能 — CIDR 匹配延迟**
- **Given**: 单 Token 含 100 条 IP 策略（最大配置）
- **When**: 高并发（QPS=1000）鉴权请求
- **Then**: IP 策略校验块增加的延迟 P99 < 1ms

**TC-N2: 并发安全 — CIDR 列表解析**
- **Given**: 多 goroutine 并发处理不同 Token 的鉴权请求（各自调用 `ParseCIDRList`）
- **When**: race detector 开启下并发测试
- **Then**: 无数据竞争，每次调用返回独立的 `[]*net.IPNet` 切片

**TC-N3: 日志记录**
- **Given**: IP 策略校验拒绝请求
- **When**: 中间件执行拒绝逻辑
- **Then**: 写入 INFO 日志，包含 `token_id`、`client_ip`、`mode`（便于审计）

---

## 【依赖与风险】

| 项 | 说明 |
|----|------|
| 依赖 Tech-系统级-可信代理感知IP提取 | 提供 `common.GetClientIP()` 和 `c.Get("client_ip")` |
| 依赖 `model/token.go` IpPolicy 字段 | Token 结构体中 `IpPolicy *IpPolicy` 字段 |
| 执行位置 | `middleware/auth.go` Token 有效性校验链末尾 |
| 风险：Q1 未确认 | `ips=[]` + `mode` 非空时，`ParseCIDRList([])` 返回空 slice，`IPMatchesCIDRList` 始终返回 false；白名单下 → 全部拒绝；黑名单下 → 全部通过。需与产品确认是否此行为符合预期 |
| 风险：Q6 未确认 | IP 拒绝事件是否写入数据库 logs（type=5）；不写库则仅依赖日志文件审计 |
| 风险：每次请求解析 CIDR | `ParseCIDRList` 每次调用均解析，高频 Key 有重复开销；若压测发现瓶颈，增加 Token 级缓存 |

---

## 【技术思路】

### 策略执行逻辑（中间件层）

```go
// new-api/middleware/auth.go — Token 有效性校验通过后追加（最后一个校验块）

// 从 Context 获取已提取的客户端 IP（由 Story 1.2.1 写入）
clientIP, _ := c.Get("client_ip")
clientIPStr, _ := clientIP.(string)

// IP 策略校验
policy := token.IpPolicy
if policy != nil && policy.Mode != "" {
    if clientIPStr == "" {
        clientIPStr = common.GetClientIP(c) // 降级：若 Context 未写入则重新提取
    }
    cidrs, err := common.ParseCIDRList(policy.Ips)
    if err != nil {
        // 降级：CIDR 解析失败，记录 WARN，不拦截
        logger.LogWarn(c.Request.Context(), fmt.Sprintf(
            "ip_policy_parse_failed: token_id=%d err=%v", token.Id, err))
    } else {
        hit := common.IPMatchesCIDRList(clientIPStr, cidrs)
        blocked := (policy.Mode == "whitelist" && !hit) ||
                   (policy.Mode == "blacklist" && hit)
        if blocked {
            logger.LogInfo(c.Request.Context(), fmt.Sprintf(
                "ip_not_allowed: token_id=%d client_ip=%s mode=%s",
                token.Id, clientIPStr, policy.Mode))
            c.JSON(http.StatusForbidden, gin.H{
                "success":    false,
                "message":    "IP not allowed",
                "error_code": "IP_NOT_ALLOWED",
            })
            c.Abort()
            return
        }
    }
}
```

### 鉴权链执行顺序

```
1. Token 存在性校验          → 不存在: 401/404
2. Token 状态校验（启用）     → 禁用: 403
3. Token 有效期校验           → 已过期: 401
4. Token 额度校验             → 额度不足: 429
5. Token 模型白名单校验       → 模型不允许: 403
6. ← Token IP 策略校验（本块） → IP 被拒: 403 + IP_NOT_ALLOWED
7. c.Next() → 转发上游/计费/日志
```

### 与旧 IP 白名单字段的关系

| 字段 | 执行时机 | 行为 |
|------|---------|------|
| 旧 `subnet` 字段（精确 IP 白名单） | 步骤 5 之前（现有逻辑） | 精确匹配，不支持 CIDR |
| 新 `IpPolicy` 字段 | 步骤 6（新增） | 支持 CIDR，白名单/黑名单双模式 |

两个字段均通过后，请求才能继续；旧字段优先执行。
