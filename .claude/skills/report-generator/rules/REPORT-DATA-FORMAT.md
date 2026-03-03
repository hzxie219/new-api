# 代码规范检查报告数据格式规范

本文档定义了 code-checker skills 与 report-generator skill 之间的标准数据格式。

## 数据格式版本

数据格式版本: 2.0

## 数据结构

### 完整数据示例

```json
{
  "version": "1.0.0",
  "metadata": {
    "language": "go",
    "mode": "incremental",
    "current_branch": "feature-auth",
    "base_branch": "main",
    "check_time": "2025-11-27 14:30:00",
    "total_files": 3,
    "total_issues": 42
  },
  "statistics": {
    "errors": 18,
    "warnings": 16,
    "suggestions": 8,
    "by_category": {
      "naming": {"errors": 4, "warnings": 3, "suggestions": 1},
      "error_handling": {"errors": 6, "warnings": 2, "suggestions": 0},
      "code_format": {"errors": 1, "warnings": 2, "suggestions": 2},
      "documentation": {"errors": 2, "warnings": 5, "suggestions": 2},
      "idioms": {"errors": 3, "warnings": 2, "suggestions": 1},
      "code_quality": {"errors": 2, "warnings": 2, "suggestions": 2}
    }
  },
  "files": [
    {
      "path": "src/app/main.go",
      "issues": [
        {
          "id": "E001",
          "level": "error",
          "category": "naming",
          "title": "包名使用下划线违反 Go 规范",
          "location": "src/app/main.go:1",
          "line_number": 1,
          "description": "包名 `dsp_bad_code_example` 包含下划线，严重违反 Go 命名规范。",
          "current_code": "package dsp_bad_code_example",
          "suggested_code": "package dspbadcode",
          "explanation": "Go 包名应该简短、小写、单个单词，不使用下划线或驼峰。",
          "reference": "Effective Go - Package names",
          "reference_url": "https://go.dev/doc/effective_go#package-names"
        },
        {
          "id": "E006",
          "level": "error",
          "category": "error_handling",
          "title": "忽略错误返回值",
          "location": "src/app/main.go:39",
          "line_number": 39,
          "description": "`ioutil.ReadAll` 的错误被忽略，可能导致后续处理异常数据。",
          "current_code": "inMsg, _ := ioutil.ReadAll(r.Body)",
          "suggested_code": "inMsg, err := ioutil.ReadAll(r.Body)\nif err != nil {\n\thttp.Error(w, fmt.Sprintf(\"failed to read request body: %v\", err), http.StatusBadRequest)\n\treturn\n}",
          "explanation": "忽略错误可能导致程序在遇到问题时无法正确处理。",
          "reference": "Effective Go - Errors",
          "reference_url": "https://go.dev/doc/effective_go#errors"
        }
      ]
    }
  ],
  "recommendations": {
    "priority_1": [
      {
        "location": "src/app/main.go:1",
        "description": "包名重命名：将 `dsp_bad_code_example` 改为符合规范的包名"
      },
      {
        "location": "src/app/main.go:39",
        "description": "处理 `ioutil.ReadAll` 的错误返回"
      }
    ],
    "priority_2": [
      {
        "location": "src/app/main.go:43",
        "description": "使用有意义的变量名"
      }
    ],
    "priority_3": [
      {
        "location": "src/app/main.go:50",
        "description": "简化函数参数"
      }
    ]
  },
  "quality_score": {
    "total": 45,
    "dimensions": {
      "naming": {"score": 4, "description": "包名使用下划线，严重违规"},
      "error_handling": {"score": 4, "description": "多处忽略错误，缺少类型断言检查"},
      "code_organization": {"score": 5, "description": "存在代码重复和未使用函数"},
      "documentation": {"score": 3, "description": "大量导出标识符缺少注释"},
      "concurrency": {"score": 6, "description": "使用了 mutex，但有局部变量误用"},
      "idioms": {"score": 5, "description": "使用了已弃用的 API"},
      "maintainability": {"score": 5, "description": "硬编码较多，缺少配置验证"},
      "performance": {"score": 6, "description": "基本合理，有性能统计"}
    }
  },
  "references": [
    {
      "title": "Effective Go",
      "url": "https://go.dev/doc/effective_go"
    },
    {
      "title": "Go Code Review Comments",
      "url": "https://go.dev/wiki/CodeReviewComments"
    },
    {
      "title": "Uber Go Style Guide",
      "url": "https://github.com/uber-go/guide/blob/master/style.md"
    }
  ]
}
```

