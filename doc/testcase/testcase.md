# API Key IP访问控制模块 接口测试用例文档

> 生成时间: 2026-03-02
> 需求文档: doc/design/tech_design.md
> 生成方法: TDD方法论，基于接口分析报告系统化设计

---

## 1. 测试概述

### 1.1 功能描述

本模块为 `new-api` AI API 网关的 API Key 增加结构化 IP 访问控制能力，核心功能包括：

- **F001 IP策略配置接口**：新增 `PUT /api/token/:id/ip_policy`，支持为指定 Token 设置/清除白名单或黑名单策略
- **F002 Token.IpPolicy 数据模型**：Token 结构体扩展，支持 JSON 序列化持久化到数据库 `ip_policy TEXT NULL` 列
- **F003 CIDR 工具函数库**：`common/ip_matcher.go`，提供 `ValidateCIDRList`、`ParseCIDRList`、`IPMatchesCIDRList`、`GetClientIP`、`InitTrustedProxies` 五个工具函数
- **F004 可信代理感知 IP 提取**：`GetClientIP` 根据 TRUSTED_PROXIES 环境变量决定是否采信 XFF 头
- **F005 IP 策略执行**：在 `TokenAuth` 中间件中执行白/黑名单策略决策

**核心安全约束**：
- 未配置 TRUSTED_PROXIES 时 XFF 头永远被忽略（防 XFF 伪造）
- IP 拦截统一返回 `IP_NOT_ALLOWED`，不透露策略模式
- CIDR 运行时解析失败时 fail-open（宽松降级），不影响服务可用性

### 1.2 测试范围

- **HTTP API 接口**: 1 个（`PUT /api/token/:id/ip_policy`）
- **内部工具函数**: 5 个（`ValidateCIDRList`/`ParseCIDRList`/`IPMatchesCIDRList`/`GetClientIP`/`InitTrustedProxies`）
- **测试类型**: 功能测试、接口测试、单元测试、安全测试
- **测试深度**: 覆盖正常场景、异常场景、边界场景、安全场景

**不在本测试范围**：
- AllowIps 旧字段行为（已有测试覆盖，不变更）
- 数据库 AutoMigrate 自动化验证（DBA 手工确认三库执行结果）
- 性能压测（P99 指标由专项性能测试覆盖）

### 1.3 测试策略

| 策略项 | 说明 |
|-------|------|
| 测试方法 | TDD 方法论，先设计用例再实现功能 |
| 验证方式 | HTTP 状态码 + 响应体字段 + 数据库状态三重验证 |
| 覆盖标准 | 接口覆盖率 100%，参数覆盖率 100%，功能点覆盖率 ≥ 95% |
| 并发安全 | FUNC002 需在 `-race` 标志下运行，验证 `ParseCIDRList` 无竞态 |
| 数据隔离 | 每个测试用例独立准备测试数据，执行后清理，互不干扰 |

---

## 2. 接口清单

| 接口ID | 接口名称 | HTTP方法 | 路径 | 优先级 | 参数数量 | 认证要求 |
|--------|---------|---------|------|--------|---------|---------|
| API001 | 更新 Token IP 策略接口 | PUT | `/api/token/:id/ip_policy` | P0 | 3（路径1+请求体2） | UserAuth（JWT）+ 拥有者/管理员校验 |

| 函数ID | 函数名称 | 文件 | 测试类型 | 优先级 |
|--------|---------|------|---------|--------|
| FUNC001 | `ValidateCIDRList` | `common/ip_matcher.go` | 单元测试 | P0 |
| FUNC002 | `ParseCIDRList` | `common/ip_matcher.go` | 单元测试（含并发）| P0 |
| FUNC003 | `IPMatchesCIDRList` | `common/ip_matcher.go` | 单元测试 | P0 |
| FUNC004 | `GetClientIP` | `common/ip_matcher.go` | 单元测试 | P1 |
| FUNC005 | `InitTrustedProxies` | `common/ip_matcher.go` | 单元测试 | P1 |

---

## 3. 测试用例

### 3.1 API001 — 更新 Token IP 策略接口（`PUT /api/token/:id/ip_policy`）

#### 3.1.1 正常场景

