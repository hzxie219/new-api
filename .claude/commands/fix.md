---
allowed-tools: Bash,Read,Write,Edit,Glob,Grep
description: 代码自动修复命令。在快速模式下自动触发,在深度模式下需手动执行。
skills: report-reader,issue-analyzer,code-fixer-go,code-fixer-python,code-fixer-java,fix-report-generator,cleanup-handler
version: 2.0
---

# 代码自动修复命令 (/fix)

## 功能概述

本命令基于 /lint 命令生成的检查报告，自动修复 error 和 warning 级别的代码问题。

**两种触发方式：**

- **快速模式（Fast）**: /lint 命令自动触发 /fix，无需人工干预
- **深度模式（Deep）**: 需要用户查看深度分析报告后手动执行 /fix

## 使用场景

### 快速模式 - 自动修复

用户执行 `/lint [branch]` 后，系统自动：
1. 执行代码检查（Lint工具 + AI分析 + 安全检查）
2. 生成检查报告并验证
3. **自动触发 /fix 修复所有问题**
4. 生成修复报告

用户无需手动执行 /fix 命令。

### 深度模式 - 手动修复

用户执行 `/lint --mode=deep [branch]` 后：
1. 执行代码检查 + 深度分析（依赖、调用链、架构）
2. 生成深度分析报告
3. **系统提示用户查看报告**
4. **用户确认后执行 /fix（自动使用最新深度报告）**

---

## ⚠️ 强约束规则

**所有代码修复必须100%遵守以下强约束规则：**

### 1. 只修复报告中的问题 ⭐⭐⭐

- ✅ **必须**: 只修复 lint 报告中明确列出的问题
- ✅ **必须**: 每个修复都关联到报告中的问题 ID（如 E001, W023）
- ❌ **禁止**: 修复报告之外发现的任何问题
- ❌ **禁止**: 进行"顺手优化"或"代码美化"
- ❌ **禁止**: 添加额外的"改进"或"重构"
- ❌ **禁止**: 修改与问题无关的代码行

**示例**：
```python
# ✅ 正确：只修复报告中的问题
报告：[E001] 第42行未处理错误返回值
修复：只修改第42行，添加错误处理

# ❌ 错误：顺手优化
报告：[E001] 第42行未处理错误返回值
修复：修改第42行 + 重命名变量 + 调整缩进 + 添加注释
```

### 2. 覆盖所有可修复问题 ⭐⭐⭐

- ✅ **必须**: 尝试修复报告中所有 error 和 warning 级别问题
- ✅ **必须**: 如果无法自动修复，必须标记为"跳过"并说明原因
- ❌ **禁止**: 遗漏报告中的任何 error 或 warning 问题
- ❌ **禁止**: 因为困难而跳过可修复的问题
- ✅ **必须**: 在修复报告中记录每个问题的处理结果（成功/失败/跳过）

**合理的跳过原因**：
- 高风险修复，需要整个函数重构
- 涉及业务逻辑，需要人工判断
- 修复后语法验证失败
- 依赖外部资源或配置

**不合理的跳过原因**：
- "这个问题比较难"
- "不确定怎么修"
- "时间不够"

### 3. 报告唯一性约束 ⭐⭐⭐

**⚠️ 适用所有模式：快速模式和深度模式都必须遵守**

**快速模式修复报告**：
- ✅ **必须**: 生成新报告前删除旧的快速模式修复报告（**无论哪种范围**）
- ✅ **必须**: 删除模式：`fix-{scope}-fast-{language}-*.md`
- ❌ **禁止**: `doc/fix/` 目录中存在多份同语言、同范围的快速模式修复报告

**深度模式修复报告**：
- ✅ **必须**: 生成新报告前删除旧的深度模式修复报告（**无论哪种范围**）
- ✅ **必须**: 删除模式：`fix-{scope}-deep-{language}-*.md`
- ❌ **禁止**: `doc/fix/` 目录中存在多份同语言、同范围的深度模式修复报告

**强制要求**：
- ⚠️ **关键**：这是强制要求，不区分快速模式还是深度模式
- ⚠️ **目标**：确保报告目录中只保留当前执行生成的最新修复报告
- ⚠️ **执行时机**：在调用 fix-report-generator 生成报告之前就必须删除

**示例**：
```bash
# 增量+快速模式：删除 fix-incremental-fast-go-*.md
# 全量+深度模式：删除 fix-full-deep-python-*.md
# 最新+快速模式：删除 fix-latest-fast-java-*.md
```

