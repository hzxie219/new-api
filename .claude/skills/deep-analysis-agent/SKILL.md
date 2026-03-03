---
name: deep-analysis-agent
description: 深度分析变更代码的依赖关系、调用链和架构一致性
version: 2.0
---

# Deep Analysis Agent Skill

## 功能概述

Deep Analysis Agent 是专门用于深度模式的分析引擎，负责对代码变更进行全面的依赖分析、调用链追踪和架构一致性检查。与快速模式不同，深度模式不仅关注变更代码本身，还会分析其对整个项目的潜在影响。

## 核心功能

### 1. 依赖关系分析

分析变更代码的依赖关系，识别可能受影响的代码模块。

#### Go 语言依赖分析
```bash
# 分析 import 依赖
- 直接导入的包
- 间接依赖的包
- 循环依赖检测

# 使用 go mod graph
go mod graph | grep "changed_module"

# 分析函数/类型依赖
- 被哪些包使用
- 使用了哪些外部类型
```

#### Python 语言依赖分析
```bash
# 分析 import 依赖
- from ... import 语句
- import 语句
- 动态导入

# 使用 pydeps 或静态分析
- 模块级依赖
- 函数级依赖
```

#### Java 语言依赖分析
```bash
# 分析 import 依赖
- 包级依赖
- 类级依赖
- Maven/Gradle 依赖树

# 使用 jdeps
jdeps -verbose changed_classes
```

### 2. 调用链追踪

追踪函数/方法的调用链，识别变更的影响范围。

#### 向上追踪（谁调用了变更的代码）
```
变更函数: processOrder()
  ↑
  被调用于: handleCheckout()
  ↑
  被调用于: CheckoutController.submit()
  ↑
  被调用于: API endpoint /api/checkout
```

**分析内容：**
- 直接调用者
- 间接调用者
- 调用频率（如果有日志数据）
- 关键路径识别

#### 向下追踪（变更代码调用了什么）
```
变更函数: processOrder()
  ↓
  调用: validateOrder()
  ↓
  调用: updateInventory()
  ↓
  调用: Database.update()
```

**分析内容：**
- 调用的内部函数
- 调用的外部 API
- 数据库操作
- 文件系统操作
- 网络请求

### 3. 架构一致性检查

检查变更代码是否符合项目的架构规范和设计模式。

#### 分层架构检查
```
检查点:
  • Controller 层是否直接访问 Database？（违反分层）
  • Service 层是否调用 Controller？（逆向依赖）
  • Model 层是否包含业务逻辑？（职责混乱）
```

#### 设计模式一致性
```
检查点:
  • 单例模式实现是否正确？
  • 工厂模式是否符合项目约定？
  • 观察者模式是否正确解耦？
```

#### 命名规范检查
```
检查点:
  • 包名/模块名是否符合约定？
  • 类名/函数名是否语义清晰？
  • 常量/变量命名是否规范？
```

### 4. 潜在副作用检测

识别代码变更可能产生的副作用。

#### 数据流分析
```
变更: 修改了 User.email 字段的验证逻辑
  ↓
潜在影响:
  • 现有用户数据可能不符合新验证规则
  • 需要数据迁移脚本
  • API 响应格式可能变化
  • 前端表单验证需要同步更新
```

#### 并发安全检查
```
检查点:
  • 共享变量访问是否加锁？
  • 是否存在竞态条件？
  • 是否正确使用了同步机制？
```

#### 性能影响分析
```
检查点:
  • 新增的循环是否会导致性能问题？
  • 数据库查询是否优化？
  • 是否引入了 N+1 查询问题？
  • 是否存在内存泄漏风险？
```

## 分析输出格式

