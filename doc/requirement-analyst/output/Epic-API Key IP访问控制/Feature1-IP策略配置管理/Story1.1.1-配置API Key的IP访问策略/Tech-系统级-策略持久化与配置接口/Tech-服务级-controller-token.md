# 【Tech-服务级】controller-token（UpdateTokenIpPolicy）

## 服务职责

实现 `PUT /api/keys/:id/ip_policy` REST 接口的 Handler 函数 `UpdateTokenIpPolicy`，提供 API Key IP 策略的配置能力。包含参数解析、mode 枚举校验、CIDR 格式校验、权限控制（Owner/Admin）和数据持久化完整链路。

**所属 Tech-系统级**：Tech-系统级-策略持久化与配置接口

---

## 详细验收条件

### AC-1: 路由正确注册

- **Given**: 服务启动
- **When**: 查看路由表（如 gin debug 模式输出）
- **Then**: `PUT /api/keys/:id/ip_policy` 已注册，路由组带有管理员或 Owner 鉴权中间件

### AC-2: 正常更新白名单策略

- **Given**: 请求方为 Token Owner（userId 匹配），Token id=1 存在
- **When**: `PUT /api/keys/1/ip_policy`，body `{"mode":"whitelist","ips":["1.2.3.4/32","10.0.0.0/8"]}`
- **Then**: HTTP 200，`{"success":true,"message":"..."}`；数据库 `tokens.ip_policy` 更新为对应 JSON

### AC-3: 正常更新黑名单策略（管理员）

- **Given**: 请求方为管理员（is_admin=true），Token 存在
- **When**: `PUT /api/keys/2/ip_policy`，body `{"mode":"blacklist","ips":["5.5.5.0/24"]}`
- **Then**: HTTP 200；策略写入成功

### AC-4: 清空策略

- **Given**: Token 已有 IP 策略
- **When**: body `{"mode":"","ips":[]}`
- **Then**: HTTP 200；数据库 `ip_policy` 字段更新为 NULL 或空策略

### AC-5: mode 非法值 — 返回 400

- **Given**: 任意有效 Token id，请求方有权限
- **When**: body `{"mode":"allow","ips":[]}`
- **Then**: HTTP 400；`message` 含 "mode must be whitelist, blacklist or empty"；数据库未写入

### AC-6: CIDR 格式错误 — 返回 400

- **Given**: 任意有效 Token id，请求方有权限
- **When**: body `{"mode":"whitelist","ips":["bad-ip","999.0.0.0/8"]}`
- **Then**: HTTP 400；`message` 包含第一个非法条目说明（如 `"invalid IP/CIDR: bad-ip"`）；数据库未写入

### AC-7: Token 不存在 — 返回 404

- **Given**: 请求 id=99999，数据库中不存在该记录
- **When**: `PUT /api/keys/99999/ip_policy`（入参合法）
- **Then**: HTTP 404

### AC-8: 非 Owner 普通用户 — 返回 403

- **Given**: 请求方为普通用户（is_admin=false）；目标 Token.UserId ≠ 请求方 userId
- **When**: `PUT /api/keys/{id}/ip_policy`
- **Then**: HTTP 403；数据库未写入

### AC-9: 入参校验顺序（防止越权探测）

- **Given**: 请求方为普通用户，目标 Token 不存在且请求方无权限
- **When**: 发起请求
- **Then**: 先做入参格式校验（400）再做权限校验（403），**Token 存在性校验在权限校验之后**（防止通过 404 探测他人 Token 是否存在）

> 注：具体校验顺序见技术实现节

### AC-10: 操作日志写入

- **Given**: 策略更新成功
- **When**: Handler 完成持久化
- **Then**: 写入操作日志（type=3），记录操作者 userId、target token_id、操作内容（mode 和 ips 条目数）

---

## 技术实现

### 代码位置

