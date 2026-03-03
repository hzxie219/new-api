---
name: api-test-analyzer
description: Use when analyzing requirements documents to extract API interface information, business rules, and testing requirements - identifies all interfaces, parameters, constraints, and dependencies for comprehensive test case generation
version: 1.0
---

# API接口测试分析器

## 概述

基于TDD方法论，深度分析需求设计文档，系统化提取API接口信息、业务规则和测试需求，为后续测试用例设计提供完整的基础数据。

## 强制约束

1. **参数完整性**: 接口定义的所有参数必须100%提取，不允许遗漏
2. **参数追踪**: 每个参数必须记录来源位置（文档章节/行号）
3. **参数验证**: 输出前必须验证参数数量与原文档一致

## 何时使用

当需要从需求设计文档中提取接口测试信息时使用此技能：

```
开始分析 → 读取需求文档 → 提取接口信息 → 分析业务规则 → 识别测试需求 → 生成分析报告
```

**使用场景**：
- 需要分析新的API接口设计文档
- 需要更新现有接口的测试需求
- 需要识别接口之间的依赖关系
- 需要提取业务规则和约束条件

**不适用场景**：
- 已经有了完整的接口清单
- 需要生成具体的测试代码
- 需要执行性能测试或安全测试

## 核心分析模式

### 模式1: 接口识别与提取

**分析前**：手动翻阅文档，容易遗漏接口信息
```text
文档中散布着各种接口描述：
- 告警规则查询接口：/risk_visualization/risk_alarm_rules/list
- 删除区域处理规则接口：/risk_visualization/home_page/deleteRegion
- 新增区域处理规则接口：/risk_visualization/home_page/addRegion
```

**分析后**：结构化的接口清单
```json
{
  "interfaces": [
    {
      "id": "API001",
      "name": "告警规则查询接口",
      "method": "GET",
      "path": "/risk_visualization/risk_alarm_rules/list",
      "description": "查询告警规则列表，支持筛选和分页",
      "module": "告警管理"
    }
  ]
}
```

### 模式2: 参数约束分析

**分析前**：参数信息分散在不同章节
```text
- limit: 分页大小，有范围限制
- skip: 跳过数量，用于分页
- name: 规则名称，支持模糊搜索
```

**分析后**：完整的参数定义
```json
{
  "parameters": {
    "limit": {
      "type": "integer",
      "required": false,
      "constraints": {
        "min": 1,
        "max": 100,
        "default": 20
      },
      "description": "每页显示条数"
    }
  }
}
```

## 快速参考

### 核心分析任务

| 任务类型 | 分析重点 | 输出格式 | 优先级 |
|---------|----------|----------|--------|
| 接口识别 | HTTP方法、路径、功能 | 接口清单JSON | 高 |
| 参数分析 | 类型、约束、验证规则 | 参数定义JSON | 高 |
| 业务规则 | 约束条件、状态转换 | 业务规则列表 | 中 |
| 依赖分析 | 接口调用关系、数据依赖 | 依赖关系图 | 中 |

### 关键文档元素识别

| 元素类型 | 识别关键词 | 提取策略 |
|---------|-----------|----------|
| HTTP方法 | GET, POST, PUT, DELETE | 直接匹配 |
| HTTP路径 | /api/, /dashboard/, /app_merge/ | 正则表达式 |
| gRPC服务 | service, rpc, message, proto, Handler | Proto定义分析 |
| gRPC方法 | Server, Request, Response, pb. | 上下文分析 |
| 消息接口 | Topic, Pulsar, Kafka, 消息格式 | 消息定义分析 |
| 参数定义 | 参数, 请求, 请求体, message字段 | 上下文分析 |
| 错误码 | 错误码, HTTP状态码, gRPC status | 模式匹配 |
| 业务规则 | 规则, 约束, 条件 | 语义分析 |

### 接口类型识别

| 接口类型 | 识别特征 | 概设章节位置 | 测试模板 |
|---------|----------|-------------|----------|
| HTTP API | POST/GET + URL路径 | 3.1 API接口 | universal_test_template |
| gRPC | service定义, rpc方法, Handler | 3.1或独立章节 | grpc_test_template |
| 消息接口 | Topic, Pulsar/Kafka, 消息格式 | 3.2 消息接口 | 需Mock消息队列 |
| 内部接口 | 模块间调用, 函数签名 | 3.3 内部接口 | 单元测试 |