### 4. 修复范围最小化 ⭐⭐⭐

- ✅ **必须**: 只修改问题所在的代码行或最小必要范围
- ❌ **禁止**: 重构整个函数或模块
- ❌ **禁止**: 修改相关但不在报告中的代码
- ❌ **禁止**: 调整与问题无关的格式、缩进、空行
- ✅ **必须**: 保持修复最小化原则

**示例**：
```go
// 报告：[E006] 第39行未处理错误返回值
// ✅ 正确的修复范围
// 修改前：
inMsg, _ := ioutil.ReadAll(r.Body)

// 修改后：
inMsg, err := ioutil.ReadAll(r.Body)
if err != nil {
    return err
}

// ❌ 错误的修复范围（过度修复）
// 修改前：
func HandleRequest(w http.ResponseWriter, r *http.Request) {
    inMsg, _ := ioutil.ReadAll(r.Body)
    // ... 其他代码
}

// 修改后：重构整个函数，添加日志，优化变量命名
func HandleRequest(writer http.ResponseWriter, request *http.Request) {
    logger.Info("Handling request")
    inputMessage, err := ioutil.ReadAll(request.Body)
    if err != nil {
        logger.Error("Failed to read body", err)
        return err
    }
    // ... 其他代码
}
```

### 5. 数据传递约束 ⭐⭐⭐ 新增

**优先使用JSON数据**（从/lint自动触发时）:
- ✅ **必须**: 优先读取 `--data-file` 参数指定的JSON数据
- ✅ **必须**: JSON数据来源：`.claude/temp/report-data-merged-{timestamp}.json`
- ✅ **优势**: 避免Markdown解析歧义，数据更准确
- ✅ **必须**: JSON数据包含完整的问题信息（已验证、已修正、已合并）

**传统Markdown解析**（手动执行时）:
- ✅ **允许**: 如果未提供 `--data-file`，解析Markdown报告
- ✅ **必须**: 使用 report-reader 解析报告
- ✅ **必须**: 验证解析结果的完整性

**数据优先级**:
```
1. --data-file 参数指定的JSON文件（最高优先级）
2. --report 参数指定的Markdown报告
3. 自动查找最新Markdown报告
```

### 6. 自动报告检测约束 ⭐⭐⭐

**报告查找优先级**：
- ✅ **必须**: 如果未指定 `--report`，自动查找最新报告
- ✅ **必须**: 优先级1：最新的深度模式报告（`lint-*-deep-*.md`）
- ✅ **必须**: 优先级2：最新的快速模式报告（`lint-*-fast-*.md`）
- ✅ **必须**: 如果都不存在，提示用户先执行 `/lint`
- ❌ **禁止**: 在找不到报告时静默失败

**错误处理**：
```bash
# 场景1：找到深度模式报告
使用报告：lint-incremental-deep-go-20251222-143000.md

# 场景2：只找到快速模式报告
使用报告：lint-incremental-fast-go-20251222-120000.md

# 场景3：未找到任何报告
❌ 错误：未找到 lint 报告
建议：先执行 /lint 命令生成报告
```

### 6. 修复验证约束 ⭐⭐⭐

**语法验证**：
- ✅ **必须**: 每次修复后验证代码语法正确性
- ✅ **必须**: Go：尝试 `go build` 或 `gofmt -e`
- ✅ **必须**: Python：使用 `ast.parse()` 验证语法
- ✅ **必须**: Java：尝试 `javac -Xstderr`
- ❌ **禁止**: 生成语法错误的代码

**修复回滚**：
- ✅ **必须**: 如果语法验证失败，立即回滚该修复
- ✅ **必须**: 记录失败原因到修复报告
- ✅ **必须**: 标记该问题为"修复失败"，建议人工处理

### 7. 备份创建约束 ⭐⭐⭐

**Git 仓库环境**：
- ✅ **必须**: 修复前创建 git stash 备份
- ✅ **必须**: 使用描述性消息：`git stash push -m "backup before fix - {timestamp}"`
- ✅ **必须**: 在修复报告中记录备份ID

**非 Git 环境**：
- ✅ **必须**: 复制文件到 `.backup/{timestamp}/` 目录
- ✅ **必须**: 在修复报告中记录备份路径
- ❌ **禁止**: 跳过备份直接修复

### 8. 修复级别约束 ⭐⭐⭐

