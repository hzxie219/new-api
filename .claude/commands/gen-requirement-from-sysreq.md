---
allowed-tools: Bash,Skill,Read,Glob,Grep
description: 将系统需求文档转换为包含代码引用的后端需求分析文档
skills: requirement-analyzer-sysreq
version: 1.0
---

# 任务名称

将系统需求文档转换为包含代码引用的后端需求分析文档

# 任务描述

## 全局约束

- **一次性自动完成**，无需用户干预
- 全程自动化，不询问用户
- **禁止编造、猜测或假设**接口路径、代码位置等信息
- 所有代码引用必须指向真实代码文件，禁止引用 `docs/` 下的文档
- 上下文窗口会在接近限制时自动压缩，不要因为Token预算问题提前停止任务
- 最终只输出一个需求分析文档

## 输入来源说明

本命令接收来自 `/system-analyst` 的输出作为输入：
- 固定路径：`doc/requirement-analyst/output/`
- 该目录包含 Epic、Feature、Story、Tech 层级的系统需求文档

## 处理步骤

### 步骤1: 验证输入

验证输入目录和代码仓库：

**验证要求**：
- 确认 `doc/requirement-analyst/output/` 目录存在
- 确认目录下包含系统需求文档（README.md、Tech*.md 等）
- 验证代码仓库存在：项目根目录下除了 `.claude`、`README.md`、`doc` 之外的其他文件夹

如果验证失败，给出明确错误提示：
```
错误：未找到系统需求文档
请先执行 /system-analyst 生成系统需求文档
```

### 步骤2: 调用技能生成需求分析文档

调用SKILL /requirement-analyzer-sysreq：

```
调用SKILL /requirement-analyzer-sysreq
```

该技能将自动完成：
- 解析系统需求文档集（Epic/Feature/Story/Tech 层级）
- 搜索真实代码并建立索引
- 功能点分类验证与代码引用补充
- 生成规范的后端需求分析文档

输出文件：`doc/requirement/requirement.md`

# 工具说明

- **Bash**: 用于执行文件扫描命令
- 调用SKILL /requirement-analyzer-sysreq 执行核心分析任务
- **Read**: 读取文件内容
- **Glob**: 查找文件
- **Grep**: 搜索文件内容

# 输入说明

## 输入目录结构

输入目录固定为：`doc/requirement-analyst/output/`

该目录由 `/system-analyst` 生成，包含分层的系统需求文档。

## 知识库索引

- **代码仓库知识库**：项目根目录下各代码仓库文件夹（排除 `.claude`、`README.md`、`doc`）的 `doc/kb/` 目录下的知识库

# 输出模板

- @.claude/skills/requirement-analyzer-sysreq/references/output-template.md

# 全局补充输入（动态补充）

## 命令参数

```
/gen-requirement-from-sysreq [补充说明]
```

- 无必选参数（输入路径固定）
- `[补充说明]` (可选): 用户提供的额外上下文或特殊要求

## 使用说明

### 前置条件

在执行本命令之前，需要先执行以下命令生成系统需求文档：
1. `/ai-business-analyst` - 将原始需求转换为 AI 需求（Epic→Feature→Story）
2. `/system-analyst` - 将 AI 需求细化为系统需求（包含 Tech 层级）

### 处理方式

- 自动读取 `doc/requirement-analyst/output/` 目录下的所有需求文档
- 如果用户提供了补充说明，在执行步骤2调用 `requirement-analyzer-sysreq` 技能时，将补充说明作为额外上下文传递
- 补充说明会影响：
  - 功能点的优先级判断
  - 复杂度评估
  - 风险分析
  - 技术方案建议

### 默认行为

如果未提供补充说明，则完全基于系统需求文档和代码仓库进行自动分析。