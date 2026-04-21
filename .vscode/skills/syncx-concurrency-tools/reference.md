# syncx 详细示例

## 1. 写锁操作 WithLock

```go
var mu sync.Mutex
var count int

// 无返回值
syncx.WithLock(&mu, func() {
    count++
})

// 带返回值
val, err := syncx.WithLockReturn[int](&mu, func() (int, error) {
    return count, nil
})

// 仅返回值
val := syncx.WithLockReturnValue[int](&mu, func() int {
    return count
})

// 自定义错误类型
val, myErr := syncx.WithLockReturnWithE[int, *MyError](&mu, func() (int, *MyError) {
    return count, nil
})

// 尝试加锁（非阻塞）
err := syncx.WithTryLock(&mu, func() {
    count++
})
if err == syncx.ErrLockNotAcquired {
    // 锁未获取到
}
```

## 2. 并发安全 Map

```go
m := syncx.NewMap[string, int]()

// 存储
m.Store("key", 42)

// 加载
val, ok := m.Load("key")

// LoadOrStore - 不存在则存储
actual, loaded := m.LoadOrStore("key", 100)

// GetOrCompute - 不存在则计算
val := m.GetOrCompute("key", func() int { return 99 })

// CAS
swapped := m.CompareAndSwap("key", 42, 43)

// 过滤
m.Filter(func(k string, v int) bool { return v > 10 })

// 批量操作
keys := m.Keys()
vals := m.Values()
size := m.Size()
isEmpty := m.IsEmpty()
```

## 3. WorkerPool 工作池

```go
pool := syncx.NewWorkerPool(4, 100)
defer pool.Close()

// 提交任务
err := pool.Submit(func() {
    // 执行任务
})

// 非阻塞提交
err := pool.SubmitNonBlocking(func() {
    // 执行任务
})
if err == syncx.ErrQueueFull {
    // 队列已满
}

pool.Wait()
fmt.Println("queue size:", pool.GetQueueSize())
```

## 4. SafeGo 安全协程启动

```go
// 基本用法 - default recover
syncx.SafeGo(func() {
    panic("oops")
}, syncx.Recover)

// 自定义panic处理器
syncx.SafeGo(func() {
    panic("oops")
}, syncx.RecoverWithHandler(func(r interface{}) {
    log.Println("recovered:", r)
}))

// 将panic转为error
var err error
syncx.SafeGo(func() {
    panic("oops")
}, syncx.RecoverToError(&err, syncx.Recover))
```

## 5. StateMachine 状态机

```go
sm := syncx.NewStateMachine("idle",
    syncx.WithAllowAnyTransition[string](),
    syncx.WithTrackHistory[string](10),
)

sm.AllowTransition("idle", "running")
sm.AllowTransition("running", "stopped")

sm.OnTransition(func(from, to string) {
    log.Printf("transition: %s -> %s", from, to)
})

if sm.CanTransitionTo("running") {
    err := sm.TransitionTo("running")
}

current := sm.CurrentState() // "running"
history := sm.GetHistory()
sm.Reset()
```

## 6. 原子类型

```go
b := syncx.NewBool(false)
b.Store(true)
b.Toggle()
v := b.Load()

i := syncx.NewInt64(0)
i.Add(1)
i.Sub(1)
swapped := i.CAS(0, 10)
s := i.String()

av := syncx.NewAtomicValue[string]("hello")
av.Store("world")
old := av.Swap("new")
swapped := av.CompareAndSwap("new", "newer")
```

## 7. Delayer 延迟器

```go
d := syncx.NewDelayer[int]()
d.SetStrategy(syncx.ExponentialDelayStrategy{})
// 每次调用 Next 获取延迟时间
delay := d.Next()
```

## 8. ParallelExecutor 并行执行器

```go
tasks := map[string]func() (string, error){
    "a": func() (string, error) { return "result-a", nil },
    "b": func() (string, error) { return "result-b", nil },
}
// 使用 NewParallelExecutor 按 key 并行执行
pe := syncx.NewParallelExecutor[string, func() (string, error), string](tasks)
```