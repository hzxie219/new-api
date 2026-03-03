# 【Story 1.2.1】请求鉴权时获取真实客户端 IP

## 1、用户需求

**用户故事**：作为**平台网关**，我需要在处理每个 API Key 鉴权请求时，正确识别发起请求的真实客户端 IP 地址——当请求经由可信代理转发时，从 `X-Forwarded-For` 中提取真实 IP；当请求直连时，使用 TCP 连接的远端 IP。这样才能保证后续 IP 策略校验（Story 1.2.2）不被代理 IP 或伪造头欺骗。

**涉及场景**：

| 场景 | IP 来源 |
|------|--------|
| 直连（无代理） | `RemoteAddr` 中的 IP |
| 经由可信代理转发 | `X-Forwarded-For` 最左侧 IP |
| 非可信来源携带 XFF | 忽略 XFF，使用 `RemoteAddr` |
| 多层代理 | 取 XFF 最左侧 IP（最原始客户端 IP） |

---

## 2、功能性需求

### 2.1 真实客户端 IP 提取

#### 2.1.1 功能描述

在鉴权中间件（`middleware/auth.go`）中，新增统一的客户端 IP 提取逻辑。通过 `common.GetClientIP(c *gin.Context)` 函数实现：

- **可信代理判断**：检查请求的 `RemoteAddr`（TCP 直连 IP）是否在全局可信代理列表（`trustedProxyCIDRs`）中
- **IP 提取规则**：
  - 若来自可信代理 → 读取 `X-Forwarded-For` 头，取最左侧（第一个）IP
  - 若非可信代理 → 使用 `RemoteAddr` 中的 IP，忽略 XFF 头
- **可信代理配置**：通过环境变量 `TRUSTED_PROXIES` 在服务启动时配置，支持单 IP 和 CIDR，逗号分隔

#### 2.1.2 正常场景

**场景 1：直连请求（无代理）**
- **Given**：`TRUSTED_PROXIES` 未配置；请求 `RemoteAddr=8.8.8.8:12345`，无 `X-Forwarded-For` 头
- **When**：鉴权流程调用 `GetClientIP(c)`
- **Then**：返回 `"8.8.8.8"`（RemoteAddr 的 IP 部分，不含端口）

**场景 2：经由可信代理转发（单层）**
- **Given**：`TRUSTED_PROXIES=192.168.1.0/24`；请求 `RemoteAddr=192.168.1.10:54321`，`X-Forwarded-For: 1.2.3.4`
- **When**：调用 `GetClientIP(c)`
- **Then**：返回 `"1.2.3.4"`（XFF 最左侧 IP）

**场景 3：多层代理转发**
- **Given**：`TRUSTED_PROXIES=192.168.1.0/24`；请求 `RemoteAddr=192.168.1.5:443`，`X-Forwarded-For: 1.2.3.4, 5.5.5.5, 192.168.1.5`
- **When**：调用 `GetClientIP(c)`
- **Then**：返回 `"1.2.3.4"`（XFF 最左侧，即最原始客户端 IP）

**场景 4：可信代理为精确 IP 配置**
- **Given**：`TRUSTED_PROXIES=10.0.0.1`（精确 IP）；请求 `RemoteAddr=10.0.0.1:80`，`X-Forwarded-For: 2.2.2.2`
- **When**：调用 `GetClientIP(c)`
- **Then**：返回 `"2.2.2.2"`

#### 2.1.3 异常场景

**场景 5：非可信代理伪造 X-Forwarded-For**
- **Given**：`TRUSTED_PROXIES=192.168.1.0/24`；请求 `RemoteAddr=8.8.8.8:443`（不在可信列表），`X-Forwarded-For: 1.2.3.4`（伪造）
- **When**：调用 `GetClientIP(c)`
- **Then**：返回 `"8.8.8.8"`（忽略伪造的 XFF，防止 IP 绕过）

**场景 6：`TRUSTED_PROXIES` 未配置**
- **Given**：`TRUSTED_PROXIES` 环境变量未设置（空值）
- **When**：任意请求到达鉴权中间件
- **Then**：始终返回 `RemoteAddr` 中的 IP（最保守策略，不信任任何代理）

**场景 7：XFF 头为空但来自可信代理**
- **Given**：`TRUSTED_PROXIES=192.168.1.0/24`；请求 `RemoteAddr=192.168.1.3:5678`，无 `X-Forwarded-For` 头
- **When**：调用 `GetClientIP(c)`
- **Then**：返回 `"192.168.1.3"`（XFF 不存在时回退到 RemoteAddr）

#### 2.1.4 安全场景

**场景 8：攻击者通过伪造 XFF 绕过 IP 白名单**
- **Given**：API Key 设置了白名单 `["1.1.1.1/32"]`；攻击者从 `8.8.8.8` 发起请求，携带伪造头 `X-Forwarded-For: 1.1.1.1`；`TRUSTED_PROXIES` 未配置
- **When**：鉴权中间件处理请求
- **Then**：`GetClientIP()` 返回 `"8.8.8.8"`，后续 IP 策略校验失败返回 403；攻击无效

