---
name: tool-runner
description: 调用外部Lint工具进行代码规范检查,输出统一格式到临时文件,支持自动安装缺失工具
version: 3.0
---

# Tool Runner Skill

## 功能概述

Tool Runner 负责根据项目语言类型自动选择和调用相应的外部 Lint 工具,并将检查结果输出到临时文件供后续步骤使用。

**核心职责**:
- 自动检测并安装缺失的 Lint 工具
- 根据检测模式（快速/深度）配置工具参数
- 调用外部 Lint 工具执行代码检查
- 解析工具输出并统一格式
- **将结果写入临时文件**: `.claude/temp/lint-results-{timestamp}.json`

---

## 输出规范 ⭐⭐⭐

### 输出文件

**必须**将检查结果写入临时文件,供 `report-generator` 读取:

```
文件路径: .claude/temp/lint-results-{timestamp}.json
文件格式: JSON (符合 REPORT-DATA-FORMAT.md 规范)
文件用途: 供 report-generator 读取并生成最终报告
```

### 输出格式 (统一格式)

**⭐ 关键要求**: tool-runner 的输出格式必须与 code-checker 保持一致,遵循 `REPORT-DATA-FORMAT.md` 规范。

```json
{
  "version": "2.0",
  "source": "external_tool",
  "tool": "golangci-lint",
  "language": "go",
  "mode": "fast",
  "timestamp": "20250115-143022",
  "metadata": {
    "language": "go",
    "mode": "incremental",
    "check_time": "2025-01-15 14:30:22",
    "total_files": 3,
    "total_issues": 10
  },
  "statistics": {
    "errors": 3,
    "warnings": 5,
    "suggestions": 2,
    "by_category": {
      "error_handling": {"errors": 2, "warnings": 1, "suggestions": 0},
      "naming": {"errors": 1, "warnings": 2, "suggestions": 1},
      "code_format": {"errors": 0, "warnings": 2, "suggestions": 1}
    }
  },
  "files": [
    {
      "path": "src/app/main.go",
      "issues": [
        {
          "id": "LINT-E001",
          "level": "error",
          "category": "error_handling",
          "title": "错误返回值未检查",
          "location": "src/app/main.go:42",
          "line_number": 42,
          "description": "函数 `ioutil.ReadAll` 返回的 error 未被检查,可能导致运行时异常。",
          "current_code": "data, _ := ioutil.ReadAll(r.Body)",
          "suggested_code": "data, err := ioutil.ReadAll(r.Body)\nif err != nil {\n    return fmt.Errorf(\"failed to read body: %w\", err)\n}",
          "explanation": "忽略错误返回值是危险的,应该始终检查并处理错误。",
          "reference": "golangci-lint - errcheck",
          "reference_url": "https://github.com/kisielk/errcheck",
          "detection_trace": {
            "tool": "golangci-lint",
            "rule": "errcheck",
            "matched_pattern": "error return value not checked",
            "detection_method": "static_analysis"
          }
        }
      ]
    }
  ],
  "tool_warnings": [
    "工具安装警告信息（如果有）"
  ]
}
```

**关键字段说明**:

| 字段 | 类型 | 必需 | 说明 |
|------|------|------|------|
| `version` | string | ✅ | 数据格式版本 (2.0) |
| `source` | string | ✅ | 数据来源标识 (`external_tool`) |
| `tool` | string | ✅ | 使用的工具名称 (golangci-lint/pylint/checkstyle) |
| `language` | string | ✅ | 检测的语言类型 (go/python/java) |
| `mode` | string | ✅ | 检测模式 (fast/deep) |
| `timestamp` | string | ✅ | 执行时间戳 |
| `metadata` | object | ✅ | 元数据信息 |
| `statistics` | object | ✅ | 统计信息 |
| `files` | array | ✅ | 文件列表及问题 |
| `tool_warnings` | array | ⚠️ | 工具安装或执行警告（可选） |

**每个问题 (issue) 必须包含**:

| 字段 | 类型 | 必需 | 说明 |
|------|------|------|------|
| `id` | string | ✅ | 问题ID (LINT-E001, LINT-W001) |
| `level` | string | ✅ | 级别 (error/warning/suggestion) |
| `category` | string | ✅ | 分类 (error_handling/naming/code_format) |
| `title` | string | ✅ | 简短标题 |
| `location` | string | ✅ | 位置 (完整路径:行号) |
| `line_number` | number | ✅ | 行号 |
| `description` | string | ✅ | 详细描述 |
| `current_code` | string | ⚠️ | 当前代码（推荐提供） |
| `suggested_code` | string | ⚠️ | 建议代码（推荐提供） |
| `explanation` | string | ⚠️ | 额外说明（推荐提供） |
| `reference` | string | ✅ | **规范来源/工具规则名称** ⭐⭐⭐ |
| `reference_url` | string | ✅ | **规范文档链接** ⭐⭐⭐ |
| `detection_trace` | object | ✅ | 检测过程追溯 |

