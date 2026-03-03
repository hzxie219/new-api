# 【Tech-系统级】可信代理感知 IP 提取

## 【技术故事】

作为系统，需要在鉴权中间件中提供统一的"真实客户端 IP 提取"能力：根据可信代理配置决定是否信任 `X-Forwarded-For` 头，从而为 IP 策略校验（Tech-系统级-IP策略执行引擎）提供可靠的客户端 IP 来源，防止攻击者通过伪造 XFF 头绕过 IP 访问控制。

**所属 Story**：Story 1.2.1 — 请求鉴权时获取真实客户端 IP

**子系统职责**：
1. 维护可信代理 IP/CIDR 列表（通过环境变量配置，服务启动时预编译为 `[]*net.IPNet`）
2. 在鉴权时判断 `RemoteAddr` 是否属于可信代理
3. 按规则提取真实客户端 IP（可信代理 → 取 XFF 最左侧 IP；非可信代理 → 取 `RemoteAddr`）

---

## 【验收条件】

### 一、功能性验收条件

**TC-1: 可信代理转发 — 正确取真实 IP**
- **Given**: `TRUSTED_PROXIES=192.168.1.0/24`；请求 `RemoteAddr=192.168.1.10:54321`，`X-Forwarded-For: 1.2.3.4, 192.168.1.10`
- **When**: 调用 `GetClientIP(c)`
- **Then**: 返回 `"1.2.3.4"`（XFF 最左侧 IP）

**TC-2: 直连请求 — 忽略 XFF**
- **Given**: `TRUSTED_PROXIES=192.168.1.0/24`；请求 `RemoteAddr=8.8.8.8:443`，`X-Forwarded-For: 1.2.3.4`
- **When**: 调用 `GetClientIP(c)`
- **Then**: 返回 `"8.8.8.8"`，忽略伪造的 XFF 头

**TC-3: `TRUSTED_PROXIES` 未配置 — 默认不信任任何代理**
- **Given**: `TRUSTED_PROXIES` 为空或未设置
- **When**: 任意请求调用 `GetClientIP(c)`
- **Then**: 始终返回 `RemoteAddr` 的 IP 部分（最保守策略）

**TC-4: XFF 多跳 — 取最左侧**
- **Given**: `TRUSTED_PROXIES=192.168.1.0/24`；请求可信；`X-Forwarded-For: 1.2.3.4, 5.5.5.5, 192.168.1.1`
- **When**: 调用 `GetClientIP(c)`
- **Then**: 返回 `"1.2.3.4"`（最左侧，最原始客户端 IP）

**TC-5: XFF 头为空 — 回退到 RemoteAddr**
- **Given**: `TRUSTED_PROXIES=192.168.1.0/24`；请求 `RemoteAddr=192.168.1.5:1234`，无 XFF 头
- **When**: 调用 `GetClientIP(c)`
- **Then**: 返回 `"192.168.1.5"`（回退到 RemoteAddr）

**TC-6: XFF 值含空格 — 正确 Trim**
- **Given**: 可信代理请求；`X-Forwarded-For: " 1.2.3.4 , 5.5.5.5"`
- **When**: 调用 `GetClientIP(c)`
- **Then**: 返回 `"1.2.3.4"`（TrimSpace 后的第一个 IP）

**TC-7: IPv6 格式 RemoteAddr 正确解析**
- **Given**: 请求 `RemoteAddr=[::1]:5678`
- **When**: 调用 `GetClientIP(c)`
- **Then**: 返回 `"::1"`，不含端口，不 panic

**TC-8: `InitTrustedProxies` 非法 CIDR 跳过**
- **Given**: `TRUSTED_PROXIES=192.168.1.0/24,invalid-cidr,10.0.0.0/8`
- **When**: 服务启动解析配置
- **Then**: 合法的 `192.168.1.0/24` 和 `10.0.0.0/8` 正常加入列表；`invalid-cidr` 跳过并记录 WARN 日志；服务正常启动

### 二、非功能性验收条件

**TC-N1: 函数性能**
- **Given**: `trustedProxyCIDRs` 包含 50 个 CIDR 条目（最大预期配置量）
- **When**: 高并发场景下连续调用 `GetClientIP(c)`（QPS=1000）
- **Then**: 函数耗时 P99 < 0.1ms，不成为鉴权链的性能瓶颈

**TC-N2: 并发安全**
- **Given**: `trustedProxyCIDRs` 为全局变量（启动时写入，运行时只读）
- **When**: 多 goroutine 并发调用 `GetClientIP` 和 `isTrustedProxy`
- **Then**: 无数据竞争（race detector 检测通过）

---

## 【依赖与风险】

| 项 | 说明 |
|----|------|
| 依赖 `common/init.go` | 调用 `InitTrustedProxies(os.Getenv("TRUSTED_PROXIES"))` |
| 被 `middleware/auth.go` 调用 | 在 Token 鉴权成功后调用 `GetClientIP()` |
| 被 Tech-系统级-IP策略执行引擎 依赖 | 提供客户端 IP 供策略校验使用 |
| 风险：配置过宽 | `TRUSTED_PROXIES=0.0.0.0/0` 导致任何人都可伪造 XFF；部署文档须强制说明 |
| 风险：内网机器被攻击 | 配置为内网段（如 `10.0.0.0/8`）时，内网任意机器均可伪造来源 IP |

---

## 【技术思路】

### 可信代理配置与初始化

```go
// new-api/common/ip_matcher.go

var trustedProxyCIDRs []*net.IPNet

// InitTrustedProxies 解析逗号分隔的 IP/CIDR 列表，在服务启动时调用一次
func InitTrustedProxies(proxies string) {
    trustedProxyCIDRs = nil // 重置（支持测试中多次调用）
    for _, p := range strings.Split(proxies, ",") {
        p = strings.TrimSpace(p)
        if p == "" {
            continue
        }
        // 单 IP 补全为 /32
        if !strings.Contains(p, "/") {
            p = p + "/32"
        }
        _, cidr, err := net.ParseCIDR(p)
        if err != nil {
            // 记录 WARN 但继续（跳过非法条目）
            continue
        }
        trustedProxyCIDRs = append(trustedProxyCIDRs, cidr)
    }
}
```

```go
// new-api/common/init.go — 在现有初始化逻辑末尾追加
common.InitTrustedProxies(os.Getenv("TRUSTED_PROXIES"))
```

### IP 提取函数

```go
// new-api/common/ip_matcher.go

func isTrustedProxy(ip net.IP) bool {
    for _, cidr := range trustedProxyCIDRs {
        if cidr.Contains(ip) {
            return true
        }
    }
    return false
}

// GetClientIP 从 gin.Context 中提取真实客户端 IP
func GetClientIP(c *gin.Context) string {
    host, _, err := net.SplitHostPort(c.Request.RemoteAddr)
    if err != nil {
        // RemoteAddr 格式异常（如无端口），回退使用原始值
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

### 部署配置示例

```bash
# 单层 Nginx 反向代理
TRUSTED_PROXIES=10.0.0.1

# 多个代理节点
TRUSTED_PROXIES=10.0.0.1,10.0.0.2

# 代理网段（k8s ingress 常见场景）
TRUSTED_PROXIES=10.0.0.0/8,172.16.0.0/12

# 不使用代理（直连，安全默认）
TRUSTED_PROXIES=
```
