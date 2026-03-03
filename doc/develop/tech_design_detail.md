# API Key IP 访问控制模块 - 详细设计文档

> 本文档使用 RFC 2119 规范语言：必须（SHALL）、应该（SHOULD）、可以（MAY）
>
> 源码分析基准版本：`21cfc1ca`
> 输入文档：`doc/design/tech_design.md`（概要设计 V1.0）、`doc/testcase/testcase.md`

---

## 1. 概述

### 1.1 为什么要做

现有 API Key 的 IP 限制（`AllowIps` 字段）仅支持精确 IP 白名单、缺乏 CIDR 网段支持和黑名单模式，同时在反向代理架构下因 `c.ClientIP()` 行为不可控导致 XFF 头可被伪造绕过，无法满足企业级 IP 访问控制需求。

### 1.2 变更内容

- **新增文件** `common/ip_matcher.go`：CIDR 工具函数库（ValidateCIDRList、ParseCIDRList、IPMatchesCIDRList、GetClientIP、InitTrustedProxies）
- **修改文件** `model/token.go`：Token struct 新增 `IpPolicy *IpPolicy` 字段；新增 `IpPolicy` 结构体及 Valuer/Scanner 实现；新增 `UpdateIpPolicy()` 方法
- **修改文件** `controller/token.go`：新增 `UpdateTokenIpPolicy` HTTP Handler；新增 `IpPolicyRequest` DTO
- **修改文件** `router/api-router.go`：在 tokenRoute 组注册 `PUT("/:id/ip_policy", ...)`
- **修改文件** `middleware/auth.go`：在 `TokenAuth()` 中插入 F004（GetClientIP）和 F005（IpPolicy 策略执行）代码块
- **修改文件** `common/init.go`：`InitEnv()` 末尾追加 `InitTrustedProxies(os.Getenv("TRUSTED_PROXIES"))` 调用
- **修改文件** `types/error.go`：新增常量 `ErrorCodeIpNotAllowed`
- **新增文件** `common/ip_matcher_test.go`：CIDR 工具函数单元测试（覆盖率目标 ≥ 85%）
- 无破坏性变更：`AllowIps` 字段和 `GetIpLimits()` 方法保持不变

### 1.3 影响范围

- **受影响的功能模块**: Token 鉴权流程（TokenAuth 中间件）、Token 管理 API（controller/token.go）
- **受影响的代码文件**:
  - `common/ip_matcher.go`（新建）
  - `common/ip_matcher_test.go`（新建）
  - `model/token.go`（修改）
  - `controller/token.go`（修改）
  - `router/api-router.go`（修改）
  - `middleware/auth.go`（修改）
  - `common/init.go`（修改）
  - `types/error.go`（修改）

---

## 2. 目标与非目标

### 2.1 目标

1. **[G1] CIDR 工具库**：新建 `common/ip_matcher.go`，提供 `ValidateCIDRList`、`ParseCIDRList`、`IPMatchesCIDRList`、`GetClientIP`、`InitTrustedProxies` 五个函数，纯标准库实现，并发安全，单元测试覆盖率 ≥ 85%
2. **[G2] Token.IpPolicy 数据模型**：Token struct 新增 `IpPolicy *IpPolicy` 字段，JSON 序列化存储到 `ip_policy TEXT NULL` 列；实现 `driver.Valuer`/`sql.Scanner`；GORM AutoMigrate 兼容三库（SQLite/MySQL/PostgreSQL）
3. **[G3] IP 策略配置 API**：新增 `PUT /api/token/:id/ip_policy` 接口，支持配置/清除白名单/黑名单策略，权限双重校验（拥有者 OR 管理员），返回 `{"success":true/false,"message":"..."}`
4. **[G4] 可信代理感知 IP 提取**：在 TokenAuth 中间件中使用 `GetClientIP(c)` 替代 `c.ClientIP()`，支持 TRUSTED_PROXIES 环境变量配置，防 XFF 伪造
5. **[G5] IP 策略执行**：在 TokenAuth 中间件 AllowIps 校验块之后插入 IpPolicy 策略校验块，支持白名单/黑名单双模式，fail-open 降级，拒绝时返回 `IP_NOT_ALLOWED`

### 2.2 非目标

明确不在本次范围内的内容：

- GET 接口查询 IpPolicy（产品 Q3 待确认，预留扩展点）
- Redis 缓存主动失效（R6 遗留问题，IpPolicy 更新后缓存依赖自然过期，最长延迟约 1 分钟）
- TRUSTED_PROXIES 运行时动态更新（需重启服务生效）
- Token 级 CIDR 解析结果缓存（当前每次请求解析，后续可扩展）
- AllowIps 旧字段的任何修改

---

## 3. 技术决策

### 3.1 决策1: CIDR 工具库的包位置与函数设计 [G1]

**决策内容**: 新建 `common/ip_matcher.go`（与现有 `common/ip.go` 同包），导出 5 个函数，`trustedProxyCIDRs` 为包级私有全局变量。

**备选方案**:

| 方案 | 优点 | 缺点 | 结论 |
|------|------|------|------|
| A: 新增 `common/ip_matcher.go` | 与现有 `common/ip.go` 风格一致；中间件可直接调用 `common.GetClientIP(c)` | 包级全局变量需要注意初始化顺序 | **采纳** |
| B: 合并到现有 `common/ip.go` | 文件更少 | 现有 ip.go 已有功能，混在一起降低内聚性 | 放弃 |
| C: 新建 `service/ip_policy.go` | 层次更清晰 | middleware 调用 service 会增加依赖；service 层通常不操作 gin.Context | 放弃 |

**选择理由**: 工具函数（无业务逻辑、无 IO）归入 `common/` 是本项目约定（参见 `doc/kb/技术知识库/代码编写指南.md:32`），且 `common/ip.go:33` 的 `IsIpInCIDRList` 函数为同类设计，风格一致性优先。

---

### 3.2 决策2: IpPolicy 持久化方案 [G2]

**决策内容**: 在 `tokens` 表新增 `ip_policy TEXT NULL` 列，通过 `driver.Valuer`/`sql.Scanner` 将 `IpPolicy` struct 序列化为 JSON 字符串存储。

**备选方案**:

