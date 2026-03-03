# 【Tech-服务级】common-ip-matcher

## 服务职责

提供 CIDR 解析、IP 命中判断的通用工具函数，供 controller 层的入参校验和 middleware 层的策略执行调用。是 IP 访问控制功能的核心算法库。

## 所属 Tech-系统级

Tech-系统级-IP 策略执行引擎

## 详细验收条件

### AC-1: CIDR 列表校验 — 合法输入

- **Given**: `ips = ["1.2.3.4", "10.0.0.0/8", "192.168.1.1/32"]`
- **When**: 调用 `ValidateCIDRList(ips)`
- **Then**: 返回 nil error

### AC-2: CIDR 列表校验 — 非法输入

- **Given**: `ips = ["999.0.0.1", "not-an-ip"]`
- **When**: 调用 `ValidateCIDRList(ips)`
- **Then**: 返回 error，message 包含非法条目内容

### AC-3: IP 命中 CIDR

- **Given**: cidrs 包含 `10.0.0.0/8`
- **When**: `IPMatchesCIDRList("10.255.255.255", cidrs)`
- **Then**: 返回 true

### AC-4: IP 不命中 CIDR

- **Given**: cidrs 包含 `10.0.0.0/8`
- **When**: `IPMatchesCIDRList("11.0.0.1", cidrs)`
- **Then**: 返回 false

### AC-5: 精确 IP 匹配（单 IP 视为 /32）

- **Given**: ips 传入 `["1.2.3.4"]`（无 CIDR 前缀）
- **When**: `ParseCIDRList(["1.2.3.4"])` 后 `IPMatchesCIDRList("1.2.3.4", ...)`
- **Then**: 返回 true；`IPMatchesCIDRList("1.2.3.5", ...)` 返回 false

### AC-6: 空列表

- **Given**: `ips = []`
- **When**: `IPMatchesCIDRList("1.2.3.4", [])`
- **Then**: 返回 false（空列表不命中任何 IP）

## 技术实现

### 代码位置

`new-api/common/ip_matcher.go`（新建文件）

### 核心代码参考

```go
package common

import (
    "fmt"
    "net"
    "strings"
)

// ValidateCIDRList 校验 IP/CIDR 列表格式，返回第一个非法条目的错误
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
    var result []*net.IPNet
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

// IPMatchesCIDRList 检查 IP 是否命中列表中的任意 CIDR
func IPMatchesCIDRList(ipStr string, cidrs []*net.IPNet) bool {
    ip := net.ParseIP(ipStr)
    if ip == nil {
        return false
    }
    for _, cidr := range cidrs {
        if cidr.Contains(ip) {
            return true
        }
    }
    return false
}
```

### 单元测试

测试文件: `new-api/common/ip_matcher_test.go`，覆盖：
- 合法 CIDR 校验通过
- 非法 CIDR 返回错误
- 精确 IP 命中/不命中
- 网段边界命中/不命中
- 空列表返回 false

### 注意事项

- 单 IP（无 `/`）自动补全为 `/32`，不要求调用方显式写 CIDR 格式
- `net.ParseCIDR` 会将主机位归零（如 `1.2.3.4/24` → `1.2.3.0/24`），校验和匹配时需注意
- 本文件不依赖任何业务包，保持纯工具函数，方便单元测试

## 监控与排障

- 该模块为纯函数，无 IO，无需特殊监控
- 若 `ParseCIDRList` 在运行时调用失败，调用方（middleware）应记录 WARN 日志并视为无策略（降级为不拦截）
