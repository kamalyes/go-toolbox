---
name: queue-structures
description: 队列与双端队列数据结构，提供泛型有界队列（自动伸缩、context 取消、超时出队）、FIFO/LIFO 队列、优先队列与环形双端队列（含 range-over-func 迭代器）；当需要带背压的任务队列、优先调度或双端缓冲时使用
---

# queue - 队列与双端队列

五种结构：`BoundedQueue[T]`（泛型，推荐）、`FIFOQueue`、`LIFOQueue`、`PriorityQueue`、`Deque`

## 先选对类型

| 类型 | 泛型 | 阻塞/取消 | 适用 |
|---|---|---|---|
| `BoundedQueue[T]` | ✅ | `context` + 超时 | **新代码首选**：带背压的生产者消费者 |
| `FIFOQueue` | ❌ `interface{}` | `context` | 先进先出，容量可自动伸缩 |
| `LIFOQueue` | ❌ `interface{}` | — | 后进先出（栈） |
| `PriorityQueue` | ❌ `interface{}` | `context` | 按优先级出队 |
| `Deque` | ❌ `interface{}` | — | 两端进出的环形缓冲，**非并发安全** |

只有 `BoundedQueue[T]` 是泛型的，其余四个都是 `interface{}`，每次入队都有一次装箱分配

## 快速开始

```go
import "github.com/kamalyes/go-toolbox/pkg/queue"
```

```go
q := queue.NewBoundedQueue[Task](16, 1024)   // minCap, maxCap
defer q.Close()

// 满了就阻塞，直到有空位或 ctx 取消
if err := q.Enqueue(ctx, task); err != nil {
    return err
}

// 阻塞出队
task, err := q.Dequeue(ctx)

// 非阻塞出队
if task, ok := q.TryDequeue(); ok { /* ... */ }

// 带超时
task, err = q.DequeueTimeout(200 * time.Millisecond)
```

双端队列与 range-over-func 迭代：
```go
d := queue.NewDeque()
d.PushBack(1); d.PushFront(0)

for item := range d.Iter() {            // 只读遍历，不消费
    fmt.Println(item)
}
for item := range d.IterPopFront() {    // 边遍历边弹出，会清空
    fmt.Println(item)
}
```

## 完整API索引

### BoundedQueue[T]

| 导出名称 | 签名 | 说明 |
|---|---|---|
| `NewBoundedQueue[T]` | `func(minCap, maxCap int) *BoundedQueue[T]` | 创建有界队列 |
| `Enqueue` | `func(ctx context.Context, item T) error` | 入队，满时阻塞到有空位或 ctx 结束 |
| `Dequeue` | `func(ctx context.Context) (T, error)` | 出队，空时阻塞 |
| `DequeueTimeout` | `func(timeout time.Duration) (T, error)` | 带超时出队 |
| `TryDequeue` | `func() (T, bool)` | 非阻塞出队 |
| `Size` / `Cap` | `func() int` | 当前长度 / 容量 |
| `IsEmpty` / `IsClosed` | `func() bool` | 是否空 / 是否已关闭 |
| `SetAutoResize` | `func(enabled bool)` | 开关自动伸缩 |
| `Close` | `func()` | 关闭队列 |
| `Stats` | `func() map[string]interface{}` | 运行统计 |

### FIFOQueue

`NewFIFOQueue(capacity int, autoResize bool) *FIFOQueue`

方法：`Enqueue(ctx, item) error`、`Dequeue(ctx) (interface{}, error)`、`IsEmpty`、
`Capacity`、`MinCapacity`、`GrowthFactor`、`SetGrowthFactor(float64) *FIFOQueue`、
`SetShrinkFactor(float64) *FIFOQueue`

### LIFOQueue

方法包含 `IsEmpty`、`Size` 等；语义为后进先出

### PriorityQueue

`NewPriorityQueue() *PriorityQueue`

