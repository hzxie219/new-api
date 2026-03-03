# 【Tech-服务级】middleware-auth（IP 提取部分）

## 服务职责

在 `new-api/middleware/auth.go` 的 Token 鉴权流程中，新增调用 `common.GetClientIP(c)` 的逻辑，为后续 IP 策略校验提供可靠的真实客户端 IP。本文件仅描述 IP 提取相关的变更，不包含策略校验逻辑（见 Tech-服务级-middleware-auth-策略执行）。

**所属 Tech-系统级**：Tech-系统级-可信代理感知 IP 提取

---

## 详细验收条件

### AC-1: 可信代理请求 — 正确提取 XFF IP

- **Given**: 环境变量 `TRUSTED_PROXIES=192.168.1.0/24`；请求 `RemoteAddr=192.168.1.5:12345`，`X-Forwarded-For: 1.2.3.4`
- **When**: Token 鉴权有效，中间件调用 `common.GetClientIP(c)`
- **Then**: 返回 `"1.2.3.4"`

### AC-2: 非可信代理 — 使用 RemoteAddr

- **Given**: `TRUSTED_PROXIES=192.168.1.0/24`；请求 `RemoteAddr=8.8.8.8:443`，`X-Forwarded-For: 1.2.3.4`
- **When**: Token 鉴权有效，调用 `common.GetClientIP(c)`
- **Then**: 返回 `"8.8.8.8"`，忽略伪造的 XFF 头

### AC-3: `TRUSTED_PROXIES` 未配置

- **Given**: 服务启动时未设置 `TRUSTED_PROXIES` 环境变量
- **When**: 任意请求调用 `GetClientIP(c)`
- **Then**: 始终返回 `RemoteAddr` 的 IP 部分

### AC-4: XFF 多跳 — 取最左侧

- **Given**: 可信代理命中；`X-Forwarded-For: 1.2.3.4, 5.5.5.5`
- **When**: 调用 `GetClientIP(c)`
- **Then**: 返回 `"1.2.3.4"`

### AC-5: RemoteAddr 带端口 — 正确解析 IP

- **Given**: `RemoteAddr=192.168.1.5:54321`
- **When**: 调用 `GetClientIP(c)`
- **Then**: 返回 `"192.168.1.5"`，不含端口号

### AC-6: IPv6 RemoteAddr 正确解析

- **Given**: `RemoteAddr=[::1]:3000`（IPv6 本地地址）
- **When**: 调用 `GetClientIP(c)`
- **Then**: 返回 `"::1"`，不 panic，不含方括号和端口

### AC-7: IP 提取在 Token 有效性校验通过后执行

- **Given**: Token 已过期或额度为 0
- **When**: 请求到达鉴权中间件
- **Then**: 中间件在 Token 有效性校验失败时直接返回错误，不执行 `GetClientIP()`（IP 提取仅对有效 Token 生效）

### AC-8: IP 提取结果存入 Context

- **Given**: Token 有效，IP 提取成功
- **When**: 调用 `GetClientIP(c)` 后
- **Then**: 将客户端 IP 存入 gin.Context（如 `c.Set("client_ip", clientIP)`），供同一请求的后续逻辑（包括策略执行和日志）复用，避免重复解析

---

## 技术实现

### 代码位置

| 文件 | 变更内容 |
|------|---------|
| `new-api/common/ip_matcher.go` | 新增 `GetClientIP`、`isTrustedProxy`、`InitTrustedProxies` |
| `new-api/common/init.go` | 追加 `InitTrustedProxies(os.Getenv("TRUSTED_PROXIES"))` |
| `new-api/middleware/auth.go` | Token 鉴权成功后调用 `GetClientIP(c)` 并存入 Context |

### 核心代码参考

```go
// new-api/common/ip_matcher.go

var trustedProxyCIDRs []*net.IPNet

func InitTrustedProxies(proxies string) {
    trustedProxyCIDRs = nil
    for _, p := range strings.Split(proxies, ",") {
        p = strings.TrimSpace(p)
        if p == "" {
            continue
        }
        if !strings.Contains(p, "/") {
            p = p + "/32"
        }
        _, cidr, err := net.ParseCIDR(p)
        if err != nil {
            // 跳过非法条目，记录 WARN 日志（调用方处理）
            continue
        }
        trustedProxyCIDRs = append(trustedProxyCIDRs, cidr)
    }
}

func isTrustedProxy(ip net.IP) bool {
    for _, cidr := range trustedProxyCIDRs {
        if cidr.Contains(ip) {
            return true
        }
    }
    return false
}

func GetClientIP(c *gin.Context) string {
    host, _, err := net.SplitHostPort(c.Request.RemoteAddr)
    if err != nil {
        host = c.Request.RemoteAddr
    }
    remoteIP := net.ParseIP(host)
    if remoteIP != nil && isTrustedProxy(remoteIP) {
        if xff := c.GetHeader("X-Forwarded-For"); xff != "" {
            first := strings.TrimSpace(strings.SplitN(xff, ",", 2)[0])
            if first != "" {
                return first
            }
        }
    }
    return host
}
```

```go
// new-api/middleware/auth.go — Token 鉴权成功后（其他有效性校验通过之后）追加

// 提取并缓存客户端真实 IP
clientIP := common.GetClientIP(c)
c.Set("client_ip", clientIP)
// （后续 IP 策略执行块从 Context 读取，无需重复调用）
```

```go
// new-api/common/init.go — 追加在现有初始化逻辑末尾
common.InitTrustedProxies(os.Getenv("TRUSTED_PROXIES"))
```

### 注意事项

1. `InitTrustedProxies` 为全局变量写入，仅在服务启动时调用一次；运行时 `GetClientIP` 只读 `trustedProxyCIDRs`，并发安全
2. `GetClientIP` 结果存入 `c.Set("client_ip", clientIP)` 后，IP 策略执行（Tech-服务级-middleware-auth-策略执行）从 Context 读取，避免重复解析
3. IPv6 地址的 `net.SplitHostPort` 可正确处理 `[::1]:port` 格式

---

## 监控与排障

| 场景 | 日志关键字 | 处理方式 |
|------|-----------|---------|
| RemoteAddr 格式异常 | `get_client_ip_warn: invalid remote_addr` | WARN，记录原始 RemoteAddr，回退使用原始值 |
| `InitTrustedProxies` 非法条目 | `trusted_proxy_parse_warn` | WARN，记录非法 CIDR 值，跳过继续初始化 |
| 怀疑 XFF 伪造（审计） | `get_client_ip: xff_ignored` | DEBUG（可选），记录被忽略的 XFF 值，便于安全审计 |
