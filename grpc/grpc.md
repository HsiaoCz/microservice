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


## 拦截器

所谓的拦截器就是在还没有进入请求之前对每个请求先做一遍预处理
grpc内置了拦截器的配置
有一个库:go-grpc/middleware
