---
allowed-tools: Bash,Read,Glob,Grep,Edit,Write
description: 代码规范检查命令。支持快速/深度模式（--mode）和增量/全量/最新范围（--scope）的灵活组合。
skills: branch-validator,language-detector,standard-loader,tool-runner,external-standards-go,external-standards-python,external-standards-java,internal-standards-go,internal-standards-python,internal-standards-java,code-checker-go,code-checker-python,code-checker-java,security-checker,deep-analysis-agent,report-generator,report-validator,report-corrector,issue-merger,code-fixer-go,code-fixer-python,code-fixer-java,fix-report-generator,cleanup-handler
version: 2.0
---

# 代码规范检查命令 (/lint)


## ⚠️ 强约束规则

**所有代码检查必须100%遵守以下强约束规则：**

### 1. 规范使用约束 ⭐⭐⭐

**内部规范优先原则**：
- ✅ **必须**: 优先使用组织内部编码规范（INTERNAL-*）
- ✅ **允许**: 使用外部规范（Effective Go / PEP 8 / Google Style Guide）作为补充
- ✅ **必须**: 通过 `standard-loader` 加载并合并内部规范和外部规范
- ⚠️ **配置**: 内部规范优先级高于外部规范（priority: 200 vs 100）
- ❌ **禁止**: 完全忽略内部规范，只使用外部规范

**外部规范控制**：
- 可通过 `standard-loader` skill 目录下的 `config.yaml` 中的 `external_standards.enabled` 控制是否启用外部规范
- 默认：`enabled: false`（只使用内部规范）

### 2. 代码范围约束 ⭐⭐⭐

**增量模式（--scope=incremental）**：
- ✅ **必须**: 严格按照 git diff 获取的行号范围检查
- ✅ **必须**: 只检查变更的代码行（新增、修改的行）
- ❌ **禁止**: 检查变更范围之外的任何代码行
- ❌ **禁止**: 检查未变更的上下文代码

**全量模式（--scope=full）**：
- ✅ **必须**: 检查项目中所有代码文件（排除规则除外）
- ✅ **必须**: 遵守 `report-generator` skill 目录下的 `rules/EXCLUDE-RULES.md` 中的排除规则

**最新模式（--scope=latest）**：
- ✅ **必须**: 只检查最近一次提交（HEAD）的变更
- ✅ **必须**: 使用 `git show HEAD` 获取变更范围

### 3. 报告唯一性约束 ⭐⭐⭐

**⚠️ 适用所有模式：快速、深度、增量、全量、最新模式都必须遵守**

**快速模式报告**：
- ✅ **必须**: 生成新报告前删除旧的快速模式报告（**无论哪种范围**）
- ✅ **必须**: 删除模式：`lint-{scope}-fast-{language}-*.md`
- ❌ **禁止**: `doc/lint/` 目录中存在多份同语言、同范围的快速模式报告

**深度模式报告**：
- ✅ **必须**: 生成新报告前删除旧的深度模式报告（**无论哪种范围**）
- ✅ **必须**: 删除模式：`lint-{scope}-deep-{language}-*.md`
- ❌ **禁止**: `doc/lint/` 目录中存在多份同语言、同范围的深度模式报告

**强制要求**：
- ⚠️ **关键**：这是强制要求，不区分快速模式还是深度模式
- ⚠️ **目标**：确保报告目录中只保留当前执行生成的最新报告
- ⚠️ **执行时机**：在调用 report-generator 生成报告之前就必须删除

### 4. 只报告违规代码 ⭐⭐⭐

- ✅ **必须**: 只报告真正违反规范的问题
- ❌ **禁止**: 报告"符合规范"、"保持当前风格"、"代码正确，无需修改"等无问题内容
- ❌ **禁止**: 在 `suggested_code` 中写"保持当前代码即可"
- ✅ **必须**: `report-generator` 在步骤0.1中过滤掉所有"符合规范"的项

### 5. 自动执行约束（快速模式）⭐⭐⭐