## 字段说明

### metadata（元数据）

| 字段 | 类型 | 必需 | 说明 |
|------|------|------|------|
| language | string | ✅ | 编程语言：`go`, `python`, `java` |
| mode | string | ✅ | 检查模式：`incremental`, `full`, `latest` |
| current_branch | string | ✅ | 当前分支名 |
| base_branch | string | ⚠️ | 基准分支名（incremental 模式必需） |
| check_time | string | ✅ | 检查时间（格式：YYYY-MM-DD HH:MM:SS） |
| total_files | number | ✅ | 检查的文件总数 |
| total_issues | number | ✅ | 发现的问题总数 |

### statistics（统计信息）

| 字段 | 类型 | 必需 | 说明 |
|------|------|------|------|
| errors | number | ✅ | 错误级别问题数量 |
| warnings | number | ✅ | 警告级别问题数量 |
| suggestions | number | ✅ | 建议级别问题数量 |
| by_category | object | ✅ | 按类别统计的问题数量 |

### files（文件列表）

每个文件对象包含：

| 字段 | 类型 | 必需 | 说明 |
|------|------|------|------|
| path | string | ✅ | 文件的完整相对路径（从项目根目录开始） |
| issues | array | ✅ | 该文件中发现的问题列表 |

### issues（问题对象）

每个问题对象包含：

| 字段 | 类型 | 必需 | 说明 |
|------|------|------|------|
| id | string | ✅ | 问题 ID（如 E001, W001, S001） |
| level | string | ✅ | 问题级别：`error`, `warning`, `suggestion` |
| category | string | ✅ | 问题类别（naming, error_handling 等） |
| title | string | ✅ | 问题标题 |
| location | string | ✅ | 问题位置（格式：完整路径:行号） |
| line_number | number | ✅ | 行号 |
| description | string | ✅ | 问题详细描述 |
| current_code | string | ⚠️ | 当前代码（可选，有代码示例时提供） |
| suggested_code | string | ⚠️ | 建议的代码（可选） |
| explanation | string | ⚠️ | 额外说明（可选） |
| reference | string | ✅ | 规范参考名称（强烈推荐，帮助开发者理解规范依据） |
| reference_url | string | ✅ | 规范参考链接（强烈推荐，提供权威规范文档） |
| detection_trace | object | ✅ | **检测过程追溯信息**（v2.0新增，记录问题如何被检测出来） |

### recommendations（整改建议）

| 字段 | 类型 | 必需 | 说明 |
|------|------|------|------|
| priority_1 | array | ✅ | 必须修改的问题（错误级别） |
| priority_2 | array | ✅ | 建议修改的问题（警告级别） |
| priority_3 | array | ✅ | 可选优化的问题（建议级别） |

每个建议对象包含：
- `location`: 问题位置（完整路径:行号）
- `description`: 简短描述

### quality_score（代码质量评分）

| 字段 | 类型 | 必需 | 说明 |
|------|------|------|------|
| total | number | ⚠️ | 总体评分（0-100，可选） |
| dimensions | object | ⚠️ | 各维度评分（可选） |

### references（规范参考）

| 字段 | 类型 | 必需 | 说明 |
|------|------|------|------|
| title | string | ✅ | 参考资料标题 |
| url | string | ✅ | 参考资料链接 |

## 问题级别定义

### error（错误）
- **严重程度**: 高
- **必须修改**: 是
- **示例**: 语法错误、严重违反规范、可能导致运行时错误

