# 问题分级映射规则 (Issue Grading Mapping Rules)

## 概述

本文档定义内部规范标记到三级分类体系（Error, Warning, Suggestion）的精确映射规则，适用于 Go 和 Python 语言。

---

## 映射优先级

问题分级遵循以下优先级顺序（从高到低）：

```
1. 问题类型分类（注释 / 代码风格 / 功能/稳定性/性能/安全）
   ↓
2. 规范标记映射（【强制】/ 【建议】）
   ↓
3. 工具标记参考（【typecheck】/ 【govet】/ 【pylint】等）
```

---

## 核心映射规则

### 规则 1: 注释类问题 → Suggestion

**判断条件**：
```
IF (规范章节 == "02. comment - 注释" OR
    规范章节 == "1.6. comment - 注释" OR
    规范描述包含 ["注释", "comment", "文档字符串", "docstring", "godoc", "TODO", "FIXME"])
THEN
    级别 = Suggestion
```

**Go 语言示例**：

| 规范ID | 规范内容 | 原标记 | 映射级别 | 理由 |
|--------|---------|--------|---------|------|
| 2.1 | 单行注释符号与内容之间，用一个空格隔开 | 【建议】【gocritic】 | **Suggestion** | 注释类问题 |
| 2.2 | 每个可导出的名字都要有注释 | 【建议】 | **Suggestion** | 注释类问题 |
| 2.2 | 包注释 | 【建议】 | **Suggestion** | 注释类问题 |
| 2.2 | 函数/方法注释 | 【建议】 | **Suggestion** | 注释类问题 |

**Python 语言示例**：

| 规范ID | 规范内容 | 原标记 | 映射级别 | 理由 |
|--------|---------|--------|---------|------|
| 1.6 | 模块文档字符串(docstring) | 【建议】 | **Suggestion** | 注释类问题 |
| 1.6 | 函数文档字符串应注释 | 【建议】 | **Suggestion** | 注释类问题 |
| 1.6 | 类文档注释(docstring) | 【建议】 | **Suggestion** | 注释类问题 |
| 1.6 | 行注释需要在行上面做注释 | 【建议】 | **Suggestion** | 注释类问题 |

---

### 规则 2: 代码风格类问题 → Warning

**判断条件**：
```
IF (规范描述包含 ["命名", "naming", "缩进", "indent", "空格", "space", "折行", "wrap",
                  "括号", "parenthesis", "编码格式", "coding_format", "驼峰", "camelCase",
                  "snake_case", "下划线", "大小写", "import顺序", "分组", "group"] AND
    问题类型 != "功能异常" AND
    问题类型 != "安全问题")
THEN
    级别 = Warning
```

**Go 语言示例**：

| 规范ID | 规范内容 | 原标记 | 映射级别 | 理由 |
|--------|---------|--------|---------|------|
| 1.3 | 缩进、括号和空格使用gofmt工具处理 | 【强制】【typecheck】 | **Warning** | 代码风格（尽管标记为【强制】，但是格式问题） |
| 1.4 | 标识符应采用MixedCaps/mixedCaps规则 | 【强制】 | **Warning** | 命名风格问题 |
| 1.4 | 包名：不应该包含下划线 | 【建议】 | **Warning** | 命名风格问题 |
| 1.4 | 文件名：全部使用小写字母，下划线分隔单词 | 【强制】 | **Warning** | 命名风格问题 |
| 1.4 | 变量/常量名：不以类型作为前后缀 | 【强制】 | **Warning** | 命名风格问题 |
| 1.5 | 编码格式必须为UTF-8 | 【强制】 | **Warning** | 编码格式风格问题 |
| 5.3 | import导入使用goimports规范 | 【强制】 | **Warning** | 导入顺序风格问题 |

**Python 语言示例**：

