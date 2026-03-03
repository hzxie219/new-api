# 数据结构需求检查（数据结构需求_check.md）

> 基于需求文档 doc/requirement/requirement.md 和项目代码检索的数据结构需求分析

## 1. 数据结构需求分类

| 数据结构 | 类型 | 来源 | 代码位置 |
|---------|------|------|---------|
| `IpPolicy` struct | 全新结构体 | requirement.md §2.2.2 F002 | 新增至 model/token.go |
| `Token.IpPolicy` 字段 | 增量字段 | requirement.md §2.2.2 F002 | model/token.go:14（Token struct） |
| `trustedProxyCIDRs` 全局变量 | 全新全局变量 | requirement.md §2.2.2 F003 | 新增至 common/ip_matcher.go |
| `TRUSTED_PROXIES` 环境变量 | 新增配置 | requirement.md §3.2 | common/init.go:InitEnv() |
| `tokens.ip_policy` 数据库列 | 增量DDL | requirement.md §3.1 | GORM AutoMigrate（model/main.go:256） |

## 2. 现有相关数据结构

| 现有结构 | 代码位置 | 关联关系 |
|---------|---------|---------|
| `Token` struct | model/token.go:14 | 新增 IpPolicy 字段（增量） |
| `AllowIps *string` | model/token.go:27 | 旧 IP 白名单字段，保持不变 |
| `net.IPNet` | Go 标准库 | `trustedProxyCIDRs` 和 ParseCIDRList 返回类型 |
| `driver.Valuer` 接口 | database/sql/driver | IpPolicy.Value() 实现 |
| `sql.Scanner` 接口 | database/sql | IpPolicy.Scan() 实现 |

## 3. 配置文件变更

| 配置项 | 类型 | 来源 | 读取位置 |
|-------|------|------|---------|
| `TRUSTED_PROXIES` 环境变量 | 新增 | requirement.md §3.2 | common/init.go:InitEnv()（新增1行） |

## 4. 数据库变更

| 表名 | 变更类型 | 变更内容 | 迁移方式 |
|-----|---------|---------|---------|
| `tokens` | ADD COLUMN | `ip_policy TEXT NULL` | GORM AutoMigrate（model/main.go:256，无需手工迁移）|
