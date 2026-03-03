# 概要设计对外接口生成命令 (Chapter 3)

## 命令描述
生成模块概要设计说明书的第3章"对外接口"，包括API接口定义和消息接口定义。本命令特别强化增量开发场景，优先检索和复用现有代码中的接口实现，确保接口设计的一致性和向后兼容性。

## 命令语法
```bash
/chapter3 $1
```

## 参数说明

- **$1**: 必填参数，需求分析文档路径，指向包含接口需求的需求文档

## 输出内容

**完成设计文档章节3**：

- `3. 对外接口`
  - `3.1 API接口`
  - `3.2 消息接口`

## 输出文件

追加到现有设计文档：`doc/design/tech_design.md`

## ⚠️ 核心强制约束（必须严格遵守）

### 1. Prompt文件步骤严格执行原则

**❌ 禁止的行为**：

- 跳过Prompt文件中定义的任何步骤或子步骤
- 合并多个步骤一起执行
- 只生成内容而不写入文件
- 不生成Prompt要求的检查点文件（`xxx_check.md`）
- 不执行Prompt要求的验证操作

**✅ 必须的行为**：

- **设计模板**: templates/design_template.md
- **API设计规范**: references/specifications/RESTful_API格式规范v3.0.md
- **设计checklist**: references/specifications/设计方案checklist.md

- **逐步执行**：严格按照Prompt文件中定义的步骤1、步骤2、步骤3...顺序执行
- **子步骤不跳过**：每个步骤内的1.1、1.2、1.3等子步骤必须全部执行
- **检查点必生成**：Prompt中要求生成的`xxx_check.md`检查点文件必须生成
- **写入必验证**：每次使用Edit/Write工具写入后，必须使用Read工具验证写入成功
- **失败必重试**：验证失败时必须重新执行写入操作，直到成功

### 增量开发优先原则
```
API接口定义优先在代码路径中检索需求文档$1中列举的接口，处理流程：
1. 在项目根目录代码目录中检索需求文档$1中列举的接口
2. 能检索到的接口，按照需求描述对已存在的接口进行修改，并输出接口设计
3. 检索不到的接口，按照需求描述对接口进行全新生成，接口设计参考规范：references/specifications/RESTful_API格式规范v3.0.md
```

### 章节完整性要求
如果需求文档中未检索到未对外提供API接口或消息接口，需明确写明"不涉及"及原因。

## 执行流程

**⚠️ 核心执行原则（必须严格遵守）**：

1. **强制写入原则**：每个步骤完成后，**必须**使用Edit/Write工具将生成的内容写入目标文件，不能只生成内容而不保存
2. **验证确认原则**：每次写入后，**必须**使用Read工具读取目标文件验证内容已正确写入
3. **失败重试原则**：如果验证失败，**必须**重新执行写入操作，直到验证通过
4. **完整性原则**：**你的上下文窗口会在接近限制时自动被压缩，因此不要因为Token预算问题提前停止任务，即使预算快用完，也要尽可能完整执行任务**
5. **检查点文件原则**：生成的检查点文件`xxx_check.md`保留，**不允许删除**
6. **检查点原则**：每个子步骤都有"完成检查点"，**必须**在该检查点执行强制写入和验证操作

**❌ 错误做法**：

- 生成内容后认为任务完成，没有写入文件
- 调用Agent生成内容后，没有将Agent返回的内容写入文件
- 跳过验证步骤，假设内容已写入

**✅ 正确做法**：
- 步骤3生成3.1 API接口 → 立即使用Edit工具写入 → 使用Read工具验证 → 确认通过
- 步骤6生成3.2 消息接口 → 立即使用Edit工具追加 → 使用Read工具验证 → 确认通过

### 步骤1：API接口需求提取与分类

**执行目标**：从需求文档中提取所有API接口需求，读取接口的代码位置，并分类识别。

**执行内容**：

#### 1. 读取需求文档
```
输入：$1 需求分析文档路径
处理：提取所有API接口相关需求
```

#### 2. 接口需求分类
从需求文档中识别和实现类型为API接口功能的需求：

