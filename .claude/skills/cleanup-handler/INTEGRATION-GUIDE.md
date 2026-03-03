# cleanup-handler 集成指南

## 概述

`cleanup-handler` skill 负责在 lint 和 fix 命令执行完成后清理临时文件，确保工作区保持干净。

## 集成位置

### 1. lint 命令集成

**位置**: `/commands/lint.md`

**快速模式流程图中添加步骤9**（在步骤8之后）：

```markdown
[步骤8] 汇总输出结果
    ├─ 显示检查报告路径
    ├─ 显示修复报告路径
    ├─ 显示修复统计
    └─ 显示备份信息
    ↓
[步骤9] 清理临时文件 (cleanup-handler)
    ├─ 删除 .claude/temp/
    ├─ 删除 .bak 文件
    └─ 显示清理结果
    ↓
✅ 完成 - 检查、修复和清理全部完成
```

**深度模式流程图中添加步骤9**（在步骤8之后）：

```markdown
[步骤8] 输出报告并提示用户
    ├─ doc/lint/lint-deep-{language}-{date}.md
    └─ 提示：请查看报告
    ↓
[步骤9] 清理临时文件 (cleanup-handler)
    ├─ 删除 .claude/temp/
    ├─ 删除 .bak 文件
    └─ 显示清理结果
    ↓
⏸️ 等待用户确认
```

**调用示例**：

```markdown
# 在步骤8完成后调用
Task(
  subagent_type="general-purpose",
  description="清理lint临时文件",
  prompt="请使用 cleanup-handler skill 清理临时文件：

命令: lint
状态: success

清理范围：
- .claude/temp/
- .claude/**/*.bak

请执行清理并显示结果。"
)
```

### 2. fix 命令集成

**位置**: `/commands/fix.md`

**更新步骤7**（原有的清理步骤）：

```markdown
[步骤7] 清理临时文件 (cleanup-handler)
    ├─ 如果修复成功：删除 .backup/
    ├─ 删除 .claude/temp/
    ├─ 删除 .bak 文件
    └─ 显示清理结果
    ↓
[步骤8] 输出修复摘要
```

**调用示例（修复成功）**：

```markdown
Task(
  subagent_type="general-purpose",
  description="清理fix临时文件",
  prompt="请使用 cleanup-handler skill 清理临时文件：

命令: fix
状态: success

清理范围：
- .backup/ (修复成功，删除备份)
- .claude/temp/
- .claude/**/*.bak

请执行清理并显示结果。"
)
```

**调用示例（修复失败）**：

```markdown
Task(
  subagent_type="general-purpose",
  description="清理fix临时文件",
  prompt="请使用 cleanup-handler skill 清理临时文件：

命令: fix
状态: failed

清理范围：
- .claude/temp/
- .claude/**/*.bak

保留：
- .backup/ (修复失败，保留备份用于回滚)

请执行清理并显示结果，提示用户备份位置。"
)
```

### 3. 在 skills 列表中引用

确保 cleanup-handler 已添加到命令的 skills 列表：

```yaml
# lint.md
skills: ...,cleanup-handler

# fix.md
skills: ...,cleanup-handler
```

## 执行流程

### Lint 命令清理流程

```
步骤8完成
  ↓
调用 cleanup-handler
  ├─ 检测到 command=lint
  ├─ 使用标准清理策略
  ├─ rm -rf .claude/temp/
  ├─ find .claude/ -name "*.bak" -delete
  └─ 输出：🧹 临时文件已清理
  ↓
完成
```

### Fix 命令清理流程（成功）

```
步骤6完成（生成报告）
  ↓
调用 cleanup-handler
  ├─ 检测到 command=fix, status=success
  ├─ 使用完整清理策略
  ├─ rm -rf .backup/
  ├─ rm -rf .claude/temp/
  ├─ find .claude/ -name "*.bak" -delete
  └─ 输出：🧹 临时文件已清理
  ↓
步骤8：输出摘要
```

### Fix 命令清理流程（失败）

```
步骤6完成（生成报告）
  ↓
调用 cleanup-handler
  ├─ 检测到 command=fix, status=failed
  ├─ 使用部分清理策略
  ├─ 保留 .backup/ 目录
  ├─ rm -rf .claude/temp/
  ├─ find .claude/ -name "*.bak" -delete
  └─ 输出：🧹 临时文件已清理
      ⚠️ 备份已保留在 .backup/ 目录
  ↓
步骤8：输出摘要
```

## 清理策略总结

| 命令 | 状态 | .backup/ | .claude/temp/ | *.bak |
|------|------|----------|---------------|-------|
| lint | success | N/A | ✅ 删除 | ✅ 删除 |
| fix  | success | ✅ 删除 | ✅ 删除 | ✅ 删除 |
| fix  | failed  | ⚠️ 保留 | ✅ 删除 | ✅ 删除 |

## 预期输出

### 成功清理

```
🧹 临时文件清理完成
  ✅ 已删除: .claude/temp/
  ✅ 已删除: .bak 文件 (0个)
  [如果是fix成功] ✅ 已删除: .backup/
```

### 保留备份（修复失败时）

```
🧹 临时文件清理完成
  ✅ 已删除: .claude/temp/
  ✅ 已删除: .bak 文件
  ⚠️ 已保留备份: .backup/ (修复失败，可用于回滚)

💡 提示：使用以下命令回滚：
   cp -r .backup/[timestamp]/* ./
```

## 实施检查清单

- [ ] 更新 lint.md 快速模式流程图（添加步骤9）
- [ ] 更新 lint.md 深度模式流程图（添加步骤9）
- [ ] 更新 lint.md 执行步骤详解（添加步骤9说明）
- [ ] 更新 fix.md 流程图（修改步骤7调用cleanup-handler）
- [ ] 更新 fix.md 执行步骤详解（步骤7改为调用cleanup-handler）
- [ ] 确认 lint.md 的 skills 列表包含 cleanup-handler
- [ ] 确认 fix.md 的 skills 列表包含 cleanup-handler
- [ ] 测试 lint 快速模式是否正确清理
- [ ] 测试 lint 深度模式是否正确清理
- [ ] 测试 fix 成功情况是否正确清理（包括 .backup/）
- [ ] 测试 fix 失败情况是否正确清理（保留 .backup/）

## 注意事项

1. **清理时机**：必须在报告生成后、输出摘要前执行
2. **错误处理**：清理失败不应中断主流程，记录错误但允许继续
3. **Windows兼容性**：使用 `rmdir /s /q` 代替 `rm -rf`
4. **验证**：清理后验证关键文件（如报告）未被误删

---

**更新时间**: 2025-12-26
