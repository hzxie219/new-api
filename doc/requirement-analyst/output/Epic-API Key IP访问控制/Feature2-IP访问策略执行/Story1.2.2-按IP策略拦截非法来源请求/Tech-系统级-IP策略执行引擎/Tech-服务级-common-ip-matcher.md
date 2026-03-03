# 【Tech-服务级】common-ip-matcher

## 服务职责

提供 CIDR 解析、格式校验、IP 命中判断和可信代理管理的通用工具函数，是整个 IP 访问控制功能的核心算法库。供 controller 层入参校验（`ValidateCIDRList`）和 middleware 层策略执行（`ParseCIDRList`、`IPMatchesCIDRList`）调用，同时提供可信代理感知 IP 提取（`GetClientIP`）。

**所属 Tech-系统级**：Tech-系统级-IP 策略执行引擎

**代码位置**：`new-api/common/ip_matcher.go`（新建文件）

---

## 详细验收条件

### 组 A：CIDR 格式校验（ValidateCIDRList）

**AC-A1: 合法 CIDR 输入 — 返回 nil**
- **Given**: `ips = ["1.2.3.4", "10.0.0.0/8", "192.168.1.1/32"]`
- **When**: 调用 `ValidateCIDRList(ips)`
- **Then**: 返回 nil（无错误）

**AC-A2: 非法 IP — 返回 error**
- **Given**: `ips = ["999.0.0.1"]`
- **When**: 调用 `ValidateCIDRList(ips)`
- **Then**: 返回 error，message 含 `"invalid IP/CIDR: 999.0.0.1"`

**AC-A3: 非 IP 字符串 — 返回 error**
- **Given**: `ips = ["not-an-ip"]`
- **When**: 调用 `ValidateCIDRList(ips)`
- **Then**: 返回 error，message 含非法条目内容

**AC-A4: 混合列表 — 遇第一个非法条目返回 error**
- **Given**: `ips = ["1.2.3.4/32", "bad-ip", "10.0.0.0/8"]`
- **When**: 调用 `ValidateCIDRList(ips)`
- **Then**: 返回关于 `"bad-ip"` 的 error；不继续校验后续条目

**AC-A5: 空列表 — 返回 nil**
- **Given**: `ips = []`
- **When**: 调用 `ValidateCIDRList(ips)`
- **Then**: 返回 nil（空列表合法）

### 组 B：CIDR 列表解析（ParseCIDRList）

**AC-B1: 合法 CIDR 解析正确**
- **Given**: `ips = ["1.2.3.4", "10.0.0.0/8"]`
- **When**: 调用 `ParseCIDRList(ips)`
- **Then**: 返回长度为 2 的 `[]*net.IPNet`，err == nil；`1.2.3.4` 对应 `1.2.3.4/32`

**AC-B2: 单 IP 自动补全为 /32**
- **Given**: `ips = ["1.2.3.4"]`（无 CIDR 前缀）
- **When**: 调用 `ParseCIDRList(ips)` 后 `IPMatchesCIDRList("1.2.3.4", cidrs)`
- **Then**: 返回 true；`IPMatchesCIDRList("1.2.3.5", cidrs)` 返回 false

**AC-B3: 空列表 — 返回空 slice**
- **Given**: `ips = []`
- **When**: 调用 `ParseCIDRList(ips)`
- **Then**: 返回空 `[]*net.IPNet`，err == nil

**AC-B4: 非法 CIDR — 返回 error**
- **Given**: `ips = ["bad"]`
- **When**: 调用 `ParseCIDRList(ips)`
- **Then**: 返回 nil 和 error（含非法条目说明）

### 组 C：IP 命中判断（IPMatchesCIDRList）

**AC-C1: IP 命中 CIDR 网段**
- **Given**: `cidrs` 包含 `10.0.0.0/8`
- **When**: `IPMatchesCIDRList("10.255.255.255", cidrs)`
- **Then**: 返回 true

**AC-C2: IP 不命中 CIDR 网段**
- **Given**: `cidrs` 包含 `10.0.0.0/8`
- **When**: `IPMatchesCIDRList("11.0.0.1", cidrs)`
- **Then**: 返回 false

**AC-C3: 边界 IP 命中**
- **Given**: `cidrs` 包含 `172.16.0.0/12`（范围 `172.16.0.0`–`172.31.255.255`）
- **When**: `IPMatchesCIDRList("172.31.255.255", cidrs)`
- **Then**: 返回 true

**AC-C4: 边界 IP 不命中**
- **Given**: `cidrs` 包含 `172.16.0.0/12`
- **When**: `IPMatchesCIDRList("172.32.0.0", cidrs)`
- **Then**: 返回 false（超出范围）

**AC-C5: 空 CIDR 列表**
- **Given**: `cidrs = []`
- **When**: `IPMatchesCIDRList("1.2.3.4", [])`
- **Then**: 返回 false（空列表不命中任何 IP）