| 方案 | 优点 | 缺点 | 结论 |
|------|------|------|------|
| A: JSON TEXT 单列（`driver.Valuer`/`Scanner`）| 无需新表；GORM AutoMigrate 仅 ADD COLUMN；三库兼容；向后兼容（旧记录 NULL）| 无法在 DB 层对单个 IP 条目做索引查询 | **采纳** |
| B: 独立 `token_ip_policies` 关联表 | 支持结构化查询（如"哪些 token 包含某 IP"）| 引入新表、外键、JOIN 查询，复杂度高；不符合本需求（无跨 Token 查询需求）| 放弃 |
| C: 直接扩展 AllowIps 字段语义 | 零新增字段 | 破坏现有 AllowIps 的精确 IP 语义；向后不兼容 | 放弃 |

**选择理由**: 本需求无"按 IP 查询哪些 Token"的场景，单列 JSON 存储简单可靠，与项目中 `ModelLimits string` 字段（GORM 类型 varchar，JSON 序列化）的模式一致。

---

### 3.3 决策3: HTTP Handler 权限校验方案 [G3]

**决策内容**: 使用 `model.GetTokenById(id)` 按 ID 查询 Token（不限制 userId），然后在 Handler 内手动校验 `token.UserId == currentUserId || model.IsAdmin(currentUserId)`。

**源码分析发现的关键约束**:

> `model.GetTokenByIds(id, userId)` 使用 `WHERE id = ? AND user_id = ?` 联合查询（`model/token.go:222`），管理员操作他人 Token 时会返回 `gorm.ErrRecordNotFound`（404），无法正常使用。
>
> 必须使用 `model.GetTokenById(id)` + `model.IsAdmin(userId)` 组合实现权限校验。

**备选方案**:

| 方案 | 优点 | 缺点 | 结论 |
|------|------|------|------|
| A: `GetTokenById` + `IsAdmin` 手动校验 | 支持管理员操作他人 Token；权限逻辑清晰 | 需要额外 `IsAdmin` 数据库查询 | **采纳** |
| B: `GetTokenByIds(id, userId)` | 利用已有函数简化代码 | 管理员使用 userId 查会 404，功能缺陷 | 放弃 |

**选择理由**: 概要设计明确要求"权限校验（拥有者/管理员）"（`tech_design.md:1.1`），方案 B 会导致 TC007（管理员操作他人 Token）测试失败。

---

### 3.4 决策4: UpdateIpPolicy Model 方法的 NULL 更新策略 [G2, G3]

**决策内容**: 新增专用 `UpdateIpPolicy()` 方法，使用 `DB.Model(&Token{}).Where("id = ?", token.Id).Update("ip_policy", rawValue)` 直接传递值（JSON 字符串或 nil）实现精确更新，避免 GORM 跳过零值的问题。

**GORM 零值问题说明**:

```
// GORM的Updates会跳过零值（nil指针是零值）
// DB.Model(token).Select("ip_policy").Updates(token)
// 当 token.IpPolicy == nil 时，GORM 不会写 NULL，而是跳过该字段
```

**备选方案**:

| 方案 | 优点 | 缺点 | 结论 |
|------|------|------|------|
| A: `Update("ip_policy", rawValue)` 精确传值 | 正确处理 NULL；代码直观 | 需要 Handler 预先序列化 JSON 或在 Model 层处理 | **采纳** |
| B: `Select("ip_policy").Updates(token)` | 利用 Valuer 接口自动序列化 | 当 IpPolicy 为 nil 时 GORM 跳过零值，无法写 NULL | 放弃 |
| C: `Save(token)` 全字段更新 | 能写 NULL | 覆盖所有字段，有误写风险；绕过 select 保护 | 放弃 |

**选择理由**: 精确指定字段、正确处理 NULL 是 `ip_policy` 字段更新的核心要求；方案 A 与项目中 `increaseTokenQuota` 等直接 `Update` 的模式一致（`model/token.go:382`）。

---

### 3.5 决策5: F005 IP 策略执行中的错误响应格式 [G5]

**决策内容**: 在 `TokenAuth` 中间件中，IP 策略拒绝请求时调用 `abortWithOpenAiMessage(c, http.StatusForbidden, "IP_NOT_ALLOWED", types.ErrorCodeIpNotAllowed)` 响应，与现有 AllowIps 拒绝格式一致。

**源码分析发现**:

> `middleware/utils.go:12` 的 `abortWithOpenAiMessage` 是 TokenAuth 中间件内唯一的响应函数，返回 OpenAI 兼容格式（`{"error": {"message":"...", "type":"new_api_error", "code":"..."}}`）。这是因为 TokenAuth 用于中继路由（relay），客户端期望 OpenAI 格式的错误响应。

**备选方案**:

| 方案 | 优点 | 缺点 | 结论 |
|------|------|------|------|
| A: `abortWithOpenAiMessage` (OpenAI 格式) | 与现有 AllowIps 拒绝响应一致；中继客户端可正确解析 | 非标准 REST 管理 API 格式，但此处是中间件 | **采纳** |
| B: `c.JSON(http.StatusForbidden, gin.H{"success":false,...})` | REST 管理格式 | 与现有 AllowIps 校验（行326）格式不一致；中继客户端无法解析 | 放弃 |

**选择理由**: TokenAuth 中间件保护的是中继路由，不是管理 API 路由，与 AllowIps 校验格式保持一致是正确的。

---

## 4. 整体架构