| 规范ID | 规范内容 | 原标记 | 映射级别 | 理由 |
|--------|---------|--------|---------|------|
| 1.1 | 模块与包命名全部使用小写 | 【强制】 | **Warning** | 命名风格问题 |
| 1.1 | 类命名采用CapWords | 【强制】 | **Warning** | 命名风格问题 |
| 1.1 | 函数与方法命名全部使用小写+下划线 | 【强制】 | **Warning** | 命名风格问题 |
| 1.1 | 变量使用全小写加下划线 | 【强制】 | **Warning** | 命名风格问题 |
| 1.1 | 常量命名使用全部大写 | 【强制】 | **Warning** | 命名风格问题 |
| 1.2 | 缩进采用4个空格 | 【强制】 | **Warning** | 缩进风格问题 |
| 1.3 | 类和top-level函数定义之间空两行 | 【强制】 | **Warning** | 空行风格问题 |
| 1.4 | 续行应采用圆括号 | 【强制】 | **Warning** | 折行风格问题 |
| 1.7 | 操作符左右各加一个空格 | 【强制】 | **Warning** | 空格风格问题 |
| 1.9 | 编码格式必须为UTF-8无BOM | 【强制】 | **Warning** | 编码格式风格问题 |
| 4.1 | 导入始终放在文件顶部 | 【强制】 | **Warning** | 导入顺序风格问题 |

---

### 规则 3: 功能/稳定性/性能/安全类问题 → Error

**判断条件**：
```
IF (规范描述包含 ["错误处理", "exception", "error", "panic", "资源管理", "resource",
                  "内存泄漏", "memory leak", "并发", "concurrent", "goroutine", "channel",
                  "安全", "security", "注入", "injection", "权限", "permission",
                  "性能", "performance", "优化", "SQL", "XSS", "nil", "空指针",
                  "数据竞态", "race", "死锁", "deadlock", "硬编码", "hardcode",
                  "弱加密", "weak crypto", "未捕获", "未处理", "unhandled"] OR
    规范章节 IN ["03. exception", "04. resource", "2.1. raise", "2.2. catch",
                  "05. concurrent", "07. security"])
THEN
    级别 = Error
```

**Go 语言 - 功能异常类**：

| 规范ID | 规范内容 | 原标记 | 映射级别 | 理由 |
|--------|---------|--------|---------|------|
| 3.1 | 如果函数需要返回错误信息，统一将错误信息作为最后一个返回值返回 | 【强制】【revive】 | **Error** | 功能异常 - 错误处理不规范 |
| 3.3 | 显式处理error，或使用空白标识符忽略 | 【强制】 | **Error** | 功能异常 - 未处理错误 |
| 3.3 | 出现错误时，应该立即处理error | 【强制】 | **Error** | 功能异常 - 错误处理不规范 |
| 3.5 | 类型断言使用comma ok样式 | 【强制】 | **Error** | 功能异常 - 可能panic |
| 3.6 | panic必须在当前Goroutine被捕获 | 【强制】 | **Error** | 功能异常 - panic处理不当 |
| 3.6 | recover必须在defer中使用 | 【强制】 | **Error** | 功能异常 - recover使用不当 |
| 4.1 | 资源使用完毕后，必须进行额外的检查以确保已经关闭 | 【强制】 | **Error** | 功能异常 - 资源泄漏 |
| 4.1 | 如果资源管理存在失败时，需要先判断错误再defer释放 | 【强制】 | **Error** | 功能异常 - 资源管理不当 |
| 4.1 | 禁止在循环中直接使用defer | 【强制】 | **Error** | 功能异常 - 资源延迟释放 |
| 7.2 | 禁止在闭包中直接调用循环变量 | 【强制】 | **Error** | 功能异常 - 指针问题 |
| 7.2 | 进行指针操作时，须判断该指针是否为nil | 【强制】 | **Error** | 功能异常 - 空指针引用 |

**Go 语言 - 稳定性/性能类**：