**API接口需求分类**：

1. 全新接口
2. 增量接口

**API接口代码位置读取**：

1. 记录API代码位置、函数位置、所属业务仓库

#### 3. 完成检查点（强制执行）

**⚠️ 关键操作：步骤1仅做分析，不需要写入目标文件**

**注意事项**：
- 步骤1的输出是内存中的需求分析结果，供后续步骤使用
- 无需写入文件，但需要确保分析结果完整准确
- 后续步骤将基于此分析结果检索和生成接口设计

**验证清单**：
- [ ] 所有API接口需求已提取
- [ ] API接口分类（全新/增量）完成
- [ ] API代码位置信息已记录

------

### 步骤2：现有接口检索与分析

**执行目标**：在`{所属业务仓库}`(代码仓库)和API索引.md中检索所有需求涉及的API接口实现，为增量开发提供基础。

**执行原则**：你的上下文窗口会在接近限制时自动被压缩，因此不要因为Token预算问题提前停止任务，即使预算快用完，也要尽可能完整执行任务。

**查阅项目知识库**：

- `{所属业务仓库}/doc/kb/技术知识库/API索引.md`了解API接口代码位置和基本定义
- 根据`{所属业务仓库}/doc/kb/技术知识库/API索引.md`检索出的代码位置，在当前目录的代码项目中检索

**执行内容**：

#### 1. API代码位置检索
- 根据需求文档中提取的代码位置信息和所属业务仓库，在`{所属业务仓库}/doc/kb/技术知识库/API索引.md`中加速检索
- 结合需求文档和API索引.md中的信息整理全部API接口表格，表格列：API路径，代码位置信息

#### 2. API接口检索内容

​	**检索的API信息如下（请求体和响应体优先级最高）：**

- **HTTP Method**: GET, POST, PUT, DELETE, PATCH, OPTIONS, HEAD
- **路径 (Path)**: 完整的 URL 路径，包括路径参数占位符
- **请求体**：API接口的请求参数结构
- **响应体**：API接口的响应体结构
- **功能摘要 (Summary)**: 1句话描述该接口的功能
- **所属分组 (Tags)**: 功能分类标签
- **代码位置**: 接口定义所在的文件路径和函数名
- **认证要求**: 是否需要认证，需要什么权限

输出API接口check表格，只需要按照表格格式输出，不需要输出其他内容，路径`doc/design/`，文件名`API接口_check.md`

#### 3. 完成检查点（强制执行）

**⚠️ 关键操作：必须将检索结果写入检查点文件**

**强制执行步骤**：

1. **使用Write工具**将步骤2检索到的API接口check表格写入检查点文件`doc/design/API接口_check.md`（不允许删除）
2. **使用Read工具**读取检查点文件，验证内容已正确写入
3. 如果验证失败，重新执行写入操作

**验证清单**：
- [ ] 检查点文件`doc/design/API接口_check.md`已创建
- [ ] 表格包含所有需求中涉及的API接口
- [ ] 每个API包含：API路径、代码位置信息
- [ ] 表格格式正确（Markdown表格格式）

**如果验证失败**：
- 检查目录`doc/design/`是否存在，不存在则创建
- 检查表格格式是否符合Markdown规范
- 重新执行写入操作直到验证通过

------

### 步骤3：生成第三章节对外接口中的3.1. API接口内容

**角色定义：**

读取专家角色定义
- 使用Read工具读取 `references/experts/api/restful-api-generator.md`
- 了解RESTful API生成专家的角色定义和核心职责
应用专家角色并执行任务

你现在是**RESTful API生成专家**，专门负责基于需求文档和API设计规范生成完整的RESTful API接口定义。

**你的核心职责**（来自 restful-api-generator.md）：

1. **API接口代码生成**
   - 根据需求生成符合规范的 RESTful API 接口代码
   - 支持全新接口开发，生成完整的路由、模型、逻辑代码
   - 支持增量开发，在现有接口基础上新增字段或功能
   - 确保生成的代码符合项目规范、可维护、安全且高性能

