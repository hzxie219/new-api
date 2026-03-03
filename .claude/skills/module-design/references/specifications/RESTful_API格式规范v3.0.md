# RESTful API格式规范v3.0

## **1 RESTful风格**

RESTful是一种API的设计风格，用于定义在CS、BS架构下暴露服务端接口。

### **1.1 核心特征**

1）**URI资源化**：URI代表资源，不包含动作。例如：/class/students

2）**动作由HTTP方法决定**：

   - GET：获取资源
   - POST：创建资源
   - PUT：整体替换资源
   - PATCH：部分更新资源
   - DELETE：删除资源

3）**无状态**：客户端无状态，服务端可维护必要状态（如登录状态）

4）**数据安全**：使用HTTPS协议加密数据

### **1.2 HTTP方法说明**

| **方法** | **描述** | **幂等** | **示例** |
| :------- | :------- | :------- | :------- |
| GET | 获取资源 | 是 | GET /class/students |
| POST | 创建资源 | 否 | POST /class/students {"name": "Jake", "age": 18} |
| PUT | 整体替换 | 是 | PUT /class/students/2 {"name": "Jim", "age": 19} |
| PATCH | 部分更新 | 否 | PATCH /class/students/2 {"age": 20} |
| DELETE | 删除资源 | 是 | DELETE /class/students/2 |

**注意**：支持_method参数覆盖HTTP方法，用于特殊场景。

## **2 命名规范**

### **2.1 URL命名**
- **规则**：全部使用小写字母，单词间无分隔符
- **正则**：`/^[0-9a-z]+$/`
- **示例**：`/mailsetting/testemail/v1`

### **2.2 变量命名**

变量是指URL查询参数（Query String）和请求体中json字段的名字。

- **原则**：尽可能使用单一名称涵盖要表达的意义，尽可能不要用多个词连接
- **规则**：驼峰法命名，大小写敏感
- **正则**：`/^_?[a-z][0-9A-Za-z]*$/`
- **特殊**：允许前缀下划线防止关键字冲突
- **示例**：`{“userName”: "zhangsan", "_enableAuth": 1, "_cache": false}`

## **3 URI格式**

为防止多产品之间接口冲突，需要改进URL加上业务前缀：`/{module}/{version}/{resources}`

**参数说明**：

- **module**：模块名
- **version**：api版本
- **resources**：资源

## **4 GET请求处理**

### **4.1 GET OVER POST**
当参数过长或复杂时，使用POST + _method=GET

### **4.2 GET语义转变**
将GET请求转换为POST创建任务的方式

### **4.3 传统GET**
- **长度限制**：建议2000字节以内
- **参数风格**：
  - 多参数：`/cars/?color=blue&type=sedan`
  - 数组：`/appointments?users=[id1,id2]`
  - 对象：`/appointments?params={users:[id1,id2], age:18}`

## **5 批量处理接口**

### **5.1 相同URI批量操作**
```json
POST /class/students
{
  "data": [
    {"name": "Jake", "age": 18},
    {"name": "Jakson", "age": 19}
  ]
}
```

### **5.2 不同URI批量查询（只读）**
```json
POST /status/batch?_method=GET
{
  "items": [
    {"path": "status/cpu", "param": {"max": 90}},
    {"path": "status/memory"},
    {"path": "status/disk"}
  ]
}
```

### **5.3 返回结果格式**
推荐使用整体错误+详细错误格式：
```json
{
  "code": "Success",
  "message": "Partial operation completed",
  "data": [
    {"code": "Success", "message": "success", "data": {"id": 1}},
    {"code": "InvalidParameter", "message": "invalid param"}
  ]
}
```

## **6 应用层错误码**

（接口设计中找不到对应文件可忽略）

遵守规范：@公共错误码定义v1.0.md

## **7 回复消息格式**

```json
{
    "code":string,              //必选，使用公共错误码定义中的错误码字符串
    "message":string,           //在code不为空时必选，code为空时可选。
    "data": object or array     //可选
}
```

**错误码示例**：
- 成功：`"code": "Success"`
- 参数错误：`"code": "InvalidParameter"`
- 权限不足：`"code": "PermissionDenied"`
- 内部错误：`"code": "InternalError"`

## **8 缓存设计**

- **禁用缓存**：`_cache=0`
- **启用缓存**：`_cache=1`
- 其他取值保留，产品线不得自行定义

## **9 异步处理**

1. 创建异步任务，返回任务ID
2. 客户端轮询任务ID获取进度
3. 返回完成度百分比（0-100）

## **10 国际化**

### **10.1 编码**
- 所有字符串采用UTF-8编码
- Content-Type：`application/json;charset=UTF-8`

### **10.2 语言**
- HTTP头部：`Accept-Language:zh-CN`
- 参数覆盖：`lang=en`

### **10.3 时间格式**
- 优先：UTC秒数
- 备选：ISO 8601格式 `"yyyy-MM-dd'T'HH:mm:ss.SSS'Z'"`

## **11 ACL权限定义**

**权限简写**：N(POST)、D(DELETE)、P(PUT)、G(GET)、H(PATCH)、*(所有权限)

**格式**：`[!]权限简写 资源路径`

**示例**：
- Admin：`* /users * /resource/partA`
- Operator：`NPG /users GH /resource/partA`
- Guest：`!NDP /users GH /resource/partA`