```
┌─────────────────────────────────────────────────────────────┐
│                    请求生命周期                               │
│                                                             │
│  服务启动:                                                   │
│  common.InitEnv() → common.InitTrustedProxies(             │
│                       os.Getenv("TRUSTED_PROXIES"))         │
│       ↓ 解析后存入 trustedProxyCIDRs []*net.IPNet            │
│                                                             │
│  请求到达 (AI 中继路由):                                     │
│  Router → [UserAuth → TokenAuth → ...] → Handler           │
│                                                             │
│  TokenAuth() 执行链:                                         │
│  ① WebSocket/Anthropic/Gemini key 统一化                    │
│  ② ValidateUserToken(key) → 验证 Token 状态                  │
│  ③ [现有] AllowIps 精确 IP 白名单校验                        │
│  ④ [新增F004] common.GetClientIP(c) → c.Set("client_ip")   │
│  ⑤ [新增F005] IpPolicy 策略校验 (白名单/黑名单)              │
│     ↓ 拒绝: abortWithOpenAiMessage 403 IP_NOT_ALLOWED       │
│  ⑥ GetUserCache → 用户状态/组权限校验                        │
│  ⑦ SetupContextForToken → 设置上下文                        │
│                                                             │
│  请求到达 (管理 API):                                        │
│  PUT /api/token/:id/ip_policy                               │
│  → [UserAuth] → UpdateTokenIpPolicy                         │
│     ① strconv.Atoi(id)                                     │
│     ② ShouldBindJSON → IpPolicyRequest{Mode, Ips}          │
│     ③ ValidateMode(mode)                                    │
│     ④ ValidateCIDRList(ips) (len<=100 先检查)               │
│     ⑤ GetTokenById(id)                                     │
│     ⑥ 权限校验: UserId==currentUser || IsAdmin              │
│     ⑦ 构建 IpPolicy / nil                                   │
│     ⑧ token.UpdateIpPolicy()                               │
│     ↓ 成功: {"success":true,"message":""}                   │
└─────────────────────────────────────────────────────────────┘

┌──────────────────── 模块依赖关系 ────────────────────────────┐
│  middleware/auth.go                                          │
│       ↓ calls                                               │
│  common/ip_matcher.go ←──────── common/init.go             │
│       ↑ uses                    (InitTrustedProxies)        │
│  model/token.go                                             │
│       ↑ uses                                               │
│  controller/token.go                                        │
│       ↑ registers                                           │
│  router/api-router.go                                       │
└─────────────────────────────────────────────────────────────┘
```

---

## 5. 模块设计

### 5.1 目录结构

```
new-api/
├── common/
│   ├── ip.go                    # 现有，保持不变
│   ├── ip_matcher.go            # 【新建】CIDR 工具函数库
│   └── ip_matcher_test.go       # 【新建】单元测试（覆盖率 ≥ 85%）
├── controller/
│   └── token.go                 # 【修改】新增 UpdateTokenIpPolicy Handler 和 IpPolicyRequest DTO
├── middleware/
│   └── auth.go                  # 【修改】在 TokenAuth() 插入 F004/F005 代码块
├── model/
│   └── token.go                 # 【修改】新增 IpPolicy struct、Token.IpPolicy 字段、UpdateIpPolicy 方法
├── router/
│   └── api-router.go            # 【修改】注册 PUT /:id/ip_policy 路由
├── types/
│   └── error.go                 # 【修改】新增 ErrorCodeIpNotAllowed 常量
└── common/
    └── init.go                  # 【修改】InitEnv() 末尾追加 InitTrustedProxies 调用
```

### 5.2 核心流程

#### 5.2.1 F001+F002+F003: 更新 IP 策略（配置接口）

```
UpdateTokenIpPolicy(c *gin.Context)
    │
    ├─[1] 解析路径参数
    │      id, err := strconv.Atoi(c.Param("id"))
    │      if err != nil → c.JSON(400, ApiError)  return
    │
    ├─[2] 绑定请求体
    │      var req IpPolicyRequest
    │      if err := c.ShouldBindJSON(&req); err != nil → ApiError(c, err)  return
    │
    ├─[3] 校验 mode
    │      if mode != "" && mode != "whitelist" && mode != "blacklist"
    │          → c.JSON(400, {"success":false,"message":"invalid mode: "+mode})  return
    │
    ├─[4] 校验 ips 条目数
    │      if len(req.Ips) > 100
    │          → c.JSON(400, {"success":false,"message":"too many IP entries, max 100"})  return
    │
    ├─[5] 校验 ips 格式（仅 mode 非空时执行）
    │      if mode != "" && len(req.Ips) > 0:
    │          if err := common.ValidateCIDRList(req.Ips); err != nil
    │              → c.JSON(400, {"success":false,"message":"invalid IP/CIDR: "+err.Error()})  return
    │
    ├─[6] 查询 Token
    │      token, err := model.GetTokenById(id)
    │      if err != nil (含 gorm.ErrRecordNotFound)
    │          → c.JSON(404, {"success":false,"message":"not found"})  return
    │
    ├─[7] 权限校验
    │      currentUserId := c.GetInt("id")
    │      if token.UserId != currentUserId && !model.IsAdmin(currentUserId)
    │          → c.JSON(403, {"success":false,"message":"forbidden"})  return
    │
    ├─[8] 构建 IpPolicy
    │      if mode == ""  → token.IpPolicy = nil   // 清除策略
    │      else           → token.IpPolicy = &model.IpPolicy{Mode: mode, Ips: req.Ips}
    │
    ├─[9] 持久化
    │      if err := token.UpdateIpPolicy(); err != nil
    │          → c.JSON(500, {"success":false,"message":"internal server error"})  return
    │
    └─[10] 成功响应
           c.JSON(200, {"success":true,"message":""})
```

#### 5.2.2 F004+F005: IP 策略执行（TokenAuth 中间件增量）

```
TokenAuth() — 在 AllowIps 校验块（行316-330）之后插入：

    // [F004] 获取真实客户端 IP（可信代理感知）
    clientIP := common.GetClientIP(c)
    c.Set("client_ip", clientIP)

    // [F005] IpPolicy 策略校验
    if token.IpPolicy != nil && token.IpPolicy.Mode != "" {
        cidrs, err := common.ParseCIDRList(token.IpPolicy.Ips)
        if err != nil {
            // fail-open: CIDR 解析失败时不拦截，记录 WARN 日志
            logger.LogWarn(c.Request.Context(), fmt.Sprintf(
                "ip_policy_parse_failed: token_id=%d err=%s", token.Id, err.Error()))
            // 继续执行，不 return
        } else {
            hit := common.IPMatchesCIDRList(clientIP, cidrs)
            blocked := (token.IpPolicy.Mode == "whitelist" && !hit) ||
                       (token.IpPolicy.Mode == "blacklist" && hit)
            if blocked {
                abortWithOpenAiMessage(c, http.StatusForbidden,
                    "IP_NOT_ALLOWED", types.ErrorCodeIpNotAllowed)
                return
            }
        }
    }
```

