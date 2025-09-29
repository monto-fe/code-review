# 🚦 Code Review 项目

<div align="center">
  <h3>基于 AI 的自动化代码审查系统</h3>
  <p>支持 GitLab Merge Request 审查，可通过 Docker Compose 一键启动服务与数据库</p>
</div>

---

## 🚀 快速启动

> ✅ 支持 **MySQL**（推荐版本：`8.0`）与 **MariaDB**（推荐版本：`10.5`）

**环境依赖**

- docker && docker compose > 2.0
- git
- 虚拟机配置最小配置
  - 支持RAG增强检索：4C16G100G
  - 不需要RAG增强检索：2C8G50G

### 1️⃣ 克隆项目

```bash
git clone https://github.com/monto-fe/code-review.git
cd code-review
```

### 2️⃣ 配置环境变量

> 假设服务器 IP 为：`192.168.1.1`

#### 方式一：使用自有数据库

```bash
# 先在数据库中执行初始化 SQL 脚本
cd mysql/init

# 编辑环境变量
vi .env

# 服务启动+RAG
docker compose up -d

# 不启动RAG服务
docker compose up -d frontend backend
```

#### 方式二：使用内置 MySQL

```bash
cp .env.mysql .env
vi .env  # 修改数据库密码等必要配置
docker compose -f docker-compose.mysql.yml up -d
```

#### 方式三：使用内置 MariaDB

```bash
cp .env.mariadb .env
vi .env
docker compose -f docker-compose.mariadb.yml up -d
```

### 3️⃣ 访问控制台

在浏览器中访问：

> 🌐 地址：[http://192.168.1.1:9003](http://192.168.1.1:9003)
> 🔐 默认账号密码：admin / 12345678

#### 控制台示意图：

![控制台](https://github.com/user-attachments/assets/4100582a-d69e-49cf-826d-b8e3d32a4d79)

首次登录请配置以下内容：

* ✅ AI 模型 Key（支持主流模型）
* ✅ GitLab Token（需具备完整权限）
* ✅ 企业微信机器人（可选）

---

## 📡 配置 Webhook（GitLab）

### 方法一：直接在管理员Gitlab账户中配置系统钩子

### 方法二：

在对应项目中：

> `Settings` ➜ `Webhooks`

添加如下地址：

```
http://192.168.1.1:9000/v1/webhook/merge
```

> ✅ 如果 GitLab 支持 Group Webhook，推荐在 Group 层级统一配置。

勾选事件类型：**Merge request events**

提交 Merge Request 时将自动触发 AI 审查。

#### Webhook 配置示意：

![Webhook配置](https://github.com/user-attachments/assets/3615a38d-53f2-4acd-9ea0-7443e4bea8bd)

#### AI 检查示例：

![Webhook触发](https://github.com/user-attachments/assets/e4bcbce3-4866-45af-924a-7dd634f0696e)

---

## ❓ 常见问题

| 问题      | 解决方案                                     |
| ------- | ---------------------------------------- |
| 镜像拉取失败  | 检查网络连接，或替换为国内镜像源                         |
| 启动后无法访问 | 检查服务器防火墙，确保开放端口 `9000` 与 `9003`          |
| 无法连接数据库 | 数据库可能启动较慢，执行 `docker compose restart` 重试 |

---

## 📞 联系方式

欢迎添加微信，拉群沟通。
微信备注：CodeReview

<img width="686" height="940" alt="微信二维码" src="https://github.com/user-attachments/assets/46a174fb-ee49-4a7a-b576-62fd47292ae1" />

欢迎提Bug/问题：[issue](https://github.com/monto-fe/code-review/issues)

## 📝 更新日志

* ✅ 支持自定义 GitLab 地址
* ✅ 支持多个 GitLab Token
* ✅ 支持 Merge request 和 Push 事件触发webhook
* ✅ 支持 webhook 推送 AI 检查结果
* ✅ 提供快速部署方案
* ✅ 集成 DeepSeek、UCloud 模型，支持切换私有/公有大模型
* ✅ 提供一键启动的 Docker 镜像（支持挂载配置）
* ✅ 优化 GitLab Token 添加/更新流程，自动刷新模型缓存
* ✅ 支持自定义提示词检测代码
* ✅ 支持RAG检索代码仓库
* ✅ 支持行级评论【版本测试16.11.5】
* ✅ 支持Qwen大模型
* ✅ 已测试gitlab版本：12.9.2、16.11.5，其他版本如果出现问题可以反馈

---

## 🔭 规划中功能

* [ ] 支持通过 `projectId + mergeId` 直接发起 AI 审查
* [ ] 支持AI模型精确到对应的Token，也支持默认
* [ ] 支持读取回复评论信息
* [ ] 完成控制台权限管理能力
* [ ] 支持Kimi k2模型
* [ ] 根据description描述信息，搜索对应的需求注意事项，进行检测