### warning（警告）
- **严重程度**: 中
- **建议修改**: 强烈建议
- **示例**: 不符合最佳实践、可能的性能问题、缺少文档

### suggestion（建议）
- **严重程度**: 低
- **可选优化**: 可选
- **示例**: 代码优化建议、可读性改进

## 问题类别（按语言）

### Go 语言
- `naming`: 命名规范
- `error_handling`: 错误处理
- `code_format`: 代码格式
- `documentation`: 注释文档
- `idioms`: Go 惯用法
- `code_quality`: 代码质量
- `concurrency`: 并发安全

### Python 语言
- `code_style`: 代码风格
- `naming`: 命名规范
- `imports`: 导入规范
- `documentation`: 注释文档
- `code_quality`: 代码质量
- `best_practices`: 最佳实践

### Java 语言
- `code_format`: 代码格式
- `naming`: 命名规范
- `imports`: 导入规范
- `documentation`: 注释文档
- `exception_handling`: 异常处理
- `best_practices`: 最佳实践

## 规范出处的重要性 📖 ⭐

**为什么需要规范出处？**

每个检查出的问题都应该标注其规范依据，原因包括：

1. **教育价值**: 帮助开发者理解为什么这是个问题，不是主观判断
2. **权威性**: 提供官方或权威规范的链接，增强说服力
3. **学习资源**: 开发者可以深入学习相关规范
4. **团队共识**: 基于公开规范建立团队编码标准
5. **可追溯性**: 方便后续审查和规范更新
6. **内部规范扩展**: 支持添加团队内部规范文档链接

**字段说明**：

- `reference`: 规范名称，如 "Effective Go - Package names"、"PEP 8 - Indentation"
- `reference_url`: 规范链接，指向权威文档的具体章节

**示例**：

```json
{
  "reference": "Effective Go - Error handling",
  "reference_url": "https://go.dev/doc/effective_go#errors"
}
```

**常用规范参考**：

