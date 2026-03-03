# API 索引

## 概览

- **接口总数**：约 180 个
- **认证方式**：Session Cookie（用户登录态）、API Token（Bearer Token / x-api-key）、管理员权限（AdminAuth）、Root权限（RootAuth）、Token ReadOnly（只读令牌）
- **基础路径**：`/api`（管理接口）、`/v1`（AI中继接口）、`/mj`（Midjourney）、`/suno`（Suno）、`/kling`（Kling视频）等

---

## 接口列表

### 系统与状态

| 方法 | 路径 | 功能说明 | 认证 | 代码位置 | 主要参数 |
|------|------|----------|------|----------|----------|
| GET | /api/setup | 获取初始化安装状态 | 否 | `new-api/controller/setup.go:GetSetup` | 无 |
| POST | /api/setup | 提交初始化安装配置 | 否 | `new-api/controller/setup.go:PostSetup` | username, password, confirmPassword, SelfUseModeEnabled, DemoSiteEnabled |
| GET | /api/status | 获取系统状态及配置 | 否 | `new-api/controller/misc.go:GetStatus` | 无 |
| GET | /api/status/test | 测试服务器与数据库连通性 | AdminAuth | `new-api/controller/misc.go:TestStatus` | 无 |
| GET | /api/uptime/status | 获取 Uptime Kuma 监控状态 | 否 | `new-api/controller/uptime_kuma.go:GetUptimeKumaStatus` | 无 |
| GET | /api/notice | 获取系统公告 | 否 | `new-api/controller/misc.go:GetNotice` | 无 |
| GET | /api/user-agreement | 获取用户协议 | 否 | `new-api/controller/misc.go:GetUserAgreement` | 无 |
| GET | /api/privacy-policy | 获取隐私政策 | 否 | `new-api/controller/misc.go:GetPrivacyPolicy` | 无 |
| GET | /api/about | 获取关于页面内容 | 否 | `new-api/controller/misc.go:GetAbout` | 无 |
| GET | /api/home_page_content | 获取首页内容 | 否 | `new-api/controller/misc.go:GetHomePageContent` | 无 |
| GET | /api/ratio_config | 获取公开倍率配置 | 否（需启用） | `new-api/controller/ratio_config.go:GetRatioConfig` | 无 |
| GET | /api/verification | 发送邮箱验证码 | 否 | `new-api/controller/misc.go:SendEmailVerification` | email（query） |
| GET | /api/reset_password | 发送密码重置邮件 | 否 | `new-api/controller/misc.go:SendPasswordResetEmail` | email（query） |
| POST | /api/user/reset | 重置密码 | 否 | `new-api/controller/misc.go:ResetPassword` | token, password |

---

### 用户认证与账户

| 方法 | 路径 | 功能说明 | 认证 | 代码位置 | 主要参数 |
|------|------|----------|------|----------|----------|
| POST | /api/user/register | 用户注册 | 否 | `new-api/controller/user.go:Register` | username, password, email, aff |
| POST | /api/user/login | 用户名密码登录 | 否 | `new-api/controller/user.go:Login` | username, password |
| POST | /api/user/login/2fa | 二次验证完成登录（2FA） | 否（需pending session） | `new-api/controller/twofa.go:Verify2FALogin` | code |
| GET | /api/user/logout | 用户退出登录 | 否 | `new-api/controller/user.go:Logout` | 无 |
| GET | /api/user/self | 获取当前用户信息 | UserAuth | `new-api/controller/user.go:GetSelf` | 无 |
| PUT | /api/user/self | 更新当前用户信息 | UserAuth | `new-api/controller/user.go:UpdateSelf` | display_name, password, email 等 |
| DELETE | /api/user/self | 注销当前用户账号 | UserAuth | `new-api/controller/user.go:DeleteSelf` | 无 |
| GET | /api/user/token | 获取用户 Access Token | UserAuth | `new-api/controller/user.go:GenerateAccessToken` | 无 |
| GET | /api/user/models | 获取当前用户可用模型 | UserAuth | `new-api/controller/user.go:GetUserModels` | 无 |
| PUT | /api/user/setting | 更新用户设置 | UserAuth | `new-api/controller/user.go:UpdateUserSetting` | billing_preference 等 |
| GET | /api/user/groups | 获取用户分组（公开） | 否 | `new-api/controller/group.go:GetUserGroups` | 无 |
| GET | /api/user/self/groups | 获取当前用户分组信息 | UserAuth | `new-api/controller/group.go:GetUserGroups` | 无 |
| GET | /api/user/aff | 获取邀请码 | UserAuth | `new-api/controller/user.go:GetAffCode` | 无 |
| POST | /api/user/aff_transfer | 邀请奖励额度转入余额 | UserAuth | `new-api/controller/user.go:TransferAffQuota` | 无 |
| GET | /api/user/oauth/bindings | 获取当前用户 OAuth 绑定列表 | UserAuth | `new-api/controller/custom_oauth.go:GetUserOAuthBindings` | 无 |
| DELETE | /api/user/oauth/bindings/:provider_id | 解绑当前用户 OAuth | UserAuth | `new-api/controller/custom_oauth.go:UnbindCustomOAuth` | provider_id（路径） |

---

### OAuth 认证

| 方法 | 路径 | 功能说明 | 认证 | 代码位置 | 主要参数 |
|------|------|----------|------|----------|----------|
| GET | /api/oauth/state | 生成 OAuth CSRF State | 否 | `new-api/controller/oauth.go:GenerateOAuthCode` | aff（query） |
| GET | /api/oauth/:provider | 标准 OAuth 回调（GitHub/Discord/OIDC/LinuxDO） | 否 | `new-api/controller/oauth.go:HandleOAuth` | code, state（query） |
| GET | /api/oauth/email/bind | 邮箱绑定回调 | 否 | `new-api/controller/oauth.go:EmailBind` | code（query） |
| GET | /api/oauth/wechat | 微信 OAuth 登录 | 否 | `new-api/controller/wechat.go:WeChatAuth` | code（query） |
| GET | /api/oauth/wechat/bind | 微信 OAuth 绑定 | 否 | `new-api/controller/wechat.go:WeChatBind` | code（query） |
| GET | /api/oauth/telegram/login | Telegram 登录回调 | 否 | `new-api/controller/telegram.go:TelegramLogin` | query params |
| GET | /api/oauth/telegram/bind | Telegram 账号绑定 | 否 | `new-api/controller/telegram.go:TelegramBind` | query params |

---

### Passkey（无密码认证）

