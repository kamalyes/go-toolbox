/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2025-06-05 16:25:18
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2025-07-23 15:10:33
 * @FilePath: \go-toolbox\pkg\syncx\task.go
 * @Description:
 * 泛型参数说明：
 *
 * T - 任务输入类型
 *    - 定义：`T` 是一个类型参数，表示任务执行函数的输入类型
 *    - 用法：在任务函数中，`input` 参数的类型为 `T`, 这意味着你可以在创建任务时指定任何类型作为输入
 *      例如，如果你创建一个任务来处理整数输入，则 `T` 可以是 `int`；如果任务需要处理字符串，则 `T` 可以是 `string`
 *
 * R - 任务结果类型
 *    - 定义：`R` 是一个类型参数，表示任务执行函数的返回结果类型
 *    - 用法：在任务函数中，返回值的类型为 `R`, 这使得你可以灵活地指定任务完成后返回的结果类型
 *      例如，如果任务执行后需要返回一个字符串，则 `R` 可以是 `string`；如果返回一个整数，则 `R` 可以是 `int`
 *
 * U - 回调结果类型
 *    - 定义：`U` 是一个类型参数，表示任务成功或失败后的回调函数的返回结果类型
 *    - 用法：在设置成功或失败回调时，回调函数的返回值类型为 `U`, 这允许你指定回调函数的返回类型
 *      例如，如果回调函数需要返回一个字符串，则 `U` 可以是 `string`；如果返回一个布尔值，则 `U` 可以是 `bool`
 *
 * 总结：
 * - `T`：表示任务输入的类型，允许灵活地定义任务需要处理的数据类型
 * - `R`：表示任务执行的结果类型，使得任务的返回值可以是任何类型
 * - `U`：表示回调函数的返回结果类型，允许在任务执行后处理结果并返回相应的数据类型
 *
 * 通过使用这些泛型参数，`Task` 和 `TaskManager` 可以在不同的上下文中使用，适应不同类型的任务和回调，增强了代码的灵活性和可重用性
 * Copyright (c) 2025 by kamalyes, All Rights Reserved.
 */
package syncx

import (
	"context"
	"fmt"
	"reflect"
	"sort"
	"sync"
	"time"
)

// TaskState 表示任务的状态
type TaskState int

const (
	Pending   TaskState = 1 << iota // 等待中
	Running                         // 运行中
	Completed                       // 已完成
	Cancelled                       // 已取消
	Failed                          // 失败

)

// 能被中途取消的状态
var cancellableStates = map[TaskState]struct{}{
	Pending: {},
	Running: {},
	Failed:  {},
}

// ExecutionMode 表示任务执行模式
type ExecutionMode int

const (
	Sequential ExecutionMode = iota // 顺序执行
	Concurrent                      // 并发执行
)

// TaskType 表示任务类型
type TaskType int

const (
	MainTask       TaskType = iota // 主任务
	DependencyTask                 // 依赖任务
)

// Task 表示一个异步任务
type Task[T any, R any, U any] struct {
	name             string                                        // 任务名称
	funcPointer      uintptr                                       // 任务函数的指针，用于循环依赖检查
	fn               func(ctx context.Context, input T) (R, error) // 任务执行的函数
	depends          []*Task[T, R, U]                              // 依赖的任务列表
	priority         int                                           // 任务优先级
	state            TaskState                                     // 任务状态
	result           R                                             // 任务执行结果
	err              error                                         // 任务执行错误
	cancel           context.CancelFunc                            // 取消函数
	timeout          time.Duration                                 // 超时时间
	successCallback  func(result R, err error) (U, error)          // 任务成功后的回调函数
	failureCallback  func(result R, err error) (U, error)          // 任务失败后的回调函数
	input            T                                             // 任务输入
	retryCount       int                                           // 当前重试次数
	retryInterval    time.Duration                                 // 重试间隔时间
	maxRetries       int                                           // 最大重试次数
	mu               sync.RWMutex                                  // 互斥锁，确保线程安全
	ctx              context.Context                               // 传入的上下文
	fnDuration       time.Duration                                 // 主任务运行时间
	callbackDuration time.Duration                                 // 回调运行时间
	callbackResult   U                                             // 存储回调结果
	callbackError    error                                         // 存储回调错误
	callbackState    TaskState                                     // 任务状态
	taskType         TaskType                                      // 任务类型（主任务或依赖任务）
}

