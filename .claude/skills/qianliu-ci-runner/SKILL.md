---
name: qianliu-ci-runner
version: 1.0
type: skill
description: 封装千流CI MCP工具的操作能力
author: 15085
tags: [ci-cd, qianliu, mcp]
dependencies:
  - mcp-qianliu-ci-server
resources:
  scripts: []
  templates: []
---

# qianliu-ci-runner

封装千流 CI MCP 工具,提供 CI/CD 流水线操作接口。

## 功能说明

提供以下千流 CI 操作能力:

| 功能分类 | 操作 | MCP工具 |
|----------|------|---------|
| 任务管理 | 初始化CI任务 | `mcp__qianliu-ci__init_task` |
| 阶段管理 | 列出可用阶段 | `mcp__qianliu-ci__list_stages` |
| 执行控制 | 运行CI阶段 | `mcp__qianliu-ci__run_task` |
| 日志管理 | 获取执行日志 | `mcp__qianliu-ci__get_logs` |

## 操作指南

### 1. 初始化 CI 任务

```yaml
调用: mcp__qianliu-ci__init_task
参数:
  - repo_url: 仓库HTTPS URL (必需)
  - branch: 分支名称 (必需)
  - service_name: 服务名称 (必需)
  - variables: 环境变量 (可选)

返回:
  - task_id: 任务ID,用于后续操作

注意: 仓库URL必须为HTTPS格式,SSH格式需先转换
```

### 2. 列出可用阶段

```yaml
调用: mcp__qianliu-ci__list_stages
参数:
  - task_id: 任务ID (必需)

返回:
  - stages: 可用阶段列表,包含名称和状态
```

### 3. 运行 CI 阶段

```yaml
调用: mcp__qianliu-ci__run_task
参数:
  - task_id: 任务ID (必需)
  - stage_name: 阶段名称 (必需)
  - wait: 是否等待完成 (可选,默认true)

返回:
  - status: 执行状态 (success/failed/running)
  - duration: 执行耗时
```

### 4. 获取执行日志

```yaml
调用: mcp__qianliu-ci__get_logs
参数:
  - task_id: 任务ID (必需)
  - stage_name: 阶段名称 (必需)
  - include_file: 是否包含日志文件下载链接 (可选)

返回:
  - logs: 日志内容
  - download_url: 日志文件下载URL (如果include_file=true)
```

## 常见错误

| 错误 | 原因 | 解决方案 |
|------|------|----------|
| 任务初始化失败 | 仓库URL格式错误 | 确保使用HTTPS格式URL |
| 阶段不存在 | 服务未配置该阶段 | 先调用list_stages确认可用阶段 |
| 执行超时 | 任务耗时超过限制 | 增加超时时间或检查任务逻辑 |
| 权限不足 | 服务未授权访问仓库 | 检查千流CI服务配置 |

## 变更日志

| 版本 | 日期 | 变更内容 | 作者 |
|------|------|----------|------|
| 1.0.0 | 2025-12-23 | 初始版本 | 15085 |