**快速模式**：
- ✅ **必须**: 所有步骤自动连续执行，从检查到修复一气呵成
- ✅ **必须**: 自动触发 `/fix` 命令修复所有问题
- ❌ **禁止**: 在步骤之间暂停并询问"是否继续"
- ✅ **必须**: 使用 TodoWrite 工具实时显示进度
- ✅ **允许**: 在需要用户决策时暂停（如：检测到多语言项目）
- ⚠️ **原则**: 用户执行 `/lint` 后，应自动完成检查和修复，直到生成最终报告

**分阶段执行机制**（避免上下文过长）：
- ✅ **推荐**: 将完整流程分为多个独立阶段，使用 Skill 工具链式调用
- ✅ **推荐**: 每个阶段职责单一，通过临时文件（`.claude/temp/`）传递数据
- ✅ **推荐**: 阶段划分示例：
  - 阶段1：准备和检查阶段（步骤1-4）
    - 步骤1: 分支验证、语言检测和参数解析
    - 步骤2: **规范加载并行外部工具检查** - 并行执行 standard-loader 和 tool-runner
    - 步骤3: AI 规范检查 (code-checker-*)
    - 步骤4: 安全编码检查 (security-checker)
    - → 生成检查结果文件到 `.claude/temp/`
  - 阶段2：报告阶段（步骤5-6）
    - 步骤5: 生成初步报告 (report-generator)
    - 步骤6: 报告验证、修正与优化 (report-validator + report-corrector + issue-merger)
    - → 生成验证后的报告到 `doc/lint/`
  - 阶段3：修复阶段（步骤7-9，快速模式自动触发）
    - 步骤7: 自动修复 (fix skill)
    - 步骤8: 汇总输出结果
    - 步骤9: 清理临时文件 (cleanup-handler)
- ✅ **优势**: 每个阶段上下文独立，避免单次执行上下文过长导致超时或性能问题
- ⚠️ **实现**: 主命令只负责参数解析和第一阶段调用，各阶段自动链式调用下一阶段

**深度模式**：
- ✅ **必须**: 检查和分析步骤自动执行
- ✅ **必须**: 生成报告后提示用户确认
- ⏸️ **暂停点**: 等待用户手动执行 `/fix` 命令
- ❌ **禁止**: 深度模式自动触发修复

### 6. 报告验证约束 ⭐⭐⭐

**report-validator 触发规则**：
- ✅ **必须**: 所有生成的报告都必须经过 `report-validator` 验证
- ✅ **必须**: 验证内容：行号准确性、问题在变更范围内、规范引用真实性
- ✅ **必须**: 如果无效率 ≥ 30%，返回 checker 重新生成报告
- ✅ **必须**: 如果无效率 < 30%，直接修正原报告
- ❌ **禁止**: 跳过报告验证步骤

### 7. 零问题处理约束（深度模式）⭐⭐⭐

**当基础检查发现0个问题时**：
- ✅ **必须**: 深度模式强制触发深度分析（`deep-analysis-agent`）
- ✅ **必须**: 分析依赖关系、调用链、架构一致性
- ❌ **禁止**: 直接生成"检查通过，无问题"报告
- ✅ **必须**: 深度分析完成后，根据结果决定是否更新报告
- ⚠️ **检查点**: 在生成报告前检查 `if total_issues == 0 && mode == deep then 执行深度分析`

**快速模式零问题处理**：
- ✅ **允许**: 快速模式可以直接生成"检查通过"报告
- ❌ **禁止**: 快速模式触发深度分析（保持快速）

### 7.5. 禁止架构分析章节约束 ⭐⭐⭐

**⚠️ 适用所有模式：快速模式和深度模式都必须遵守**

**报告内容约束**：
- ❌ **严格禁止**: 在报告中包含任何架构分析相关的独立章节
- ❌ **严格禁止**: 包含以下内容：
  - "依赖分析结果"章节
  - "调用链分析"章节
  - "架构一致性"章节
  - "潜在副作用"章节
  - 任何关于依赖树、调用图、架构评估的系统级分析
- ✅ **必须**: 报告只包含"📋 检测过程追溯"章节（问题级别，非系统级别）

**关键区别**：
- ✅ **允许**: 问题级别的检测过程（如何发现这个问题）
  - 示例：记录某个问题是通过什么规则、什么步骤检测出来的
