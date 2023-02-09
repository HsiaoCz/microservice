# grpc

首先 关于protobuf的使用
可以阅读这些内容:[https://www.liwenzhou.com/posts/Go/protobuf/]

protobuf的时间戳类型
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
在go里面使用:
```go
addTime:=timeStamppb.New(time.Now())
```

## grpc metadata

grpc让我们可以像本地调用一样实现远程调用，对于每一次远程调用，都可能会产生一些有用的数据，
这些数据可以通过metadata来传递，metadata是以key-value的形式存储数据，其中key是string类型
value是[]string类型
metadata可以使得client和server端能够为对方提供关于本次调用的一些信息，就像一次http请求的requestHeader和ResponseHeader一样

http中header的生命周期是一次http请求，那么metadata的生命周期就是一次rpc调用
有一些需要注意的点是：metadata 中的键不能以grpc-开头
metadata一般是在请求和响应过程中，需要但不是具体业务的信息，比如说身份验证


新建metadata
MD实际是map,key是string,value是[]string
```go
type MD map[string][]string
```

新建的时候可以像创建普通的map类型一样使用new来创建
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

发送metadata
```go
md:=metadata.Pairs("key","val")

// 新建一个有metadata的context
ctx:=metadata.NewOutgoingContext(context.Background(),md)

//单向rpc
response,err:=client.SomeRPC(ctx,someRequest)
```

接受matedata
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

## grpc错误码
go语言使用的gRPC status定义在
```go
import "google.golang.org/grpc/status"
```
rpc服务应该返回nil或来自status.Status类型的错误，客户端可以直接访问

创建错误，通常使用status.New()创建一个status.Status,通过类型的Err将它转换成error
```go
//创建status.Status
st:=status.New(codes.NotFound,"some description")
err:=st.Err() //转为error类型

// 或者直接使用status.Error
err:=status.Error(codes.NotFound,"some description")
```

为错误添加详细信息
使用status.WithDetails,它可以添加多个proto.Message
```go
st := status.New(codes.ResourceExhausted, "Request limit exceeded.")
ds, _ := st.WithDetails(
	// proto.Message
)
return nil, ds.Err()
```
客户端可以通过首先将普通error类型转换回status.Status，然后使用status.Details来读取这些详细信息
```go
s := status.Convert(err)
for _, d := range s.Details() {
	// ...
}
```

## 拦截器

gRPC 为在每个 ClientConn/Server 基础上实现和安装拦截器提供了一些简单的 API。 拦截器拦截每个 RPC 调用的执行。
用户可以使用拦截器进行日志记录、身份验证/授权、指标收集以及许多其他可以跨 RPC 共享的功能。
在 gRPC 中，拦截器根据拦截的 RPC 调用类型可以分为两类。第一个是普通拦截器（一元拦截器），它拦截普通RPC 调用。
另一个是流拦截器，它处理流式 RPC 调用。而客户端和服务端都有自己的普通拦截器和流拦截器类型。
因此，在 gRPC 中总共有四种不同类型的拦截器。
可以将拦截器理解成中间件

客户端拦截器
普通拦截器、一元拦截器
```go
func(ctx context.Context, method string, req, reply interface{}, cc *ClientConn, invoker UnaryInvoker, opts ...CallOption) error
```
一元拦截器的实现通常可以分为三个部分: 调用 RPC 方法之前（预处理）、调用 RPC 方法（RPC调用）和调用 RPC 方法之后（调用后）。

预处理：用户可以通过检查传入的参数(如 RPC 上下文、方法字符串、要发送的请求和 CallOptions 配置)来获得有关当前 RPC 调用的信息。
RPC调用：预处理完成后，可以通过执行invoker执行 RPC 调用。
调用后：一旦调用者返回应答和错误，用户就可以对 RPC 调用进行后处理。通常，它是关于处理返回的响应和错误的。 若要在 ClientConn 上安装一元拦截器，请使用DialOptionWithUnaryInterceptor的DialOption配置 Dial 。

流拦截器
StreamClientInterceptor是客户端流拦截器的类型。它的函数签名是

func(ctx context.Context, desc *StreamDesc, cc *ClientConn, method string, streamer Streamer, opts ...CallOption) (ClientStream, error)
流拦截器的实现通常包括预处理和流操作拦截。

预处理：类似于上面的一元拦截器。
流操作拦截：流拦截器并没有事后进行 RPC 方法调用和后处理，而是拦截了用户在流上的操作。首先，拦截器调用传入的streamer以获取 ClientStream，然后包装 ClientStream 并用拦截逻辑重载其方法。
最后，拦截器将包装好的 ClientStream 返回给用户进行操作。
若要为 ClientConn 安装流拦截器，请使用WithStreamInterceptor的 DialOption 配置 Dial。

server端拦截器
服务器端拦截器与客户端类似，但提供的信息略有不同。

普通拦截器/一元拦截器
UnaryServerInterceptor是服务端的一元拦截器类型，它的函数签名是

func(ctx context.Context, req interface{}, info *UnaryServerInfo, handler UnaryHandler) (resp interface{}, err error)
服务端一元拦截器具体实现细节和客户端版本的类似。

若要为服务端安装一元拦截器，请使用 UnaryInterceptor 的ServerOption配置 NewServer。

流拦截器
StreamServerInterceptor是服务端流式拦截器的类型，它的签名如下：
```go
func(srv interface{}, ss ServerStream, info *StreamServerInfo, handler StreamHandler) error

```
实现细节类似于客户端流拦截器部分。

若要为服务端安装流拦截器，请使用 StreamInterceptor 的ServerOption来配置 NewServer。

go社区里有一些常用的grpc中间件
go-grpc/middleware[https://github.com/grpc-ecosystem/go-grpc-middleware]

## grpc 名称解析

具体内容:[https://www.liwenzhou.com/posts/Go/name-resolving-and-load-balancing-in-grpc/]
名称解析器（name resolver）可以看作是一个 map[service-name][]backend-ip。它接收一个服务名称，并返回后端的 IP 列表。
gRPC中根据目标字符串中的scheme选择名称解析器。