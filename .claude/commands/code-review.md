---
name: code-review
version: 1.0
type: command
description: 审查Merge Request代码,检测bug和规范违规
author: 15085
tags: [code-review, gitlab, quality]
allowed-tools: [Read, Bash, Task]
skills:
  - gitlab-executor
  - code-analyzer
inputs:
  - name: merge_request_url
    type: string
    required: true
    description: Merge Request的URL或IID
  - name: no_publish
    type: boolean
    required: false
    description: 仅生成报告不发布到GitLab
outputs:
  - name: review_result
    type: object
    required: true
    description: 审查结果,包含发现的问题列表
  - name: comment_url
    type: string
    required: false
    description: 发布的评论URL(如果发布)
estimated_duration: 10m
---

# code-review

对给定的 Merge Request 进行代码审查,检测bug、安全问题和项目规范违规。

## 使用前提

1. 已配置GitLab访问权限
2. 项目中存在 `doc/kb/` 知识库目录(可选)
3. MR处于可审查状态(未关闭、非草稿)

## 全局注意事项

- 使用 `mcp__gitlab` 相关的工具与 GitLab 交互（例如,获取 Merge Request,创建评论）。**不要使用 web fetch**。
- **在开始之前创建一个待办事项列表。**
- 必须引用并链接每个问题（例如,如果引用项目知识库文档,请包含指向它的链接）。

## 执行步骤

### 步骤 0: 参数解析

1. 检查参数是否包含 `--no-publish`
2. 设置 `publish_to_gitlab` 标志(默认为 false)

### 步骤 1: MR状态检查

启动一个 **haiku agent** 检查以下情况:
- Merge Request 是否已关闭
- Merge Request 是否为草稿
- Merge Request 是否不需要代码审查（例如自动生成的 MR,明显正确的微小更改）
- 是否已对此 Merge Request 提交过代码审查

如果任何条件为真,停止并退出。

**注意：仍需审查 Claude 生成的 MR。**

### 步骤 2: 加载项目知识库

启动一个 **haiku agent** 返回 `doc/kb` 目录下所有相关项目知识库 md 文件的**路径列表（不是内容）**。

### 步骤 3: 获取代码Diff

启动一个 **haiku agent** 执行以下任务:

#### a. 获取 MR 基本信息
- 调用 `mcp__gitlab__get_merge_request`
- 提取: source_branch, target_branch, head_sha, title, description

#### b. 拉取远程分支(只读操作)
```bash
git fetch origin <source_branch>
git fetch origin <target_branch>
```
- **不更新本地分支引用,直接使用 origin/<branch> 远程引用**

#### c. 智能识别代码文件
- 获取变更文件列表: `git diff --name-only origin/<target_branch>...origin/<source_branch>`
- 对每个文件进行检查:
  * 验证文件在源分支中存在: `git cat-file -e origin/<source_branch>:"$file"`
  * 获取文件内容并判断类型: `git show origin/<source_branch>:"$file" | file -`
  * 统计文件行数,只保留 ≤500 行的文件
- 将符合条件的文件路径保存到数组中
- **注意**: 使用数组和引号正确处理包含空格的文件路径

#### d. 获取代码文件的 diff
- 如果有符合条件的代码文件,获取这些文件的 diff
- 使用远程引用: `git diff origin/<target_branch>...origin/<source_branch> -- <文件列表>`
- **注意**: 正确引用文件路径数组,避免空格导致的错误

#### e. 返回数据给步骤 4 的审查 agents
- MR 元数据（title, description, head_sha, source_branch, target_branch）
- Git diff 内容（只包含识别出的代码文件）
- 代码文件列表（已审查的文件路径）
- 跳过的文件列表（因过大或非代码文件而跳过的文件及原因）

**重要说明**:
- 使用 **Haiku agent**
- 使用 `file` 命令智能识别代码文件
- **如果 diff 总行数超过 10000 行,需进一步过滤到 < 300 行的文件**
- 在返回数据中明确列出所有跳过的文件及原因

### 步骤 4: 并行代码审查

并行启动 4 个 agent:

| Agent | 类型 | 职责 |
|-------|------|------|
| Agent 1 | Sonnet | 项目知识库合规性审查 (第一部分) |
| Agent 2 | Sonnet | 项目知识库合规性审查 (第二部分) |
| Agent 3 | Opus | 明显 Bug 扫描 |
| Agent 4 | Opus | 深度问题查找 (安全、逻辑错误) |

**输入数据**（从步骤 3 接收）:
- MR 元数据（title, description, head_sha）
- Git diff 内容（只包含代码文件）
- 代码文件列表

**Agent 职责详情**:

**Agent 1 + 2: 项目知识库合规性审查（Sonnet agents）**
- 输入: 步骤 3 的 diff
- 任务: 检查变更是否符合 doc/kb 目录下的项目知识库文档规范
- 输出: 项目知识库违规问题列表

**Agent 3: 明显 Bug 扫描（Opus agent）**
- 输入: 步骤 3 的 diff
- 任务: 扫描明显的 bug。只关注 diff 本身,不阅读额外的上下文。只标记重大的 bug；忽略挑剔的小问题和可能的误报。**不要标记如果不看 git diff 之外的上下文就无法验证的问题**
- 输出: Bug 列表

**Agent 4: 深度问题查找（Opus agent）**
- 输入: 步骤 3 的 diff
- 任务: 查找安全问题、逻辑错误等深层问题,只查找已更改代码范围内的问题
- 输出: 深层问题列表