## 概设文档章节指引

基于性价比分析，以下是概设文档中各章节的阅读优先级：

### P0 必读章节（核心信息源）

| 章节 | 内容 | 提取目标 |
|------|------|----------|
| 3.对外接口 | HTTP/gRPC接口定义、参数、响应 | 接口清单、参数定义、约束条件 |
| 2.1需求跟踪 | 功能点与接口映射关系 | 接口分类、优先级、测试范围 |

### P1 次要章节（补充信息）

| 章节 | 内容 | 提取目标 |
|------|------|----------|
| 5.数据结构设计 | 参数对象、枚举定义 | 复杂参数类型、枚举值约束 |
| 1.2定义和缩写 | 业务术语解释 | 术语理解、命名规范 |

### 可跳过章节

- 1.1编写目的、1.3参考资料（文档元信息）
- 4.设计方法、4.3架构影响（架构决策，非测试必需）
- 6.变更记录（历史信息）

## 实现方法

### 步骤1: 文档预处理
```markdown
1. 使用Read工具读取概设文档
2. 优先定位"3.对外接口"章节，提取所有接口定义
3. 阅读"2.1需求跟踪"，建立功能点与接口映射
4. 按需查阅"5.数据结构设计"获取复杂参数定义
```

### 步骤2: 接口信息提取
```markdown
1. 使用Grep工具搜索接口模式：
   - 搜索"接口路径"、"apiPath"、"HTTP方法"
   - 搜索"/api/"、"/risk_visualization/"等URL模式
   - 搜索"GET"、"POST"、"PUT"、"DELETE"等HTTP方法

2. 提取接口元数据：
   - 接口名称和功能描述
   - HTTP方法和完整路径
   - 请求参数和响应格式
   - 权限要求和前置条件

3. 结构化存储：
   - 生成接口清单JSON
   - 建立接口索引
   - 标记接口分类和模块
```

### 步骤3: 参数约束分析
```markdown
1. 参数识别：
   - 从接口文档中提取所有参数
   - 区分路径参数、查询参数、请求体参数
   - 识别必填参数和可选参数

2. 约束提取：
   - 数据类型（string, integer, boolean, array）
   - 长度限制（minLength, maxLength）
   - 数值范围（minimum, maximum）
   - 格式要求（email, date, uuid）
   - 枚举值（enum）

3. 验证规则：
   - 必填验证（required）
   - 格式验证（pattern）
   - 业务规则验证
   - 权限验证
```

### 步骤4: 业务规则分析
```markdown
1. 规则识别：
   - 搜索"规则"、"约束"、"条件"等关键词
   - 分析业务流程和状态转换
   - 识别异常处理和错误码

2. 规则分类：
   - 数据约束：唯一性、完整性、关联性
   - 业务约束：状态转换、权限控制、时间限制
   - 技术约束：性能要求、安全要求

3. 验收条件提取：
   - 功能验收标准
   - 性能指标要求
   - 安全合规要求
```

### 步骤5: 依赖关系分析
```markdown
1. 接口依赖：
   - 分析接口调用链路
   - 识别前置条件依赖
   - 确定执行顺序

2. 数据依赖：
   - 分析数据流向
   - 识别共享数据
   - 确定数据一致性要求

3. 服务依赖：
   - 识别外部服务调用
   - 分析服务可用性要求
   - 确定降级策略
```

## 常见错误

### 错误1: 接口信息不完整
**问题**：只提取了接口路径，遗漏了HTTP方法或参数
```text
❌ 错误：只提取 "/risk_visualization/home_page/addRegion"
✅ 正确：提取 "POST /risk_visualization/home_page/addRegion" 及完整参数
```

**解决方案**：建立完整的接口信息检查清单
- [ ] HTTP方法
- [ ] 完整URL路径
- [ ] 请求参数定义
- [ ] 响应格式说明
- [ ] 错误处理说明

### 错误2: 参数约束忽略
**问题**：提取参数但忽略约束条件
```text
❌ 错误：limit 参数类型为 integer
✅ 正确：limit 参数类型为 integer，范围 1-100，默认值 20
```

**解决方案**：系统化检查参数属性
- [ ] 数据类型
- [ ] 长度/范围限制
- [ ] 默认值
- [ ] 必填/可选
- [ ] 格式要求

