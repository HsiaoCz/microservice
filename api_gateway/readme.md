## api 网关

api 网关应该有什么功能？

1.服务的路由：动态路由，负载均衡 2.服务发现 3.限流，熔断，降级 4.流量的管理：黑白名单，反爬的策略

目前的网关有很多
api 网关对性能的要求很高

## kong 的安装和配置

kong 是一个开源的 api 网关，它是一个针对 Api 的一个管理工具，

安装前置的 postgres

```docker
docker run -d --name kong-database -p 5432:5432 -e "POSTGRES_USER=kong"
-e "POSTGRES_DB=kong" -e "POSTGRES_PASSWORD=kong" -e "POSTGRES_DB=kong" postgres:12
```

运行 kong

```docker
docker run --rm -e "KONG_DATABASE=postgres" -e "KONG_PG_HOST=172.24.107.79"
-e "KONG_PG_PASSWORD=kong" -e "POSTGRES_USER=kong" -e "KONG_CASSANDRA_CONTACT_POINTS=kong-database"
kong kong migrations bootstarp
```
