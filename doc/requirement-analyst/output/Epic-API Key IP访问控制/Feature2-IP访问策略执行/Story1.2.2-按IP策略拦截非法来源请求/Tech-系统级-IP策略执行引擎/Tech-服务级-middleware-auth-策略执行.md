# 【Tech-服务级】middleware-auth（IP 策略执行部分）

## 服务职责

在 `new-api/middleware/auth.go` 的 Token 鉴权流程中，读取 Token 的 `IpPolicy` 字段，结合 `common.GetClientIP()` 提供的真实客户端 IP，调用 CIDR 匹配工具执行白名单/黑名单校验。对违规请求立即返回 HTTP 403 + `IP_NOT_ALLOWED`，并调用 `c.Abort()` 终止后续请求处理链。

**所属 Tech-系统级**：Tech-系统级-IP 策略执行引擎

---

## 详细验收条件

### AC-1: 白名单命中 — 请求通过

- **Given**: Token IpPolicy `{mode:"whitelist", ips:["1.2.3.4/32"]}`；客户端 IP `1.2.3.4`
- **When**: 中间件执行 IP 策略校验块
- **Then**: 不产生拒绝响应；调用 `c.Next()` 继续处理请求

### AC-2: 白名单未命中 — 返回 403

- **Given**: Token IpPolicy `{mode:"whitelist", ips:["1.2.3.4/32"]}`；客户端 IP `8.8.8.8`
- **When**: 中间件执行 IP 策略校验块
- **Then**:
  - HTTP 状态码 403
  - 响应体: `{"success":false,"message":"IP not allowed","error_code":"IP_NOT_ALLOWED"}`
  - 调用 `c.Abort()` 后立即 `return`，不继续执行后续中间件和 handler

### AC-3: 黑名单命中 — 返回 403

- **Given**: Token IpPolicy `{mode:"blacklist", ips:["5.5.5.0/24"]}`；客户端 IP `5.5.5.100`
- **When**: 中间件执行 IP 策略校验块
- **Then**: HTTP 403 + `IP_NOT_ALLOWED`；`c.Abort()` + `return`

### AC-4: 黑名单未命中 — 请求通过

- **Given**: Token IpPolicy `{mode:"blacklist", ips:["5.5.5.0/24"]}`；客户端 IP `8.8.8.8`
- **When**: 中间件执行 IP 策略校验块
- **Then**: 不产生拒绝响应；调用 `c.Next()`

### AC-5: 无策略 — 跳过校验

- **Given**: Token IpPolicy 为 nil 或 mode 为空字符串
- **When**: 中间件执行 IP 策略校验块
- **Then**: 跳过校验，直接继续后续逻辑

### AC-6: IP 策略在 Token 有效性校验之后执行

- **Given**: Token 已过期（ValidBefore < now）；Token 有白名单策略；客户端 IP 在白名单内
- **When**: 请求到达鉴权中间件
- **Then**: 鉴权链在"有效期校验"步骤返回错误（如 HTTP 401），不进入"IP 策略校验"块

### AC-7: CIDR 解析失败 — 降级不拦截

- **Given**: Token `ip_policy` 包含非法 CIDR（如数据库被手工修改）
- **When**: 中间件调用 `common.ParseCIDRList` 返回 error
- **Then**:
  - 记录 WARN 日志，包含 token_id 和错误信息
  - **不拦截请求**，调用 `c.Next()` 继续处理（宽松降级，保证服务可用性）

### AC-8: IP 拦截事件日志记录

- **Given**: IP 策略校验失败，请求被拦截
- **When**: 中间件执行拒绝逻辑
- **Then**: 写入 INFO 日志，包含 token_id、client_ip、mode（便于审计和告警）

### AC-9: `c.Abort()` 后不继续执行后续链

- **Given**: 白名单未命中，调用 `c.JSON(403, ...) + c.Abort() + return`
- **When**: 检查后续中间件和 handler 是否被执行（如 relay 转发、计费 handler）
- **Then**: 均不再执行（gin 的 `Abort()` 机制保证）

