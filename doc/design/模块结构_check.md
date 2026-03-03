# 模块代码结构检查（模块结构_check.md）

> 基于代码仓库检索的现有模块结构

## M1 - common/ip_matcher.go（全新文件）

不存在，需完整新建。参考 `common/ip.go:33` 中的 `IsIpInCIDRList` 实现风格。

## M2 - model/token.go（增量）

| 现有内容 | 行号 | 增量内容 |
|---------|------|---------|
| Token struct 定义 | 14 | 新增 `IpPolicy *IpPolicy` 字段 |
| AllowIps *string | 27 | 保持不变 |
| GetIpLimits() 方法 | 38 | 保持不变 |
| 无 IpPolicy struct | - | 新增 IpPolicy struct + Valuer/Scanner |

## M3 - controller/token.go（增量）

| 现有内容 | 行号 | 增量内容 |
|---------|------|---------|
| UpdateToken 函数 | 225 | 参考此函数新增 UpdateTokenIpPolicy |
| 无 UpdateTokenIpPolicy | - | 全新函数 |

路由注册：`router/api-router.go:241-250` tokenRoute 路由组，新增 `tokenRoute.PUT("/:id/ip_policy", controller.UpdateTokenIpPolicy)`

## M4 - middleware/auth.go（增量）

| 现有内容 | 行号 | 增量内容 |
|---------|------|---------|
| TokenAuth() 函数 | 248 | 在 AllowIps 块（316-330行）之后插入 GetClientIP + IpPolicy 策略块 |
| AllowIps 白名单逻辑 | 316-330 | 保持不变 |

初始化：`common/init.go:InitEnv()` 函数末尾追加 `common.InitTrustedProxies(os.Getenv("TRUSTED_PROXIES"))`