| 方法 | 路径 | 功能说明 | 认证 | 代码位置 | 主要参数 |
|------|------|----------|------|----------|----------|
| POST | /api/user/passkey/login/begin | 开始 Passkey 登录流程 | 否 | `new-api/controller/passkey.go:PasskeyLoginBegin` | 无 |
| POST | /api/user/passkey/login/finish | 完成 Passkey 登录流程 | 否 | `new-api/controller/passkey.go:PasskeyLoginFinish` | WebAuthn assertion body |
| GET | /api/user/passkey | 获取 Passkey 状态 | UserAuth | `new-api/controller/passkey.go:PasskeyStatus` | 无 |
| POST | /api/user/passkey/register/begin | 开始 Passkey 注册流程 | UserAuth | `new-api/controller/passkey.go:PasskeyRegisterBegin` | 无 |
| POST | /api/user/passkey/register/finish | 完成 Passkey 注册 | UserAuth | `new-api/controller/passkey.go:PasskeyRegisterFinish` | WebAuthn attestation body |
| POST | /api/user/passkey/verify/begin | 开始 Passkey 验证 | UserAuth | `new-api/controller/passkey.go:PasskeyVerifyBegin` | 无 |
| POST | /api/user/passkey/verify/finish | 完成 Passkey 验证 | UserAuth | `new-api/controller/passkey.go:PasskeyVerifyFinish` | WebAuthn assertion body |
| DELETE | /api/user/passkey | 删除用户 Passkey | UserAuth | `new-api/controller/passkey.go:PasskeyDelete` | 无 |

---

### 两步验证（2FA）

| 方法 | 路径 | 功能说明 | 认证 | 代码位置 | 主要参数 |
|------|------|----------|------|----------|----------|
| GET | /api/user/2fa/status | 获取当前用户 2FA 状态 | UserAuth | `new-api/controller/twofa.go:Get2FAStatus` | 无 |
| POST | /api/user/2fa/setup | 初始化 2FA 设置（生成 secret/QR） | UserAuth | `new-api/controller/twofa.go:Setup2FA` | 无 |
| POST | /api/user/2fa/enable | 启用 2FA | UserAuth | `new-api/controller/twofa.go:Enable2FA` | code |
| POST | /api/user/2fa/disable | 禁用 2FA | UserAuth | `new-api/controller/twofa.go:Disable2FA` | code |
| POST | /api/user/2fa/backup_codes | 重新生成备用码 | UserAuth | `new-api/controller/twofa.go:RegenerateBackupCodes` | code |
| POST | /api/verify | 通用安全验证（2FA 或 Passkey） | UserAuth | `new-api/controller/secure_verification.go:UniversalVerify` | method, code |

---

### 签到

| 方法 | 路径 | 功能说明 | 认证 | 代码位置 | 主要参数 |
|------|------|----------|------|----------|----------|
| GET | /api/user/checkin | 获取用户签到状态和历史 | UserAuth | `new-api/controller/checkin.go:GetCheckinStatus` | month（query） |
| POST | /api/user/checkin | 执行每日签到 | UserAuth | `new-api/controller/checkin.go:DoCheckin` | 无 |

---

### 充值与支付

| 方法 | 路径 | 功能说明 | 认证 | 代码位置 | 主要参数 |
|------|------|----------|------|----------|----------|
| GET | /api/user/topup/info | 获取充值配置信息（支付方式、最低金额等） | UserAuth | `new-api/controller/topup.go:GetTopUpInfo` | 无 |
| GET | /api/user/topup/self | 获取当前用户充值记录 | UserAuth | `new-api/controller/topup.go:GetUserTopUps` | 无 |
| POST | /api/user/topup | 使用兑换码充值 | UserAuth | `new-api/controller/topup.go:TopUp` | key |
| POST | /api/user/pay | 发起 Epay 支付 | UserAuth | `new-api/controller/topup.go:RequestEpay` | amount, pay_type |
| POST | /api/user/amount | 计算 Epay 支付金额 | UserAuth | `new-api/controller/topup.go:RequestAmount` | amount |
| POST | /api/user/stripe/pay | 发起 Stripe 支付 | UserAuth | `new-api/controller/topup_stripe.go:RequestStripePay` | amount |
| POST | /api/user/stripe/amount | 计算 Stripe 支付金额 | UserAuth | `new-api/controller/topup_stripe.go:RequestStripeAmount` | amount |
| POST | /api/user/creem/pay | 发起 Creem 支付 | UserAuth | `new-api/controller/topup_creem.go:RequestCreemPay` | product_id |
| GET | /api/user/epay/notify | Epay 支付回调通知（GET） | 否 | `new-api/controller/topup.go:EpayNotify` | query params |
| POST | /api/user/epay/notify | Epay 支付回调通知（POST） | 否 | `new-api/controller/topup.go:EpayNotify` | body params |
| POST | /api/stripe/webhook | Stripe Webhook 回调 | 否 | `new-api/controller/topup_stripe.go:StripeWebhook` | Stripe-Signature header |
| POST | /api/creem/webhook | Creem Webhook 回调 | 否 | `new-api/controller/topup_creem.go:CreemWebhook` | body |

---

### 订阅管理（用户端）

| 方法 | 路径 | 功能说明 | 认证 | 代码位置 | 主要参数 |
|------|------|----------|------|----------|----------|
| GET | /api/subscription/plans | 获取可用订阅计划列表 | UserAuth | `new-api/controller/subscription.go:GetSubscriptionPlans` | 无 |
| GET | /api/subscription/self | 获取当前用户订阅信息 | UserAuth | `new-api/controller/subscription.go:GetSubscriptionSelf` | 无 |
| PUT | /api/subscription/self/preference | 更新订阅计费偏好 | UserAuth | `new-api/controller/subscription.go:UpdateSubscriptionPreference` | billing_preference |
| POST | /api/subscription/epay/pay | 发起 Epay 订阅支付 | UserAuth | `new-api/controller/subscription_payment_epay.go:SubscriptionRequestEpay` | plan_id, pay_type |
| POST | /api/subscription/stripe/pay | 发起 Stripe 订阅支付 | UserAuth | `new-api/controller/subscription_payment_stripe.go:SubscriptionRequestStripePay` | plan_id |
| POST | /api/subscription/creem/pay | 发起 Creem 订阅支付 | UserAuth | `new-api/controller/subscription_payment_creem.go:SubscriptionRequestCreemPay` | plan_id |
| GET | /api/subscription/epay/notify | Epay 订阅回调通知（GET） | 否 | `new-api/controller/subscription_payment_epay.go:SubscriptionEpayNotify` | query params |
| POST | /api/subscription/epay/notify | Epay 订阅回调通知（POST） | 否 | `new-api/controller/subscription_payment_epay.go:SubscriptionEpayNotify` | body params |
| GET | /api/subscription/epay/return | Epay 订阅同步回跳（GET） | 否 | `new-api/controller/subscription_payment_epay.go:SubscriptionEpayReturn` | query params |
| POST | /api/subscription/epay/return | Epay 订阅同步回跳（POST） | 否 | `new-api/controller/subscription_payment_epay.go:SubscriptionEpayReturn` | body params |

---

### 管理员 - 用户管理

