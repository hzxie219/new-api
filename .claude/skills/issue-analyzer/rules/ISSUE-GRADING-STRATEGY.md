# 问题分级策略 (Issue Grading Strategy)

## 概述

本文档定义代码规范检查中问题的分类分级标准，适用于所有支持的编程语言（Go、Python、Java等）。

## 三级分类体系

### 级别定义

| 级别 | 英文名称 | 严重程度 | 修复优先级 | 典型影响 |
|------|---------|---------|-----------|---------|
| **错误** | Error | 高 | 必须修复 | 功能异常、稳定性问题、性能问题、安全漏洞 |
| **告警** | Warning | 中 | 建议修复 | 代码风格不统一、可读性差、潜在维护问题 |
| **建议** | Suggestion | 低 | 可选修复 | 注释缺失、文档不完善、轻微优化建议 |

---

## 分级规则

### 规则1：注释类问题 → 建议 (Suggestion)

**适用范围**：所有与注释、文档相关的问题

**示例**：
- ✅ 缺少函数注释
- ✅ 缺少包/模块注释
- ✅ 注释格式不规范
- ✅ TODO/FIXME 注释未处理
- ✅ 注释与代码不一致

**Go 语言示例**：
```go
// ❌ Suggestion: 缺少函数注释
func processOrder(id int) error {
    // ...
}

// ✅ 修复后
// processOrder 处理指定ID的订单
func processOrder(id int) error {
    // ...
}
```

**Python 语言示例**：
```python
# ❌ Suggestion: 缺少函数docstring
def process_order(order_id):
    pass

# ✅ 修复后
def process_order(order_id):
    """处理指定ID的订单"""
    pass
```

---

### 规则2：代码风格类问题 → 告警 (Warning)

**适用范围**：所有影响代码可读性和一致性但不影响功能的问题

**示例**：
- ✅ 驼峰命名不规范（camelCase vs snake_case）
- ✅ 缩进不一致
- ✅ 空格使用不规范
- ✅ 行长度超过限制
- ✅ 导入顺序不规范
- ✅ 变量/函数命名不清晰
- ✅ 代码重复（DRY原则）

**Go 语言示例**：
```go
// ❌ Warning: 变量命名应使用驼峰式
var user_name string

// ✅ 修复后
var userName string

// ❌ Warning: 包名不应使用下划线
package dsp_bad_code

// ✅ 修复后
package dspbadcode
```

**Python 语言示例**：
```python
# ❌ Warning: 函数命名应使用snake_case
def getUserName():
    pass

# ✅ 修复后
def get_user_name():
    pass

# ❌ Warning: 赋值符号周围缺少空格
USER_SYNC_ENV_YAML='/etc/user_sync/env.yaml'

# ✅ 修复后
USER_SYNC_ENV_YAML = '/etc/user_sync/env.yaml'
```

---

### 规则3：功能/稳定性/性能/安全类问题 → 错误 (Error)

**适用范围**：所有可能导致程序异常、不稳定、性能下降或安全风险的问题

#### 3.1 功能异常类

**示例**：
- ✅ 未处理的错误返回值
- ✅ 空指针引用风险
- ✅ 类型转换错误
- ✅ 逻辑错误（死代码、无限循环等）
- ✅ 资源未关闭（文件、数据库连接等）

**Go 语言示例**：
```go
// ❌ Error: 未处理错误返回值
file, _ := os.Open("data.txt")

// ✅ 修复后
file, err := os.Open("data.txt")
if err != nil {
    return err
}
defer file.Close()

// ❌ Error: 空指针引用风险
var user *User
fmt.Println(user.Name)  // panic risk

// ✅ 修复后
var user *User
if user != nil {
    fmt.Println(user.Name)
}
```

**Python 语言示例**：
```python
# ❌ Error: 未捕获异常
def read_config():
    with open('config.json') as f:
        return json.load(f)

# ✅ 修复后
def read_config():
    try:
        with open('config.json') as f:
            return json.load(f)
    except (IOError, json.JSONDecodeError) as e:
        logger.error(f"Failed to read config: {e}")
        return {}
```

#### 3.2 稳定性问题类

**示例**：
- ✅ 并发竞态条件
- ✅ 死锁风险
- ✅ 内存泄漏
- ✅ 资源耗尽风险

