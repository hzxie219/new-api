# 知识库配置agent-rules指南

生成将知识库配置到 agent-rules 系统的配置示例。

**参数解析**

解析命令参数$ARGMENTS:
- `--kb-path` {知识库路径} 指定知识库所在目录。默认：doc/kb
- `--product` {产品名称} 指定产品类型（scc/scp/skyops/dmp等）。可选，如不指定则提供通用配置模板

---

## 执行流程

### 步骤 1：收集项目信息

从当前项目收集以下信息：

1. **项目目录名称**：使用 git 提取仓库名称，如果 git 不可用则使用当前目录名

2. **Git 仓库信息**：
   - 仓库地址（git remote -v，提取 origin 的 fetch URL）
   - 当前分支名（git branch --show-current）

3. **知识库路径**：{知识库路径参数}，默认 doc/kb

4. **仓库能力描述**：
   - 读取 `{知识库路径}/仓库概览.md 和 外部接入指南.md` 文件
   - 将该内容转化为的能力描述（说明这个仓库能做什么, 提供什么具体能力，用于 skills 描述）
   - 如果无法读取，使用占位符 "[请根据仓库功能填写能力描述]"

### 步骤 2：生成配置示例

直接输出以下内容（使用收集的真实信息替换占位符）：

```
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
📋 知识库配置示例
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

以下配置可直接复制到对应产品的配置文件中：

配置文件位置：
• SCC 产品：/home/agent-rules/configs/knowledge/scc.yaml
• SCP 产品：/home/agent-rules/configs/knowledge/scp.yaml
• SkyOps 产品：/home/agent-rules/configs/knowledge/skyops.yaml
• DMP 产品：/home/agent-rules/configs/knowledge/dmp.yaml
• 新产品：/home/agent-rules/configs/knowledge/{新产品名}.yaml

━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
```

然后输出 YAML 配置块：

```yaml
profiles:
  {项目目录名称}:
    name: "{项目目录名称}"
    description: "{仓库能力描述 - 从仓库概览和外部接入指南提取的能力说明}"
    tags:
      - {产品标签}  # 修改为：SCC/SCP/SkyOps/DMP 等
    knowledge:
      - name: "{项目目录名称}"
        repo: "{Git仓库地址}"
        branch: "{当前分支}"
        remote_path: "{知识库路径}"
        # 注意：普通知识库不需要 is_product_kb 字段（默认为 false）
        # 如需配置为产品级知识库，添加：is_product_kb: true
```

**重要说明**：
- 如果用户指定了 `--product` 参数，则在 tags 中自动填入对应的产品标签，并在上方提示具体的配置文件路径
- `description` 字段会被用作 skills 的描述，请确保准确反映仓库能力
- **普通知识库不需要 is_product_kb 字段**（默认值为 false，无需显式指定）
- 只有产品级知识库才需要添加 `is_product_kb: true`

---

## 输出要求

- 直接输出配置示例，不要冗长的说明
- 使用分隔线和图标使内容清晰
- 配置示例使用代码块格式，方便复制
- 所有真实信息（仓库地址、分支、能力描述等）必须从实际项目中提取

