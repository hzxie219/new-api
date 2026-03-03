# 【Tech-系统级】策略持久化与配置接口

## 技术故事

作为系统，需要将 API Key 的 IP 策略（mode + ips）持久化到数据库，并对外提供 REST 接口供前端和管理 API 调用，完成策略的增删改。

## 所属 Story

Story 1.1.1 — 配置 API Key 的 IP 访问策略

## 子系统职责

本子系统负责：
1. `tokens` 表新增 `ip_policy` 字段（存储 JSON 格式的策略对象）
2. 提供 `PUT /api/keys/{id}/ip_policy` REST 接口，包含入参校验和权限控制
3. 数据库迁移兼容 SQLite / MySQL / PostgreSQL 三种数据库

## 验收条件（Tech 级）

### TC-1: 数据库字段正确存储策略

- **Given**: `tokens` 表存在，且已完成数据库迁移
- **When**: 通过 `PUT /api/keys/{id}/ip_policy` 提交有效策略 `{"mode":"whitelist","ips":["1.2.3.4/32"]}`
- **Then**: 该 Token 记录的 `ip_policy` 字段存储正确的 JSON，`mode="whitelist"`，`ips` 包含 `"1.2.3.4/32"`

### TC-2: 接口参数校验 — mode 非法值

- **Given**: 接口接收请求
- **When**: body 中 `mode="invalid_mode"`
- **Then**: 返回 HTTP 400，错误信息中说明 mode 只允许 `whitelist`/`blacklist`/空值

### TC-3: 接口参数校验 — CIDR 格式错误

- **Given**: 接口接收请求
- **When**: body 中 `ips` 包含 `"999.0.0.0/8"`
- **Then**: 返回 HTTP 400，错误信息指明非法 IP/CIDR 条目

### TC-4: 权限控制 — 非 Owner 操作

- **Given**: 请求方 userId=100，目标 Token 的 owner userId=200
- **When**: 普通用户调用 `PUT /api/keys/{id}/ip_policy`
- **Then**: 返回 HTTP 403，策略不被修改（管理员可操作任意 Token）

### TC-5: 清空策略

- **Given**: Token 已有 IP 策略
- **When**: 提交 `{"mode":"","ips":[]}`
- **Then**: `ip_policy` 字段更新为空/null，后续请求不再进行 IP 校验

## 技术思路

### 数据模型

在 `new-api/model/token.go` 的 `Token` 结构体中新增字段：

```go
type IpPolicy struct {
    Mode string   `json:"mode"` // "whitelist" | "blacklist" | ""
    Ips  []string `json:"ips"`  // CIDR 或精确 IP 列表
}

type Token struct {
    // ...已有字段...
    IpPolicy *IpPolicy `json:"ip_policy" gorm:"type:text;serializer:json"`
}
```

数据库存储为 TEXT 类型（JSON 序列化），兼容 SQLite/MySQL/PostgreSQL。

### 接口设计

```
PUT /api/keys/{id}/ip_policy
Authorization: Bearer <admin_or_owner_token>
Content-Type: application/json

Request Body:
{
  "mode": "whitelist" | "blacklist" | "",
  "ips": ["1.2.3.4/32", "10.0.0.0/8"]
}

Response 200:
{
  "success": true,
  "message": "IP 策略更新成功"
}

Response 400:
{
  "success": false,
  "message": "invalid CIDR: 999.0.0.0/8"
}
```

### 依赖关系

- 依赖 `new-api/model/token.go`: Token 结构体
- 依赖 `new-api/common/ip_matcher.go`（新增）: CIDR 格式校验
- 路由注册在 `new-api/router/api-router.go`

## 风险与依赖

- 数据库迁移：`AutoMigrate` 会自动 ADD COLUMN，不影响存量 Token（ip_policy 默认 null，等同于无策略）
- 与现有 Token IP 白名单字段的共存方案需确认（是否废弃旧字段？）
