# 语言检测排除规则

本文档定义了在项目语言检测时需要排除的目录和文件，确保语言检测的准确性。

## 为什么需要排除规则？

语言检测的目的是识别项目的**主要开发语言**，因此需要排除：
- 第三方依赖包（可能包含不同语言的代码）
- 工具和构建输出（可能包含生成的代码）
- 测试和示例代码（可能使用多种语言）

## 核心排除规则

### 1. 依赖目录（最高优先级）⭐⭐⭐

这些目录包含第三方依赖，会严重干扰语言检测：

```
# Node.js / JavaScript
node_modules/
bower_components/
jspm_packages/

# Go
vendor/
Godeps/

# Python
venv/
.venv/
env/
.env/
virtualenv/
site-packages/
*.egg-info/

# Java
target/
lib/
libs/
.gradle/
build/

# Ruby
vendor/bundle/
```

### 2. 构建和输出目录

```
dist/
build/
out/
bin/
obj/
*.build/
```

### 3. 版本控制和工具配置

```
.git/
.svn/
.hg/
.claude/
.vscode/
.idea/
```

### 4. 测试和示例目录

```
test/
tests/
testing/
__tests__/
spec/
specs/
examples/
sample/
samples/
demo/
demos/
```

### 5. 文档目录

```
docs/
doc/
documentation/
```

## 检测逻辑

```python
def should_exclude_for_language_detection(path):
    """
    判断路径是否应该在语言检测时排除

    Args:
        path: 文件或目录路径

    Returns:
        bool: True 表示应该排除
    """
    exclude_patterns = [
        # 依赖目录
        "node_modules/", "vendor/", "venv/", ".venv/",
        "virtualenv/", "target/", "lib/", "libs/",

        # 构建输出
        "dist/", "build/", "out/", "bin/", "obj/",

        # 工具配置
        ".git/", ".svn/", ".claude/", ".vscode/", ".idea/",

        # 测试和示例
        "test/", "tests/", "__tests__/", "spec/", "specs/",
        "examples/", "sample/", "demo/",

        # 文档
        "docs/", "doc/", "documentation/"
    ]

    for pattern in exclude_patterns:
        if pattern in path:
            return True

    return False
```

## 使用示例

### 示例 1: 检测 Node.js 项目

```bash
# 项目结构
project/
├── src/
│   ├── index.js           # ✅ 检测
│   └── utils.js           # ✅ 检测
├── node_modules/          # ❌ 排除（依赖）
│   └── express/
├── test/                  # ❌ 排除（测试）
│   └── index.test.js
└── package.json

# 检测结果: JavaScript (基于 src/ 目录)
```

### 示例 2: 检测 Go 项目

```bash
# 项目结构
project/
├── cmd/
│   └── main.go            # ✅ 检测
├── internal/
│   └── service.go         # ✅ 检测
├── vendor/                # ❌ 排除（依赖）
│   └── github.com/
├── test/                  # ❌ 排除（测试）
│   └── main_test.go
└── go.mod

# 检测结果: Go (基于 cmd/ 和 internal/)
```

### 示例 3: 检测 Python 项目

```bash
# 项目结构
project/
├── src/
│   ├── app.py             # ✅ 检测
│   └── utils.py           # ✅ 检测
├── venv/                  # ❌ 排除（虚拟环境）
│   └── lib/
├── tests/                 # ❌ 排除（测试）
│   └── test_app.py
└── requirements.txt

# 检测结果: Python (基于 src/)
```

## 与代码检查排除规则的区别

| 用途 | 语言检测排除规则 | 代码检查排除规则 |
|------|---------------|---------------|
| **文件位置** | `language-detector/rules/LANGUAGE-DETECTION-EXCLUDE.md` | `report-generator/rules/EXCLUDE-RULES.md` |
| **目的** | 确保准确识别主要开发语言 | 确保只检查需要规范的源代码 |
| **范围** | **精简**（只排除影响语言检测的目录） | **完整**（排除所有不需要检查的文件） |
| **示例** | 排除 `node_modules/`, `vendor/` | 排除 `.env`, `.gitignore`, `*.md`, 配置文件等 |

**关键区别**：
- 语言检测排除规则：**轻量级**，只关注依赖和构建目录
- 代码检查排除规则：**全面**，包含配置、文档、测试、工具文件等

## 执行流程

```
language-detector 执行时:
  ↓
应用 LANGUAGE-DETECTION-EXCLUDE.md 规则
  ↓
扫描项目目录
  ↓
统计代码文件类型 (.go, .py, .java, .js等)
  ↓
返回主要语言
```

---

**最后更新**: 2025-12-26
