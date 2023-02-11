## 熔断、限流、降级

服务雪崩

服务雪崩可以分为三个阶段:
1、服务提供者不可用
2、重试加大请求流量
3、服务调用者不可用

服务雪崩的每个阶段都可以由不同的原因造成
总结如下:
1、服务不可用:硬件故障、程序 bug、缓存击穿、用户大量请求
2、重试加大流量：用户重试、代码逻辑重试
3、服务调用者不可用：同步等待造成的资源耗尽

应对策略：
1、应用扩容：增加机器数量、升级硬件规格
2、流控：限流，关闭重试
3、缓存：缓存预加载
4、服务降级：服务接口拒绝服务、页面拒绝服务、延迟持久化、随机拒绝服务
5、服务熔断

所谓限流，比如说我的网站本来只能允许 1k 的并发访问
但是现在 咔来了 2k，那么就要考虑限流
可以直接拒绝访问，也可以让用户排队访问
这么一来用户的体验就降级了

熔断：比如服务 A 访问服务 B 服务，这个时候 B 服务很慢，B 服务压力过大，导致了出现了不少错误的情况
调用方很容易出现一个问题：每次都超时 2k,如果这个时候数据库出现了问题，超时重试，网络流量 2k 直接变成了 3k
这就让原本满负荷的 b 服务雪上加霜，如果这个时候调用方有一种机制：比如说 1、发现了大部分请求很慢 50%的服务都很慢
或者发现 50%的服务都发生了错误，这个时候就可以熔断 就像保险丝一样，断开服务

限流熔断都会导致服务降级

## 熔断限流技术

Hystrix netflix 开源的熔断技术
Sentinel 阿里开源的熔断技术

Sentinel 基于信号量隔离，熔断降级的策略 基于响应时间或失败比率 基于滑动窗口的实时指标实现
限流：基于 QPS，支持基于调用关系的限流

### sentinel 限流基于 qps 限流

```go
package main

// sentinel 基于qps的限流
// 所谓qps就是每秒钟请求的通过数
import (
	"fmt"
	"log"

	sentinel "github.com/alibaba/sentinel-golang/api"
	"github.com/alibaba/sentinel-golang/core/base"
	"github.com/alibaba/sentinel-golang/core/flow"
)

func main() {
	// 先初始化
	err := sentinel.InitDefault()
	if err != nil {
		log.Fatalf("初始化异常:%v\n", err)
	}
	// 配置限流规则
	// 	Resource：资源名，即规则的作用目标。
	// TokenCalculateStrategy: 当前流量控制器的Token计算策略。Direct表示直接使用字段 Threshold 作为阈值；WarmUp表示使用预热方式计算Token的阈值。
	// ControlBehavior: 表示流量控制器的控制策略；Reject表示超过阈值直接拒绝，Throttling表示匀速排队。
	// Threshold: 表示流控阈值；如果字段 StatIntervalInMs 是1000(也就是1秒)，那么Threshold就表示QPS，流量控制器也就会依据资源的QPS来做流控。
	// RelationStrategy: 调用关系限流策略，CurrentResource表示使用当前规则的resource做流控；AssociatedResource表示使用关联的resource做流控，关联的resource在字段 RefResource 定义；
	// RefResource: 关联的resource；
	// WarmUpPeriodSec: 预热的时间长度，该字段仅仅对 WarmUp 的TokenCalculateStrategy生效；
	// WarmUpColdFactor: 预热的因子，默认是3，该值的设置会影响预热的速度，该字段仅仅对 WarmUp 的TokenCalculateStrategy生效；
	// MaxQueueingTimeMs: 匀速排队的最大等待时间，该字段仅仅对 Throttling ControlBehavior生效；
	// StatIntervalInMs: 规则对应的流量控制器的独立统计结构的统计周期。如果StatIntervalInMs是1000，也就是统计QPS。
	// 这里特别强调一下 StatIntervalInMs 和 Threshold 这两个字段，这两个字段决定了流量控制器的灵敏度。
	// 以 Direct + Reject 的流控策略为例，流量控制器的行为就是在 StatIntervalInMs 周期内，允许的最大请求数量是Threshold。
	// 比如如果 StatIntervalInMs 是 10000，Threshold 是10000，那么流量控制器的行为就是控制该资源10s内运行最多10000次访问。
	_, err = flow.LoadRules([]*flow.Rule{
		{
			Resource:               "some-test", // 规则的作用目标，资源名
			TokenCalculateStrategy: flow.Direct, // 当前流量控制器的Token计算策略
			ControlBehavior:        flow.Reject, // 表示流量控制器的控制策略
			Threshold:              10,          // 表示流量阈值 1s有十个流量进来
			StatIntervalInMs:       1000,        // 规则对应的流量控制器的独立统计结构的统计周期
		},
	})
    // StatIntervalInMs和Threshold加在一起表示单位时间内可以有多少个流量进来
    // StatIntervalInMs 1000代表1s  Threshold 10表示10个流量
	if err != nil {
		log.Fatalf("加载规则失败:%v", err)
	}

	// 使用 sentinel.Entry表示一个流控的入口点
	// 第一个参数是资源名 表示入口点使用资源名的规则做检查
	// sentinel.WithTrafficType 配置入口或出口的流量控制
	// base.Inbound表示入口的流量控制
	// 这里模拟一下流量
	for i := 0; i < 12; i++ {
		e, b := sentinel.Entry("some-test", sentinel.WithTrafficType(base.Inbound))
		if b != nil {
			fmt.Println("限流了")
		} else {
			fmt.Println("检查通过")
			e.Exit()
		}
	}
}
```

### sentinel 预热和冷启动

