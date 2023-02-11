package main

import (
	"fmt"
	"log"
	"math/rand"
	"time"

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
	if err != nil {
		log.Fatalf("加载规则失败:%v", err)
	}
	ch := make(chan struct{})
	// 使用 sentinel.Entry表示一个流控的入口点
	// 第一个参数是资源名 表示入口点使用资源名的规则做检查
	// sentinel.WithTrafficType 配置入口或出口的流量控制
	// base.Inbound表示入口的流量控制
	// 这里模拟一下流量

	// 这里启动多个线程来模拟一下
	// 计算一下通过的流量，总共的流量，没通过的流量
	var sumTotal int   // 总共的流量
	var passTotal int  // 通过的流量
	var blockTotal int // 没通过的流量
	for i := 0; i < 100; i++ {
		go func() {
			for {
				sumTotal++
				e, b := sentinel.Entry("some-test", sentinel.WithTrafficType(base.Inbound))
				if b != nil {
					// fmt.Println("限流了")
					blockTotal++
					time.Sleep(time.Duration(rand.Uint64()%10) * time.Millisecond)
				} else {
					passTotal++
					time.Sleep(time.Duration(rand.Uint64()%10) * time.Millisecond)
					e.Exit()
				}
			}
		}()
	}

	go func() {
		// 这里再统计一下，过去一秒钟总共产生了多少个
		// 总共通过了多少个
		// 总共没通过多少个
		var oldTotal int //过去一秒钟总共产生了多少个
		var oldPass int  // 总共通过了多少个
		var oldBlock int // 总共没通过多少个

		for {
			oneSecond := sumTotal - oldTotal
			oldTotal = sumTotal
			oneSecondPass := passTotal - oldPass
			oldPass = passTotal

			oneSecondBlock := blockTotal - oldBlock
			oldBlock = blockTotal
			time.Sleep(time.Second)
			fmt.Println("过去一秒钟的总数:", oneSecond, "过去一秒通过了多少:", oneSecondPass, "过去一秒没通过的数目:", oneSecondBlock)
		}
	}()

	<-ch
}