| 方法 | 路径 | 功能说明 | 认证 | 代码位置 | 主要参数 |
|------|------|----------|------|----------|----------|
| GET | /api/user/ | 获取所有用户列表（分页） | AdminAuth | `new-api/controller/user.go:GetAllUsers` | page, page_size |
| GET | /api/user/search | 搜索用户 | AdminAuth | `new-api/controller/user.go:SearchUsers` | keyword |
| GET | /api/user/:id | 获取指定用户详情 | AdminAuth | `new-api/controller/user.go:GetUser` | id（路径） |
| POST | /api/user/ | 创建用户 | AdminAuth | `new-api/controller/user.go:CreateUser` | username, password, role, group 等 |
| PUT | /api/user/ | 更新用户信息 | AdminAuth | `new-api/controller/user.go:UpdateUser` | id, username, quota 等 |
| DELETE | /api/user/:id | 删除用户 | AdminAuth | `new-api/controller/user.go:DeleteUser` | id（路径） |
| POST | /api/user/manage | 用户管理操作（启用/禁用/封禁等） | AdminAuth | `new-api/controller/user.go:ManageUser` | id, action |
| GET | /api/user/topup | 获取所有充值记录 | AdminAuth | `new-api/controller/topup.go:GetAllTopUps` | page, page_size |
| POST | /api/user/topup/complete | 管理员手动完成充值订单 | AdminAuth | `new-api/controller/topup.go:AdminCompleteTopUp` | trade_no |
| GET | /api/user/:id/oauth/bindings | 获取指定用户 OAuth 绑定列表 | AdminAuth | `new-api/controller/custom_oauth.go:GetUserOAuthBindingsByAdmin` | id（路径） |
| DELETE | /api/user/:id/oauth/bindings/:provider_id | 管理员解绑用户 OAuth | AdminAuth | `new-api/controller/custom_oauth.go:UnbindCustomOAuthByAdmin` | id, provider_id（路径） |
| DELETE | /api/user/:id/bindings/:binding_type | 管理员清除用户特定类型绑定 | AdminAuth | `new-api/controller/user.go:AdminClearUserBinding` | id, binding_type（路径） |
| DELETE | /api/user/:id/reset_passkey | 管理员重置用户 Passkey | AdminAuth | `new-api/controller/passkey.go:AdminResetPasskey` | id（路径） |
| GET | /api/user/2fa/stats | 获取全站 2FA 使用统计 | AdminAuth | `new-api/controller/twofa.go:Admin2FAStats` | 无 |
| DELETE | /api/user/:id/2fa | 管理员禁用指定用户 2FA | AdminAuth | `new-api/controller/twofa.go:AdminDisable2FA` | id（路径） |

---

### 管理员 - 订阅管理

| 方法 | 路径 | 功能说明 | 认证 | 代码位置 | 主要参数 |
|------|------|----------|------|----------|----------|
| GET | /api/subscription/admin/plans | 获取所有订阅计划（管理员） | AdminAuth | `new-api/controller/subscription.go:AdminListSubscriptionPlans` | 无 |
| POST | /api/subscription/admin/plans | 创建订阅计划 | AdminAuth | `new-api/controller/subscription.go:AdminCreateSubscriptionPlan` | name, price, quota, duration 等 |
| PUT | /api/subscription/admin/plans/:id | 更新订阅计划 | AdminAuth | `new-api/controller/subscription.go:AdminUpdateSubscriptionPlan` | id（路径），plan fields |
| PATCH | /api/subscription/admin/plans/:id | 更新订阅计划状态（启用/禁用） | AdminAuth | `new-api/controller/subscription.go:AdminUpdateSubscriptionPlanStatus` | id（路径），enabled |
| POST | /api/subscription/admin/bind | 绑定订阅（管理员手动） | AdminAuth | `new-api/controller/subscription.go:AdminBindSubscription` | user_id, plan_id |
| GET | /api/subscription/admin/users/:id/subscriptions | 获取用户订阅列表 | AdminAuth | `new-api/controller/subscription.go:AdminListUserSubscriptions` | id（路径） |
| POST | /api/subscription/admin/users/:id/subscriptions | 为用户创建订阅 | AdminAuth | `new-api/controller/subscription.go:AdminCreateUserSubscription` | id（路径），plan_id |
| POST | /api/subscription/admin/user_subscriptions/:id/invalidate | 使用户订阅失效 | AdminAuth | `new-api/controller/subscription.go:AdminInvalidateUserSubscription` | id（路径） |
| DELETE | /api/subscription/admin/user_subscriptions/:id | 删除用户订阅 | AdminAuth | `new-api/controller/subscription.go:AdminDeleteUserSubscription` | id（路径） |

---

### 渠道管理（管理员）