| 规范ID | 规范内容 | 原标记 | 映射级别 | 理由 |
|--------|---------|--------|---------|------|
| 4.2 | 避免不必要的循环引用 | 【强制】 | **Error** | 稳定性 - 内存泄漏 |
| 4.3 | 禁止在闭包中直接调用循环变量 | 【强制】 | **Error** | 稳定性 - 数据竞态 |
| 4.3 | map和slice不是并发安全的，多goroutine读写时必须同步 | 【强制】 | **Error** | 稳定性 - 并发安全 |
| 7.5 | 遍历较大元素，优先使用下标遍历 | 【强制】 | **Error** | 性能 - 元素拷贝开销 |
| 7.7 | switch语句必须有default | 【强制】 | **Error** | 功能异常 - 分支遗漏 |
| 7.9 | 对slice进行索引操作时，必须判断长度是否合法 | 【强制】 | **Error** | 功能异常 - 越界风险 |
| 7.11 | channel通常size应为1或是无缓冲的 | 【强制】 | **Error** | 稳定性 - 并发设计问题 |
| 7.11 | channel必须使用make初始化 | 【强制】 | **Error** | 功能异常 - nil channel阻塞 |
| 7.11 | 禁止多次关闭channel或向已关闭channel写入 | 【强制】 | **Error** | 功能异常 - panic风险 |
| 7.12 | 不要在代码中泄漏goroutine | 【强制】 | **Error** | 稳定性 - goroutine泄漏 |
| 7.12 | 创建协程池时，需对协程池进行最大数量做限制 | 【强制】 | **Error** | 稳定性 - 资源耗尽 |

**Go 语言 - 安全类**：

| 规范ID | 规范内容 | 原标记 | 映射级别 | 理由 |
|--------|---------|--------|---------|------|
| 7.2 | 在进行数字运算操作时，需要做好长度限制 | 【强制】 | **Error** | 安全 - 整数溢出 |
| 7.3 | 除了0和1，不要使用魔法数字 | 【强制】【gomnd】 | **Error** | 安全 - 代码可维护性（可能导致逻辑错误） |
| 7.14 | 格式化输出字符串时，如果输入数据来自外部，需使用%q进行安全转义 | 【强制】 | **Error** | 安全 - 注入风险 |
| 7.28 | 使用crypto/rand代替math/rand生成随机数 | 【强制】 | **Error** | 安全 - 弱随机数 |
| 7.29 | 禁止使用unsafe包，cgo场景例外 | 【强制】 | **Error** | 安全 - 内存安全破坏 |

**Python 语言 - 功能异常类**：

| 规范ID | 规范内容 | 原标记 | 映射级别 | 理由 |
|--------|---------|--------|---------|------|
| 2.1 | 对外部读取数据增加try/except处理 | 【强制】 | **Error** | 功能异常 - 未捕获异常 |
| 2.1 | 抛出异常时使用raise ValueError('message') | 【强制】 | **Error** | 功能异常 - 异常抛出不规范 |
| 2.2 | 不要使用except:捕获所有异常 | 【强制】 | **Error** | 功能异常 - 异常捕获过宽 |
| 2.2 | 异常必须同级的try catch捕获 | 【强制】 | **Error** | 功能异常 - 异常无法捕获 |
| 2.3 | 将无论异常与否都应执行的代码放在finally里 | 【强制】 | **Error** | 功能异常 - 资源清理不当 |
| 2.3 | 严禁在finally中使用return/break/continue | 【强制】 | **Error** | 功能异常 - 流控制异常 |
| 2.3 | 严禁在exception里通过pass忽略异常 | 【强制】 | **Error** | 功能异常 - 错误被隐藏 |
| 2.5 | assert只能用来保证内部正确性 | 【强制】 | **Error** | 功能异常 - 断言滥用 |
| 2.5 | 断言中禁止调用有副作用的函数 | 【强制】 | **Error** | 功能异常 - 副作用未执行 |
| 3.1 | 禁止通过默认参数实现每次调用获取不同值 | 【强制】 | **Error** | 功能异常 - 默认参数陷阱 |
| 3.1 | 不要使用可变对象作为函数默认值 | 【强制】 | **Error** | 功能异常 - 默认参数共享 |
| 4.2 | 代码必须检查if __name__ == '__main__' | 【强制】 | **Error** | 功能异常 - 模块导入执行 |
| 4.4 | __init__.py禁止实现业务逻辑 | 【强制】 | **Error** | 功能异常 - 导入副作用 |

