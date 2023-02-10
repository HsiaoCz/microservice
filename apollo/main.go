package main

import (
	"fmt"
	"log"

	"github.com/philchia/agollo/v4"
)

// go 接入apollo
func main() {
	agollo.Start(&agollo.Conf{
		AppID:           "SampleApp",
		Cluster:         "dev",
		NameSpaceNames:  []string{"application.properties", "shopping_cart.yaml"},
		MetaAddr:        "http://localhost:8080",
		AccesskeySecret: "b8ceb3ec62f34030b1b1fd9a431e420b",
	})

	agollo.OnUpdate(func(ce *agollo.ChangeEvent) {
		//监听配置变更
		log.Printf("ce:%#v\n", ce)
	})
	log.Println("初始化配置成功")

	//从默认的application.properties命名空间获取key的值
	val := agollo.GetString("timeout")
	log.Println(val)
	// 获取命名空间下所有的key
	keys := agollo.GetAllKeys(agollo.WithNamespace("shopping_cart.yaml"))
	fmt.Println(keys)

	// 获取指定一个命名空间下key的值
	other := agollo.GetString("content", agollo.WithNamespace("shopping_cart.yaml"))
	log.Println(other)

	//获取命名空间下的所有内容
	namespaceContent := agollo.GetContent(agollo.WithNamespace("shopping_cart.yaml"))
	log.Println(namespaceContent)
}
