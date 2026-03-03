# API Key IP 访问控制模块 - 编码任务

## 任务概览

| 属性 | 值 |
|------|-----|
| **功能名称** | api-key-ip-access-control |
| **设计文档** | `doc/develop/tech_design_detail.md` |
| **任务数量** | 5 |

---

## 任务清单

### [x] Task 1: CIDR 工具函数库

**目标**：新建 `common/ip_matcher.go`，实现 5 个 CIDR/IP 工具函数，纯 Go 标准库实现，无外部依赖，并发安全。

**涉及文件**：
- `common/ip_matcher.go` - 新建，实现全部工具函数

**核心功能**：
- 包级私有全局变量 `trustedProxyCIDRs []*net.IPNet`（服务启动时初始化，运行时只读）
- `ValidateCIDRList(ips []string) error`：逐条校验 IP/CIDR 格式，返回首个非法条目的 error，空列表返回 nil
- `ParseCIDRList(ips []string) ([]*net.IPNet, error)`：将字符串列表解析为 `[]*net.IPNet`，单 IP 自动补全 `/32`（IPv4）或 `/128`（IPv6），每次返回独立切片（无共享状态）
- `IPMatchesCIDRList(ipStr string, cidrs []*net.IPNet) bool`：判断 IP 是否命中列表，IP 无效或 cidrs 为空返回 false
- `GetClientIP(c *gin.Context) string`：从 gin.Context 提取真实客户端 IP；RemoteAddr 属于 `trustedProxyCIDRs` 时采信 XFF 最左侧 IP，否则使用 RemoteAddr；永远返回纯 IP 字符串（不含端口）
- `InitTrustedProxies(proxies string)`：解析逗号分隔的 IP/CIDR 字符串初始化 trustedProxyCIDRs，非法条目跳过并记录 `log.Printf("[WARN] ...")` 日志，空字符串时设为 nil
- 私有辅助函数 `isTrustedProxy(ip net.IP) bool`：判断 IP 是否属于 trustedProxyCIDRs

**设计约束**：
- 必须使用 `net.ParseCIDR`、`net.ParseIP`、`(*net.IPNet).Contains()` 标准库，不引入第三方依赖
- `GetClientIP` 使用 `net.SplitHostPort` 分离 RemoteAddr 的 IP 和端口
- XFF 解析：`strings.Split(xff, ",")` 取 `[0]` 并 `strings.TrimSpace`

---

### [x] Task 2: CIDR 工具函数单元测试

**目标**：新建 `common/ip_matcher_test.go`，对 Task 1 的 5 个函数实现全面单元测试，覆盖率 ≥ 85%，通过 `-race` 并发安全验证。

**涉及文件**：
- `common/ip_matcher_test.go` - 新建，与被测文件同包（`package common`）

**核心功能**：
- 框架：Go 标准库 `testing` + `github.com/stretchr/testify/assert`，命名规范：`Test{FuncName}_{场景}`
- `TestValidateCIDRList_*`：覆盖空列表、合法 IPv4/IPv6/CIDR、非法格式（IP 段超范围、前缀超范围、非 IP 字符串、空字符串）、混合合法与非法（返回首个非法）
- `TestParseCIDRList_*`：覆盖空列表、单 IP 自动补全 /32、合法 CIDR 解析、含非法条目报错、每次返回独立切片（修改 slice1 不影响 slice2）
- `TestParseCIDRList_ConcurrentSafety`：100 goroutine 并发调用，通过 `go test -race` 无数据竞争
- `TestIPMatchesCIDRList_*`：覆盖命中/未命中、空 cidrs、网段首末地址边界、IPv6、无效 IP 字符串
- `TestGetClientIP_*`：使用 `httptest.NewRequest` 构建 gin.Context，覆盖无可信代理忽略 XFF、可信代理采信 XFF 最左侧 IP、可信代理 XFF 为空回退 RemoteAddr、非可信代理忽略 XFF
- `TestInitTrustedProxies_*`：覆盖空字符串设 nil、单 IP、CIDR、多条逗号分隔、含非法条目跳过

**依赖**：Task 1

---

### [x] Task 3: Token.IpPolicy 数据模型扩展

**目标**：在 `model/token.go` 中扩展 Token 数据模型，新增 IpPolicy 结构体及其序列化实现，新增专用 `UpdateIpPolicy()` 方法。

**涉及文件**：
- `model/token.go` - 修改，新增数据结构、字段、方法

**核心功能**：
- 新增 `IpPolicy` 结构体（位于 Token struct 定义之前）：
  ```go
  type IpPolicy struct {
      Mode string   `json:"mode"`
      Ips  []string `json:"ips"`
  }
  ```