```json
{
  "analysis_mode": "deep",
  "timestamp": "2025-12-22T10:30:00Z",
  "changed_files": [
    "src/service/order.go",
    "src/model/user.py"
  ],

  "dependency_analysis": {
    "direct_dependencies": [
      {
        "file": "src/service/order.go",
        "depends_on": [
          "src/model/order.go",
          "src/repository/order_repo.go"
        ]
      }
    ],
    "reverse_dependencies": [
      {
        "file": "src/service/order.go",
        "used_by": [
          "src/controller/order_controller.go",
          "src/api/checkout.go"
        ]
      }
    ],
    "circular_dependencies": []
  },

  "call_chain_analysis": {
    "upstream_calls": [
      {
        "function": "processOrder",
        "called_by": [
          {
            "function": "handleCheckout",
            "file": "src/controller/checkout.go",
            "line": 45
          }
        ]
      }
    ],
    "downstream_calls": [
      {
        "function": "processOrder",
        "calls": [
          {
            "function": "validateOrder",
            "file": "src/service/validator.go",
            "type": "internal"
          },
          {
            "function": "Database.Update",
            "type": "external",
            "risk": "high"
          }
        ]
      }
    ]
  },

  "architecture_consistency": {
    "layer_violations": [
      {
        "severity": "error",
        "file": "src/controller/user.go",
        "line": 78,
        "issue": "Controller directly accesses Database layer",
        "suggestion": "Use Service layer instead"
      }
    ],
    "pattern_violations": [],
    "naming_violations": [
      {
        "severity": "warning",
        "file": "src/model/user.py",
        "line": 12,
        "issue": "Function name 'getData' not following snake_case convention",
        "suggestion": "Rename to 'get_data'"
      }
    ]
  },

  "side_effect_detection": {
    "data_migration_needed": true,
    "api_breaking_changes": [
      {
        "endpoint": "/api/user/profile",
        "field": "email",
        "change": "Validation rule changed",
        "impact": "Existing clients may fail validation"
      }
    ],
    "concurrency_issues": [],
    "performance_risks": [
      {
        "severity": "warning",
        "file": "src/service/order.go",
        "line": 123,
        "issue": "Loop contains database query (potential N+1 problem)",
        "suggestion": "Use batch query or join"
      }
    ]
  },

  "summary": {
    "total_issues": 5,
    "critical": 1,
    "high": 2,
    "medium": 2,
    "low": 0,
    "files_analyzed": 15,
    "dependencies_checked": 42
  }
}
```

## 分析策略

### 快速模式 vs 深度模式对比

| 分析维度 | 快速模式 | 深度模式 |
|---------|---------|---------|
| **代码范围** | 仅变更行 | 变更行 + 依赖代码 |
| **依赖分析** | 不分析 | 完整依赖树 |
| **调用链** | 不追踪 | 双向追踪 |
| **架构检查** | 基础检查 | 全面检查 |
| **副作用** | 不检测 | 深度检测 |
| **执行时间** | < 1分钟 | 5-10分钟 |

### 分析深度配置

可以通过配置文件调整分析深度：

```yaml
# deep-analysis-config.yml
analysis:
  dependency:
    max_depth: 3              # 依赖分析最大深度
    include_test_files: false # 是否包含测试文件

  call_chain:
    max_upstream_depth: 5     # 向上追踪最大深度
    max_downstream_depth: 5   # 向下追踪最大深度
    include_stdlib: false     # 是否包含标准库调用

  architecture:
    check_layering: true      # 分层检查
    check_patterns: true      # 设计模式检查
    check_naming: true        # 命名规范检查

  performance:
    check_loops: true         # 检查循环性能
    check_queries: true       # 检查数据库查询
    check_memory: true        # 检查内存使用
```

## 技术实现

### Go 语言分析工具
```bash
# 依赖分析
go list -m all
go mod graph

# 调用分析
go build -toolexec=callgraph

# 静态分析
staticcheck
go vet
```

### Python 语言分析工具
```bash
# 依赖分析
pipdeptree
pydeps

# 调用分析
pyan (Python Call Graph)

# 静态分析
mypy --strict
bandit (security)
```

### Java 语言分析工具
```bash
# 依赖分析
jdeps
mvn dependency:tree

# 调用分析
JDepend
Classycle

# 静态分析
PMD
FindBugs/SpotBugs
```

## 执行流程

