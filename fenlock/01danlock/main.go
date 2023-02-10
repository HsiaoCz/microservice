package main

import (
	"fmt"
	"sync"
)

// 我们看一下并发访问资源不加锁会发生什么

// 全局变量
var counter int
var wg sync.WaitGroup
var lock sync.Mutex

func main() {
	for i := 0; i < 1000; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			lock.Lock()
			counter++
			lock.Unlock()
		}()
	}
	wg.Wait()
	fmt.Println(counter)
}

// 并发访问资源 如果不加锁，每次都会得到不一样的结果
// 想要的到正确的结果就需要加锁
