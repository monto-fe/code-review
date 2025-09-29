import gitWebhookImg from '@/assets/images/git-webhook.png';
import codeReviewPreviewImg from '@/assets/images/code-review-preview.png';
import gitlabTokenImg from '@/assets/images/gitlab-token.png';

const gitlabConfigDoc = (IP: string, domain: string) => `

## 1. 配置AI模型「若已配置请跳过」

*模型支持情况及Key的获取*

UCloud：https://console.ucloud.cn/modelverse/experience/api-keys
- DeepSeek-V3-0324
- Qwen/Qwen3-Coder
- moonshotai/Kimi-K2-Instruct

DeepSeek：https://platform.deepseek.com/api_keys
- deepseek-chat

Qwen：https://bailian.console.aliyun.com/?tab=model#/api-key
- qwen3-coder-plus

## 2. 配置Gitlab Webhook 配置示例「若已配置可以跳过」

如果Gitlab有的版本支持Group Webhook，推荐在 Group 层级统一配置，在对应项目/项目组/系统钩子中，如果项目组或系统钩子中配置过了，就不要再处理了

配置位置：Settings ➜ Webhooks

添加如下地址：${IP}/webhook/merge

勾选事件类型：**Merge request events**、**Push events**

提交 Merge Request 或推送 Push 时将自动触发 AI 审查。

![Webhook配置](${gitWebhookImg})


## 3. 配置Gitlab Token

**获取Token**

![Gitlab Token](${gitlabTokenImg}) 

**录入Token**

配置地址：${domain}/aicodecheck/GitlabToken


## 4. AI 检查示例

![Webhook触发](${codeReviewPreviewImg})

`;

export { gitlabConfigDoc };