2. **RESTful API v3.0规范遵循**
   - URI格式标准：`/{module}/{version}/{resources}`
   - HTTP方法语义：GET/POST/PUT/PATCH/DELETE
   - 标准响应格式：code/message/data结构
   - 批量操作支持

3. **FastAPI框架技术栈**
   - 使用FastAPI定义路由和接口
   - 使用Pydantic进行数据验证和序列化
   - 使用async/await实现异步处理
   - 集成SQLAlchemy ORM进行数据库操作

**前置条件**：读取步骤2中生成的`doc/design/API接口_check.md`；如果没有检索到API接口需求，3.1章节输出内容为"不涉及"

**执行目标**：

- 严格按照API规范进行API接口设计：`references/specifications/RESTful_API格式规范v3.0.md`
- **根据API代码位置检索生成的API接口check表格`doc/design/API接口_check.md`，按行的顺序，一个API一个API的在当前路径的项目代码中检索，获取API的基本信息（优先获取请求体和响应体结构），并按照下方生成内容结构逐个生成**
- **生成完整的API接口设计**

**增量接口设计原则**：检索到的增量接口，根据步骤2中检索的请求体和响应体参数，将增量开发中的新增参数加入请求参数或者响应参数里，并且列出已存在的请求和响应体字段

**全新接口设计原则**：严格按照API规范进行API接口设计：`references/specifications/RESTful_API格式规范v3.0.md`

#### 1. 检索代码仓库中的API接口

根据API代码位置检索生成的API接口check表格`doc/design/API接口_check.md`，按行顺序执行，逐个 API的在当前路径的项目代码中检索，获取API的基本信息（优先获取请求体和响应体结构），并按照下方生成内容结构逐个生成

**详细检索步骤**：

##### 1.1 读取API接口check文件

使用Read工具读取`doc/design/API接口_check.md`文件，获取所有API的文件路径和行号信息。

**文件格式示例**：

```markdown
| API名称 | 文件路径 | 行号范围 |
|---------|---------|---------|
| 创建xxx配置 | xxxxServer/routers/xxx_router.py/.java/.php/.go | 45-120 |
| 查询xxx列表 | xxxServer/routers/xxx_router.py/.java/.php/.go | 122-180 |
```

##### 1.2 逐个API精确搜索

针对每个API，按以下顺序使用Grep工具精确搜索关键信息：

###### 1.2.1 搜索路由定义（获取HTTP方法和路径）

根据文件扩展名选择对应的搜索模式：

**Python (FastAPI/Flask/Django)**
```python
# FastAPI路由装饰器
Grep(
  pattern="@router\\.(get|post|put|delete|patch)\\([\"']|@app\\.(get|post|put|delete|patch)\\([\"']",
  path="<API文件路径>",
  output_mode="content",
  -A=5
)

# Django URL配置
Grep(
  pattern="path\\([\"'].*[\"'],.*views\\.",
  path="<API文件路径>",
  output_mode="content",
  -A=3
)
```

**Java (Spring Boot)**
```java
// Spring注解
Grep(
  pattern="@(GetMapping|PostMapping|PutMapping|DeleteMapping|PatchMapping|RequestMapping)\\([\"']",
  path="<API文件路径>",
  output_mode="content",
  -A=5
)
```

**Go (Gin/Echo/Mux)**
```go
// Gin/Echo路由注册
Grep(
  pattern="\\.(GET|POST|PUT|DELETE|PATCH)\\([\"']|router\\.(Handle|HandleFunc)\\(",
  path="<API文件路径>",
  output_mode="content",
  -A=5
)
```

**PHP (Laravel/Symfony)**
```php
// Laravel路由
Grep(
  pattern="Route::(get|post|put|delete|patch)\\([\"']",
  path="<API文件路径>",
  output_mode="content",
  -A=5
)
```

**示例输出**：
```python
# Python FastAPI
@router.post("/aaa/bbb")
async def create_xxx(request: xxxRequest):

# Java Spring
@PostMapping("/aaa/bbb")
public ResponseEntity<xxxResponse> createXXX(@RequestBody xxxRequest request) {

# Go Gin
router.POST("/aaa/bbb", func(c *gin.Context) {

# PHP Laravel
Route::post('/aaa/bbb', [xxxController::class, 'create']);
```

