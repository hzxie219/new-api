# 【Tech-服务级】model-token

## 服务职责

负责 `tokens` 表的数据模型扩展与数据库迁移。新增 `IpPolicy` 结构体和 `ip_policy` 字段，支持 IP 访问策略的持久化存储，实现 JSON 序列化/反序列化以兼容 SQLite/MySQL/PostgreSQL 三种数据库。

**所属 Tech-系统级**：Tech-系统级-策略持久化与配置接口

---

## 详细验收条件

### AC-1: Token 结构体新增 IpPolicy 字段

- **Given**: 完成 `model/token.go` 代码变更
- **When**: 检查 `Token` 结构体定义
- **Then**:
  - 包含 `IpPolicy *IpPolicy` 字段（指针类型，nil 表示无策略）
  - GORM tag 使用 `gorm:"type:text"`（TEXT 类型，三数据库通用）
  - JSON tag 为 `json:"ip_policy"`

### AC-2: IpPolicy 结构体定义正确

- **Given**: 完成代码变更
- **When**: 检查 `IpPolicy` 结构体
- **Then**:
  - `Mode string`：JSON tag `json:"mode"`，取值语义为 `"whitelist"`/`"blacklist"`/`""`
  - `Ips []string`：JSON tag `json:"ips"`，元素为 IP 或 CIDR 字符串

### AC-3: driver.Valuer 实现正确（写库序列化）

- **Given**: `IpPolicy{Mode:"whitelist", Ips:["1.2.3.4/32"]}`
- **When**: GORM 保存 Token 时调用 `Value()`
- **Then**: 返回 JSON 字节序列 `{"mode":"whitelist","ips":["1.2.3.4/32"]}`，写入数据库；使用 `common.Marshal`，不直接使用 `encoding/json`

### AC-4: sql.Scanner 实现正确（读库反序列化）

- **Given**: 数据库 `ip_policy` 列存储值为 `{"mode":"whitelist","ips":["1.2.3.4/32"]}`
- **When**: GORM 读取 Token 时调用 `Scan()`
- **Then**: `IpPolicy.Mode == "whitelist"`，`IpPolicy.Ips == ["1.2.3.4/32"]`，无数据丢失

### AC-5: Scanner 兼容 []byte 和 string 两种类型

- **Given**: 不同数据库驱动返回 `ip_policy` 字段的类型不同（MySQL 返回 `[]byte`，PostgreSQL/SQLite 可能返回 `string`）
- **When**: 调用 `Scan()` 方法
- **Then**: 两种类型均能正确反序列化；若值为 nil/NULL，`IpPolicy` 字段保持 nil

### AC-6: 数据库迁移兼容三数据库

- **Given**: SQLite / MySQL / PostgreSQL 任一数据库，`tokens` 表已存在（无 `ip_policy` 列）
- **When**: 服务启动执行 `model.AutoMigrate`
- **Then**: `tokens` 表成功 ADD COLUMN `ip_policy TEXT`，存量行该列为 NULL；服务正常启动，无迁移错误

### AC-7: 指针语义 — nil 表示无策略

- **Given**: 存量 Token（`ip_policy` 列为 NULL）
- **When**: 读取该 Token 到 Go 结构体
- **Then**: `token.IpPolicy == nil`；鉴权流程中判断 `IpPolicy == nil` 可跳过 IP 校验

---

## 技术实现

### 代码位置

| 文件 | 变更内容 |
|------|---------|
| `new-api/model/token.go` | 新增 `IpPolicy` 结构体及其方法，在 `Token` 结构体中新增字段 |
| `new-api/model/main.go` | 确认 `AutoMigrate` 参数列表已包含 `&Token{}`（无需新增，仅确认） |

### 核心代码

```go
// new-api/model/token.go

import (
    "database/sql/driver"
    "github.com/songquanpeng/one-api/common"
)

type IpPolicy struct {
    Mode string   `json:"mode"` // "whitelist" | "blacklist" | ""
    Ips  []string `json:"ips"`  // CIDR 或精确 IP 列表，nil 等同于空列表
}

// Value 实现 driver.Valuer，写库时序列化为 JSON bytes
func (p IpPolicy) Value() (driver.Value, error) {
    b, err := common.Marshal(p)
    if err != nil {
        return nil, err
    }
    return string(b), nil // TEXT 类型存 string
}

// Scan 实现 sql.Scanner，读库时从 JSON 反序列化
func (p *IpPolicy) Scan(value interface{}) error {
    if value == nil {
        return nil
    }
    switch v := value.(type) {
    case []byte:
        return common.Unmarshal(v, p)
    case string:
        return common.UnmarshalJsonStr(v, p)
    }
    return nil
}

type Token struct {
    // ...已有字段保持不变...
    IpPolicy *IpPolicy `json:"ip_policy" gorm:"type:text"`
}
```

### 注意事项

1. **必须使用 `common.Marshal / common.Unmarshal`**，不得直接使用 `encoding/json`（项目 Rule 1）
2. 字段为指针类型 `*IpPolicy`，nil 指针表示"未配置策略"，与 `{mode:"",ips:[]}` 语义等价
3. 参考同类实现：`new-api/model/channel.go` 中 JSON 字段的 `Value/Scan` 模式
4. `AutoMigrate` 会自动 ADD COLUMN，不需要手动 SQL；对 SQLite 的 `ALTER TABLE` 限制已由 GORM 处理

---

## 监控与排障

| 场景 | 日志关键字 | 处理方式 |
|------|-----------|---------|
| AutoMigrate 失败 | `[FATAL] database migration` | 服务启动失败，检查数据库连接和权限 |
| Scan 反序列化失败 | `[WARN] ip_policy scan failed` | 返回 nil IpPolicy（容错降级），记录 token_id 和原始值 |
| Value 序列化失败 | `[ERROR] ip_policy marshal failed` | 写库失败，返回 HTTP 500 |