#### 5.2.3 GetClientIP 决策逻辑

```
GetClientIP(c *gin.Context) string:
    remoteAddr := c.Request.RemoteAddr
    host, _, err := net.SplitHostPort(remoteAddr)
    if err != nil → host = remoteAddr   // 无端口情况

    remoteIP := net.ParseIP(host)
    if remoteIP == nil || !isTrustedProxy(remoteIP):
        return host   // 不可信或无法解析 → 直接返回 RemoteAddr

    // RemoteAddr 属于可信代理 → 采信 XFF 最左侧 IP
    xff := c.Request.Header.Get("X-Forwarded-For")
    if xff == ""  → return host   // XFF 为空 → 回退 RemoteAddr
    parts := strings.Split(xff, ",")
    leftmost := strings.TrimSpace(parts[0])
    if leftmost == "" → return host
    return leftmost
```

### 5.3 接口定义

#### 5.3.1 `common/ip_matcher.go` — 函数签名

```go
package common

import (
    "net"
    "strings"
    "github.com/gin-gonic/gin"
)

// trustedProxyCIDRs 包级私有全局变量，服务启动时初始化一次，运行时只读（并发安全）
var trustedProxyCIDRs []*net.IPNet

// InitTrustedProxies 解析逗号分隔的 IP/CIDR 字符串，初始化 trustedProxyCIDRs。
// 非法条目跳过并记录警告日志。
// 调用时机: common.InitEnv() 末尾，服务启动时执行一次。
func InitTrustedProxies(proxies string)

// ValidateCIDRList 逐条校验 IP/CIDR 格式合法性。
// 返回第一个非法条目的 error（包含非法条目字符串）；全部合法或空列表返回 nil。
func ValidateCIDRList(ips []string) error

// ParseCIDRList 将 IP/CIDR 字符串列表解析为 []*net.IPNet。
// 单个 IP（无前缀）自动补全为 /32（IPv4）或 /128（IPv6）。
// 每次调用返回独立切片（无共享状态），并发安全。
// 返回 error 时调用方 SHALL 执行 fail-open 降级（不拦截请求）。
func ParseCIDRList(ips []string) ([]*net.IPNet, error)

// IPMatchesCIDRList 判断给定 IP 字符串是否命中 CIDR 列表中任意一条。
// 无效 IP 字符串或 nil/空 cidrs 均返回 false。
func IPMatchesCIDRList(ipStr string, cidrs []*net.IPNet) bool

// GetClientIP 从 gin.Context 中提取真实客户端 IP。
// 若 RemoteAddr 属于 trustedProxyCIDRs，则采信 X-Forwarded-For 头的最左侧 IP；否则使用 RemoteAddr。
// 永远返回纯 IP 字符串（不含端口）。
func GetClientIP(c *gin.Context) string

// isTrustedProxy 判断给定 IP 是否在 trustedProxyCIDRs 中（包内私有）。
func isTrustedProxy(ip net.IP) bool
```

#### 5.3.2 `model/token.go` — 新增数据结构

```go
// IpPolicy 定义 API Key 的 IP 访问策略。
// 通过 driver.Valuer/sql.Scanner 实现 JSON 序列化，存储在 tokens.ip_policy TEXT 列。
type IpPolicy struct {
    Mode string   `json:"mode"` // "whitelist" | "blacklist"
    Ips  []string `json:"ips"`  // IP 或 CIDR 列表
}

// Value 实现 driver.Valuer，将 IpPolicy 序列化为 JSON 字符串写入数据库。
// nil 指针时返回 (nil, nil) → 数据库写入 NULL。
func (p *IpPolicy) Value() (driver.Value, error)

// Scan 实现 sql.Scanner，从数据库 JSON 字符串反序列化为 IpPolicy。
// 兼容 []byte 和 string 两种驱动返回类型（跨数据库驱动差异）。
func (p *IpPolicy) Scan(value interface{}) error

// UpdateIpPolicy 仅更新 Token 的 ip_policy 字段（精确单字段更新）。
// 同步删除 Redis 缓存（强制下次请求从 DB 读取最新策略）。
func (token *Token) UpdateIpPolicy() (err error)
```

**Token struct 新增字段**（位于 `AllowIps` 字段之后）：

```go
type Token struct {
    // ... 现有字段保持不变 ...
    AllowIps  *string   `json:"allow_ips" gorm:"default:''"`              // 旧字段，不变
    IpPolicy  *IpPolicy `json:"ip_policy,omitempty" gorm:"type:text;column:ip_policy"` // 新增
    // ...
}
```

#### 5.3.3 `controller/token.go` — DTO 和 Handler 签名

```go
// IpPolicyRequest 是 PUT /api/token/:id/ip_policy 的请求体绑定 DTO。
type IpPolicyRequest struct {
    Mode string   `json:"mode"` // 策略模式：whitelist | blacklist | ""（清除）
    Ips  []string `json:"ips"`  // CIDR 列表；mode=="" 时忽略
}

// UpdateTokenIpPolicy 为指定 Token 设置或清除 IP 访问策略。
func UpdateTokenIpPolicy(c *gin.Context)
```

---

## 6. 需求规格

### 6.1 需求1: CIDR 工具函数库 [F001/F003] [G1]

系统**必须**在 `common/ip_matcher.go` 中实现以下函数，全部使用 Go 标准库 `net` 包，不引入外部依赖。

#### 接口定义

- **文件**: `common/ip_matcher.go`
- **包**: `common`

#### 函数规格

**ValidateCIDRList(ips []string) error**

| 条件 | 行为 |
|------|------|
| `ips` 为 nil 或空 | 返回 nil（视为合法） |
| 所有条目均为合法 IPv4/IPv6 精确地址或 CIDR | 返回 nil |
| 存在非法条目（格式错误、IP 段超范围等）| 返回 `fmt.Errorf("%s", invalidEntry)`（首个非法条目）|

实现伪代码：

