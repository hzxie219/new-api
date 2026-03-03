# 报告生成流程数据传递协议

## 概述

本文档定义了 report-generator、report-validator、report-corrector、issue-merger 之间的数据传递协议,确保数据准确传递。

## 设计原则

1. **单一数据源**: 使用 JSON 文件作为数据传递媒介
2. **格式统一**: 所有中间数据遵循 REPORT-DATA-FORMAT.md 规范
3. **可追溯**: 每个阶段保留中间数据文件
4. **可恢复**: 任何阶段失败可以从中间数据恢复

## 数据流转图

```
code-checker-* 输出
    ↓
[1] report-generator (生成初步报告)
    ├─ 输入: 标准化 JSON 数据
    ├─ 输出1: .claude/temp/report-data-{timestamp}.json (保存原始数据)
    └─ 输出2: doc/lint/lint-xxx-draft.md (初步 Markdown)
    ↓
[2] report-validator (验证报告准确性)
    ├─ 输入: .claude/temp/report-data-{timestamp}.json
    ├─ 验证: 行号、范围、规范引用
    └─ 输出: .claude/temp/validation-result-{timestamp}.json
    ↓
[3] report-corrector (修正无效问题)
    ├─ 输入1: .claude/temp/validation-result-{timestamp}.json (验证结果)
    ├─ 输入2: .claude/temp/report-data-{timestamp}.json (原始数据)
    ├─ 操作:
    │   ├─ 策略A: 无效率<30% → 删除无效问题,更新数据
    │   └─ 策略B: 无效率≥30% → 返回 REGENERATE,终止流程
    └─ 输出: .claude/temp/report-data-corrected-{timestamp}.json
    ↓
[4] issue-merger (合并同源问题)
    ├─ 输入: .claude/temp/report-data-corrected-{timestamp}.json
    ├─ 操作: 识别并合并同源问题
    └─ 输出: .claude/temp/report-data-merged-{timestamp}.json
    ↓
[5] 最终报告生成
    ├─ 输入: .claude/temp/report-data-merged-{timestamp}.json
    ├─ 删除: doc/lint/lint-xxx-draft.md (初步报告)
    └─ 输出: doc/lint/lint-{scope}-{mode}-{language}-{date}.md (最终报告)
```

## 文件命名规范

### 临时数据文件

```
.claude/temp/
├── report-data-{timestamp}.json              # 原始数据
├── validation-result-{timestamp}.json        # 验证结果
├── report-data-corrected-{timestamp}.json    # 修正后数据
└── report-data-merged-{timestamp}.json       # 合并后数据
```

### 报告文件

```
doc/lint/
├── lint-{scope}-{mode}-{language}-{date}-draft.md  # 初步报告(临时)
└── lint-{scope}-{mode}-{language}-{date}.md        # 最终报告
```

## 数据格式定义

### 1. report-data-{timestamp}.json

**格式**: 完全符合 REPORT-DATA-FORMAT.md v2.0

**关键字段**:
```json
{
  "version": "1.0.0",
  "metadata": {
    "timestamp": "20251227-143000",
    "language": "go",
    "mode": "incremental",
    ...
  },
  "files": [
    {
      "path": "src/app/main.go",
      "issues": [
        {
          "id": "E001",
          "level": "error",
          "line_number": 1,
          ...
        }
      ]
    }
  ]
}
```

### 2. validation-result-{timestamp}.json

**格式**:
```json
{
  "report_path": "doc/lint/lint-xxx-draft.md",
  "source_data": ".claude/temp/report-data-{timestamp}.json",
  "validation_time": "2025-12-27 14:35:00",
  "total_issues": 42,
  "valid_issues": 39,
  "invalid_issues": 3,
  "invalid_rate": 0.071,
  "consecutive_invalid": 1,
  "action": "FIX",
  "invalid_list": [
    {
      "issue_id": "E012",
      "file": "src/app/main.go",
      "line": 125,
      "reason": "行号超出文件范围",
      "type": "line_number_invalid"
    }
  ],
  "valid_issue_ids": ["E001", "E002", "E003", ...]
}
```

### 3. report-data-corrected-{timestamp}.json

**格式**: 与 report-data-{timestamp}.json 相同,但移除了无效问题

**新增字段**:
```json
{
  "version": "1.0.0",
  "correction_info": {
    "source_data": ".claude/temp/report-data-{timestamp}.json",
    "validation_result": ".claude/temp/validation-result-{timestamp}.json",
    "removed_issues": ["E012", "W025", "S003"],
    "removed_count": 3,
    "correction_time": "2025-12-27 14:36:00"
  },
  "metadata": {
    "total_issues": 39  // 更新后的总数
  },
  "files": [...]  // 已移除无效问题
}
```

### 4. report-data-merged-{timestamp}.json

**格式**: 与 report-data-corrected-{timestamp}.json 类似,但包含合并信息

**新增字段**:
```json
{
  "version": "1.0.0",
  "merge_info": {
    "source_data": ".claude/temp/report-data-corrected-{timestamp}.json",
    "merged_groups": 5,
    "original_count": 39,
    "merged_count": 34,
    "merge_time": "2025-12-27 14:37:00"
  },
  "files": [
    {
      "path": "src/app/main.go",
      "issues": [
        {
          "id": "E001-GROUP",  // 合并后的问题组
          "is_merged": true,
          "merged_count": 3,
          "original_ids": ["E001", "E002", "E003"],
          "locations": [
            {"file": "src/app/main.go", "line": 1},
            {"file": "src/app/service.go", "line": 5},
            {"file": "src/utils/helper.go", "line": 12}
          ],
          ...
        }
      ]
    }
  ]
}
```