###### 1.2.2 搜索请求体模型（获取请求参数结构）

根据编程语言选择对应的搜索模式：

**Python (Pydantic/Dataclass)**
```python
# Pydantic模型定义
Grep(
  pattern="class.*Request.*\\(BaseModel\\):|class.*Request.*\\(Schema\\):",
  path="<API文件路径或models目录>",
  output_mode="content",
  -A=20
)

# 函数参数类型注解
Grep(
  pattern="(async )?def.*\\(.*request:.*Request|def.*\\(.*\\*\\*kwargs",
  path="<API文件路径>",
  output_mode="content",
  -A=3
)
```

**Java (POJO/DTO)**
```java
// 请求DTO类
Grep(
  pattern="(public |)class.*Request\\s*\\{|@RequestBody.*Request",
  path="<API文件路径或dto/entity目录>",
  output_mode="content",
  -A=30
)

// 方法参数注解
Grep(
  pattern="@RequestBody|@RequestParam|@PathVariable",
  path="<API文件路径>",
  output_mode="content",
  -A=3
)
```

**Go (Struct)**
```go
// 请求结构体
Grep(
  pattern="type.*Request struct \\{",
  path="<API文件路径或models目录>",
  output_mode="content",
  -A=20
)

// JSON绑定
Grep(
  pattern="c\\.ShouldBindJSON\\(|c\\.Bind\\(",
  path="<API文件路径>",
  output_mode="content",
  -B=3,
  -A=5
)
```

**PHP (Array/Object)**
```php
// 请求验证规则
Grep(
  pattern="\\$request->validate\\(|protected \\$rules|public function rules\\(",
  path="<API文件路径>",
  output_mode="content",
  -A=15
)

// 请求类
Grep(
  pattern="class.*Request extends FormRequest",
  path="<API文件路径>",
  output_mode="content",
  -A=20
)
```

###### 1.2.3 搜索响应体模型（获取响应参数结构）

根据编程语言选择对应的搜索模式：

**Python**
```python
# Response模型定义
Grep(
  pattern="class.*Response.*\\(BaseModel\\):|class.*Response.*\\(Schema\\):",
  path="<API文件路径或models目录>",
  output_mode="content",
  -A=15
)

# return语句中的响应结构
Grep(
  pattern="return.*\\{.*[\"']code[\"']|return Response\\(|return JSONResponse\\(",
  path="<API文件路径>",
  output_mode="content",
  -B=2,
  -A=5
)
```

**Java**
```java
// 响应DTO类
Grep(
  pattern="(public |)class.*Response\\s*\\{|ResponseEntity<",
  path="<API文件路径或dto目录>",
  output_mode="content",
  -A=20
)

// 返回语句
Grep(
  pattern="return ResponseEntity\\.|return.*Response\\.|return Result\\.",
  path="<API文件路径>",
  output_mode="content",
  -B=3,
  -A=3
)
```

**Go**
```go
// 响应结构体
Grep(
  pattern="type.*Response struct \\{",
  path="<API文件路径或models目录>",
  output_mode="content",
  -A=15
)

// JSON响应
Grep(
  pattern="c\\.JSON\\(|c\\.IndentedJSON\\(",
  path="<API文件路径>",
  output_mode="content",
  -B=2,
  -A=5
)
```

**PHP**
```php
// 响应资源类
Grep(
  pattern="class.*Resource extends JsonResource|return response\\(\\)->json\\(",
  path="<API文件路径>",
  output_mode="content",
  -A=15
)

// 返回JSON
Grep(
  pattern="return.*->json\\(|return response\\(\\)",
  path="<API文件路径>",
  output_mode="content",
  -B=2,
  -A=5
)
```

###### 1.2.4 搜索接口功能说明（获取文档字符串）

根据编程语言选择对应的搜索模式：