### 规范来源映射 ⭐⭐⭐

**关键要求**: 每个 Lint 工具检测出的问题都必须提供规范来源和文档链接。

#### Go 语言工具规范映射

**golangci-lint** 包含多个 linter,需要映射到具体的规范来源:

| linter 规则 | reference | reference_url |
|------------|-----------|---------------|
| errcheck | golangci-lint - errcheck | https://github.com/kisielk/errcheck |
| govet | golangci-lint - govet | https://pkg.go.dev/cmd/vet |
| staticcheck | golangci-lint - staticcheck | https://staticcheck.io/ |
| gosimple | golangci-lint - gosimple | https://github.com/dominikh/go-tools/tree/master/simple |
| ineffassign | golangci-lint - ineffassign | https://github.com/gordonklaus/ineffassign |
| unused | golangci-lint - unused | https://github.com/dominikh/go-tools/tree/master/unused |
| golint | golangci-lint - golint | https://github.com/golang/lint |
| gofmt | golangci-lint - gofmt | https://pkg.go.dev/cmd/gofmt |

#### Python 语言工具规范映射

**pylint** 规则映射:

| pylint 代码 | reference | reference_url |
|------------|-----------|---------------|
| C0103 | pylint - C0103 (Invalid name) | https://pylint.pycqa.org/en/latest/messages/convention/c0103.html |
| C0114 | pylint - C0114 (Missing module docstring) | https://pylint.pycqa.org/en/latest/messages/convention/c0114.html |
| E0401 | pylint - E0401 (Import error) | https://pylint.pycqa.org/en/latest/messages/error/e0401.html |
| W0612 | pylint - W0612 (Unused variable) | https://pylint.pycqa.org/en/latest/messages/warning/w0612.html |

**flake8** 规则映射:

| flake8 代码 | reference | reference_url |
|------------|-----------|---------------|
| E501 | flake8 - E501 (Line too long) | https://www.flake8rules.com/rules/E501.html |
| F401 | flake8 - F401 (Imported but unused) | https://www.flake8rules.com/rules/F401.html |
| W503 | flake8 - W503 (Line break before operator) | https://www.flake8rules.com/rules/W503.html |

**black** 规则映射:

| black 检测 | reference | reference_url |
|-----------|-----------|---------------|
| 格式不符 | black - Code formatting | https://black.readthedocs.io/en/stable/ |

#### Java 语言工具规范映射

**checkstyle** 规则映射:

| checkstyle 规则 | reference | reference_url |
|----------------|-----------|---------------|
| LineLength | checkstyle - LineLength | https://checkstyle.sourceforge.io/config_sizes.html#LineLength |
| MemberName | checkstyle - MemberName | https://checkstyle.sourceforge.io/config_naming.html#MemberName |
| JavadocMethod | checkstyle - JavadocMethod | https://checkstyle.sourceforge.io/config_javadoc.html#JavadocMethod |

**spotbugs** 规则映射:

| spotbugs 代码 | reference | reference_url |
|--------------|-----------|---------------|
| NP_NULL_ON_SOME_PATH | spotbugs - NP_NULL_ON_SOME_PATH | https://spotbugs.readthedocs.io/en/stable/bugDescriptions.html#np-null-on-some-path |
| DLS_DEAD_LOCAL_STORE | spotbugs - DLS_DEAD_LOCAL_STORE | https://spotbugs.readthedocs.io/en/stable/bugDescriptions.html#dls-dead-local-store |

### 工具输出转换规则

**转换步骤**:

1. **解析工具原始输出** (JSON/XML/文本)
2. **提取问题信息**
   - 文件路径、行号、列号
   - 规则代码 (如 errcheck, E501, NP_NULL_ON_SOME_PATH)
   - 问题描述
3. **映射到标准格式**
   - `level`: 根据工具的严重程度映射 (error/warning/info → error/warning/suggestion)
   - `category`: 根据规则类型映射 (见下文分类映射表)
   - `title`: 提取或生成简短标题
   - `reference`: 使用规范映射表查找
   - `reference_url`: 使用规范映射表查找
