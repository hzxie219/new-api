# 【Tech-系统级】可信代理感知 IP 提取

## 技术故事

作为系统，需要在鉴权中间件中提供统一的"真实客户端 IP 提取"能力：根据可信代理配置决定是否信任 `X-Forwarded-For` 头，从而为 IP 策略校验提供可靠的客户端 IP。

## 所属 Story

Story 1.2.1 — 请求鉴权时获取真实客户端 IP

## 子系统职责

本子系统负责：
1. 维护可信代理 IP/CIDR 列表配置（支持环境变量或系统配置读取）
2. 在鉴权时判断 `RemoteAddr` 是否在可信代理列表中
3. 按规则提取真实客户端 IP（可信代理 → 取 XFF 最左侧 IP；非可信代理 → 取 `RemoteAddr`）

## 验收条件（Tech 级）

### TC-1: 可信代理转发，正确取真实 IP

- **Given**: `TRUSTED_PROXIES=192.168.1.0/24`；请求 `RemoteAddr=192.168.1.10`，`X-Forwarded-For: 1.2.3.4, 192.168.1.10`
- **When**: 执行 IP 提取逻辑
- **Then**: 返回 `1.2.3.4`

### TC-2: 直连请求，忽略 XFF

- **Given**: `TRUSTED_PROXIES=192.168.1.0/24`；请求 `RemoteAddr=8.8.8.8`，`X-Forwarded-For: 1.2.3.4`
- **When**: 执行 IP 提取逻辑
- **Then**: 返回 `8.8.8.8`，忽略伪造的 XFF 头

### TC-3: 未配置可信代理，默认行为

- **Given**: `TRUSTED_PROXIES` 未设置（空）
- **When**: 任意请求到达鉴权中间件
- **Then**: 始终使用 `RemoteAddr` 作为客户端 IP（最保守策略）

### TC-4: XFF 包含多跳

- **Given**: 可信代理配置正确；`X-Forwarded-For: 1.2.3.4, 5.5.5.5, 192.168.1.1`
- **When**: 执行 IP 提取逻辑
- **Then**: 返回 `1.2.3.4`（最左侧，即最原始的客户端 IP）

## 技术思路

### 可信代理配置

通过环境变量 `TRUSTED_PROXIES` 配置，逗号分隔，支持 CIDR：
```
TRUSTED_PROXIES=127.0.0.1,192.168.0.0/16,10.0.0.0/8
```

在 `new-api/common/init.go` 中解析并存为全局变量，启动时预编译 CIDR 为 `net.IPNet` 列表。

### IP 提取函数

建议新增 `new-api/common/ip_matcher.go`：

```go
// GetClientIP 从 gin.Context 中提取真实客户端 IP
func GetClientIP(c *gin.Context) string {
    remoteIP := net.ParseIP(strings.Split(c.Request.RemoteAddr, ":")[0])
    if isTrustedProxy(remoteIP) {
        // 取 XFF 最左侧有效 IP
        if xff := c.GetHeader("X-Forwarded-For"); xff != "" {
            ips := strings.Split(xff, ",")
            if clientIP := strings.TrimSpace(ips[0]); clientIP != "" {
                return clientIP
            }
        }
    }
    return remoteIP.String()
}
```

### 依赖关系

- 依赖 `new-api/common/init.go`: 可信代理列表全局变量
- 由 `new-api/middleware/auth.go` 调用

## 风险与依赖

- `TRUSTED_PROXIES` 默认值建议谨慎：默认为空（最安全）或仅信任 `127.0.0.1`
- IPv6 地址的 RemoteAddr 格式（`[::1]:port`）需要特殊解析