**Python 语言 - 稳定性/性能类**：

| 规范ID | 规范内容 | 原标记 | 映射级别 | 理由 |
|--------|---------|--------|---------|------|
| 4.9 | 禁止在闭包中直接进行变量绑定 | 【强制】 | **Error** | 稳定性 - 闭包变量陷阱 |
| 4.18 | 禁止使用全局变量（特殊场景例外） | 【强制】 | **Error** | 稳定性 - 全局状态污染 |
| 4.19 | 全局变量需要进行初始化 | 【强制】 | **Error** | 功能异常 - 未初始化变量 |
| 4.19 | 类成员变量需要在构造函数中初始化 | 【强制】 | **Error** | 功能异常 - 未初始化变量 |
| 4.20 | 在遍历列表时不允许删除元素 | 【强制】 | **Error** | 功能异常 - 迭代器失效 |
| 4.21 | 使用x += y代替x = x + y（可变对象） | 【强制】 | **Error** | 性能 - 不必要的拷贝 |
| 4.26 | 文件、套接字等资源使用完后应显式关闭 | 【强制】 | **Error** | 功能异常 - 资源泄漏 |
| 5.1 | 当父任务启动协程时，结束时必须进行协程回收 | 【强制】 | **Error** | 稳定性 - 协程泄漏 |
| 5.1 | 禁止在协程内部调用阻塞函数 | 【强制】 | **Error** | 稳定性 - 协程阻塞 |
| 5.2 | 涉及大对象频繁更新，禁止使用multiprocessing | 【强制】 | **Error** | 性能 - 序列化开销 |

**Python 语言 - 安全类**：

| 规范ID | 规范内容 | 原标记 | 映射级别 | 理由 |
|--------|---------|--------|---------|------|
| 4.25 | 避免使用魔数，用常量代替 | 【强制】 | **Error** | 安全 - 代码可维护性（可能导致逻辑错误） |
| 7.2 | 不建议对密码进行硬编码 | 【强制】 | **Error** | 安全 - 密码泄露 |
| 7.2 | 不建议对生成密码的函数默认参数进行硬编码 | 【强制】 | **Error** | 安全 - 密码泄露 |
| 7.4 | 以root权限执行，权限过高 | 【强制】 | **Error** | 安全 - 权限过高 |
| 7.5 | eval可以任意执行字符串，导致命令注入 | 【强制】 | **Error** | 安全 - 命令注入 |
| 7.5 | input会从标准输入读取并执行代码 | 【强制】 | **Error** | 安全 - 命令注入 |
| 7.5 | 禁止使用linux命令通配符 | 【强制】 | **Error** | 安全 - 命令注入 |
| 7.6 | random模块不能用于安全加密 | 【强制】 | **Error** | 安全 - 弱随机数 |
| 7.6 | DSA密钥大小低于1024位是不安全的 | 【强制】 | **Error** | 安全 - 弱加密 |
| 7.7 | 不使用系统自带的xml库，建议使用defusedxml | 【强制】 | **Error** | 安全 - XML解析漏洞 |
| 7.8 | request使用时未开启ssl认证 | 【强制】 | **Error** | 安全 - SSL认证缺失 |
| 7.9 | 不使用sql字符串拼接的查询语句 | 【强制】 | **Error** | 安全 - SQL注入 |
| 7.10 | 设置flask模板jinja2的autoescape为false | 【强制】 | **Error** | 安全 - XSS漏洞 |
| 4.33 | 不使用yaml_load反序列化执行程序 | 【强制】 | **Error** | 安全 - 代码注入 |
| 4.33 | 不使用pickle模块序列化 | 【强制】 | **Error** | 安全 - 代码注入 |

---

## 特殊场景处理规则

### 场景 1: 【强制】标记但属于代码风格

**规则**：即使标记为【强制】，如果问题本质是代码风格（命名、缩进、空格等），仍然归类为 **Warning**

