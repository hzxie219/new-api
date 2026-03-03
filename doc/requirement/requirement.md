# API Key IP 访问控制 — 后端需求分析文档

> 基于系统需求文档深度分析生成
> 生成时间: 2026-03-02
> 分析维度: 功能细化 | 预期效果 | 依赖关系

## 文档说明

本文档通过深度分析 `doc/requirement-analyst/output/` 目录下的 15 个系统需求文档，结合对 `new-api` 代码仓库真实源码的检索，将每个功能需求细化为清晰的功能描述，包含代码引用位置、输入输出规范、边界条件等信息，可作为系统设计和开发的直接输入。

---

## 1. 需求背景和目标

### 1.1 背景概述

当前 API Key 系统（`model/token.go`）已有基础的 IP 访问控制能力：`AllowIps *string` 字段（第 27 行）支持精确 IP 白名单（逐行分割，无 CIDR 支持），鉴权中间件（`middleware/auth.go` 第 316–330 行）在 Token 有效性校验后对 AllowIps 进行匹配。

然而，现有机制存在以下不足：

1. **不支持 CIDR 网段**：企业用户通常使用整个子网段（如 `10.0.0.0/8`），逐 IP 录入成本极高
2. **缺少黑名单模式**：只有白名单，无法快速封堵攻击来源 IP
3. **无可信代理感知**：使用 `c.ClientIP()`（gin 内置，基于 `RemoteAddr`），在反向代理架构下无法获取真实客户端 IP，且未保护 X-Forwarded-For 头的伪造攻击
4. **新旧功能缺乏统一化**：旧字段（`AllowIps`）与未来可能扩展的策略能力之间无结构化连接

本次需求在不破坏现有 `AllowIps` 逻辑的前提下，新增 `IpPolicy` 结构化字段，支持 CIDR 白名单/黑名单双模式，并补齐可信代理感知 IP 提取能力。

### 1.2 用户角色和场景

| 角色 | 核心场景 |
|------|---------|
| API Key 使用者（开发者/企业用户） | 期望所配置的 Key 只能从企业 IP 段使用，防止 Key 泄露后被外部滥用 |
| 平台管理员 | 为每个 Key 配置 IP 策略（白名单/黑名单），执行合规审计 |
| 平台运营者 | 在安全事件中快速通过黑名单封堵攻击来源 IP 段 |

### 1.3 核心痛点

- **痛点 1**：CIDR 支持缺失 — 企业网络使用 `10.0.0.0/8`、`172.16.0.0/12` 等大型子网段，精确 IP 列举不可操作
  - 影响范围：企业用户（API Key 使用者、管理员）
  - 严重程度：高

- **痛点 2**：无黑名单能力 — 攻击来源 IP 已知时，只能下线整个 Key 而无法精准封堵
  - 影响范围：平台运营者、安全团队
  - 严重程度：高

- **痛点 3**：IP 提取不可信 — 反向代理（Nginx/CDN）后 `c.ClientIP()` 返回代理 IP 而非真实客户端 IP；XFF 头可被伪造绕过白名单
  - 影响范围：所有使用反向代理的部署环境
  - 严重程度：高（安全漏洞）

### 1.4 用户旅程

```
管理员配置 IP 策略:
  PUT /api/token/:id/ip_policy → 校验 CIDR 格式 → 持久化到 tokens.ip_policy → 策略生效

API 请求鉴权流程（含 IP 策略执行）:
  接收请求
  → 提取真实客户端 IP（可信代理感知）
  → Token 有效性校验（存在/启用/有效期/额度/模型白名单）
  → 旧 AllowIps 白名单校验
  → 新 IpPolicy 策略校验（白名单/黑名单）
  → 通过 → 转发上游
  → 拒绝 → HTTP 403 + IP_NOT_ALLOWED
```

### 1.5 需求目标

- **目标 1**：支持 CIDR 格式的 IP 白名单 — 成功标准：能配置 `10.0.0.0/8`，`10.x.x.x` 范围内的请求均通过
- **目标 2**：支持 CIDR 格式的 IP 黑名单 — 成功标准：黑名单 IP 段的请求返回 HTTP 403 + `IP_NOT_ALLOWED`
- **目标 3**：可信代理感知 IP 提取 — 成功标准：XFF 伪造攻击（未经可信代理）被正确忽略，IP 校验生效
- **目标 4**：向后兼容 — 成功标准：未配置 `IpPolicy` 的 Token，鉴权行为与现有逻辑完全一致

---

## 2. 功能需求详细分析（后端）

> **🚨 重要约束**: 本章节只包含功能性需求（业务功能、接口、数据处理逻辑等），非功能性需求（性能、安全、可靠性等）在第 5 章描述

### 2.1 功能概览

