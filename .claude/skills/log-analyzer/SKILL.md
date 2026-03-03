---
name: log-analyzer
version: 1.0
type: skill
description: 专业日志分析专家，快速定位错误并提供解决建议
author: 15085
tags: [log, analysis, troubleshooting]
dependencies: [bash, grep, awk]
resources:
  scripts: []
  templates: []
---

# log-analyzer

专业的日志分析专家，用于解析和诊断指定的日志文件，快速定位错误、异常和问题，提供清晰的分析报告和解决建议。适用于应用日志、错误日志、系统日志等各类日志文件的分析。

## 功能说明

### 核心能力

| 能力 | 描述 |
|------|------|
| 日志解析与识别 | 自动识别日志格式（JSON、plaintext、structured logs、syslog 等），解析时间戳、日志级别、消息内容、堆栈跟踪，识别常见日志框架（Log4j、Winston、Logrus、Python logging 等） |
| 问题诊断 | 快速定位错误和异常，分析堆栈跟踪并识别根本原因，识别重复出现的问题模式，提取关键错误信息 |
| 统计分析 | 按日志级别（ERROR、WARN、INFO 等）分类统计，识别最频繁出现的错误，统计总日志条目数 |
| 报告生成 | 生成结构化的分析报告，包含行动建议 |

### 支持的日志格式

| 格式类型 | 示例 |
|----------|------|
| JSON 日志 | `{"timestamp":"2024-01-01T12:00:00Z","level":"ERROR","message":"..."}` |
| 结构化文本 | `2024-01-01 12:00:00 [ERROR] module - message` |
| Python 日志 | `2024-01-01 12:00:00,123 - app.module - ERROR - message` |
| Java 日志 | `2024-01-01 12:00:00.123 ERROR [main] com.app.Module - message` |
| Apache/Nginx | 标准访问日志格式 |

### 常见日志格式示例

**应用程序 JSON 日志**
```json
{"timestamp":"2024-01-01T12:00:00Z","level":"ERROR","message":"Database connection failed","error":"ECONNREFUSED","stack":"..."}
```

**结构化文本日志**
```
2024-01-01 12:00:00.123 [ERROR] app.module - Connection timeout after 30s
```

**Python 日志**
```
2024-01-01 12:00:00,123 - app.module - ERROR - Connection failed
Traceback (most recent call last):
  File "app.py", line 123, in connect
    ...
```

**Java 日志**
```
2024-01-01 12:00:00.123 ERROR [main] com.app.Module - Error processing request
java.lang.NullPointerException: Cannot invoke method on null object
    at com.app.Module.process(Module.java:123)
    ...
```

## 操作指南

### 步骤 1: 日志概览与格式识别

1. 查看日志文件基本信息：
   ```bash
   wc -l <log_file>
   ```

2. 快速预览日志格式：
   ```bash
   head -n 50 <log_file>
   tail -n 50 <log_file>
   ```

3. 识别要点：
   - 检查是否为 JSON 格式
   - 识别时间戳格式
   - 确定日志级别关键词（ERROR、WARN、INFO、DEBUG 等）
   - 识别结构化字段

### 步骤 2: 错误分析

1. 定位所有错误：
   ```bash
   grep -i "error\|exception\|fatal\|critical" <log_file>
   ```

2. 统计各类型错误数量：
   ```bash
   grep -i "error" <log_file> | wc -l
   grep -i "exception" <log_file> | wc -l
   grep -i "fatal" <log_file> | wc -l
   ```

3. 分析错误类型和频率：
   ```bash
   grep -i "error" <log_file> | awk -F'error' '{print $2}' | sort | uniq -c | sort -rn | head -20
   ```

4. 获取错误上下文：
   ```bash
   grep -A 10 -B 5 "具体错误关键词" <log_file>
   ```

### 步骤 3: 日志级别统计

```bash
grep -o '\[ERROR\]\|\[WARN\]\|\[INFO\]\|\[DEBUG\]' <log_file> | sort | uniq -c
```

### 步骤 4: 生成分析报告

分析完成后，提供结构化报告，包括：

| 报告章节 | 内容 |
|----------|------|
| 执行摘要 | 日志文件名、总条目数、时间范围 |
| 日志级别统计 | 各级别（ERROR、WARN、INFO、DEBUG）数量及占比 |
| 错误分析 | 错误总数、Top 5-10 错误类型，每个关键错误包含：错误消息、出现次数、典型堆栈跟踪、首次/最后出现时间 |
| 关键发现 | 主要问题总结、错误集中时间段、异常模式 |
| 行动建议 | 🔴立即处理（严重问题）、🟡后续优化（建议改进）、🔵需进一步调查（待确认问题） |

