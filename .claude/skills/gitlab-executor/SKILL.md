---
name: gitlab-executor
version: 1.0
type: skill
description: 封装GitLab MCP工具的操作能力
author: 15085
tags: [gitlab, mcp, api]
dependencies:
  - mcp-gitlab-server
resources:
  scripts: []
  templates: []
---

# gitlab-executor

封装 GitLab MCP 工具,提供统一的 GitLab 操作接口。

## 功能说明

提供以下 GitLab 操作能力:

| 功能分类 | 操作 | MCP工具 |
|----------|------|---------|
| 分支管理 | 创建分支 | `mcp__gitlab__create_branch` |
| 文件操作 | 提交文件 | `mcp__gitlab__create_or_update_file` |
| 文件操作 | 批量推送 | `mcp__gitlab__push_files` |
| MR管理 | 创建MR | `mcp__gitlab__create_merge_request` |
| MR管理 | 获取MR详情 | `mcp__gitlab__get_merge_request` |
| MR管理 | 获取MR变更 | `mcp__gitlab__get_merge_request_diffs` |
| MR管理 | 合并MR | `mcp__gitlab__merge_merge_request` |
| 评论管理 | 创建评论 | `mcp__gitlab__create_note` |
| 评论管理 | 创建讨论线程 | `mcp__gitlab__create_merge_request_thread` |
| 评论管理 | 创建草稿评论 | `mcp__gitlab__create_draft_note` |
| 评论管理 | 发布草稿 | `mcp__gitlab__bulk_publish_draft_notes` |
| 提交管理 | 列出提交 | `mcp__gitlab__list_commits` |
| 提交管理 | 获取提交详情 | `mcp__gitlab__get_commit` |
| 提交管理 | 获取提交差异 | `mcp__gitlab__get_commit_diff` |

## 操作指南

### 1. 创建分支

```yaml
调用: mcp__gitlab__create_branch
参数:
  - branch: 新分支名称
  - ref: 源分支或commit SHA (可选,默认为默认分支)
```

### 2. 创建 Merge Request

```yaml
调用: mcp__gitlab__create_merge_request
参数:
  - title: MR标题
  - source_branch: 源分支
  - target_branch: 目标分支
  - description: MR描述 (可选)
```

### 3. 获取 MR 变更

```yaml
调用: mcp__gitlab__get_merge_request_diffs
参数:
  - merge_request_iid: MR的IID
```

### 4. 添加代码评论

```yaml
调用: mcp__gitlab__create_merge_request_thread
参数:
  - merge_request_iid: MR的IID
  - body: 评论内容
  - position: 代码位置信息 (可选)
```

### 5. 合并 MR

```yaml
调用: mcp__gitlab__merge_merge_request
参数:
  - merge_request_iid: MR的IID
  - should_remove_source_branch: 是否删除源分支 (可选)
```

## 常见错误

| 错误 | 原因 | 解决方案 |
|------|------|----------|
| 401 Unauthorized | Token无效或过期 | 检查GITLAB_TOKEN配置 |
| 404 Not Found | 项目或资源不存在 | 检查project_id和资源ID |
| 403 Forbidden | 无权限操作 | 检查Token权限范围 |
| 409 Conflict | 资源冲突(如分支已存在) | 使用不同名称或先删除 |

## 变更日志

| 版本 | 日期 | 变更内容 | 作者 |
|------|------|----------|------|
| 1.0.0 | 2025-12-23 | 初始版本 | 15085 |