// TaskManager 管理所有的任务
type TaskManager[T any, R any, U any] struct {
	tasks               map[string]*Task[T, R, U] // 存储所有任务的映射
	mu                  sync.RWMutex              // 互斥锁，确保并发安全
	wg                  sync.WaitGroup            // 等待组，用于等待所有任务完成
	concurrency         int                       // 最大并发数
	history             map[string][]TaskHistory  // 任务执行历史
	maxHistorySize      int                       // 最大历史行数
	trunUpFunc          func() (R, error)         // 启动时执行的函数
	trunDownFunc        func() (R, error)         // 关闭时执行的函数
	dependExecutionMode ExecutionMode             // 依赖任务的执行模式
	taskChan            chan *Task[T, R, U]       // 任务通道
	isReleased          bool                      // 标识任务是否已释放资源
}

// TaskHistory 记录任务执行的历史信息
type TaskHistory struct {
	Timestamp             time.Time     // 任务执行的时间戳
	State                 TaskState     // 任务的状态（如成功、失败等）
	Result                interface{}   // 任务执行的结果
	Error                 error         // 任务执行过程中发生的错误
	FnExecutionTime       time.Duration // 任务函数的执行时间
	CallbackExecutionTime time.Duration // 回调函数的执行时间
	TaskType              TaskType      // 任务类型（主任务或依赖任务）
}

// NewTaskManager 创建一个新的 TaskManager
func NewTaskManager[T any, R any, U any](concurrency int) *TaskManager[T, R, U] {
	tm := &TaskManager[T, R, U]{
		tasks:          make(map[string]*Task[T, R, U]), // 初始化任务
		history:        make(map[string][]TaskHistory),  // 初始化历史记录
		maxHistorySize: -1,
		concurrency:    concurrency,
		wg:             sync.WaitGroup{},                       // 显式初始化 WaitGroup
		taskChan:       make(chan *Task[T, R, U], concurrency), // 初始化任务通道
		isReleased:     true,
	}
	return tm
}

// NewTaskWithOptions 创建一个新的任务
func NewTaskWithOptions[T any, R any, U any](name string, fn func(ctx context.Context, input T) (R, error), input T, ctx context.Context, maxRetries int, retryInterval time.Duration) *Task[T, R, U] {
	return &Task[T, R, U]{
		name:          name,                          // 任务名称
		fn:            fn,                            // 任务执行的函数
		ctx:           ctx,                           // 存储传入的上下文
		input:         input,                         // 任务的输入数据
		maxRetries:    maxRetries,                    // 最大重试次数
		retryInterval: retryInterval,                 // 使用传入的重试间隔时间
		state:         Pending,                       // 任务状态默认为等待中
		callbackState: Pending,                       // 回调任务状态默认为等待中
		funcPointer:   reflect.ValueOf(fn).Pointer(), // 获取函数指针
		taskType:      MainTask,                      // 任务类型，默认为主任务
	}
}

// NewTask 创建一个新的任务，使用背景上下文
func NewTask[T any, R any, U any](name string, fn func(ctx context.Context, input T) (R, error), input T) *Task[T, R, U] {
	// 调用 NewTaskWithOptions，并传入背景上下文和默认的重试参数
	return NewTaskWithOptions[T, R, U](name, fn, input, context.Background(), 3, 1*time.Second)
}

// AddTask 添加一个任务到 TaskManager
func (tm *TaskManager[T, R, U]) AddTask(task *Task[T, R, U]) {
	WithLock(&tm.mu, func() {
		tm.tasks[task.name] = task
	})
}

// SetMaxHistorySize 设置最大历史记录行到 TaskManager
func (tm *TaskManager[T, R, U]) SetMaxHistorySize(maxHistorySize int) {
	WithLock(&tm.mu, func() {
		tm.maxHistorySize = maxHistorySize
	})
}

// AddDependency 添加依赖关系
func (tk *Task[T, R, U]) AddDependency(dep *Task[T, R, U]) *Task[T, R, U] {
	visited := make(map[uintptr]bool)
	// 检查新依赖是否会导致循环依赖
	if err := dep.checkCircularDependency(visited, tk); err != nil {
		panic(err)
	}

	// 如果没有循环依赖，添加依赖任务
	return WithLockReturnValue(&tk.mu, func() *Task[T, R, U] {
		dep.taskType = DependencyTask
		tk.depends = append(tk.depends, dep)
		return tk
	})
}