| 功能编号 | 功能名称 | 类型 | 实现类型 | 所属仓库 | 相关现有接口 | 需要方案设计 | 需要方案选型 | 优先级 | 复杂度 | 依赖功能 |
|---------|---------|------|---------|---------|------------|------------|------------|-------|-------|---------|
| F001 | IP 策略配置接口 | 全新功能 | 接口功能 | new-api | - | 否 | 否 | P0 | 中 | F002, F003 |
| F002 | Token.IpPolicy 数据模型扩展 | 增量功能 | 后台功能 | new-api | - | 否 | 否 | P0 | 低 | - |
| F003 | CIDR 工具函数库与可信代理管理 | 全新功能 | 后台功能 | new-api | - | 否 | 否 | P0 | 中 | - |
| F004 | 可信代理感知 IP 提取（鉴权中间件扩展） | 增量功能 | 后台功能 | new-api | - | 否 | 否 | P0 | 低 | F003 |
| F005 | IP 策略执行（鉴权中间件扩展） | 增量功能 | 后台功能 | new-api | - | 否 | 否 | P0 | 中 | F002, F003, F004 |

**说明**：
- 所有功能均在 `new-api` 单一仓库内实现，无跨仓库依赖
- 无需方案选型：全部使用 Go 标准库 `net` 包实现 CIDR 解析与匹配，无新增外部依赖
- F002 是其他功能的基础（数据字段），F003 是工具层（被 F001/F005 调用），F004 是 F005 的前置

---

### 2.2 功能详细描述

#### 2.2.1 接口功能

---

##### 功能 F001: IP 策略配置接口（全新功能）

###### 功能描述

**功能类型**：全新功能
**功能实现类型**：接口功能
**业务目标**：为管理员/Key 拥有者提供 HTTP 接口，对指定 API Key 设置或清除 IP 访问策略（白名单/黑名单 + CIDR 列表）
**技术目标**：在 `controller/token.go` 中新增 `UpdateTokenIpPolicy` handler，在 `router/api-router.go` 中注册新路由 `PUT /api/token/:id/ip_policy`

**参考功能**：
- 可参考功能：现有 `UpdateToken` handler — 代码位置：`controller/token.go:225`
- 参考价值：Token 权限校验模式（先查 Token 是否存在，再校验是否属于当前用户/是否为管理员，再执行更新）
- 现有 token 路由注册位置：`router/api-router.go:241-250`（`tokenRoute` 路由组，挂载于 `/api/token/`，已使用 `middleware.UserAuth()`）

###### 输入定义

**新增接口**：`PUT /api/token/:id/ip_policy`

| 输入项 | 类型 | 来源 | 必填 | 约束条件 | 示例 |
|-------|------|------|------|---------|------|
| `id` | Integer | URL 路径参数 | 是 | > 0，目标 Token 的数据库 ID | `42` |
| `mode` | String | 请求体 JSON | 是 | 枚举值：`"whitelist"`、`"blacklist"`、`""` (清除策略) | `"whitelist"` |
| `ips` | Array\<String\> | 请求体 JSON | 是 | 每条为合法 IPv4/IPv6 或 CIDR；条目数 0–100；空数组合法 | `["1.2.3.4", "10.0.0.0/8"]` |

**请求体示例**：

```json
{
  "mode": "whitelist",
  "ips": ["1.2.3.4", "10.0.0.0/8", "192.168.1.0/24"]
}
```

**输入校验规则**（严格按顺序执行）：

1. **路径参数 `id`**：必须为正整数；Token 必须存在（否则 404）
2. **权限校验**：当前登录用户必须是 Token 的拥有者（`token.UserId == currentUserId`）或管理员（否则 403）
3. **`mode` 校验**：值必须为 `"whitelist"`、`"blacklist"` 或 `""`（否则 400）
4. **`ips` 格式校验**：调用 `common.ValidateCIDRList(req.Ips)` 校验每条格式（否则 400）；单 IP（无 `/` 前缀）自动视为 `/32`
5. **条目数校验**：`len(req.Ips) <= 100`（否则 400）

注：当 `mode == ""` 时，清除策略，`ips` 值被忽略（可传任意值）

###### 输出定义

**成功响应（HTTP 200）**：

```json
{
  "success": true,
  "message": ""
}
```

**失败响应**：

| HTTP 状态码 | 错误信息 | 触发条件 | 处理建议 |
|-----------|---------|---------|---------|
| 400 | `"invalid mode: xxx"` | `mode` 不在枚举值内 | 检查 mode 字段值 |
| 400 | `"invalid IP/CIDR: xxx"` | `ips` 中有非法 CIDR 条目 | 检查 ips 格式 |
| 400 | `"too many IP entries"` | `len(ips) > 100` | 减少条目数量 |
| 403 | `"forbidden"` | 非 Token 拥有者且非管理员 | 确认操作权限 |
| 404 | `"not found"` | Token ID 不存在 | 确认 Token ID |
| 500 | `"internal server error"` | 数据库写入失败 | 检查数据库连接 |

###### 边界条件处理

