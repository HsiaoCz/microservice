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
