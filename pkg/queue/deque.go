/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2025-08-21 16:01:08
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2025-08-22 09:59:16
 * @FilePath: \go-toolbox\pkg\queue\deque.go
 * @Description:
 * Deque 是一个双端队列（double-ended queue）实现，支持从两端插入和删除元素。
 * 该实现使用环形缓冲区来优化存储和访问效率。队列支持动态扩容和缩容，
 * 以适应不同的使用场景，提供了多种操作方法，包括 Push、Pop、Iterate 等
 * 适合需要高效插入和删除的场景，如任务调度、缓存等
 * Copyright (c) 2025 by kamalyes, All Rights Reserved.
 */
package queue

import "errors"

// minCapacity 是双端队列可能拥有的最小容量必须是 2 的幂
// 以便使用位运算：x % n == x & (n - 1)
const minCapacity = 16

// Deque 表示双端队列数据结构的单个实例Deque 实例包含
// 指定类型的项
type Deque struct {
	buf   []interface{} // 存储队列元素的缓冲区，使用切片实现
	head  int           // 队列头部索引，指向队列的第一个元素
	tail  int           // 队列尾部索引，指向下一个插入位置
	count int           // 当前队列中元素的数量
}

// NewDeque 创建并返回一个新的 Deque 实例
// 该函数初始化一个双端队列，分配最小容量的缓冲区，并设置头、尾索引和元素计数
func NewDeque() *Deque {
	return &Deque{
		buf:   make([]interface{}, minCapacity), // 初始化缓冲区，分配最小容量
		head:  0,                                // 初始化头部索引为 0
		tail:  0,                                // 初始化尾部索引为 0
		count: 0,                                // 初始化元素计数为 0
	}
}

// Cap 返回 Deque 的当前容量如果 q 为 nil，q.Cap() 返回零
func (q *Deque) Cap() int {
	if q == nil {
		return 0
	}
	return len(q.buf)
}

// Len 返回当前存储在队列中的元素数量如果 q 为 nil，
// q.Len() 返回零
func (q *Deque) Len() int {
	if q == nil {
		return 0
	}
	return q.count
}

// PushBack 将元素追加到队列的末尾当使用 PopFront 删除元素时实现 FIFO，
// 当使用 PopBack 删除元素时实现 LIFO
func (q *Deque) PushBack(elem interface{}) {
	q.growIfFull()

	q.buf[q.tail] = elem
	// 计算新的尾部位置
	q.tail = q.next(q.tail)
	q.count++
}

// PushFront 在队列的前面插入元素
func (q *Deque) PushFront(elem interface{}) {
	q.growIfFull()

	// 计算新的头部位置
	q.head = q.prev(q.head)
	q.buf[q.head] = elem
	q.count++
}

// PopFront 从队列的前面移除并返回元素
// 当与 PushBack 一起使用时实现 FIFO如果队列为空，则调用会 panic
func (q *Deque) PopFront() interface{} {
	if q.count <= 0 {
		panic("deque: PopFront() 在空队列上调用")
	}
	ret := q.buf[q.head]
	var zero interface{}
	q.buf[q.head] = zero
	// 计算新的头部位置
	q.head = q.next(q.head)
	q.count--

	q.shrinkIfExcess()
	return ret
}

// IterPopFront 返回一个迭代器，该迭代器从双端队列的前面迭代移除项目
// 这比一次移除一个项目更有效，因为它避免了中间的调整大小
// 如果需要调整大小，则仅在迭代结束时进行一次
func (q *Deque) IterPopFront() func(func(interface{}) bool) {
	return func(yield func(interface{}) bool) {
		if q.Len() == 0 {
			return
		}
		var zero interface{}
		for q.count != 0 {
			ret := q.buf[q.head]
			q.buf[q.head] = zero
			q.head = q.next(q.head)
			q.count--
			if !yield(ret) {
				break
			}
		}
		q.shrinkToFit()
	}
}

// PopBack 从队列的末尾移除并返回元素
// 当与 PushBack 一起使用时实现 LIFO如果队列为空，则调用会 panic
func (q *Deque) PopBack() interface{} {
	if q.count <= 0 {
		panic("deque: PopBack() 在空队列上调用")
	}

	// 计算新的尾部位置
	q.tail = q.prev(q.tail)

	// 移除尾部的值
	ret := q.buf[q.tail]
	var zero interface{}
	q.buf[q.tail] = zero
	q.count--

	q.shrinkIfExcess()
	return ret
}