**Python**
```python
# 函数文档字符串
Grep(
  pattern="(async )?def [a-zA-Z_][a-zA-Z0-9_]*\\(",
  path="<API文件路径>",
  output_mode="content",
  -A=10  # 查看函数定义后的docstring
)
```

**Java**
```java
// JavaDoc注释
Grep(
  pattern="/\\*\\*|@Api|@ApiOperation",
  path="<API文件路径>",
  output_mode="content",
  -A=8
)

// 方法签名
Grep(
  pattern="(public|protected|private).*(get|post|put|delete|create|update|query)[A-Z]",
  path="<API文件路径>",
  output_mode="content",
  -B=5,
  -A=3
)
```

**Go**
```go
// 函数注释
Grep(
  pattern="^// [A-Z].*|^/\\*",
  path="<API文件路径>",
  output_mode="content",
  -A=5
)

// 函数签名
Grep(
  pattern="func [a-zA-Z][a-zA-Z0-9]*\\(",
  path="<API文件路径>",
  output_mode="content",
  -B=3,
  -A=3
)
```

**PHP**
```php
# PHPDoc注释
Grep(
  pattern="/\\*\\*|@param|@return",
  path="<API文件路径>",
  output_mode="content",
  -A=8
)

# 方法签名
Grep(
  pattern="(public|protected|private) function [a-zA-Z]",
  path="<API文件路径>",
  output_mode="content",
  -B=5,
  -A=3
)
```

###### 1.2.5 搜索认证装饰器（获取权限要求）

根据编程语言选择对应的搜索模式：

**Python**
```python
# 认证装饰器
Grep(
  pattern="@require_auth|@permission|@login_required|Depends\\(.*auth|Security\\(",
  path="<API文件路径>",
  output_mode="content",
  -B=2,
  -A=2
)
```

**Java**
```java
# Spring Security注解
Grep(
  pattern="@PreAuthorize|@Secured|@RolesAllowed|@PermitAll",
  path="<API文件路径>",
  output_mode="content",
  -B=2,
  -A=2
)
```

**Go**
```go
# 中间件调用
Grep(
  pattern="AuthRequired|RequireAuth|AuthMiddleware|JWTMiddleware",
  path="<API文件路径>",
  output_mode="content",
  -B=2,
  -A=2
)
```

**PHP**
```php
# Laravel中间件
Grep(
  pattern="->middleware\\(|@middleware|protected \\$middleware",
  path="<API文件路径>",
  output_mode="content",
  -B=2,
  -A=2
)
```

##### 1.3 信息汇总

将检索到的信息按照以下结构整理，并追加输出写入`doc/design/API接口_check.md`：

| 信息项 | 检索方法 | 提取内容 |
|-------|---------|---------|
| HTTP方法 | 路由装饰器 | GET/POST/PUT/DELETE/PATCH |
| API路径 | 路由装饰器参数 | /xxx/config |
| 功能说明 | 函数文档字符串 | 创建xxx配置 |
| 请求体 | Request模型或参数注解 | 字段名、类型、说明 |
| 响应体 | Response模型或return语句 | code、message、data结构 |
| 认证要求 | 装饰器或Depends | 需要的权限级别 |
| 代码位置 | API接口_check.md | 文件路径:行号范围 |

**检索失败的降级处理**：

- 如果Grep搜索未找到请求体/响应体模型，使用Read工具读取完整文件内容，手动分析
- 如果模型定义在其他文件，使用Glob搜索models目录：`models/**/*request*.py`
- 记录检索失败的API，在步骤2中根据需求文档手动设计

#### 2. 生成：API接口设计

基于检索后的内容生成API接口设计：

- API接口总览表格
- 每个API接口设计

**生成内容结构**（强制要求每个API接口按照以下格式生成内容，不需要输入其他内容）：

**API接口总览**

| API路径  | 功能类型 | 描述 | 所属业务仓库 | 代码位置       |
| -------- | -------- | ---- | ------------ | -------------- |
| /aaa/bbb | 增量     | xxx  | xxxServer    | xxx.py line:20 |

**接口1: xxx接口(增量or全量)**

**接口信息**:

- **方法**: POST
- **路径**: /aaa/bbb/xxx
- **功能**: xxx功能
- **认证**: xxx鉴权码
- **代码位置**: aaa/ccc/xxx.py/java/go/php:1xx-2xx line

| 请求体 | 参数名 | 类型   | 方向 | 说明 |
| ------ | ------ | ------ | ---- | ---- |
|        | id     | int    | in   | id   |
|        | type   | string | in   | 类型 |
|        | name   | string | in   | 名称 |
| 请求体 |        |        |      |      |
| 返回值 |        |        |      |      |

```
{
    "id": 1,
    "type": "xxx",
    "name": "xxx"
}
```

| 响应体 | 参数名  | 类型   | 方向 | 说明   |
| ------ | ------- | ------ | ---- | ------ |
|        | code    | string | out  | 状态码 |
|        | message | string | out  | 消息   |
|        | data    | string | out  | 数据   |

```
{
    "code": "success",
    "message": "success",
    "data": {
        "name": "xxxx",
        "type": "xxx",
        "desc": ""
    }
}
```

#### 3. 完成检查点（强制执行）

**⚠️ 关键操作：必须将3.1章节内容写入目标文件**

**强制执行步骤**：

1. **使用Edit工具**将步骤3生成的完整3.1 API接口章节（包括所有API接口设计）追加到目标文件`doc/design/tech_design.md`
2. **使用Read工具**读取目标文件，验证3.1章节已正确写入
3. 如果验证失败，重新执行写入操作

**验证清单**：
- [ ] 3.1章节标题"### 3.1. API接口"存在于目标文件
- [ ] 所有需求中的API接口都已生成设计
- [ ] 每个API接口包含：接口信息、请求体表格+JSON示例、响应体表格+JSON示例
- [ ] 增量接口的新增参数已标注清楚
- [ ] 全新接口符合RESTful规范
- [ ] 表格格式正确，JSON格式正确

**如果验证失败**：
- 检查目标文件路径是否正确
- 检查Edit工具的old_string定位是否准确
- 确认所有API接口都已包含在内容中
- 验证表格和JSON格式是否符合Markdown规范
- 重新执行写入操作直到验证通过

------

### 步骤4：API接口完整性检查

#### 1. 检查点（强制执行）

确保以下内容完整：

**API索引.md**：

- [ ] **需求中所有接口都已完成设计**
- [ ] **每个API的请求体和响应体都已描述清楚**
- [ ] **增量接口的新增参数已加入请求参数或者响应参数**
- [ ] **全新接口设计符合references/specifications/RESTful_API格式规范v3.0.md**
- [ ] **每个接口的基本信息**完整（方法、路径、功能说明、认证要求、代码位置、主要参数）

**⚠️ 关键操作：验证3.1章节完整性**

**强制执行步骤**：
1. **使用Read工具**读取目标文件`doc/design/tech_design.md`
2. 逐项检查步骤4中的完整性清单
3. 如果任何检查项未通过，回到步骤3重新生成并写入

**验证清单（全部通过才能继续）**：
- [ ] 需求中所有API接口都已设计
- [ ] 每个API的请求体和响应体都已描述清楚
- [ ] 增量接口的新增参数已标注
- [ ] 全新接口符合RESTful规范
- [ ] 每个接口基本信息完整

**如果验证失败**：
- 定位缺失或不完整的接口
- 回到步骤3补充对应接口的设计
- 使用Edit工具将补充内容追加到目标文件
- 重新执行步骤4验证直到全部通过

---

### 步骤5：需求中消息接口的检索

**执行目标**：从需求文档中提取消息接口需求，并在现有代码和文档中检索消息接口实现。

**执行内容**：

#### 1. 消息接口需求提取

从需求文档$1中识别消息接口需求和消息接口的所属业务仓库：

**消息接口识别关键词**：
- Pulsar、Kafka、RabbitMQ、消息队列、MQ
- Topic、消息、事件、通知
- Producer（生产者）、Consumer（消费者）
- 异步处理、事件驱动

**消息接口需求分类**：

