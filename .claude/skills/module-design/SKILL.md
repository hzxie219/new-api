---
name: module-design
description: 根据需求文档自动生成完整的模块概要设计说明书。按照企业标准模板,系统化地生成包含设计方法参考、介绍、设计任务书、对外接口、概要说明(含DFX设计)、数据结构设计、流程设计、总结、变更控制和修订记录等完整章节的设计文档。适用于:(1) 需要基于需求文档生成概要设计 (2) 需要系统化的设计文档生成流程 (3) 需要符合企业标准的设计文档 (4) 用户提到"概要设计"、"模块设计"、"设计文档"、"生成设计"等关键词 (5) 需要包含架构设计、接口设计、DFX特性设计(安全性、可靠性、可测试性等)的完整设计文档
version: 1.0
---

# Module Design - 模块概要设计生成工具

## 重要说明

**Skill 包含的参考文件**：本 skill 的所有模板和规范文件都位于 skill 内部，在执行过程中可以直接读取这些文件。当 prompt 中提到以下文件时，请在 skill 的 references或templates 目录中查找：

- `design_template.md` - 设计文档标准模板
- `设计方案checklist.md` - 设计检查清单
- `RESTful_API格式规范v3.0.md` - API接口设计规范
- `FMEA分析输出格式定义.md` - FMEA分析格式规范
- `Chapter0-2.md` 至 `Chapter7-10.md` - 各章节生成指南

## 功能概述

module-design skill 是一个专业的模块概要设计文档生成工具,能够根据需求分析文档自动生成符合标准的完整概要设计说明书。

### 核心特性

- **系统化流程**: 按照10个标准步骤依次生成设计文档的各个章节
- **企业标准**: 严格遵循企业标准模板规范
- **增量开发支持**: 优先检索和复用现有代码,支持增量开发场景
- **DFX设计完整**: 包含安全性、可靠性、可测试性等7个维度的DFX特性设计
- **质量门禁**: 内置评审机制,确保设计质量达标
- **检查点机制**: 每个步骤生成检查点文件,保证设计过程可追溯

### 生成的文档结构

```
模块概要设计说明书/
├── 0. 设计方法参考
├── 1. 介绍
├── 2. 设计任务书
├── 3. 对外接口
├── 4. 概要说明
│   ├── 4.1. 背景描述
│   ├── 4.2. 方案选型
│   ├── 4.3. 静态结构
│   ├── 4.4. 对软件总体架构的影响
│   ├── 4.5. 概要流程
│   ├── 4.6. 关键特性设计(DFX)
│   │   ├── 4.6.1. 安全性设计
│   │   ├── 4.6.2. 可靠性设计
│   │   ├── 4.6.3. 可测试性设计
│   │   ├── 4.6.4. 可调试性设计
│   │   ├── 4.6.5. 可运维性设计
│   │   ├── 4.6.6. 可扩展性设计
│   │   ├── 4.6.7. 可复用性设计
│   │   ├── 4.6.8. 系统隐私设计
│   │   └── 4.6.9. 跨平台设计和平台差异处理
│   └── 4.7. 方案风险分析
├── 5. 数据结构设计
├── 6. 流程设计
├── 7. 总结
├── 9. 变更控制
└── 10. 修订记录
```

## 使用方法

### 基本用法

当用户提供需求文档并要求生成概要设计时,使用此skill:

```
请根据需求文档 [需求文档路径] 生成模块概要设计说明书
```

### 触发场景

- 用户明确要求生成"概要设计"、"模块设计"、"设计文档"
- 用户提供了需求文档并需要进行系统设计
- 用户需要符合企业标准的设计文档
- 用户需要包含DFX特性设计(安全性、可靠性等)的完整设计

### 参数说明

- **需求文档路径**: 可选参数，指向需求分析文档的路径
  - 如果提供：使用指定的需求文档路径
  - 如果不提供：自动使用默认路径 `doc/requirement/requirement.md`

## 执行流程

### 步骤1: 生成目标设计文档
**模板文件**: `templates/design_template.md`