4. **生成ID**: `LINT-{level_prefix}{序号}` (如 LINT-E001, LINT-W001)
5. **添加检测追溯信息**

**问题分类映射**:

| 工具规则类型 | category | 说明 |
|------------|----------|------|
| errcheck, error handling | error_handling | 错误处理 |
| naming conventions | naming | 命名规范 |
| code style, formatting | code_format | 代码格式 |
| documentation, comments | documentation | 文档注释 |
| unused variables/imports | code_quality | 代码质量 |
| concurrency, race conditions | concurrency | 并发安全 |

**级别映射**:

| 工具级别 | 标准级别 | 说明 |
|---------|---------|------|
| error | error | 错误 |
| warning | warning | 警告 |
| info / convention | suggestion | 建议 |

---

## 支持的工具

| 语言 | 工具 | 输出格式 | 说明 |
|------|------|---------|------|
| **Go** | golangci-lint | JSON | 综合性 Lint 工具,集成多个 linter |
| **Python** | pylint | JSON | 静态分析工具,检查规范和错误 |
| **Python** | flake8 | 文本 | 轻量级风格检查工具 |
| **Python** | black | 文本 | 代码格式化检查（--check 模式） |
| **Java** | checkstyle | XML | 基于 Google Style Guide 的风格检查 |
| **Java** | spotbugs | XML | 静态代码分析,检查潜在 bug |

---

## 核心功能

### 1. 自动安装工具 ⭐

**检测并自动安装缺失的工具**,无需用户手动安装。

**安装机制**:
1. 执行前检测工具是否已安装（`which`/`command`）
2. 调用平台对应的安装脚本（`scripts/install-tools.sh` 或 `install-tools.bat`）
3. 使用深信服内部镜像源加速安装（PyPI 和 Go Proxy）
4. 验证安装结果,失败则记录警告并跳过该工具

**安装脚本位置**:
```
skills/tool-runner/scripts/
├── install-tools.sh      # Linux/macOS
└── install-tools.bat     # Windows
```

**镜像源配置**:
- Python: `http://mirrors.sangfor.org/pypi/simple`
- Go: `http://mirrors.sangfor.org/nexus/repository/go-proxy-group`

**优雅降级**:
- 安装失败时跳过该工具,继续执行其他工具
- 在输出 JSON 的 `warnings` 字段中记录失败信息
- 最大化检查覆盖率

### 2. 工具参数配置

根据检测模式（快速/深度）配置不同的工具参数:

**快速模式**:
- 只检查变更的文件
- 使用轻量级配置（`--fast`）
- 快速返回结果

**深度模式**:
- 检查变更文件 + 关联依赖文件
- 使用完整配置（`--enable-all`）
- 更深入的分析

### 3. 结果解析和格式转换 ⭐⭐⭐

**关键职责**: 将各工具的原始输出转换为统一的标准格式。

**转换流程**:

1. **解析原始输出**
   - JSON 格式 (golangci-lint, pylint) → 直接解析
   - XML 格式 (checkstyle, spotbugs) → 解析 XML 并转换为 JSON
   - 文本格式 (flake8, black) → 正则提取信息

2. **规范映射** ⭐⭐⭐
   - 根据工具规则代码查找规范映射表
   - 获取 `reference` (规范名称)
   - 获取 `reference_url` (文档链接)
   - 如果映射表中没有,使用默认值: `{tool} - {rule_code}`

3. **字段转换**
   - `severity` → `level` (error/warning/info → error/warning/suggestion)
   - `rule` → `category` (使用分类映射表)
   - 生成 `id`: LINT-{E/W/S}{序号}
   - 生成 `title`: 基于消息提取简短标题
   - 生成 `location`: {file}:{line}

4. **补充信息**
   - `current_code`: 读取文件的指定行
   - `suggested_code`: 基于工具建议生成(如果有)
   - `explanation`: 提取或生成额外说明
   - `detection_trace`: 记录检测工具和规则

5. **写入临时文件**
   - 按照 REPORT-DATA-FORMAT.md 规范组织数据
   - 写入 `.claude/temp/lint-results-{timestamp}.json`

**示例转换**:

**原始输出 (golangci-lint)**:
```json
{
  "Issues": [{
    "FromLinter": "errcheck",
    "Text": "Error return value of `ioutil.ReadAll` is not checked",
    "Pos": {
      "Filename": "src/app/main.go",
      "Line": 42
    },
    "Severity": "error"
  }]
}
```

