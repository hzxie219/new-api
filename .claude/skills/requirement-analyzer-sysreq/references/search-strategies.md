# 代码搜索策略和索引建立

本文档提供详细的真实代码搜索策略和代码索引建立方法。

## 了解项目架构（必须优先执行）

**目的**：快速了解项目整体架构和功能域。

扫描项目根目录下各代码仓库（排除 `.claude`、`README.md`、`doc`、`agent-rules`）的 `doc/kb/` 目录下的知识库文件（如果存在）：
- 仓库概览.md：了解项目用途和核心功能域
- 仓库架构.md：了解技术栈和系统架构
- 仓库依赖.md：了解项目的依赖关系
- 业务知识库：了解业务功能和业务流程
- 技术知识库：了解有哪些API接口、数据模型等技术实现
- 测试知识库：了解测试相关信息

**注意**：这些文档仅用于快速理解架构，**不能作为最终引用！**

## 在实际代码仓库中搜索真实代码（必须执行）

**确定代码仓库位置**：
- 使用项目根目录作为代码仓库根目录
- 识别项目根目录下所有代码仓库文件夹（排除 `.claude`、`README.md`、`doc`、`agent-rules`）
- 在每个代码仓库文件夹中搜索代码文件

### 1. 搜索API路由文件

**搜索策略**：
- 使用Glob：`**/routes/**/*.py` 或 `**/routers/**/*.py` 或 `**/api/**/*.py`
- 使用Grep搜索接口装饰器：`@app.route` 或 `@router.post` 或 `@router.get` 或 `@router.put` 或 `@router.delete` 或 `@api_route`
- 使用Read工具读取代码文件，确认接口路径和行号
- 记录每个接口的真实代码位置

**🚨 接口路径验证强约束**：
- ✅ **必须使用Read工具读取代码文件，直接查看接口定义**
- ✅ **必须从代码装饰器中提取真实的接口路径**（如：`@router.post("/dashboard/app/import/app_list")`）
- ✅ **必须记录准确的行号**（从Read工具结果中获取）
- ❌ **禁止凭记忆或推测编造接口路径**
- ❌ **禁止从文档中复制接口路径**（文档可能过时）
- ❌ **禁止假设接口路径存在**（必须在代码中找到）
- 🔍 **验证方法**：使用Grep工具搜索接口路径字符串，确认其存在于代码文件中

**示例**：
```python
# 在 project/api/routes/app.py:156 中找到
@router.post("/api/app/whitelist")
async def add_whitelist(request: WhitelistRequest):
    ...
```

记录为：
- POST /api/app/whitelist: 添加白名单 [代码位置: project/api/routes/app.py:156]

### 2. 搜索数据表定义

**搜索策略**：
- 使用Glob：`**/models/**/*.py` 或 `**/model/**/*.py` 或 `**/entity/**/*.py`
- 使用Grep搜索表定义关键词：`__tablename__` 或 `Table(` 或 `CREATE TABLE` 或 `create_table`
- 记录每个表的真实代码位置、数据库类型（根据项目实际使用的数据库）、Schema
- **注意**：只记录表名、数据库类型和Schema，不记录字段信息

**示例**：
```python
# 在 project/models/app.py:23 中找到
class AppWhitelist(Base):
    __tablename__ = 'app_whitelist'
    __table_args__ = {'schema': 'public'}
```

记录为：
- app_whitelist: 白名单表 [数据库: PostgreSQL, Schema: public, 代码位置: project/models/app.py:23]

### 3. 搜索配置文件

**搜索策略**：
- 使用Glob：`**/config/**/*.yaml` 或 `**/config/**/*.yml` 或 `**/config/**/*.json` 或 `**/conf/**/*.yaml`
- 记录配置文件路径和简要说明配置文件的作用
- **注意**：只列出需求相关的配置文件，如果没有则标注"不涉及"

**示例**：
- project/config/app_config.yaml: 应用配置，用于超时和重试设置

### 4. 搜索Pulsar消息格式

**搜索策略**：
- 使用Glob：`**/message/**/*.py` 或 `**/event/**/*.py` 或 `**/producer/**/*.py` 或 `**/consumer/**/*.py`
- 使用Grep搜索Pulsar相关关键词：`pulsar` 或 `Producer` 或 `Consumer` 或 `topic` 或 `send_message`
- 记录消息Topic、消息用途（传递什么业务数据）、代码位置
- **注意**：只列出需求相关的消息格式，如果没有则标注"不涉及"

**示例**：
- app-event-topic: 传递应用事件数据 [代码位置: project/message/app_events.py:45]

### 5. 搜索业务逻辑文件

**搜索策略**：
- 使用Glob：`**/services/**/*.py` 或 `**/logic/**/*.py`
- 使用Grep搜索关键业务类或函数
- 记录关键业务组件的真实代码位置

**示例**：
- WhitelistService: 白名单管理服务 [代码位置: project/services/whitelist_service.py:12]

## 建立代码索引

搜索完成后，将发现的内容组织成结构化索引：

```
代码索引：

接口索引（按功能分类）：
【功能域A相关接口】
- HTTP方法 /api/path/to/endpoint1: 接口功能说明 [代码位置: 项目名/路径/文件.py:行号]
- HTTP方法 /api/path/to/endpoint2: 接口功能说明 [代码位置: 项目名/路径/文件.py:行号]

【功能域B相关接口】
- HTTP方法 /api/path/to/endpoint3: 接口功能说明 [代码位置: 项目名/路径/文件.py:行号]

数据表索引：
{数据库类型A}表：
- 表名: 功能说明 [数据库: {数据库类型}, Schema: xxx, 代码位置: 项目名/路径/文件.py:行号]
{数据库类型B}表：
- 表名: 功能说明 [数据库: {数据库类型}, Schema: xxx, 代码位置: 项目名/路径/文件.py:行号]

配置文件索引：
- 配置文件路径: 配置作用说明
- 配置文件路径: 配置作用说明

消息格式索引：
- 消息Topic: 消息用途说明 [代码位置: 项目名/路径/文件.py:行号]
- 消息Topic: 消息用途说明 [代码位置: 项目名/路径/文件.py:行号]

业务组件索引：
- 组件名称: 组件功能说明 [代码位置: 项目名/路径/文件.py:行号]
```

## 接口分类原则

- 根据业务功能域对接口进行分类（如：数据导入、列表查询、批量操作、权限管理等）
- 同一分类下的接口按照业务流程顺序排列（如：下载模板 → 数据校验 → 执行导入）
- 分类名称应清晰明确，使用【】标记，反映业务功能而非技术实现
- 每个分类下至少包含1个接口，最多不超过10个接口
- 如果某个功能域接口过多，可以进一步细分子分类

## 重要提示

- 使用 `Glob` 工具查找代码文件
- 使用 `Grep` 工具搜索关键代码
- 使用 `Read` 工具读取具体代码，记录准确的行号
- 所有引用必须是真实代码文件，不能是docs/文档
