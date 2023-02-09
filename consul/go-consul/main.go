package main

import (
	"fmt"
	"log"

	"github.com/hashicorp/consul/api"
)

// 对服务的注册
func Register(address string, port int, name string, tags []string, id string) error {
	cfg := api.DefaultConfig()
	// 这里的address需要填condul注册的地址
	cfg.Address = "172.27.109.169:8500"

	client, err := api.NewClient(cfg)
	if err != nil {
		log.Fatalln(err)
	}
	// 生成对应的检查对象
	// 这里的Http也可以拼接一下
	check := &api.AgentServiceCheck{
		HTTP:                           "http://172.27.109.169:3002/health",
		Timeout:                        "5s",
		Interval:                       "5s",
		DeregisterCriticalServiceAfter: "10s",
	}
	//生成注册对象
	registeration := new(api.AgentServiceRegistration)
	registeration.Name = name
	registeration.ID = id
	registeration.Port = port
	registeration.Tags = tags
	registeration.Address = address
	registeration.Check = check
	err = client.Agent().ServiceRegister(registeration)
	if err != nil {
		log.Fatalln(err)
	}
	return err
}

// 服务的发现
// 获取所有的服务
func AllServices() {
	cfg := api.DefaultConfig()
	cfg.Address = "172.27.109.169:8500"

	client, err := api.NewClient(cfg)
	if err != nil {
		log.Fatalln(err)
	}
	data, err := client.Agent().Services()
	if err != nil {
		log.Fatalln(err)
	}
	for key := range data {
		fmt.Println(key)
	}
}

// 发现某个服务
// 这里其实也可以传递一个参数
func FilterService() {
	cfg := api.DefaultConfig()
	cfg.Address = "172.27.109.169:8500"

	client, err := api.NewClient(cfg)
	if err != nil {
		log.Fatalln(err)
	}
	data, err := client.Agent().ServicesWithFilter(`Service=="Hello"`)
	if err != nil {
		log.Fatalln(err)
	}
	for key := range data {
		fmt.Println(key)
	}
}

func main() {
	err := Register("172.27.109.169", 3002, "Hello", []string{"hello", "Hsiao", "web", "gin"}, "12138")
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println("服务注册成功")
	AllServices()
	FilterService()
}