**创建目标文件 `doc/design/tech_design.md`**

复制模板文件到目标位置,并使用固定名称:design_template.md



### 步骤2: 生成第0-2章 (基础章节)

**参考文件**: [Chapter0-2.md](references/chapters/Chapter0-2.md)

生成设计方法参考、介绍和设计任务书章节,包括:
- 需求分析与功能点识别
- 设计原则和方法论
- 需求跟踪表
- 模块整体目标

### 步骤3: 生成第3章 (对外接口)

**参考文件**: [Chapter3.md](references/chapters/Chapter3.md)

生成对外接口设计,包括:
- API接口设计(优先检索现有接口)
- 消息接口设计
- 接口规范和示例

### 步骤4: 生成第4章 (概要说明-架构部分)

**参考文件**: [Chapter4.md](references/chapters/Chapter4.md)

生成架构设计核心内容,包括:
- 背景描述(工作原理、应用场景、对手分析)
- 方案选型(多方案对比、选型结论)
- 静态结构(架构图、模块关系)
- 对软件总体架构的影响
- 概要流程(核心流程图)

### 步骤5: 生成第4章 (概要说明-DFX部分)

**参考文件**: [Chapter4-DFX.md](references/chapters/Chapter4-DFX.md)

生成DFX特性设计,包括:
- 安全性设计(STRIDE威胁建模)
- 可靠性设计(FMEA分析)
- 可测试性设计
- 可调试性设计
- 可运维性设计
- 可扩展性设计
- 可复用性设计
- 系统隐私设计
- 跨平台设计
- 方案风险分析

### 步骤6: 生成第5章 (数据结构设计)

**参考文件**: [Chapter5.md](references/chapters/Chapter5.md)

生成数据结构设计,包括:
- 配置文件定义
- 全局数据结构定义
- 数据表结构定义(如涉及)

### 步骤7: 生成第6章 (流程设计)

**参考文件**: [Chapter6.md](references/chapters/Chapter6.md)

为每个子模块生成详细流程设计,包括:
- 静态结构(类图、职责定义)
- 处理流程(流程图、时序图)
- 关键算法描述
- 数据结构定义
- 函数列表
- 设计要点检视

### 步骤8: 生成第7/9/10章 (总结与管理)

**参考文件**: [Chapter7-10.md](references/chapters/Chapter7-10.md)

生成总结和管理章节,包括:
- 关联分析(对老模块、老版本、相关产品的影响)
- 遗留问题解决
- 变更控制
- 修订记录

### 步骤9: 跳过第8章

第8章"业务逻辑相关的测试用例"由专门的测试设计工具负责,不在本skill范围内。

### 步骤10: 概要设计文档评审与质量门禁

基于七维度评分体系进行质量评审:
- 红线检查(6项必检项)
- 七维度评分(设计目标达成、可靠性、可测试性等)
- 分章节评审
- 评级判定(A/B+/B/C)
- 迭代优化(可选)

## 核心特性说明

### 增量开发优先

所有设计步骤都优先检索现有代码和知识库:
- API接口设计: 优先在知识库和代码中检索现有接口
- 数据结构设计: 优先在知识库中检索现有数据结构
- 流程设计: 优先检索现有模块实现
- 兼容性优先: 确保增量开发的向后兼容性

### DFX设计完整性

严格按照7个维度串行执行DFX设计:
1. 安全性设计(STRIDE威胁建模)
2. 可靠性设计(FMEA分析、过载保护)
3. 可测试性设计(测试覆盖率>80%)
4. 可调试性设计(日志、链路追踪)
5. 可运维性设计(黄金指标监控)
6. 可扩展性设计(扩展点、版本兼容)
7. 可复用性设计(公用代码识别)

### 质量门禁

- **红线门禁**: 6项红线检查必须全部通过
- **发布门禁**: 总评分必须≥85分(B+级)
- **优秀标准**: 总评分≥90分(A级)
- **迭代上限**: 最多3次迭代优化

### 检查点机制

