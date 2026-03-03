---
skills: task-executor,coder
version: 1.0
description: |
  任务执行 - 解析 feature-task.md、调度编码任务、更新任务状态。
  触发场景：用户指定需求目录进行编码、要求按需求文档完成任务清单、要求执行任务列表。
  核心能力：按顺序阅读需求文档，解析任务列表，调用 coder skill 进行编码实现，维护任务状态。
allowed-tools: Read, Glob, Grep, Edit, Write, Bash, Skill, Task
---
# 任务执行

根据任务清单执行编码任务，自动调度 coder skill 完成实现。


# 任务描述

## 全局约束

- 调用SKILL /task-executor

## 处理步骤

### 调用SKILL /task-executor

```
输入文件：
  - doc/develop/tech_design_detail.md（详细设计文档）
  - doc/develop/feature-task.md（任务清单）
输出：编码实现 + 任务状态更新
```

#### 项目知识库参考

根据任务类型查阅对应知识库：
- **新增 API**：参考 `doc/kb/技术知识库/API索引.md`
- **业务逻辑**：参考 `doc/kb/业务知识库/`
- **数据结构**：参考 `doc/kb/技术知识库/数据结构.md`
- **测试相关**：参考 `doc/kb/测试知识库/单元测试.md`
- **编码规范**：参考 `.claude/skills/coder/references/`