**级别修复规则**：
- ✅ **必须**: `--level=error` 只修复 Error 级别
- ✅ **必须**: `--level=warning` 修复 Warning **和** Error 级别（层级修复）
- ✅ **必须**: `--level=all` 修复所有级别（Error + Warning + Suggestion）
- ❌ **禁止**: `--level=warning` 只修复 Warning 而忽略 Error

**默认行为**：
- ✅ **必须**: 快速模式默认修复 Error + Warning 级别（不含 Suggestion）
- ✅ **必须**: 深度模式默认修复所有级别（Error + Warning + Suggestion）
- ✅ **允许**: 用户通过 `--level` 参数覆盖默认行为

### 9. 临时文件清理约束 ⭐⭐⭐

**必须清理的文件**：
- ✅ **必须**: 清理 `.claude/temp/` 目录
- ✅ **必须**: 清理 `.bak` 备份文件
- ✅ **必须**: 清理中间生成的临时数据
- ❌ **禁止**: 清理用户的源代码文件
- ❌ **禁止**: 清理 `doc/lint/` 和 `doc/fix/` 报告

**保留的文件**：
- ✅ **必须**: 保留修复后的源代码
- ✅ **必须**: 保留 fix 报告：`doc/fix/fix-*.md`
- ✅ **必须**: 保留 lint 报告：`doc/lint/lint-*.md`
- ✅ **必须**: 保留 git stash 备份或 `.backup/` 目录

### 10. 快速模式自动触发约束 ⭐⭐⭐

**快速模式行为**：
- ✅ **必须**: 由 `/lint` 命令自动触发
- ✅ **必须**: 无需用户手动执行
- ✅ **必须**: 自动修复所有 error 和 warning 问题
- ❌ **禁止**: 在快速模式下等待用户确认

**深度模式行为**：
- ❌ **禁止**: 自动触发修复
- ✅ **必须**: 等待用户查看报告后手动执行
- ✅ **必须**: 支持选择性修复（使用 `--level` 和 `--files` 参数）

---

## 基本语法

```bash
# 快速模式下 - 通常不需要手动执行
# /lint 会自动触发修复

# 深度模式下 - 手动执行修复
/fix                                  # 自动使用最新深度报告
/fix --report lint-deep-{language}-{date}.md  # 或明确指定报告

# 通用用法
/fix                                    # 使用最新报告
/fix --report {report-filename}          # 指定报告
/fix --level {error|warning|all}        # 只修复特定级别
/fix --files {file1,file2,...}          # 只修复特定文件
```

### 参数说明

- **--report**（可选）: 指定要使用的lint报告文件名，默认使用最新报告

- **--level**（可选）: 指定修复级别，支持层级修复
  - `error`: 只修复 Error 级别问题
  - `warning`: 修复 Warning **及以上**级别问题（Warning + Error）
  - `all` 或省略: 修复所有级别问题（Error + Warning + Suggestion）
  - **默认值**: `all`

- **--files**（可选）: 只修复指定的文件（逗号分隔）

### 问题级别说明

本系统使用三级分类体系，详见 [ISSUE-GRADING-STRATEGY.md](rules/ISSUE-GRADING-STRATEGY.md)：

| 级别 | 说明 | 典型问题 |
|------|------|---------|
| **Error** | 功能异常、稳定性、性能、安全问题 | 未处理错误、空指针、SQL注入、资源泄漏 |
| **Warning** | 代码风格、命名规范问题 | 驼峰命名、缩进格式、导入顺序 |
| **Suggestion** | 注释文档问题 | 缺少注释、TODO未处理 |

### 修复级别层级关系

```
--level=error    →  修复: [Error]
--level=warning  →  修复: [Error, Warning]  (⚠️ 注意：包含Error)
--level=all      →  修复: [Error, Warning, Suggestion]
```

**重要**：`--level=warning` 会修复 Warning **和** Error 两个级别，不是只修复 Warning。

## 使用示例

### 快速模式工作流

```bash
# 1. 执行检查（自动修复）
/lint develop

# 系统自动完成：
# ├─ 代码检查（Lint + AI + 安全）
# ├─ 生成检查报告
# ├─ 自动修复所有问题
# └─ 生成修复报告

# 2. 查看结果
# - doc/lint/lint-fast-python-20251222.md (检查报告)
# - doc/fix/fix-fast-python-20251222.md (修复报告)
```

### 深度模式工作流