| 方法 | 路径 | 功能说明 | 认证 | 代码位置 | 主要参数 |
|------|------|----------|------|----------|----------|
| GET | /api/channel/ | 获取所有渠道列表（分页） | AdminAuth | `new-api/controller/channel.go:GetAllChannels` | page, page_size, status |
| GET | /api/channel/search | 搜索渠道 | AdminAuth | `new-api/controller/channel.go:SearchChannels` | keyword |
| GET | /api/channel/models | 渠道可用模型列表 | AdminAuth | `new-api/controller/channel.go:ChannelListModels` | 无 |
| GET | /api/channel/models_enabled | 已启用渠道的模型列表 | AdminAuth | `new-api/controller/channel.go:EnabledListModels` | 无 |
| GET | /api/channel/:id | 获取指定渠道详情 | AdminAuth | `new-api/controller/channel.go:GetChannel` | id（路径） |
| POST | /api/channel/:id/key | 获取渠道密钥（需安全验证） | RootAuth+SecureVerification | `new-api/controller/channel.go:GetChannelKey` | id（路径） |
| GET | /api/channel/test | 测试所有渠道连通性 | AdminAuth | `new-api/controller/channel-test.go:TestAllChannels` | 无 |
| GET | /api/channel/test/:id | 测试指定渠道连通性 | AdminAuth | `new-api/controller/channel-test.go:TestChannel` | id（路径） |
| GET | /api/channel/update_balance | 更新所有渠道余额 | AdminAuth | `new-api/controller/channel.go:UpdateAllChannelsBalance` | 无 |
| GET | /api/channel/update_balance/:id | 更新指定渠道余额 | AdminAuth | `new-api/controller/channel.go:UpdateChannelBalance` | id（路径） |
| POST | /api/channel/ | 创建渠道 | AdminAuth | `new-api/controller/channel.go:AddChannel` | name, type, key, models 等 |
| PUT | /api/channel/ | 更新渠道信息 | AdminAuth | `new-api/controller/channel.go:UpdateChannel` | id, name, key 等 |
| DELETE | /api/channel/disabled | 删除所有禁用渠道 | AdminAuth | `new-api/controller/channel.go:DeleteDisabledChannel` | 无 |
| DELETE | /api/channel/:id | 删除指定渠道 | AdminAuth | `new-api/controller/channel.go:DeleteChannel` | id（路径） |
| POST | /api/channel/batch | 批量删除渠道 | AdminAuth | `new-api/controller/channel.go:DeleteChannelBatch` | ids |
| POST | /api/channel/fix | 修复渠道能力表 | AdminAuth | `new-api/controller/channel.go:FixChannelsAbilities` | 无 |
| POST | /api/channel/tag/disabled | 按标签禁用渠道 | AdminAuth | `new-api/controller/channel.go:DisableTagChannels` | tag |
| POST | /api/channel/tag/enabled | 按标签启用渠道 | AdminAuth | `new-api/controller/channel.go:EnableTagChannels` | tag |
| PUT | /api/channel/tag | 编辑渠道标签 | AdminAuth | `new-api/controller/channel.go:EditTagChannels` | tag, new_tag |
| GET | /api/channel/fetch_models/:id | 从上游获取指定渠道模型列表 | AdminAuth | `new-api/controller/channel.go:FetchUpstreamModels` | id（路径） |
| POST | /api/channel/fetch_models | 批量从上游获取模型 | AdminAuth | `new-api/controller/channel.go:FetchModels` | channel_ids |
| POST | /api/channel/batch/tag | 批量设置渠道标签 | AdminAuth | `new-api/controller/channel.go:BatchSetChannelTag` | ids, tag |
| GET | /api/channel/tag/models | 获取标签对应模型列表 | AdminAuth | `new-api/controller/channel.go:GetTagModels` | tag |
| POST | /api/channel/copy/:id | 复制渠道 | AdminAuth | `new-api/controller/channel.go:CopyChannel` | id（路径） |
| POST | /api/channel/multi_key/manage | 管理多密钥渠道 | AdminAuth | `new-api/controller/channel.go:ManageMultiKeys` | channel_id, keys |
| POST | /api/channel/codex/oauth/start | 发起 Codex OAuth 授权 | AdminAuth | `new-api/controller/codex_oauth.go:StartCodexOAuth` | base_url |
| POST | /api/channel/codex/oauth/complete | 完成 Codex OAuth 授权 | AdminAuth | `new-api/controller/codex_oauth.go:CompleteCodexOAuth` | input |
| POST | /api/channel/:id/codex/oauth/start | 为指定渠道发起 Codex OAuth | AdminAuth | `new-api/controller/codex_oauth.go:StartCodexOAuthForChannel` | id（路径） |
| POST | /api/channel/:id/codex/oauth/complete | 完成指定渠道 Codex OAuth | AdminAuth | `new-api/controller/codex_oauth.go:CompleteCodexOAuthForChannel` | id（路径），input |
| POST | /api/channel/:id/codex/refresh | 刷新 Codex 渠道凭证 | AdminAuth | `new-api/controller/codex_oauth.go:RefreshCodexChannelCredential` | id（路径） |
| GET | /api/channel/:id/codex/usage | 获取 Codex 渠道用量 | AdminAuth | `new-api/controller/codex_usage.go:GetCodexChannelUsage` | id（路径） |
| POST | /api/channel/ollama/pull | 触发 Ollama 拉取模型 | AdminAuth | `new-api/controller/channel.go:OllamaPullModel` | channel_id, model |
| POST | /api/channel/ollama/pull/stream | 流式触发 Ollama 拉取模型 | AdminAuth | `new-api/controller/channel.go:OllamaPullModelStream` | channel_id, model |
| DELETE | /api/channel/ollama/delete | 删除 Ollama 模型 | AdminAuth | `new-api/controller/channel.go:OllamaDeleteModel` | channel_id, model |
| GET | /api/channel/ollama/version/:id | 获取 Ollama 版本信息 | AdminAuth | `new-api/controller/channel.go:OllamaVersion` | id（路径） |

---

### 令牌管理（API Token）

| 方法 | 路径 | 功能说明 | 认证 | 代码位置 | 主要参数 |
|------|------|----------|------|----------|----------|
| GET | /api/token/ | 获取当前用户令牌列表（分页） | UserAuth | `new-api/controller/token.go:GetAllTokens` | page, page_size |
| GET | /api/token/search | 搜索令牌 | UserAuth | `new-api/controller/token.go:SearchTokens` | keyword, token |
| GET | /api/token/:id | 获取指定令牌详情 | UserAuth | `new-api/controller/token.go:GetToken` | id（路径） |
| POST | /api/token/ | 创建令牌 | UserAuth | `new-api/controller/token.go:AddToken` | name, quota, expired_time, model_limits 等 |
| PUT | /api/token/ | 更新令牌信息 | UserAuth | `new-api/controller/token.go:UpdateToken` | id, name, quota 等 |
| DELETE | /api/token/:id | 删除令牌 | UserAuth | `new-api/controller/token.go:DeleteToken` | id（路径） |
| POST | /api/token/batch | 批量删除令牌 | UserAuth | `new-api/controller/token.go:DeleteTokenBatch` | ids |

---

### 用量查询

| 方法 | 路径 | 功能说明 | 认证 | 代码位置 | 主要参数 |
|------|------|----------|------|----------|----------|
| GET | /api/usage/token/ | 按令牌查询用量信息 | TokenAuthReadOnly | `new-api/controller/user.go:GetTokenUsage` | 无 |

---

### 兑换码管理（管理员）

| 方法 | 路径 | 功能说明 | 认证 | 代码位置 | 主要参数 |
|------|------|----------|------|----------|----------|
| GET | /api/redemption/ | 获取所有兑换码（分页） | AdminAuth | `new-api/controller/redemption.go:GetAllRedemptions` | page, page_size |
| GET | /api/redemption/search | 搜索兑换码 | AdminAuth | `new-api/controller/redemption.go:SearchRedemptions` | keyword |
| GET | /api/redemption/:id | 获取指定兑换码详情 | AdminAuth | `new-api/controller/redemption.go:GetRedemption` | id（路径） |
| POST | /api/redemption/ | 创建兑换码 | AdminAuth | `new-api/controller/redemption.go:AddRedemption` | name, quota, count |
| PUT | /api/redemption/ | 更新兑换码 | AdminAuth | `new-api/controller/redemption.go:UpdateRedemption` | id, name, quota 等 |
| DELETE | /api/redemption/invalid | 删除所有已使用兑换码 | AdminAuth | `new-api/controller/redemption.go:DeleteInvalidRedemption` | 无 |
| DELETE | /api/redemption/:id | 删除指定兑换码 | AdminAuth | `new-api/controller/redemption.go:DeleteRedemption` | id（路径） |

---

### 日志与使用数据