**AC-C6: IP 格式非法**
- **Given**: `ipStr = "not-an-ip"`
- **When**: `IPMatchesCIDRList("not-an-ip", cidrs)`
- **Then**: 返回 false（不 panic，容错处理）

**AC-C7: 多 CIDR 条目 — 命中任意一条即返回 true**
- **Given**: `cidrs` 包含 `["10.0.0.0/8", "192.168.0.0/16"]`
- **When**: `IPMatchesCIDRList("192.168.1.1", cidrs)`
- **Then**: 返回 true（命中第二条）

---

## 技术实现

### 完整代码

```go
// new-api/common/ip_matcher.go

package common

import (
    "fmt"
    "net"
    "strings"
)

// ============= CIDR 工具函数 =============

// ValidateCIDRList 校验 IP/CIDR 列表格式，返回第一个非法条目的 error
func ValidateCIDRList(ips []string) error {
    for _, ip := range ips {
        entry := ip
        if !strings.Contains(ip, "/") {
            entry = ip + "/32" // 单 IP 补全为 /32
        }
        if _, _, err := net.ParseCIDR(entry); err != nil {
            return fmt.Errorf("invalid IP/CIDR: %s", ip)
        }
    }
    return nil
}

// ParseCIDRList 将 IP/CIDR 字符串列表解析为 []*net.IPNet
func ParseCIDRList(ips []string) ([]*net.IPNet, error) {
    result := make([]*net.IPNet, 0, len(ips))
    for _, ip := range ips {
        entry := ip
        if !strings.Contains(ip, "/") {
            entry = ip + "/32"
        }
        _, cidr, err := net.ParseCIDR(entry)
        if err != nil {
            return nil, fmt.Errorf("invalid IP/CIDR: %s", ip)
        }
        result = append(result, cidr)
    }
    return result, nil
}

// IPMatchesCIDRList 检查 IP 字符串是否命中列表中的任意 CIDR
func IPMatchesCIDRList(ipStr string, cidrs []*net.IPNet) bool {
    ip := net.ParseIP(ipStr)
    if ip == nil {
        return false // 非法 IP 字符串，不命中
    }
    for _, cidr := range cidrs {
        if cidr.Contains(ip) {
            return true
        }
    }
    return false
}

// ============= 可信代理管理 =============

var trustedProxyCIDRs []*net.IPNet

// InitTrustedProxies 解析 TRUSTED_PROXIES 环境变量，在服务启动时调用一次
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
            // 非法条目跳过，调用方（init.go）记录 WARN 日志
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

// GetClientIP 从 gin.Context 中提取真实客户端 IP
// 若来自可信代理，取 X-Forwarded-For 最左侧 IP；否则取 RemoteAddr
func GetClientIP(c *gin.Context) string {
    host, _, err := net.SplitHostPort(c.Request.RemoteAddr)
    if err != nil {
        host = c.Request.RemoteAddr // 格式异常，回退
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

### 单元测试覆盖要求

测试文件：`new-api/common/ip_matcher_test.go`

| 测试组 | 覆盖场景 |
|--------|---------|
| `TestValidateCIDRList` | 合法 CIDR、非法 IP、非法字符串、混合列表、空列表 |
| `TestParseCIDRList` | 合法列表、单 IP 补全、空列表、非法条目 |
| `TestIPMatchesCIDRList` | 命中/不命中、边界 IP、空列表、非法 IP 字符串、多条目 |
| `TestGetClientIP` | 直连、可信代理、XFF 伪造、多跳、IPv6、XFF 空 |
| `TestInitTrustedProxies` | 合法 CIDR、非法条目跳过、空配置、单 IP |

### 注意事项

1. **`net.ParseCIDR` 行为**：会将主机位归零（如 `1.2.3.4/24` → `1.2.3.0/24`）；这意味着 `ParseCIDRList(["1.2.3.4/24"])` 后，`IPMatchesCIDRList("1.2.3.100", ...)` 返回 true；行为符合预期（CIDR 表示网段而非单主机）
2. **单 IP 补全为 /32**：`net.ParseCIDR("1.2.3.4/32")` 不归零主机位，因此 `ParseCIDRList(["1.2.3.4"])` 后精确匹配 `1.2.3.4`，不匹配 `1.2.3.5`
3. **本文件无 IO 操作**：纯函数（除全局变量 `trustedProxyCIDRs`），便于单元测试
4. **`gin` 包引用**：`GetClientIP` 依赖 `gin.Context`，如需在非 gin 场景使用，需重构为接受 `*http.Request` 的版本

---

## 监控与排障

| 场景 | 处理方式 |
|------|---------|
| `ParseCIDRList` 运行时失败 | 调用方（middleware）降级不拦截，记录 WARN 日志；说明数据库数据被篡改 |
| `ValidateCIDRList` 入参校验失败 | 调用方（controller）返回 HTTP 400，用户修正输入 |
| `GetClientIP` 返回异常 IP 格式 | 记录 WARN 日志，后续 `IPMatchesCIDRList` 对非法 IP 返回 false（降级不拦截） |