**转换后 (标准格式)**:
```json
{
  "id": "LINT-E001",
  "level": "error",
  "category": "error_handling",
  "title": "错误返回值未检查",
  "location": "src/app/main.go:42",
  "line_number": 42,
  "description": "函数 `ioutil.ReadAll` 返回的 error 未被检查,可能导致运行时异常。",
  "current_code": "data, _ := ioutil.ReadAll(r.Body)",
  "reference": "golangci-lint - errcheck",
  "reference_url": "https://github.com/kisielk/errcheck",
  "detection_trace": {
    "tool": "golangci-lint",
    "rule": "errcheck",
    "matched_pattern": "Error return value not checked",
    "detection_method": "static_analysis"
  }
}
```

### 4. 错误处理

- **工具未安装**: 自动尝试安装,失败则跳过并记录警告
- **工具执行失败**: 记录错误详情,继续执行其他工具
- **配置文件缺失**: 使用默认配置
- **文件不存在**: 跳过该文件,记录警告
- **执行超时**: 终止命令,记录超时日志,继续执行其他工具

---

## 执行流程

```
读取上下文：.claude/temp/lint-context-{timestamp}.json
  ├─ 获取：language, mode, scope, files (含 check_lines)
  └─ 获取：timestamp（用于输出文件命名）
  ↓
选择工具（golangci-lint / pylint / checkstyle 等）
  ↓
检测并安装缺失工具
  ↓
构建命令参数（根据模式和检查范围）
  ↓
执行工具（并行执行多个工具）
  ↓
解析输出（统一格式）
  ↓
写入临时文件: .claude/temp/lint-results-{timestamp}.json ⭐
```

---

## 配置文件

工具配置文件位于:

```
skills/tool-runner/configs/
├── go/
│   └── .golangci.yml
├── python/
│   ├── .pylintrc
│   ├── .flake8
│   └── pyproject.toml
└── java/
    ├── checkstyle.xml
    └── spotbugs.xml
```

**优先级**: 项目根目录配置 > skill 配置 > 工具默认配置

---

## 与其他 Skills 的协作

```
language-detector → 检测语言类型
                         ↓
                   tool-runner
                         ↓
        写入: .claude/temp/lint-results-{timestamp}.json
                         ↓
                  report-generator → 读取临时文件生成报告
```

---

## 超时配置 ⭐⭐⭐

**所有 Bash 命令调用必须设置 timeout 参数,防止工具卡死**

### 推荐超时阈值

| 操作类型 | 推荐超时 | 说明 |
|---------|---------|------|
| 工具检测 | 5秒 | which/command 检查 |
| 工具安装 | 300秒 | pip install/go install |
| Go Lint | 600秒 | golangci-lint（工具内部 10 分钟） |
| Python Lint | 300秒 | pylint/flake8 运行 |
| Java Lint | 300秒 | checkstyle/spotbugs |

### 工具特定配置

**golangci-lint**:
```bash
# ✅ 推荐配置
golangci-lint run \
    --timeout=10m \
    --max-issues-per-linter=100 \
    --max-same-issues=5 \
    --fast \
    --out-format=json
```

**pylint**:
```bash
# ✅ 推荐配置（显式指定文件,避免递归扫描整个文件系统）
find src/ -name "*.py" -type f | head -n 200 | xargs pylint \
    --output-format=json \
    --max-line-length=120
```

### Bash 调用示例

```python
# ✅ 正确：设置超时
Bash(
    command="golangci-lint run --timeout=10m src/",
    timeout=660000,  # 11分钟（工具10分钟+1分钟buffer）
    description="执行Go代码Lint检查"
)

# ✅ 正确：工具检测
Bash(
    command="which golangci-lint",
    timeout=5000,  # 5秒
    description="检查工具是否安装"
)
```

---

## 限制和注意事项

1. **工具依赖**: 支持自动安装,但安装失败仍需用户手动安装
2. **超时要求**: 所有工具调用必须设置合理的 timeout 参数
3. **网络要求**: 自动安装依赖网络连接访问深信服镜像源
4. **权限要求**: 安装工具可能需要系统权限（pip3、go install）
5. **平台限制**: Java 工具在某些平台可能需要手动安装
6. **输出格式**: 部分工具可能不支持 JSON 输出,需要文本解析
7. **临时文件**: 必须将结果写入 `.claude/temp/` 目录,由 `cleanup-handler` 清理