```go
func ValidateCIDRList(ips []string) error {
    for _, ip := range ips {
        if strings.Contains(ip, "/") {
            if _, _, err := net.ParseCIDR(ip); err != nil {
                return fmt.Errorf("%s", ip)
            }
        } else {
            if net.ParseIP(ip) == nil {
                return fmt.Errorf("%s", ip)
            }
        }
    }
    return nil
}
```

**ParseCIDRList(ips []string) ([]*net.IPNet, error)**

| 条件 | 行为 |
|------|------|
| `ips` 为 nil 或空 | 返回空 `[]*net.IPNet{}`，nil |
| 单个精确 IP（无 `/` 前缀）| 自动补全：IPv4 → `/32`，IPv6 → `/128` |
| 合法 CIDR | 解析并追加到结果切片 |
| 存在解析失败条目 | 返回 nil，error（调用方 **必须** fail-open） |
| 成功 | 返回独立 `[]*net.IPNet` 切片（每次调用新分配，无共享状态）|

实现伪代码：

```go
func ParseCIDRList(ips []string) ([]*net.IPNet, error) {
    result := make([]*net.IPNet, 0, len(ips))
    for _, ipStr := range ips {
        if !strings.Contains(ipStr, "/") {
            // 补全 CIDR 前缀
            if net.ParseIP(ipStr).To4() != nil {
                ipStr = ipStr + "/32"
            } else {
                ipStr = ipStr + "/128"
            }
        }
        _, network, err := net.ParseCIDR(ipStr)
        if err != nil {
            return nil, err
        }
        result = append(result, network)
    }
    return result, nil
}
```

**IPMatchesCIDRList(ipStr string, cidrs []*net.IPNet) bool**

| 条件 | 行为 |
|------|------|
| `ipStr` 为空或不合法 IP | 返回 false |
| `cidrs` 为 nil 或空 | 返回 false |
| `ipStr` 匹配 cidrs 中任意一条 | 返回 true |
| 均不匹配 | 返回 false |

**GetClientIP(c *gin.Context) string**

| RemoteAddr 是否可信 | XFF 头 | 返回值 |
|---------------------|--------|--------|
| 否（trustedProxyCIDRs 为 nil 或不匹配）| 任意 | RemoteAddr 纯 IP |
| 是 | 非空 | XFF 最左侧 IP（TrimSpace 后）|
| 是 | 空 | RemoteAddr 纯 IP（回退）|

**InitTrustedProxies(proxies string)**

| 输入 | 行为 |
|------|------|
| 空字符串 | `trustedProxyCIDRs = nil` |
| 逗号分隔的合法 IP/CIDR 列表 | 解析并存入 `trustedProxyCIDRs` |
| 含非法条目 | 跳过非法条目，`log.Printf("[WARN] invalid trusted proxy entry: %s", entry)` |

#### 场景: F003 成功校验

- **前置条件**: 服务已启动，ip_matcher.go 已存在
- **当**: 调用 `ValidateCIDRList(["1.2.3.4", "10.0.0.0/8"])`
- **则**: 返回 nil；调用 `ParseCIDRList(["1.2.3.4", "10.0.0.0/8"])` 返回 2 个 `*net.IPNet`，nil

#### 场景: F003 非法 CIDR 校验

- **前置条件**: 同上
- **当**: 调用 `ValidateCIDRList(["1.2.3.4", "999.0.0.1/33"])`
- **则**: 返回包含 `"999.0.0.1/33"` 的 error

---

### 6.2 需求2: Token.IpPolicy 数据模型 [F002] [G2]

系统**必须**在 `model/token.go` 中实现 IpPolicy 结构体的序列化和持久化。

#### 接口定义

- **文件**: `model/token.go`
- **数据库**: GORM AutoMigrate 在 `tokens` 表执行 `ADD COLUMN ip_policy TEXT NULL`

#### 数据结构

```go
type IpPolicy struct {
    Mode string   `json:"mode"`
    Ips  []string `json:"ips"`
}
```

#### Valuer 实现规格

```go
// Value 必须实现 driver.Valuer 接口（指针接收者，处理 nil 情况）
func (p *IpPolicy) Value() (driver.Value, error) {
    if p == nil {
        return nil, nil  // 返回 SQL NULL
    }
    bytes, err := common.Marshal(p)  // 必须使用 common.Marshal（Rule 1）
    if err != nil {
        return nil, err
    }
    return string(bytes), nil
}
```

#### Scanner 实现规格

```go
// Scan 必须实现 sql.Scanner 接口（指针接收者）
func (p *IpPolicy) Scan(value interface{}) error {
    if value == nil {
        return nil  // NULL 列 → IpPolicy 保持零值（不影响上层 nil 指针）
    }
    var bytes []byte
    switch v := value.(type) {
    case []byte:
        bytes = v
    case string:
        bytes = []byte(v)
    default:
        return fmt.Errorf("IpPolicy.Scan: unsupported type %T", value)
    }
    return common.Unmarshal(bytes, p)  // 必须使用 common.Unmarshal（Rule 1）
}
```

#### UpdateIpPolicy 方法规格

```go
// UpdateIpPolicy 仅更新 ip_policy 单字段（精确更新，正确处理 NULL）。
// 更新成功后删除 Redis 缓存（token.Key 必须已加载，由 Handler 在 GetTokenById 后调用）。
func (token *Token) UpdateIpPolicy() (err error) {
    defer func() {
        // Redis 缓存失效：删除缓存，强制下次鉴权从 DB 读取最新 IpPolicy
        if shouldUpdateRedis(true, err) {
            gopool.Go(func() {
                if err := cacheDeleteToken(token.Key); err != nil {
                    common.SysLog("failed to delete token cache after UpdateIpPolicy: " + err.Error())
                }
            })
        }
    }()
    // 直接传值，正确处理 NULL（nil IpPolicy → SQL NULL）
    // token.IpPolicy 实现了 driver.Valuer，nil 时 Value() 返回 (nil, nil)
    return DB.Model(&Token{}).Where("id = ?", token.Id).
        Update("ip_policy", token.IpPolicy).Error
}
```

> **注意**：由于 `driver.Valuer` 在 `*IpPolicy` 上实现，GORM 在调用 `Update("ip_policy", token.IpPolicy)` 时会调用 `token.IpPolicy.Value()`，nil 时返回 `(nil, nil)` 对应 SQL NULL。

