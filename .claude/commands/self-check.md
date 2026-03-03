---
name: self-check
version: 1.0
type: command
description: 完整CI/CD工作流,从代码推送到千流CI执行和代码审查
author: 15085
tags: [ci-cd, gitlab, qianliu]
allowed-tools: [Read, Edit, Bash, Task]
skills:
  - gitlab-executor
  - qianliu-ci-runner
  - log-analyzer
inputs:
  - name: branch
    type: string
    required: false
    description: 要推送的分支名,默认使用当前分支
  - name: commit_message
    type: string
    required: false
    description: 提交信息,默认引导用户输入
  - name: service_name
    type: string
    required: false
    description: 千流CI服务名,默认为apolloserver
  - name: target_branch
    type: string
    required: false
    description: MR目标分支,默认自动推断当前分支的基础分支
outputs:
  - name: mr_url
    type: string
    required: false
    description: 创建的Merge Request URL
  - name: ci_result
    type: object
    required: true
    description: CI执行结果摘要
estimated_duration: 15m
---

# self-check

整合GitLab代码管理和千流CI系统的完整自动化流程。

## 使用前提

1. 已配置GitLab访问权限
2. 已配置千流CI MCP服务
3. 本地Git仓库已关联远程仓库

## 执行步骤

### 步骤 0: 参数准备

1. **解析用户参数**:
   - 分支名: `$1` (如果未提供,使用当前分支)
   - 提交消息: `$2` (如果未提供,后续自动生成)
   - 服务名称: 尝试从参数获取,否则自动推断:
     - 获取远程 URL: `git remote get-url origin`
     - 提取仓库名 (例如 `.../apolloserver.git` -> `apolloserver`)
   - 目标分支: 尝试从参数获取,否则自动推断:
     - 首先尝试获取当前分支的上游分支: `git config --get "branch.$(git branch --show-current).merge" | sed 's|refs/heads/||'`
     - 如果失败,获取远程默认分支: `git symbolic-ref refs/remotes/origin/HEAD | sed 's|refs/remotes/origin/||'`
     - 如果仍然失败,询问用户

2. **输出确认信息**:
   - 显示解析后的: 服务名、目标分支、当前分支

### 步骤 1: 智能推送代码

1. 检查状态: `git status --porcelain`
2. **如有修改**:
   - `git add .`
   - **生成提交消息** (如果用户未提供):
     - 查看变更: `git diff --cached`
     - **请根据变更内容生成一个简洁的提交摘要** (例如: "Fix: 修复了空指针异常" 或 "Feat: 添加了支付接口")
   - `git commit -m "<commit_message>"`
   - `git push origin <branch>`
3. **无修改**:
   - 自动跳过推送,直接进入下一步

### 步骤 2: 初始化千流CI

1. 获取仓库URL: `git remote get-url origin`

2. **转换SSH URL为HTTPS** (千流CI仅支持HTTPS):
   ```bash
   # git@code.sangfor.org:ADS/DSP/Analysis/apolloserver.git
   # 转换为:
   # https://code.sangfor.org/ADS/DSP/Analysis/apolloserver.git
   # 转换规则: 将 "git@" 替换为 "https://", 将第一个 ":" 替换为 "/"
   ```

3. 调用: `mcp__qianliu-ci__init_task(repo_url, branch, service_name, variables)`

4. 显示 task_id

### 步骤 3: 选择CI阶段

1. **获取可用阶段**:
   - 调用: `mcp__qianliu-ci__list_stages(task_id)`
   - 解析返回的阶段列表,显示所有可用阶段

2. **显示阶段列表**:
   ```
   发现以下CI阶段:
   1. 编译打包
   2. 单元测试
   3. 集成测试
   4. 代码检查
   5. 覆盖率检查

   执行选项:
   [A] 执行全部阶段
   [1-5] 单独执行指定阶段 (可多选,用逗号分隔,如: 1,3,5)

   请选择执行方式:
   ```

3. **处理用户选择**:
   - 选择A: 执行所有可用阶段
   - 选择数字: 单独执行指定的阶段
   - 生成最终的 `selected_stages` 列表用于步骤4

### 步骤 4: 执行CI阶段

**循环执行逻辑**:
```
for each stage in selected_stages:
  while true:
    执行: mcp__qianliu-ci__run_task(task_id, stage_name, wait=true)

    if 成功:
      显示 ✓ 进入下一阶段
    else:
      获取日志: mcp__qianliu-ci__get_logs(task_id, stage_name)
      询问: [修复后重试] [跳过阶段] [终止流程]

      if 重试:
        询问是否推送新代码
        if 需要: 返回步骤1
        continue
      elif 跳过:
        break
      else:
        终止流程
```

**超时处理**: 30分钟后询问 [继续等待] [转为后台执行]

### 步骤 5a: CI成功 - 自动创建MR并审查