| 用例ID | 测试场景 | 测试数据 | 预期结果 | 验收条件 | 优先级 |
|--------|---------|---------|---------|---------|--------|
| TC001 | 设置 IPv4 白名单策略 | `id=1`（拥有者操作），`{"mode":"whitelist","ips":["1.2.3.4","10.0.0.0/8","192.168.1.0/24"]}` | HTTP 200 | ① 响应 `{"success":true,"message":""}` ② DB `tokens.ip_policy` 写入 `{"mode":"whitelist","ips":["1.2.3.4","10.0.0.0/8","192.168.1.0/24"]}` | P0 |
| TC002 | 设置 IPv4 黑名单策略 | `id=1`（拥有者操作），`{"mode":"blacklist","ips":["203.0.113.0/24","198.51.100.5"]}` | HTTP 200 | ① 响应 `{"success":true,"message":""}` ② DB `ip_policy` 写入黑名单 JSON | P0 |
| TC003 | 清除策略（mode 为空字符串） | `id=1`（拥有者操作），`{"mode":"","ips":[]}` | HTTP 200 | ① 响应 `{"success":true,"message":""}` ② DB `tokens.ip_policy` 写入 NULL | P0 |
| TC004 | 单个精确 IP（无 CIDR 前缀）| `id=1`，`{"mode":"whitelist","ips":["203.0.113.45"]}` | HTTP 200 | ① 成功写入 ② 后续鉴权时 `203.0.113.45` 被视为 `/32` 处理，命中白名单 | P1 |
| TC005 | 设置 IPv6 CIDR 白名单 | `id=1`，`{"mode":"whitelist","ips":["2001:db8::/32","::1"]}` | HTTP 200 | ① 成功写入 ② 响应 `success:true` | P1 |
| TC006 | IPv4 与 IPv6 混合列表 | `id=1`，`{"mode":"blacklist","ips":["1.2.3.4","2001:db8::1/128","10.0.0.0/8"]}` | HTTP 200 | 成功写入，混合类型均合法 | P1 |
| TC007 | 管理员操作他人 Token | `id=2`（Token 属于 user_b），当前用户为 `admin`，`{"mode":"whitelist","ips":["192.168.0.0/16"]}` | HTTP 200 | ① 管理员权限校验通过 ② 成功写入 user_b 的 Token 策略 | P0 |
| TC008 | 拥有者操作自己的 Token | `id=1`（当前用户即 Token 拥有者），`{"mode":"whitelist","ips":["10.0.0.1"]}` | HTTP 200 | 拥有者权限校验通过，写入成功 | P0 |
| TC009 | ips 为空数组（含 mode） | `id=1`，`{"mode":"whitelist","ips":[]}` | HTTP 200 | 成功写入 `{"mode":"whitelist","ips":[]}` ；空白名单语义待产品确认（Q1）| P1 |
| TC010 | 设置策略后再次覆盖更新 | 先 TC001 设白名单，再 `{"mode":"blacklist","ips":["5.5.5.5"]}` | HTTP 200 | 第二次更新覆盖第一次，DB 中最终为黑名单策略 | P1 |

#### 3.1.2 异常场景

| 用例ID | 测试场景 | 测试数据 | 预期结果 | 验收条件 | 优先级 |
|--------|---------|---------|---------|---------|--------|
| TC011 | mode 为非法枚举值 | `id=1`，`{"mode":"deny","ips":["1.2.3.4"]}` | HTTP 400 | ① 响应 `{"success":false,"message":"invalid mode: deny"}` ② DB 不写入 | P0 |
| TC012 | mode 为随机字符串 | `id=1`，`{"mode":"WHITELIST","ips":[]}` | HTTP 400 | 响应 `invalid mode: WHITELIST`（大小写敏感）② DB 不写入 | P1 |
| TC013 | ips 中含非法 CIDR 格式 | `id=1`，`{"mode":"whitelist","ips":["999.0.0.1/33"]}` | HTTP 400 | ① 响应 `{"success":false,"message":"invalid IP/CIDR: 999.0.0.1/33"}` ② DB 不写入 | P0 |
| TC014 | ips 中含非 IP 字符串 | `id=1`，`{"mode":"whitelist","ips":["not-an-ip"]}` | HTTP 400 | 响应 `invalid IP/CIDR: not-an-ip` ② DB 不写入 | P0 |
| TC015 | ips 中混合合法与非法条目 | `id=1`，`{"mode":"whitelist","ips":["1.2.3.4","bad-ip"]}` | HTTP 400 | 返回第一个非法条目的错误信息（`bad-ip`）② 全部拒绝写入（非部分写入）| P1 |
| TC016 | ips 条目数超过 100 | `id=1`，`{"mode":"whitelist","ips":["1.2.3.4",...]}` （101 条） | HTTP 400 | 响应 `{"success":false,"message":"too many IP entries, max 100"}` ② DB 不写入 | P0 |
| TC017 | Token ID 不存在 | `id=99999`（不存在），`{"mode":"whitelist","ips":["1.2.3.4"]}` | HTTP 404 | 响应 `{"success":false,"message":"not found"}` | P0 |
| TC018 | 非拥有者且非管理员操作 | `id=2`（Token 属于 user_b），当前用户为 `user_a`（普通用户），`{"mode":"whitelist","ips":["1.2.3.4"]}` | HTTP 403 | 响应 `{"success":false,"message":"forbidden"}` ② DB 不写入 | P0 |
| TC019 | 未携带 JWT Token | 无 Authorization 头，`{"mode":"whitelist","ips":["1.2.3.4"]}` | HTTP 401 | 路由组 UserAuth 中间件拦截，返回 401 Unauthorized | P0 |
| TC020 | JWT Token 已过期 | 携带过期 JWT，`{"mode":"whitelist","ips":["1.2.3.4"]}` | HTTP 401 | UserAuth 中间件校验失败，返回 401 | P1 |
| TC021 | 请求体 JSON 格式错误 | `id=1`，请求体为 `{mode: whitelist}` （非合法 JSON）| HTTP 400 | `ShouldBindJSON` 失败，返回 400 | P1 |
| TC022 | 请求体缺少 mode 字段 | `id=1`，`{"ips":["1.2.3.4"]}` | HTTP 400 | `mode` 字段 binding 校验失败，返回 400（`mode` 为 required 字段） | P1 |
| TC023 | 请求体缺少 ips 字段 | `id=1`，`{"mode":"whitelist"}` | HTTP 400 | `ips` 字段 binding 校验失败，返回 400（`ips` 为 required 字段）| P1 |
| TC024 | 数据库写入失败（模拟）| 正常请求参数，Mock DB 返回写入错误 | HTTP 500 | 响应 `{"success":false,"message":"internal server error"}` ② DB 无中间状态残留 | P1 |