**Go 语言示例**：
```go
// ❌ Error: 并发访问未加锁
type Counter struct {
    count int
}

func (c *Counter) Increment() {
    c.count++  // race condition
}

// ✅ 修复后
type Counter struct {
    mu    sync.Mutex
    count int
}

func (c *Counter) Increment() {
    c.mu.Lock()
    defer c.mu.Unlock()
    c.count++
}
```

**Python 语言示例**：
```python
# ❌ Error: 资源未正确关闭
def process_file():
    f = open('data.txt')
    data = f.read()
    return data  # file not closed

# ✅ 修复后
def process_file():
    with open('data.txt') as f:
        return f.read()
```

#### 3.3 性能问题类

**示例**：
- ✅ N+1 查询问题
- ✅ 不必要的循环嵌套
- ✅ 低效的算法复杂度
- ✅ 内存使用不当

**Go 语言示例**：
```go
// ❌ Error: N+1 查询问题
func GetOrdersWithUsers(orderIDs []int) []Order {
    orders := []Order{}
    for _, id := range orderIDs {
        order := db.GetOrder(id)        // 第1次查询
        user := db.GetUser(order.UserID) // N次查询
        order.User = user
        orders = append(orders, order)
    }
    return orders
}

// ✅ 修复后：使用JOIN或批量查询
func GetOrdersWithUsers(orderIDs []int) []Order {
    return db.Query(`
        SELECT o.*, u.*
        FROM orders o
        JOIN users u ON o.user_id = u.id
        WHERE o.id IN (?)
    `, orderIDs)
}
```

**Python 语言示例**：
```python
# ❌ Error: 重复计算
def process_items(items):
    for item in items:
        if len(items) > 100:  # len()在每次循环中重复计算
            process(item)

# ✅ 修复后
def process_items(items):
    items_count = len(items)
    for item in items:
        if items_count > 100:
            process(item)
```

#### 3.4 安全问题类

**示例**：
- ✅ SQL 注入风险
- ✅ XSS 跨站脚本风险
- ✅ 路径穿越漏洞
- ✅ 密码/密钥硬编码
- ✅ 不安全的加密算法
- ✅ 敏感信息泄露

**Go 语言示例**：
```go
// ❌ Error: SQL注入风险
func GetUser(username string) (*User, error) {
    query := "SELECT * FROM users WHERE username = '" + username + "'"
    return db.Query(query)
}

// ✅ 修复后：使用参数化查询
func GetUser(username string) (*User, error) {
    query := "SELECT * FROM users WHERE username = ?"
    return db.Query(query, username)
}

// ❌ Error: 密钥硬编码
const apiKey = "sk-1234567890abcdef"

// ✅ 修复后：使用环境变量
var apiKey = os.Getenv("API_KEY")
```

**Python 语言示例**：
```python
# ❌ Error: SQL注入风险
def get_user(username):
    query = f"SELECT * FROM users WHERE username = '{username}'"
    return db.execute(query)

# ✅ 修复后
def get_user(username):
    query = "SELECT * FROM users WHERE username = %s"
    return db.execute(query, (username,))

# ❌ Error: 使用弱加密算法
import md5
hash = md5.md5(password.encode()).hexdigest()

# ✅ 修复后
import hashlib
hash = hashlib.sha256(password.encode()).hexdigest()
```

---

## 内部规范标记映射

### Go 语言规范标记

| 规范标记 | 映射级别 | 说明 |
|---------|---------|------|
| 【强制】 | Error | 必须遵守的规则，违反可能导致功能/安全/稳定性问题 |
| 【建议】 | Warning/Suggestion | 推荐遵守的规则，违反不影响功能但影响代码质量 |
| 【建议】+ 注释相关 | Suggestion | 如果是注释类问题，降级为 Suggestion |
| 【工具标记】如【typecheck】 | 根据问题类型 | 工具检测的问题按上述规则1-3分类 |

**映射逻辑**：
```
if 问题类型 == "注释相关":
    级别 = Suggestion
else if 问题类型 == "代码风格":
    级别 = Warning
else if 问题类型 in ["功能异常", "稳定性", "性能", "安全"]:
    级别 = Error
else if 规范标记 == "【强制】":
    级别 = Error
else if 规范标记 == "【建议】":
    级别 = Warning
```

### Python 语言规范标记

| 规范标记 | 映射级别 | 说明 |
|---------|---------|------|
| 【强制】 | Error | 必须遵守的规则 |
| 【建议】 | Warning/Suggestion | 推荐遵守的规则 |
| PEP 8 规则 | Warning | 代码风格问题（如启用外部规范） |

---

## 修复策略

### 快速模式 (Fast Mode)

