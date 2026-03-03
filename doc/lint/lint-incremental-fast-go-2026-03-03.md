# Go 代码规范检查报告

**模式**: 快速模式 (fast)
**范围**: 增量检查 (incremental) — 对比 HEAD 未提交变更
**语言**: Go
**日期**: 2026-03-03
**分支**: main

---

## 📊 检查概要

| 类型 | 数量 |
|------|------|
| 错误 (error) | 6 |
| 警告 (warning) | 1 |
| 建议 (suggestion) | 0 |
| 安全问题 | 0 |
| **合计** | **7** |

**检查文件** (已排除测试文件 `common/ip_matcher_test.go`):

- `common/ip_matcher.go` (新增文件)
- `common/init.go` (第 124 行)
- `controller/token.go` (第 294–356 行)
- `middleware/auth.go` (第 332–355 行)
- `model/token.go` (第 4, 15–52, 67, 360–376 行)
- `router/api-router.go` (第 251 行)
- `types/error.go` (第 67 行)

---

## 🚨 错误问题（Error）

### E1 · 命名规范：导出标识符中的缩写词应全大写

**规则**: 内部规范 1.4.2【强制】 — 当使用缩写或特殊术语时，全大写或者全小写保持和导出规则一致；导出标识符中的缩写词必须全部大写。

本次变更引入了 5 处导出标识符使用了 `Ip`（应为 `IP`），需批量修正：

#### E1.1 — `model/token.go:17` 类型名

```go
// 当前（违规）
type IpPolicy struct {

// 修正
type IPPolicy struct {
```

#### E1.2 — `model/token.go:67` 结构体字段与类型引用

```go
// 当前（违规）
IpPolicy           *IpPolicy      `json:"ip_policy,omitempty" gorm:"type:text;column:ip_policy"`

// 修正
IPPolicy           *IPPolicy      `json:"ip_policy,omitempty" gorm:"type:text;column:ip_policy"`
```

#### E1.3 — `model/token.go:360` 方法名

```go
// 当前（违规）
func (token *Token) UpdateIpPolicy() (err error) {

// 修正
func (token *Token) UpdateIPPolicy() (err error) {
```

#### E1.4 — `controller/token.go:295` 类型名

```go
// 当前（违规）
type IpPolicyRequest struct {

// 修正
type IPPolicyRequest struct {
```

#### E1.5 — `controller/token.go:301` 函数名

```go
// 当前（违规）
func UpdateTokenIpPolicy(c *gin.Context) {

// 修正
func UpdateTokenIPPolicy(c *gin.Context) {
```

> **批量修复说明**: 以上 5 处同源问题需同步修改。`middleware/auth.go` 中对 `token.IpPolicy` 的引用、`router/api-router.go` 中对 `controller.UpdateTokenIpPolicy` 的引用也需一并更新。

---

### E2 · 命名规范：未导出标识符中的缩写词应全小写

**规则**: 内部规范 1.4.2【强制】

**位置**: `common/ip_matcher.go:14`

```go
// 当前（违规）
var trustedProxyCIDRs []*net.IPNet

// 修正
var trustedProxyCidrs []*net.IPNet
```

> 变量名为未导出（小写开头），其中 `CIDR` 为缩写词，应与导出规则一致使用全小写 → `Cidrs`。
> 同时 `isTrustedProxy`、`GetClientIP` 函数内部对 `trustedProxyCIDRs` 的所有引用也需同步更新。

---

### E3 · 魔数：直接使用字面量数值

**规则**: 内部规范 7.3.1【强制】 — 除了 0 和 1，不要使用魔法数字。

**位置**: `controller/token.go:319`

```go
// 当前（违规）
if len(req.Ips) > 100 {
    c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "too many IP entries, max 100"})
    return
}

// 修正
const maxIPPolicyEntries = 100
// ...
if len(req.Ips) > maxIPPolicyEntries {
    c.JSON(http.StatusBadRequest, gin.H{
        "success": false,
        "message": fmt.Sprintf("too many IP entries, max %d", maxIPPolicyEntries),
    })
    return
}
```

---

## ⚠️ 警告问题（Warning）

### W1 · 命名规范：未导出全局变量建议使用 `_` 前缀

**规则**: 内部规范 1.4.8【建议】 — 对于未导出的全局变量/常量，应该使用 `_` 前缀，避免在同一个包中的其他文件中意外使用错误的值。

**位置**: `common/ip_matcher.go:14`

```go
// 当前
var trustedProxyCIDRs []*net.IPNet

// 建议修正（结合 E2 的缩写词修正）
var _trustedProxyCidrs []*net.IPNet
```

---

## 🔒 安全检查

未发现安全问题。

- X-Forwarded-For 仅从受信代理（`trustedProxyCIDRs`）读取，设计安全
- IP/CIDR 输入通过 `net.ParseCIDR` / `net.ParseIP` 严格校验，无注入风险
- `UpdateTokenIpPolicy` 已校验用户所有权与管理员身份
- 数据库操作通过 GORM ORM 执行，无 SQL 注入风险
- 无硬编码凭据

---

## 📋 检测过程追溯

| 问题 ID | 文件 | 行号 | 检测依据 | 检测步骤 |
|---------|------|------|---------|---------|
| E1.1–E1.5 | model/token.go, controller/token.go | 17, 67, 360, 295, 301 | 内部规范 1.4.2【强制】 | AI 规范检查：遍历新增代码，识别导出标识符中 `Ip` 缩写词未使用全大写 `IP` |
| E2 | common/ip_matcher.go | 14 | 内部规范 1.4.2【强制】 | AI 规范检查：新增文件全文扫描，`trustedProxyCIDRs` 为未导出变量，`CIDR` 缩写应全小写 |
| E3 | controller/token.go | 319 | 内部规范 7.3.1【强制】 | AI 规范检查：新增代码行检测到字面量 `100`，且未在同一语句中说明其含义 |
| W1 | common/ip_matcher.go | 14 | 内部规范 1.4.8【建议】 | AI 规范检查：新增文件扫描，未导出全局变量缺少 `_` 前缀 |

---

## 🔧 整改建议

**优先处理（error 级别）**:

1. **批量重命名（E1）**: 将所有 `IpPolicy` → `IPPolicy`、`UpdateIpPolicy` → `UpdateIPPolicy`、`IpPolicyRequest` → `IPPolicyRequest`、`UpdateTokenIpPolicy` → `UpdateTokenIPPolicy`，同时更新所有调用处（middleware、router）。
2. **变量重命名（E2）**: `trustedProxyCIDRs` → `trustedProxyCidrs`，同步更新 `ip_matcher.go` 内所有引用。
3. **提取常量（E3）**: 在 `controller/token.go` 中新增 `const maxIPPolicyEntries = 100`，替换魔数。

**建议处理（warning 级别）**:

4. **加下划线前缀（W1）**: 结合 E2 的修正，将变量改为 `_trustedProxyCidrs`。