| 方法 | 路径 | 功能说明 | 认证 | 代码位置 | 主要参数 |
|------|------|----------|------|----------|----------|
| GET | /api/log/ | 获取所有请求日志（分页） | AdminAuth | `new-api/controller/log.go:GetAllLogs` | type, start_timestamp, end_timestamp, username, token_name, model_name, channel, group |
| DELETE | /api/log/ | 删除历史日志 | AdminAuth | `new-api/controller/log.go:DeleteHistoryLogs` | target_timestamp |
| GET | /api/log/stat | 获取日志统计数据 | AdminAuth | `new-api/controller/log.go:GetLogsStat` | start_timestamp, end_timestamp |
| GET | /api/log/self/stat | 获取当前用户日志统计 | UserAuth | `new-api/controller/log.go:GetLogsSelfStat` | start_timestamp, end_timestamp |
| GET | /api/log/search | 搜索所有日志 | AdminAuth | `new-api/controller/log.go:SearchAllLogs` | keyword |
| GET | /api/log/self | 获取当前用户请求日志 | UserAuth | `new-api/controller/log.go:GetUserLogs` | type, start_timestamp, end_timestamp, token_name, model_name |
| GET | /api/log/self/search | 当前用户日志搜索 | UserAuth | `new-api/controller/log.go:SearchUserLogs` | keyword |
| GET | /api/log/token | 按令牌查询日志 | TokenAuthReadOnly | `new-api/controller/log.go:GetLogByKey` | 无 |
| GET | /api/log/channel_affinity_usage_cache | 获取渠道亲和性用量缓存统计 | AdminAuth | `new-api/controller/log.go:GetChannelAffinityUsageCacheStats` | 无 |
| GET | /api/data/ | 获取所有用户配额日期数据 | AdminAuth | `new-api/controller/usedata.go:GetAllQuotaDates` | start_timestamp, end_timestamp, username |
| GET | /api/data/self | 获取当前用户配额日期数据 | UserAuth | `new-api/controller/usedata.go:GetUserQuotaDates` | start_timestamp, end_timestamp |

---

### 分组管理（管理员）

| 方法 | 路径 | 功能说明 | 认证 | 代码位置 | 主要参数 |
|------|------|----------|------|----------|----------|
| GET | /api/group/ | 获取所有可用分组 | AdminAuth | `new-api/controller/group.go:GetGroups` | 无 |

---

### 预填组管理（管理员）

| 方法 | 路径 | 功能说明 | 认证 | 代码位置 | 主要参数 |
|------|------|----------|------|----------|----------|
| GET | /api/prefill_group/ | 获取预填组列表 | AdminAuth | `new-api/controller/prefill_group.go:GetPrefillGroups` | type（query） |
| POST | /api/prefill_group/ | 创建预填组 | AdminAuth | `new-api/controller/prefill_group.go:CreatePrefillGroup` | name, type, content |
| PUT | /api/prefill_group/ | 更新预填组 | AdminAuth | `new-api/controller/prefill_group.go:UpdatePrefillGroup` | id, name, content |
| DELETE | /api/prefill_group/:id | 删除预填组 | AdminAuth | `new-api/controller/prefill_group.go:DeletePrefillGroup` | id（路径） |

---

### Midjourney 任务

| 方法 | 路径 | 功能说明 | 认证 | 代码位置 | 主要参数 |
|------|------|----------|------|----------|----------|
| GET | /api/mj/self | 获取当前用户 MJ 任务列表 | UserAuth | `new-api/controller/midjourney.go:GetUserMidjourney` | page, page_size |
| GET | /api/mj/ | 获取所有 MJ 任务列表（管理员） | AdminAuth | `new-api/controller/midjourney.go:GetAllMidjourney` | page, page_size |

---

### 异步任务（视频/Suno等）

| 方法 | 路径 | 功能说明 | 认证 | 代码位置 | 主要参数 |
|------|------|----------|------|----------|----------|
| GET | /api/task/self | 获取当前用户任务列表 | UserAuth | `new-api/controller/task.go:GetUserTask` | platform, task_id, status, start_timestamp, end_timestamp |
| GET | /api/task/ | 获取所有任务列表（管理员） | AdminAuth | `new-api/controller/task.go:GetAllTask` | platform, task_id, status, channel_id, start_timestamp, end_timestamp |

---

### 供应商元数据管理（管理员）

| 方法 | 路径 | 功能说明 | 认证 | 代码位置 | 主要参数 |
|------|------|----------|------|----------|----------|
| GET | /api/vendors/ | 获取所有供应商列表（分页） | AdminAuth | `new-api/controller/vendor_meta.go:GetAllVendors` | page, page_size |
| GET | /api/vendors/search | 搜索供应商 | AdminAuth | `new-api/controller/vendor_meta.go:SearchVendors` | keyword |
| GET | /api/vendors/:id | 获取指定供应商详情 | AdminAuth | `new-api/controller/vendor_meta.go:GetVendorMeta` | id（路径） |
| POST | /api/vendors/ | 创建供应商 | AdminAuth | `new-api/controller/vendor_meta.go:CreateVendorMeta` | name, logo 等 |
| PUT | /api/vendors/ | 更新供应商信息 | AdminAuth | `new-api/controller/vendor_meta.go:UpdateVendorMeta` | id, name 等 |
| DELETE | /api/vendors/:id | 删除供应商 | AdminAuth | `new-api/controller/vendor_meta.go:DeleteVendorMeta` | id（路径） |

---

### 模型元数据管理（管理员）

| 方法 | 路径 | 功能说明 | 认证 | 代码位置 | 主要参数 |
|------|------|----------|------|----------|----------|
| GET | /api/models/ | 获取所有模型元数据（分页） | AdminAuth | `new-api/controller/model_meta.go:GetAllModelsMeta` | page, page_size |
| GET | /api/models/search | 搜索模型元数据 | AdminAuth | `new-api/controller/model_meta.go:SearchModelsMeta` | keyword, vendor |
| GET | /api/models/:id | 获取指定模型元数据 | AdminAuth | `new-api/controller/model_meta.go:GetModelMeta` | id（路径） |
| POST | /api/models/ | 创建模型元数据 | AdminAuth | `new-api/controller/model_meta.go:CreateModelMeta` | model_name, vendor, description 等 |
| PUT | /api/models/ | 更新模型元数据 | AdminAuth | `new-api/controller/model_meta.go:UpdateModelMeta` | id, model_name 等 |
| DELETE | /api/models/:id | 删除模型元数据 | AdminAuth | `new-api/controller/model_meta.go:DeleteModelMeta` | id（路径） |
| GET | /api/models/missing | 获取配置缺失的模型列表 | AdminAuth | `new-api/controller/missing_models.go:GetMissingModels` | 无 |
| GET | /api/models/sync_upstream/preview | 预览上游模型同步结果 | AdminAuth | `new-api/controller/model_sync.go:SyncUpstreamPreview` | 无 |
| POST | /api/models/sync_upstream | 同步上游模型数据 | AdminAuth | `new-api/controller/model_sync.go:SyncUpstreamModels` | locale |

---

### 模型部署管理（管理员，io.net）