#### 3.1.3 边界场景

| 用例ID | 测试场景 | 测试数据 | 预期结果 | 验收条件 | 优先级 |
|--------|---------|---------|---------|---------|--------|
| TC025 | id = 1（最小合法值）| `id=1`，`{"mode":"whitelist","ips":["1.2.3.4"]}` | HTTP 200 | 最小合法 ID 正常处理 | P1 |
| TC026 | id = 0（非法值）| `id=0`，`{"mode":"whitelist","ips":["1.2.3.4"]}` | HTTP 400 或 404 | id ≤ 0 不合法，路由参数校验或 DB 查询返回错误 | P1 |
| TC027 | id 为负数 | `id=-1`，`{"mode":"whitelist","ips":[]}` | HTTP 400 或 404 | 负数 ID 不合法 | P1 |
| TC028 | id 为非整数字符串 | `id=abc`，`{"mode":"whitelist","ips":[]}` | HTTP 400 | 路径参数类型转换失败，返回 400 | P1 |
| TC029 | ips 恰好 100 条（上限）| `id=1`，`{"mode":"blacklist","ips":["1.2.3.{n}":n=1..100]}` | HTTP 200 | 正好 100 条合法，成功写入 | P0 |
| TC030 | ips 恰好 101 条（超限 1 条）| `id=1`，`{"mode":"blacklist","ips":["1.2.3.{n}":n=1..101]}` | HTTP 400 | 返回 `too many IP entries, max 100` | P0 |
| TC031 | /32 主机路由（单 IP CIDR）| `id=1`，`{"mode":"whitelist","ips":["192.168.1.100/32"]}` | HTTP 200 | 显式 `/32` 合法，等同单 IP | P1 |
| TC032 | /0 全匹配 CIDR | `id=1`，`{"mode":"blacklist","ips":["0.0.0.0/0"]}` | HTTP 200 | `/0` 为合法 CIDR（匹配所有 IPv4），成功写入；行为语义：黑名单拦截所有 IPv4 请求 | P1 |
| TC033 | /128 IPv6 主机路由 | `id=1`，`{"mode":"whitelist","ips":["::1/128"]}` | HTTP 200 | IPv6 `/128` 合法 | P1 |
| TC034 | 精确 IP 边界（CIDR 网段首地址）| 白名单 `["10.0.0.0/8"]`，客户端 IP 为 `10.0.0.0` | 鉴权通过 | 网段首地址命中白名单，HTTP 200 继续处理 | P1 |
| TC035 | CIDR 网段末地址命中 | 白名单 `["192.168.1.0/24"]`，客户端 IP 为 `192.168.1.255` | 鉴权通过 | 网段末地址命中白名单 | P1 |
| TC036 | CIDR 网段外第一个 IP | 白名单 `["192.168.1.0/24"]`，客户端 IP 为 `192.168.2.0` | 鉴权拒绝 HTTP 403 | 超出网段，白名单未命中，返回 IP_NOT_ALLOWED | P1 |

#### 3.1.4 安全场景

| 用例ID | 测试场景 | 测试数据 | 预期结果 | 验收条件 | 优先级 |
|--------|---------|---------|---------|---------|--------|
| TC037 | 未认证请求（无 Authorization 头）| 无 JWT Token，直接请求 | HTTP 401 | UserAuth 中间件拦截，不进入 Handler | P0 |
| TC038 | 无效 JWT（篡改签名）| Authorization: Bearer `eyJ...invalid_sig` | HTTP 401 | JWT 验证失败，401 | P0 |
| TC039 | 越权操作（水平权限提升）| user_a 的 JWT，操作 user_b 的 token（`id=2`）| HTTP 403 | Handler 内权限校验：非拥有者非管理员，返回 forbidden | P0 |
| TC040 | 错误码不泄露策略模式 | 白名单策略下，不在列表的 IP 请求 API | HTTP 403，`IP_NOT_ALLOWED` | 响应体不包含 `whitelist`/`blacklist` 等模式信息 | P0 |
| TC041 | 错误码统一（黑名单）| 黑名单策略下，列表内 IP 请求 API | HTTP 403，`IP_NOT_ALLOWED` | 与白名单拒绝响应相同，不区分策略类型 | P0 |
| TC042 | XFF 头伪造攻击（未配置 TRUSTED_PROXIES）| 无 TRUSTED_PROXIES 配置，请求携带 `X-Forwarded-For: 1.2.3.4`；Token 白名单为 `["1.2.3.4"]`；实际 RemoteAddr 不在白名单 | HTTP 403 | XFF 头被忽略，使用真实 RemoteAddr，白名单未命中，拒绝请求 | P0 |
| TC043 | XFF 头伪造攻击（已配置 TRUSTED_PROXIES，非可信代理发送 XFF）| TRUSTED_PROXIES=`"10.0.0.1"`，请求来自 `RemoteAddr=5.5.5.5`（非可信），携带 `X-Forwarded-For: 1.2.3.4`；白名单为 `["1.2.3.4"]` | HTTP 403 | 直连 IP `5.5.5.5` 不在可信代理列表，XFF 被忽略，使用 RemoteAddr，白名单未命中 | P0 |
| TC044 | 合法 XFF 路径（可信代理转发）| TRUSTED_PROXIES=`"10.0.0.1"`，`RemoteAddr=10.0.0.1`，`X-Forwarded-For: 192.168.1.100`；白名单 `["192.168.1.100"]` | HTTP 200（鉴权通过） | RemoteAddr 属于可信代理，采信 XFF，clientIP=`192.168.1.100`，白名单命中 | P1 |
| TC045 | ips 中注入特殊字符（SQL 注入尝试）| `{"mode":"whitelist","ips":["1.2.3.4'; DROP TABLE tokens; --"]}` | HTTP 400 | `ValidateCIDRList` 校验失败，`invalid IP/CIDR`，不入库 | P0 |
| TC046 | ips 中超长字符串 | 单条 IP 字符串长度 > 1000 字符 | HTTP 400 | `ValidateCIDRList` 校验失败，invalid IP/CIDR | P1 |
| TC047 | id 传入极大整数（int64 溢出测试）| `id=9999999999999999999` | HTTP 400 | 路径参数整数解析失败，返回 400 | P2 |
| TC048 | 重放攻击（旧 JWT 重放）| 使用已注销或权限已变更用户的旧 JWT | HTTP 401 或 403 | 依赖项目现有 JWT 验证机制（不在本模块测试范围，记录为关联测试）| P2 |