// checkCircularDependency 检查任务依赖是否存在循环
func (tk *Task[T, R, U]) checkCircularDependency(visited map[uintptr]bool, newDep *Task[T, R, U]) error {
	// 如果当前任务已经被访问过，说明存在循环依赖
	if visited[tk.funcPointer] {
		return fmt.Errorf("circular dependency detected for task: %s", tk.name)
	}

	// 将当前任务标记为已访问
	visited[tk.funcPointer] = true

	// 检查新依赖是否与当前任务形成循环
	if newDep != nil && newDep.funcPointer == tk.funcPointer {
		return fmt.Errorf("circular dependency detected: task %s cannot depend on itself", newDep.name)
	}

	// 递归检查所有依赖的任务
	for _, dep := range tk.depends {
		if err := dep.checkCircularDependency(visited, newDep); err != nil {
			return err // 如果发现循环依赖，返回错误
		}
	}

	// 从访问记录中删除当前任务，表示这个任务的检查已经完成
	delete(visited, tk.funcPointer)

	return nil // 没有发现循环依赖，返回 nil
}

// SetPriority 设置任务优先级，并返回当前任务以支持链式调用
func (tk *Task[T, R, U]) SetPriority(priority int) *Task[T, R, U] {
	return WithLockReturnValue(&tk.mu, func() *Task[T, R, U] {
		tk.priority = priority
		return tk
	})
}

// SetTimeout 设置任务超时时间，并返回当前任务以支持链式调用
func (tk *Task[T, R, U]) SetTimeout(timeout time.Duration) *Task[T, R, U] {
	return WithLockReturnValue(&tk.mu, func() *Task[T, R, U] {
		tk.timeout = timeout
		return tk
	})
}

// SetSuccessCallback 设置任务成功后的回调函数，并返回当前任务以支持链式调用
func (tk *Task[T, R, U]) SetSuccessCallback(callback func(R, error) (U, error)) *Task[T, R, U] {
	return WithLockReturnValue(&tk.mu, func() *Task[T, R, U] {
		tk.successCallback = callback
		return tk
	})
}

// SetFailureCallback 设置任务失败后的回调函数，并返回当前任务以支持链式调用
func (tk *Task[T, R, U]) SetFailureCallback(callback func(R, error) (U, error)) *Task[T, R, U] {
	return WithLockReturnValue(&tk.mu, func() *Task[T, R, U] {
		tk.failureCallback = callback
		return tk
	})
}

// SetRetryInterval 设置任务失败后重试间隔时间，并返回当前任务以支持链式调用
func (tk *Task[T, R, U]) SetRetryInterval(retryInterval time.Duration) *Task[T, R, U] {
	return WithLockReturnValue(&tk.mu, func() *Task[T, R, U] {
		tk.retryInterval = retryInterval
		return tk
	})
}

// SetMaxRetries 设置最大重试次数，并返回当前任务以支持链式调用
func (tk *Task[T, R, U]) SetMaxRetries(count int) *Task[T, R, U] {
	return WithLockReturnValue(&tk.mu, func() *Task[T, R, U] {
		tk.maxRetries = count
		return tk
	})
}

// GetName 获取任务名称
func (tk *Task[T, R, U]) GetName() string {
	return WithRLockReturnValue(&tk.mu, func() string {
		return tk.name
	})
}

// GetState 获取任务状态
func (tk *Task[T, R, U]) GetState() TaskState {
	return WithRLockReturnValue(&tk.mu, func() TaskState {
		return tk.state
	})
}

// GetCallbackState 获取回调状态
func (tk *Task[T, R, U]) GetCallbackState() TaskState {
	return WithRLockReturnValue(&tk.mu, func() TaskState {
		return tk.callbackState
	})
}

// GetResult 获取任务结果
func (tk *Task[T, R, U]) GetResult() R {
	return WithRLockReturnValue(&tk.mu, func() R {
		return tk.result
	})
}

// GetError 获取任务执行错误
func (tk *Task[T, R, U]) GetError() error {
	return WithRLockReturnValue(&tk.mu, func() error {
		return tk.err
	})
}

// GetRetryCount 获取任务失败后重试次数
func (tk *Task[T, R, U]) GetRetryCount() int {
	return WithRLockReturnValue(&tk.mu, func() int {
		return tk.retryCount
	})
}

// GetMaxRetries 获取最大重试次数
func (tk *Task[T, R, U]) GetMaxRetries() int {
	return WithRLockReturnValue(&tk.mu, func() int {
		return tk.maxRetries
	})
}

// GetFnDuration 获取主任务运行时间
func (tk *Task[T, R, U]) GetFnDuration() time.Duration {
	return WithRLockReturnValue(&tk.mu, func() time.Duration {
		return tk.fnDuration
	})
}