| 边界条件 | 处理方式 |
|---------|---------|
| `mode == ""` 且 `ips` 非空 | 清除策略（`ips` 忽略），不报错 |
| `mode` 非空且 `ips == []` | 暂定：等同无策略（待确认 Q1）；当前实现存储 `{mode, ips:[]}` |
| 单 IP 无 `/` 前缀（如 `"1.2.3.4"`） | 校验时自动补全为 `"1.2.3.4/32"`，合法通过 |
| IPv6 地址（如 `"::1/128"`） | `net.ParseCIDR` 原生支持，合法通过 |
| Token 属于他人，当前用户为管理员 | 允许操作（管理员可管理任意 Token） |
| 数据库并发写入（同 Token 并发更新策略） | GORM 写入为原子操作，最后写入者生效 |

---

#### 2.2.2 后台功能

---

##### 功能 F002: Token.IpPolicy 数据模型扩展（增量功能）

###### 功能描述

**功能类型**：增量功能
**功能实现类型**：后台功能
**业务目标**：为 `Token` 结构体增加 `IpPolicy` 字段，将 IP 策略（模式+CIDR列表）持久化到数据库
**技术目标**：在 `model/token.go` 中定义 `IpPolicy` 结构体并实现 `driver.Valuer` / `sql.Scanner` 接口，通过 JSON 序列化存储到 `tokens.ip_policy TEXT` 列

**相关现有功能**：
- 现有模型：`Token` struct — 代码位置：`model/token.go:14`
- 现有字段：`AllowIps *string` (第 27 行)，`GetIpLimits()` 方法 (第 38 行) — 旧精确 IP 白名单，**保持不变，不删除**
- 增量方向：结构体新增字段 + 新结构体定义 + 序列化实现
- 兼容性注意：GORM AutoMigrate 会在 `tokens` 表中 ADD COLUMN，不影响现有行；旧行 `ip_policy` 为 NULL，中间件逻辑中 `token.IpPolicy == nil` 时跳过策略校验（向后兼容）

**AutoMigrate 已包含 `&Token{}`**：
- 代码位置：`model/main.go:256` — `DB.AutoMigrate(&Token{}, ...)` 已存在
- 无需修改 `model/main.go`，添加字段后 AutoMigrate 自动执行 ADD COLUMN

**新增内容**：

1. `IpPolicy` 结构体（新）：
   ```
   type IpPolicy struct {
       Mode string   `json:"mode"` // "whitelist" | "blacklist" | ""
       Ips  []string `json:"ips"`  // IP/CIDR 列表
   }
   ```

2. `Token` 结构体扩展（增量，在 `model/token.go:14` 的 Token struct 中新增字段）：
   ```
   IpPolicy *IpPolicy `gorm:"type:text;column:ip_policy" json:"ip_policy,omitempty"`
   ```

3. `IpPolicy.Value()` 实现 `driver.Valuer`：JSON 序列化为字符串存入数据库
4. `IpPolicy.Scan()` 实现 `sql.Scanner`：从数据库字符串反序列化为 `IpPolicy` 结构体

**数据库兼容性**：
- `type:text` 在 SQLite/MySQL/PostgreSQL 均支持（不使用 JSONB，不使用 MySQL/PG 专有类型）
- SQLite 不支持 ALTER COLUMN，但 `AutoMigrate` 的 ADD COLUMN 操作在三数据库均支持

###### 边界条件处理

| 边界条件 | 处理方式 |
|---------|---------|
| 旧记录 `ip_policy` 为 NULL | GORM 扫描为 `nil`，中间件判断 `token.IpPolicy == nil` 跳过校验，行为与旧版一致 |
| `IpPolicy.Ips` 为空切片 | `Value()` 序列化为 `{"mode":"whitelist","ips":[]}`，`Scan()` 反序列化后 `Ips` 为 `[]string{}`（非 nil） |
| JSON 序列化失败 | `Value()` 返回 error，数据库写入失败，controller 层返回 500 |
| 数据库字段被手工篡改为非法 JSON | `Scan()` 返回 error，GORM 扫描失败，该 Token 鉴权时 `IpPolicy` 为 nil，降级为无策略（宽松降级） |

---

##### 功能 F003: CIDR 工具函数库与可信代理管理（全新功能）

###### 功能描述

**功能类型**：全新功能
**功能实现类型**：后台功能
**业务目标**：提供 CIDR 格式校验、解析、IP 命中判断和可信代理管理的通用工具函数，是整个 IP 访问控制功能的核心算法库
**技术目标**：新建 `common/ip_matcher.go` 文件，实现 5 个函数和 1 个包级全局变量，全部使用 Go 标准库 `net` 包，无外部依赖

**参考功能**：
- 可参考功能：现有 `IsIpInCIDRList` — 代码位置：`common/ip.go:33`
- 参考价值：`net.ParseCIDR` + `cidr.Contains(ip)` 的基础模式；但现有函数签名为 `(ip net.IP, cidrList []string)`，每次调用均解析 CIDR 字符串，性能次优；新实现拆分为 Parse（一次解析）+ Match（多次复用）

**新增内容（`common/ip_matcher.go`，全新文件）**：

