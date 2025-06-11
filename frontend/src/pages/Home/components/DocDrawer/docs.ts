const gitlabConfigDoc = (IP: string) => `
✅ 如果Gitlab有的版本支持Group Webhook，推荐在 Group 层级统一配置。

## 1. 在对应项目中：

配置位置：Settings ➜ Webhooks

添加如下地址：http://${IP}/webhook/merge

勾选事件类型：**Merge request events**

提交 Merge Request 时将自动触发 AI 审查。

## 2. Webhook 配置示意：

![Webhook配置](https://picture.questionlearn.cn/blog/picture/1746626508783.png)

## 3. AI 检查示例：

![Webhook触发](https://picture.questionlearn.cn/blog/picture/1746626303888.png)

---
`;

export { gitlabConfigDoc };