// IterPopBack 返回一个迭代器，该迭代器从双端队列的末尾迭代移除项目
// 这比一次移除一个项目更有效，因为它避免了中间的调整大小
// 如果需要调整大小，则仅在迭代结束时进行一次
func (q *Deque) IterPopBack() func(func(interface{}) bool) {
	return func(yield func(interface{}) bool) {
		if q.Len() == 0 {
			return
		}
		var zero interface{}
		for q.count != 0 {
			q.tail = q.prev(q.tail)
			ret := q.buf[q.tail]
			q.buf[q.tail] = zero
			q.count--
			if !yield(ret) {
				break
			}
		}
		q.shrinkToFit()
	}
}

// Front 返回队列前面的元素这是 PopFront 返回的元素
// 如果队列为空，则调用会返回错误信息
func (q *Deque) Front() (interface{}, error) {
	if q.count <= 0 {
		return nil, errors.New("deque: Front() 在空队列上调用")
	}
	return q.buf[q.head], nil
}

// Back 返回队列末尾的元素这是 PopBack 返回的元素
// 如果队列为空，则调用会返回错误信息
func (q *Deque) Back() (interface{}, error) {
	if q.count <= 0 {
		return nil, errors.New("deque: PopBack() 在空队列上调用")
	}
	return q.buf[q.prev(q.tail)], nil
}

// At 返回队列中索引 i 处的元素，而不移除该元素
// 此方法仅接受非负索引值At(0) 指的是第一个元素，
// 与 Front() 相同At(Len()-1) 指的是最后一个元素，
// 与 Back() 相同如果索引无效，调用会 panic
func (q *Deque) At(i int) interface{} {
	q.checkRange(i)
	return q.buf[(q.head+i)&(len(q.buf)-1)]
}

// Set 将项目分配给队列中索引 i 的位置
// Set 的索引与 At 相同，但执行相反的操作
// 如果索引无效，调用会 panic
func (q *Deque) Set(i int, item interface{}) {
	q.checkRange(i)
	q.buf[(q.head+i)&(len(q.buf)-1)] = item
}

// Iter 返回一个迭代器，用于遍历 Deque 中的所有项目，
// 从前（索引 0）到后（索引 Len()-1）依次返回每个项目
// 在迭代过程中修改 Deque 会导致 panic
func (q *Deque) Iter() func(func(interface{}) bool) {
	return func(yield func(interface{}) bool) {
		origHead := q.head
		origTail := q.tail
		head := origHead
		for i := 0; i < q.Len(); i++ {
			if q.head != origHead || q.tail != origTail {
				panic("deque: 在迭代过程中修改了队列")
			}
			if !yield(q.buf[head]) {
				return
			}
			head = q.next(head)
		}
	}
}

// RIter 返回一个反向迭代器，用于遍历 Deque 中的所有项目，
// 从后（索引 Len()-1）到前（索引 0）依次返回每个项目
// 在迭代过程中修改 Deque 会导致 panic
func (q *Deque) RIter() func(func(interface{}) bool) {
	return func(yield func(interface{}) bool) {
		origHead := q.head
		origTail := q.tail
		tail := origTail
		for i := 0; i < q.Len(); i++ {
			if q.head != origHead || q.tail != origTail {
				panic("deque: 在迭代过程中修改了队列")
			}
			tail = q.prev(tail)
			if !yield(q.buf[tail]) {
				return
			}
		}
	}
}

// Clear 移除队列中的所有元素，但保留当前容量
// 这在高频率重复使用队列时非常有用，以避免垃圾回收
// 只要仅添加项目，队列就不会被调整为更小的尺寸
// 只有在移除项目时，队列才会被调整为更小的尺寸
func (q *Deque) Clear() {
	if q.Len() == 0 {
		return
	}
	q.count = 0
	q.head = 0
	q.tail = 0
	for i := range q.buf {
		q.buf[i] = nil // 清空
	}
}

// Grow 如果需要，增加双端队列的容量，以保证可以容纳另 n
// 个项目在 Grow(n) 之后，至少可以向队列中写入 n 个项目，
// 而无需再次分配如果 n 为负数，Grow 会 panic
func (q *Deque) Grow(n int) {
	if n < 0 {
		panic("deque.Grow: 负数计数")
	}
	c := q.Cap()
	l := q.Len()
	// 如果已经足够大
	if n <= c-l {
		return
	}

	if c == 0 {
		c = minCapacity
	}

	newLen := l + n
	for c < newLen {
		c <<= 1
	}
	if l == 0 {
		q.buf = make([]interface{}, c)
		q.head = 0
		q.tail = 0
	} else {
		q.resize(c)
	}
}