// GetCallbackDuration 获取回调运行时间
func (tk *Task[T, R, U]) GetCallbackDuration() time.Duration {
	return WithRLockReturnValue(&tk.mu, func() time.Duration {
		return tk.callbackDuration
	})
}

// GetCallbackResult 获取回调结果
func (tk *Task[T, R, U]) GetCallbackResult() U {
	return WithRLockReturnValue(&tk.mu, func() U {
		return tk.callbackResult
	})
}

// GetCallbackError 获取回调错误
func (tk *Task[T, R, U]) GetCallbackError() error {
	return WithRLockReturnValue(&tk.mu, func() error {
		return tk.callbackError
	})
}

// SetTrunUp 设置启动时执行的函数
func (tm *TaskManager[T, R, U]) SetTrunUp(fn func() (R, error)) {
	WithLock(&tm.mu, func() {
		tm.trunUpFunc = fn
	})
}

// SetTrunDown 设置关闭时执行的函数
func (tm *TaskManager[T, R, U]) SetTrunDown(fn func() (R, error)) {
	WithLock(&tm.mu, func() {
		tm.trunDownFunc = fn
	})
}

// SetDependExecutionMode 设置依赖任务的执行模式，并返回当前任务以支持链式调用
func (tm *TaskManager[T, R, U]) SetDependExecutionMode(executionMode ExecutionMode) {
	WithLock(&tm.mu, func() {
		tm.dependExecutionMode = executionMode
	})
}

// GetDependExecutionMode 获取依赖任务的执行模式
func (tm *TaskManager[T, R, U]) GetDependExecutionMode() ExecutionMode {
	return WithLockReturnValue(&tm.mu, func() ExecutionMode {
		return tm.dependExecutionMode
	})
}

// TrunUp 启动任务管理器并返回结果和错误
func (tm *TaskManager[T, R, U]) TrunUp() (result R, err error) {
	return WithLockReturn(&tm.mu, func() (R, error) {
		fmt.Println("Task Manager is starting up...")
		if tm.trunUpFunc != nil {
			return tm.trunUpFunc() // 调用设置的启动函数
		}
		return result, err // 返回零值和 nil
	})
}

// TrunDown 关闭任务管理器并返回结果和错误
func (tm *TaskManager[T, R, U]) TrunDown() (result R, err error) {
	return WithLockReturn(&tm.mu, func() (R, error) {
		fmt.Println("Task Manager is shutting down...")
		if tm.trunDownFunc != nil {
			return tm.trunDownFunc() // 调用设置的关闭函数并返回结果和错误
		}
		return result, err // 返回零值和 nil
	})
}

// GetDepends 获取当前任务的所有依赖任务
func (tk *Task[T, R, U]) GetDepends() []*Task[T, R, U] {
	return WithRLockReturnValue(&tk.mu, func() []*Task[T, R, U] {
		return tk.depends
	})
}

// GetDependencyStates 获取所有依赖任务的状态
func (tk *Task[T, R, U]) GetDependencyStates() map[string]TaskState {
	dependencyStates := make(map[string]TaskState)
	WithRLock(&tk.mu, func() {
		for _, dep := range tk.depends {
			dependencyStates[dep.name] = dep.state
		}
	})
	return dependencyStates
}

// worker 处理任务的工作函数
func (tm *TaskManager[T, R, U]) worker() {
	for task := range tm.taskChan {
		tm.runTask(task)
		tm.wg.Done() // 完成任务后减少等待组计数
	}
}

// Run 执行所有任务
func (tm *TaskManager[T, R, U]) Run() {
	WithLock(&tm.mu, func() {
		if !tm.isReleased {
			panic("must call release() to release resources first")
		}
	})
	// 创建并启动 Worker 池
	for i := 0; i < tm.concurrency; i++ {
		go tm.worker()
	}

	// 创建并按优先级排序任务
	taskSlice := make([]*Task[T, R, U], 0, len(tm.tasks))
	for _, task := range tm.tasks {
		taskSlice = append(taskSlice, task)
	}

	// 按优先级排序任务
	sort.Slice(taskSlice, func(i, j int) bool {
		return taskSlice[i].priority > taskSlice[j].priority
	})

	// 将任务发送到队列
	for _, task := range taskSlice {
		tm.wg.Add(1) // 增加等待组计数
		tm.taskChan <- task
	}

	// 等待所有任务完成
	tm.wg.Wait()

	// 关闭任务通道，表示不再有新任务
	close(tm.taskChan)
}

