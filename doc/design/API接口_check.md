# API接口检索结果（API接口_check.md）

> 基于需求文档 doc/requirement/requirement.md 检索的 API 接口信息

## API接口分类

### 全新接口（需新增）

| API名称 | HTTP方法 | 路径 | 代码位置（需新增） | 说明 |
|---------|---------|-----|----------------|------|
| 更新Token IP策略 | PUT | `/api/token/:id/ip_policy` | `controller/token.go`（新增 `UpdateTokenIpPolicy` 函数）| 为指定 Token 设置或清除 IP 访问策略 |

### 增量接口（无）

本需求无增量接口（不修改现有接口的参数或响应）。

## 现有Token路由组信息

代码位置：`router/api-router.go:241-250`

```
tokenRoute := apiRouter.Group("/token")
tokenRoute.Use(middleware.UserAuth())

现有路由：
- GET    /api/token/         → controller.GetAllTokens
- GET    /api/token/search   → controller.SearchTokens
- GET    /api/token/:id      → controller.GetToken
- POST   /api/token/         → controller.AddToken
- PUT    /api/token/         → controller.UpdateToken
- DELETE /api/token/:id      → controller.DeleteToken
- POST   /api/token/batch    → controller.DeleteTokenBatch
```

新路由注册方式：
```go
tokenRoute.PUT("/:id/ip_policy", controller.UpdateTokenIpPolicy)
```

## 认证要求

- 使用 `middleware.UserAuth()` 中间件（已在路由组级别挂载）
- Handler 内部需额外校验：当前用户必须为 Token 拥有者（`token.UserId == currentUserId`）或管理员

## 检索结论

- **全新接口**：1 个（`PUT /api/token/:id/ip_policy`）
- **增量接口**：0 个
- **消息接口**：0 个（本需求不使用消息队列）
