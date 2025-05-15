# Code Review 项目

基于 AI 的自动化代码审查系统，支持 GitLab Merge Request 审查，当前默认通过 Docker 启动后端服务和数据库。

---

## 🚀 快速启动

> 默认使用 Docker 启动 MariaDB 与 API 服务。

### 1. 克隆项目

```bash
git clone https://github.com/richLpf/code-review.git
```

---

### 2. 修改前端 API 地址

进入前端目录，配置 API 地址为你部署的服务器 IP：

```bash
cd code-review/frontend
vi .env.production
```

修改如下字段（以 `192.168.1.1` 为例）：

```env
VITE_APP_APIHOST=http://192.168.1.1:9000/v1
```

---

### 3. 启动服务

```bash
cd code-review
docker compose up -d
```

---

### 4. 打开控制台页面

在浏览器中访问：

> 📍 地址：[http://192.168.1.1:9003](http://192.168.1.1:9003)
> 🔐 账号密码：admin / 12345678

控制台界面示例：

![控制台](https://picture.questionlearn.cn/blog/picture/1746626120106.png)

* 进入后请添加：

  * AI 模型 Key（支持特定模型）
  * GitLab Token（需要具有全部权限）

---

### 5. 配置 GitLab Webhook

前往对应项目：

> `Settings` ➜ `Webhooks`

添加如下地址：

```
http://192.168.1.1:9000/v1/webhook/merge
```

> ✅ 如 GitLab 支持 Group Webhook，也可在 Group 层级统一配置。

勾选事件类型：**Merge request events**

* 当提交 Merge Request 时，AI 自动触发审查。

界面示例：

![Webhook配置](https://picture.questionlearn.cn/blog/picture/1746626508783.png)

AI 检查示例：

![Webhook触发](https://picture.questionlearn.cn/blog/picture/1746626303888.png)

---

## 🛠 常见问题

| 问题      | 解决方案                                     |
| ------- | ---------------------------------------- |
| 镜像拉取失败  | 检查是否有外网，或更换镜像源                           |
| 启动后无法访问 | 检查防火墙是否开放 `9000` 与 `9003` 端口             |
| 无法连接数据库 | 尝试执行 `docker compose restart`（数据库可能启动较慢） |

---

## 🔧 使用MySQL数据库启动

使用mysql数据库启动。

```
cd code-review
git checkout feature/mysql
cd docker-compose -f docker-compose.mysql.yml up -d
```

## 🔧 使用已有数据库（MySQL）

如果自己有mysql数据库，则可以使用以下命令启动。

数据库需要导入预制的Sql脚本。

```
cd code-review
git checkout feature/mysql
cd docker-compose -f docker-compose.without.mysql.yml up -d
```

修改配置文件`backend/ecosystem.config.js`

```
env_production: {
    PORT: 9000,
    NODE_ENV: 'production',
    DOMAIN: 'http://localhost:9000',
    DB_HOST: "your_mysql_host",
    DB_PORT: 3306,
    DB_USER: "root",
    DB_PASSWORD: "mysql123456",
    DB_DATABASE: "ucode_review"
}
```
---

## 📌 重要变更记录

* ✅ 支持自定义 GitLab 地址
* ✅ 集成 UCloud 模型，未来支持多种大模型（支持公有/私有模型切换）

---

## 🔭 开发中功能（规划）

* [ ] 支持仅输入 `projectId` 和 `mergeId` 即可进行 AI 审查
* [ ] 支持监听 GitLab Push 事件
* [ ] 支持消息推送机器人通知
* [ ] 提供一键启动的 Docker 镜像版本（挂载配置）
* [ ] 提供内部 CI/CD 集成示例（如 `.gitlab-ci.yml`）