| 方法 | 路径 | 功能说明 | 认证 | 代码位置 | 主要参数 |
|------|------|----------|------|----------|----------|
| GET | /api/deployments/settings | 获取模型部署配置 | AdminAuth | `new-api/controller/deployment.go:GetModelDeploymentSettings` | 无 |
| POST | /api/deployments/settings/test-connection | 测试 io.net 连接 | AdminAuth | `new-api/controller/deployment.go:TestIoNetConnection` | api_key |
| GET | /api/deployments/ | 获取所有部署列表 | AdminAuth | `new-api/controller/deployment.go:GetAllDeployments` | page, page_size |
| GET | /api/deployments/search | 搜索部署 | AdminAuth | `new-api/controller/deployment.go:SearchDeployments` | keyword |
| POST | /api/deployments/test-connection | 测试部署连接 | AdminAuth | `new-api/controller/deployment.go:TestIoNetConnection` | api_key |
| GET | /api/deployments/hardware-types | 获取硬件类型列表 | AdminAuth | `new-api/controller/deployment.go:GetHardwareTypes` | 无 |
| GET | /api/deployments/locations | 获取可用部署地区 | AdminAuth | `new-api/controller/deployment.go:GetLocations` | 无 |
| GET | /api/deployments/available-replicas | 获取可用副本数 | AdminAuth | `new-api/controller/deployment.go:GetAvailableReplicas` | 无 |
| POST | /api/deployments/price-estimation | 估算部署费用 | AdminAuth | `new-api/controller/deployment.go:GetPriceEstimation` | hardware_type, replicas, duration 等 |
| GET | /api/deployments/check-name | 检查集群名称可用性 | AdminAuth | `new-api/controller/deployment.go:CheckClusterNameAvailability` | name（query） |
| POST | /api/deployments/ | 创建部署 | AdminAuth | `new-api/controller/deployment.go:CreateDeployment` | name, model, hardware_type 等 |
| GET | /api/deployments/:id | 获取部署详情 | AdminAuth | `new-api/controller/deployment.go:GetDeployment` | id（路径） |
| GET | /api/deployments/:id/logs | 获取部署日志 | AdminAuth | `new-api/controller/deployment.go:GetDeploymentLogs` | id（路径） |
| GET | /api/deployments/:id/containers | 列出部署容器 | AdminAuth | `new-api/controller/deployment.go:ListDeploymentContainers` | id（路径） |
| GET | /api/deployments/:id/containers/:container_id | 获取容器详情 | AdminAuth | `new-api/controller/deployment.go:GetContainerDetails` | id, container_id（路径） |
| PUT | /api/deployments/:id | 更新部署配置 | AdminAuth | `new-api/controller/deployment.go:UpdateDeployment` | id（路径），config fields |
| PUT | /api/deployments/:id/name | 更新部署名称 | AdminAuth | `new-api/controller/deployment.go:UpdateDeploymentName` | id（路径），name |
| POST | /api/deployments/:id/extend | 延长部署时间 | AdminAuth | `new-api/controller/deployment.go:ExtendDeployment` | id（路径），duration |
| DELETE | /api/deployments/:id | 删除部署 | AdminAuth | `new-api/controller/deployment.go:DeleteDeployment` | id（路径） |

---

### 系统配置（Root 管理员）

| 方法 | 路径 | 功能说明 | 认证 | 代码位置 | 主要参数 |
|------|------|----------|------|----------|----------|
| GET | /api/option/ | 获取系统配置选项列表 | RootAuth | `new-api/controller/option.go:GetOptions` | 无 |
| PUT | /api/option/ | 更新系统配置选项 | RootAuth | `new-api/controller/option.go:UpdateOption` | key, value |
| GET | /api/option/channel_affinity_cache | 获取渠道亲和性缓存统计 | RootAuth | `new-api/controller/channel_affinity_cache.go:GetChannelAffinityCacheStats` | 无 |
| DELETE | /api/option/channel_affinity_cache | 清除渠道亲和性缓存 | RootAuth | `new-api/controller/channel_affinity_cache.go:ClearChannelAffinityCache` | all, rule_name（query） |
| POST | /api/option/rest_model_ratio | 重置模型倍率 | RootAuth | `new-api/controller/option.go:ResetModelRatio` | 无 |
| POST | /api/option/migrate_console_setting | 迁移控制台配置 | RootAuth | `new-api/controller/console_migrate.go:MigrateConsoleSetting` | 无 |

---

### 自定义 OAuth 提供商管理（Root 管理员）

| 方法 | 路径 | 功能说明 | 认证 | 代码位置 | 主要参数 |
|------|------|----------|------|----------|----------|
| POST | /api/custom-oauth-provider/discovery | 通过 Well-Known 发现 OAuth 端点 | RootAuth | `new-api/controller/custom_oauth.go:FetchCustomOAuthDiscovery` | well_known |
| GET | /api/custom-oauth-provider/ | 获取所有自定义 OAuth 提供商 | RootAuth | `new-api/controller/custom_oauth.go:GetCustomOAuthProviders` | 无 |
| GET | /api/custom-oauth-provider/:id | 获取指定 OAuth 提供商详情 | RootAuth | `new-api/controller/custom_oauth.go:GetCustomOAuthProvider` | id（路径） |
| POST | /api/custom-oauth-provider/ | 创建自定义 OAuth 提供商 | RootAuth | `new-api/controller/custom_oauth.go:CreateCustomOAuthProvider` | name, slug, client_id, client_secret 等 |
| PUT | /api/custom-oauth-provider/:id | 更新 OAuth 提供商配置 | RootAuth | `new-api/controller/custom_oauth.go:UpdateCustomOAuthProvider` | id（路径），provider fields |
| DELETE | /api/custom-oauth-provider/:id | 删除 OAuth 提供商 | RootAuth | `new-api/controller/custom_oauth.go:DeleteCustomOAuthProvider` | id（路径） |

---

### 性能管理（Root 管理员）

| 方法 | 路径 | 功能说明 | 认证 | 代码位置 | 主要参数 |
|------|------|----------|------|----------|----------|
| GET | /api/performance/stats | 获取服务性能统计（内存/缓存/磁盘） | RootAuth | `new-api/controller/performance.go:GetPerformanceStats` | 无 |
| DELETE | /api/performance/disk_cache | 清除磁盘缓存 | RootAuth | `new-api/controller/performance.go:ClearDiskCache` | 无 |
| POST | /api/performance/reset_stats | 重置性能统计数据 | RootAuth | `new-api/controller/performance.go:ResetPerformanceStats` | 无 |
| POST | /api/performance/gc | 手动触发 GC | RootAuth | `new-api/controller/performance.go:ForceGC` | 无 |

---

### 倍率同步（Root 管理员）