### 针对不同日志格式的命令

**JSON 日志**
```bash
cat <log_file> | jq 'select(.level=="ERROR")'
cat <log_file> | jq -r 'select(.level=="ERROR") | .message' | sort | uniq -c | sort -rn
```

**结构化文本日志**
```bash
grep "\[ERROR\]" <log_file> | awk -F' - ' '{print $2}' | sort | uniq -c | sort -rn
```

**Apache/Nginx 访问日志**
```bash
awk '$9 >= 400' <log_file>
awk '$9 >= 400 {print $9}' <log_file> | sort | uniq -c | sort -rn
```

**多文件分析**
```bash
grep -h "ERROR" *.log | sort | uniq -c | sort -rn
```

## 常见错误

| 错误场景 | 原因 | 解决方案 |
|----------|------|----------|
| 无法识别日志格式 | 日志格式非标准或混合格式 | 先使用 `head` 命令预览，手动确定格式 |
| grep 返回空结果 | 日志级别关键词大小写不匹配 | 使用 `-i` 参数忽略大小写 |
| 统计结果不准确 | 日志格式中有特殊字符 | 调整 `awk` 分隔符或使用正则表达式 |
| 大文件处理缓慢 | 日志文件过大 | 使用 `tail -n` 限制分析范围或按时间段过滤 |
| jq 解析失败 | JSON 格式不完整或有语法错误 | 先用 `head` 检查格式，使用 `-R` 选项处理原始输入 |

## 分析技巧

1. **先识别格式**: 了解日志格式才能准确提取信息
2. **关注严重性**: ERROR 和 FATAL 优先于 WARN
3. **寻找模式**: 重复错误往往指向系统性问题
4. **提取上下文**: 错误前后的日志提供重要线索
5. **量化结果**: 用具体数字说明问题严重程度

## 报告示例

以下是一个完整的分析报告示例：

```markdown
## 日志分析报告

**分析文件**: application.log
**日志条目总数**: 150,000
**时间范围**: 2024-01-01 00:00:00 - 2024-01-01 23:59:59

---

### 日志级别统计

| 级别  | 数量    | 占比   |
|-------|---------|--------|
| ERROR | 234     | 0.16%  |
| WARN  | 1,234   | 0.82%  |
| INFO  | 148,532 | 99.02% |

---

### 错误分析

**错误总数**: 234

#### Top 5 错误类型

1. **Database connection timeout** (89 次, 38.0%)
```
   java.sql.SQLTimeoutException: Connection timeout after 30s
   ```
   - 首次出现: 2024-01-01 08:23:15
   - 最后出现: 2024-01-01 18:45:32
   - **建议**: 检查数据库连接池配置,增加连接超时时间或优化慢查询

2. **Invalid user token** (45 次, 19.2%)
   ```
   TokenValidationException: JWT token expired
   ```
   - 首次出现: 2024-01-01 09:12:03
   - 最后出现: 2024-01-01 19:23:11
   - **建议**: 检查 token 过期时间配置,考虑实现自动刷新机制

3. **Service unavailable** (34 次, 14.5%)
   ```
   ServiceException: External API returned 503
   ```
   - 首次出现: 2024-01-01 10:05:22
   - 最后出现: 2024-01-01 16:33:45
   - **建议**: 添加重试机制和熔断器

---

### 关键发现

1. 数据库连接超时错误占所有错误的 38%,是最主要的问题
2. 错误主要集中在 08:00-10:00 和 17:00-19:00 时间段
3. 未发现 FATAL 级别错误

---

### 行动建议

#### 🔴 立即处理
- 优先解决数据库连接超时问题,影响范围最大
- 检查 08:00-10:00 时段的负载情况

#### 🟡 后续优化
- 优化 token 验证逻辑,减少过期错误
- 添加服务健康检查和自动重试机制

#### 🔵 需进一步调查
- 分析 Service unavailable 错误的具体原因
- 检查是否存在外部依赖服务问题
   ```

## 变更记录

| 版本 | 日期 | 变更内容 | 作者 |
|------|------|----------|------|
| 1.0.0 | 2025-12-23 | 初始版本，从 agents/log_analyzer.md 迁移为标准 Skill 格式 | 15085 |