### Go 语言
- [Effective Go](https://go.dev/doc/effective_go)
- [Go Code Review Comments](https://go.dev/wiki/CodeReviewComments)
- [Uber Go Style Guide](https://github.com/uber-go/guide/blob/master/style.md)

### Python 语言
- [PEP 8 - Style Guide for Python Code](https://pep8.org/)
- [Google Python Style Guide](https://google.github.io/styleguide/pyguide.html)
- [The Zen of Python (PEP 20)](https://www.python.org/dev/peps/pep-0020/)

### Java 语言
- [Google Java Style Guide](https://google.github.io/styleguide/javaguide.html)
- [Oracle Java Code Conventions](https://www.oracle.com/java/technologies/javase/codeconventions-contents.html)

**内部规范**：

团队可以添加自己的内部规范文档，例如：
```json
{
  "reference": "团队编码规范 - 错误处理",
  "reference_url": "https://internal-docs.company.com/coding-standards/error-handling"
}
```

---

## 检测过程追溯（v2.0 新增）⭐

### detection_trace 对象结构

每个问题都必须包含 `detection_trace` 字段，记录该问题是如何被检测出来的：

```json
{
  "detection_trace": {
    "source": "internal_standard",
    "detection_method": "ai_analysis",
    "matched_rule": {
      "rule_id": "INTERNAL-GO-1.4",
      "rule_chapter": "01. style - 风格规范 / 1.4. label - 标识符命名",
      "rule_content": "包名：不应该包含下划线",
      "rule_marker": "【建议】",
      "mapped_level": "warning"
    },
    "matching_process": {
      "step1": "读取代码第1行：package dsp_bad_code",
      "step2": "提取包名：dsp_bad_code",
      "step3": "检测到下划线：_",
      "step4": "匹配规则：INTERNAL-GO-1.4（包名不应包含下划线）",
      "step5": "应用分级策略：代码风格类 → Warning"
    },
    "confidence": "high",
    "timestamp": "2025-12-22 10:23:45"
  }
}
```

### 字段说明

| 字段 | 类型 | 必需 | 说明 |
|------|------|------|------|
| source | string | ✅ | 检测来源：`internal_standard`(内部规范) / `external_standard`(外部规范) / `lint_tool`(Lint工具) / `deep_analysis`(深度分析) |
| detection_method | string | ✅ | 检测方法：`ai_analysis`(AI分析) / `tool_output`(工具输出) / `dependency_analysis`(依赖分析) / `call_chain_analysis`(调用链分析) |
| matched_rule | object | ✅ | 命中的规则详情 |
| matching_process | object | ✅ | 匹配过程的详细步骤 |
| confidence | string | ✅ | 置信度：`high` / `medium` / `low` |
| timestamp | string | ✅ | 检测时间戳 |

### matched_rule 对象（内部规范）

```json
{
  "rule_id": "INTERNAL-GO-1.4",
  "rule_chapter": "01. style - 风格规范 / 1.4. label - 标识符命名",
  "rule_content": "包名：不应该包含下划线",
  "rule_marker": "【建议】",
  "mapped_level": "warning"
}
```

### matched_rule 对象（Lint工具）

```json
{
  "tool_name": "golangci-lint",
  "tool_rule": "errcheck",
  "tool_output": "Error return value of `ioutil.ReadAll` is not checked",
  "internal_mapping": "INTERNAL-GO-3.3（显式处理error）",
  "mapped_level": "error"
}
```

### matched_rule 对象（深度分析）

```json
{
  "analysis_type": "dependency_analysis",
  "finding": "函数依赖了未导入的包",
  "related_standard": "INTERNAL-GO-5.3（import规范）",
  "mapped_level": "error"
}
```

### matching_process 示例

**示例1：内部规范匹配**

```json
{
  "step1": "读取 src/app/main.go:42 的代码",
  "step2": "提取变量名：user_name",
  "step3": "检测命名风格：snake_case（包含下划线）",
  "step4": "匹配内部规范：INTERNAL-GO-1.4（变量应使用驼峰命名）",
  "step5": "检查规范标记：【强制】",
  "step6": "应用分级策略：代码风格类 → Warning（尽管是【强制】）",
  "step7": "生成建议：将 user_name 改为 userName"
}
```

**示例2：Lint工具检测**

```json
{
  "step1": "执行 golangci-lint run --enable=errcheck",
  "step2": "工具输出：main.go:39: Error return value of `ioutil.ReadAll` is not checked",
  "step3": "解析工具输出：文件=main.go, 行号=39, 规则=errcheck",
  "step4": "查找内部规范映射：errcheck → INTERNAL-GO-3.3",
  "step5": "读取内部规范：INTERNAL-GO-3.3【强制】显式处理error",
  "step6": "应用分级策略：错误处理类 → Error",
  "step7": "验证：读取第39行代码确认未处理错误"
}
```

**示例3：深度分析检测**

```json
{
  "step1": "分析函数 ProcessOrder 的依赖关系",
  "step2": "发现调用了 database.Query()",
  "step3": "检查 database 包是否已导入",
  "step4": "发现：database 未在 import 中声明",
  "step5": "匹配规范：INTERNAL-GO-5.3（import规范）",
  "step6": "应用分级策略：功能异常类 → Error"
}
```

### 实时输出格式

检测过程中应实时输出追溯信息：

```
🔍 检测文件: src/app/main.go

📍 第1行: package dsp_bad_code
  ├─ 提取包名: dsp_bad_code
  ├─ 检测问题: 包含下划线 '_'
  ├─ 匹配规则: INTERNAL-GO-1.4（包名不应包含下划线）
  ├─ 规范标记: 【建议】
  ├─ 问题类型: 代码风格类
  ├─ 分级结果: Warning
  └─ 问题ID: W001 ✓

📍 第39行: inMsg, _ := ioutil.ReadAll(r.Body)
  ├─ Lint工具: golangci-lint/errcheck
  ├─ 工具输出: Error return value is not checked
  ├─ 内部映射: INTERNAL-GO-3.3（显式处理error）
  ├─ 规范标记: 【强制】
  ├─ 问题类型: 错误处理类
  ├─ 分级结果: Error
  └─ 问题ID: E006 ✓

📊 检测统计:
  - 扫描行数: 150
  - 发现问题: 12个
  - Error: 5个
  - Warning: 7个
  - Suggestion: 0个
```

## 路径格式要求

**关键要求**：所有文件路径和位置信息必须使用从项目根目录开始的完整相对路径。

### 正确示例
```
path: "src/app/services/user_service.go"
location: "src/app/services/user_service.go:25"
```

### 错误示例
```
path: "user_service.go"  ❌
location: "行 25"  ❌
location: "user_service.go:25"  ❌
```

## 版本兼容性

- 向后兼容：新版本可以处理旧版本的数据
- 版本字段：数据中包含 `version` 字段标识格式版本
- 扩展性：可以添加新字段，但不应删除现有必需字段

## 报告格式优化要求（v2.0 新增）⭐

### 1. 问题排序规则

报告中的问题必须按以下顺序展示：

**优先级1：按问题级别排序**
- Error（错误）→ Warning（告警）→ Suggestion（建议）

**优先级2：同级别内按文件路径排序**
- 字母顺序排列文件路径

**优先级3：同文件内按行号排序**
- 从小到大排列行号

**示例顺序**：
```
## 🔴 Error 级别问题

### src/app/main.go

#### 第39行: 未处理错误返回值
...

#### 第125行: 空指针引用风险
...

### src/service/auth.go

#### 第56行: SQL注入风险
...

## 🟡 Warning 级别问题

### src/app/main.go

#### 第1行: 包名使用下划线
...

## 🔵 Suggestion 级别问题
...
```

### 2. 同行多问题合并展示

**规则**：同一行代码如果有多个问题，必须合并在一起展示，不能分散

**错误示例**：
```markdown
#### 第42行: 变量命名不规范
**代码**：`var user_name string`
...

（其他问题）

#### 第42行: 缺少变量注释
**代码**：`var user_name string`
...
```

**正确示例**：
```markdown
#### 第42行: 多个问题

**当前代码**：
```go
var user_name string
```

**问题列表**：

**1. [W023] 变量命名不规范**
- **级别**: Warning
- **分类**: 代码风格 - 命名规范
- **问题描述**: 变量名使用了snake_case，应使用camelCase
- **规范依据**: INTERNAL-GO-1.4【强制】变量应使用驼峰命名
- **建议修改**: `var userName string`

**2. [S015] 缺少变量注释**
- **级别**: Suggestion
- **分类**: 注释文档
- **问题描述**: 导出变量缺少注释说明
- **规范依据**: INTERNAL-GO-2.2【建议】每个可导出的名字都要有注释
- **建议修改**:
```go
// userName 存储用户名称
var userName string
```
\`\`\`

### 3. 检测过程追溯章节

报告末尾必须新增"检测过程追溯"章节，记录所有问题的检测过程：

```markdown
---

## 📋 检测过程追溯

本章节记录所有问题的完整检测过程，供验证和审查。

### E006 - 未处理错误返回值

**问题位置**: src/app/main.go:39

**检测来源**: Lint工具 (golangci-lint)

**检测过程**:
1. ✓ 执行 `golangci-lint run --enable=errcheck`
2. ✓ 工具输出: `main.go:39: Error return value of 'ioutil.ReadAll' is not checked`
3. ✓ 解析输出: 文件=main.go, 行号=39, 规则=errcheck
4. ✓ 查找内部规范映射: errcheck → INTERNAL-GO-3.3
5. ✓ 读取规范内容: "显式处理error，或使用空白标识符忽略"
6. ✓ 检查规范标记: 【强制】
7. ✓ 应用分级策略: 错误处理类 → Error
8. ✓ 验证: 读取第39行代码确认 `inMsg, _ := ioutil.ReadAll(r.Body)` 确实忽略了错误

**置信度**: High

**时间戳**: 2025-12-22 10:23:45

---

### W001 - 包名使用下划线

**问题位置**: src/app/main.go:1

**检测来源**: 内部规范 (INTERNAL-GO-1.4)

**检测过程**:
1. ✓ 读取代码第1行: `package dsp_bad_code`
2. ✓ 提取包名: `dsp_bad_code`
3. ✓ 检测命名风格: 发现下划线 `_`
4. ✓ 匹配规则: INTERNAL-GO-1.4 "包名：不应该包含下划线"
5. ✓ 检查规范章节: 01. style - 风格规范 / 1.4. label - 标识符命名
6. ✓ 检查规范标记: 【建议】
7. ✓ 应用分级策略: 代码风格类 → Warning

**置信度**: High

**时间戳**: 2025-12-22 10:23:42

---

### D042 - 依赖未声明（深度分析）

**问题位置**: src/service/order.go:156

**检测来源**: 深度分析 (dependency_analysis)

**检测过程**:
1. ✓ 分析函数 `ProcessOrder` 的依赖关系
2. ✓ 提取函数体中的包调用: `database.Query()`
3. ✓ 检查文件顶部的 import 声明
4. ✓ 发现: `database` 包未在 import 中声明
5. ✓ 匹配规范: INTERNAL-GO-5.3 "使用绝对导入，禁止使用相对导入"
6. ✓ 应用分级策略: 功能异常类 → Error

**置信度**: Medium (需人工确认是否为隐式导入)

**时间戳**: 2025-12-22 10:24:12

---

## 追溯统计

- **总问题数**: 28个
- **检测来源分布**:
  - 内部规范: 15个
  - Lint工具: 8个
  - 深度分析: 5个
- **平均置信度**: High (92%)
- **检测耗时**: 12.3秒
```

### 4. 规则命中说明强化

每个问题都必须清晰标注：
- ✅ 规则来源（内部规范/外部规范/Lint工具/深度分析）
- ✅ 具体规则ID和内容
- ✅ 规范章节路径
- ✅ 规范标记（【强制】/【建议】）
- ✅ 分级依据（为什么是这个级别）

**完整示例**：

```markdown
#### [E006] 未处理错误返回值

**级别**: Error
**分类**: 错误处理
**位置**: src/app/main.go:39

**当前代码**:
```go
inMsg, _ := ioutil.ReadAll(r.Body)
```

**问题描述**: `ioutil.ReadAll` 的错误被忽略，可能导致后续处理异常数据。

**规范依据**:
- **来源**: Lint工具 (golangci-lint/errcheck) + 内部规范
- **Lint规则**: errcheck - 检查未处理的错误返回值
- **内部规范**: INTERNAL-GO-3.3
- **规范章节**: 03. exception - 错误处理 / 3.3. dealwith - 错误处理
- **规范内容**: "显式处理error，或使用空白标识符忽略。【强制】"
- **规范标记**: 【强制】
- **分级依据**: 错误处理类问题，可能导致功能异常 → Error级别

**建议修改**:
```go
inMsg, err := ioutil.ReadAll(r.Body)
if err != nil {
    http.Error(w, fmt.Sprintf("failed to read request body: %v", err), http.StatusBadRequest)
    return
}
```

**参考文档**: [Effective Go - Errors](https://go.dev/doc/effective_go#errors)
```

## 使用示例

### Checker Skill 生成数据

```markdown
检查完成后，将结果组织成标准数据格式，然后调用 report-generator：

1. 收集所有问题
2. 按文件组织
3. 按级别排序（Error → Warning → Suggestion）
4. 合并同行多问题
5. 记录每个问题的 detection_trace
6. 计算统计信息
7. 准备建议
8. 调用 report-generator 并传递 JSON 数据
```

### Report Generator 接收数据

```markdown
接收 JSON 格式的检查结果：
1. 验证数据格式
2. 按排序规则组织问题
3. 合并同行问题
4. 生成问题详情章节
5. 生成检测过程追溯章节
6. 输出标准化的 Markdown 报告
```

---

**最后更新**: 2025-12-26