// TokenCalculateStrategy: 当前流量控制器的 Token 计算策略。Direct 表示直接使用字段 Threshold 作为阈值；WarmUp 表示使用预热方式计算 Token 的阈值。
// ControlBehavior: 表示流量控制器的控制策略；Reject 表示超过阈值直接拒绝，Throttling 表示匀速排队。

warm_up 表示预热的方式控制流量，所谓的预热：
WarmUp 方式，即预热/冷启动方式。当系统长期处于低水位的情况下，当流量突然增加时，直接把系统拉升到高水位可能瞬间把系统压垮。通过"冷启动"，让通过的流量缓慢增加，在一定时间内逐渐增加到阈值上限，给冷系统一个预热的时间，避免冷系统被压垮。这块设计和 Java 类似，可以参考限流-冷启动文档

简单来说，当我们配置了每秒钟的流量阈值，而平常的流量处在低水平，如果突然流量上来了，系统直接拉到最高水平，显然对系统是不好的
冷启动的意思就是让流量缓慢攀升到最高水平，而不至于让系统压力过大

```go
// 先初始化
	err := sentinel.InitDefault()
	if err != nil {
		log.Fatalf("初始化异常:%v\n", err)
	}
	// TokenCalculateStrategy: 这里我们配置成冷启动
	// ControlBehavior:        flow.Reject, // 表示流量控制器的控制策略,当超出阈值之后怎么办，flow.Reject
	// 表示直接拒绝
	_, err = flow.LoadRules([]*flow.Rule{
		{
			Resource:               "some-test", // 规则的作用目标，资源名
			TokenCalculateStrategy: flow.WarmUp, // 当前流量控制器的Token计算策略,flow.WarmUp表示冷启动
			ControlBehavior:        flow.Reject, // 表示流量控制器的控制策略
			Threshold:              1000,        // 表示流量阈值 1s有十个流量进来
			// StatIntervalInMs:       1000,        // 规则对应的流量控制器的独立统计结构的统计周期
			WarmUpPeriodSec: 60, // 表示预热的时长  表示多长时间达到上限
			// WarmUpColdFactor: 3, //这里表示预热的因子，表示预热的速度
		},
	})
```

### ControlBehavior

字段 ControlBehavior 表示表示流量控制器的控制行为，目前 Sentinel 支持两种控制行为：

Reject：表示如果当前统计周期内，统计结构统计的请求数超过了阈值，就直接拒绝。
Throttling：表示匀速排队的统计策略。它的中心思想是，以固定的间隔时间让请求通过。当请求到来的时候，如果当前请求距离上个通过的请求通过的时间间隔不小于预设值，则让当前请求通过；否则，计算当前请求的预期通过时间，如果该请求的预期通过时间小于规则预设的 timeout 时间，则该请求会等待直到预设时间到来通过（排队等待处理）；若预期的通过时间超出最大排队时长，则直接拒接这个请求。

匀速排队方式会严格控制请求通过的间隔时间，也即是让请求以均匀的速度通过，对应的是漏桶算法。该方式的作用如下图所示：

这种方式主要用于处理间隔性突发的流量，例如消息队列。想象一下这样的场景，在某一秒有大量的请求到来，而接下来的几秒则处于空闲状态，我们希望系统能够在接下来的空闲期间逐渐处理这些请求，而不是在第一秒直接拒绝多余的请求。

> 简单来说匀速排队等待就是在一定的时间间隔内允许请求通过，如果 qps=2，那么 sentinel 会计算，每隔 500ms 通过一个请求

以下规则代表每 100ms 最多通过一个请求，多余的请求将会排队等待通过，若排队时队列长度大于 500ms 则直接拒绝：

```go
{
	Resource:          "some-test",
        TokenCalculateStrategy: flow.Direct,
	ControlBehavior:   flow.Throttling, // 流控效果为匀速排队
        Threshold:             10, // 请求的间隔控制在 1000/10=100 ms
	MaxQueueingTimeMs: 500, // 最长排队等待时间
}
```

上面 Threshold 是 10，Sentinel 默认使用 1s 作为控制周期，表示 1 秒内 10 个请求匀速排队，所以排队时间就是 1000ms/10 = 100ms；

特别地，MaxQueueingTimeMs 设为 0 时代表不允许排队，只控制请求时间间隔，多余的请求将会直接拒绝。

### sentinel 的熔断接口
[https://sentinelguard.io/zh-cn/docs/golang/circuit-breaking.html]


### gin集成sentiel

将这一部分拿过来整成一个函数，放在一个单独的文件里
注意把log改一下 改成自己用的
```go
func InitSentinel(){
	err := sentinel.InitDefault()
	if err != nil {
		log.Fatalf("初始化异常:%v\n", err)
	}
	// 这里的配置应该从nacos里读取，或者从别的配置中心读取
	_, err = flow.LoadRules([]*flow.Rule{
		{
			Resource:               "some-test", // 这里可以改成自己服务的名字
			TokenCalculateStrategy: flow.Direct, // 当前流量控制器的Token计算策略
			ControlBehavior:        flow.Reject, // 表示流量控制器的控制策略
			Threshold:              10,          // 表示流量阈值 1s有十个流量进来
			StatIntervalInMs:       1000,        // 规则对应的流量控制器的独立统计结构的统计周期
		},
	})
	if err != nil {
		log.Fatalf("加载规则失败:%v", err)
	}
}
```

然后将这一行代码整到需要限流的操作之前

```go
e, b := sentinel.Entry("改成服务的名字", sentinel.WithTrafficType(base.Inbound))
if b!=nil{
	// 这里代表被限流了
	// 将信息返回给前端，提示限流了
	c.JSON(http.StatusTooManyRequest,gin.H{
		"msg":"请求过于频繁，请稍后重试"
	})
	return
}
// 这个exit在操作执行完成之后
e.Exit()
```