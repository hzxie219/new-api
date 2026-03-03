---
name: ai-business-analyst
argument-hint: [input-path（可选，原始需求目录）]
version: 1.0
description: 将任意格式的原始需求转化为结构化AI需求文档（Epic→Feature→Story）
allowed-tools: [Skill, Read, Write, Glob]
skills: [ai-requirement-analyst]
tags: [需求分析, AI需求, 需求转换]
inputs:
  - name: input_path
    type: string
    required: false
    default: doc/ai-requirement-analyst/input/
    description: 原始需求输入目录路径
---

# ai-business-analyst - AI需求分析师

将任意格式的原始需求转化为结构化AI需求文档。

## 任务描述

调用 `ai-requirement-analyst --input-path=$1` 技能，完成原始需求到AI需求的转换工作。

## 使用方式

### 基本用法
```bash
/ai-business-analyst --input-path=$1
```

## 处理流程

1. 确认输入目录中存在原始需求文档
2. 调用 Skill 工具执行 `ai-requirement-analyst --input-path=$1` 技能
3. 验证输出目录中生成的AI需求文档

## 工具说明

- **Skill**: 调用 `ai-requirement-analyst --input-path=$1` 技能执行核心分析任务
- **Read/Write**: 读写输入输出文件
- **Glob**: 扫描文件目录

## 输出说明

生成的AI需求文档位于 `doc/ai-requirement-analyst/output/`，包含：
- Epic → Feature → Story 三级结构
- 用户故事和简化验收条件
- 需要进一步细化的重点标注