- ❌ **禁止**: 系统级别的架构分析（整体依赖关系、调用链图、架构评估）
  - 示例：不要生成"依赖分析结果"展示模块间依赖关系树

**执行要求**：
- ⚠️ **强制要求**: 无论快速模式还是深度模式，都不能包含架构分析章节
- ⚠️ **报告章节**: 只允许包含语言规范问题、安全问题、检测过程追溯、整改建议等
- ⚠️ **deep-analysis-agent**: 即使调用了深度分析，也不在报告中展示架构分析结果

### 8. 文件排除约束 ⭐⭐⭐

**必须遵守的排除规则**（详见 `language-detector/rules/EXCLUDE-RULES.md`）：
- ✅ **必须**: 在步骤1.3获取文件列表时应用排除规则
- ✅ **必须**: 遵守 `language-detector/rules/EXCLUDE-RULES.md` 中定义的所有排除规则
- ❌ **禁止**: 检查配置目录、测试文件、第三方依赖、构建产物等
- ⚠️ **执行时机**: 在获取变更文件列表时就应该过滤，不要等到后续步骤

### 9. 修复触发约束 ⭐⭐⭐

**快速模式修复触发**：
- ✅ **必须**: 报告验证通过后，立即自动调用 `/fix` 命令
- ✅ **必须**: 使用 Skill 工具调用：`Skill(skill="fix", args="--level=warning")`
- ✅ **必须**: 修复所有 error 和 warning 级别问题（不含 suggestion）
- ✅ **必须**: 生成修复报告到 `doc/fix/`
- ✅ **必须**: 显示修复统计（成功/失败/跳过）
- ❌ **禁止**: 跳过自动修复步骤
- ❌ **禁止**: 在生成 lint 报告后就停止，必须继续执行修复

**深度模式修复触发**：
- ❌ **禁止**: 自动触发修复
- ✅ **必须**: 输出报告路径和提示信息
- ✅ **必须**: 提示用户："查看报告后执行 `/fix` 进行修复"

**关键执行点（快速模式）**：
```markdown
步骤6（报告验证）完成后：
  ↓
立即执行：Skill(skill="fix", args="--level=warning")
  ↓
等待修复完成并获取修复报告
  ↓
输出汇总信息（检查报告 + 修复报告 + 统计）
```

### 10. 分支验证约束 ⭐⭐⭐

**增量模式分支验证**：
- ✅ **必须**: 执行检查前先验证分支是否存在（使用 `branch-validator` skill）
- ✅ **必须**: 分支名精确匹配（区分大小写）
- ✅ **允许**: 自动纠正大小写错误并显示警告
- ✅ **必须**: 分支不存在时显示建议并立即终止执行
- ❌ **禁止**: 使用不存在的分支继续执行

**全量模式和最新模式**：
- ✅ **允许**: 不需要分支验证（不依赖分支参数）

---

## 快速模式工作流程

快速模式适合日常开发，所有步骤自动执行，无需人工干预。

### 执行架构：分阶段自动链式调用

