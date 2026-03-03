# 【Tech-服务级】controller-token

## 服务职责

实现 `PUT /api/keys/{id}/ip_policy` REST 接口，提供 API Key IP 策略的配置能力，包含参数校验、权限控制和数据持久化。

## 所属 Tech-系统级

Tech-系统级-策略持久化与配置接口

## 详细验收条件

### AC-1: 接口路由正确注册

- **Given**: 服务启动
- **When**: 查看路由表
- **Then**: `PUT /api/keys/:id/ip_policy` 已注册，需管理员或 Token Owner 鉴权

### AC-2: 正常更新策略

- **Given**: 请求方为 Token Owner 或管理员，id 对应的 Token 存在
- **When**: `PUT /api/keys/1/ip_policy` body `{"mode":"blacklist","ips":["5.5.5.0/24"]}`
- **Then**: HTTP 200，数据库对应 Token 的 `ip_policy` 已更新

### AC-3: mode 非法值校验

- **Given**: 任意有效 Token id
- **When**: body `{"mode":"allow","ips":[]}`
- **Then**: HTTP 400，message 包含 "mode must be whitelist or blacklist"

### AC-4: CIDR 格式校验

- **Given**: 任意有效 Token id
- **When**: body `{"mode":"whitelist","ips":["bad-ip"]}`
- **Then**: HTTP 400，message 包含非法 IP 条目说明

### AC-5: Token 不存在

- **Given**: 请求 id=99999（不存在）
- **When**: 发起请求
- **Then**: HTTP 404

### AC-6: 非 Owner 被拒绝

- **Given**: 请求方 userId ≠ Token.UserId，且请求方非管理员
- **When**: 发起请求
- **Then**: HTTP 403

## 技术实现

### 代码位置

- Handler: `new-api/controller/token.go`（新增 `UpdateTokenIpPolicy` 函数）
- 路由注册: `new-api/router/api-router.go`（在 Token 路由组中追加）
- CIDR 校验: 调用 `new-api/common/ip_matcher.go` 的 `ValidateCIDRList` 函数

### 核心代码参考

```go
// new-api/controller/token.go

type UpdateIpPolicyRequest struct {
    Mode string   `json:"mode"`
    Ips  []string `json:"ips"`
}

func UpdateTokenIpPolicy(c *gin.Context) {
    id, err := strconv.Atoi(c.Param("id"))
    if err != nil {
        common.ApiErrorMsg(c, "invalid id")
        return
    }

    var req UpdateIpPolicyRequest
    if err := common.UnmarshalBodyReusable(c, &req); err != nil {
        common.ApiErrorMsg(c, "invalid request body")
        return
    }

    // mode 校验
    if req.Mode != "" && req.Mode != "whitelist" && req.Mode != "blacklist" {
        common.ApiErrorMsg(c, "mode must be whitelist, blacklist or empty")
        return
    }

    // CIDR 格式校验
    if err := common.ValidateCIDRList(req.Ips); err != nil {
        common.ApiErrorMsg(c, err.Error())
        return
    }

    // 权限校验 & 持久化...
    common.ApiSuccess(c, nil)
}
```

### 路由注册

```go
// new-api/router/api-router.go
tokenRoute.PUT("/:id/ip_policy", controller.UpdateTokenIpPolicy)
```

### 注意事项

- 使用 `common.UnmarshalBodyReusable` 解析请求体（项目规范）
- 使用 `common.ApiSuccess / common.ApiErrorMsg` 返回响应（项目规范）
- 日志使用 `logger.LogInfo(ctx, ...)` 记录操作

## 监控与排障

- 操作日志：策略更新成功时写入 type=3（管理操作日志）
- 错误日志关键字: `update_ip_policy failed`