| 函数/变量 | 签名 | 调用方 | 用途 |
|---------|-----|-------|------|
| `ValidateCIDRList` | `func ValidateCIDRList(ips []string) error` | F001 controller 层入参校验 | 返回第一个非法条目的 error，合法返回 nil |
| `ParseCIDRList` | `func ParseCIDRList(ips []string) ([]*net.IPNet, error)` | F005 中间件策略执行 | 将 IP/CIDR 字符串列表解析为 `[]*net.IPNet` |
| `IPMatchesCIDRList` | `func IPMatchesCIDRList(ipStr string, cidrs []*net.IPNet) bool` | F005 中间件策略执行 | 判断 IP 是否命中列表中任意 CIDR |
| `trustedProxyCIDRs` | `var trustedProxyCIDRs []*net.IPNet` | 包内 `isTrustedProxy` | 全局存储可信代理 CIDR 列表（服务启动时初始化一次） |
| `InitTrustedProxies` | `func InitTrustedProxies(proxies string)` | F004 `common/init.go` | 解析 `TRUSTED_PROXIES` 环境变量，服务启动时调用一次 |
| `GetClientIP` | `func GetClientIP(c *gin.Context) string` | F004 中间件 | 可信代理感知 IP 提取：若直连返回 RemoteAddr；若来自可信代理取 XFF 最左侧 IP |

**核心逻辑说明**：

- **`ValidateCIDRList` / `ParseCIDRList`**：单 IP（无 `/` 前缀）自动补全为 `/32`，再调用 `net.ParseCIDR` 解析；遇第一个非法条目立即返回 error
- **`IPMatchesCIDRList`**：调用 `net.ParseIP(ipStr)` 解析 IP，若非法返回 false（容错）；遍历 `cidrs` 调用 `cidr.Contains(ip)`，任一命中返回 true
- **`GetClientIP`**：解析 `c.Request.RemoteAddr` 获取直接连接方 IP；若该 IP 在 `trustedProxyCIDRs` 中，取 `X-Forwarded-For` 头的最左侧字段（`SplitN(xff, ",", 2)[0]`）；若 XFF 为空则回退到 RemoteAddr
- **`InitTrustedProxies`**：按逗号分割输入字符串，逐条解析（单 IP 补全 `/32`），非法条目跳过，更新 `trustedProxyCIDRs` 全局变量

**测试文件**：`common/ip_matcher_test.go`（新建）

###### 边界条件处理

| 边界条件 | 处理方式 |
|---------|---------|
| 单 IP 无 `/` 前缀（`"1.2.3.4"`） | `ValidateCIDRList` / `ParseCIDRList` 均自动补全为 `"1.2.3.4/32"` |
| `ips = []`（空列表） | `ValidateCIDRList` 返回 nil；`ParseCIDRList` 返回空 slice、nil error |
| `IPMatchesCIDRList` 传入非法 IP 字符串 | `net.ParseIP` 返回 nil，函数返回 false，不 panic |
| `IPMatchesCIDRList` 传入空 CIDR 列表 | 遍历为空，返回 false |
| TRUSTED_PROXIES 包含非法 CIDR 条目 | `InitTrustedProxies` 跳过非法条目，其余合法条目正常加载 |
| TRUSTED_PROXIES 为空字符串 | `InitTrustedProxies("")` 将 `trustedProxyCIDRs` 设为 nil（空，无可信代理） |
| XFF 头为空字符串 | `GetClientIP` 回退到 `RemoteAddr` |
| XFF 头为多跳（`"client, proxy1, proxy2"`） | 取最左侧（`SplitN(xff, ",", 2)[0]` → `"client"`） |
| IPv6 地址 | `net.ParseCIDR` 和 `net.ParseIP` 原生支持 IPv6，行为一致 |
| `net.ParseCIDR` 将主机位归零（`"1.2.3.4/24"` → `"1.2.3.0/24"`） | 符合 CIDR 语义（匹配整个网段），`ValidateCIDRList` 不报错 |

---

##### 功能 F004: 可信代理感知 IP 提取（鉴权中间件扩展）（增量功能）

###### 功能描述

**功能类型**：增量功能
**功能实现类型**：后台功能
**业务目标**：在每次 Token 鉴权时，使用可信代理感知方式提取真实客户端 IP，并存入 gin Context，供后续 IP 策略执行（F005）使用，避免重复提取
**技术目标**：
1. 在 `middleware/auth.go:TokenAuth()` 的 Token 有效性校验链末尾、IP 策略执行块之前，添加 `GetClientIP` 调用并写入 Context
2. 在 `common/init.go:InitEnv()` 函数末尾，添加 `InitTrustedProxies(os.Getenv("TRUSTED_PROXIES"))` 调用，实现服务启动时一次性初始化可信代理配置

**相关现有功能**：
- 现有文件：`middleware/auth.go`，`TokenAuth()` — 代码位置：`middleware/auth.go:248`
- 现有 IP 提取：`c.ClientIP()` — 代码位置：`middleware/auth.go:318`（gin 内置，不区分代理）
- 现有初始化：`common/init.go:InitEnv()` — 代码位置：`common/init.go:31`
- 增量方向：
  - `middleware/auth.go`：在现有 AllowIps 校验块（第 316–330 行）之后插入 `GetClientIP` 调用
  - `common/init.go`：在 `InitEnv` 函数末尾追加 1 行 `InitTrustedProxies` 调用

