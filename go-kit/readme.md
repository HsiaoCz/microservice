## go-kit

go kit 是一个 golang 的工具集，可以帮助构建微服务

go-kit 构建的服务分为三层： 1.传输层(Transport layer) 2.端点层(Endpoint layer) 3.服务层(service layer)

请求从第一层进入服务，向下流到第三层，响应则相反

### transports

传输域绑定到具体的传输协议,如 HTTP 或者是 gRPC，可以在单个微服务中支持原有的 HTTP API 或者新增的 RPC 服务

当实现 REST 式的 HTTP API 时，你的路由是在 HTTP 传输中定义的。最常见的路由定义在 HTTP 路由器函数中，如下所示:

```go
r.Methods("POST").Path("/profiles/").Handler(httptransport.NewServer(
		e.PostProfileEndpoint,
		decodePostProfileRequest,
		encodeResponse,
		options...,
))
```

### Endpoints

端点就像控制器上的动作/处理程序; 它是安全性和抗脆弱性逻辑的所在。如果实现两种传输(HTTP 和 gRPC) ，则可能有两种将请求发送到同一端点的方法。

### Services

服务（指 Go kit 中的 service 层）是实现所有业务逻辑的地方。服务层通常将多个端点粘合在一起。在 Go kit 中，服务层通常被抽象为接口，这些接口的实现包含业务逻辑。Go kit 服务层应该努力遵守整洁架构或六边形架构。也就是说，业务逻辑不需要了解端点（尤其是传输域）概念：你的服务层不应该关心 HTTP 头或 gRPC 错误代码。

### Middlewares

Go kit 试图通过使用中间件（或装饰器）模式来执行严格的关注分离（separation of concerns）。中间件可以包装端点或服务以添加功能，比如日志记录、速率限制、负载平衡或分布式跟踪。围绕一个端点或服务链接多个中间件是很常见的。

### demo1

假如有个需求，简单无比的需求：
传入用户的 ID,来获取用户的用户名

先创建 service 层，再创建 endpoints 层，最后创建 transport 层

service 层定义业务类，接口

```go
package service

type IUserService interface {
	GetName(userID int) string
}

type UserService struct{}

func (u *UserService) GetName(userID int) string {
	if userID == 101 {
		return "bob"
	}
	return "alex"
}
```

创建 Endpoint
用来定义 Request 和 Response 格式,并可以使用装饰器包装函数，以此来实现各种中间件嵌套

```go
package service

type UserRequest struct {
	Uid int `json:"uid"`
}

type UserResponse struct {
	Result string `json:"result"`
}
func GenUserEndpoint(UserService IUserService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		// 使用类型断言获取请求
		r := request.(UserRequest)
		result := UserService.GetName(r.Uid)
		return UserResponse{Result: result}, nil
	}
}
```

最后创建transport