每个步骤生成检查点文件,保证过程可追溯:
- `API接口_check.md`: API接口检索结果
- `功能点方案选型_check.md`: 方案选型决策点
- `DFX设计输入_check.md`: DFX设计前置分析
- `数据结构需求_check.md`: 数据结构需求分析
- `模块划分_check.md`: 模块划分结果
- `模块结构_check.md`: 模块代码检索结果

## 注意事项

### 执行原则

1. **逐步执行**: 严格按照步骤顺序执行,不跳过任何步骤
2. **强制写入**: 每个步骤完成后必须写入目标文件
3. **验证确认**: 每次写入后必须验证内容已正确写入
4. **失败重试**: 验证失败时必须重新执行写入操作
5. **完整性保证**: 即使token预算接近限制,也要尽可能完整执行

### 设计规范

- **模板遵循**: 严格按照企业标准模板生成
- **图表规范**: 使用Mermaid格式绘制所有图表
- **命名规范**: 保持与现有代码风格一致
- **兼容性优先**: 增量开发时优先考虑向后兼容

### 质量要求

- **需求覆盖**: 需求覆盖率必须100%
- **威胁建模**: STRIDE 6类威胁必须全覆盖
- **FMEA分析**: 关键业务流程必须完成FMEA分析
- **测试覆盖**: 测试覆盖率目标>80%
- **监控覆盖**: 黄金指标必须全覆盖

## 相关文档

### Reference文件说明

- **Chapter0-2.md**: 基础章节生成指南(第0-2章) - [references/chapters/Chapter0-2.md](references/chapters/Chapter0-2.md)
- **Chapter3.md**: 对外接口设计指南(第3章) - [references/chapters/Chapter3.md](references/chapters/Chapter3.md)
- **Chapter4.md**: 架构设计指南(第4.1-4.5章) - [references/chapters/Chapter4.md](references/chapters/Chapter4.md)
- **Chapter4-DFX.md**: DFX特性设计指南(第4.6-4.7章) - [references/chapters/Chapter4-DFX.md](references/chapters/Chapter4-DFX.md)
- **Chapter5.md**: 数据结构设计指南(第5章) - [references/chapters/Chapter5.md](references/chapters/Chapter5.md)
- **Chapter6.md**: 流程设计指南(第6章) - [references/chapters/Chapter6.md](references/chapters/Chapter6.md)
- **Chapter7-10.md**: 总结与管理章节指南(第7/9/10章) - [references/chapters/Chapter7-10.md](references/chapters/Chapter7-10.md)

每个reference文件都包含详细的执行步骤、生成内容结构、验证清单和注意事项。

### 模板文件

设计文档模板已包含在 skill 中:
- **design_template.md**: 概要设计文档标准模板，定义了完整的文档结构和章节要求（位于 [templates/](templates/) 目录）

### 规范文档

以下规范文档已包含在 skill 中，作为设计过程的参考依据:
- **RESTful_API格式规范v3.0.md**: API接口设计规范，定义了接口命名、参数格式、响应结构等标准（[specifications/RESTful_API格式规范v3.0.md](references/specifications/RESTful_API格式规范v3.0.md)）
- **FMEA分析输出格式定义.md**: FMEA(故障模式与影响分析)格式规范，用于可靠性设计（[specifications/FMEA分析输出格式定义.md](references/specifications/FMEA分析输出格式定义.md)）
- **设计方案checklist.md**: 设计检查清单，用于设计质量评审和验证（[specifications/设计方案checklist.md](references/specifications/设计方案checklist.md)）

### 专家代理文件

skill 中包含多个领域专家的定义文件，在设计过程中可按需引用：

#### 架构设计专家

- **architect-review.md**: 架构评审专家，负责对产品的架构设计以及模块设计进行评审 - [references/experts/architecture/architect-review.md](references/experts/architecture/architect-review.md)
- **system-design-architect.md**: 系统设计架构师，负责企业级系统的概要设计文档编写、架构设计、算法设计和测试用例规划 - [references/experts/architecture/system-design-architect.md](references/experts/architecture/system-design-architect.md)
- **solution-selection-expert.md**: 方案选型专家，负责识别技术选型决策点，完成多方案对比分析，提供技术选型建议和决策支持 - [references/experts/architecture/solution-selection-expert.md](references/experts/architecture/solution-selection-expert.md)