- 实现 `func (p *IpPolicy) Value() (driver.Value, error)`（指针接收者）：p 为 nil 时返回 `(nil, nil)` → SQL NULL；否则调用 `common.Marshal(p)` 序列化为 JSON 字符串（必须使用 common.Marshal，Rule 1）
- 实现 `func (p *IpPolicy) Scan(value interface{}) error`：兼容 `[]byte` 和 `string` 两种驱动返回类型；调用 `common.Unmarshal(bytes, p)`（必须使用 common.Unmarshal，Rule 1）
- Token struct 新增字段（紧接 AllowIps 字段之后）：
  ```go
  IpPolicy *IpPolicy `json:"ip_policy,omitempty" gorm:"type:text;column:ip_policy"`
  ```
- 新增 `func (token *Token) UpdateIpPolicy() (err error)`：
  - 使用 `DB.Model(&Token{}).Where("id = ?", token.Id).Update("ip_policy", token.IpPolicy).Error` 精确更新单字段（避免 GORM 跳过零值问题）
  - `defer` 中使用 `shouldUpdateRedis(true, err)` + `gopool.Go(func() { cacheDeleteToken(token.Key) })` 异步删除 Redis 缓存（与 Token.Update() 模式一致）
  - 需要导入 `database/sql/driver` 包（用于 `driver.Value` 类型）

**设计约束**：
- 旧记录 ip_policy 为 NULL → GORM Scan 时 IpPolicy 保持 nil → TokenAuth 跳过策略（向后兼容）
- GORM AutoMigrate 已在 `model/main.go:256` 包含 `&Token{}`，无需修改

**依赖**：Task 1（`common.Marshal`/`common.Unmarshal`）

---

### [x] Task 4: IP 策略配置 HTTP 接口

**目标**：在 `controller/token.go` 新增 `UpdateTokenIpPolicy` Handler 实现 `PUT /api/token/:id/ip_policy` 接口，并在路由文件注册路由。

**涉及文件**：
- `controller/token.go` - 修改，新增 DTO 和 Handler
- `router/api-router.go` - 修改，注册新路由

**核心功能**：
- 新增 `IpPolicyRequest` DTO（在 controller/token.go 中，与 TokenBatch 等 DTO 放在一起）：
  ```go
  type IpPolicyRequest struct {
      Mode string   `json:"mode"`
      Ips  []string `json:"ips"`
  }
  ```
- 实现 `UpdateTokenIpPolicy(c *gin.Context)` Handler，处理流程：
  1. `strconv.Atoi(c.Param("id"))` 解析路径参数，失败返回 `common.ApiError(c, err)`
  2. `c.ShouldBindJSON(&req)` 绑定请求体，失败返回 `common.ApiError(c, err)`
  3. 校验 mode：`mode != "" && mode != "whitelist" && mode != "blacklist"` → `c.JSON(400, gin.H{"success":false,"message":"invalid mode: "+req.Mode})`
  4. 校验条目数：`len(req.Ips) > 100` → `c.JSON(400, gin.H{"success":false,"message":"too many IP entries, max 100"})`
  5. 校验 CIDR 格式（仅 mode 非空时）：`common.ValidateCIDRList(req.Ips)` 失败 → `c.JSON(400, gin.H{"success":false,"message":"invalid IP/CIDR: "+err.Error()})`
  6. 查询 Token：`model.GetTokenById(id)` 失败 → `c.JSON(404, gin.H{"success":false,"message":"not found"})`
  7. 权限校验：`token.UserId != c.GetInt("id") && !model.IsAdmin(c.GetInt("id"))` → `c.JSON(403, gin.H{"success":false,"message":"forbidden"})`
  8. 构建 IpPolicy：mode=="" 时 `token.IpPolicy = nil`；否则 `token.IpPolicy = &model.IpPolicy{Mode: req.Mode, Ips: req.Ips}`
  9. 持久化：`token.UpdateIpPolicy()` 失败 → `c.JSON(500, gin.H{"success":false,"message":"internal server error"})`
  10. 成功：`c.JSON(http.StatusOK, gin.H{"success":true,"message":""})`
- `router/api-router.go`：在 tokenRoute 组最后追加 `tokenRoute.PUT("/:id/ip_policy", controller.UpdateTokenIpPolicy)`

**设计约束**：
- 权限校验必须使用 `model.GetTokenById(id)`（不是 `GetTokenByIds`），后者会因 userId 联合查询导致管理员操作他人 Token 返回 404
- 响应格式与现有 Token 接口一致（`{"success":true/false,"message":"..."}`）

**依赖**：Task 1（`common.ValidateCIDRList`）、Task 3（`model.IpPolicy`、`token.UpdateIpPolicy()`）

---

