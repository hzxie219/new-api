---
skills: design-refiner
version: 1.0
description: |
  设计细化 - 将概要设计文档细化成详细设计文档。
  使用场景：用户描述新功能需求需要生成详细设计文档、用户提供概要设计草案需要完善成详细设计、需要基于项目知识库生成符合架构规范的设计方案。
  核心能力：设计讨论（双向澄清直到理解），深度源码分析（探索代码、识别可复用资源），生成详细设计文档。
allowed-tools: Read, Glob, Grep, Write, Task, AskUserQuestion
---

# 设计细化

将概要设计文档细化为详细设计文档。


# 任务描述

## 全局约束

- 调用SKILL /design-refiner

## 处理步骤

### 调用SKILL /design-refiner

```
输入文件 doc/design/tech_design.md、doc/testcase/testcase.md
输出文件：doc/develop/tech_design_detail.md
```

#### 项目知识库参考

根据需求类型查阅对应知识库：
- **新增 API**：参考 `doc/kb/技术知识库/API索引.md`
- **业务逻辑**：参考 `doc/kb/业务知识库/`
- **数据结构**：参考 `doc/kb/技术知识库/数据结构.md`
- **测试相关**：参考 `doc/kb/测试知识库/单元测试.md`