1. 显示CI成功摘要
2. **自动创建MR**:
   - 标题: 使用提交消息或 "CI Passed: <branch>"
   - 描述: 包含 CI 任务 ID 和执行结果摘要
   - 执行: `mcp__gitlab__create_merge_request(project_id, source_branch, target_branch, title, description)`

3. **询问是否进行代码审查**:
   ```
   MR已创建成功!是否进行代码审查?

   选项:
   [T] 生成文本审查报告(默认,不添加到GitLab)
   [Y] 进行审查并添加评论到GitLab
   [N] 跳过审查
   ```

4. **如果用户选择[T] 生成文本报告(默认)**:
   - 执行: `/code-review <mr_url> --no-publish`
   - 该命令会自动执行完整的代码审查流程:
     - 检查MR状态(是否已关闭、是否为草稿等)
     - 加载项目知识库文档(doc/kb/)
     - 多维度审查代码变更(知识库合规性、bug检测、逻辑问题等)
     - 验证发现的问题并过滤误报
   - **但不在GitLab上发布评论**,仅将审查结果以文本格式显示在终端

5. **如果用户选择[Y] 进行审查并添加评论**:
   - 执行: `/code-review <mr_url>`
   - 该命令会自动执行完整的代码审查流程(同步骤4)
   - 在GitLab MR上发布审查评论

6. **如果用户选择[N] 跳过**:
   - 直接输出后续指引:
     - 显示 MR 链接: "MR Created: <mr_url>"
     - 显示 **Ready to Merge**

### 步骤 5b: CI失败 - 智能修复流程

1. **失败信息收集**
   - 自动识别所有失败阶段
   - 显示失败摘要和耗时

2. **智能分析选项**
   ```
   ❌ CI执行失败 - 共有 N 个阶段失败

   失败阶段:
   1. [阶段名称] (耗时: X分Y秒)
      - 错误摘要: [简要错误信息]

   智能分析选项:
   [A] 智能分析并修复所有问题（必需先下载日志）
   [B] 逐个处理失败阶段（必需先下载日志）
   [C] 仅下载完整日志
   [D] 重新运行整个流程
   [Q] 退出流程

   请选择操作:
   ```

3. **自动日志处理**
   - 直接调用 `mcp__qianliu-ci__get_logs(task_id, stage_name, include_file=True)`
   - 获取日志下载URL（格式: `http://<server_ip>:<port>/api/logs/download/<log_id>`）
   - 下载完整日志文件到本地:
     ```bash
     # 创建日志目录
     mkdir -p "./logs/${task_id}_$(date +%Y%m%d_%H%M%S)"

     # 下载日志文件（使用动态获取的下载URL）
     curl -o "${log_dir}/${stage_name}_error_job${job_id}.log" "${download_url}"
     ```

4. **集成日志分析**
   - 调用SKILL /log-analyzer 分析下载的日志文件：
   - 获取结构化分析结果，包括错误原因和修复建议

5. **错误分类与修复建议**
   - 根据分析结果提供具体修复步骤
   - 识别可自动修复的问题：
     - 依赖问题 → 自动执行包安装命令
     - 权限问题 → 提供chmod/chown建议
     - 配置问题 → 显示配置文件修改建议
   - 提供手动修复指导
   - 询问是否应用修复并重试

## 使用示例

```bash
# 完整参数
/self-check feat-payment "添加支付功能" payment-service develop-DSP3.0.31

# 部分参数
/self-check feat-payment "添加支付功能"

# 最简模式
/self-check feat-payment

# 交互模式
/self-check
```

## 验收

| 验收项 | 验收标准 |
|--------|----------|
| 代码推送 | 代码成功推送到指定分支 |
| CI执行 | 所有选定阶段执行完成 |
| MR创建 | CI成功后MR创建成功 |
| 代码审查 | 审查结果正确发布到MR(如选择) |

## 关联资源

### GitLab MCP工具
- `mcp__gitlab__create_branch` - 创建分支
- `mcp__gitlab__push_files` - 推送文件
- `mcp__gitlab__create_merge_request` - 创建MR
- `mcp__gitlab__get_merge_request_diffs` - 获取MR代码变更
- `mcp__gitlab__create_merge_request_thread` - 添加代码行评论
- `mcp__gitlab__create_draft_note` - 创建草稿评论
- `mcp__gitlab__bulk_publish_draft_notes` - 批量发布草稿评论
- `mcp__gitlab__merge_merge_request` - 合并MR

### 千流CI工具
- `mcp__qianliu-ci__init_task` - 初始化任务
- `mcp__qianliu-ci__list_stages` - 列出阶段
- `mcp__qianliu-ci__run_task` - 运行阶段
- `mcp__qianliu-ci__get_logs` - 获取日志

## 变更记录

| 版本 | 日期 | 变更内容 | 作者 |
|------|------|----------|------|
| 1.0.0 | 2025-12-23 | 初始版本,从Claude Code命令迁移 | 15085 |