```
用户执行: /lint [branch]  或  /lint --mode=fast [branch]
    ↓
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
阶段1：准备和检查阶段（上下文独立）
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
[步骤1] 分支验证、语言检测和参数解析
    ├─ 1.1 验证分支名（branch-validator）⭐
    │   ├─ 增量模式：验证分支是否存在且大小写正确
    │   ├─ 自动纠正大小写错误（警告提示）
    │   └─ 分支不存在：显示建议并终止执行
    ├─ 1.2 检测项目语言（language-detector）
    ├─ 1.3 获取变更文件列表（严格限定变更行）⭐ 应用排除规则
    │   ├─ 使用 git diff 获取变更文件列表
    │   ├─ **应用排除规则**：过滤掉不应检查的文件和目录
    │   │   └─ 规则定义：参考 `language-detector/rules/EXCLUDE-RULES.md`
    │   └─ 输出：过滤后的文件列表（只包含源代码文件）
    ├─ 1.4 解析变更行号范围（构建完整上下文数据）
    └─ 1.5 输出上下文到：.claude/temp/lint-context-{timestamp}.json
        ├─ 包含：language, mode, scope, current_branch, base_branch
        ├─ 包含：files (文件路径 + check_lines 行号范围)
        └─ 包含：timestamp（供后续步骤使用相同时间戳）
    ↓ 🚀 自动继续
[步骤2] 规范加载并行外部工具检查 ⭐⭐⭐ 并行执行
    ├─ 分支A：规范加载 (standard-loader)
    │   ├─ 读取上下文：.claude/temp/lint-context-{timestamp}.json
    │   ├─ 2.1 加载外部规范（Effective Go / PEP 8 / Google Style Guide）
    │   ├─ 2.2 加载内部规范（组织自定义编码规范）
    │   ├─ 2.3 合并规范（内部规范优先级 200 > 外部规范优先级 100）
    │   └─ 2.4 输出合并后的规范到：.claude/temp/standards-{timestamp}.json
    │
    └─ 分支B：调用外部 Lint 工具 (tool-runner) ⭐ 独立执行
        ├─ 读取上下文：.claude/temp/lint-context-{timestamp}.json
        ├─ Go: golangci-lint
        ├─ Python: pylint + flake8 + black
        ├─ Java: checkstyle + spotbugs
        └─ 输出结果到：.claude/temp/lint-results-{timestamp}.json
    ↓ 🚀 自动继续
[步骤3] AI 规范检查 (code-checker-*)
    ├─ 读取上下文：.claude/temp/lint-context-{timestamp}.json
    ├─ 读取步骤2-A加载的规范：.claude/temp/standards-{timestamp}.json
    ├─ 基于规范进行代码检查
    └─ 输出结果到：.claude/temp/ai-check-results-{timestamp}.json
    ↓ 🚀 自动继续
[步骤4] 安全编码检查 (security-checker)
    ├─ 读取上下文：.claude/temp/lint-context-{timestamp}.json
    ├─ 16大类安全问题检测
    ├─ 独立于AI规范检查，专注于安全漏洞
    └─ 输出结果到：.claude/temp/security-check-results-{timestamp}.json
    ↓ 自动调用阶段2
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
阶段2：报告阶段（上下文独立）⭐ 基于JSON数据流转
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
[步骤5] 生成初步报告 (report-generator)
    ├─ 读取Lint工具结果：.claude/temp/lint-results-{timestamp}.json
    ├─ 读取AI检查结果：.claude/temp/ai-check-results-{timestamp}.json
    ├─ 读取安全检查结果：.claude/temp/security-check-results-{timestamp}.json
    ├─ 合并所有检查结果
    ├─ **输出1**：保存结构化数据 → .claude/temp/report-data-{timestamp}.json
    └─ **输出2**：生成初步Markdown → doc/lint/lint-{scope}-fast-{language}-{date}-draft.md
    ↓ 🚀 自动继续（通过JSON数据传递）
[步骤6] 报告验证、修正与优化 ⭐ 串行依赖，基于JSON
    ├─ 6.1 验证报告准确性 (report-validator)
    │   ├─ **输入**：读取结构化数据 → .claude/temp/report-data-{timestamp}.json
    │   ├─ 验证行号准确性和代码匹配性
    │   ├─ 验证问题在变更范围内（增量模式）
    │   ├─ 验证规范引用真实性
    │   ├─ 计算无效问题比例
    │   └─ **输出**：验证结果 → .claude/temp/validation-result-{timestamp}.json
    ├─ 6.2 修正报告问题 (report-corrector)
    │   ├─ **输入1**：验证结果 → .claude/temp/validation-result-{timestamp}.json
    │   ├─ **输入2**：原始数据 → .claude/temp/report-data-{timestamp}.json
    │   ├─ 无效率<30%：删除无效问题，更新JSON数据
    │   ├─ 无效率≥30%：返回REGENERATE，终止流程，返回阶段1重新执行
    │   └─ **输出**：修正后数据 → .claude/temp/report-data-corrected-{timestamp}.json
    ├─ 6.3 优化报告内容 (issue-merger)
    │   ├─ **输入**：修正后数据 → .claude/temp/report-data-corrected-{timestamp}.json
    │   ├─ 识别并合并同源问题（相同规范+类别+修复方式）
    │   ├─ 生成批量修复建议
    │   └─ **输出**：合并后数据 → .claude/temp/report-data-merged-{timestamp}.json
    └─ 6.4 生成最终报告
        ├─ **输入**：合并后数据 → .claude/temp/report-data-merged-{timestamp}.json
        ├─ 删除初步报告：doc/lint/lint-{scope}-fast-{language}-{date}-draft.md
        └─ **输出**：最终Markdown报告 → doc/lint/lint-{scope}-fast-{language}-{date}.md
    ↓ 自动调用阶段3（快速模式特有）
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
阶段3：修复阶段（上下文独立，快速模式自动触发）⭐⭐⭐ 基于JSON数据
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
[步骤7] 自动修复 - 调用SKILL /fix
    ├─ ⚠️ 关键：使用 Skill 工具调用 fix skill
    ├─ ⚠️ **数据传递**：传递合并后的JSON数据，避免重新解析Markdown ⭐ 新增
    ├─ 命令：Skill(skill="fix", args="--level=warning --data-file=.claude/temp/report-data-merged-{timestamp}.json")
    ├─ 7.1 读取问题数据：.claude/temp/report-data-merged-{timestamp}.json ⭐ 修改
    ├─ 7.2 创建备份 (git stash)
    ├─ 7.3 修复所有 error 和 warning 级别问题
    ├─ 7.4 验证修复结果（语法检查）
    └─ 7.5 生成修复报告：doc/fix/fix-{scope}-fast-{language}-{date}.md
    ↓ 🚀 自动继续
[步骤8] 汇总输出结果
    ├─ 显示检查报告路径：doc/lint/lint-{scope}-fast-{language}-{date}.md
    ├─ 显示修复报告路径：doc/fix/fix-{scope}-fast-{language}-{date}.md
    ├─ 显示修复统计：成功/失败/跳过
    └─ 显示备份信息（如果创建了备份）
    ↓ 🚀 自动继续
[步骤9] 清理临时文件 (cleanup-handler)
    ├─ 删除上下文数据：.claude/temp/lint-context-{timestamp}.json
    ├─ 删除规范数据：.claude/temp/standards-{timestamp}.json
    ├─ 删除检查结果：.claude/temp/lint-results-{timestamp}.json
    ├─ 删除AI检查结果：.claude/temp/ai-check-results-{timestamp}.json
    ├─ 删除安全检查结果：.claude/temp/security-check-results-{timestamp}.json
    ├─ **删除报告数据链**：⭐ 新增
    │   ├─ .claude/temp/report-data-{timestamp}.json
    │   ├─ .claude/temp/validation-result-{timestamp}.json
    │   ├─ .claude/temp/report-data-corrected-{timestamp}.json
    │   └─ .claude/temp/report-data-merged-{timestamp}.json
    ├─ 删除初步报告：doc/lint/lint-{scope}-fast-{language}-{date}-draft.md
    ├─ 删除 .bak 文件
    └─ 显示清理结果
    ↓
✅ 完成 - 检查、修复和清理全部完成
```

