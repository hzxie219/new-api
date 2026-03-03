# 【Tech-服务级】middleware-auth（IP 提取部分）

## 服务职责

在 `new-api/middleware/auth.go` 的 Token 鉴权流程中，新增获取真实客户端 IP 的逻辑：根据可信代理配置，决定是否信任 `X-Forwarded-For` 头并提取真实 IP，为后续 IP 策略校验提供可靠的客户端 IP 来源。

## 所属 Tech-系统级

Tech-系统级-可信代理感知 IP 提取

## 详细验收条件

### AC-1: 可信代理请求 — 正确提取 XFF IP

- **Given**: 环境变量 `TRUSTED_PROXIES=192.168.1.0/24`；请求 `RemoteAddr=192.168.1.5:12345`，`X-Forwarded-For: 1.2.3.4`
- **When**: 中间件提取客户端 IP
- **Then**: 返回 `1.2.3.4`

### AC-2: 非可信代理 — 使用 RemoteAddr

- **Given**: 环境变量 `TRUSTED_PROXIES=192.168.1.0/24`；请求 `RemoteAddr=8.8.8.8:443`，`X-Forwarded-For: 1.2.3.4`
- **When**: 中间件提取客户端 IP
- **Then**: 返回 `8.8.8.8`，忽略 XFF

### AC-3: TRUSTED_PROXIES 未配置

- **Given**: 未设置环境变量
- **When**: 任意请求
- **Then**: 始终返回 `RemoteAddr`（默认不信任任何代理）

### AC-4: XFF 多跳 — 取最左侧

- **Given**: 可信代理命中；`X-Forwarded-For: 1.2.3.4, 5.5.5.5`
- **When**: 提取客户端 IP
- **Then**: 返回 `1.2.3.4`

### AC-5: RemoteAddr 带端口 — 正确解析 IP

- **Given**: `RemoteAddr=192.168.1.5:54321`
- **When**: 解析 RemoteAddr
- **Then**: 正确提取 IP 部分 `192.168.1.5`，不含端口

## 技术实现

### 代码位置

- 工具函数: `new-api/common/ip_matcher.go`（新增 `GetClientIP`、`isTrustedProxy`）
- 全局变量: `new-api/common/init.go`（解析 `TRUSTED_PROXIES` 环境变量）
- 调用位置: `new-api/middleware/auth.go` Token 鉴权函数中

### 核心代码参考

```go
// new-api/common/ip_matcher.go

var trustedProxyCIDRs []*net.IPNet

// InitTrustedProxies 在 init.go 中调用
func InitTrustedProxies(proxies string) {
    for _, p := range strings.Split(proxies, ",") {
        p = strings.TrimSpace(p)
        if p == "" { continue }
        _, cidr, err := net.ParseCIDR(p)
        if err == nil {
            trustedProxyCIDRs = append(trustedProxyCIDRs, cidr)
        }
    }
}

func isTrustedProxy(ip net.IP) bool {
    for _, cidr := range trustedProxyCIDRs {
        if cidr.Contains(ip) { return true }
    }
    return false
}

func GetClientIP(c *gin.Context) string {
    host, _, err := net.SplitHostPort(c.Request.RemoteAddr)
    if err != nil { host = c.Request.RemoteAddr }
    remoteIP := net.ParseIP(host)
    if remoteIP != nil && isTrustedProxy(remoteIP) {
        if xff := c.GetHeader("X-Forwarded-For"); xff != "" {
            if first := strings.TrimSpace(strings.Split(xff, ",")[0]); first != "" {
                return first
            }
        }
    }
    return host
}
```

### 注意事项

- `TRUSTED_PROXIES` 解析在服务启动时完成（`common/init.go`），避免每次请求重复解析
- IPv6 地址格式 `[::1]:port` 需使用 `net.SplitHostPort` 解析（代码已覆盖）

## 监控与排障

- 若 IP 提取失败（RemoteAddr 格式异常），记录 WARN 日志并回退到空字符串
- 日志关键字: `get_client_ip_failed`