#### DFX特性设计专家

- **security-design-expert.md**: 安全性设计专家，负责模块的安全性特性设计，包括安全分析、权限控制、数据加密和防护机制 - [references/experts/dfx/security-design-expert.md](references/experts/dfx/security-design-expert.md)
- **reliability-design-expert.md**: 可靠性设计专家，负责模块的可靠性特性设计，包括过载保护、故障监控、FMEA分析和容错设计 - [references/experts/dfx/reliability-design-expert.md](references/experts/dfx/reliability-design-expert.md)
- **maintainability-design-expert.md**: 可维护性设计专家，负责模块的可维护性特性设计，包括可调试性、可运维性、可扩展性和可复用性设计 - [references/experts/dfx/maintainability-design-expert.md](references/experts/dfx/maintainability-design-expert.md)
- **testability-design-expert.md**: 可测试性设计专家，负责模块的可测试性特性设计，包括单元测试、集成测试、性能测试和自动化测试方案 - [references/experts/dfx/testability-design-expert.md](references/experts/dfx/testability-design-expert.md)

#### API设计专家

- **restful-api-generator.md**: RESTful API生成专家，基于需求文档和API设计规范生成完整的RESTful API接口定义 - [references/experts/api/restful-api-generator.md](references/experts/api/restful-api-generator.md)
- **restful-api-reviewer.md**: API评审专家，对API设计进行全面评审，确保符合RESTful规范和最佳实践 - [references/experts/api/restful-api-reviewer.md](references/experts/api/restful-api-reviewer.md)

#### 开发优化专家

- **tdd-performance-optimizer.md**: TDD与性能优化专家，精通TDD方法论、性能分析、算法优化和代码质量保证 - [references/experts/development/tdd-performance-optimizer.md](references/experts/development/tdd-performance-optimizer.md)

#### 需求分析专家

- **requirements-analyzer.md**: 需求分析专家，负责分析、澄清和记录软件需求，完成需求分解和结构化 - [references/experts/requirements/requirements-analyzer.md](references/experts/requirements/requirements-analyzer.md)

详细说明请参考 [experts/README.md](references/experts/README.md)

## 示例用法

### 场景1: 全新模块设计

```
用户: 请根据需求文档 doc/requirements/用户认证模块需求.md 生成概要设计说明书

助手: [启动 module-design skill]
- 步骤1: 创建目标设计文档
- 步骤2: 分析需求,生成第0-2章
- 步骤3-8: 依次生成各章节内容
- 步骤10: 执行质量评审
- 输出: 完整的概要设计说明书
```

### 场景2: 增量功能设计

```
用户: 基于现有的支付模块,生成支付宝支付方式的概要设计

助手: [启动 module-design skill]
- 自动检索现有支付模块代码
- 识别需要扩展的接口和数据结构
- 生成增量设计文档,标注新增/修改部分
- 确保向后兼容性
```

### 场景3: DFX设计专项

```
用户: 为订单模块补充DFX特性设计

助手: [启动 module-design skill,聚焦第4.6章]
- 完成STRIDE安全威胁建模
- 完成FMEA可靠性分析
- 设计测试方案(覆盖率>80%)
- 设计监控方案(黄金指标)
- 输出完整的DFX设计章节
```

## 总结

module-design skill 是一个强大的企业级概要设计文档生成工具,它:

- ✅ **系统化**: 10步标准流程,覆盖设计全生命周期
- ✅ **标准化**: 严格遵循企业标准模板规范
- ✅ **增量友好**: 优先复用现有代码,支持增量开发
- ✅ **质量保证**: 内置质量门禁,确保设计达标
- ✅ **可追溯**: 检查点机制,设计过程完全可追溯
- ✅ **DFX完整**: 7个维度完整覆盖非功能性设计

通过使用此skill,可以快速生成高质量、符合企业标准的模块概要设计说明书,大幅提升设计效率和质量。