#### 场景: 向后兼容（旧 Token）

- **前置条件**: 数据库中已有旧 Token 记录，`ip_policy` 列为 NULL
- **当**: GORM 查询旧 Token
- **则**: `token.IpPolicy` 为 nil；TokenAuth 中间件跳过 F005 校验（`token.IpPolicy == nil` 条件）

#### 场景: 清除策略

- **前置条件**: Token 已配置白名单策略（ip_policy 为 JSON 字符串）
- **当**: Handler 收到 `{"mode":"","ips":[]}` 请求，`token.IpPolicy = nil`，调用 `token.UpdateIpPolicy()`
- **则**: `ip_policy` 列更新为 NULL；下次鉴权时跳过 IP 策略校验

---

### 6.3 需求3: IP 策略配置 HTTP 接口 [F001] [G3]

系统**必须**提供 `PUT /api/token/:id/ip_policy` 接口，允许 Token 拥有者或管理员配置 IP 访问策略。

#### 接口定义

- **路径**: `PUT /api/token/:id/ip_policy`
- **Content-Type**: `application/json`
- **认证**: UserAuth（JWT Bearer Token，由路由组中间件处理）

#### 路由注册

在 `router/api-router.go:251` 的 `tokenRoute` 组中追加：

```go
tokenRoute.PUT("/:id/ip_policy", controller.UpdateTokenIpPolicy)
```

#### 请求参数

> **重要**：必须明确每个参数的必填/非必填状态

| 参数名 | 类型 | 位置 | 必填 | 约束 | 说明 |
|--------|------|------|------|------|------|
| `id` | integer | path | 是 | > 0，可解析为 int | 目标 Token 的数据库 ID |
| `mode` | string | body | 是 | `whitelist` \| `blacklist` \| `""` | 策略模式；`""` 清除策略 |
| `ips` | []string | body | 是（可为空数组）| len ≤ 100；每条为合法 IP 或 CIDR | CIDR 列表；mode=="" 时忽略 |

#### 响应参数

| 参数名 | 类型 | 说明 |
|--------|------|------|
| `success` | bool | 操作是否成功 |
| `message` | string | 错误信息（成功时为空字符串） |

#### 业务逻辑说明

> **重要**：业务逻辑判断必须明确说明数据来源和计算规则

**[mode 校验]** 的判断**必须**基于以下规则：
- 数据来源：请求体 JSON `mode` 字段
- 计算规则：`mode != "" && mode != "whitelist" && mode != "blacklist"` → 返回 400

**[ips 条目数校验]** 的判断**必须**基于以下规则：
- 数据来源：请求体 JSON `ips` 字段，`len(req.Ips)`
- 计算规则：`len(req.Ips) > 100` → 返回 400；注意：条目数校验在 ValidateCIDRList 之前执行

**[CIDR 格式校验]** 的判断**必须**基于以下规则：
- 仅在 `mode != ""` 且 `len(req.Ips) > 0` 时执行
- 调用 `common.ValidateCIDRList(req.Ips)`，返回 error 时携带具体非法条目

**[权限校验]** 的判断**必须**基于以下规则：
- 数据来源：`c.GetInt("id")` 获取当前 JWT 用户 ID；`model.GetTokenById(id)` 查询 Token 拥有者
- 计算规则：`token.UserId != currentUserId && !model.IsAdmin(currentUserId)` → 返回 403
- `model.IsAdmin(userId)` 定义在 `model/user.go:715`，通过 `user.Role >= common.RoleAdminUser` 判断

#### 场景: 成功设置白名单策略

- **前置条件**: Token ID=1 存在，当前用户为该 Token 拥有者
- **当**: `PUT /api/token/1/ip_policy`，`{"mode":"whitelist","ips":["1.2.3.4","10.0.0.0/8"]}`
- **则**: HTTP 200，`{"success":true,"message":""}`；DB `ip_policy` = `{"mode":"whitelist","ips":["1.2.3.4","10.0.0.0/8"]}`；Redis 缓存被删除

#### 场景: 清除策略

- **前置条件**: Token 已配置策略
- **当**: `PUT /api/token/1/ip_policy`，`{"mode":"","ips":[]}`
- **则**: HTTP 200；DB `ip_policy` = NULL；Token 鉴权恢复为无策略模式

#### 场景: 管理员操作他人 Token

- **前置条件**: Token ID=2 属于 user_b；当前用户为 admin
- **当**: admin 请求 `PUT /api/token/2/ip_policy`，`{"mode":"blacklist","ips":["5.5.5.5"]}`
- **则**: HTTP 200；`model.IsAdmin(adminUserId)` 返回 true；权限校验通过

#### 场景: 非拥有者非管理员

- **前置条件**: Token ID=2 属于 user_b；当前用户为 user_a（普通用户）
- **当**: user_a 请求 `PUT /api/token/2/ip_policy`
- **则**: HTTP 403，`{"success":false,"message":"forbidden"}`

#### 场景: 非法 CIDR 格式

- **前置条件**: 当前用户为 Token 拥有者
- **当**: `PUT /api/token/1/ip_policy`，`{"mode":"whitelist","ips":["999.0.0.1/33"]}`
- **则**: HTTP 400，`{"success":false,"message":"invalid IP/CIDR: 999.0.0.1/33"}`；DB 不写入

---

### 6.4 需求4: 可信代理感知 IP 提取 [F004] [G4]

系统**必须**扩展 `TokenAuth()` 中间件，在 AllowIps 校验块（`middleware/auth.go:330`）之后插入 F004 代码块。

#### 插入位置和内容

```go
// [F004] 真实客户端 IP 提取（插入位置：AllowIps 校验块之后，GetUserCache 之前）
clientIP := common.GetClientIP(c)
c.Set("client_ip", clientIP)
```

同时，在 `common/init.go:InitEnv()` 函数末尾（`initConstantEnv()` 调用之后）追加：

```go
// [F004] 初始化可信代理列表
common.InitTrustedProxies(os.Getenv("TRUSTED_PROXIES"))
```

#### 场景: 无可信代理配置（默认行为）