**示例**：
- Go: 1.4 标识符应采用MixedCaps规则【强制】 → **Warning**（命名风格）
- Python: 1.1 函数命名全部使用小写+下划线【强制】 → **Warning**（命名风格）
- Go: 1.3 缩进使用gofmt工具处理【强制】 → **Warning**（格式风格）

### 场景 2: 【建议】标记但涉及功能/安全

**规则**：即使标记为【建议】，如果问题涉及功能异常、稳定性、性能或安全，仍然归类为 **Error**

**示例**：
- Go: 4.2 基于切片创建子切片时需注意内存泄漏【建议】 → **Error**（内存泄漏）
- Python: 2.1 except和finally子句中不得抛出新的未知异常【建议】 → **Error**（异常覆盖）
- Python: 2.4 LOG.exception只能用于except作用域内【建议】 → **Error**（日志错误）

### 场景 3: 工具标记的处理

**规则**：工具标记（如【typecheck】、【govet】、【pylint】）仅作为参考，最终级别由问题类型决定

**示例**：
- Go: 1.3 缩进处理【强制】【typecheck】 → **Warning**（虽有typecheck，但是格式问题）
- Go: 3.1 错误信息作为最后返回值【强制】【revive】 → **Error**（虽有revive，但是功能问题）
- Go: 7.3 魔法数字【强制】【gomnd】 → **Error**（虽有gomnd，但可能导致逻辑错误）

### 场景 4: 同时满足多个规则

**规则**：当一个问题同时满足多个规则时，取最高级别

**优先级**：Error > Warning > Suggestion

**示例**：
```
问题："函数未添加注释，且函数名不符合命名规范"

规则1判断：未添加注释 → Suggestion
规则2判断：命名不规范 → Warning

最终级别：Warning（取最高级别）
```

### 场景 5: 外部规范（PEP 8, Effective Go）问题

**规则**：
- 如果 `external_standards.enabled = false`，外部规范问题应被 report-validator 移除
- 如果 `external_standards.enabled = true`：
  - 注释类 → Suggestion
  - 代码风格类 → Warning
  - 功能/安全类 → Error

**示例**：
```python
# 外部规范启用时
PEP 8: 行长度超过79字符 → Warning（代码风格）
PEP 8: 禁止使用eval() → Error（安全问题）
```

---

## 完整映射算法

```python
def classify_issue(issue, rule, internal_standard):
    """
    对问题进行分级

    Args:
        issue: 问题对象，包含描述、代码片段等
        rule: 规范规则对象，包含ID、标记、描述等
        internal_standard: 内部规范配置

    Returns:
        str: "error" | "warning" | "suggestion"
    """

    # 步骤1：检查是否为注释类问题
    comment_keywords = ["注释", "comment", "文档", "document", "docstring",
                       "godoc", "TODO", "FIXME", "文档字符串"]

    if (rule.chapter in ["02. comment", "1.6. comment"] or
        any(kw in rule.description for kw in comment_keywords)):
        return "suggestion"

    # 步骤2：检查是否为功能/稳定性/性能/安全问题
    critical_keywords = [
        # 功能异常
        "错误处理", "exception", "error", "panic", "未处理", "unhandled",
        "空指针", "nil", "null", "越界", "overflow", "未捕获", "未初始化",

        # 稳定性
        "资源管理", "resource", "内存泄漏", "memory leak", "goroutine泄漏",
        "并发", "concurrent", "数据竞态", "race", "死锁", "deadlock",
        "channel", "goroutine", "协程", "线程",

        # 性能
        "性能", "performance", "优化", "N+1", "循环", "拷贝",

        # 安全
        "安全", "security", "注入", "injection", "SQL", "XSS", "CSRF",
        "权限", "permission", "硬编码", "hardcode", "密码", "password",
        "加密", "crypto", "随机数", "random", "root", "unsafe"
    ]

    critical_chapters = [
        "03. exception", "04. resource", "05. concurrent",
        "2.1. raise", "2.2. catch", "2.3. exception",
        "07. security", "5.1. thread_coroutine", "5.2. process"
    ]

    if (rule.chapter in critical_chapters or
        any(kw in rule.description for kw in critical_keywords) or
        any(kw in issue.description for kw in critical_keywords)):
        return "error"

    # 步骤3：检查是否为代码风格问题
    style_keywords = [
        "命名", "naming", "缩进", "indent", "空格", "space", "折行", "wrap",
        "括号", "parenthesis", "编码格式", "coding_format", "驼峰", "camelCase",
        "snake_case", "下划线", "大小写", "import顺序", "分组", "group",
        "格式化", "format", "gofmt", "black", "空行", "blank"
    ]

    if any(kw in rule.description for kw in style_keywords):
        return "warning"

    # 步骤4：根据规范标记决定
    if rule.marker == "【强制】":
        # 强制标记，但前面的检查都没命中，默认为warning
        return "warning"
    else:  # 【建议】
        return "warning"


def apply_classification(lint_report):
    """
    对整个lint报告中的所有问题进行分级
    """
    for issue in lint_report.issues:
        # 获取对应的规范规则
        rule = get_rule_by_id(issue.rule_id)

        # 分级
        issue.level = classify_issue(issue, rule, internal_standards)

        # 记录分级依据
        issue.classification_reason = get_classification_reason(issue, rule)

    return lint_report
```