```bash
# 1. 执行深度检查
/lint --mode=deep develop

# 系统完成：
# ├─ 代码检查（Lint + AI + 安全）
# ├─ 深度分析（依赖 + 调用链 + 架构）
# └─ 生成深度报告

# 2. 查看深度报告（可选）
cat doc/lint/lint-deep-go-20251222.md

# 3. 确认后执行修复（自动检测最新深度报告）
/fix

# 或明确指定报告
/fix --report lint-deep-go-20251222.md

# 4. 查看修复报告
cat doc/fix/fix-deep-go-20251222.md
```

### 场景示例

#### 场景1：快速模式日常开发

```bash
# 提交前快速检查和修复
/lint main

# 无需其他操作，系统自动：
# ✅ 检查代码
# ✅ 修复问题
# ✅ 生成报告
```

#### 场景2：深度模式重构确认

```bash
# 重大重构后的深度检查
/lint --mode=deep main

# 查看深度报告（包含依赖分析、调用链、架构检查）
cat doc/lint/lint-deep-python-20251222.md

# 确认修复计划
/fix  # 自动使用最新深度报告
```

#### 场景3：只修复 Error 级别（关键问题优先）

```bash
# 深度模式下只修复关键问题
/lint --mode=deep develop

# 只修复 Error 级别（功能异常、安全漏洞等）
/fix --report lint-deep-java-20251222.md --level error

# 结果：
# ✅ 修复 12 个 Error 级别问题
# ⏭️ 跳过 25 个 Warning 级别问题
# ⏭️ 跳过 8 个 Suggestion 级别问题
```

#### 场景4：修复 Warning 及以上级别

```bash
# 深度模式下修复所有非建议类问题
/lint --mode=deep main

# 修复 Warning + Error 级别
/fix --report lint-deep-python-20251222.md --level warning

# 结果：
# ✅ 修复 12 个 Error 级别问题
# ✅ 修复 25 个 Warning 级别问题
# ⏭️ 跳过 8 个 Suggestion 级别问题（注释类）
```

#### 场景5：修复所有问题（默认行为）

```bash
# 快速模式 - 自动修复所有级别
/lint main
# 等同于：/lint main && /fix --level=all

# 深度模式 - 手动修复所有级别
/lint --mode=deep main
/fix --report lint-deep-go-20251222.md --level=all
# 或直接省略 --level 参数（默认为 all）
/fix --report lint-deep-go-20251222.md

# 结果：
# ✅ 修复 12 个 Error 级别问题
# ✅ 修复 25 个 Warning 级别问题
# ✅ 修复 8 个 Suggestion 级别问题
```

#### 场景6：只修复特定文件

```bash
# 深度模式下只修复几个关键文件
/fix --report lint-deep-go-20251222.md --files src/app/main.go,src/service/auth.go
```

#### 场景7：组合使用 --level 和 --files

```bash
# 只修复特定文件的 Error 级别问题
/fix --report lint-deep-python-20251222.md --level error --files src/core/database.py,src/core/cache.py
```

## 工作流程

```
开始
  │
[步骤1] 读取问题数据 ⭐ 优化：支持JSON和Markdown两种方式
  ├─ **方式1：直接读取JSON数据**（优先，从/lint自动触发时使用）⭐ 新增
  │  ├─ 参数：--data-file=.claude/temp/report-data-merged-{timestamp}.json
  │  ├─ 读取合并后的问题数据
  │  └─ 优势：避免Markdown解析，数据更准确
  │
  └─ **方式2：解析Markdown报告**（传统方式，手动执行时使用）
     ├─ 查找报告文件 (report-reader)：
     │  ├─ 如果指定 --report，使用指定报告
     │  └─ 否则，自动查找最新报告（优先深度模式）
     └─ 解析报告内容：
        ├─ 提取问题列表（ID、级别、位置、描述、修复建议）
  │  ├─ 提取检查元数据（语言、分支、范围、模式）
  │  └─ 验证报告格式完整性
  └─ 输出结构化问题数据供后续步骤使用
  ↓
[步骤2] 分析和筛选问题 (issue-analyzer)
  ├─ 按 --level 参数筛选问题级别
  ├─ 按 --files 参数筛选文件
  ├─ 问题优先级排序（error > warning）
  └─ 评估修复风险
  ↓
[步骤3] 创建备份
  ├─ git stash（如果在git仓库）
  └─ 或复制文件到 .backup/
  ↓
[步骤4] 执行修复 (code-fixer-*)
  ├─ Go: code-fixer-go
  ├─ Python: code-fixer-python
  └─ Java: code-fixer-java
  ↓
[步骤5] 验证修复结果
  ├─ 检查语法正确性
  ├─ 统计修复成功/失败/跳过数量
  └─ 如有失败，记录详情
  ↓
[步骤6] 生成修复报告 (fix-report-generator)
  ├─ 快速模式: doc/fix/fix-fast-{language}-{date}.md
  └─ 深度模式: doc/fix/fix-deep-{language}-{date}.md
  ↓
[步骤7] 清理临时文件 (cleanup-handler)
  ├─ 如果修复成功：删除 .backup/
  ├─ 删除 .claude/temp/
  ├─ 删除 .bak 文件
  └─ 显示清理结果
  ↓
[步骤8] 输出修复摘要
  ├─ 显示修复统计
  ├─ 显示报告路径
  └─ 提供后续建议
```