| 文件 | 变更内容 |
|------|---------|
| `new-api/controller/token.go` | 新增 `UpdateIpPolicyRequest` 结构体和 `UpdateTokenIpPolicy` 函数 |
| `new-api/router/api-router.go` | 在 Token 路由组中追加 `PUT /:id/ip_policy` |

### 核心代码参考

```go
// new-api/controller/token.go

type UpdateIpPolicyRequest struct {
    Mode string   `json:"mode"`
    Ips  []string `json:"ips"`
}

func UpdateTokenIpPolicy(c *gin.Context) {
    // 步骤1: 解析 id
    id, err := strconv.Atoi(c.Param("id"))
    if err != nil || id <= 0 {
        common.ApiErrorMsg(c, http.StatusBadRequest, "invalid token id")
        return
    }

    // 步骤2: 解析请求体
    var req UpdateIpPolicyRequest
    if err := common.UnmarshalBodyReusable(c, &req); err != nil {
        common.ApiErrorMsg(c, http.StatusBadRequest, "invalid request body")
        return
    }

    // 步骤3: 校验 mode 枚举值
    if req.Mode != "" && req.Mode != "whitelist" && req.Mode != "blacklist" {
        common.ApiErrorMsg(c, http.StatusBadRequest,
            "mode must be whitelist, blacklist or empty")
        return
    }

    // 步骤4: 校验 CIDR 格式
    if err := common.ValidateCIDRList(req.Ips); err != nil {
        common.ApiErrorMsg(c, http.StatusBadRequest, err.Error())
        return
    }

    // 步骤5: 查询 Token（在权限校验之前，但仅在格式校验之后）
    userId := c.GetInt("id")
    isAdmin := c.GetBool("isAdmin")

    token, err := model.GetTokenById(id)
    if err != nil || token == nil {
        common.ApiErrorMsg(c, http.StatusNotFound, "token not found")
        return
    }

    // 步骤6: 权限校验（Owner 或 Admin）
    if !isAdmin && token.UserId != userId {
        common.ApiErrorMsg(c, http.StatusForbidden, "forbidden")
        return
    }

    // 步骤7: 构造 IpPolicy 并持久化
    var policy *model.IpPolicy
    if req.Mode != "" {
        policy = &model.IpPolicy{Mode: req.Mode, Ips: req.Ips}
    }
    if err := model.UpdateTokenIpPolicy(id, policy); err != nil {
        common.ApiErrorMsg(c, http.StatusInternalServerError, "update failed")
        return
    }

    // 步骤8: 写操作日志
    // logger.LogInfo(c.Request.Context(), fmt.Sprintf(
    //     "update_ip_policy: user_id=%d token_id=%d mode=%s ips_count=%d",
    //     userId, id, req.Mode, len(req.Ips)))

    common.ApiSuccess(c, nil)
}
```

### 路由注册

```go
// new-api/router/api-router.go
// 在现有 tokenRoute 路由组（带用户鉴权）中追加：
tokenRoute.PUT("/:id/ip_policy", controller.UpdateTokenIpPolicy)
```

### 注意事项

1. 使用 `common.UnmarshalBodyReusable` 解析请求体（项目规范）
2. 使用 `common.ApiSuccess / common.ApiErrorMsg` 返回响应（项目规范）
3. 权限校验放在 Token 存在性校验之后（AC-9 说明）：防止通过 404/403 差异探测他人 Token 存在性
4. `model.UpdateTokenIpPolicy` 建议为单独的 model 方法，避免在 controller 中直接操作数据库

---

## 监控与排障

| 场景 | 日志关键字 | 处理方式 |
|------|-----------|---------|
| 策略更新成功 | `update_ip_policy` | INFO 级别，记录 user_id、token_id、mode、ips_count |
| 数据库写入失败 | `update_ip_policy_failed` | ERROR 级别，返回 HTTP 500 |
| CIDR 校验失败（频繁） | `invalid_ip_policy_format` | 统计指标，排查调用方格式问题 |
