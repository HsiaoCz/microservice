## 链路追踪

为什么需要链路追踪？

一个请求可能会经过多个服务，但是这个请求失败，
就需要排查是哪个服务出现了问题，这种排查过程往往很慢很繁琐

jaeger

安装 jaeger

```docker
docker run --rm --name jaeger -p6831:6831/udp -p16686:16686 jaegertracing/all-in-one:latest
```

然后访问：
172.21.17.5:16686

不过这种启动方式，数据在内存中，一旦重新启动，数据就没了

### jaeger 的架构

我们使用 jaeger 的 client，对应每个语言都会提供一个 client
使用 client 将信息以 UDP 的形式发送到本地 jaeger-agent
jaeger 的代理，jaeger 的代理会将信息推到 jaeger-collector(接收器)
接收器会将信息存储到数据库中
我们在网页上看到的 jaeger UI 它和一个 jaeger-query（查询组件）交互
这个查询组件会从数据库中将信息查询出来

在使用的时候，只需要调用 jaeger 的 client 就可以了

jaeger 组成
jaeger Client 为不同语言实现了符合 OpenTracing 标准的 SDK,应用程序通过 API 写入数据
client library 把 trace 信息按照应用程序指定的采样策略传递给 jaeger-agent

agent 它是一个监听在 udp 端口上接收 span 数据的网络守护进程，它会将数据批量发送给 collector,它被设计成一个基础组件，部署到所有的宿主机上。
agent 将 client library 和 collector 解耦，为 client library 屏蔽了路由和发现 collector 的细节
collector 接收 jaeger-agent 发送来的数据，然后将数据写入后端存储，collector 被设计成无状态的组件，可以运行任意数量的 jaeger-collector

data-store 后端存储被设计成一个可插拔的组件，支持数据写入 es 等
Query:接收查询请求，然后从后端存储系统中检索 trace 并通过 UI 进行展示，Query 是无状态的，可以启动多个实例

分布式系统的链路追踪核心要点就三个：代码埋点，数据存储，查询展示

### python 语言集成 jaeger

先下载:

```python
pip install jaeger-client
```

先弄清楚一些东西：
在分布式链路跟踪中有两个重要的概念：跟踪（trace）和 跨度（ span）。trace 是请求在分布式系统中的整个链路视图，span 则代表整个链路中不同服务内部的视图，span 组合在一起就是整个 trace 的视图。

在整个请求的调用链中，请求会一直携带 traceid 往下游服务传递，每个服务内部也会生成自己的 spanid 用于生成自己的内部调用视图，并和 traceid 一起传递给下游服务。

traceid 在请求的整个调用链中始终保持不变，所以在日志中可以通过 traceid 查询到整个请求期间系统记录下来的所有日志。请求到达每个服务后，服务都会为请求生成 spanid，而随请求一起从上游传过来的上游服务的 spanid 会被记录成 parent-spanid 或者叫 pspanid。当前服务生成的 spanid 随着请求一起再传到下游服务时，这个 spanid 又会被下游服务当做 pspanid 记录。

### go 集成 jaeger

```go
	"github.com/uber/jaeger-client-go"
	jaegercfg "github.com/uber/jaeger-client-go/config"
```
