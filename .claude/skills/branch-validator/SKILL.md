---
name: "branch-validator"
description: "Git 分支名验证器，验证分支是否存在且大小写正确"
version: 2.0
tools:
  - Bash
---

您是 Git 分支验证专家，负责验证用户提供的分支名是否正确存在。

## 核心职责

1. **精确验证分支名**：验证分支是否存在（区分大小写）
2. **检测大小写错误**：发现并提示分支名大小写错误
3. **提供修正建议**：自动纠正大小写错误或提供可用分支列表
4. **支持本地和远程分支**：同时检查本地分支和远程分支

## 为什么需要这个 Skill？

### 问题背景

在 Windows 系统上，Git 默认配置为 `core.ignorecase = true`，导致：

```bash
# 实际分支名
develop-DSP3.0.33

# 用户输入（大小写错误）
develop-DSp3.0.33

# Git 不会报错！静默接受错误的分支名
git diff develop-DSp3.0.33...feature  # ✅ 执行成功
git rev-parse develop-DSp3.0.33        # ✅ 返回正确的 commit hash
```

这会导致：
- 使用了错误的分支名但没有任何警告
- 可能在某些环境下产生不一致的行为
- 报告文件命名使用错误的分支名
- 难以追溯和调试问题

## 工作流程

```
输入: branch_name
  ↓
1. 获取所有分支列表
  ↓
2. 精确匹配（区分大小写）
  ↓
3a. 匹配成功 → 返回验证通过
  ↓
3b. 匹配失败 → 检查是否有大小写不同的类似分支
  ↓
4a. 找到类似分支 → 自动纠正 + 警告
  ↓
4b. 未找到 → 报错 + 提供可用分支列表
```

## 验证算法

### 步骤 1: 获取所有分支

```bash
# 获取本地和远程分支
git branch -a
```

输出示例：
```
* feature-DSP3.0.33-dynamic-grading
  develop-DSP3.0.33
  main
  remotes/origin/develop-DSP3.0.33
  remotes/origin/feature-DSP3.0.33-dynamic-grading
  remotes/origin/main
```

### 步骤 2: 解析分支名

从输出中提取分支名：
- 去除 `* ` 前缀（当前分支标记）
- 去除 `remotes/origin/` 前缀（远程分支）
- 去重（本地和远程可能有同名分支）

### 步骤 3: 精确匹配（区分大小写）

```python
if branch_name in branch_list:
    return {
        "valid": True,
        "original_name": branch_name,
        "corrected_name": branch_name,
        "message": "分支名验证通过"
    }
```

### 步骤 4: 大小写不敏感匹配

如果精确匹配失败，尝试大小写不敏感匹配：

```python
for actual_branch in branch_list:
    if actual_branch.lower() == branch_name.lower():
        return {
            "valid": False,
            "auto_corrected": True,
            "original_name": branch_name,
            "corrected_name": actual_branch,
            "message": f"分支名大小写已自动纠正: {branch_name} → {actual_branch}"
        }
```

### 步骤 5: 完全不匹配

如果都不匹配，返回错误和建议：

```python
# 查找相似分支（编辑距离 <= 3）
similar_branches = find_similar_branches(branch_name, branch_list, max_distance=3)

return {
    "valid": False,
    "auto_corrected": False,
    "original_name": branch_name,
    "corrected_name": None,
    "message": f"分支 '{branch_name}' 不存在",
    "suggestions": similar_branches or ["main", "master", "develop"]
}
```

## 输入参数

### 参数接收方式

本 skill 从用户的 `/lint` 命令中接收分支名参数。支持以下格式：

1. **位置参数**: `/lint develop-DSP3.0.33`
   - 自动提取第一个位置参数作为分支名

2. **命名参数**: `/lint --branch=develop-DSP3.0.33`
   - 从 `--branch` 参数提取分支名

3. **默认值**: 如果未提供，默认使用 `main`

### 参数定义

- **branch** (可选): 要验证的分支名
  - 类型: 字符串
  - 默认值: `main`
  - 示例: `develop-DSP3.0.33`, `feature-auth`, `main`

## 输出格式

### 成功验证

```json
{
  "valid": true,
  "original_name": "develop-DSP3.0.33",
  "corrected_name": "develop-DSP3.0.33",
  "message": "✅ 分支名验证通过"
}
```

### 自动纠正（大小写错误）

```json
{
  "valid": false,
  "auto_corrected": true,
  "original_name": "develop-DSp3.0.33",
  "corrected_name": "develop-DSP3.0.33",
  "message": "⚠️ 分支名大小写已自动纠正",
  "warning": "您输入: develop-DSp3.0.33\n已纠正为: develop-DSP3.0.33\n\nGit 分支名区分大小写，建议使用正确的分支名。"
}
```

### 分支不存在

```json
{
  "valid": false,
  "auto_corrected": false,
  "original_name": "develop-DSP3.0.32",
  "corrected_name": null,
  "message": "❌ 分支 'develop-DSP3.0.32' 不存在",
  "suggestions": [
    "develop-DSP3.0.33",
    "develop-DSP3.0.31",
    "main"
  ],
  "help": "请使用以下命令查看所有分支:\ngit branch -a"
}
```