## Skill 实现要求

### report-generator

**职责**:
1. 接收 code-checker 的 JSON 数据
2. **保存数据**: 写入 `.claude/temp/report-data-{timestamp}.json`
3. 生成 Markdown: 写入 `doc/lint/lint-xxx-draft.md`
4. 返回: 文件路径和 timestamp

**输出**:
```json
{
  "draft_report": "doc/lint/lint-incremental-fast-go-20251227-draft.md",
  "data_file": ".claude/temp/report-data-20251227143000.json",
  "timestamp": "20251227143000"
}
```

### report-validator

**职责**:
1. 读取数据: `.claude/temp/report-data-{timestamp}.json`
2. 验证所有问题
3. 保存结果: `.claude/temp/validation-result-{timestamp}.json`
4. 返回: 验证结果和建议操作

**输入**:
```json
{
  "data_file": ".claude/temp/report-data-20251227143000.json",
  "timestamp": "20251227143000"
}
```

**输出**:
```json
{
  "action": "FIX",
  "validation_result": ".claude/temp/validation-result-20251227143000.json",
  "invalid_rate": 0.071,
  "invalid_count": 3
}
```

### report-corrector

**职责**:
1. 读取验证结果: `.claude/temp/validation-result-{timestamp}.json`
2. 读取原始数据: `.claude/temp/report-data-{timestamp}.json`
3. 删除无效问题,更新统计
4. 保存修正数据: `.claude/temp/report-data-corrected-{timestamp}.json`
5. 返回: 修正结果

**输入**:
```json
{
  "validation_result": ".claude/temp/validation-result-20251227143000.json",
  "data_file": ".claude/temp/report-data-20251227143000.json",
  "timestamp": "20251227143000"
}
```

**输出**:
```json
{
  "action": "FIXED",
  "corrected_data": ".claude/temp/report-data-corrected-20251227143000.json",
  "removed_count": 3,
  "remaining_issues": 39
}
```

### issue-merger

**职责**:
1. 读取修正数据: `.claude/temp/report-data-corrected-{timestamp}.json`
2. 识别并合并同源问题
3. 保存合并数据: `.claude/temp/report-data-merged-{timestamp}.json`
4. 返回: 合并结果

**输入**:
```json
{
  "corrected_data": ".claude/temp/report-data-corrected-20251227143000.json",
  "timestamp": "20251227143000"
}
```

**输出**:
```json
{
  "merged_data": ".claude/temp/report-data-merged-20251227143000.json",
  "merged_groups": 5,
  "original_count": 39,
  "final_count": 34
}
```

### 最终报告生成

**职责**:
1. 读取合并数据: `.claude/temp/report-data-merged-{timestamp}.json`
2. 删除初步报告: `doc/lint/lint-xxx-draft.md`
3. 生成最终报告: `doc/lint/lint-{scope}-{mode}-{language}-{date}.md`

**输入**:
```json
{
  "merged_data": ".claude/temp/report-data-merged-20251227143000.json",
  "timestamp": "20251227143000"
}
```

**输出**:
```json
{
  "final_report": "doc/lint/lint-incremental-fast-go-20251227.md",
  "total_issues": 34,
  "merged_groups": 5
}
```

## 错误处理

### 验证失败 (无效率≥30%)

```
report-corrector 返回:
{
  "action": "REGENERATE",
  "reason": "无效问题比例过高: 35%",
  "recommendation": "重新执行步骤1-4"
}

流程处理:
1. 删除所有临时文件
2. 删除 draft 报告
3. 返回到阶段1重新执行
```

### 中间步骤失败

```
任何步骤失败时:
1. 保留所有临时数据文件 (用于调试)
2. 输出详细错误信息
3. 提示用户检查相关文件
```

## 清理策略

### 成功完成

```
cleanup-handler 清理:
✅ 删除: .claude/temp/report-data-*.json
✅ 删除: .claude/temp/validation-result-*.json
✅ 删除: doc/lint/*-draft.md
✅ 保留: doc/lint/lint-{scope}-{mode}-{language}-{date}.md
```

### 失败时

```
保留所有中间文件用于调试:
✅ 保留: .claude/temp/report-data-*.json
✅ 保留: .claude/temp/validation-result-*.json
✅ 保留: doc/lint/*-draft.md
⚠️ 提示: 用户可手动检查这些文件
```

## 优势

1. **准确性保证**:
   - 统一的 JSON 格式,避免 Markdown 解析错误
   - 每个阶段都基于结构化数据操作

2. **可追溯性**:
   - 每个阶段保留中间数据
   - 可以追溯任何修改

3. **可恢复性**:
   - 任何阶段失败可以从中间数据恢复
   - 方便调试和问题定位

4. **性能优化**:
   - 避免重复解析 Markdown
   - JSON 解析更快更准确

5. **易于测试**:
   - 每个 skill 的输入输出都是标准 JSON
   - 可以单独测试每个 skill

## 实施步骤

1. 修改 report-generator: 增加 JSON 数据保存功能
2. 修改 report-validator: 从 JSON 读取数据
3. 修改 report-corrector: 操作 JSON 而非 Markdown
4. 修改 issue-merger: 从 JSON 读取和输出数据
5. 增加最终报告生成步骤
6. 更新 cleanup-handler: 清理临时 JSON 文件

---

**关键原则**: 数据在 JSON 中流转,Markdown 只作为最终展示!