```markdown
1. 【Pulsar消息接口】
   - 标识方式：Pulsar、Topic等关键词
   - 识别要点：
     * Topic名称
     * 消息方向（生产/消费/双向）
     * 消息格式定义
     * 消费者组配置
     * 所属业务仓库

2. 【其他消息系统】
   - Kafka、RabbitMQ等其他消息队列系统
```

**消息接口需求统计**：
```markdown
| 消息类型 | 数量 | 需求编号列表 | 增量/全新 | 所属业务仓库 |
|---------|------|------------|----------|---------|
| Pulsar消息 | XX个 | REQ-MSG-001, REQ-MSG-003 | 增量 | xxxServer |
```

#### 2. 现有消息接口检索

**检索范围**(所属业务仓库源自：1. 消息接口需求提取)：

- `{所属业务仓库}`中的消息定义
- data_channel模块中的消息格式

**检索内容**：

##### 2.1 Pulsar消息接口检索
```
搜索目标：
  - Topic定义（在代码中搜索Topic常量或配置）
  - Producer代码（搜索pulsar.Producer相关代码）
  - Consumer代码（搜索pulsar.Consumer相关代码）
  - 消息格式定义（Pydantic Model或dataclass）

输出信息：
  - Topic名称
  - 消息方向（Producer/Consumer）
  - 消息格式（当前定义）
  - 代码位置（文件路径和行号）
  - 消费者组（如果是Consumer）
```

**检索结果示例**：
```markdown
### 已检索到的Pulsar消息
| Topic名称 | 消息方向 | 消息格式 | 代码位置 | 说明 | 所属业务仓库 |
|----------|---------|---------|---------|------|-------|
| xxx-change | Producer | xxxChangeMessage | data_xxx/xxx/xxx.py:15 | xxx变更通知 | xxxServer |
| audit-log-event | Producer | AuditLogMessage | xxx/xxx/xxx.py:20 | xxx事件 | xxxServer |
| asset-sync-request | Consumer | AssetSyncMessage | xxx/xxx/mq/xxx.py:45 | xxx同步请求 | xxxServer |

### 未检索到的Pulsar消息（需全新设计）
| Topic名称（需求中） | 消息方向 | 需求描述 | 设计类型 | 所属业务仓库 |
|------------------|---------|---------|---------|--------|
| xxxx-event | Producer | xxx事件通知 | 全新设计 | xxxServer |
```

#### 3. 增量/全新消息接口识别

对每个消息接口需求进行增量/全新判断：

```markdown
| 消息名称 | 消息类型 | 增量/全新 | 现有位置（如增量） | 修改类型 | 所属业务仓库 |
|---------|---------|---------|------------------|---------|--------|
| xxx变更通知 | Pulsar | 增量 | xxx/xxx.py:15 | 消息格式扩展 | xxxServer |
| xxx事件 | Pulsar | 全新 | - | - | xxxServer |
```

#### 4. 步骤5完成检查点（强制执行）

**⚠️ 关键操作：步骤5仅做检索，不需要写入目标文件**

**注意事项**：

- 步骤5的输出是内存中的消息接口检索结果，供步骤6使用
- 无需写入文件，但需要确保检索结果完整准确
- 步骤6将基于此检索结果生成消息接口设计

**验证清单**：
- [ ] 所有消息接口需求已提取
- [ ] 现有消息接口已检索完成
- [ ] 增量/全新消息接口已分类
- [ ] 检索结果结构清晰（Topic名称、消息方向、消息格式、代码位置）

---

### 步骤6：生成第三章节对外接口中的3.2 消息接口内容

**前置条件**：步骤5在需求中检索到消息接口需求；如果没有检索到消息接口需求，3.2章节输出内容为"不涉及"

**执行目标**：基于步骤5的检索结果和需求，生成完整的消息接口设计。

对于检索到的现有消息接口，进行增量设计，**增量消息接口设计原则**：
- 基于现有消息格式进行扩展

- 保持向后兼容性，新增字段使用可选字段

- 明确标注新增字段

- 评估对现有消费者的影响

  