- **前置条件**: 未设置 `TRUSTED_PROXIES` 环境变量
- **当**: 请求携带 `X-Forwarded-For: 1.2.3.4`，RemoteAddr 为 `203.0.113.5:12345`
- **则**: `clientIP = "203.0.113.5"`（XFF 头被忽略）；攻击者无法通过伪造 XFF 绕过 IP 白名单

#### 场景: 配置可信代理

- **前置条件**: `TRUSTED_PROXIES=10.0.0.1`
- **当**: RemoteAddr 为 `10.0.0.1:8080`，`X-Forwarded-For: 192.168.1.100`
- **则**: `clientIP = "192.168.1.100"`（采信 XFF 最左侧 IP）

---

### 6.5 需求5: IP 策略执行 [F005] [G5]

系统**必须**在 `TokenAuth()` 中间件 F004 代码块之后插入 F005 IP 策略校验块。

#### F005 代码块规格

```go
// [F005] IpPolicy 策略校验（插入位置：F004 代码块之后，GetUserCache 之前）
if token.IpPolicy != nil && token.IpPolicy.Mode != "" {
    cidrs, parseErr := common.ParseCIDRList(token.IpPolicy.Ips)
    if parseErr != nil {
        // fail-open: 解析失败时记录 WARN 日志，不拦截请求
        logger.LogWarn(c.Request.Context(), fmt.Sprintf(
            "ip_policy_parse_failed: token_id=%d err=%s", token.Id, parseErr.Error()))
    } else {
        clientIPForPolicy := c.GetString("client_ip")  // 读取 F004 写入的 Context 值
        hit := common.IPMatchesCIDRList(clientIPForPolicy, cidrs)
        blocked := (token.IpPolicy.Mode == "whitelist" && !hit) ||
                   (token.IpPolicy.Mode == "blacklist" && hit)
        if blocked {
            abortWithOpenAiMessage(c, http.StatusForbidden,
                "IP_NOT_ALLOWED", types.ErrorCodeIpNotAllowed)
            return  // 必须 return，防止后续逻辑执行
        }
    }
}
```

#### 新增错误码

在 `types/error.go` 的 `ErrorCode` 常量块中追加：

```go
// IP access control
ErrorCodeIpNotAllowed ErrorCode = "ip_not_allowed"
```

#### 场景: 白名单模式 — IP 在列表内（放行）

- **前置条件**: Token IpPolicy = `{Mode:"whitelist", Ips:["192.168.1.0/24"]}`；clientIP = `192.168.1.100`
- **当**: TokenAuth 执行 F005
- **则**: `hit = true`；`blocked = false`；请求继续（c.Next()）

#### 场景: 白名单模式 — IP 不在列表（拦截）

- **前置条件**: Token IpPolicy = `{Mode:"whitelist", Ips:["192.168.1.0/24"]}`；clientIP = `10.0.0.1`
- **当**: TokenAuth 执行 F005
- **则**: `hit = false`；`blocked = true`；返回 HTTP 403，错误 code = `ip_not_allowed`

#### 场景: 黑名单模式 — IP 在列表内（拦截）

- **前置条件**: Token IpPolicy = `{Mode:"blacklist", Ips:["5.5.5.5"]}`；clientIP = `5.5.5.5`
- **当**: TokenAuth 执行 F005
- **则**: `hit = true`；`blocked = true`；返回 HTTP 403，code = `ip_not_allowed`

#### 场景: CIDR 解析失败（fail-open 降级）

- **前置条件**: `ip_policy` 列存储了非法 JSON（数据异常）；Scan 成功但 Ips 包含非法条目
- **当**: `ParseCIDRList` 返回 error
- **则**: 记录 WARN 日志（`ip_policy_parse_failed`）；请求**继续**（不拦截）

#### 场景: IpPolicy 为 nil（旧 Token）

- **前置条件**: Token 无 IpPolicy（`token.IpPolicy == nil`）
- **当**: TokenAuth 执行 F005
- **则**: `if token.IpPolicy != nil` 条件不满足；直接跳过，无任何额外开销

---

### 6.6 非功能性需求

#### 性能要求

- `GetClientIP`：P99 < 0.1ms（纯内存操作，无 IO）
- `IPMatchesCIDRList`：P99 < 1ms（n ≤ 100 条 CIDR，O(n) 线性遍历）
- `ValidateCIDRList`：P99 < 10ms（配置接口，含 DB 查询总计 < 200ms）
- `ParseCIDRList`：每次调用返回独立切片，无共享状态，`go test -race` 无数据竞争

#### 可靠性要求

- **fail-open 降级**：`ParseCIDRList` 失败时**必须**记录 WARN 日志后继续请求，不拦截
- **向后兼容**：`IpPolicy == nil` 的所有现有 Token 鉴权行为与升级前完全一致
- **DB 原子性**：`UpdateIpPolicy` 失败时**必须**不留中间状态，响应 HTTP 500

---

## 7. 必须复用的现有资源

> **此章节的目的**：确保编码阶段能够复用现有资源，避免重复造轮子。

### 7.1 常量

| 常量名 | 定义位置 | 用途 |
|--------|---------|------|
| `common.RoleAdminUser` | `common/constants.go`（或类似文件）| 判断管理员角色：`user.Role >= common.RoleAdminUser` |
| `common.RedisEnabled` | `common/`（constants.go 或 redis.go）| `shouldUpdateRedis` 内部使用，已由现有函数封装 |
| `types.ErrorCodeAccessDenied` | `types/error.go:66` | 参考：AllowIps 校验拒绝使用此码；F005 使用新增的 `ErrorCodeIpNotAllowed` |

### 7.2 工具函数

