# 【Tech-服务级】model-token

## 服务职责

负责 `tokens` 表的数据模型扩展与数据库迁移，新增 `ip_policy` 字段以支持 IP 访问策略的持久化存储。

## 所属 Tech-系统级

Tech-系统级-策略持久化与配置接口

## 详细验收条件

### AC-1: Token 结构体新增 IpPolicy 字段

- **Given**: 完成代码变更
- **When**: 检查 `new-api/model/token.go` 中的 `Token` 结构体
- **Then**:
  - 包含 `IpPolicy *IpPolicy` 字段
  - gorm tag 使用 `type:text` 或等价（确保三数据库兼容）
  - JSON tag 为 `json:"ip_policy"`

### AC-2: IpPolicy 类型定义正确

- **Given**: 完成代码变更
- **When**: 检查 `IpPolicy` 结构体
- **Then**:
  - `Mode string`：取值 `"whitelist"`/`"blacklist"`/`""`
  - `Ips []string`：IP 或 CIDR 字符串列表
  - 实现 `driver.Valuer` 和 `sql.Scanner` 接口（JSON 序列化/反序列化）

### AC-3: 数据库迁移兼容三数据库

- **Given**: SQLite / MySQL / PostgreSQL 环境下启动服务
- **When**: 服务启动触发 `model/main.go` 中的 `AutoMigrate`
- **Then**: `tokens` 表成功新增 `ip_policy` 列（TEXT 类型），存量数据行该列默认为 NULL

### AC-4: IpPolicy JSON 序列化正确

- **Given**: `IpPolicy{Mode:"whitelist", Ips:["1.2.3.4/32"]}`
- **When**: 保存到数据库后读取
- **Then**: 读取结果 Mode 和 Ips 与写入一致，无数据丢失

## 技术实现

### 代码位置

- 模型定义: `new-api/model/token.go`
- 迁移注册: `new-api/model/main.go`（在 `AutoMigrate` 列表中确认 `&Token{}` 已包含）

### 核心代码参考

```go
// new-api/model/token.go

type IpPolicy struct {
    Mode string   `json:"mode"`
    Ips  []string `json:"ips"`
}

// 实现 driver.Valuer — 写入数据库时序列化为 JSON
func (p IpPolicy) Value() (driver.Value, error) {
    return common.Marshal(p)
}

// 实现 sql.Scanner — 从数据库读取时反序列化
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

### 注意事项

- 必须使用 `common.Marshal / common.Unmarshal`，不得直接使用 `encoding/json`（项目 Rule 1）
- 字段为指针类型 `*IpPolicy`，nil 表示未设置策略（与空策略 `{mode:"",ips:[]}` 语义相同）
- 参考同类实现: `new-api/model/channel.go:71`（`ChannelInfo` 的 JSON 字段处理）

## 监控与排障

- 迁移失败时日志关键字: `[FATAL] database migration`
- 读取时反序列化失败会返回 nil IpPolicy（容错处理），建议补充 WARN 日志
