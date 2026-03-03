# Claude Agents 文档

## 概述

本目录包含了用于Claude Code的各个专业代理（agent）定义文件。每个agent都有特定的职责和专业能力，可以在任务处理时被调用。

## Agent 列表

### 1. architect-review（架构评审专家）
- **文件**：architect-review.md
- **模型**：opus
- **职责**：负责对产品的架构设计和模块设计进行评审，并给出评审报告
- **工具**：Read

### 2. maintainability-design-expert（可维护性设计专家）
- **文件**：maintainability-design-expert.md
- **模型**：sonnet
- **职责**：负责模块的可维护性特性设计，包括可调试性、可运维性、可扩展性和可复用性设计
- **工具**：Read, Write, Edit, Bash

### 3. reliability-design-expert（可靠性设计专家）
- **文件**：reliability-design-expert.md
- **模型**：sonnet
- **职责**：负责模块的可靠性特性设计，包括过载保护、故障监控、FMEA分析和容错设计
- **工具**：Read, Write, Edit, Bash

### 4. requirements-analyzer（需求分析专家）
- **文件**：requirements-analyzer.md
- **模型**：sonnet
- **职责**：当用户需要分析、澄清或记录软件需求时使用此代理
- **工具**：所有可用工具

### 5. security-design-expert（安全性设计专家）
- **文件**：security-design-expert.md
- **模型**：sonnet
- **职责**：负责模块的安全性特性设计，包括安全分析、权限控制、数据加密和防护机制
- **工具**：Read, Write, Edit, Bash

### 6. solution-selection-expert（方案选型专家）
- **文件**：solution-selection-expert.md
- **模型**：sonnet
- **职责**：负责识别技术选型决策点，完成多方案对比分析，提供技术选型建议和决策支持
- **工具**：Read, Write, Edit, Bash

### 7. system-design-architect（系统设计架构师）
- **文件**：system-design-architect.md
- **模型**：opus
- **职责**：负责企业级系统的概要设计文档编写、架构设计、算法设计和测试用例规划
- **工具**：所有可用工具

### 8. tdd-performance-optimizer（TDD和性能优化专家）
- **文件**：tdd-performance-optimizer.md
- **模型**：opus
- **职责**：测试驱动开发和性能优化专家，精通TDD方法论、性能分析、算法优化、代码质量保证
- **工具**：所有可用工具

### 9. testability-design-expert（可测试性设计专家）
- **文件**：testability-design-expert.md
- **模型**：sonnet
- **职责**：负责模块的可测试性特性设计，包括单元测试、集成测试、性能测试和自动化测试方案
- **工具**：Read, Write, Edit, Bash

## 使用方法

在Task工具中使用时，通过`subagent_type`参数指定要使用的agent：

```python
Task(
    subagent_type="system-design-architect",
    description="生成系统设计文档",
    prompt="请为XX功能生成系统设计方案..."
)
```

## 格式规范

每个agent文件必须包含：
1. YAML前置元数据（frontmatter）
2. 必填字段：
   - `name`: agent名称（唯一标识符）
   - `description`: agent描述
   - `model`: 使用的模型（sonnet/opus/haiku）
   - `tools`: 可用工具列表（可选）

## 维护说明

- 新增agent时，请确保：
  - 文件名与name字段一致
  - 使用kebab-case命名
  - 包含完整的YAML前置元数据
  - 提供清晰的职责描述

- 修改agent时：
  - 保持向后兼容
  - 更新本README文档
  - 测试agent功能正常

## 注意事项

1. 所有agent文件都已经过格式验证，确保可以被Task工具正确引用
2. 每个agent都有特定的适用场景，请根据任务需求选择合适的agent
3. Opus模型适用于复杂设计任务，Sonnet适用于标准任务，Haiku适用于简单快速任务