## 快速模式 vs 深度模式

| 特性 | 快速模式 | 深度模式 |
|------|---------|---------|
| **触发方式** | /lint 自动触发 | 用户手动执行 /fix |
| **用户干预** | 无需干预 | 需查看报告后确认 |
| **报告检测** | 自动使用最新报告 | 自动使用最新深度报告 |
| **报告类型** | 基础检查报告 | 深度分析报告 |
| **修复范围** | 所有 error + warning | 可选择性修复 |
| **适用场景** | 日常开发、CI/CD | 重构、发布前检查 |
| **报告文件** | `fix-fast-*.md` | `fix-deep-*.md` |

## 执行步骤详解

### 步骤1: 查找和读取报告

#### 1.1 查找报告

如果没有指定 `--report`，自动查找最新报告：

```bash
# 在 doc/lint/ 目录中查找最新报告
ls -t doc/lint/lint-*.md 2>/dev/null | head -1
```

**查找优先级**：
1. 如果存在深度模式报告（`lint-*-deep-*.md`），优先使用最新的深度报告
2. 否则使用最新的快速模式报告（`lint-*-fast-*.md`）
3. 如果都不存在，提示用户先执行 lint 命令

**输出示例**（找到报告）：
```
使用报告：lint-deep-go-20251222-143000.md
```

**输出示例**（未找到报告）：
```
❌ 错误：未找到 lint 报告

报告目录: doc/lint/
当前状态: 目录为空或不存在

建议：
- 快速模式：执行 /lint [branch]（会自动修复）
- 深度模式：执行 /lint --mode=deep [branch]，然后手动 /fix

查看 lint 命令帮助：
/lint --help
```

#### 1.2 解析报告内容 (report-reader)

**内部实现**（对用户透明，自动执行）：

通过 `report-reader` skill 解析报告，提取结构化问题数据：

**解析内容**：
- ✅ 提取问题列表（ID、级别、位置、描述、修复建议）
- ✅ 提取检查元数据（语言、分支、范围、模式）
- ✅ 验证报告格式完整性
- ✅ 统计问题数量（按级别分类）

**输出结构化数据**：
```json
{
  "source_report": "lint-fast-go-20251222.md",
  "mode": "fast",
  "language": "go",
  "total_issues": 28,
  "issues": [
    {
      "id": "E001",
      "level": "error",
      "category": "naming",
      "file": "src/app/main.go",
      "line": 42,
      "message": "...",
      "suggestion": "..."
    }
  ]
}
```

### 步骤2: 分析和筛选问题

#### 2.1 调用SKILL /issue-analyzer

通过 `issue-analyzer` skill 分析问题、评估风险、生成修复计划。

**筛选逻辑**：
- `--level error`: 只保留 error 级别问题
- `--level warning`: 只保留 warning 级别问题
- `--level all`: 保留所有 error + warning（默认）
- `--files`: 只保留指定文件的问题

**问题优先级排序**：
1. Error 级别 > Warning 级别
2. 安全问题 > 规范问题
3. 低风险 > 高风险

**输出示例**：
```
修复计划：
- 待修复问题：18个
  - Error: 12个
  - Warning: 6个
- 自动修复：17个
- 跳过修复：1个（高风险，需人工处理）

开始修复...
```

### 步骤3: 创建备份

**备份策略**：
```bash
# 如果在 git 仓库中
git stash push -m "backup before fix"

# 否则复制文件
mkdir -p .backup/{timestamp}/
cp {modified_files} .backup/{timestamp}/
```