### 错误3: 业务规则遗漏
**问题**：只关注技术接口，忽略业务规则
```text
❌ 错误：只识别接口，不考虑业务约束
✅ 正确：识别"删除区域时不能有绑定的告警规则"等业务规则
```

**解决方案**：主动搜索业务约束
- 搜索"规则"、"约束"、"限制"、"条件"等关键词
- 分析异常处理和错误场景
- 理解业务流程和状态转换

## 真实世界影响

### 提升测试覆盖率
通过系统化的接口分析，确保100%的接口识别率，避免遗漏重要接口的测试用例。

### 提高测试质量
完整的参数约束分析确保测试用例覆盖所有边界条件和异常场景。

### 减少返工成本
准确的业务规则分析确保测试用例符合实际业务需求，减少后期修改。

## 输出模板

### 接口清单模板（HTTP API）
```json
{
  "analysis_metadata": {
    "document_path": "doc/design/tech_design.md",
    "analysis_time": "2025-01-XX",
    "total_interfaces": 15,
    "total_parameters": 45,
    "analysis_coverage": "100%",
    "interface_types": {
      "http_api": 12,
      "grpc": 2,
      "message": 1
    }
  },
  "parameter_tracking": {
    "extracted_count": 45,
    "verified": true,
    "verification_method": "逐项对照原文档"
  },
  "interfaces": [
    {
      "id": "API001",
      "type": "http_api",
      "name": "告警规则查询接口",
      "method": "GET",
      "path": "/risk_visualization/risk_alarm_rules/list",
      "description": "查询告警规则列表，支持筛选和分页",
      "module": "告警管理",
      "priority": "P0",
      "parameters": {
        "query": [
          {
            "name": "limit",
            "type": "integer",
            "required": false,
            "constraints": {
              "min": 1,
              "max": 100,
              "default": 20
            },
            "description": "每页显示条数",
            "source": {
              "section": "3.2.1 告警规则查询",
              "line": 125
            }
          }
        ],
        "total_count": 5
      },
      "response": {
        "success": "200 OK",
        "error_codes": ["400", "401", "403", "500"]
      },
      "business_rules": [
        "只能查询用户有权限的告警规则",
        "分页查询限制最大100条"
      ],
      "dependencies": [],
      "security_requirements": [
        "需要用户登录认证",
        "需要告警管理权限"
      ]
    }
  ],
  "summary": {
    "by_type": {
      "http_api": 12,
      "grpc": 2,
      "message": 1
    },
    "by_method": {
      "GET": 8,
      "POST": 4,
      "PUT": 2,
      "DELETE": 1
    },
    "by_priority": {
      "P0": 10,
      "P1": 4,
      "P2": 1
    },
    "by_module": {
      "告警管理": 8,
      "区域管理": 4,
      "规则管理": 3
    }
  }
}
```

### gRPC接口模板
```json
{
  "id": "GRPC001",
  "type": "grpc",
  "name": "数据流处理服务",
  "service": "DataFlowService",
  "handler": "DataFlowHandler",
  "proto_package": "dataflow",
  "methods": [
    {
      "name": "ProcessData",
      "request_type": "ProcessDataRequest",
      "response_type": "ProcessDataResponse",
      "parameters": [
        {
          "name": "data_id",
          "type": "string",
          "required": true,
          "description": "数据标识符"
        }
      ]
    }
  ],
  "test_template": "grpc_test_template"
}
```

### 消息接口模板
```json
{
  "id": "MSG001",
  "type": "message",
  "name": "内置分类同步消息",
  "topic": "data-identify-config-sync",
  "broker": "pulsar",
  "message_format": {
    "type": "category",
    "operation": "Modify|Delete",
    "data": []
  },
  "producer_module": "apolloserver",
  "consumer_module": "taskinference"
}
```

### 业务规则模板
```json
{
  "business_rules": [
    {
      "rule_id": "BR001",
      "description": "删除区域时不能有绑定的告警规则",
      "applicable_interfaces": [
        "DELETE /risk_visualization/home_page/deleteRegion"
      ],
      "constraint_type": "business_constraint",
      "validation_method": "precondition_check",
      "error_scenario": "区域存在绑定的告警规则时删除失败"
    }
  ]
}
```