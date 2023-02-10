## 分布式的配置中心

我们之前的开发过程中 使用的配置都是本地配置

本地配置在微服务中就不好用了，服务部署在不同的地方
如果我们需要修改配置，就得在不同的地方去挨个修改每个配置文件

这时候我们就需要将配置文件集中管理

Apollo（阿波罗）是携程开源的一款可靠的分布式配置管理中心，它能够集中化管理应用不同环境、不同集群的配置，配置修改后能够实时推送到应用端，并且具备规范的权限、流程治理等特性，适用于微服务配置管理场景。

apollo 官方提供了方便学习使用的 docker-quick-start 环境

```bash
git clone https://github.com/apolloconfig/apollo.git

cd apollo/scripts/docker-quick-start/

cd到docker-quick-start目录下执行
docker-compose up即可启动
```

启动之后访问：localhost:8070 查看 apollo 管理后台
登录的用户名和密码:apollo admin

Apollo 支持 4 个维度管理 Key-Value 格式的配置：

application (应用)
environment (环境)
cluster (集群)
namespace (命名空间)

这里可以看 apollo 使用指南

### 1、apollo 使用指南

apollo 使用指南[https://www.apolloconfig.com/#/zh/usage/apollo-user-guide]

普通应用:
独立运行的程序，比如 web 程序，带有 main 函数的程序

公共组件：
发布的类库，不能独立运行的程序，比如 java 的 jar 包
.Net 的 dll 文件

### 2、go 接入 apollo

使用 go 接入 apollo 有很多库
这里使用:
[https://github.com/philchia/agollo]文档看这

```go
go get -u github.com/philchia/agollo/v4
```

还可以使用:

```go
go get -u github.com/shima-park/agollo
```

[https://github.com/shima-park/agollo]看文档在这里