```
开始
  │
  ├─→ 1. 接收输入
  │     ├─→ 变更文件列表
  │     ├─→ 项目语言类型
  │     ├─→ 分析配置
  │     └─→ 基础检查结果（来自 tool-runner）
  │
  ├─→ 2. 依赖关系分析
  │     ├─→ 构建依赖图
  │     ├─→ 识别直接依赖
  │     ├─→ 识别间接依赖
  │     ├─→ 检测循环依赖
  │     └─→ 标记受影响模块
  │
  ├─→ 3. 调用链追踪
  │     ├─→ 向上追踪调用者
  │     │   ├─→ 直接调用者
  │     │   └─→ 间接调用者（递归）
  │     │
  │     └─→ 向下追踪被调用
  │         ├─→ 内部函数调用
  │         ├─→ 外部 API 调用
  │         └─→ 系统调用（DB/IO/Network）
  │
  ├─→ 4. 架构一致性检查
  │     ├─→ 分层架构验证
  │     │   ├─→ Controller → Service → Repository 顺序
  │     │   └─→ 禁止逆向依赖
  │     │
  │     ├─→ 设计模式检查
  │     │   ├─→ 单例模式实现
  │     │   ├─→ 工厂模式规范
  │     │   └─→ 依赖注入使用
  │     │
  │     └─→ 命名规范检查
  │         ├─→ 包/模块命名
  │         ├─→ 类/函数命名
  │         └─→ 变量/常量命名
  │
  ├─→ 5. 潜在副作用检测
  │     ├─→ 数据流分析
  │     │   ├─→ 数据模型变更影响
  │     │   ├─→ API 兼容性分析
  │     │   └─→ 数据迁移需求
  │     │
  │     ├─→ 并发安全检查
  │     │   ├─→ 竞态条件检测
  │     │   ├─→ 死锁风险分析
  │     │   └─→ 线程安全验证
  │     │
  │     └─→ 性能影响分析
  │         ├─→ 算法复杂度
  │         ├─→ 数据库查询优化
  │         └─→ 内存使用评估
  │
  ├─→ 6. 整合分析结果
  │     ├─→ 合并所有检测结果
  │     ├─→ 问题优先级排序
  │     ├─→ 生成影响范围报告
  │     └─→ 提供修复建议
  │
  └─→ 7. 输出结果
        ├─→ JSON 格式输出
        └─→ 传递给 report-generator
```

## 与其他 Skills 的协作

```
tool-runner (基础检查结果)
      │
      └─→ deep-analysis-agent
            │
            ├─→ 读取项目结构
            ├─→ 分析依赖关系
            ├─→ 追踪调用链
            ├─→ 检查架构一致性
            ├─→ 检测副作用
            │
            └─→ 输出深度分析结果
                  │
                  └─→ report-generator (生成深度报告)
                        │
                        └─→ report-validator (验证准确性)
```

## 使用场景

### 场景 1: 重构核心业务逻辑
```
变更: 重构订单处理流程
  ↓
深度分析:
  • 识别所有调用订单处理的模块
  • 分析对库存、支付、通知系统的影响
  • 检查是否破坏了现有API契约
  • 评估性能变化
  ↓
输出: 完整的影响范围报告和风险评估
```

### 场景 2: 数据模型变更
```
变更: 修改 User 表结构，添加新字段
  ↓
深度分析:
  • 识别所有访问 User 模型的代码
  • 检查 ORM 映射是否需要更新
  • 分析现有数据迁移需求
  • 检查 API 响应格式变化
  ↓
输出: 数据迁移脚本建议和 API 变更清单
```

### 场景 3: 架构调整
```
变更: 引入新的缓存层
  ↓
深度分析:
  • 检查缓存层是否符合分层架构
  • 分析缓存失效对系统的影响
  • 检查并发访问的线程安全性
  • 评估性能提升和内存开销
  ↓
输出: 架构一致性报告和性能评估
```

## 扩展指南

### 添加新的分析维度

1. 在配置文件中定义新的分析选项
2. 实现对应的分析逻辑
3. 更新输出 JSON 格式
4. 更新文档

### 集成新的分析工具

1. 在技术实现部分添加工具说明
2. 实现工具调用接口
3. 实现结果解析逻辑
4. 更新配置文件

## 限制和注意事项

1. **执行时间**：深度分析耗时较长，不适合快速反馈场景
2. **依赖完整性**：需要完整的项目代码才能准确分析
3. **动态特性**：无法完全分析动态语言的运行时行为
4. **外部依赖**：对第三方库的分析能力有限