// Rotate 将双端队列旋转 n 步，从前到后如果 n 为负，则从后到前旋转
// 让 Deque 提供 Rotate 可以避免使用仅 Pop 和 Push 方法实现旋转时可能发生的调整大小
// 如果 q.Len() 为 1 或更少，或 q 为 nil，则 Rotate 不执行任何操作
func (q *Deque) Rotate(n int) {
	if q.Len() <= 1 {
		return
	}
	// 旋转 q.count 的倍数等同于不旋转
	n %= q.count
	if n == 0 {
		return
	}

	modBits := len(q.buf) - 1
	// 如果缓冲区没有空闲空间，仅移动头部和尾部索引
	if q.head == q.tail {
		// 使用位运算计算新的头部和尾部
		q.head = (q.head + n) & modBits
		q.tail = q.head
		return
	}

	var zero interface{}

	if n < 0 {
		// 从后到前旋转
		for ; n < 0; n++ {
			// 使用位运算计算新的头部和尾部
			q.head = (q.head - 1) & modBits
			q.tail = (q.tail - 1) & modBits
			// 将尾部值放到头部，并移除尾部的值
			q.buf[q.head] = q.buf[q.tail]
			q.buf[q.tail] = zero
		}
		return
	}

	// 从前到后旋转
	for ; n > 0; n-- {
		// 将头部值放到尾部，并移除头部的值
		q.buf[q.tail] = q.buf[q.head]
		q.buf[q.head] = zero
		// 使用位运算计算新的头部和尾部
		q.head = (q.head + 1) & modBits
		q.tail = (q.tail + 1) & modBits
	}
}

// Index 返回满足 f(item) 的第一个项在 Deque 中的索引，
// 如果没有满足条件的项，则返回 -1如果 q 为 nil，则总是返回 -1
// 搜索是线性的，从索引 0 开始
func (q *Deque) Index(f func(interface{}) bool) int {
	if q.Len() > 0 {
		modBits := len(q.buf) - 1
		for i := 0; i < q.count; i++ {
			if f(q.buf[(q.head+i)&modBits]) {
				return i
			}
		}
	}
	return -1
}

// RIndex 与 Index 相同，但从后向前搜索返回的索引
// 从前向后，其中索引 0 是 Front() 返回的项目的索引
func (q *Deque) RIndex(f func(interface{}) bool) int {
	if q.Len() > 0 {
		modBits := len(q.buf) - 1
		for i := q.count - 1; i >= 0; i-- {
			if f(q.buf[(q.head+i)&modBits]) {
				return i
			}
		}
	}
	return -1
}

// Insert 用于将元素插入队列中的指定索引位置如果索引无效，
// 调用会 panic
func (q *Deque) Insert(i int, item interface{}) {
	if i < 0 || i > q.count {
		panic("deque: 索引超出范围")
	}
	if i == 0 {
		q.PushFront(item)
		return
	}
	if i == q.count {
		q.PushBack(item)
		return
	}

	q.growIfFull()

	// 计算插入位置的尾部索引
	insertPos := (q.head + i) & (len(q.buf) - 1)
	// 移动元素以为新元素腾出空间
	for j := q.count; j > i; j-- {
		q.buf[(q.head+j)&(len(q.buf)-1)] = q.buf[(q.head+j-1)&(len(q.buf)-1)]
	}
	q.buf[insertPos] = item
	q.count++
	q.tail = (q.tail + 1) & (len(q.buf) - 1)
}

// checkRange 检查索引是否在有效范围内
func (q *Deque) checkRange(i int) {
	if i < 0 || i >= q.count {
		panic("deque: 索引超出范围")
	}
}

// next 返回下一个索引
func (q *Deque) next(i int) int {
	return (i + 1) & (len(q.buf) - 1)
}

// prev 返回前一个索引
func (q *Deque) prev(i int) int {
	return (i - 1) & (len(q.buf) - 1)
}

// growIfFull 检查队列是否已满，如果已满则增长容量
func (q *Deque) growIfFull() {
	if q.count == len(q.buf) {
		q.Grow(1)
	}
}

// shrinkIfExcess 如果当前元素数量远小于容量，缩小容量
func (q *Deque) shrinkIfExcess() {
	if q.count < len(q.buf)/4 && len(q.buf) > minCapacity {
		q.resize(len(q.buf) / 2)
	}
}

// shrinkToFit 将队列的容量调整为当前元素数量
func (q *Deque) shrinkToFit() {
	if q.count < len(q.buf) {
		q.resize(q.count)
	}
}

// resize 重新分配队列的缓冲区以适应新的容量
func (q *Deque) resize(newCap int) {
	newBuf := make([]interface{}, newCap)
	if q.count > 0 {
		for i := 0; i < q.count; i++ {
			newBuf[i] = q.buf[(q.head+i)&(len(q.buf)-1)]
		}
	}
	q.buf = newBuf
	q.head = 0
	q.tail = q.count
}