### AC-10: 响应格式与现有鉴权错误一致

- **Given**: IP 策略校验失败
- **When**: 检查响应体
- **Then**: 格式为 `{"success":false,"message":"...","error_code":"IP_NOT_ALLOWED"}`，与 Token 过期、额度不足等鉴权错误响应格式相同

---

## 技术实现

### 代码位置

`new-api/middleware/auth.go` — 在 Token 鉴权函数的有效性校验链末尾追加 IP 策略校验块

### 鉴权链完整顺序

```go
// new-api/middleware/auth.go — Token 鉴权函数（伪代码，展示插入位置）

func tokenAuth(c *gin.Context) {
    // 1. 提取 Token Key
    // 2. 查询 Token（不存在 → 401）
    // 3. Token 状态校验（禁用 → 403）
    // 4. Token 有效期校验（已过期 → 401）
    // 5. Token 额度校验（额度不足 → 429）
    // 6. Token 模型白名单校验（模型不允许 → 403）

    // ← 提取真实客户端 IP（Story 1.2.1 的变更）
    clientIP := common.GetClientIP(c)
    c.Set("client_ip", clientIP)

    // ← 7. IP 策略校验（本文件新增块）
    if token.IpPolicy != nil && token.IpPolicy.Mode != "" {
        cidrs, err := common.ParseCIDRList(token.IpPolicy.Ips)
        if err != nil {
            logger.LogWarn(c.Request.Context(), fmt.Sprintf(
                "ip_policy_parse_failed: token_id=%d err=%v", token.Id, err))
            // 降级：不拦截，继续
        } else {
            hit := common.IPMatchesCIDRList(clientIP, cidrs)
            blocked := (token.IpPolicy.Mode == "whitelist" && !hit) ||
                       (token.IpPolicy.Mode == "blacklist" && hit)
            if blocked {
                logger.LogInfo(c.Request.Context(), fmt.Sprintf(
                    "ip_not_allowed: token_id=%d client_ip=%s mode=%s",
                    token.Id, clientIP, token.IpPolicy.Mode))
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

    // 8. 设置 Context 并调用 c.Next()
    c.Set("id", token.UserId)
    c.Set("token_id", token.Id)
    // ...
    c.Next()
}
```

### 注意事项

1. **调用 `c.Abort()` 后必须立即 `return`**：否则函数体后续代码仍会执行（gin 的 `Abort` 只阻止后续中间件，不终止当前函数）
2. **客户端 IP 从 Context 读取（优先）**：若 Story 1.2.1 的变更已将 IP 写入 `c.Set("client_ip", ...)` 则复用，避免重复调用 `GetClientIP`
3. **降级不拦截的语义**：CIDR 解析失败时不拦截，意味着配置错误会导致策略失效（而非服务不可用），需通过 WARN 日志及时发现
4. **错误码统一**：无论白名单未命中还是黑名单命中，均返回 `IP_NOT_ALLOWED`，不区分，防止攻击者通过不同响应推断策略类型

---

## 监控与排障

| 场景 | 日志关键字 | 级别 | 排障方法 |
|------|-----------|------|---------|
| IP 被拒绝（正常拦截） | `ip_not_allowed` | INFO | 搜索 token_id 或 client_ip 定位来源 |
| CIDR 解析失败（降级） | `ip_policy_parse_failed` | WARN | 搜索 token_id，检查数据库 `ip_policy` 字段值 |
| 频繁 IP 拒绝（疑似攻击） | `ip_not_allowed` 频率异常 | 告警 | 结合 client_ip 分析攻击来源 |
| IP 策略未生效（策略配置后仍放行） | 无 `ip_not_allowed` 日志 | 排查 | 检查 Token.IpPolicy 是否正确写入；检查 TRUSTED_PROXIES 配置；检查鉴权链执行顺序 |