---

## 分级统计规则

在报告中统计问题时，按级别分组：

```markdown
## 问题统计
- 总问题数: 47个
  - 🔴 Error: 12个 (25.5%)
    - 功能异常: 5个
    - 稳定性: 3个
    - 性能: 2个
    - 安全: 2个
  - 🟡 Warning: 28个 (59.6%)
    - 命名规范: 15个
    - 代码格式: 8个
    - 导入顺序: 5个
  - 🔵 Suggestion: 7个 (14.9%)
    - 缺少注释: 5个
    - TODO未处理: 2个
```

---

## 验证和测试

### code-checker 验证点

1. **每个问题都必须有明确的level字段**：`"level": "error" | "warning" | "suggestion"`
2. **分级依据必须记录**：`"classification_reason": "注释类问题"` 或 `"功能异常 - 未处理错误"`
3. **同类问题级别一致**：相同规范ID的问题应有相同级别

### report-validator 验证点

1. **验证分级准确性**：
   - 注释类问题必须为 Suggestion
   - 安全问题必须为 Error
   - 命名风格问题必须为 Warning

2. **验证外部规范问题**：
   - 如果 `external_standards.enabled = false`，PEP 8 / Effective Go 问题应被移除
   - 如果启用，应按规则正确分级

3. **验证问题一致性**：
   - 相同规范ID的问题应有相同级别
   - 跨文件的相同问题应有相同级别

---

## 常见问题 (FAQ)

### Q1: 为什么【强制】标记的格式问题是 Warning 而不是 Error？

**A**: 问题类型优先于规范标记。格式问题（缩进、空格、命名）不影响程序功能，只影响代码可读性和一致性，因此归类为 Warning。

### Q2: 如何处理一个问题同时涉及多个分类？

**A**: 取最高级别。例如，"函数未处理错误且命名不规范" → Error（因为包含"未处理错误"这个功能异常问题）。

### Q3: 外部规范（PEP 8, Effective Go）的问题如何分级？

**A**:
- 如果 `external_standards.enabled = false`，这些问题不应出现在报告中（被 report-validator 移除）
- 如果启用，按照规则1-3分级：注释→Suggestion，风格→Warning，功能/安全→Error

### Q4: 工具标记（如【typecheck】、【govet】）如何影响分级？

**A**: 工具标记仅作为参考，不直接决定级别。最终级别由问题类型决定。

### Q5: 如何判断一个问题是"代码风格"还是"功能异常"？

**A**:
- **代码风格**：不影响程序运行结果，只影响代码可读性和一致性（如命名、缩进、空格）
- **功能异常**：可能导致程序错误、崩溃或行为异常（如未处理错误、空指针、资源泄漏）

---

**最后更新**: 2025-12-26
  - 包含特殊场景处理规则