**关键：我们需要的是高信号问题。** 这意味着:
- 会导致运行时行为不正确的客观 bug
- 清晰、明确的违反项目知识库文档的情况,你可以引用被破坏的确切规则

**我们不需要**:
- 主观的担忧或 "建议"
- 项目知识库文档未明确要求的风格偏好
- "可能" 是问题的潜在问题
- 任何需要解释或判断的事情

**如果你不确定一个问题是否真实,请不要标记它。误报会侵蚀信任并浪费审查者的时间。**

**注意**:
- 各 agent 完全独立并行运行
- 只接收 diff,不需要自己调用 git 或 API
- 只关注代码文件的变更
- 应告知每个 agent MR 的标题和描述,提供有关作者意图的上下文

### 步骤 5: 问题验证

对于步骤 4 中 agent 3 和 4 发现的每个问题,启动并行 subagent 来验证该问题。

这些 subagent 应获取 MR 标题和描述以及问题的描述。agent 的工作是审查问题,以高置信度验证所述问题确实是一个问题。

**验证示例**:
- 如果标记了 "变量未定义" 这样的问题,subagent 的工作将是验证代码中确实如此
- 对于项目知识库问题,agent 应验证被违反的项目知识库规则是否适用于此文件且确实被违反

**Agent 分配**:
- **Opus subagents**: 验证 bug 和逻辑问题
- **Sonnet subagents**: 验证项目知识库违规

### 步骤 6: 过滤误报

过滤掉步骤 5 中未验证的任何问题。这一步将为我们的审查提供高信号问题列表。

**在步骤 4 和 5 中评估问题时使用此列表（这些是误报,不要标记）**:
- 预先存在的问题
- 看起来是 bug 但实际上正确的内容
- 高级工程师不会标记的过于挑剔的小问题
- linter 会捕获的问题（不要运行 linter 来验证）
- 一般代码质量问题（例如,缺乏测试覆盖率,一般安全问题）,除非项目知识库文档中明确要求
- 项目知识库文档中提到但在代码中明确静默的问题（例如,通过 lint 忽略注释）

### 步骤 7: 发布审查结果

根据 `publish_to_gitlab` 标志处理:
- **true**: 在 Merge Request 上发表评论
- **false**: 将审查结果以文本格式输出到终端,不发布到 GitLab

**撰写评论时请遵循以下准则**:
1. 保持输出简短
2. 避免使用表情符号
3. 为每个问题链接并引用相关的代码、文件和 URL
4. 引用项目知识库违规时,必须引用相关文档中被违反的确切文本（例如,doc/kb/style.md 指出:"变量名使用 snake_case"）

**对于最终评论,请严格遵循以下格式**:

```markdown
## Code review

Found N issues:

1. <bug 的简要描述> (doc/kb/xxx.md 指出: "<文档中的确切引用>")

<链接到文件和行,包含完整的 sha1 + 行范围以提供上下文,例如 https://gitlab.com/anthropics/claude-code/blob/1d54823877c4de72b2316a64032a54afc404e619/README.md#L13-L17>

2. <bug 的简要描述> (doc/kb/yyy.md 指出: "<文档中的确切引用>")

<链接到文件和行,包含完整的 sha1 + 行范围以提供上下文>

3. <bug 的简要描述> (由于 <文件和代码片段> 导致的 bug)

<链接到文件和行,包含完整的 sha1 + 行范围以提供上下文>

🤖 Generated with [Claude Code](https://claude.ai/code)

<sub>- If this code review was useful, please react with 👍. Otherwise, react with 👎.</sub>
```

**未发现问题时的评论格式**:

```markdown
## Auto code review

No issues found. Checked for bugs and project knowledge base compliance.

🤖 Generated with [Claude Code](https://claude.ai/code)
```

## 链接格式规范

链接代码时,请严格遵循以下格式,否则 Markdown 预览将无法正确渲染:

**示例**: `https://gitlab.com/anthropics/claude-code/blob/c21d3c10bc8e898b7ac1a2d745bdc9bc4e423afe/package.json#L10-L15`

**规则**:
1. **需要完整的 git sha** - 不能使用短 sha
2. **必须提供完整的 sha** - 像 `https://gitlab.com/owner/repo/blob/$(git rev-parse HEAD)/foo/bar` 这样的命令将不起作用,因为你的评论将直接在 Markdown 中渲染
3. **仓库名称必须与你正在审查代码的仓库匹配**
4. **文件名后跟 # 符号**
5. **行范围格式为 L[start]-L[end]**
6. **在你评论的行之前和之后至少提供 1 行上下文**,以该行为中心（例如,如果你评论第 5-6 行,你应该链接到 `L4-L7`）

## 验收

| 验收项 | 验收标准 |
|--------|----------|
| MR状态检查 | 正确识别不可审查的MR |
| Diff获取 | 成功获取代码变更 |
| 问题检测 | 检出真实问题,误报率低 |
| 结果发布 | 评论格式正确,链接有效 |

## 关联资源

- GitLab工具: `mcp__gitlab__get_merge_request`, `mcp__gitlab__create_note`
- 项目知识库: `doc/kb/` 目录下的规范文档

## 变更记录

| 版本 | 日期 | 变更内容 | 作者 |
|------|------|----------|------|
| 1.0.0 | 2025-12-23 | 初始版本,从Claude Code命令迁移 | 15085 |
