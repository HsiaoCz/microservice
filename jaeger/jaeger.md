## 链路追踪

为什么需要链路追踪？

一个请求可能会经过多个服务，但是这个请求失败，
就需要排查是哪个服务出现了问题，这种排查过程往往很慢很繁琐

jaeger

安装jaeger

```docker
docker run --rm --name jaeger -p6831:6831/udp -p16686:16686 jaegertracing/all-in-one:latest
```

然后访问：
172.21.17.5:16686

不过这种启动方式，数据在内存中，一旦重新启动，数据就没了

### jaeger的架构

我们使用jaeger的client，对应每个语言都会提供一个client
使用client将信息以UDP的形式发送到本地jaeger-agent
jaeger的代理，jaeger的代理会将信息推到jaeger-collector(接收器)
接收器会将信息存储到数据库中
我们在网页上看到的jaeger UI它和一个jaeger-query（查询组件）交互
这个查询组件会从数据库中将信息查询出来

在使用的时候，只需要调用jaeger 的client就可以了

jaeger 组成
jaeger Client 为不同语言实现了符合OpenTracing标准的SDK,应用程序通过API写入数据
client library把trace信息按照应用程序指定的采样策略传递给jaeger-agent

agent 它是一个监听在udp端口上接收span数据的网络守护进程，它会将数据批量发送给collector,它被设计成一个基础组件，部署到所有的宿主机上。
agent将client library和collector解耦，为client library屏蔽了路由和发现collector的细节
collector 接收jaeger-agent发送来的数据，然后将数据写入后端存储，collector被设计成无状态的组件，可以运行任意数量的jaeger-collector

data-store 后端存储被设计成一个可插拔的组件，支持数据写入es等
Query:接收查询请求，然后从后端存储系统中检索trace并通过UI进行展示，Query是无状态的，可以启动多个实例

分布式系统的链路追踪核心要点就三个：代码埋点，数据存储，查询展示

### python语言集成jaeger

```python

```

