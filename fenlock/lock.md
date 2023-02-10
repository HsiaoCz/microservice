## 分布式锁

为什么需要分布式锁
比如用户下单，我们需要锁住 uid 防止重复下单
比如库存扣减，要锁住库存，防止超卖
比如余额扣减，锁住账户，防止并发操作

在分布式系统中同一个资源往往需要分布式锁来保证变更资源的一致性

分布式锁需要具备哪些特性:
1、排他性：
锁的基本特性，并且只能被第一个持有者持有
2、防死锁
高并发场景下临界资源一旦发生死锁非常难以排查，通常可以设置超时时间到期自动释放锁来规避
3、可重入
持有者支持可重入，防止锁持有者再次重入锁时被释放
4、高性能高可用
锁是代码运行的关键前置节点，一旦不可用则业务直接报故障了，高并发场景下，高性能高可用是基础要求

### 1、单机并发的锁

在单机程序并发修改全局变量的时候，需要加锁以创建临界区

```go
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
```

之前的例子 我们使用 lock 来加锁
如果在某些场景下，我们需要一个任务有单一的执行者，而不是像计数场景那样，所有的 goroutine 都执行成功，我们需要后来的 groutine 在抢锁失败之后，自动放弃流程，这时候就需要 trylock

trylock:尝试加锁，加锁成功后执行后续流程，如果加锁失败的话也不会阻塞，而是直接返回加锁的结果

```go
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
```

因为我们的逻辑限定每个 goroutine 只有成功执行了 Lock 才会继续执行后续逻辑，因此在 Unlock 时可以保证 Lock 结构体中的 channel 一定是空，从而不会阻塞，也不会失败。上面的代码使用了大小为 1 的 channel 来模拟 trylock，理论上还可以使用标准库中的 CAS 来实现相同的功能且成本更低，读者可以自行尝试。

在单机系统中，trylock 并不是一个好选择。因为大量的 goroutine 抢锁可能会导致 CPU 无意义的资源浪费。有一个专有名词用来描述这种抢锁的场景：活锁。

活锁指的是程序看起来在正常执行，但实际上 CPU 周期被浪费在抢锁，而非执行任务上，从而程序整体的执行效率低下。活锁的问题定位起来要麻烦很多。所以在单机场景下，不建议使用这种锁。

### 2、基于 redis 的 setnx 的分布式锁

在分布式场景下也需要这种抢占的逻辑
比如防止重复下单，订单超卖等等
redis 提供了一个 setnx 命令

```go
func incr() {
	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	var lockKey = "counter_lock"
	var counterKey = "counter"

	// lock
	resp := client.SetNX(lockKey, 1, time.Second*5)
	lockSuccess, err := resp.Result()

	if err != nil || !lockSuccess {
		fmt.Println(err, "lock result: ", lockSuccess)
		return
	}

	// counter ++
	getResp := client.Get(counterKey)
	cntValue, err := getResp.Int64()
	if err == nil {
		cntValue++
		resp := client.Set(counterKey, cntValue, 0)
		_, err := resp.Result()
		if err != nil {
			// log err
			println("set value error!")
		}
	}
	println("current counter is ", cntValue)

	delResp := client.Del(lockKey)
	unlockSuccess, err := delResp.Result()
	if err == nil && unlockSuccess > 0 {
		println("unlock success!")
	} else {
		println("unlock failed", err)
	}
}

func main() {
	var wg sync.WaitGroup
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			incr()
		}()
	}
	wg.Wait()
}
```

通过代码和执行结果可以看到，我们远程调用 setnx 实际上和单机的 trylock 非常相似，如果获取锁失败，那么相关的任务逻辑就不应该继续向前执行。

setnx 很适合在高并发场景下，用来争抢一些“唯一”的资源。比如交易撮合系统中卖家发起订单，而多个买家会对其进行并发争抢。这种场景我们没有办法依赖具体的时间来判断先后，因为不管是用户设备的时间，还是分布式场景下的各台机器的时间，都是没有办法在合并后保证正确的时序的。哪怕是我们同一个机房的集群，不同的机器的系统时间可能也会有细微的差别。

所以，我们需要依赖于这些请求到达 Redis 节点的顺序来做正确的抢锁操作。如果用户的网络环境比较差，那也只能自求多福了。

### 3、基于 zookeeper 的分布式锁

```go
func main() {
	c, _, err := zk.Connect([]string{"127.0.0.1"}, time.Second) //*10)
	if err != nil {
		panic(err)
	}
	l := zk.NewLock(c, "/lock", zk.WorldACL(zk.PermAll))
	err = l.Lock()
	if err != nil {
		panic(err)
	}
	println("lock succ, do your business logic")

	time.Sleep(time.Second * 10)

	// do some thing
	l.Unlock()
	println("unlock succ, finish business logic")
}
```

基于 ZooKeeper 的锁与基于 Redis 的锁的不同之处在于 Lock 成功之前会一直阻塞，这与我们单机场景中的 mutex.Lock 很相似。

其原理也是基于临时 Sequence 节点和 watch API，例如我们这里使用的是/lock 节点。Lock 会在该节点下的节点列表中插入自己的值，只要节点下的子节点发生变化，就会通知所有 watch 该节点的程序。这时候程序会检查当前节点下最小的子节点的 id 是否与自己的一致。如果一致，说明加锁成功了。

这种分布式的阻塞锁比较适合分布式任务调度场景，但不适合高频次持锁时间短的抢锁场景。按照 Google 的 Chubby 论文里的阐述，基于强一致协议的锁适用于粗粒度的加锁操作。这里的粗粒度指锁占用时间较长。我们在使用时也应思考在自己的业务场景中使用是否合适。

### 4、基于 ETCD 的分布式锁

这里的 ETCD 包似乎出了点问题

```go
func main() {
	m, err := etcdsync.New("/lock", 10, []string{"http://127.0.0.1:2379"})
	if m == nil || err != nil {
		log.Printf("etcdsync.New failed")
		return
	}
	err = m.Lock()
	if err != nil {
		log.Printf("etcdsync.Lock failed")
		return
	}

	log.Printf("etcdsync.Lock OK")
	log.Printf("Get the lock. Do something here.")

	err = m.Unlock()
	if err != nil {
		log.Printf("etcdsync.Unlock failed")
	} else {
		log.Printf("etcdsync.Unlock OK")
	}
}
```

etcd 中没有像 ZooKeeper 那样的 Sequence 节点。所以其锁实现和基于 ZooKeeper 实现的有所不同。在上述示例代码中使用的 etcdsync 的 Lock 流程是：

先检查/lock 路径下是否有值，如果有值，说明锁已经被别人抢了
如果没有值，那么写入自己的值。写入成功返回，说明加锁成功。写入时如果节点被其它节点写入过了，那么会导致加锁失败，这时候到 3
watch /lock 下的事件，此时陷入阻塞
当/lock 路径下发生事件时，当前进程被唤醒。检查发生的事件是否是删除事件（说明锁被持有者主动 unlock），或者过期事件（说明锁过期失效）。如果是的话，那么回到 1，走抢锁流程。
值得一提的是，在 etcdv3 的 API 中官方已经提供了可以直接使用的锁 API，读者可以查阅 etcd 的文档做进一步的学习。