**新增内容**：

```
// middleware/auth.go — 在 Token 有效性校验通过后（原 AllowIps 块之后）插入
clientIP := common.GetClientIP(c)
c.Set("client_ip", clientIP)
```

```
// common/init.go — 在 InitEnv() 末尾追加
common.InitTrustedProxies(os.Getenv("TRUSTED_PROXIES"))
```

**新增环境变量**：

| 环境变量 | 格式 | 默认值 | 说明 |
|---------|------|-------|------|
| `TRUSTED_PROXIES` | 逗号分隔的 IP 或 CIDR 列表，如 `"10.0.0.1,172.16.0.0/12"` | 空（无可信代理） | 配置可信反向代理 IP 段；空时 XFF 头始终被忽略，使用 RemoteAddr |

###### 边界条件处理

| 边界条件 | 处理方式 |
|---------|---------|
| `TRUSTED_PROXIES` 未配置 | `InitTrustedProxies("")` 将 `trustedProxyCIDRs` 置 nil；所有请求使用 `RemoteAddr` |
| `TRUSTED_PROXIES` 包含非法条目 | 跳过非法条目（`InitTrustedProxies` 内已处理），合法条目生效；建议在 `InitEnv` 中记录 WARN 日志 |
| `c.Request.RemoteAddr` 格式异常 | `net.SplitHostPort` 失败，`GetClientIP` 直接使用 `RemoteAddr` 字符串回退 |
| 后续 IP 策略执行（F005）读取 Context | 通过 `c.Get("client_ip")` 读取，避免重复调用 `GetClientIP` |

---

##### 功能 F005: IP 策略执行（鉴权中间件扩展）（增量功能）

###### 功能描述

**功能类型**：增量功能
**功能实现类型**：后台功能
**业务目标**：在每次 Token 鉴权时，根据 Token 绑定的 `IpPolicy`（白名单或黑名单）对真实客户端 IP 进行匹配，拦截不符合策略的请求并返回标准化 HTTP 403 响应
**技术目标**：在 `middleware/auth.go:TokenAuth()` 中（紧随 F004 IP 提取之后），新增 IP 策略校验块

**相关现有功能**：
- 现有文件：`middleware/auth.go`，`TokenAuth()` — 代码位置：`middleware/auth.go:248`
- 现有 AllowIps 逻辑：第 316–330 行，使用旧 `GetIpLimits()` 和 `common.IsIpInCIDRList()`，**保持不变**
- 增量方向：在原 AllowIps 块之后（约第 330 行）和 `SetupContextForToken` 之前插入新 IP 策略块

**完整鉴权链执行顺序**（插入后）：

```
1. Token 存在性校验          → 不存在: 401/404
2. Token 状态校验（启用）    → 禁用: 403
3. Token 有效期校验          → 已过期: 401
4. Token 额度校验            → 额度不足: 429
5. Token 模型白名单校验      → 模型不允许: 403
   [旧 AllowIps 精确 IP 白名单校验（现有逻辑，middleware/auth.go:316-330）]
   [← F004: clientIP := common.GetClientIP(c); c.Set("client_ip", clientIP)]
6. ← IP 策略校验（F005 新增块）→ IP 被拒: HTTP 403 + IP_NOT_ALLOWED
7. c.Next() → 转发上游/计费/日志
```

**新增策略块逻辑**：

```
从 Context 读取 clientIP（F004 写入，避免重复提取）
若 token.IpPolicy == nil 或 token.IpPolicy.Mode == ""：
    跳过策略校验，继续
若调用 common.ParseCIDRList(token.IpPolicy.Ips) 返回 error：
    写 WARN 日志（含 token_id + error）
    不拦截（宽松降级，服务可用性优先）
    继续
否则：
    hit := common.IPMatchesCIDRList(clientIPStr, cidrs)
    blocked := (mode == "whitelist" && !hit) || (mode == "blacklist" && hit)
    若 blocked：
        写 INFO 日志（含 token_id + client_ip + mode）
        返回 HTTP 403 + {"success":false,"message":"IP not allowed","error_code":"IP_NOT_ALLOWED"}
        c.Abort() + return
    否则继续
```

**响应格式**（IP 拦截时）：

```json
{
  "success": false,
  "message": "IP not allowed",
  "error_code": "IP_NOT_ALLOWED"
}
```

**日志格式**：

| 事件 | 级别 | 关键字 | 字段 |
|-----|------|-------|-----|
| IP 被拦截（正常鉴权拒绝） | INFO | `ip_not_allowed` | `token_id`, `client_ip`, `mode` |
| CIDR 解析失败（降级） | WARN | `ip_policy_parse_failed` | `token_id`, `error` |

###### 边界条件处理

