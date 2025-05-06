# code review

## 快速启动服务

> 服务默认使用docker启动mariadbAPI服务

[修改代码，通过前端录入] 1. 修改backendend/.env.production，配置deepseek

1. cd code-review && docker compose up -d 即可启动服务

2. 打开web端，在web上中录入gitlab的token

登陆账号和密码：admin/12345678

登陆前端，录入token和ai token

3. 在gitlab的webhook中配置对应的url

如果gitlab支持Group的webhook，可以为Group配置webhook

## 如果使用已有的数据库，需要修改对应的配置信息

### 使用mysql

## 重要改动
- 减少操作难度，将ai的配置录入数据库，读取和切换
- 支持不同的gitlab url
- 目前支持UCloud的模型，需要支持多种大模型，可以进行切换【公有的大模型和私有的大模型】

## 待开发
- 支持录入projectId和mergeId即可进行ai检查
- 支持push事件
- 支持推送机器人信息
- 支持完整的docker镜像启动服务，配置文件通过挂载的方式获取
- 提供内部版本（.gitlab-ci.yml）