// resetTask 初始化任务的所有状态
func (tm *TaskManager[T, R, U]) resetTask(task *Task[T, R, U]) {
	task.state = Pending
	task.err = nil
	task.retryCount = 0
	task.fnDuration = 0
	task.callbackDuration = 0
	task.result = *new(R)         // 使用指针来初始化结果
	task.callbackResult = *new(U) // 使用指针来初始化回调结果
	task.callbackError = nil
	task.callbackState = Pending
}

// resetTasks 递归重置任务及其依赖任务的状态
func (tm *TaskManager[T, R, U]) resetTasks(task *Task[T, R, U]) {
	tm.resetTask(task) // 重置当前任务
	for _, dep := range task.depends {
		tm.resetTasks(dep) // 递归重置依赖任务
	}
}

// Release 释放资源
func (tm *TaskManager[T, R, U]) Release() {
	WithLock(&tm.mu, func() {
		for _, task := range tm.tasks {
			tm.resetTasks(task) // 合并重置逻辑
		}
		tm.isReleased = true
		tm.wg = sync.WaitGroup{}                                // 重置等待组
		tm.taskChan = make(chan *Task[T, R, U], tm.concurrency) // 重新初始化任务通道
	})
}

// cancelTaskAndDependencies 递归取消任务及其依赖
func (tm *TaskManager[T, R, U]) cancelTaskAndDependencies(task *Task[T, R, U]) {
	WithLock(&tm.mu, func() {
		// 检查任务是否已经取消
		if task.state == Cancelled {
			return
		}

		// 只在可取消的状态下执行
		if _, canCancel := cancellableStates[task.state]; canCancel {
			if task.cancel != nil {
				task.cancel() // 调用取消函数
			}
			task.state = Cancelled // 设置状态为已取消
		}
	})

	// 并发取消所有依赖任务
	for _, depTask := range task.depends {
		tm.wg.Add(1)
		go func(dep *Task[T, R, U]) {
			defer tm.wg.Done()
			tm.cancelTaskAndDependencies(dep)
		}(depTask)
	}
}

// CancelAll 取消所有任务
func (tm *TaskManager[T, R, U]) CancelAll() {
	for _, task := range tm.tasks {
		tm.wg.Add(1)
		go func(t *Task[T, R, U]) {
			defer tm.wg.Done()
			tm.cancelTaskAndDependencies(t) // 使用通用的取消逻辑
		}(task)
	}
}

// Cancel 取消某个任务
func (tm *TaskManager[T, R, U]) Cancel(taskName string) {
	if task, exists := tm.tasks[taskName]; exists {
		tm.cancelTaskAndDependencies(task) // 使用通用的取消逻辑
	}
}

// runTask 执行任务的具体逻辑
func (tm *TaskManager[T, R, U]) runTask(task *Task[T, R, U]) {
	// 检查任务是否已经被取消、完成或失败
	if task.state == Cancelled || task.state == Completed || task.state == Failed {
		return
	}

	// 根据依赖的执行模式执行依赖任务
	switch tm.dependExecutionMode {
	case Sequential:
		// 顺序执行依赖任务
		for _, dep := range task.depends {
			tm.runTask(dep) // 递归执行依赖任务
			if dep.state == Failed {
				task.state = Failed                                                    // 如果依赖任务失败，主任务也标记为失败
				task.err = fmt.Errorf("dependency '%s' failed: %w", dep.name, dep.err) // 创建新的错误信息,赋值依赖任务的错误信息给主任务
				return
			}
		}
	case Concurrent:
		var wg sync.WaitGroup
		for _, dep := range task.depends {
			wg.Add(1)
			go func(dep *Task[T, R, U]) {
				defer wg.Done()
				tm.runTask(dep) // 并发执行依赖任务
			}(dep)
		}
		wg.Wait() // 等待所有依赖任务完成

		// 检查依赖任务的状态
		for _, dep := range task.depends {
			if dep.state == Failed {
				task.state = Failed // 如果依赖任务失败，主任务也标记为失败
				task.err = dep.err  // 赋值依赖任务的错误信息给主任务
				return
			}
		}
	}

	// 开始主任务的执行
	task.result, task.err = tm.runWithRetries(task)

	// 调用回调函数并获取结果与错误
	task.callbackResult, task.callbackError = tm.invokeCallback(task, task.result, task.err)
	tm.logHistory(task) // 记录任务历史
}

