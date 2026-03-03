---
allowed-tools: Read,Glob,Grep,Edit,Write,Task
description: 基于测试用例文档，生成符合规范的接口测试代码
skills: test-code-generator,test-validator
version: 2.0
---

# 接口测试代码生成器

## 任务描述

基于TDD方法论，通过阅读测试用例文档（testcase.md），自动生成符合测试用例要求的接口测试代码。

> **前置条件**: 需要先执行 `/generate-interface-testcases` 生成测试用例文档

### 全局约束

1. **依赖前置产物**: 必须基于 testcase.md 生成代码
2. **用例一一对应**: testcase.md 中的每个测试用例必须有对应的测试代码
3. **测试框架**:
   - Python HTTP API: pytest + TestClient
   - Go gRPC: go test + 直接调用Handler
4. **代码可执行**: 生成的测试代码必须能直接运行
5. **代码规范**: 严格遵循项目现有代码风格

### 处理步骤

#### 步骤1: 测试代码生成
调用SKILL /test-code-generator
- 输入: doc/testcase/testcase.md
- 处理:
  - 阶段一: 测试数据设计（Mock数据、测试参数）
  - 阶段二: 测试代码生成（根据项目语言选择模板）
- 输出: 完整测试代码文件 + 测试数据文件
- 模板:
  - Python: @skills/test-code-generator/templates/universal_test_template.md
  - Go: @skills/test-code-generator/templates/grpc_test_template.md

#### 步骤2: 质量检视
调用SKILL /test-validator
- 输入: 生成的测试代码
- 输出: 质量检视报告、改进建议
- 方式: AI静态检视（TDD阶段接口未实现）

## 工具说明

### 必需工具
- Read: 读取测试用例文档、代码模板
- Write: 生成测试代码文件
- Glob/Grep: 查找相关文件和代码结构

### 可选工具
- Edit: 修改现有模板文件
- Task: 协调技能执行流程

## 输入说明

### 知识库索引
- 测试用例文档: doc/testcase/testcase.md（由 generate-interface-testcases 生成）
- 需求设计文档: doc/design/tech_design.md（参考）
- 测试代码模板: @skills/test-code-generator/templates/
- 现有测试代码: testing/

## 输出模板

### Python HTTP API测试
```
testing/
├── test_router/
│   └── test_{module}_router_BVT.py
└── data_template/
    └── mock_{module}_data.py
```

### Go gRPC测试
```
internal/service_process/{Service}/
└── {service}_test.go
```

### 测试文件结构
```
1. 文件头部注释
2. 导入语句
3. 测试数据/Mock数据
4. 测试类/函数
5. 正常场景测试 (@pytest.mark.BVT / 无后缀)
6. 异常场景测试 (@pytest.mark.Level1 / _InvalidInput)
7. 边界场景测试 (@pytest.mark.Level1 / _Boundary)
8. 安全场景测试 (@pytest.mark.Level2)
```

## 补充输入(动态补充)

### 执行时需要确认的信息
1. 测试用例文档路径（默认: doc/testcase/testcase.md）
2. 自动识别项目语言后确认（Python/Go）
3. 测试代码输出目录确认
4. 模块名称

## 成功标准

- [ ] 用例100%代码化（testcase.md中的每个用例都有对应测试代码）
- [ ] 测试代码可直接执行
- [ ] 代码风格符合项目规范
- [ ] Mock数据真实有效
- [ ] 测试数据隔离良好

## 注意事项

1. **增量开发**: 如果模块已存在，仅添加新测试
2. **命名规范**: 遵循项目现有命名规范
3. **依赖导入**: 确保正确导入所有依赖
4. **测试隔离**: 每个测试函数独立运行

## 完整工作流

```
概设文档 → [/generate-interface-testcases] → testcase.md → [/generate-interface-tests] → 测试代码
```

## 版本信息

- 版本: v2.0
- 适用项目: DSP数据安全平台
- 维护者: AI Native Dev Team