## 执行示例

### 示例 1: 分支名正确

**输入**: `develop-DSP3.0.33`

**执行**:
```bash
git branch -a | grep -E "^\s*(\\*\\s+)?develop-DSP3.0.33$|remotes/origin/develop-DSP3.0.33$"
```

**输出**:
```
✅ 分支名验证通过

分支名: develop-DSP3.0.33
类型: 本地分支 + 远程分支
```

### 示例 2: 大小写错误

**输入**: `develop-DSp3.0.33`

**执行**:
```bash
# 1. 精确匹配失败
# 2. 大小写不敏感匹配成功
```

**输出**:
```
⚠️ 分支名大小写已自动纠正

您输入: develop-DSp3.0.33
已纠正: develop-DSP3.0.33

Git 分支名区分大小写，建议使用正确的分支名。

🔄 已使用纠正后的分支名继续执行
```

### 示例 3: 分支不存在

**输入**: `develop-DSP3.0.32`

**执行**:
```bash
# 精确匹配和大小写匹配都失败
# 查找相似分支
```

**输出**:
```
❌ 分支 'develop-DSP3.0.32' 不存在

💡 您可能想使用以下分支:
  1. develop-DSP3.0.33
  2. develop-DSP3.0.31
  3. main

查看所有分支:
  git branch -a

查看远程分支:
  git branch -r
```

## 集成到 lint command

在 lint command 的步骤 1 中调用：

### 调用时机

- **增量模式** (`--scope=incremental`): 必须调用
- **全量模式** (`--scope=full`): 不需要调用
- **最新模式** (`--scope=latest`): 不需要调用

### 调用方式

```python
# 1. 解析用户输入
# 从 /lint [branch] 中提取分支名
user_input = "develop-DSP3.0.33"  # 用户输入的位置参数
# 或从 /lint --branch=develop-DSP3.0.33 中提取
# 如果都没有，使用默认值 "main"

# 2. 调用 branch-validator
result = Skill(
    skill="branch-validator",
    args=f"{user_input}"  # 直接传递分支名，无需 --branch= 前缀
)

# 3. 处理验证结果
if result["auto_corrected"]:
    # 大小写被自动纠正
    print(f"⚠️ 分支名已纠正: {result['original_name']} → {result['corrected_name']}")
    branch_name = result["corrected_name"]
    # 继续执行，使用纠正后的分支名

elif result["valid"]:
    # 验证通过
    branch_name = result["corrected_name"]
    # 继续执行

else:
    # 分支不存在
    print(f"❌ {result['message']}")
    print(f"💡 建议: {result['suggestions']}")
    # 终止执行
    return
```

### 完整步骤示例

```markdown
### 步骤 1: 分支验证、语言检测和参数解析

1. **解析命令参数**
   - 从用户输入 `/lint develop-DSP3.0.33` 提取分支名
   - 如果是全量模式或最新模式，跳过分支验证

2. **验证分支名**（仅增量模式）
   - 调用 branch-validator skill
   - 处理验证结果（通过/纠正/终止）
   - 使用验证后的分支名继续执行

3. 检测项目语言类型（Go/Python/Java）
4. 获取变更文件列表
5. 解析变更行号范围
```

## 使用场景

### 场景 1: 日常开发检查

```bash
/lint develop-DSP3.0.33
```

系统自动调用 `branch-validator`，验证分支名正确性。

### 场景 2: 用户输入错误

```bash
/lint develop-DSp3.0.33  # 大小写错误
```

**系统输出**:
```
⚠️ 分支名大小写已自动纠正
   develop-DSp3.0.33 → develop-DSP3.0.33

继续执行代码检查...
```

### 场景 3: 分支不存在

```bash
/lint old-branch
```

**系统输出**:
```
❌ 分支 'old-branch' 不存在

建议使用:
  1. develop-DSP3.0.33
  2. main

已终止执行
```

## 错误处理

### Git 仓库不存在

```bash
fatal: not a git repository
```

**处理**: 返回错误，提示用户当前目录不是 Git 仓库

### Git 命令失败

```bash
git: command not found
```

**处理**: 返回错误，提示用户安装 Git

### 分支列表为空

**处理**: 返回错误，提示用户仓库可能损坏或未初始化

## 性能优化

- **缓存分支列表**: 在单次 lint 执行中缓存分支列表（避免多次调用 git branch）
- **快速失败**: 精确匹配成功后立即返回
- **异步执行**: 分支验证可与其他无关步骤并行执行

## 配置选项（可选）

### 严格模式

```json
{
  "strict_mode": true,  // 发现大小写错误时报错而不是自动纠正
  "allow_remote_only": true,  // 允许只在远程存在的分支
  "similarity_threshold": 3  // 相似分支的最大编辑距离
}
```

---

**重要**: 本 skill 是代码检查流程的第一道防线，确保后续操作使用正确的分支名。