| 边界条件 | 处理方式 |
|---------|---------|
| `token.IpPolicy == nil` | 跳过策略校验（向后兼容旧记录） |
| `token.IpPolicy.Mode == ""` | 跳过策略校验（策略已被清除） |
| `token.IpPolicy.Ips == []`（空列表），mode 非空 | `ParseCIDRList([])` 返回空 slice；`IPMatchesCIDRList` 返回 false；白名单下全部拒绝，黑名单下全部通过（待确认 Q1） |
| CIDR 解析失败（数据库被篡改） | 宽松降级：写 WARN 日志，不拦截，继续处理 |
| `c.Abort()` 调用后 | 必须立即 `return`，防止函数体后续代码继续执行（gin 的 Abort 只阻止后续中间件，不终止当前函数） |
| Token 已过期/禁用时 | 鉴权链在步骤 2/3 已返回错误，不进入 IP 策略校验块（IP 策略仅对有效 Token 生效） |
| XFF 伪造（未配置 TRUSTED_PROXIES） | F004 的 `GetClientIP` 返回 RemoteAddr（真实连接方 IP），XFF 头被忽略，攻击无效 |

---

## 3. 静态结构分析

> **本章节说明**: 分析本需求需要修改或新增的数据表和配置

### 3.1 数据表变更分析

#### 关系型数据库（SQLite / MySQL / PostgreSQL 三库兼容）

**需要修改的表**：

| 表名 | 数据库 | 修改内容 | 修改目的 | 兼容性 |
|-----|-------|---------|---------|-------|
| `tokens` | SQLite / MySQL / PostgreSQL | 新增列 `ip_policy TEXT NULL` | 存储 JSON 序列化的 `IpPolicy` 结构（mode + ips 列表），供鉴权中间件在每次请求时读取 | 通过 GORM `AutoMigrate`（`model/main.go:256` 已包含 `&Token{}`）自动执行 ADD COLUMN，无需手工迁移脚本 |

**迁移方式**：
- 无需修改 `model/main.go`
- `Token` struct 添加 `IpPolicy *IpPolicy \`gorm:"type:text;column:ip_policy"\`` 字段后，GORM AutoMigrate 在三数据库执行 `ALTER TABLE tokens ADD COLUMN ip_policy TEXT`
- SQLite 支持 ADD COLUMN（不支持 ALTER COLUMN，本次无 ALTER COLUMN 操作）
- 旧记录 `ip_policy = NULL`，鉴权时 `token.IpPolicy == nil` 触发跳过逻辑，向后兼容

**不需要新增的表**：本需求通过扩展现有 `tokens` 表字段实现，不引入新表。

---

### 3.2 配置文件变更分析

本需求不修改任何传统配置文件（无 YAML/JSON 配置文件变更）。

**新增环境变量**（在 `common/init.go:InitEnv()` 中读取）：

| 环境变量 | 读取位置 | 格式 | 默认值 | 用途 |
|---------|---------|------|-------|------|
| `TRUSTED_PROXIES` | `common/init.go:InitEnv()`（新增 1 行）| 逗号分隔的 IP 或 CIDR 字符串 | 空字符串 | 指定可信反向代理 IP 段，用于控制是否信任 X-Forwarded-For 头 |

**变更内容**：在 `common/init.go:InitEnv()` 函数末尾追加：

```go
common.InitTrustedProxies(os.Getenv("TRUSTED_PROXIES"))
```

---

## 4. 消息格式分析

本需求不涉及修改或新增任何消息队列（Pulsar/Kafka）消息格式。所有 IP 策略的配置和执行均在同步 HTTP 请求处理链中完成。

---

## 5. 非功能需求详细说明

> **🚨 重要**: 非功能性需求只在本章描述，不在第 2 章"功能需求详细分析"中出现

### 5.1 性能需求

| 指标类型 | 功能 | 具体要求 | 测量方法 |
|---------|------|---------|---------|
| 接口响应时间 | F001（IP 策略配置接口） | P99 < 200ms（含数据库写入） | 压测工具统计 |
| 中间件延迟 | F004（GetClientIP） | P99 < 0.1ms（纯内存操作） | race detector 下 benchmark |
| 中间件延迟 | F005（IP 策略校验） | P99 < 1ms（O(n)，n ≤ 100 条 CIDR） | race detector 下 benchmark |
| 可信代理匹配 | F004（isTrustedProxy） | P99 < 0.1ms（通常代理数 < 10） | 与 GetClientIP 合并测量 |
| 内存占用 | F003（trustedProxyCIDRs 全局变量） | < 1KB（< 50 条 CIDR） | 进程内存监控 |
| CIDR 解析频次 | F005 | 每次请求调用 `ParseCIDRList`（无缓存） | 如压测发现瓶颈，可增加 Token 级 CIDR 缓存（TTL=60s） |

### 5.2 安全性需求