**默认行为**：自动修复所有问题（Error + Warning + Suggestion）

```bash
/lint main
# 自动执行修复，无需人工干预
```

### 深度模式 (Deep Mode)

**默认行为**：生成报告，需用户手动确认后修复

```bash
/lint --mode=deep main
# 查看报告
/fix --report lint-incremental-deep-go-20251222.md
```

### 修复级别参数 (--level)

**语法**：`/fix --level=<error|warning|all>`

**级别说明**：

| 参数 | 修复范围 | 说明 |
|------|---------|------|
| `--level=error` | 仅 Error | 只修复错误级别问题 |
| `--level=warning` | Warning + Error | 修复告警及以上级别（告警 + 错误） |
| `--level=all` 或省略 | All | 修复所有问题（错误 + 告警 + 建议） |

**使用示例**：

```bash
# 快速模式 - 自动修复所有问题
/lint main
# 等同于：/lint main && /fix --level=all

# 深度模式 - 只修复错误级别
/lint --mode=deep main
/fix --report lint-incremental-deep-go-20251222.md --level=error

# 深度模式 - 修复告警及以上（告警 + 错误）
/fix --report lint-incremental-deep-python-20251222.md --level=warning

# 深度模式 - 修复所有问题
/fix --report lint-full-deep-go-20251222.md --level=all
```

**层级修复逻辑**：

```
--level=error    →  修复: [Error]
--level=warning  →  修复: [Error, Warning]
--level=all      →  修复: [Error, Warning, Suggestion]
```

---

## 报告中的问题呈现

### 问题标识格式

```markdown
#### [E001] 未处理错误返回值
**级别**: Error
**分类**: 功能异常
...

#### [W015] 变量命名不规范
**级别**: Warning
**分类**: 代码风格
...

#### [S003] 缺少函数注释
**级别**: Suggestion
**分类**: 注释文档
...
```

### 问题统计格式

```markdown
## 问题统计
- 总问题数: 28个
  - 🔴 Error: 5个
  - 🟡 Warning: 15个
  - 🔵 Suggestion: 8个
```

---

## 特殊场景处理

### 场景1：Lint 工具与内部规范冲突

**原则**：内部规范优先

**示例**：
- Lint 工具报告行长度超过 79 字符（PEP 8）
- 内部规范未强制要求行长度限制
- **结果**：如果外部规范已禁用，该问题应被 report-validator 移除

### 场景2：同一问题多个级别标记

**原则**：取最高级别

**示例**：
- 问题既是"代码风格"又是"功能异常"
- **结果**：归类为 Error

### 场景3：外部规范启用时的级别

**原则**：外部规范问题默认为 Warning，除非涉及安全/功能

**示例**：
- PEP 8 行长度限制 → Warning
- PEP 8 + 安全问题（如 eval 使用） → Error

---

## 语言特定规则

### Go 语言

**Error 级别常见问题**：
- 未检查 error 返回值
- 空指针引用
- goroutine 泄漏
- 数据竞态
- defer 使用不当
- context 未传递

**Warning 级别常见问题**：
- 包名使用下划线
- 变量名不符合驼峰规范
- import 顺序不规范
- 代码重复

**Suggestion 级别常见问题**：
- 缺少包注释
- 缺少导出函数注释
- TODO 注释未处理

### Python 语言

**Error 级别常见问题**：
- 未捕获异常
- 资源未关闭
- SQL 注入风险
- 使用 eval/exec
- 不安全的反序列化

**Warning 级别常见问题**：
- 函数名不符合 snake_case
- 赋值符号周围缺少空格
- import 语句顺序不规范
- 行长度超过限制（如启用 PEP 8）

**Suggestion 级别常见问题**：
- 缺少 docstring
- 缺少类型注解
- TODO/FIXME 注释

---

## 验证和质量保证

### report-validator 的职责

1. **验证问题级别准确性**
   - 检查问题分类是否符合本策略
   - 检查级别标记是否正确

2. **去除误报**
   - 移除不符合分级规则的问题
   - 调整级别错误的问题

3. **重新生成阈值**
   - 如果无效问题 ≥ 30%，触发重新生成
   - 如果连续 5 个无效问题，触发重新生成

### code-checker 的职责

1. **正确应用分级规则**
   - 按规则 1-3 对每个问题分级
   - 正确映射内部规范标记

2. **一致性保证**
   - 同类问题应有相同级别
   - 跨文件的相同问题级别一致

---

**最后更新**: 2025-12-26