| 方法 | 路径 | 功能说明 | 认证 | 代码位置 | 主要参数 |
|------|------|----------|------|----------|----------|
| GET | /api/ratio_sync/channels | 获取可同步倍率的渠道列表 | RootAuth | `new-api/controller/ratio_sync.go:GetSyncableChannels` | 无 |
| POST | /api/ratio_sync/fetch | 从上游拉取并同步倍率 | RootAuth | `new-api/controller/ratio_sync.go:FetchUpstreamRatios` | channel_ids, source |

---

### 定价与模型信息（公开/用户）

| 方法 | 路径 | 功能说明 | 认证 | 代码位置 | 主要参数 |
|------|------|----------|------|----------|----------|
| GET | /api/pricing | 获取模型定价及分组信息 | 可选UserAuth | `new-api/controller/pricing.go:GetPricing` | 无 |
| GET | /api/models | 获取当前用户可用模型列表 | UserAuth | `new-api/controller/model.go:DashboardListModels` | 无 |

---

### Dashboard 兼容接口（OpenAI Billing）

| 方法 | 路径 | 功能说明 | 认证 | 代码位置 | 主要参数 |
|------|------|----------|------|----------|----------|
| GET | /dashboard/billing/subscription | 获取订阅/余额信息（OpenAI 兼容） | TokenAuth | `new-api/controller/billing.go:GetSubscription` | 无 |
| GET | /v1/dashboard/billing/subscription | 获取订阅/余额信息（OpenAI v1 兼容） | TokenAuth | `new-api/controller/billing.go:GetSubscription` | 无 |
| GET | /dashboard/billing/usage | 获取用量信息（OpenAI 兼容） | TokenAuth | `new-api/controller/billing.go:GetUsage` | start_date, end_date |
| GET | /v1/dashboard/billing/usage | 获取用量信息（OpenAI v1 兼容） | TokenAuth | `new-api/controller/billing.go:GetUsage` | start_date, end_date |

---

### AI 中继接口（/v1，OpenAI 兼容）