// runWithRetries 函数处理任务的重试逻辑
func (tm *TaskManager[T, R, U]) runWithRetries(t *Task[T, R, U]) (result R, err error) {
	// 在重试次数小于最大重试次数的情况下进行重试
	for t.retryCount < t.maxRetries {
		// 如果任务被取消，直接返回
		if t.state == Cancelled {
			return result, nil
		}

		select {
		// 检查上下文是否已经被取消
		case <-t.ctx.Done():
			t.state = Cancelled // 将任务状态设置为取消
			return result, nil
		default:
			t.state = Running // 将任务状态设置为正在运行
			startTime := time.Now()
			// 执行任务函数
			result, err = t.fn(t.ctx, t.input)
			t.fnDuration = time.Since(startTime) // 记录任务执行时间
			// 如果没有错误，表示任务成功完成
			if err == nil {
				t.state = Completed // 任务成功完成
				return result, nil
			}

			// 处理任务错误
			tm.handleTaskError(t, err)
			time.Sleep(t.retryInterval) // 等待重试间隔
		}
	}
	return result, err // 返回最终结果和错误
}

// handleTaskError 错误处理
func (tm *TaskManager[T, R, U]) handleTaskError(t *Task[T, R, U], err error) {
	t.state = Failed // 任务执行失败状态
	t.err = err      // 任务执行失败错误信息
	t.retryCount++   // 增加重试计数
}

// invokeCallback 处理回调的执行
// 根据任务的执行结果调用相应的回调函数（成功或失败），并记录相关的执行时间和状态
func (tm *TaskManager[T, R, U]) invokeCallback(t *Task[T, R, U], result R, err error) (callbackResult U, callbackError error) {
	// 如果主任务已经出错，则不执行回调，直接返回
	if t.err != nil {
		return
	}

	// 定义一个映射，将错误状态映射到相应的回调函数
	// 如果没有错误，使用成功回调；如果有错误，使用失败回调
	callbackMap := map[bool]func(R, error) (U, error){
		false: t.successCallback, // err == nil
		true:  t.failureCallback, // err != nil
	}

	// 根据错误状态获取相应的回调函数
	callback := callbackMap[err != nil]

	// 如果回调函数存在，执行它
	if callback != nil {
		// 设置回调状态为正在运行
		t.callbackState = Running

		// 记录回调开始执行的时间
		callbackStartTime := time.Now()

		// 调用回调函数并处理返回值
		callbackResult, callbackError = callback(result, err)

		// 记录回调的运行时间
		t.callbackDuration = time.Since(callbackStartTime)

		// 存储回调的结果和错误
		t.callbackResult = callbackResult
		t.callbackError = callbackError

		// 设置回调状态为完成
		t.callbackState = Completed

		// 处理回调函数返回的错误
		if t.callbackError != nil {
			// 如果回调执行失败，设置任务状态为失败
			t.state = Failed
			t.callbackState = Failed // 更新回调状态为失败
		}
	}

	// 返回回调的结果和错误
	return callbackResult, callbackError
}

// GetTasks 获取所有任务
func (tm *TaskManager[T, R, U]) GetTasks() map[string]*Task[T, R, U] {
	return WithRLockReturnValue(&tm.mu, func() map[string]*Task[T, R, U] {
		return tm.tasks
	})
}

// logHistory 记录任务执行历史
func (tm *TaskManager[T, R, U]) logHistory(task *Task[T, R, U]) {
	WithLock(&tm.mu, func() {
		history := TaskHistory{
			Timestamp:             time.Now(),
			State:                 task.state,
			Result:                task.result,
			Error:                 task.err,
			FnExecutionTime:       task.fnDuration,
			CallbackExecutionTime: task.callbackDuration,
			TaskType:              task.taskType,
		}
		tm.history[task.name] = append(tm.history[task.name], history)
		// 限制历史记录的大小
		if tm.maxHistorySize > 0 && len(tm.history[task.name]) > tm.maxHistorySize {
			tm.history[task.name] = tm.history[task.name][1:] // 删除最旧的记录
		}
	})
}

// GetTaskHistory 获取任务执行历史
func (tm *TaskManager[T, R, U]) GetTaskHistory(taskName string) []TaskHistory {
	return WithRLockReturnValue(&tm.mu, func() []TaskHistory {
		return tm.history[taskName]
	})
}

// GetConcurrency 获取并发数设置
func (tm *TaskManager[T, R, U]) GetConcurrency() int {
	return WithRLockReturnValue(&tm.mu, func() int {
		return tm.concurrency
	})
}