### 快速模式特点

- ✅ 严格限定变更行，不检查未变更代码
- ✅ 集成外部 Lint 工具（golangci-lint / pylint / checkstyle 等）
- ✅ **规范加载**：优先使用内部规范，外部规范作为补充
- ✅ 自动去除误报
- ✅ **自动执行修复，无需人工干预**
- ✅ **分阶段执行，上下文独立**，避免单次执行超时
- ✅ 适合 CI/CD 流程
- ❌ 不进行深度关联分析

## 深度模式工作流程

深度模式适合重大重构或发布前检查，包含深度分析，修复前需人工确认。

### 执行架构：分阶段自动链式调用

```
用户执行: /lint --mode=deep [branch]
    ↓
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
阶段1：准备、检查和深度分析阶段（上下文独立）
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
[步骤1] 分支验证、语言检测和参数解析
    ├─ 1.1 验证分支名（branch-validator）⭐
    │   ├─ 增量模式：验证分支是否存在且大小写正确
    │   ├─ 自动纠正大小写错误（警告提示）
    │   └─ 分支不存在：显示建议并终止执行
    ├─ 1.2 检测项目语言（language-detector）
    ├─ 1.3 获取变更文件 + 关联依赖文件（深度模式特有）⭐ 应用排除规则
    │   ├─ 使用 git diff 获取变更文件列表
    │   ├─ 分析关联依赖文件（深度模式）
    │   ├─ **应用排除规则**：过滤掉不应检查的文件和目录
    │   │   └─ 规则定义：参考 `language-detector/rules/EXCLUDE-RULES.md`
    │   └─ 输出：过滤后的文件列表（只包含源代码文件）
    ├─ 1.4 解析变更行号范围（构建完整上下文数据）
    └─ 1.5 输出上下文到：.claude/temp/lint-context-{timestamp}.json
        ├─ 包含：language, mode, scope, current_branch, base_branch
        ├─ 包含：files (文件路径 + check_lines 行号范围)
        └─ 包含：timestamp（供后续步骤使用相同时间戳）
    ↓ 🚀 自动继续
[步骤2] 规范加载并行外部工具检查 ⭐⭐⭐ 并行执行
    ├─ 分支A：规范加载 (standard-loader)
    │   ├─ 读取上下文：.claude/temp/lint-context-{timestamp}.json
    │   ├─ 2.1 加载外部规范（Effective Go / PEP 8 / Google Style Guide）
    │   ├─ 2.2 加载内部规范（组织自定义编码规范）
    │   ├─ 2.3 合并规范（内部规范优先级 200 > 外部规范优先级 100）
    │   └─ 2.4 输出合并后的规范到：.claude/temp/standards-{timestamp}.json
    │
    └─ 分支B：调用外部 Lint 工具 (tool-runner) ⭐ 独立执行
        ├─ 读取上下文：.claude/temp/lint-context-{timestamp}.json
        ├─ 使用完整参数，更深入检查
        └─ 输出结果到：.claude/temp/lint-results-{timestamp}.json
    ↓ 🚀 自动继续
[步骤3] AI 规范检查 (code-checker-*)
    ├─ 读取上下文：.claude/temp/lint-context-{timestamp}.json
    ├─ 读取步骤2-A加载的规范：.claude/temp/standards-{timestamp}.json
    ├─ 检查变更代码 + 关联代码（深度模式范围更广）
    └─ 输出结果到：.claude/temp/ai-check-results-{timestamp}.json
    ↓ 🚀 自动继续
[步骤4] 安全编码检查 (security-checker)
    ├─ 读取上下文：.claude/temp/lint-context-{timestamp}.json
    ├─ 16大类安全问题检测
    ├─ 独立于AI规范检查，专注于安全漏洞
    └─ 输出结果到：.claude/temp/security-check-results-{timestamp}.json
    ↓ 🚀 自动继续
[步骤5] 深度分析 (deep-analysis-agent) ⭐ 深度模式特有
    ├─ 读取上下文：.claude/temp/lint-context-{timestamp}.json
    ├─ 5.1 依赖关系分析（直接/间接依赖）
    ├─ 5.2 调用链追踪（向上/向下）
    ├─ 5.3 架构一致性检查（分层、设计模式）
    ├─ 5.4 潜在副作用检测（数据流、并发、性能）
    └─ 输出结果到：.claude/temp/deep-analysis-{timestamp}.json
    ↓ 自动调用阶段2
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
阶段2：报告阶段（上下文独立）⭐ 基于JSON数据流转
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
[步骤6] 生成深度分析报告 (report-generator)
    ├─ 读取Lint工具结果：.claude/temp/lint-results-{timestamp}.json
    ├─ 读取AI检查结果：.claude/temp/ai-check-results-{timestamp}.json
    ├─ 读取安全检查结果：.claude/temp/security-check-results-{timestamp}.json
    ├─ 读取深度分析结果：.claude/temp/deep-analysis-{timestamp}.json
    ├─ 合并所有检查结果
    ├─ **输出1**：保存结构化数据 → .claude/temp/report-data-{timestamp}.json
    └─ **输出2**：生成初步Markdown → doc/lint/lint-{scope}-deep-{language}-{date}-draft.md
    ↓ 🚀 自动继续（通过JSON数据传递）
[步骤7] 报告验证、修正与优化 ⭐ 串行依赖，基于JSON
    ├─ 7.1 验证报告准确性 (report-validator)
    │   ├─ **输入**：读取结构化数据 → .claude/temp/report-data-{timestamp}.json
    │   ├─ 验证行号准确性和代码匹配性
    │   ├─ 验证问题在变更范围内
    │   ├─ 验证规范引用真实性
    │   ├─ 计算无效问题比例
    │   └─ **输出**：验证结果 → .claude/temp/validation-result-{timestamp}.json
    ├─ 7.2 修正报告问题 (report-corrector)
    │   ├─ **输入1**：验证结果 → .claude/temp/validation-result-{timestamp}.json
    │   ├─ **输入2**：原始数据 → .claude/temp/report-data-{timestamp}.json
    │   ├─ 无效率<30%：删除无效问题，更新JSON数据
    │   ├─ 无效率≥30%：返回REGENERATE，终止流程，返回阶段1重新执行
    │   └─ **输出**：修正后数据 → .claude/temp/report-data-corrected-{timestamp}.json
    ├─ 7.3 优化报告内容 (issue-merger)
    │   ├─ **输入**：修正后数据 → .claude/temp/report-data-corrected-{timestamp}.json
    │   ├─ 合并同源问题
    │   ├─ 生成批量修复建议
    │   └─ **输出**：合并后数据 → .claude/temp/report-data-merged-{timestamp}.json
    └─ 7.4 生成最终报告
        ├─ **输入**：合并后数据 → .claude/temp/report-data-merged-{timestamp}.json
        ├─ 删除初步报告：doc/lint/lint-{scope}-deep-{language}-{date}-draft.md
        └─ **输出**：最终Markdown报告 → doc/lint/lint-{scope}-deep-{language}-{date}.md
    ↓ 🚀 自动继续
[步骤8] 输出报告并提示用户
    ├─ 显示报告路径：doc/lint/lint-{scope}-deep-{language}-{date}.md
    ├─ 显示检查结果概要：
    │   ├─ 基础问题统计（error/warning/suggestion）
    │   ├─ 依赖影响范围
    │   ├─ 调用链追踪结果
    │   └─ 架构一致性问题
    └─ 提示：请查看报告，确认后手动执行修复
    ↓ 🚀 自动继续
[步骤9] 清理临时文件 (cleanup-handler)
    ├─ 删除上下文数据：.claude/temp/lint-context-{timestamp}.json
    ├─ 删除规范数据：.claude/temp/standards-{timestamp}.json
    ├─ 删除检查结果：.claude/temp/lint-results-{timestamp}.json
    ├─ 删除AI检查结果：.claude/temp/ai-check-results-{timestamp}.json
    ├─ 删除安全检查结果：.claude/temp/security-check-results-{timestamp}.json
    ├─ 删除深度分析结果：.claude/temp/deep-analysis-{timestamp}.json
    ├─ **删除报告数据链**：⭐ 新增
    │   ├─ .claude/temp/report-data-{timestamp}.json
    │   ├─ .claude/temp/validation-result-{timestamp}.json
    │   ├─ .claude/temp/report-data-corrected-{timestamp}.json
    │   └─ .claude/temp/report-data-merged-{timestamp}.json
    ├─ 删除初步报告：doc/lint/lint-{scope}-deep-{language}-{date}-draft.md
    ├─ 删除 .bak 文件
    └─ 显示清理结果
    ↓
⏸️ 等待用户确认（深度模式暂停点）
    ↓
用户手动执行: /fix --report lint-{scope}-deep-{language}-{date}.md
    ↓
✅ 完成修复
```

### 深度模式特点

- ✅ 检查范围：变更行 + 关联代码
- ✅ **规范加载**：优先使用内部规范，外部规范作为补充
- ✅ 依赖分析：识别受影响模块
- ✅ 调用链追踪：双向追踪函数调用
- ✅ 架构检查：分层规范、设计模式一致性
- ✅ 副作用检测：数据流、并发、性能风险
- ✅ **分阶段执行，上下文独立**，避免单次执行超时
- ✅ **需人工确认后修复**
- ✅ 适合重大重构或发布前检查

---

**本命令提供两种模式的代码规范检查，平衡了开发效率和代码质量保障需求。**