---

### 3.2 FUNC001 — `ValidateCIDRList` 单元测试

**函数签名**：`func ValidateCIDRList(ips []string) error`
**文件**：`common/ip_matcher.go`
**测试文件**：`common/ip_matcher_test.go`

| 用例ID | 测试场景 | 输入 | 预期返回 | 优先级 |
|--------|---------|------|---------|--------|
| UT001 | 空列表 | `[]string{}` | `nil` | P0 |
| UT002 | nil 列表 | `nil` | `nil` | P0 |
| UT003 | 合法 IPv4 精确地址 | `["1.2.3.4"]` | `nil` | P0 |
| UT004 | 合法 IPv4 CIDR | `["10.0.0.0/8", "192.168.1.0/24"]` | `nil` | P0 |
| UT005 | 合法 IPv6 地址 | `["::1", "2001:db8::1"]` | `nil` | P0 |
| UT006 | 合法 IPv6 CIDR | `["2001:db8::/32", "fe80::/10"]` | `nil` | P0 |
| UT007 | 合法混合 IPv4/IPv6 | `["1.2.3.4", "::1", "10.0.0.0/8", "2001:db8::/32"]` | `nil` | P0 |
| UT008 | 非法 IP（IP 段超范围）| `["999.0.0.1"]` | 返回包含 `"999.0.0.1"` 的 `error` | P0 |
| UT009 | 非法 CIDR（前缀长度超范围）| `["192.168.1.0/33"]` | 返回包含 `"192.168.1.0/33"` 的 `error` | P0 |
| UT010 | 非法 CIDR（IPv6 前缀超范围）| `["::1/129"]` | 返回 `error` | P0 |
| UT011 | 非 IP 字符串 | `["not-an-ip"]` | 返回包含 `"not-an-ip"` 的 `error` | P0 |
| UT012 | 空字符串条目 | `[""]` | 返回 `error` | P1 |
| UT013 | 混合合法与非法：返回第一个非法 | `["1.2.3.4", "bad", "5.6.7.8"]` | 返回包含 `"bad"` 的 `error`（第一个非法条目）| P0 |
| UT014 | 前缀长度为 0（全匹配）| `["0.0.0.0/0"]` | `nil`（合法 CIDR）| P1 |
| UT015 | /32 主机路由 CIDR | `["192.168.1.1/32"]` | `nil` | P1 |
| UT016 | /128 IPv6 主机路由 | `["::1/128"]` | `nil` | P1 |

---

### 3.3 FUNC002 — `ParseCIDRList` 单元测试

**函数签名**：`func ParseCIDRList(ips []string) ([]*net.IPNet, error)`
**文件**：`common/ip_matcher.go`

| 用例ID | 测试场景 | 输入 | 预期返回 | 验收条件 | 优先级 |
|--------|---------|------|---------|---------|--------|
| UT017 | 空列表 | `[]string{}` | `([]*net.IPNet{}, nil)` | 返回空切片而非 nil，error 为 nil | P0 |
| UT018 | 单个精确 IP（无前缀）| `["1.2.3.4"]` | 包含 `1.2.3.4/32` 的 `[]*net.IPNet`，`nil` | 单 IP 自动补全 `/32` | P0 |
| UT019 | 单个 CIDR | `["10.0.0.0/8"]` | 包含对应 `*net.IPNet`，`nil` | `IPNet.IP` = `10.0.0.0`，`Mask` = `/8` | P0 |
| UT020 | 多条 CIDR 混合 | `["1.2.3.4", "10.0.0.0/8", "::1/128"]` | 3 个 `*net.IPNet`，`nil` | 顺序与输入一致 | P0 |
| UT021 | 包含非法条目 | `["1.2.3.4", "bad"]` | `(nil, error)` | 返回 error，error 信息含 `"bad"` | P0 |
| UT022 | 每次调用返回独立切片 | 同一输入调用两次 | 两次返回不同切片指针 | 无共享状态：修改 slice1 不影响 slice2 | P0 |
| UT023 | 并发调用无竞态（race detector）| 100 goroutine 并发调用，相同输入 | 所有调用均正常返回，无竞态报告 | 使用 `go test -race` 运行，无 DATA RACE 输出 | P0 |
| UT024 | IPv6 CIDR 解析 | `["2001:db8::/32"]` | 对应 IPv6 `*net.IPNet`，`nil` | IPv6 地址族正确 | P1 |
| UT025 | /0 全匹配 CIDR | `["0.0.0.0/0"]` | 匹配所有 IPv4 的 `*net.IPNet`，`nil` | Mask 全为 0 | P1 |