| 安全维度 | 要求 | 覆盖功能 |
|---------|------|---------|
| 防 XFF 伪造 | `GetClientIP` 仅在请求来自可信代理时采信 XFF 头，否则使用 RemoteAddr | F003, F004 |
| 权限控制 | `UpdateTokenIpPolicy` 校验当前用户必须为 Token 拥有者或管理员（双重校验） | F001 |
| 输入格式校验 | CIDR 字符串必须通过 `ValidateCIDRList` 校验，非法格式返回 400，拒绝入库 | F001, F003 |
| 防策略信息泄露 | 无论白名单未命中还是黑名单命中，错误码统一为 `IP_NOT_ALLOWED`，不区分模式 | F005 |
| 鉴权链完整性 | IP 策略校验在 Token 有效性校验之后执行，无效 Token 不进入 IP 策略块 | F005 |
| 请求链终止 | 拦截时 `c.Abort()` 后立即 `return`，防止后续 handler（计费、转发等）执行 | F005 |
| 降级安全语义 | CIDR 解析失败时宽松降级（不拦截），避免因数据库数据错误导致服务不可用 | F005 |

### 5.3 可靠性需求

| 可靠性指标 | 要求 |
|-----------|------|
| 降级行为 | CIDR 解析失败时不拦截请求，保证服务可用性；须记录 WARN 日志 |
| 请求隔离 | IP 策略校验失败不影响其他 Token 的请求处理（无全局状态污染） |
| 数据库写入失败 | F001 接口返回 HTTP 500，策略未更新，不留中间状态 |
| 向后兼容 | `IpPolicy == nil` 的旧 Token，鉴权行为与现有逻辑完全一致 |

### 5.4 可观测性需求

| 观测维度 | 要求 |
|---------|------|
| IP 拦截事件日志 | INFO 级别，包含 `token_id`、`client_ip`、`mode`，关键字 `ip_not_allowed`（供审计和告警） |
| CIDR 解析失败日志 | WARN 级别，包含 `token_id`、`error`，关键字 `ip_policy_parse_failed`（供排障） |
| 频繁拦截告警 | `ip_not_allowed` 频率异常时可配置告警（结合 `client_ip` 分析攻击来源） |
| IP 策略未生效排查 | 若无 `ip_not_allowed` 日志但策略未生效，需排查：Token.IpPolicy 字段写入是否正确；TRUSTED_PROXIES 配置；鉴权链顺序 |

### 5.5 可维护性需求

- **单元测试**：新建 `common/ip_matcher_test.go`，覆盖 `ValidateCIDRList`、`ParseCIDRList`、`IPMatchesCIDRList`、`GetClientIP`、`InitTrustedProxies` 全部函数及边界场景
- **集成测试**：扩展 `middleware/auth_test.go`，覆盖白名单/黑名单拦截（含响应格式验证）、无策略通过、CIDR 解析失败降级、鉴权链顺序
- **并发安全测试**：在 race detector 模式下运行并发测试（`ParseCIDRList` 每次返回独立 slice，无共享状态）
- **代码位置清晰**：`ip_matcher.go` 为独立纯函数文件（除 `trustedProxyCIDRs` 全局变量外无 IO 操作），便于单独测试

### 5.6 可扩展性需求

- **TRUSTED_PROXIES 配置化**：通过环境变量配置，部署时按需调整，不需要重新编译
- **Token 级缓存扩展点**：若压测发现 `ParseCIDRList` 频繁解析成为瓶颈，可在不修改接口的情况下在 `middleware/auth.go` 内增加 Token 级 CIDR 缓存（TTL=60s）
- **条目数上限可配置**：当前硬编码 100 条，后续可通过配置项调整

---

## 6. 风险和依赖分析

### 6.1 技术风险

| 风险 ID | 风险描述 | 影响功能 | 概率 | 影响 | 应对策略 |
|--------|---------|---------|------|------|---------|
| R001 | 每次请求调用 `ParseCIDRList` 重复解析（无缓存），高频 Token 有额外开销 | F005 | 中 | 中 | 压测发现瓶颈后增加 Token 级 CIDR 缓存（TTL=60s） |
| R002 | `TRUSTED_PROXIES` 配置错误（信任了非代理 IP），导致攻击者可伪造 XFF 绕过 IP 策略 | F003, F004 | 低 | 高 | 部署文档明确指导；仅配置明确受控的代理 IP/CIDR |
| R003 | 数据库 `ip_policy` 字段被手工篡改为非法 JSON，导致 Token 鉴权时 `Scan` 失败 | F002, F005 | 低 | 低 | 中间件降级不拦截（宽松降级）+ WARN 日志，不影响服务可用性 |
| R004 | SQLite AutoMigrate ADD COLUMN 在某些版本下行为差异 | F002 | 低 | 中 | 在 SQLite/MySQL/PostgreSQL 三库验证 AutoMigrate 结果 |

### 6.2 业务风险（待确认问题）

| 编号 | 问题 | 影响功能 | 优先级 | 暂定行为 | 需确认方 |
|-----|------|---------|--------|---------|---------|
| Q1 | `ips=[]` 且 `mode` 非空时的语义（空白名单是拒绝一切还是等同无策略？） | F001, F005 | 🔴 高 | 暂定：白名单下拒绝一切（`IPMatchesCIDRList` 始终 false）；黑名单下通过一切 | 产品确认 |
| Q2 | 单 Key IP 策略条目数量上限 | F001 | 🟡 中 | 暂定 100 条 | 产品确认 |
| Q3 | 是否需要 `GET /api/token/:id/ip_policy` 查询接口 | F001 | 🟡 中 | 当前不实现，可后续补充 | 产品确认 |
| Q4 | 普通用户是否可设置黑名单模式（还是仅管理员可用？） | F001 | 🔴 高 | 暂定：仅管理员可设置黑名单；普通用户只能设置白名单 | 产品确认 |
| Q6 | IP 拒绝事件是否写入数据库 `logs` 表（`type=5`），供管理后台审计 | F005 | 🟡 中 | 暂定：仅写日志文件，不入库 | 产品确认 |