| 方法 | 路径 | 功能说明 | 认证 | 代码位置 | 主要参数 |
|------|------|----------|------|----------|----------|
| GET | /v1/models | 列出可用模型 | TokenAuth | `new-api/controller/model.go:ListModels` | 无 |
| GET | /v1/models/:model | 获取指定模型信息 | TokenAuth | `new-api/controller/model.go:RetrieveModel` | model（路径） |
| POST | /v1/chat/completions | 聊天补全（主要接口） | TokenAuth | `new-api/controller/relay.go:Relay` | model, messages, stream 等 |
| POST | /v1/completions | 文本补全 | TokenAuth | `new-api/controller/relay.go:Relay` | model, prompt, max_tokens 等 |
| POST | /v1/messages | Claude 格式消息接口 | TokenAuth | `new-api/controller/relay.go:Relay` | model, messages, max_tokens 等 |
| POST | /v1/responses | OpenAI Responses API | TokenAuth | `new-api/controller/relay.go:Relay` | model, input 等 |
| POST | /v1/responses/compact | OpenAI Responses API（紧凑模式） | TokenAuth | `new-api/controller/relay.go:Relay` | model, input 等 |
| POST | /v1/images/generations | 图像生成 | TokenAuth | `new-api/controller/relay.go:Relay` | model, prompt, n, size 等 |
| POST | /v1/images/edits | 图像编辑 | TokenAuth | `new-api/controller/relay.go:Relay` | model, image, prompt 等 |
| POST | /v1/edits | 文本编辑（旧接口） | TokenAuth | `new-api/controller/relay.go:Relay` | model, input, instruction |
| POST | /v1/embeddings | 文本向量嵌入 | TokenAuth | `new-api/controller/relay.go:Relay` | model, input |
| POST | /v1/audio/transcriptions | 音频转写 | TokenAuth | `new-api/controller/relay.go:Relay` | model, file, language 等 |
| POST | /v1/audio/translations | 音频翻译 | TokenAuth | `new-api/controller/relay.go:Relay` | model, file |
| POST | /v1/audio/speech | 文字转语音 | TokenAuth | `new-api/controller/relay.go:Relay` | model, input, voice 等 |
| POST | /v1/rerank | 文档重排序 | TokenAuth | `new-api/controller/relay.go:Relay` | model, query, documents 等 |
| POST | /v1/moderations | 内容审核 | TokenAuth | `new-api/controller/relay.go:Relay` | model, input |
| POST | /v1/engines/:model/embeddings | Gemini 引擎嵌入 | TokenAuth | `new-api/controller/relay.go:Relay` | model（路径），input |
| POST | /v1/models/*path | Gemini API 通用路径 | TokenAuth | `new-api/controller/relay.go:Relay` | model（路径），body |
| GET | /v1/realtime | OpenAI Realtime API（WebSocket） | TokenAuth | `new-api/controller/relay.go:Relay` | model（query） |
| POST | /pg/chat/completions | Playground 聊天补全 | UserAuth | `new-api/controller/playground.go:Playground` | model, messages 等 |

---

### Gemini 原生接口（/v1beta）

| 方法 | 路径 | 功能说明 | 认证 | 代码位置 | 主要参数 |
|------|------|----------|------|----------|----------|
| GET | /v1beta/models | 列出 Gemini 模型 | TokenAuth | `new-api/controller/model.go:ListModels` | 无 |
| GET | /v1beta/openai/models | 列出 OpenAI 兼容模型（Gemini兼容路由） | TokenAuth | `new-api/controller/model.go:ListModels` | 无 |
| POST | /v1beta/models/*path | Gemini 原生 API（generateContent/streamGenerateContent等） | TokenAuth | `new-api/controller/relay.go:Relay` | model（路径），contents 等 |

---

### Midjourney 中继接口

| 方法 | 路径 | 功能说明 | 认证 | 代码位置 | 主要参数 |
|------|------|----------|------|----------|----------|
| GET | /mj/image/:id | 获取 MJ 图片 | 否 | `new-api/relay/relay_mj.go:RelayMidjourneyImage` | id（路径） |
| POST | /mj/submit/imagine | 提交 MJ 绘图任务 | TokenAuth | `new-api/controller/relay.go:RelayMidjourney` | prompt, base64Array 等 |
| POST | /mj/submit/action | 提交 MJ 操作（upscale/variation等） | TokenAuth | `new-api/controller/relay.go:RelayMidjourney` | taskId, customId |
| POST | /mj/submit/shorten | 提交提示词缩短任务 | TokenAuth | `new-api/controller/relay.go:RelayMidjourney` | prompt |
| POST | /mj/submit/modal | 提交 MJ 弹窗任务 | TokenAuth | `new-api/controller/relay.go:RelayMidjourney` | taskId |
| POST | /mj/submit/change | 提交 MJ 变换任务 | TokenAuth | `new-api/controller/relay.go:RelayMidjourney` | taskId, action, index |
| POST | /mj/submit/simple-change | 提交简单变换任务 | TokenAuth | `new-api/controller/relay.go:RelayMidjourney` | content |
| POST | /mj/submit/describe | 提交图片描述任务 | TokenAuth | `new-api/controller/relay.go:RelayMidjourney` | base64, url |
| POST | /mj/submit/blend | 提交图片混合任务 | TokenAuth | `new-api/controller/relay.go:RelayMidjourney` | base64Array |
| POST | /mj/submit/edits | 提交图片编辑任务 | TokenAuth | `new-api/controller/relay.go:RelayMidjourney` | taskId, prompt, mask 等 |
| POST | /mj/submit/video | 提交 MJ 视频任务 | TokenAuth | `new-api/controller/relay.go:RelayMidjourney` | taskId |
| POST | /mj/notify | MJ 回调通知 | TokenAuth | `new-api/controller/relay.go:RelayMidjourney` | 无 |
| GET | /mj/task/:id/fetch | 查询 MJ 任务状态 | TokenAuth | `new-api/controller/relay.go:RelayMidjourney` | id（路径） |
| GET | /mj/task/:id/image-seed | 获取任务图片 seed | TokenAuth | `new-api/controller/relay.go:RelayMidjourney` | id（路径） |
| POST | /mj/task/list-by-condition | 按条件查询 MJ 任务 | TokenAuth | `new-api/controller/relay.go:RelayMidjourney` | ids, channelId 等 |
| POST | /mj/insight-face/swap | 人脸替换 | TokenAuth | `new-api/controller/relay.go:RelayMidjourney` | sourceBase64, targetBase64 |
| POST | /mj/submit/upload-discord-images | 上传图片到 Discord | TokenAuth | `new-api/controller/relay.go:RelayMidjourney` | base64Array |
| GET | /:mode/mj/image/:id | 模式化 MJ 图片获取 | 否 | `new-api/relay/relay_mj.go:RelayMidjourneyImage` | mode, id（路径） |
| POST | /:mode/mj/submit/* | 模式化 MJ 任务提交 | TokenAuth | `new-api/controller/relay.go:RelayMidjourney` | mode（路径），task body |

---

### Suno 音乐生成接口

| 方法 | 路径 | 功能说明 | 认证 | 代码位置 | 主要参数 |
|------|------|----------|------|----------|----------|
| POST | /suno/submit/:action | 提交 Suno 任务 | TokenAuth | `new-api/controller/relay.go:RelayTask` | action（路径），prompt, style 等 |
| POST | /suno/fetch | 批量查询 Suno 任务状态 | TokenAuth | `new-api/controller/relay.go:RelayTaskFetch` | ids |
| GET | /suno/fetch/:id | 查询单个 Suno 任务状态 | TokenAuth | `new-api/controller/relay.go:RelayTaskFetch` | id（路径） |

---

### 视频生成接口

| 方法 | 路径 | 功能说明 | 认证 | 代码位置 | 主要参数 |
|------|------|----------|------|----------|----------|
| POST | /v1/video/generations | 提交视频生成任务 | TokenAuth | `new-api/controller/relay.go:RelayTask` | model, prompt 等 |
| GET | /v1/video/generations/:task_id | 查询视频生成任务状态 | TokenAuth | `new-api/controller/relay.go:RelayTaskFetch` | task_id（路径） |
| POST | /v1/videos/:video_id/remix | 视频混剪任务 | TokenAuth | `new-api/controller/relay.go:RelayTask` | video_id（路径），prompt 等 |
| POST | /v1/videos | 创建视频（OpenAI 兼容） | TokenAuth | `new-api/controller/relay.go:RelayTask` | model, prompt 等 |
| GET | /v1/videos/:task_id | 查询视频任务状态（OpenAI 兼容） | TokenAuth | `new-api/controller/relay.go:RelayTaskFetch` | task_id（路径） |
| GET | /v1/videos/:task_id/content | 代理获取视频内容 | TokenAuth 或 UserAuth | `new-api/controller/video_proxy.go:VideoProxy` | task_id（路径） |
| POST | /kling/v1/videos/text2video | Kling 文本生成视频 | TokenAuth | `new-api/controller/relay.go:RelayTask` | model, prompt 等 |
| POST | /kling/v1/videos/image2video | Kling 图片生成视频 | TokenAuth | `new-api/controller/relay.go:RelayTask` | model, image, prompt 等 |
| GET | /kling/v1/videos/text2video/:task_id | 查询 Kling 文生视频任务 | TokenAuth | `new-api/controller/relay.go:RelayTaskFetch` | task_id（路径） |
| GET | /kling/v1/videos/image2video/:task_id | 查询 Kling 图生视频任务 | TokenAuth | `new-api/controller/relay.go:RelayTaskFetch` | task_id（路径） |
| POST | /jimeng/ | 即梦官方 API（提交/查询任务） | TokenAuth | `new-api/controller/relay.go:RelayTask` | Action（query），body |

---

## 快速查找指南

**修改接口步骤**：
1. 在表格中找到目标接口
2. 查看"代码位置"列，找到文件和函数名（格式：`new-api/文件路径:函数名`）
3. 打开对应文件直接修改

**了解接口详细参数**：
1. 找到代码位置
2. 查看请求体的结构定义（如 `LoginRequest`、`SetupRequest` 等 struct）
3. 查看 `dto/` 目录下的数据传输对象定义

**认证方式说明**：
- `否`：完全公开，无需认证
- `UserAuth`：需要用户登录态（Session Cookie）
- `AdminAuth`：需要管理员角色
- `RootAuth`：需要 Root（超级管理员）角色
- `TokenAuth`：需要 API Token（Bearer Token 或 x-api-key）
- `TokenAuthReadOnly`：只读 API Token

---

## 常用接口快速链接

- 系统状态：`GET /api/status` → `new-api/controller/misc.go:GetStatus`
- 用户登录：`POST /api/user/login` → `new-api/controller/user.go:Login`
- 聊天补全（主要 AI 接口）：`POST /v1/chat/completions` → `new-api/controller/relay.go:Relay`
- 获取令牌列表：`GET /api/token/` → `new-api/controller/token.go:GetAllTokens`
- 创建令牌：`POST /api/token/` → `new-api/controller/token.go:AddToken`
- 获取定价信息：`GET /api/pricing` → `new-api/controller/pricing.go:GetPricing`
- 图像生成：`POST /v1/images/generations` → `new-api/controller/relay.go:Relay`
- 获取所有日志（管理员）：`GET /api/log/` → `new-api/controller/log.go:GetAllLogs`
- 创建渠道（管理员）：`POST /api/channel/` → `new-api/controller/channel.go:AddChannel`
- 系统选项配置（Root）：`PUT /api/option/` → `new-api/controller/option.go:UpdateOption`

---

<!-- 扫描进度 (请勿删除此注释，用于增量更新)
执行次数: 1
扫描时间: 2026-03-02T00:00:00Z
完成状态: 已完成
已扫描目录: router/, controller/
已扫描文件数: 58
下次建议: 无（所有路由文件及控制器已全量扫描）
-->