### [x] Task 5: TokenAuth 中间件 IP 策略执行

**目标**：扩展 `middleware/auth.go` 中的 `TokenAuth()` 函数，插入 F004（可信代理感知 IP 提取）和 F005（IP 策略校验）代码块；在 `common/init.go` 注册启动时初始化；在 `types/error.go` 新增错误码常量。

**涉及文件**：
- `middleware/auth.go` - 修改，在 TokenAuth() AllowIps 校验块之后插入两个代码块
- `common/init.go` - 修改，InitEnv() 末尾追加初始化调用
- `types/error.go` - 修改，新增 ErrorCodeIpNotAllowed 常量

**核心功能**：

**[types/error.go]**：在 `ErrorCode` 常量块（`types/error.go:66` 附近）的 client request error 区域追加：
```go
ErrorCodeIpNotAllowed ErrorCode = "ip_not_allowed"
```

**[common/init.go]**：在 `InitEnv()` 函数末尾（`initConstantEnv()` 调用之后）追加：
```go
common.InitTrustedProxies(os.Getenv("TRUSTED_PROXIES"))
```
（注意：`initConstantEnv()` 调用在 `init.go` 内部，`InitTrustedProxies` 在 `common` 包内，同包调用直接写 `InitTrustedProxies(os.Getenv("TRUSTED_PROXIES"))`）

**[middleware/auth.go — F004]**：在 AllowIps 校验块结束后（`logger.LogDebug(c, "Client IP %s passed the token IP restrictions check", clientIp)` 之后，第 329-330 行之间）插入：
```go
// [F004] 可信代理感知 IP 提取
clientIP := common.GetClientIP(c)
c.Set("client_ip", clientIP)
```

**[middleware/auth.go — F005]**：在 F004 代码块之后，`userCache, err := model.GetUserCache(token.UserId)` 之前插入：
```go
// [F005] IpPolicy 策略校验
if token.IpPolicy != nil && token.IpPolicy.Mode != "" {
    cidrs, parseErr := common.ParseCIDRList(token.IpPolicy.Ips)
    if parseErr != nil {
        logger.LogWarn(c.Request.Context(), fmt.Sprintf(
            "ip_policy_parse_failed: token_id=%d err=%s", token.Id, parseErr.Error()))
    } else {
        clientIPForPolicy := c.GetString("client_ip")
        hit := common.IPMatchesCIDRList(clientIPForPolicy, cidrs)
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

**设计约束**：
- F005 `blocked` 判断后必须 `return`，不能只 `Abort()`（`abortWithOpenAiMessage` 内部已调用 `c.Abort()`，但外层 `return` 必须存在以防止后续代码执行）
- fail-open 语义：`ParseCIDRList` 失败时记录 WARN 日志后**继续**请求，不调用 `Abort()`
- `logger.LogWarn` 如包中不存在，使用 `logger.SysLog` 或参考 `middleware/auth.go` 中现有日志调用方式

**依赖**：Task 1（`common.GetClientIP`、`common.ParseCIDRList`、`common.IPMatchesCIDRList`、`common.InitTrustedProxies`）、Task 3（`token.IpPolicy` 字段）

---

## 任务依赖关系

```
Task 1: CIDR 工具函数库（独立，最优先实现）
    ├──→ Task 2: CIDR 工具函数单元测试（依赖 Task 1，可与 Task 3 并行）
    └──→ Task 3: Token.IpPolicy 数据模型扩展（依赖 Task 1）
              ├──→ Task 4: IP 策略配置 HTTP 接口（依赖 Task 1 + Task 3，可与 Task 5 并行）
              └──→ Task 5: TokenAuth 中间件 IP 策略执行（依赖 Task 1 + Task 3，可与 Task 4 并行）
```

**执行建议**：
1. **Task 1** 必须最先完成（其他 4 个 Task 均直接或间接依赖它）
2. Task 1 完成后，**Task 2 和 Task 3 可并行开发**
3. Task 3 完成后，**Task 4 和 Task 5 可并行开发**
4. Task 2 不阻塞 Task 4/Task 5，可在任意时间完成

---

## 整体完成标准

- [x] 所有 Task 完成（Task 1-5）
- [x] `go build ./...` 无编译错误
- [x] `go test -race ./common/ -run TestValidateCIDRList` 通过
- [x] `go test -race ./common/ -run TestParseCIDRList` 通过（含并发安全）
- [x] `go test -coverprofile=coverage.out ./common/` 覆盖率 ≥ 85%
- [x] 代码符合项目编码规范（common.Marshal/Unmarshal 规则，不直接使用 encoding/json）
- [x] 无破坏性变更：AllowIps 旧字段和 GetIpLimits() 方法保持不变
