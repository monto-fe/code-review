# code review

## 快速启动服务

> 服务默认使用docker启动mariadb API服务

1. 克隆项目
```
git clone https://github.com/richLpf/code-review.git

```
   
2. 修改控制台请求的api地址
```
cd code-review/frontend

vi .env.production

# api接口域名，假如部署的机器ip为192.168.1.1
VITE_APP_APIHOST = http://192.168.1.1:9000/v1
```

3. 启动服务

```
cd code-review && docker compose up -d 

```

4. 打开web端

访问地址：http://192.168.1.1:9003
账号密码：admin/12345678

![控制台](https://picture.questionlearn.cn/blog/picture/1746626120106.png)

添加AI模型key（特定Key）和Gitlab token(Token需要所有权限)

5. 在gitlab的项目的setting => webhook中配置对应的url

> 如果gitlab支持Group的webhook，可以为Group配置webhook

webhook地址：http://192.168.1.1:9000/v1/webhook/merge

![webhook配置](https://picture.questionlearn.cn/blog/picture/1746626508783.png)

目前只支持`Merge request events`方式触发

提交merge，将会触发AI检查

![webhook](https://picture.questionlearn.cn/blog/picture/1746626303888.png)


## 常见问题

- 镜像拉不下来
  - 需要外网或修改镜像地址
- 启动后访问失败
  - 检查防火墙是否开启9000和9003端口
- 启动后连接数据库失败
  - 可以重启一下docker compose restart,可能数据库启动慢，导致服务没有正常连接

## 如果使用已有的数据库，需要修改对应的配置信息

### 使用mysql

需要对backend项目进行简单修改，待同步

## 重要改动
- 支持不同的gitlab url
- 目前支持UCloud的模型，需要支持多种大模型，可以进行切换【公有的大模型和私有的大模型】

## 待开发
- 支持录入projectId和mergeId即可进行ai检查
- 支持push事件
- 支持推送机器人信息
- 支持完整的docker镜像启动服务，配置文件通过挂载的方式获取
- 提供内部版本（.gitlab-ci.yml）