| 方法 | 签名 | 说明 |
|---|---|---|
| `Enqueue` | `func(ctx context.Context, item interface{}, priority int) error` | 带优先级入队 |
| `Dequeue` | `func(ctx context.Context) (interface{}, error)` | 按优先级出队 |
| `Len` / `Size` | `func() int` | 元素个数 |
| `IsEmpty` | `func() bool` | 是否空 |
| `Push` / `Pop` / `Less` / `Swap` | — | `heap.Interface` 实现，**不要直接调用** |

`Push`/`Pop`/`Less`/`Swap` 是为 `container/heap` 准备的，直接调用会绕过锁与不变量，
请一律用 `Enqueue`/`Dequeue`

### Deque（环形双端队列）

`NewDeque() *Deque`

| 方法 | 签名 | 说明 |
|---|---|---|
| `PushFront` / `PushBack` | `func(elem interface{})` | 两端入队 |
| `PopFront` / `PopBack` | `func() interface{}` | 两端出队 |
| `Front` / `Back` | `func() (interface{}, error)` | 查看两端元素，不弹出 |
| `At` / `Set` | `func(i int) interface{}` / `func(i int, item interface{})` | 按下标读写 |
| `Insert` | `func(i int, item interface{})` | 指定位置插入 |
| `Index` / `RIndex` | `func(f func(interface{}) bool) int` | 正向 / 反向查找，返回下标 |
| `Rotate` | `func(n int)` | 旋转 n 位 |
| `Grow` | `func(n int)` | 预扩容 |
| `Len` / `Cap` / `Clear` | — | 长度 / 容量 / 清空 |
| `Iter` / `RIter` | `func() func(func(interface{}) bool)` | 正向 / 反向遍历，**不消费** |
| `IterPopFront` / `IterPopBack` | `func() func(func(interface{}) bool)` | 边遍历边弹出，**会清空队列** |

### 接口

`Queue` —— 队列的公共接口

## 注意事项

- **`Deque` 不是并发安全的**，多 goroutine 访问要自己加锁；
  `BoundedQueue`/`FIFOQueue`/`PriorityQueue` 内部有锁
- **`IterPopFront` / `IterPopBack` 是消费型迭代器**，走完一遍队列就空了
  想重复遍历请用 `Iter` / `RIter`
- 除 `BoundedQueue[T]` 外都是 `interface{}` 容器，入队有装箱分配；
  元素是小值类型且吞吐敏感时优先 `BoundedQueue[T]`
- `Enqueue` 在队列满时**阻塞**而不是报错，这是背压设计；
  不想阻塞请自己带可取消的 `ctx`，或改用 `TryDequeue` 的对偶逻辑
- `PopFront` / `PopBack` 在空队列上的返回值是 `nil`，与「存了一个 nil 元素」无法区分；
  需要区分请先判 `Len()`，或用返回 error 的 `Front` / `Back`

## 本包有五个基准测试目前是无效的，读数前先看这里

这不是推测，是可复现的：

```
$ go test ./pkg/queue/ -run='^$' -bench='BenchmarkDequeIterPopFront$|BenchmarkFIFOQueueDequeue'
BenchmarkDequeIterPopFront-10    1000000000    0.2309 ns/op
BenchmarkFIFOQueueDequeue-10     --- FAIL: 队列为空
```

**`BenchmarkFIFOQueueDequeue` 会直接让 `go test -bench=.` 失败。**
它预填 1000 个元素后出队 `b.N` 次，而 `b.N` 会远超 1000，
队列一空 `Dequeue` 返回错误、`b.Fatalf` 终止

**`BenchmarkDequeIterPopFront` / `IterPopBack` / `IterPopFrontPartial` / `IterPopBackPartial`
测的是空循环。** 它们在循环外预填 1000 个元素，但 `IterPopFront()` 在第一次迭代就把队列排空，
第 2 到第 10 亿次迭代都在遍历一个空队列
`0.23 ns/op` 差不多就是一个时钟周期——**基准报出接近 0.2 ns 时，基本都意味着循环体什么都没做**

两者的修法相同：把夹具搬进循环，并用 `b.StopTimer()` / `b.StartTimer()` 把重填的开销排除在计时外，
或者按 `b.N` 预填足量元素
在修好之前，**不要引用这五个数字**，也不要拿它们判断 `Deque` 的迭代器性能