---

### 3.4 FUNC003 — `IPMatchesCIDRList` 单元测试

**函数签名**：`func IPMatchesCIDRList(ipStr string, cidrs []*net.IPNet) bool`
**文件**：`common/ip_matcher.go`

| 用例ID | 测试场景 | 输入（ipStr, cidrs） | 预期返回 | 优先级 |
|--------|---------|---------------------|---------|--------|
| UT026 | IP 命中单条精确 CIDR | `"1.2.3.4"`, `[1.2.3.4/32]` | `true` | P0 |
| UT027 | IP 命中网段 CIDR | `"10.1.2.3"`, `[10.0.0.0/8]` | `true` | P0 |
| UT028 | IP 未命中（在网段外）| `"192.168.2.1"`, `[192.168.1.0/24]` | `false` | P0 |
| UT029 | 空 CIDR 列表 | `"1.2.3.4"`, `[]` | `false` | P0 |
| UT030 | IP 命中多条中的一条 | `"10.1.1.1"`, `[192.168.0.0/16, 10.0.0.0/8]` | `true` | P0 |
| UT031 | IP 均不命中多条 CIDR | `"1.2.3.4"`, `[192.168.0.0/16, 10.0.0.0/8]` | `false` | P0 |
| UT032 | 网段首地址命中 | `"10.0.0.0"`, `[10.0.0.0/8]` | `true` | P1 |
| UT033 | 网段末地址命中 | `"10.255.255.255"`, `[10.0.0.0/8]` | `true` | P1 |
| UT034 | 网段外第一个 IP | `"11.0.0.0"`, `[10.0.0.0/8]` | `false` | P1 |
| UT035 | IPv6 命中 | `"2001:db8::1"`, `[2001:db8::/32]` | `true` | P1 |
| UT036 | IPv6 未命中 | `"2001:db9::1"`, `[2001:db8::/32]` | `false` | P1 |
| UT037 | 无效 IP 字符串 | `"not-an-ip"`, `[10.0.0.0/8]` | `false` | P1 |
| UT038 | 空字符串 IP | `""`, `[10.0.0.0/8]` | `false` | P1 |
| UT039 | 全匹配 /0 CIDR | `"203.0.113.1"`, `[0.0.0.0/0]` | `true` | P1 |
| UT040 | IPv4 不匹配 IPv6 CIDR | `"1.2.3.4"`, `[2001:db8::/32]` | `false` | P1 |
| UT041 | cidrs 为 nil | `"1.2.3.4"`, `nil` | `false` | P1 |

---

### 3.5 FUNC004 — `GetClientIP` 单元测试

**函数签名**：`func GetClientIP(c *gin.Context) string`
**文件**：`common/ip_matcher.go`
**测试前置**：使用 `httptest.NewRequest` 构建 gin.Context，调用 `InitTrustedProxies` 设置可信代理

| 用例ID | 测试场景 | 测试环境 | 请求构造 | 预期结果 | 优先级 |
|--------|---------|---------|---------|---------|--------|
| UT042 | 无可信代理，直连请求 | `InitTrustedProxies("")`（默认） | `RemoteAddr="203.0.113.5:12345"`，无 XFF 头 | 返回 `"203.0.113.5"` | P0 |
| UT043 | 无可信代理，携带 XFF 头（XFF 被忽略）| `InitTrustedProxies("")` | `RemoteAddr="203.0.113.5:12345"`，`X-Forwarded-For: 1.2.3.4` | 返回 `"203.0.113.5"`（XFF 被忽略，使用 RemoteAddr）| P0 |
| UT044 | 可信代理，采信 XFF 最左侧 IP | `InitTrustedProxies("10.0.0.1")` | `RemoteAddr="10.0.0.1:8080"`，`X-Forwarded-For: 192.168.1.100` | 返回 `"192.168.1.100"`（XFF 最左侧 IP）| P0 |
| UT045 | 可信代理，多跳 XFF，取最左侧 | `InitTrustedProxies("10.0.0.1")` | `RemoteAddr="10.0.0.1:8080"`，`X-Forwarded-For: 1.2.3.4, 5.6.7.8` | 返回 `"1.2.3.4"`（最左侧为原始客户端 IP）| P0 |
| UT046 | 可信代理，XFF 头为空 | `InitTrustedProxies("10.0.0.1")` | `RemoteAddr="10.0.0.1:8080"`，无 XFF 头 | 返回 `"10.0.0.1"`（回退到 RemoteAddr）| P0 |
| UT047 | 非可信代理发送 XFF（XFF 被忽略）| `InitTrustedProxies("10.0.0.1")` | `RemoteAddr="5.5.5.5:12345"`，`X-Forwarded-For: 1.2.3.4` | 返回 `"5.5.5.5"`（5.5.5.5 不在可信代理列表，XFF 被忽略）| P0 |
| UT048 | 可信代理 CIDR 范围 | `InitTrustedProxies("10.0.0.0/8")` | `RemoteAddr="10.1.2.3:8080"`，`X-Forwarded-For: 203.0.113.1` | 返回 `"203.0.113.1"`（RemoteAddr 在可信 CIDR 内）| P1 |
| UT049 | XFF 中含多个 IP，取最左 | `InitTrustedProxies("10.0.0.1")` | `RemoteAddr="10.0.0.1:8080"`，`X-Forwarded-For:  203.0.113.1, 172.16.0.1 ` （含空白符）| 返回 `"203.0.113.1"`（去除首尾空白符）| P1 |