对于需求中的新消息接口，进行全新设计，**全新消息接口设计原则**：

- 遵循现有消息命名规范
- 消息格式清晰、完整
- 考虑未来扩展性
- 明确生产者和消费者

#### 1. 生成：消息接口设计

**生成内容结构**：

**消息接口1**

| 名称    | *msg_hdr*                        |       |
| ------- | -------------------------------- | ----- |
| 说明    | *该二进制消息的用途，注意事项等* |       |
| 字段名  | 类型/长度(单位bit)               | 说明  |
| param_1 | 32                               | *x*xx |
| param_2 | 24                               | xxx   |

**完整消息示例**：

```json
{
  "param_1": "xxxxx",
  "param_2": "xxxx"
}
```

#### 2. 完成检查点（强制执行）

**⚠️ 关键操作：必须将3.2章节内容写入目标文件**

**强制执行步骤**：
1. **使用Edit工具**将步骤6生成的完整3.2 消息接口章节（包括所有消息接口设计）追加到目标文件`doc/design/tech_design.md`
2. **使用Read工具**读取目标文件，验证3.2章节已正确写入
3. 如果验证失败，重新执行写入操作

**验证清单**：
- [ ] 3.2章节标题"### 3.2. 消息接口"存在于目标文件（或"不涉及"说明）
- [ ] 所有需求中的消息接口都已生成设计
- [ ] 每个消息接口包含：名称、说明、字段表格、完整JSON示例
- [ ] 增量消息接口的新增字段已标注清楚
- [ ] 全新消息接口设计完整
- [ ] 表格格式正确，JSON格式正确

**如果验证失败**：
- 检查目标文件路径是否正确
- 检查Edit工具的old_string定位是否准确
- 确认所有消息接口都已包含在内容中
- 验证表格和JSON格式是否符合Markdown规范
- 重新执行写入操作直到验证通过

------

### 步骤7：章节完整性检查与文件验证

**⚠️ 关键操作：验证第3章所有章节已写入目标文件**

#### 1. 内容完整性检查
检查章节是否都存在内容


#### 2. API接口章节检查

- [ ] 每个API接口的详细设计存在
- [ ] 接口设计包含：接口路径、HTTP方法、功能说明、请求参数、响应格式、业务规则
- [ ] 表格格式符合Markdown规范
- [ ] JSON示例格式正确

#### 3. 消息接口章节检查：

- [ ] 每个消息接口的详细设计存在
- [ ] 消息设计包含：Topic名称、消息方向、消息格式、字段说明、消息示例
- [ ] 表格格式符合Markdown规范

#### 4. 验证失败处理

**如果章节缺失**：
1. 确认是步骤3（API接口）还是步骤6（消息接口）的内容缺失
2. 返回对应步骤重新执行写入操作
3. 再次执行本验证步骤

---

## 注意事项

### 增量开发要点
- **优先检索**：在生成新接口前，必须先在代码仓库中检索现有实现
- **谨慎修改**：修改现有接口时，必须评估向后兼容性
- **明确标注**：清晰标注接口是"新增"、"修改"还是"复用"
- **兼容性优先**：修改接口时优先考虑向后兼容，避免破坏性变更

### RESTful规范遵循：references/specifications/RESTful_API格式规范v3.0.md
- **URI资源化**：URI使用名词而非动词，如 `/regions` 而非 `/getRegions`
- **HTTP方法语义**：正确使用GET/POST/PUT/PATCH/DELETE
- **统一响应格式**：所有接口使用统一的响应格式
- **错误码规范**：使用标准错误码，避免自定义错误码

### 接口设计原则
- **单一职责**：每个接口只负责一个明确的功能
- **幂等性**：GET/PUT/DELETE操作必须是幂等的
- **安全性**：所有接口需要认证和授权
- **性能要求**：明确每个接口的性能目标

### 文档完整性
- **参数完整**：每个参数必须有类型、必填性、说明、示例、约束
- **响应完整**：包含成功和各种失败场景的响应示例
- **业务规则**：明确接口的业务规则和约束
- **测试用例**：提供基本的测试用例

