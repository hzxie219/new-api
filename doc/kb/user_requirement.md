API Key 的 IP 黑名单控制
背景描述
事故/盗用场景下，只有 RPM 无法阻止“来源异常”的请求；需要按来源 IP 精准限制 Key 的可用范围。
需求
●提供接口 PUT /api/keys/{id}/ip_policy 配置 IP 策略：mode=whitelist|blacklist + ips（支持 CIDR）。
●业务接口鉴权时获取客户端 IP：仅在请求来自“可信代理”时才使用 X-Forwarded-For，否则使用直连远端 IP。
●校验规则：
○whitelist：客户端 IP 不命中即拒绝
○blacklist：客户端 IP 命中即拒绝
●拒绝访问：HTTP 403，固定错误码 IP_NOT_ALLOWED。
验收标准
●配置 whitelist 仅允许 1.2.3.4/32：
○客户端 IP 为 1.2.3.4 调用业务接口成功
○客户端 IP 为 8.8.8.8 调用业务接口返回 403 且错误码正确
●配置 blacklist 禁止某网段后，命中该网段的请求必定返回 403 且错误码正确。