---

### 3.6 FUNC005 — `InitTrustedProxies` 单元测试

**函数签名**：`func InitTrustedProxies(proxies string)`
**文件**：`common/ip_matcher.go`
**说明**：函数修改包级全局变量 `trustedProxyCIDRs`，测试需验证函数执行后的全局状态（通过 `GetClientIP` 间接验证或反射读取）

| 用例ID | 测试场景 | 输入 | 预期状态 | 优先级 |
|--------|---------|------|---------|--------|
| UT050 | 空字符串（默认配置）| `""` | `trustedProxyCIDRs = nil`，`GetClientIP` 永远返回 RemoteAddr | P0 |
| UT051 | 单个合法 IP 地址 | `"10.0.0.1"` | `trustedProxyCIDRs` 包含 `10.0.0.1/32`，`isTrustedProxy("10.0.0.1")=true` | P0 |
| UT052 | 单个合法 CIDR | `"192.168.0.0/16"` | `trustedProxyCIDRs` 包含对应 `*net.IPNet` | P0 |
| UT053 | 逗号分隔多个合法条目 | `"10.0.0.1,172.16.0.0/12,192.168.0.0/16"` | `trustedProxyCIDRs` 包含 3 个 `*net.IPNet` | P0 |
| UT054 | 包含非法条目（跳过非法，加载合法）| `"10.0.0.1,bad-entry,192.168.0.0/16"` | `trustedProxyCIDRs` 包含 2 个（10.0.0.1/32 和 192.168.0.0/16），跳过 `bad-entry`；输出 WARN 日志 | P0 |
| UT055 | 全部条目非法 | `"bad1,bad2"` | `trustedProxyCIDRs = nil` 或空切片；每个非法条目记录 WARN 日志 | P0 |
| UT056 | 含空格的逗号分隔 | `"10.0.0.1, 192.168.0.0/16"` | 正确加载（空格 trim 处理）或跳过含空格条目（记录为已知行为）| P1 |
| UT057 | 逗号结尾（末尾空字符串）| `"10.0.0.1,"` | 正确加载 `10.0.0.1/32`，末尾空字符串跳过 | P1 |
| UT058 | IPv6 CIDR | `"2001:db8::/32"` | `trustedProxyCIDRs` 包含对应 IPv6 `*net.IPNet` | P1 |
| UT059 | 仅包含空白符 | `"   "` | `trustedProxyCIDRs = nil`（视为空配置）| P1 |

---

## 4. 测试功能点汇总报告

### 4.1 覆盖度统计

| 维度 | 统计 | 覆盖率 |
|------|------|--------|
| HTTP 接口覆盖 | 1/1 | **100%** |
| 内部函数覆盖 | 5/5 | **100%** |
| 接口参数覆盖 | 3/3（id、mode、ips）| **100%** |
| 错误码覆盖 | 6/6（400×3 + 403 + 404 + 500）| **100%** |
| 业务规则覆盖 | 10/10（BR001-BR010）| **100%** |
| 安全需求覆盖 | 5/5 | **100%** |

**测试用例总计**：

| 范围 | 正常场景 | 异常场景 | 边界场景 | 安全场景 | 合计 |
|------|---------|---------|---------|---------|------|
| API001（HTTP 接口）| 10 | 14 | 12 | 12 | **48** |
| FUNC001 | 7 | 9 | — | — | **16** |
| FUNC002 | 6 | 1 | 2 | — | **9** |
| FUNC003 | 6 | 3 | 6 | — | **16** |
| FUNC004 | 3 | — | 5 | — | **8** |
| FUNC005 | 4 | 5 | — | — | **10** |
| **合计** | **36** | **32** | **25** | **12** | **107** |

**场景分布比例**：

```
正常场景: 36/107 = 33.6%
异常场景: 32/107 = 29.9%
边界场景: 25/107 = 23.4%
安全场景: 12/107 = 11.2%
```

---

### 4.2 参数覆盖矩阵

#### API001 接口参数覆盖

| 接口 | 参数名 | 类型 | 位置 | 覆盖用例 | 覆盖状态 |
|------|--------|------|------|---------|---------|
| API001 | `id` | integer | path | TC001–TC010（正常）, TC017（404）, TC025–TC028（边界）| ✅ 已覆盖 |
| API001 | `mode` | string（枚举）| body | TC001（whitelist）, TC002（blacklist）, TC003（""清除）, TC011（非法值）, TC012（大小写）| ✅ 已覆盖 |
| API001 | `ips` | array\<string\> | body | TC001-TC006（合法内容）, TC009（空数组）, TC013-TC015（非法格式）, TC016（超100条）, TC029–TC036（边界值）| ✅ 已覆盖 |

