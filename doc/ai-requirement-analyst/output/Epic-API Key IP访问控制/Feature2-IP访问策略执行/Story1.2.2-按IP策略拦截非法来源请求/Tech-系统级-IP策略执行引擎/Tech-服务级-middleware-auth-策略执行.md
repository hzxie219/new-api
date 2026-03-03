# 【Tech-服务级】middleware-auth（策略执行部分）

## 服务职责

在 `new-api/middleware/auth.go` 的 Token 鉴权流程中，读取 Token 的 IP 策略，结合 `GetClientIP()` 提供的真实客户端 IP，调用 CIDR 匹配工具执行白名单/黑名单校验，对违规请求立即返回 HTTP 403 + `IP_NOT_ALLOWED`，并终止请求处理链。

## 所属 Tech-系统级

Tech-系统级-IP 策略执行引擎

## 详细验收条件

### AC-1: 白名单命中 — 请求通过

- **Given**: Token IpPolicy `{mode:"whitelist", ips:["1.2.3.4/32"]}`；客户端 IP `1.2.3.4`
- **When**: 中间件执行 IP 策略校验
- **Then**: 调用 `c.Next()` 继续处理，不产生拒绝响应

### AC-2: 白名单未命中 — 返回 403

- **Given**: Token IpPolicy `{mode:"whitelist", ips:["1.2.3.4/32"]}`；客户端 IP `8.8.8.8`
- **When**: 中间件执行 IP 策略校验
- **Then**:
  - HTTP 状态码 403
  - 响应体: `{"success":false,"message":"IP not allowed","error_code":"IP_NOT_ALLOWED"}`
  - 调用 `c.Abort()`，不继续执行后续中间件和 handler

### AC-3: 黑名单命中 — 返回 403

- **Given**: Token IpPolicy `{mode:"blacklist", ips:["5.5.5.0/24"]}`；客户端 IP `5.5.5.100`
- **When**: 中间件执行 IP 策略校验
- **Then**: HTTP 403 + `IP_NOT_ALLOWED`，`c.Abort()`

### AC-4: 黑名单未命中 — 请求通过

- **Given**: Token IpPolicy `{mode:"blacklist", ips:["5.5.5.0/24"]}`；客户端 IP `8.8.8.8`
- **When**: 中间件执行 IP 策略校验
- **Then**: 调用 `c.Next()`，请求正常处理

### AC-5: 无策略 — 跳过校验

- **Given**: Token IpPolicy 为 nil 或 mode 为空字符串
- **When**: 中间件执行 IP 策略校验
- **Then**: 跳过校验，调用 `c.Next()`

### AC-6: IP 策略校验在 Token 有效性校验之后执行

- **Given**: Token 已过期或额度为 0
- **When**: 请求到达鉴权中间件
- **Then**: 先返回过期/额度不足错误，不执行 IP 策略校验（IP 策略仅对有效 Token 生效）

## 技术实现

### 代码位置

`new-api/middleware/auth.go` — 在现有 Token 鉴权逻辑（额度/有效期/IP 白名单检查）完成后追加

### 核心代码参考

```go
// new-api/middleware/auth.go — Token 鉴权函数内，Token 有效性校验通过后追加

// IP 策略校验
if token.IpPolicy != nil && token.IpPolicy.Mode != "" {
    clientIP := common.GetClientIP(c)
    cidrs, err := common.ParseCIDRList(token.IpPolicy.Ips)
    if err != nil {
        // CIDR 解析失败：降级处理，记录警告，不拦截
        logger.LogWarn(c.Request.Context(), "ip_policy parse failed: "+err.Error())
    } else {
        hit := common.IPMatchesCIDRList(clientIP, cidrs)
        blocked := (token.IpPolicy.Mode == "whitelist" && !hit) ||
                   (token.IpPolicy.Mode == "blacklist" && hit)
        if blocked {
            logger.LogInfo(c.Request.Context(), fmt.Sprintf(
                "ip_not_allowed: token_id=%d client_ip=%s mode=%s",
                token.Id, clientIP, token.IpPolicy.Mode,
            ))
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

### 执行位置

在 Token 鉴权链中的位置（建议顺序）：
1. Token 存在性校验
2. Token 状态校验（启用/禁用）
3. Token 有效期校验
4. Token 额度校验
5. Token 模型白名单校验
6. **← Token IP 策略校验（本文件新增）**
7. `c.Next()`

### 注意事项

- CIDR 解析失败时**降级不拦截**（宽松策略），避免因配置异常导致服务不可用
- 降级时必须记录 WARN 日志，便于运维发现配置问题
- 错误响应格式与现有 Token 鉴权错误保持一致（`success/message/error_code` 三字段）
- 调用 `c.Abort()` 后必须 `return`，确保函数退出

## 监控与排障

- IP 拦截事件日志关键字: `ip_not_allowed`，可通过 request_id 追踪
- 建议在日志中记录: token_id、client_ip、mode，便于审计
- 若 `ERROR_LOG_ENABLED=true`，可将 IP 拒绝事件写入数据库 logs 表（type=5）便于后台查询