| 函数名 | 定义位置 | 用途 |
|--------|---------|------|
| `common.Marshal(v any)` | `common/json.go` | **必须**用于 IpPolicy.Value() 序列化（Rule 1）|
| `common.Unmarshal(data []byte, v any)` | `common/json.go` | **必须**用于 IpPolicy.Scan() 反序列化（Rule 1）|
| `common.ApiError(c, err)` | `common/gin.go:181` | UpdateTokenIpPolicy Handler 中的 DB 错误响应 |
| `common.SysLog(msg)` | `common/`（log.go 或类似）| 缓存删除失败的日志记录（参考 token.go 模式）|
| `model.GetTokenById(id int)` | `model/token.go:226` | 按 ID 查询 Token（不限制 userId），F001 Handler 必须使用此函数 |
| `model.IsAdmin(userId int)` | `model/user.go:715` | 判断当前用户是否为管理员，F001 Handler 权限校验使用 |
| `shouldUpdateRedis(fromDB bool, err error) bool` | `model/utils.go:110` | UpdateIpPolicy 方法内部使用（与 Update() 模式一致）|
| `cacheDeleteToken(key string)` | `model/token_cache.go:21` | UpdateIpPolicy 成功后删除 Redis 缓存 |
| `gopool.Go(func())` | `github.com/bytedance/gopkg/util/gopool` | 异步执行缓存删除（与 Update() 模式一致）|
| `logger.LogWarn(ctx, msg)` | `logger/` | F005 fail-open 时记录 WARN 日志 |
| `abortWithOpenAiMessage(c, statusCode, message, code)` | `middleware/utils.go:12` | F005 IP 策略拒绝时的响应函数（在 middleware 包内可直接调用）|
| `common.GetClientIP` | — | **本次新建函数**，F004 调用 |

### 7.3 数据结构

| 结构体名 | 定义位置 | 用途 |
|----------|---------|------|
| `model.Token` | `model/token.go:14` | 本次扩展目标，新增 `IpPolicy *IpPolicy` 字段 |
| `gin.Context` | `github.com/gin-gonic/gin` | GetClientIP 参数类型；F004 写入 `"client_ip"` 上下文键 |

### 7.4 代码模式

| 模式名 | 参考位置 | 应用场景 |
|--------|---------|---------|
| `DB.Model(&Token{}).Where("id = ?", id).Update("field", value)` | `model/token.go:383-389`（increaseTokenQuota）| UpdateIpPolicy 单字段精确更新 |
| `defer func() { if shouldUpdateRedis(...) { gopool.Go(func()...) } }()` | `model/token.go:274-283`（Update()方法）| UpdateIpPolicy 异步缓存失效 |
| `strconv.Atoi(c.Param("id"))` + `c.GetInt("id")` | `controller/token.go:51-53`（GetToken）| UpdateTokenIpPolicy 解析参数 |
| `c.ShouldBindJSON(&req)` + `common.ApiError(c, err)` | `controller/token.go:229-231`（UpdateToken）| UpdateTokenIpPolicy 请求绑定 |
| `c.JSON(http.StatusOK, gin.H{"success": true, "message": ""})` | `controller/token.go:203-207`（AddToken）| UpdateTokenIpPolicy 成功响应 |

---

## 8. 测试规格

### 8.1 单元测试文件

**文件**: `common/ip_matcher_test.go`
**包**: `common`（与被测文件同包，白盒测试）
**框架**: `testing` + `testify/assert`（参见 `doc/kb/测试知识库/单元测试.md`）
**运行命令**: `go test -v -race ./common/ -run TestValidateCIDRList`

**测试函数命名规范**（遵循项目约定）：

```
Test{FuncName}_{场景描述}
// 示例：
TestValidateCIDRList_Empty
TestValidateCIDRList_ValidIPv4CIDR
TestValidateCIDRList_InvalidCIDRFormat
TestParseCIDRList_SingleIPAutoComplete
TestParseCIDRList_ConcurrentSafety
TestIPMatchesCIDRList_HitNetworkBoundary
TestGetClientIP_NoTrustedProxy
TestGetClientIP_TrustedProxyWithXFF
TestInitTrustedProxies_SkipInvalidEntry
```

**覆盖率目标**: ≥ 85%（通过 `go test -coverprofile=coverage.out ./common/` 验证）

**并发安全测试**（TC023 对应 UT023）：

```go
func TestParseCIDRList_ConcurrentSafety(t *testing.T) {
    ips := []string{"1.2.3.4", "10.0.0.0/8"}
    var wg sync.WaitGroup
    for i := 0; i < 100; i++ {
        wg.Add(1)
        go func() {
            defer wg.Done()
            result, err := ParseCIDRList(ips)
            require.NoError(t, err)
            assert.Len(t, result, 2)
        }()
    }
    wg.Wait()
}
// 运行: go test -race ./common/ -run TestParseCIDRList_ConcurrentSafety
```

### 8.2 集成测试覆盖点

以下场景需要通过 HTTP 集成测试覆盖（使用 `httptest.NewRecorder` + SQLite 内存模式）：

| 测试用例 ID | 测试点 |
|-----------|--------|
| TC001 | 设置白名单策略 → DB 验证 ip_policy 写入 |
| TC003 | 清除策略 → DB 验证 ip_policy 为 NULL |
| TC007 | 管理员操作他人 Token → 权限校验通过 |
| TC018 | 非拥有者非管理员 → 403 forbidden |
| TC042 | XFF 伪造防护（无 TRUSTED_PROXIES）→ 使用 RemoteAddr |

---

## 9. 遗留问题

| 编号 | 问题 | 对本文档的影响 |
|------|------|--------------|
| Q1 | 白名单 mode 下 `ips=[]` 时语义（拦截所有 or 放行所有？）待产品确认 | TC009 / 6.5 场景中 `hit=false`；当前实现：空白名单拦截所有（`!hit=true → blocked`）；确认后若需放行，需在 F005 增加 `len(Ips)==0` 短路条件 |
| Q4 | 黑名单策略下管理员 IP 是否豁免 | 当前设计：管理员无豁免，黑名单对所有调用者一视同仁；确认后若需豁免，需在 F005 增加管理员身份判断 |
| R6 | Redis 缓存主动失效延迟 | `UpdateIpPolicy` 调用 `cacheDeleteToken` 删除缓存；由于 Token 鉴权流程优先读 Redis（有 TTL），更新后最长可能有 `RedisKeyCacheSeconds()` 秒的延迟；当前接受此风险，测试 TC010 需等待缓存过期后验证 |

---

## 10. 变更记录

| 版本 | 日期 | 变更说明 |
|------|------|---------|
| v1.0 | 2026-03-02 | 初稿，基于概要设计 V1.0 和测试用例文档 V1.0 生成；覆盖 F001-F005 全部需求；源码分析版本 `21cfc1ca` |
