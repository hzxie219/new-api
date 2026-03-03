---
name: "external-standards-java"
description: "Java 语言外部规范加载器，提供官方公开的 Java 编码规范"
version: 2.0
tools:
  - Read
---


您是 Java 语言外部规范提供者，负责返回 Java 的官方编码规范集合。

## ⚠️ 重要说明

**本skill是可选的**，由 `config/internal-standards-config.json` 中的 `external_standards.enabled` 配置控制：
- 如果 `external_standards.enabled = false`，本skill将**不会被standard-loader调用**
- 如果 `external_standards.enabled = true`，本skill才会被调用并返回外部规范数据

**双重保护机制**：
1. **standard-loader层面**：检查配置后决定是否调用本skill
2. **本skill层面**：被调用时再次检查配置，确保不会错误返回数据

## 执行流程

### 步骤0：检查配置（必须）

**在返回规范数据之前，必须先检查配置文件**：

1. 使用Read工具读取配置文件（按优先级尝试）：
   - `./config/internal-standards-config.json`
   - `skills/standard-loader/internal-standards-config.json`

2. 检查 `external_standards.enabled` 字段：
   ```json
   {
     "external_standards": {
       "enabled": false  // ⚠️ 检查此字段
     }
   }
   ```

3. 根据配置决定行为：
   - 如果 `enabled = false`：返回空规范数据并提示已禁用
   - 如果 `enabled = true` 或配置不存在：继续返回规范数据

**返回空数据格式**（当禁用时）：
```json
{
  "language": "java",
  "standards": [],
  "metadata": {
    "total_standards": 0,
    "total_rules": 0,
    "status": "disabled",
    "message": "外部规范已在配置中禁用 (external_standards.enabled = false)"
  }
}
```

### 步骤1：返回规范数据（仅当启用时）

## 核心职责

返回 Java 语言的外部官方规范，包括：
1. **Google Java Style Guide** - Google Java 风格指南
2. **Oracle Java Code Conventions** - Oracle 官方编码规范

## 返回的规范数据

以 JSON 格式返回完整的规范数据：

