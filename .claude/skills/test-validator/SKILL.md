---
name: test-validator
description: 通过AI静态检视验证测试代码是否符合testcase设计意图，确保代码规范和测试覆盖率
version: 2.1
---

# 测试代码检视技能

## 概述

通过AI静态检视验证生成的测试代码质量，确保：
- 测试代码完全符合 testcase 的设计意图
- 测试场景覆盖完整（正常、异常、边界、安全）
- 代码符合项目规范和命名规范
- 断言逻辑正确完整

> **注意**：本技能不执行代码，仅进行静态代码检视，适用于TDD流程中接口未实现的场景。

## 强制约束

1. **用例对应检查**: 每个 testcase 必须有对应的测试代码
2. **验收条件覆盖**: testcase 的验收条件必须转化为断言
3. **输出检视报告**: 必须生成结构化的检视报告

## 实现方法

### 步骤1: 准备检视材料
```markdown
1. 读取 testcase 文档（doc/testcase/testcase.md）
2. 读取生成的测试代码文件
3. 读取项目测试代码模板（参考）
```

### 步骤2: Testcase 符合度检视
```markdown
检查项：
- 每个 testcase 是否都有对应的测试方法
- 测试数据是否与 testcase 一致
- 验收条件是否都转化为断言
```

### 步骤3: 代码规范检视
```markdown
检查项：
- 文件名：test_{module}_router_BVT.py
- 类名：Test{ModuleName}Router
- 方法名：test_tc_{module}_{function}_{sequence}_{description}
- setup/teardown 是否完整
- docstring 是否规范
```

### 步骤4: 测试覆盖度检视
```markdown
检查项：
- 正常场景：基本成功流程、参数组合
- 异常场景：参数校验失败、业务逻辑异常
- 边界场景：数值边界、字符串边界
- 安全场景：权限验证、参数注入
```

### 步骤5: 断言逻辑检视
```markdown
检查项：
- HTTP 状态码断言
- 响应数据结构断言
- 业务逻辑结果断言
- 错误信息断言（异常场景）
```

### 步骤6: 生成检视报告

## 输出模板

### 检视报告 (JSON)
```json
{
  "validation_info": {
    "generated_time": "2025-01-XX",
    "validation_method": "AI静态代码检视"
  },
  "testcase_compliance": {
    "total_testcases": 40,
    "implemented_testcases": 40,
    "compliance_rate": "100%",
    "missing_testcases": []
  },
  "coverage_analysis": {
    "scenario_coverage": {
      "normal": "100%",
      "exception": "95%",
      "boundary": "90%",
      "security": "100%"
    },
    "priority_coverage": {
      "P0": "100%",
      "P1": "95%",
      "P2": "85%"
    }
  },
  "quality_scores": {
    "code_quality": {"score": 28, "max": 30},
    "coverage": {"score": 33, "max": 35},
    "assertion": {"score": 18, "max": 20},
    "test_data": {"score": 13, "max": 15}
  },
  "overall_score": 92,
  "validation_result": "通过",
  "recommendations": [],
  "critical_issues": []
}
```

### 评分标准

| 维度 | 满分 | 优秀 | 良好 | 可接受 |
|------|------|------|------|--------|
| 代码质量 | 30 | 27+ | 24+ | 21+ |
| 测试覆盖 | 35 | 33+ | 30+ | 26+ |
| 断言质量 | 20 | 18+ | 16+ | 14+ |
| 测试数据 | 15 | 14+ | 12+ | 11+ |

**总分评价**：90+ 优秀 | 80+ 良好 | 70+ 可接受 | <70 需改进

## 常见错误

### 错误1: 用例遗漏
```text
❌ 错误：testcase 有12个用例，代码只有8个测试函数
✅ 正确：每个 testcase 都有对应的测试函数
```

### 错误2: 断言不完整
```text
❌ 错误：只验证状态码200
✅ 正确：验证状态码 + 响应结构 + 关键字段值 + 错误信息
```

### 错误3: docstring 缺失
```text
❌ 错误：def test_001(self): ...
✅ 正确：包含测试目标、前置条件、测试步骤、期望结果的 docstring
```

## 检视示例

### 好的实践
```python
def test_tc_risk_list_query_001_normal(self):
    """
    测试目标：风险告警列表正常查询
    前置条件：测试数据已初始化
    期望结果：HTTP 200，返回正确数据结构
    """
    response = client.post("/api/risk_list", json={"limit": 20})
    assert response.status_code == 200, f"期望200，实际{response.status_code}"
    assert "total" in response.json(), "响应应包含total字段"
```

### 不好的实践
```python
def test_001(self):  # 命名不规范，无docstring
    response = client.post("/api/risk_list", json={"limit": 20})
    assert response.status_code == 200  # 断言不完整，无错误提示
```

## 输入输出

### 输入
```json
{
  "testcase_document": "doc/testcase/testcase.md",
  "test_files": ["testing/test_router/test_xxx_router_BVT.py"]
}
```

### 输出
- 结构化检视报告（JSON格式）
- 问题清单和改进建议

## 注意事项

1. **客观公正**：基于事实进行检视，不带主观偏见
2. **全面细致**：关注代码的各个方面，不遗漏关键点
3. **建设性意见**：提供具体的、可执行的改进建议
4. **局限性说明**：静态检视无法验证运行时行为，需后续执行验证
