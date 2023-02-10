package main

import (
	"fmt"
	"sync"
)

// 使用大小为1的channel来模拟trylock

type Lock struct {
	c chan struct{}
}

// NewLock generate a try lock
func NewLock() Lock {
	var l Lock
	l.c = make(chan struct{}, 1)
	l.c <- struct{}{}
	return l
}

// Lock try lock,return lock result
func (l Lock) Lock() bool {
	lockResult := false
	select {
	case <-l.c:
		lockResult = true
	default:
	}
	return lockResult
}

// unlock
func (l Lock) Unlock() {
	l.c <- struct{}{}
}

var counter int

func main() {
	var l = NewLock()
	var wg sync.WaitGroup
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			if !l.Lock() {
				fmt.Println("lock failed")
				return
			}
			counter++
			fmt.Println("current counter", counter)
			l.Unlock()
		}()
	}
	wg.Wait()
}
