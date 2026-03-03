# 【Tech-系统级】IP 策略执行引擎

## 技术故事

作为系统，需要在 API Key 鉴权流程中，根据 Key 绑定的 IP 策略（whitelist/blacklist + CIDR 列表）对请求来源 IP 进行匹配判断，对违规请求立即拒绝并返回标准化的 403 响应。

## 所属 Story

Story 1.2.2 — 按 IP 策略拦截非法来源请求

## 子系统职责

本子系统负责：
1. CIDR 解析与 IP 命中判断的工具函数（`common` 层）
2. 在鉴权中间件中读取 Token 的 IP 策略并调用工具函数进行校验
3. 违规时返回 HTTP 403 + 错误码 `IP_NOT_ALLOWED`

## 验收条件（Tech 级）

### TC-1: 白名单 — 命中允许通过

- **Given**: Token IpPolicy: `{mode:"whitelist", ips:["1.2.3.4/32"]}`；客户端 IP: `1.2.3.4`
- **When**: 执行 IP 策略校验
- **Then**: 校验通过，不拒绝请求

### TC-2: 白名单 — 未命中拒绝

- **Given**: Token IpPolicy: `{mode:"whitelist", ips:["1.2.3.4/32"]}`；客户端 IP: `8.8.8.8`
- **When**: 执行 IP 策略校验
- **Then**: 校验失败，返回 HTTP 403，body: `{"success":false,"message":"IP not allowed","error_code":"IP_NOT_ALLOWED"}`

### TC-3: 黑名单 — 命中拒绝

- **Given**: Token IpPolicy: `{mode:"blacklist", ips:["5.5.5.0/24"]}`；客户端 IP: `5.5.5.100`
- **When**: 执行 IP 策略校验
- **Then**: 校验失败，返回 HTTP 403 + `IP_NOT_ALLOWED`

### TC-4: 黑名单 — 未命中允许通过

- **Given**: Token IpPolicy: `{mode:"blacklist", ips:["5.5.5.0/24"]}`；客户端 IP: `8.8.8.8`
- **When**: 执行 IP 策略校验
- **Then**: 校验通过

### TC-5: 无策略 — 跳过校验

- **Given**: Token IpPolicy 为 nil 或 mode 为空
- **When**: 执行 IP 策略校验
- **Then**: 跳过，不影响请求处理

### TC-6: CIDR 匹配精度

- **Given**: ips 列表包含 `"10.0.0.0/8"`；分别测试 IP `10.255.255.255`（应命中）和 `11.0.0.1`（不应命中）
- **When**: 执行 CIDR 匹配
- **Then**: 匹配结果正确

## 技术思路

### CIDR 匹配工具（新增）

`new-api/common/ip_matcher.go`:

```go
// ParseCIDRList 预解析 CIDR 列表（建议在策略加载时调用，缓存结果）
func ParseCIDRList(ips []string) ([]*net.IPNet, error)

// IPMatchesCIDRList 检查 IP 是否命中列表中任意 CIDR
func IPMatchesCIDRList(ip string, cidrs []*net.IPNet) bool
```

### 策略执行逻辑（中间件层）

在 `new-api/middleware/auth.go` Token 鉴权通过后追加：

```go
policy := token.IpPolicy
if policy != nil && policy.Mode != "" {
    clientIP := common.GetClientIP(c)
    cidrs, _ := common.ParseCIDRList(policy.Ips)
    hit := common.IPMatchesCIDRList(clientIP, cidrs)

    blocked := (policy.Mode == "whitelist" && !hit) ||
               (policy.Mode == "blacklist" && hit)
    if blocked {
        c.JSON(http.StatusForbidden, gin.H{
            "success":    false,
            "message":    "IP not allowed",
            "error_code": "IP_NOT_ALLOWED",
        })
        c.Abort()
        return
    }
}
```

### 依赖关系

- 依赖 Tech-系统级-可信代理感知IP提取 提供的 `common.GetClientIP()`
- 依赖 `new-api/model/token.go` 的 `IpPolicy` 字段
- 执行位置：`new-api/middleware/auth.go` Token 鉴权成功后

## 风险与依赖

- CIDR 解析建议在策略写入时做一次预校验（入参校验阶段），避免运行时解析失败
- 性能：每次请求均需解析 CIDR 列表，建议对 Token 的已解析 CIDR 做轻量缓存（同一 Token 在短时间内请求频繁时复用）