**输出示例**：
```
已创建备份：git stash
```

### 步骤4: 执行修复

根据语言类型调用对应的 fixer skill：

**Go**: `code-fixer-go`
**Python**: `code-fixer-python`
**Java**: `code-fixer-java`

**修复流程**：
1. 读取源文件
2. 根据问题位置和类型应用修复规则
3. 验证修复后的代码语法
4. 如果验证失败，回滚该修复
5. 记录修复结果

**输出示例**（实时）：
```
正在修复...
  ✅ [1/18] src/app/main.go:42 - 包名规范化
  ✅ [2/18] src/app/main.go:85 - 添加错误处理
  ✅ [3/18] src/utils/helper.go:23 - 变量命名优化
  ...
  ⏭️ [18/18] src/service/auth.go:156 - 跳过（高风险，需人工处理）
```

### 步骤5: 验证修复结果

**验证内容**：
- 语法正确性（编译/解析）
- 修复是否成功应用
- 是否引入新问题

**输出示例**：
```
修复验证通过
   - 成功: 17个
   - 失败: 0个
   - 跳过: 1个
```

### 步骤6: 生成修复报告

通过 `fix-report-generator` skill 生成修复报告。

**报告文件命名**：
- **快速模式**: `fix-fast-{language}-{YYYYMMDD-HHMMSS}.md`
- **深度模式**: `fix-deep-{language}-{YYYYMMDD-HHMMSS}.md`

**报告内容**：
- 修复概要和统计
- 所有修复详情（修复前后对比）
- 未修复问题及原因
- 后续建议

### 步骤7: 清理临时文件

通过 `cleanup-handler` skill 清理临时文件，根据修复状态采用不同策略：

**修复成功时（完整清理）**：
```bash
# 使用 Task 工具调用 cleanup-handler
Task(
  subagent_type="general-purpose",
  description="清理fix临时文件",
  prompt="调用SKILL /cleanup-handler 清理临时文件：

命令: fix
状态: success

清理范围：
- .backup/ (修复成功，删除备份)
- .claude/temp/
- .claude/**/*.bak

请执行清理并显示结果。"
)
```

**修复失败时（部分清理）**：
```bash
# 使用 Task 工具调用 cleanup-handler
Task(
  subagent_type="general-purpose",
  description="清理fix临时文件",
  prompt="调用SKILL /cleanup-handler 清理临时文件：

命令: fix
状态: failed

清理范围：
- .claude/temp/
- .claude/**/*.bak

保留：
- .backup/ (修复失败，保留备份用于回滚)

请执行清理并显示结果，提示用户备份位置。"
)
```

**清理策略**：
- ✅ 修复成功：删除所有临时文件和备份（.backup/、.claude/temp/、*.bak）
- ⚠️ 修复失败：保留 .backup/ 目录供手动回滚，删除其他临时文件

**保留文件**：
- ✅ 修复后的源代码
- ✅ fix 报告：`doc/fix/fix-*.md`
- ✅ lint 报告：`doc/lint/lint-*.md`
- ⚠️ .backup/ 目录（仅修复失败时保留）

### 步骤8: 输出修复摘要

**输出示例**：
```
代码修复完成！

修复统计：
- 源报告: lint-fast-go-20251222.md
- 检测模式: 快速模式（自动修复）
- 修复语言: Go
- 修复时间: 2025-12-22 10:30:00
- 修复文件: 5个
- 修复问题: 18个
  - ✅ 成功: 17个 (94.4%)
  - ❌ 失败: 0个 (0%)
  - ⏭️ 跳过: 1个 (5.6%)

详细报告: doc/fix/fix-fast-go-20251222.md

临时文件已清理

建议:
- ✅ Error 级别问题已全部修复
- ⚠️ 1个 Warning 问题需人工处理（高风险）
```

## 修复报告格式

### 快速模式修复报告