### 6.3 依赖关系图

```
[API Key IP 访问控制]
  ├── F001 (IP策略配置接口)
  │     ├── 依赖 → F002 (Token.IpPolicy 字段，用于存储策略)
  │     └── 依赖 → F003 (ValidateCIDRList，用于入参校验)
  ├── F002 (Token.IpPolicy 数据模型扩展)
  │     └── 依赖 → GORM AutoMigrate（model/main.go:256 已包含 &Token{}）
  ├── F003 (CIDR 工具函数库，全新文件 common/ip_matcher.go)
  │     └── 依赖 → Go 标准库 net 包（无外部依赖）
  ├── F004 (可信代理感知 IP 提取)
  │     └── 依赖 → F003 (GetClientIP, InitTrustedProxies)
  └── F005 (IP 策略执行)
        ├── 依赖 → F002 (Token.IpPolicy 字段)
        ├── 依赖 → F003 (ParseCIDRList, IPMatchesCIDRList)
        └── 依赖 → F004 (clientIP 已存入 gin Context)

[被依赖情况]
  └── 所有使用 API Key 的请求路由 ← 依赖 middleware/auth.go（F004/F005 修改此文件）
```

### 6.4 外部依赖详细说明

| 依赖项 | 类型 | 必需性 | 说明 |
|-------|------|-------|------|
| Go 标准库 `net` 包 | 工具库 | 必需 | `net.ParseCIDR`、`net.ParseIP`、`(*net.IPNet).Contains()`、`net.SplitHostPort`；无需引入外部依赖 |
| GORM AutoMigrate | ORM 框架 | 必需 | 已在项目中使用，`model/main.go:254` 已包含 `&Token{}`，无需新增配置 |
| gin Framework | HTTP 框架 | 必需 | `GetClientIP` 依赖 `*gin.Context`，已在项目中全面使用 |
| 数据库（SQLite/MySQL/PostgreSQL） | 存储 | 必需 | 新增 `ip_policy TEXT NULL` 列，三数据库均通过 GORM 兼容处理 |

### 6.5 影响范围评估

**本需求实现后的影响**：

- ✅ 正向影响：
  - 每个 API Key 可独立配置 CIDR 级别的 IP 访问策略，满足企业合规需求
  - 黑名单模式支持在安全事件中快速封堵攻击来源 IP 段，无需下线 Key
  - 可信代理感知消除了现有架构中 IP 校验不可信的安全漏洞
  - 所有已有 Key（`IpPolicy == nil`）行为完全不变，零停机升级

- ⚠️ 需要注意：
  - 部署时需根据实际网络架构配置 `TRUSTED_PROXIES` 环境变量；未配置时 XFF 头永远被忽略（适合直连部署，不适合经过代理的部署）
  - Q1/Q4 两个高优先级问题需在开发前与产品确认，否则影响 F001/F005 实现逻辑
  - 每次请求新增 `ParseCIDRList` 调用，对已配置 `IpPolicy` 的 Token 有轻微性能影响（预估 < 1ms，但需压测验证）

- 📋 建议同步的团队：
  - 产品团队：确认 Q1（空白名单语义）和 Q4（普通用户黑名单权限）
  - 运维团队：更新部署文档，说明 `TRUSTED_PROXIES` 环境变量的配置方式和安全注意事项
  - QA 团队：根据 `doc/requirement-analyst/output/` 下的 AC 文档覆盖测试场景

### 6.6 推荐开发顺序

```
阶段 1: 工具函数（可独立开发和单测，无其他依赖）
  → common/ip_matcher.go（F003）
  → common/ip_matcher_test.go

阶段 2: 数据层（依赖阶段 1 的 ValidateCIDRList）
  → model/token.go — IpPolicy 字段扩展（F002）
  → 三数据库迁移验证（SQLite/MySQL/PostgreSQL）

阶段 3: 接口层（依赖阶段 1 + 2）
  → controller/token.go — UpdateTokenIpPolicy（F001）
  → router/api-router.go — 路由注册（F001）

阶段 4: 执行层（依赖阶段 1 + 2）
  → common/init.go — TRUSTED_PROXIES 初始化（F004）
  → middleware/auth.go — IP 提取 + IP 策略校验（F004 + F005）

阶段 5: 集成测试
  → 端到端：白名单/黑名单拦截验证
  → 安全测试：XFF 伪造攻击场景（TRUSTED_PROXIES 未配置）
  → 性能测试：鉴权延迟压测（含 ParseCIDRList 开销）
```