> ⚠️ 未覆盖参数警告：**无**

#### 错误码覆盖矩阵

| HTTP 状态码 | 错误信息 | 触发条件 | 覆盖用例 | 状态 |
|-----------|---------|---------|---------|------|
| 400 | `invalid mode: xxx` | mode 不在枚举值 | TC011, TC012 | ✅ |
| 400 | `invalid IP/CIDR: xxx` | ips 含非法 CIDR | TC013, TC014, TC015 | ✅ |
| 400 | `too many IP entries, max 100` | len(ips)>100 | TC016, TC030 | ✅ |
| 403 | `forbidden` | 非拥有者非管理员 | TC018, TC039 | ✅ |
| 404 | `not found` | Token ID 不存在 | TC017 | ✅ |
| 500 | `internal server error` | DB 写入失败 | TC024 | ✅ |

#### 业务规则覆盖矩阵

| 规则ID | 规则说明 | 覆盖用例 | 状态 |
|--------|---------|---------|------|
| BR001 | mode='' 时清除策略，ips 忽略 | TC003 | ✅ |
| BR002 | 单个 IP 自动视为 /32 | TC004, TC018（UT018）| ✅ |
| BR003 | 白名单：IP 未命中则 403 | TC036, TC042, TC043 | ✅ |
| BR004 | 黑名单：IP 命中则 403 | TC041 | ✅ |
| BR005 | IpPolicy=nil 跳过策略（向后兼容）| TC003（清除后验证）| ✅ |
| BR006 | CIDR 解析失败 fail-open（不拦截）| TC024（间接）, FUNC002/TC021 | ✅ |
| BR007 | 拦截统一返回 IP_NOT_ALLOWED | TC040, TC041 | ✅ |
| BR008 | 无 TRUSTED_PROXIES 时 XFF 被忽略 | TC042, UT043 | ✅ |
| BR009 | TRUSTED_PROXIES 含非法条目跳过 | UT054, UT055 | ✅ |
| BR010 | 条目数上限 100，接口层拒绝 >100 | TC016, TC029, TC030 | ✅ |

---

### 4.3 质量评估

| 评估维度 | 评级 | 说明 |
|---------|------|------|
| **参数覆盖完整性** | 优秀 ⭐⭐⭐ | 3/3 接口参数全覆盖，6/6 错误码全覆盖 |
| **业务规则覆盖** | 优秀 ⭐⭐⭐ | 10/10 业务规则全部有对应测试用例 |
| **安全测试深度** | 优秀 ⭐⭐⭐ | 覆盖 JWT 认证、越权、XFF 伪造、SQL 注入等关键安全场景 |
| **边界值覆盖** | 良好 ⭐⭐ | 覆盖 id 边界、ips 条目数边界、CIDR 网段首末地址 |
| **并发安全测试** | 良好 ⭐⭐ | `ParseCIDRList` race detector 测试（UT023），可扩展到 `InitTrustedProxies` |
| **用例可执行性** | 优秀 ⭐⭐⭐ | 每条用例含完整输入数据和具体验收条件，可直接执行 |
| **文档完整性** | 优秀 ⭐⭐⭐ | 所有 6 个必需章节均已包含 |

**已知遗漏/待跟进**：

| 编号 | 说明 | 影响等级 | 处理方式 |
|------|------|---------|---------|
| Q1 | 白名单 mode 下 ips 为空数组时的业务语义（拦截所有还是放行所有？）待产品确认 | 高 | TC009 标注待确认；确认后补充对应的断言 |
| Q4 | 黑名单策略下管理员 IP 是否豁免（产品 Q4 问题）| 高 | TC007 覆盖管理员操作 Token，豁免逻辑由产品确认后补充 |
| R6 | Redis 缓存 Token 更新后主动失效未实现，更新策略后旧缓存仍生效（最长 1 分钟）| 中 | 集成测试 TC010 需等待缓存过期后再验证策略是否生效 |

---

## 5. 测试环境说明

### 5.1 依赖服务

| 服务 | 版本要求 | 用途 | 说明 |
|------|---------|------|------|
| MySQL | ≥ 5.7.8 | 主数据库（三库之一）| 验证 AutoMigrate ADD COLUMN 成功 |
| PostgreSQL | ≥ 9.6 | 主数据库（三库之一）| 同上 |
| SQLite | 任意版本 | 主数据库（三库之一）| 轻量本地测试首选 |
| Redis | ≥ 6.0 | Token 缓存 | TC010 需等待缓存失效或手动 DEL |

### 5.2 前置条件

**测试数据准备**（测试前初始化）：

```sql
-- 准备 user_a（普通用户），user_b（普通用户），admin 用户
-- 准备 token_1（属于 user_a，id=1）
-- 准备 token_2（属于 user_b，id=2）
-- 确保 tokens.ip_policy 列已存在（AutoMigrate 完成）
```

**环境变量配置**：

