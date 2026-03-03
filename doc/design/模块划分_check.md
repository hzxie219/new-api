# 模块划分检查（模块划分_check.md）

> 基于需求文档 doc/requirement/requirement.md 和第4章架构划分的子模块分析

## 模块划分结果

| 模块编号 | 模块名称 | 功能覆盖 | 类型 | 文件位置 |
|---------|---------|---------|------|---------|
| M1 | CIDR工具函数库 | F003（ValidateCIDRList/ParseCIDRList/IPMatchesCIDRList/GetClientIP/InitTrustedProxies） | 全新 | common/ip_matcher.go |
| M2 | Token数据模型扩展 | F002（IpPolicy struct + Valuer/Scanner） | 增量 | model/token.go |
| M3 | IP策略配置接口 | F001（UpdateTokenIpPolicy Handler + 路由注册） | 全新 | controller/token.go + router/api-router.go |
| M4 | 鉴权中间件IP策略扩展 | F004（GetClientIP调用）+ F005（IpPolicy策略执行块） | 增量 | middleware/auth.go + common/init.go |

## 子模块设计顺序

```
M1（CIDR工具函数库）→ M2（数据模型）→ M3（配置接口）→ M4（鉴权中间件）
```

依赖关系：M1 被 M2/M3/M4 依赖；M2 被 M3/M4 依赖；M4 依赖 M1+M2。