```markdown
# {Language} 代码修复报告（快速模式）

## 修复概要
- 检测模式: 快速模式（自动修复）
- 源报告: lint-fast-go-20251222.md
- 修复时间: 2025-12-22 10:30:00
- 修复文件: 5个
- 修复问题: 18个

## 修复统计
- ✅ 成功: 17个 (94.4%)
- ❌ 失败: 0个 (0%)
- ⏭️ 跳过: 1个 (5.6%)

## 修复详情

### 文件: src/app/main.go

#### [E001] 包名规范化 ✅
- **位置**: 第 42 行
- **级别**: Error
- **修复前**:
  ```go
  package dsp_bad_code
  ```
- **修复后**:
  ```go
  package dspbadcode
  ```

... (其他修复详情)

## 未修复问题

### [W015] 函数复杂度过高 ⏭️
- **文件**: src/service/auth.go
- **位置**: 第 156 行
- **原因**: 高风险，需要重构整个函数，建议人工处理
- **建议**: 将函数拆分为多个小函数

## 后续建议
- ✅ Error 级别问题已全部修复
- ⚠️ 1个 Warning 问题需人工处理
```

### 深度模式修复报告

```markdown
# {Language} 代码修复报告（深度模式）

## 修复概要
- 检测模式: 深度模式（人工确认修复）
- 源报告: lint-deep-python-20251222.md
- 深度分析: 依赖关系 + 调用链 + 架构检查
- 修复时间: 2025-12-22 15:45:00

## 修复统计
[同快速模式]

## 深度分析影响评估

### 依赖影响
- 修复影响 8 个模块
- 反向依赖检查通过

### 调用链影响
- 修复影响 5 个上游调用者
- 无破坏性变更

### 架构一致性
- 修复后符合分层架构
- 设计模式一致性保持

## 修复详情
[同快速模式]

## 未修复问题
[同快速模式]

## 后续建议
- ✅ 基础问题已修复
- ✅ 架构一致性已改善
- 💡 建议重新运行深度分析验证影响
```

## 错误处理

### 常见错误及解决方案

#### 1. 未找到报告

```
❌ 错误：未找到 lint 报告

可能原因：
1. 尚未执行代码规范检查
2. 报告目录为空

建议：
- 快速模式：执行 /lint [branch]（会自动修复）
- 深度模式：执行 /lint --mode=deep [branch]，然后手动 /fix
```

#### 2. 指定的报告不存在

```
❌ 错误：报告文件不存在

指定的报告: lint-deep-go-xxx.md
报告目录: doc/lint/

建议：
- 检查报告文件名是否正确
- 查看可用报告：ls doc/lint/
```

#### 3. 报告中没有需要修复的问题

```
ℹ️ 提示：报告中没有需要修复的问题

可能原因：
1. 代码检查未发现 error/warning 问题
2. 使用了 --level 参数，过滤掉了所有问题

建议：
- 查看 lint 报告内容
- 调整 --level 参数
```

#### 4. 修复失败

```
⚠️ 警告：部分问题修复失败

修复失败的问题：
1. src/app/main.go:125 - 语法验证失败
   原因：修复后代码存在语法错误
   建议：请手动检查并修复

已自动回滚失败的修复
```

## 备份和回滚

### 备份机制

**Git 仓库**：
```bash
git stash list  # 查看备份
git stash pop   # 恢复最近的备份
```

**非 Git 环境**：
```
.backup/
└── 20251222-103000/          # 时间戳目录
    ├── src/app/main.go       # 原文件
    └── src/utils/helper.py
```

### 回滚方法

```bash
# 查看备份
ls -la .backup/

# 回滚单个文件
cp .backup/20251222-103000/src/app/main.go src/app/main.go

# 回滚所有文件
cp -r .backup/20251222-103000/* ./
```

## 强约束规则

### 1. 只修复报告中的问题

- ✅ **必须**: 只修复 lint 报告中明确列出的问题
- ❌ **禁止**: 修复报告之外的任何问题
- ❌ **禁止**: 进行"顺手优化"或"代码美化"

### 2. 覆盖所有可修复问题

- ✅ **必须**: 尝试修复报告中所有 error 和 warning 问题
- ✅ **必须**: 如果无法自动修复，必须标记为"跳过"并说明原因
- ❌ **禁止**: 遗漏报告中的任何问题

### 3. 修复范围限制

- ✅ **必须**: 只修改问题所在的代码行或最小必要范围
- ❌ **禁止**: 重构整个函数或模块
- ❌ **禁止**: 修改相关但不在报告中的代码

## 扩展性

### 添加新语言支持

1. 创建 `skills/code-fixer-{language}/SKILL.md`
2. 在本 command 的 skills meta 信息中添加
3. 在步骤4中添加语言判断逻辑

---

**本命令提供智能化的代码修复能力，在快速模式下实现自动化修复，在深度模式下提供可控的修复流程。**