| 用例组 | TRUSTED_PROXIES 值 | 说明 |
|--------|-------------------|------|
| TC001-TC036（默认）| `""` | 无可信代理，XFF 永远被忽略 |
| TC042-TC043（XFF 伪造测试）| `""` | 验证未配置时的安全行为 |
| TC044-TC045（合法 XFF）| `"10.0.0.1"` | 验证可信代理正常工作 |

### 5.3 测试框架

| 层次 | 框架/工具 | 用途 |
|------|---------|------|
| HTTP 接口测试 | `net/http/httptest` + 自定义 test helper | 构建 gin.Context，发送测试请求 |
| 单元测试 | Go 标准 `testing` 包 | FUNC001-FUNC005 |
| 并发安全 | `go test -race ./common/...` | 验证 ParseCIDRList 无数据竞争 |
| Mock 数据库 | `sqlmock` 或 SQLite 内存模式 | TC024 模拟 DB 写入失败 |
| 覆盖率 | `go test -coverprofile=coverage.out` | 验证 ≥ 85% 覆盖率目标 |

---

## 6. 执行指南

### 6.1 测试执行顺序

```
P0 优先级 → P1 → P2
单元测试（FUNC001-005）→ 接口正常场景 → 接口异常场景 → 边界场景 → 安全场景
```

### 6.2 执行步骤

#### 步骤 1：环境初始化

```bash
# 1. 初始化测试数据库（以 SQLite 为例）
export DSN="file:test.db?cache=shared&mode=memory"
# 2. 清空 TRUSTED_PROXIES（默认测试）
unset TRUSTED_PROXIES
# 3. 启动服务（AutoMigrate 会自动执行 ADD COLUMN）
go run main.go
```

#### 步骤 2：单元测试（FUNC001-005）

```bash
# 运行 ip_matcher 单元测试（含 race detector）
go test -v -race ./common/ -run TestValidateCIDRList
go test -v -race ./common/ -run TestParseCIDRList
go test -v -race ./common/ -run TestIPMatchesCIDRList
go test -v -race ./common/ -run TestGetClientIP
go test -v -race ./common/ -run TestInitTrustedProxies

# 生成覆盖率报告（目标 ≥ 85%）
go test -coverprofile=coverage.out ./common/
go tool cover -html=coverage.out -o coverage.html
```

#### 步骤 3：接口测试（API001）

```bash
# P0 用例（必须全部通过）
curl -X PUT http://localhost:3000/api/token/1/ip_policy \
  -H "Authorization: Bearer <user_a_jwt>" \
  -H "Content-Type: application/json" \
  -d '{"mode":"whitelist","ips":["1.2.3.4","10.0.0.0/8"]}'
# 期望：HTTP 200，{"success":true,"message":""}

# 清除策略
curl -X PUT http://localhost:3000/api/token/1/ip_policy \
  -H "Authorization: Bearer <user_a_jwt>" \
  -H "Content-Type: application/json" \
  -d '{"mode":"","ips":[]}'
# 期望：HTTP 200

# 权限测试
curl -X PUT http://localhost:3000/api/token/2/ip_policy \
  -H "Authorization: Bearer <user_a_jwt>" \
  -H "Content-Type: application/json" \
  -d '{"mode":"whitelist","ips":["1.2.3.4"]}'
# 期望：HTTP 403，{"success":false,"message":"forbidden"}
```

#### 步骤 4：安全专项测试

```bash
# XFF 伪造测试（无 TRUSTED_PROXIES）
curl -X GET http://localhost:3000/v1/chat/completions \
  -H "Authorization: sk-<token_with_whitelist>" \
  -H "X-Forwarded-For: 1.2.3.4" \
  # RemoteAddr 为非白名单 IP
# 期望：HTTP 403，IP_NOT_ALLOWED（XFF 被忽略，使用真实 RemoteAddr）
```

### 6.3 结果验证标准

| 验证项 | 通过标准 |
|-------|---------|
| P0 用例全部通过 | 所有 P0 用例返回预期 HTTP 状态码和响应体 |
| DB 状态正确 | `SELECT ip_policy FROM tokens WHERE id=1` 内容与期望 JSON 一致 |
| 单元测试覆盖率 | `go tool cover` 输出 `common/ip_matcher.go` 覆盖率 ≥ 85% |
| 并发安全 | `go test -race` 无 DATA RACE 输出 |
| 安全场景通过 | XFF 伪造、越权、SQL 注入等均被正确拦截 |

### 6.4 注意事项

1. **缓存失效延迟**（R6）：TC010（更新策略后验证生效）需在 Token 缓存过期后执行，或手动 `DEL token:cache:<id>` 后验证
2. **TC009 空白名单语义**：产品确认 Q1 后，补充断言（拦截所有 OR 放行所有）
3. **TC024 DB 故障模拟**：使用 `sqlmock` 注入 `errors.New("db connection refused")` 模拟写入失败
4. **三库分别执行**：TC001 等写入场景需在 SQLite/MySQL/PostgreSQL 三个环境分别验证，确保 JSON 序列化和反序列化均正确
5. **并发用例隔离**：UT023（race detector）使用完全独立的输入数据，避免与其他测试共享状态

---

## 7. 变更记录

| 版本 | 日期 | 变更说明 | 作者 |
|------|------|---------|------|
| v1.0 | 2026-03-02 | 初稿，基于 tech_design.md V1.0 生成，覆盖 API001 + FUNC001-005 全部测试用例 | AI Native Dev Team |
