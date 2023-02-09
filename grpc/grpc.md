# grpc

首先 关于 protobuf 的使用
可以阅读这些内容:[https://www.liwenzhou.com/posts/Go/protobuf/]

protobuf 的时间戳类型
在使用之前需要先引入以下定义的包

```protobuf
import "google/protobuf/timestamp.proto"

message HelloRequest{
   string name =1;
   string url =2;
   map<string,string> scop=3;
   google.protobuf.Timestamp addTime=4;
}

```

使用的时候也需要使用包名.的方式
在 go 里面使用:

```go
addTime:=timeStamppb.New(time.Now())
```

## grpc metadata

grpc 让我们可以像本地调用一样实现远程调用，对于每一次远程调用，都可能会产生一些有用的数据，
这些数据可以通过 metadata 来传递，metadata 是以 key-value 的形式存储数据，其中 key 是 string 类型
value 是[]string 类型
metadata 可以使得 client 和 server 端能够为对方提供关于本次调用的一些信息，就像一次 http 请求的 requestHeader 和 ResponseHeader 一样

http 中 header 的生命周期是一次 http 请求，那么 metadata 的生命周期就是一次 rpc 调用
有一些需要注意的点是：metadata 中的键不能以 grpc-开头
metadata 一般是在请求和响应过程中，需要但不是具体业务的信息，比如说身份验证

新建 metadata
MD 实际是 map,key 是 string,value 是[]string

```go
type MD map[string][]string
```

新建的时候可以像创建普通的 map 类型一样使用 new 来创建

```go
md:=metadata.New(map[string][]string{"key1":"jr","key2":"mf"})

//这里key不区分大小写，会统一小写
//需要注意的是 即使是key和value也是使用,隔开
md:=metadata.Pairs(
    "key1","vali",
    "key1","valis", //key1:[]string{"vali","valis"}
    "key2","va3",
)

```

发送 metadata

```go
md:=metadata.Pairs("key","val")

// 新建一个有metadata的context
ctx:=metadata.NewOutgoingContext(context.Background(),md)

//单向rpc
response,err:=client.SomeRPC(ctx,someRequest)
```

接受 matedata

```go
func(s *Server)SomeRPC(ctx context.Context,in *pb.SomeRequest)(*pb.SomeResponse,error){
    md,ok:=metadata.FormIncomingContext(ctx)

}
```

元数据可以存储二进制的数据

```go
md:=metadata.Pairs(
	"key","string value",
	"key-bin",string([]byte{96,102}), //二进制发送前会进行(base64)编码 ，收到后会进行解码
	)
```

## grpc 错误码

go 语言使用的 gRPC status 定义在

```go
import "google.golang.org/grpc/status"
```

rpc 服务应该返回 nil 或来自 status.Status 类型的错误，客户端可以直接访问

创建错误，通常使用 status.New()创建一个 status.Status,通过类型的 Err 将它转换成 error

```go
//创建status.Status
st:=status.New(codes.NotFound,"some description")
err:=st.Err() //转为error类型

// 或者直接使用status.Error
err:=status.Error(codes.NotFound,"some description")
```

为错误添加详细信息
使用 status.WithDetails,它可以添加多个 proto.Message

```go
st := status.New(codes.ResourceExhausted, "Request limit exceeded.")
ds, _ := st.WithDetails(
	// proto.Message
)
return nil, ds.Err()
```

客户端可以通过首先将普通 error 类型转换回 status.Status，然后使用 status.Details 来读取这些详细信息

```go
s := status.Convert(err)
for _, d := range s.Details() {
	// ...
}
```

## 拦截器

gRPC 为在每个 ClientConn/Server 基础上实现和安装拦截器提供了一些简单的 API。 拦截器拦截每个 RPC 调用的执行。
用户可以使用拦截器进行日志记录、身份验证/授权、指标收集以及许多其他可以跨 RPC 共享的功能。
在 gRPC 中，拦截器根据拦截的 RPC 调用类型可以分为两类。第一个是普通拦截器（一元拦截器），它拦截普通 RPC 调用。
另一个是流拦截器，它处理流式 RPC 调用。而客户端和服务端都有自己的普通拦截器和流拦截器类型。
因此，在 gRPC 中总共有四种不同类型的拦截器。
可以将拦截器理解成中间件

客户端拦截器
普通拦截器、一元拦截器

```go
func(ctx context.Context, method string, req, reply interface{}, cc *ClientConn, invoker UnaryInvoker, opts ...CallOption) error
```

一元拦截器的实现通常可以分为三个部分: 调用 RPC 方法之前（预处理）、调用 RPC 方法（RPC 调用）和调用 RPC 方法之后（调用后）。

预处理：用户可以通过检查传入的参数(如 RPC 上下文、方法字符串、要发送的请求和 CallOptions 配置)来获得有关当前 RPC 调用的信息。
RPC 调用：预处理完成后，可以通过执行 invoker 执行 RPC 调用。
调用后：一旦调用者返回应答和错误，用户就可以对 RPC 调用进行后处理。通常，它是关于处理返回的响应和错误的。 若要在 ClientConn 上安装一元拦截器，请使用 DialOptionWithUnaryInterceptor 的 DialOption 配置 Dial 。

流拦截器
StreamClientInterceptor 是客户端流拦截器的类型。它的函数签名是

func(ctx context.Context, desc *StreamDesc, cc *ClientConn, method string, streamer Streamer, opts ...CallOption) (ClientStream, error)
流拦截器的实现通常包括预处理和流操作拦截。

预处理：类似于上面的一元拦截器。
流操作拦截：流拦截器并没有事后进行 RPC 方法调用和后处理，而是拦截了用户在流上的操作。首先，拦截器调用传入的 streamer 以获取 ClientStream，然后包装 ClientStream 并用拦截逻辑重载其方法。
最后，拦截器将包装好的 ClientStream 返回给用户进行操作。
若要为 ClientConn 安装流拦截器，请使用 WithStreamInterceptor 的 DialOption 配置 Dial。

server 端拦截器
服务器端拦截器与客户端类似，但提供的信息略有不同。

普通拦截器/一元拦截器
UnaryServerInterceptor 是服务端的一元拦截器类型，它的函数签名是

func(ctx context.Context, req interface{}, info \*UnaryServerInfo, handler UnaryHandler) (resp interface{}, err error)
服务端一元拦截器具体实现细节和客户端版本的类似。

若要为服务端安装一元拦截器，请使用 UnaryInterceptor 的 ServerOption 配置 NewServer。

流拦截器
StreamServerInterceptor 是服务端流式拦截器的类型，它的签名如下：

```go
func(srv interface{}, ss ServerStream, info *StreamServerInfo, handler StreamHandler) error

```

实现细节类似于客户端流拦截器部分。

若要为服务端安装流拦截器，请使用 StreamInterceptor 的 ServerOption 来配置 NewServer。

go 社区里有一些常用的 grpc 中间件
go-grpc/middleware[https://github.com/grpc-ecosystem/go-grpc-middleware]

## grpc 名称解析

具体内容:[https://www.liwenzhou.com/posts/Go/name-resolving-and-load-balancing-in-grpc/]
名称解析器（name resolver）可以看作是一个 map[service-name][]backend-ip。它接收一个服务名称，并返回后端的 IP 列表。
gRPC 中根据目标字符串中的 scheme 选择名称解析器。

## grpc gateway

## grpc transcoding

## 负载均衡

当并发量上来之后，一般会添加服务器
所谓负载均衡就是控制请求访问哪台服务器，缓解单台服务器的压力
负载均衡的目的是保证服务的高可用性

这里的负载均衡主要是 grpc 中使用负载均衡

负载均衡的策略：

1.集中式的负载均衡

2.进程内负载均衡

3.独立连接负载均衡

常用的负载均衡算法：

1.随机算法
使用完全随机的方式决定由哪个服务器来处理请求，这样每个服务器被选中的概率是一样的
但是缺点是请求满的服务器或者性能差的服务器与性能好的服务器概率一致

2。轮询算法
请求依次分配给每台服务器，优点是实现简单，每台服务器获得请求的概率一致，缺点是不分好赖

3.加权随机算法
在随机的基础上加权重，性能高的权重大，性能小的权重小，是对随机和轮询的改进

4.加权轮询
这种算法和加权随机算法类似，只是在轮询算法的基础上加入了权重，对不同处理能力的服务器分配不同的权重，这样轮询过程中不同服务器被询问到的概率也会不同，很好的对集群性能做了优化。但是这种算法显然也存在着缺点，一个问题就是长链接的维持和命中率，在轮询算法中，即使是相同的请求的处理也不会有特殊处理，依旧会采取轮询的策略，这样显然会造成命中率的下降，即相同的请求会被分配到不同的服务器中。另外，服务器的权重都是静态配置的，当服务器的性能发生变化时应变能力不足，且容灾能力较差，试想，一台高性能的服务器万一出现问题，由于其权重较大，瞬间就会有大量请求失效。

5.hash 算法
那么既然提到了以上两个问题，有没有更好的算法来解决这些问题呢，这时就要提到我们的 hash 算法了。hash 算法大家都不陌生，将请求的 url 或者是 ip 进行 hash，并选取对应数字的服务器进行分配。像下图所示，这种问题很好地解决了命中率和长链接的问题，因为显然相同的 url 或是 ip 会有相同的 hash 值并会被分配到相同的服务器中。但另一个问题并没有解决，就是容灾能力，当一台服务器性能下降甚至宕机的时候，就会有大量的请求被阻塞甚至拒绝。

6.一致性 hash
那么如何提高容灾能力呢，这时就不得不提到我们的一致性 hash 算法，或者说是 hash 环了。hash 环相信大家也不陌生，简单介绍一下 hash 环的工作过程。

(1) 将机器根据 hash 函数映射到环上；

(2) 将数据桶根据 hash 函数映射到环上；

(3) 据数据映射到桶的位置顺时针找到第一台机器将该桶放到该机器上；

(4) 当某台机器坏掉时，类似(3)将存储在该机器上的数据顺时针找到下一台机器；

(5) 当增加机器时，将该机器与前一台机器(逆时针)之间的桶存储在新增机器上并从原来机器上移除。

很显然这样的算法相比 hash 算法大大提高了容灾能力，当某台机器宕机时，相同的请求会继续同样分配到下一台机器中。

7.最小连接数法
最小连接数算法是一种动态负载均衡策略，所谓动态负载均衡策略，其实就是指在分配服务器时会考虑每台服务器当时的状态来选择与哪台服务器建立链接，而在最小连接数算法中，起到决定性作用的正是每台服务器的链接数，每次分配新的请求时，会选取当前连接数最少的服务器。相应的这种算法也有改进版的加权最小连接数算法，既考虑了当前服务器的连接数也考虑了当前服务器的性能，并对其进行加权。

8.最小响应时间算法
最小 RT 算法也是一种很常用的动态负载均衡策略，与最小连接数算法不同的地方就在于选取服务器时的决定性因素变为了每台服务器的 rt（response time，响应时间），这种算法实时考虑了每台服务器当前的状态，如果一台服务器性能很好，但是已经达到了很高的 qps，那么它的 rt 还是会上升，接下来的请求再次分配给它的概率就会变低了，这样就很好地利用了服务器的当前状态来决定请求的合理分配。

grpc 的负载均衡

grpc 中的负载均衡是基于每次调用的，而不是基于每个连接，即使所有的连接都来自一个客户端，也希望在所有的服务器之间负载均衡

grpc-go 内置支持 pick_first(默认值)和 round_robin 两种策略

pick_first 是 gRPC 负载均衡的默认值，因此不需要设置。pick_first 会尝试连接取到的第一个服务端地址，如果连接成功，则将其用于所有 RPC，如果连接失败，则尝试下一个地址(并继续这样做，直到一个连接成功)。因此，所有的 RPC 将被发送到同一个后端。所有接收到的响应都显示相同的后端地址。
round_robin 连接到它所看到的所有地址，并按顺序一次向每个 server 发送一个 RPC。例如，我们现在注册有两个 server，第一个 RPC 将被发送到 server-1，第二个 RPC 将被发送到 server-2，第三个 RPC 将再次被发送到 server-1。

grpc 客户端使用 grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy":"round_robin"}`)设置负载均衡策略

```go
conn, err := grpc.Dial(
	"q1mi:///resolver.liwenzhou.com",
	grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy":"round_robin"}`), // 这里设置初始策略
	grpc.WithTransportCredentials(insecure.NewCredentials()),
)
```

### go 使用 grpc 负载均衡

首先要理解一个东西：
名称解析：name resolving
它会将域名给解析成 ip 地址，并将这些 ip 地址发送到负载均衡器
由负载均衡器来对这些服务进行连接操作

比如有一个域名：http://wwww.examplename.com 被名称解析器解析器解析出了这四个 ip
192.167.12.14,
192.168.13.15,
192.168.23.14,
192.168.24.15
负载均衡器会将请求在这四个服务之间负载均衡

名称解析器使用:
grpc 默认的名称解析从 DNS 解析:

```go
conn, err := grpc.Dial("dns:///localhost:8972",
	grpc.WithTransportCredentials(insecure.NewCredentials()),
)
```

解析的语法：
dns:///localhost:8932

除了从 dns 名称解析 还可以从 consul 做名称解析
github 地址[https://github.com/mbobakov/grpc-consul-resolver]

使用:

```go
package main
import _ "github.com/mbobakov/grpc-consul-resolver"
// ...
conn, err := grpc.Dial(
		// consul服务
		"consul://192.168.1.11:8500/hello?wait=14s",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
```

当然名称解析也可以自定义
有了名称解析之后使用负载均衡策略

```go
conn, err := grpc.Dial(
	"consul://192.168.1.11:8500/hello?wait=14s",
	grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy":"round_robin"}`), // 这里设置初始策略
	grpc.WithTransportCredentials(insecure.NewCredentials()),
)
```

## 服务注册与发现

在微服务的开发过程中，往往服务与服务之间需要互相调用，比如商城的用户服务和库存服务，如果每个服务都在本地维护别的服务的 ip 等信息，一旦某个服务发生了更改，别的服务就需要重新部署

注册中心：每个服务将自己的信息注册到注册中心，当一个服务访问另外的服务的时候，先从服务注册中心拉取别的服务的配置。

服务注册中心的健康检查，所谓的健康检查就是定期检查某个服务是不是通的

consul 支持健康检查 consul 的分布式一致性算法:raft

使用 docker 安装:

```docker
docker pull consul
```

运行：

```docker
docker run -d -p 8500:8500 -p 8300:8300 -p 8301:8301 -p 8302:8302 -p 8600:8600/udp consul consul agent -dev -client=0.0.0.0
```
访问地址:
wsl:172.20.115.6:8500 网页端访问consul
有图形化界面

如果希望一个容器在重启docker的时候能够自动重启 需要以下命令:
```docker
docker container update  --restart=always 容器id
```

这里要注意的是： consul有两个默认端口 一个是8500 这个端口是HTTP的端口 还有一个是8600 这个是dns的端口

访问dns consul提供dns功能，可以让我们通过dig命令行来测试，consul dns端口是8600 命令行:
```bash
dig @172.20.115.6 -8600 consul.service.consul SRV
这里 consul.service.consul 是域名自己起的
```

这里不禁要问：DNS是干嘛的？ 简单来说 我们访问一个页面 比如说:www.taobao.com 访问淘宝的服务器 实际上是访问不到的 当我们访问这个域名的时候，浏览器会拿着这个域名去dns去查询 dns查询该域名的ip地址，将ip地址返回给浏览器 浏览器拿着ip地址访问服务

我们在windos system32 driver 找到host 将ip地址注册到里面 浏览器就可以不走DNS，直接访问服务

关于注册中心与dns
一个服务在调用另外一个服务的时候， 我们需要知道另外一个服务的URL接口 如果是第三方的来调用 走网关 如果我们将注册中心伪装成一个DNS 一个服务在注册中心注册的时候，就生成一个域名 那么网关只需要开放一个DNS查询功能就可以了

### 1.consul服务注册