**场景 9：内部网络请求携带外网 IP 伪造 XFF**
- **Given**：`TRUSTED_PROXIES=10.0.0.0/8`；攻击者控制内网机器 `10.1.2.3`，伪造 `X-Forwarded-For: 5.5.5.5`
- **When**：鉴权中间件处理请求
- **Then**：由于 `RemoteAddr=10.1.2.3` 在可信代理列表，`GetClientIP()` 返回 `"5.5.5.5"`；后续 IP 策略按 `5.5.5.5` 校验
- **风险说明**：内网机器若被攻击者控制，可伪造来源 IP；需确保可信代理仅配置真正可信的代理服务器地址

#### 2.1.5 边界场景

**场景 10：RemoteAddr 包含 IPv6 地址**
- **Given**：请求 `RemoteAddr=[::1]:12345`（IPv6 本地地址）
- **When**：调用 `GetClientIP(c)`
- **Then**：正确解析 IP 部分 `"::1"`，不含端口，不出现解析错误

**场景 11：RemoteAddr 格式异常（无端口）**
- **Given**：`RemoteAddr=8.8.8.8`（无端口号，格式异常）
- **When**：调用 `GetClientIP(c)`
- **Then**：回退使用完整 `RemoteAddr` 字符串作为 IP，记录 WARN 日志

**场景 12：XFF 头值含多余空格**
- **Given**：`X-Forwarded-For: " 1.2.3.4 , 5.5.5.5"`（含空格）
- **When**：调用 `GetClientIP(c)`
- **Then**：正确 Trim 后返回 `"1.2.3.4"`

#### 2.1.6 约束条件

| 约束 | 规格 |
|------|------|
| `TRUSTED_PROXIES` 格式 | 逗号分隔的精确 IP 或 CIDR 列表，如 `127.0.0.1,192.168.0.0/16` |
| `TRUSTED_PROXIES` 默认值 | 空（不信任任何代理），运维需显式配置 |
| XFF 取值规则 | 取最左侧（第一个非空）IP，代表最原始客户端 |
| IPv6 支持 | `GetClientIP` 须正确解析 IPv6 格式的 RemoteAddr，CIDR 匹配当前版本仅保证 IPv4 |
| 函数调用时机 | 在 Token 有效性校验（存在、启用、有效期、额度）通过后调用，用于 IP 策略校验 |
| 全局变量初始化 | `trustedProxyCIDRs` 在服务启动时（`common/init.go`）初始化，运行时不可动态修改 |

#### 2.1.7 依赖与风险

| 项 | 说明 |
|----|------|
| 依赖 `common/init.go` | 解析并初始化 `trustedProxyCIDRs` 全局变量 |
| 依赖 `common/ip_matcher.go` | `GetClientIP`、`InitTrustedProxies`、`isTrustedProxy` 函数 |
| 被 Story 1.2.2 依赖 | 提供 `GetClientIP()` 结果供 IP 策略校验使用 |
| 风险：TRUSTED_PROXIES 配置错误 | 配置过宽（如 `0.0.0.0/0`）导致 XFF 被任意伪造；部署文档须明确警示 |

---

## 3、非功能性需求

### 3.1 性能分析

| 指标 | 目标值 | 说明 |
|------|--------|------|
| `GetClientIP` 函数耗时 | < 0.1ms | 纯内存操作（CIDR 匹配 + 字符串解析），无 IO |
| `InitTrustedProxies` 初始化耗时 | < 10ms | 启动时一次性执行，不影响请求处理 |
| 内存占用（`trustedProxyCIDRs`） | < 1KB（< 50 个 CIDR 条目） | 全局变量，常驻内存，量级极小 |

### 3.2 可靠性

- `GetClientIP` 不依赖数据库或外部服务，无 IO 操作，本身不存在可用性风险
- `RemoteAddr` 格式解析失败时降级处理（返回原始字符串），不 panic
- `InitTrustedProxies` 中非法 CIDR 条目跳过并记录 WARN 日志，不影响服务启动

### 3.3 可维护性

- `GetClientIP`、`isTrustedProxy`、`InitTrustedProxies` 集中在 `common/ip_matcher.go`，便于单独测试和维护
- `TRUSTED_PROXIES` 配置变更需重启服务（无热加载），部署文档须说明

### 3.4 安全性

| 安全需求 | 措施 |
|---------|------|
| 防 XFF 伪造 | 仅信任来自可信代理的 XFF 头，直连请求忽略 XFF |
| 默认安全 | `TRUSTED_PROXIES` 默认为空（不信任任何代理），需显式配置 |
| 可信代理最小化 | 仅将实际部署的 Nginx/LB IP 加入列表，避免配置过宽 |

### 3.5 可测试性

- `GetClientIP` 有独立单元测试，覆盖直连/可信代理/XFF 伪造/多层代理/IPv6 等场景
- `InitTrustedProxies` 有单元测试，覆盖合法 CIDR、非法 CIDR 跳过、空配置等场景
- 测试文件：`new-api/common/ip_matcher_test.go`
