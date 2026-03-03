# 【Tech-系统级】策略持久化与配置接口

## 【技术故事】

作为系统，需要将 API Key 的 IP 策略（mode + ips）持久化到数据库，并对外提供 REST 接口供前端和管理 API 调用，完成策略的增删改。同时需要兼容 SQLite/MySQL/PostgreSQL 三种数据库，确保存量 Token 不受影响。

**所属 Story**：Story 1.1.1 — 配置 API Key 的 IP 访问策略

**子系统职责**：
1. `tokens` 表新增 `ip_policy` 字段（TEXT 类型，JSON 序列化，三数据库兼容）
2. 提供 `PUT /api/keys/{id}/ip_policy` REST 接口，包含 mode 枚举校验、CIDR 格式校验和权限控制
3. 数据库迁移通过 GORM `AutoMigrate` 自动完成，存量数据默认 NULL（无策略）

---

## 【验收条件】

### 一、功能性验收条件

**TC-1: 数据库字段正确存储白名单策略**
- **Given**: `tokens` 表已完成 AutoMigrate，存在 `ip_policy` 列
- **When**: 通过 `PUT /api/keys/1/ip_policy` 提交 `{"mode":"whitelist","ips":["1.2.3.4/32","10.0.0.0/8"]}`
- **Then**: 数据库对应行 `ip_policy` 列存储为合法 JSON，反序列化后 `mode="whitelist"`，`ips=["1.2.3.4/32","10.0.0.0/8"]`

**TC-2: 数据库字段正确存储黑名单策略**
- **Given**: 同上
- **When**: 提交 `{"mode":"blacklist","ips":["5.5.5.0/24"]}`
- **Then**: `ip_policy` 反序列化后 `mode="blacklist"`，`ips=["5.5.5.0/24"]`

**TC-3: 清空策略写入 NULL**
- **Given**: Token 已有 IP 策略
- **When**: 提交 `{"mode":"","ips":[]}`
- **Then**: `ip_policy` 更新为 NULL 或空策略；后续读取时 `IpPolicy == nil`

**TC-4: mode 非法值校验**
- **Given**: 接口接收请求
- **When**: body `{"mode":"invalid","ips":[]}`
- **Then**: HTTP 400，message 包含 "mode must be whitelist, blacklist or empty"；数据库未写入

**TC-5: CIDR 格式校验 — 非法 IP**
- **Given**: 接口接收请求
- **When**: body `{"mode":"whitelist","ips":["999.0.0.0/8"]}`
- **Then**: HTTP 400，message 包含 "invalid IP/CIDR: 999.0.0.0/8"；数据库未写入

**TC-6: CIDR 格式校验 — 非 IP 字符串**
- **Given**: 接口接收请求
- **When**: body `{"mode":"whitelist","ips":["not-an-ip"]}`
- **Then**: HTTP 400，message 包含非法条目说明

**TC-7: 权限控制 — 非 Owner 普通用户**
- **Given**: 请求方 userId=100，目标 Token owner userId=200，请求方非管理员
- **When**: `PUT /api/keys/{id}/ip_policy`
- **Then**: HTTP 403；数据库未写入

**TC-8: 权限控制 — 管理员可操作任意 Token**
- **Given**: 请求方为管理员（is_admin=true），目标 Token 属于任意用户
- **When**: `PUT /api/keys/{id}/ip_policy`（合法入参）
- **Then**: HTTP 200；策略写入成功

**TC-9: Token 不存在**
- **Given**: 请求 id 在数据库中不存在
- **When**: `PUT /api/keys/99999/ip_policy`
- **Then**: HTTP 404

**TC-10: 三数据库兼容迁移**
- **Given**: SQLite / MySQL / PostgreSQL 任意一种数据库，服务首次启动
- **When**: 服务启动触发 AutoMigrate
- **Then**: `tokens` 表成功新增 `ip_policy` 列（TEXT 类型），存量行该列默认 NULL，服务正常启动

### 二、非功能性验收条件

**TC-N1: 接口响应时间**
- **Given**: 数据库正常，单次策略写入请求
- **When**: 并发 10 请求下
- **Then**: P99 响应时间 < 200ms

**TC-N2: 并发写入安全**
- **Given**: 同一 Token 被并发请求更新 IP 策略
- **When**: 10 个并发请求同时更新
- **Then**: 最终数据库中保存的是其中某一次请求的有效策略，不出现数据损坏（GORM 事务保证）

**TC-N3: 存量 Token 不受影响**
- **Given**: 已有存量 Token 记录（无 `ip_policy` 列）
- **When**: 执行 AutoMigrate
- **Then**: 存量 Token 行的 `ip_policy` 列为 NULL，现有请求鉴权流程不受影响

---

## 【依赖与风险】

| 项 | 说明 |
|----|------|
| 依赖 `model/token.go` | 新增 `IpPolicy` 结构体和字段 |
| 依赖 `common/ip_matcher.go` | CIDR 校验工具函数 `ValidateCIDRList` |
| 依赖 `router/api-router.go` | 路由注册 `PUT /:id/ip_policy` |
| 风险：Q1 待确认 | `ips=[]` + `mode` 非空的行为，影响 TC-3 的预期结果 |
| 风险：旧白名单字段 | 与 `tokens.subnet` 等旧字段的共存和优先级需在执行侧（Feature 2）明确 |

---

## 【技术思路】

### 数据模型设计

```go
// new-api/model/token.go

type IpPolicy struct {
    Mode string   `json:"mode"` // "whitelist" | "blacklist" | ""
    Ips  []string `json:"ips"`  // CIDR 或精确 IP 列表
}

// 实现 driver.Valuer（写库时序列化为 JSON）
func (p IpPolicy) Value() (driver.Value, error) {
    return common.Marshal(p)
}

// 实现 sql.Scanner（读库时反序列化）
func (p *IpPolicy) Scan(value interface{}) error {
    switch v := value.(type) {
    case []byte:
        return common.Unmarshal(v, p)
    case string:
        return common.UnmarshalJsonStr(v, p)
    default:
        return nil
    }
}

type Token struct {
    // ...已有字段...
    IpPolicy *IpPolicy `json:"ip_policy" gorm:"type:text"`
}
```

### 接口设计

```
PUT /api/keys/:id/ip_policy
Authorization: Bearer <admin_or_owner_token>
Content-Type: application/json

Request:
{
  "mode": "whitelist" | "blacklist" | "",
  "ips": ["1.2.3.4", "10.0.0.0/8"]
}

Response 200:
{"success":true,"message":"IP 策略更新成功"}

Response 400:
{"success":false,"message":"invalid IP/CIDR: xxx"}

Response 403:
{"success":false,"message":"forbidden"}

Response 404:
{"success":false,"message":"token not found"}
```

### 路由注册位置

```go
// new-api/router/api-router.go
// 在 Token 路由组（需管理员或 Owner 鉴权）中追加：
tokenRoute.PUT("/:id/ip_policy", controller.UpdateTokenIpPolicy)
```

### 校验逻辑顺序

1. 解析 `id` 参数（非法时 400）
2. 解析请求体（JSON 格式错误时 400）
3. 校验 `mode` 枚举值（非法时 400）
4. 校验 `ips` CIDR 格式（逐条校验，第一个非法项时 400 + 具体错误）
5. 查询 Token 是否存在（不存在时 404）
6. 校验操作权限（非 Owner 且非管理员时 403）
7. 持久化（DB 写入失败时 500）
8. 写入操作日志（type=3）