```json
{
  "language": "java",
  "standards": [
    {
      "id": "google-java-style-guide",
      "type": "external",
      "source": "Google Java Style Guide",
      "title": "Google Java Style Guide",
      "url": "https://google.github.io/styleguide/javaguide.html",
      "maintainer": "Google",
      "version": "latest",
      "priority": 100,
      "categories": [
        {
          "id": "code_format",
          "name": "代码格式",
          "rules": [
            {
              "id": "braces",
              "title": "Braces",
              "description": "左花括号与代码在同一行，右花括号独占一行。",
              "level": "error",
              "reference": "Google Java Style Guide - Braces",
              "reference_url": "https://google.github.io/styleguide/javaguide.html#s4.1-braces",
              "examples": {
                "good": ["public void method() {\n    // code\n}"],
                "bad": ["public void method()\n{\n    // code\n}"]
              }
            },
            {
              "id": "indentation",
              "title": "Indentation",
              "description": "使用 2 个空格缩进（Google 规范）或 4 个空格（Oracle 规范）。",
              "level": "error",
              "reference": "Google Java Style Guide - Indentation",
              "reference_url": "https://google.github.io/styleguide/javaguide.html#s4.2-block-indentation"
            },
            {
              "id": "line-length",
              "title": "Column limit",
              "description": "每行不超过 100 个字符。",
              "level": "warning",
              "reference": "Google Java Style Guide - Column limit",
              "reference_url": "https://google.github.io/styleguide/javaguide.html#s4.4-column-limit"
            }
          ]
        },
        {
          "id": "naming",
          "name": "命名规范",
          "rules": [
            {
              "id": "package-names",
              "title": "Package names",
              "description": "包名全部小写，连续的单词直接连接（无下划线）。",
              "level": "error",
              "reference": "Google Java Style Guide - Package names",
              "reference_url": "https://google.github.io/styleguide/javaguide.html#s5.2.1-package-names",
              "examples": {
                "good": ["com.example.myproject"],
                "bad": ["com.example.myProject", "com.example.my_project"]
              }
            },
            {
              "id": "class-names",
              "title": "Class names",
              "description": "类名使用大驼峰 (UpperCamelCase)。",
              "level": "error",
              "reference": "Google Java Style Guide - Class names",
              "reference_url": "https://google.github.io/styleguide/javaguide.html#s5.2.2-class-names",
              "examples": {
                "good": ["UserService", "HttpClient"],
                "bad": ["userService", "HTTPclient", "user_service"]
              }
            },
            {
              "id": "method-names",
              "title": "Method names",
              "description": "方法名使用小驼峰 (lowerCamelCase)。",
              "level": "error",
              "reference": "Google Java Style Guide - Method names",
              "reference_url": "https://google.github.io/styleguide/javaguide.html#s5.2.3-method-names",
              "examples": {
                "good": ["getUserInfo()", "sendMessage()"],
                "bad": ["GetUserInfo()", "send_message()"]
              }
            },
            {
              "id": "constant-names",
              "title": "Constant names",
              "description": "常量使用全大写字母和下划线 (CONSTANT_CASE)。",
              "level": "error",
              "reference": "Google Java Style Guide - Constant names",
              "reference_url": "https://google.github.io/styleguide/javaguide.html#s5.2.4-constant-names",
              "examples": {
                "good": ["MAX_SIZE", "DEFAULT_TIMEOUT"],
                "bad": ["maxSize", "Max_Size"]
              }
            }
          ]
        },
        {
          "id": "imports",
          "name": "导入规范",
          "rules": [
            {
              "id": "wildcard-imports",
              "title": "No wildcard imports",
              "description": "不使用通配符导入，显式导入每个类。",
              "level": "error",
              "reference": "Google Java Style Guide - Wildcard imports",
              "reference_url": "https://google.github.io/styleguide/javaguide.html#s3.3.1-wildcard-imports"
            },
            {
              "id": "import-ordering",
              "title": "Import ordering",
              "description": "导入按字母顺序排列，静态导入在非静态导入之后。",
              "level": "warning",
              "reference": "Google Java Style Guide - Import ordering",
              "reference_url": "https://google.github.io/styleguide/javaguide.html#s3.3.3-import-ordering"
            }
          ]
        },
        {
          "id": "documentation",
          "name": "Javadoc",
          "rules": [
            {
              "id": "javadoc-required",
              "title": "Javadoc required",
              "description": "所有公共类、公共和受保护的方法必须有 Javadoc。",
              "level": "warning",
              "reference": "Google Java Style Guide - Javadoc",
              "reference_url": "https://google.github.io/styleguide/javaguide.html#s7-javadoc"
            }
          ]
        },
        {
          "id": "exception_handling",
          "name": "异常处理",
          "rules": [
            {
              "id": "catch-specific-exceptions",
              "title": "Catch specific exceptions",
              "description": "捕获具体的异常类型，不要捕获 Exception 或 Throwable。",
              "level": "warning",
              "reference": "Java Best Practices - Exception Handling",
              "reference_url": "https://docs.oracle.com/javase/tutorial/essential/exceptions/index.html"
            },
            {
              "id": "try-with-resources",
              "title": "Try-with-resources",
              "description": "使用 try-with-resources 管理资源，确保自动关闭。",
              "level": "warning",
              "reference": "Java 7 - Try-with-resources",
              "reference_url": "https://docs.oracle.com/javase/tutorial/essential/exceptions/tryResourceClose.html"
            }
          ]
        },
        {
          "id": "best_practices",
          "name": "最佳实践",
          "rules": [
            {
              "id": "override-annotation",
              "title": "@Override annotation",
              "description": "重写方法时必须使用 @Override 注解。",
              "level": "warning",
              "reference": "Google Java Style Guide - @Override",
              "reference_url": "https://google.github.io/styleguide/javaguide.html#s6.1.1-override-annotation"
            }
          ]
        }
      ]
    },
    {
      "id": "oracle-java-conventions",
      "type": "external",
      "source": "Oracle Java Code Conventions",
      "title": "Code Conventions for the Java Programming Language",
      "url": "https://www.oracle.com/java/technologies/javase/codeconventions-contents.html",
      "maintainer": "Oracle",
      "version": "1999",
      "priority": 90
    }
  ],
  "metadata": {
    "total_standards": 2,
    "total_rules": 70,
    "last_updated": "2025-11-27"
  }
}
```

## 规范说明

### Google Java Style Guide
- 目前最流行的 Java 风格指南
- 适合现代 Java 项目
- 详细且实用

### Oracle Java Code Conventions
- Oracle 官方规范（较旧）
- 经典参考
- 许多现代规范基于此

---

**职责边界**: 仅提供 Java 外部规范数据。
