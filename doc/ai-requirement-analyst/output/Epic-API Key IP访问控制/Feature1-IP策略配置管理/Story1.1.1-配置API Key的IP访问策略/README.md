### 【Story 1.1.1】配置 API Key 的 IP 访问策略

#### 用户故事说明

作为 **平台用户（或管理员）**，我想要通过接口为指定 API Key 设置 IP 访问策略（白名单或黑名单模式 + IP/CIDR 列表），以便在 Key 泄露或遭受异常来源请求时，能精准限制该 Key 的可用 IP 范围。

#### 简化验收条件（A/C）

##### 【正常场景】设置白名单策略

- **Given**：用户持有有效的管理员/用户凭证，且目标 Key ID 存在并属于该用户
- **When**：调用 `PUT /api/keys/{id}/ip_policy`，body 为 `{"mode":"whitelist","ips":["1.2.3.4/32","10.0.0.0/8"]}`
- **Then**：接口返回 200，该 Key 的 IP 策略更新成功；后续使用该 Key 发起的请求将按白名单模式校验

##### 【正常场景】清空 IP 策略（恢复无限制）

- **Given**：目标 Key 已设置过 IP 策略
- **When**：调用 `PUT /api/keys/{id}/ip_policy`，body 为 `{"mode":"","ips":[]}`（或 mode 为空）
- **Then**：接口返回 200，该 Key 的 IP 策略被清除，不再进行 IP 校验

##### 【主要异常场景】CIDR 格式错误

- **Given**：用户提交的 IP 列表中包含非法格式，如 `"999.0.0.0/8"` 或 `"not-an-ip"`
- **When**：调用 `PUT /api/keys/{id}/ip_policy`
- **Then**：接口返回 400，错误信息明确指出哪条 IP/CIDR 格式非法，策略不被保存

##### 【主要异常场景】无权操作他人 Key

- **Given**：请求方为普通用户，目标 Key 属于其他用户
- **When**：调用 `PUT /api/keys/{id}/ip_policy`
- **Then**：接口返回 403，策略不被修改

#### 备注说明

- mode 取值仅限 `whitelist` 或 `blacklist`，其他值应返回 400
- ips 列表为空时，若 mode 不为空，含义待确认（建议：等同于清空策略或拒绝所有，需与产品确认）
- 单个 Key 的 IP 策略条目数量上限待确认（建议默认 100 条）
- 🔴 **重点细化**：接口权限边界——普通用户是否可为自己的 Key 设置黑名单？管理员是否可操作任意用户的 Key？
- 🟡 **待确认**：是否需要支持查询当前 Key 的 IP 策略（GET /api/keys/{id}/ip_policy）？
- 🟢 **已明确**：支持 CIDR 格式（如 `10.0.0.0/8`）；拒绝时错误码为 `IP_NOT_ALLOWED`
