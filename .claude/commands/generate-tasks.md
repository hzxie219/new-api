---
skills: task-generator
version: 1.0
description: |
  任务生成 - 将详细设计文档转换成具体的可执行任务清单。
  使用场景：用户已有详细设计文档需要生成编码任务、用户确认设计后需要生成 feature-task.md 指导 AI 编码、需要将复杂功能拆解为可独立执行的任务单元。
  核心能力：读取设计文档，确认单测需求，分析功能点，任务拆解，依赖分析，生成 feature-task.md。
allowed-tools: Read, Glob, Grep, Write, AskUserQuestion
---

# 任务生成

将详细设计文档转换为可执行任务清单。


# 任务描述

## 全局约束

- 调用SKILL /task-generator

## 处理步骤

### 调用SKILL /task-generator

```
输入文件：doc/develop/tech_design_detail.md
输出文件：doc/develop/feature-task.md
```

#### 项目知识库参考

根据需求类型查阅对应知识库：
- **新增 API**：参考 `doc/kb/技术知识库/API索引.md`
- **业务逻辑**：参考 `doc/kb/业务知识库/`
- **数据结构**：参考 `doc/kb/技术知识库/数据结构.md`
- **测试相关**：参考 `doc/kb/测试知识库/单元测试.md`
