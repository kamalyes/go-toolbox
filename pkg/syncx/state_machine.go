/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2025-12-28 00:00:00
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2026-01-08 08:56:48
 * @FilePath: \go-toolbox\pkg\syncx\state_machine.go
 * @Description: 状态机实现
 *
 * Copyright (c) 2025 by kamalyes, All Rights Reserved.
 */

package syncx

import (
	"fmt"
	"sync"
	"time"
)

// StateMachine 状态机,用于管理状态转换和验证
// S 必须是可比较的类型 (comparable),通常是 string、int 或自定义枚举类型
//
// 使用示例:
//
//	type ConnectionState string
//	const (
//	    StateDisconnected ConnectionState = "disconnected"
//	    StateConnecting   ConnectionState = "connecting"
//	    StateConnected    ConnectionState = "connected"
//	)
//
//	sm := NewStateMachine(StateDisconnected)
//	sm.AllowTransition(StateDisconnected, StateConnecting)
//	sm.AllowTransition(StateConnecting, StateConnected)
//	sm.AllowTransition(StateConnected, StateDisconnected)
//
//	sm.OnTransition(func(from, to ConnectionState) {
//	    fmt.Printf("State changed: %s -> %s\n", from, to)
//	})
//
//	if err := sm.TransitionTo(StateConnecting); err != nil {
//	   处理错误
//	}

// StateTransition 状态转换记录
type StateTransition[S comparable] struct {
	From      S             // 转换前状态
	To        S             // 转换后状态
	Timestamp time.Time     // 转换时间
	Duration  time.Duration // 在 From 状态停留的时间
}

type StateMachine[S comparable] struct {
	currentState       S                    // 当前状态
	transitions        map[S]map[S]struct{} // 允许的状态转换 (from -> [to...])
	onTransition       []func(from, to S)   // 状态转换回调
	onEnter            map[S][]func(from S) // 进入状态回调
	onExit             map[S][]func(to S)   // 离开状态回调
	mu                 sync.RWMutex         // 读写锁
	allowAny           bool                 // 是否允许任意转换
	history            []StateTransition[S] // 状态转换历史记录
	maxHistory         int                  // 最大历史记录数量,0表示不限制
	trackHistory       bool                 // 是否追踪历史
	timeFormat         string               // 历史记录的时间格式,默认为 time.RFC3339
	lastTransitionTime time.Time            // 上次状态转换时间
}

// StateMachineOption 状态机配置选项
type StateMachineOption[S comparable] func(*StateMachine[S])

// WithAllowAnyTransition 允许任意状态转换
func WithAllowAnyTransition[S comparable]() StateMachineOption[S] {
	return func(sm *StateMachine[S]) {
		sm.allowAny = true
	}
}

// WithTrackHistory 启用历史记录追踪
// maxHistory 为最大历史记录数量,0表示不限制
//
// 示例:
//
//	sm := NewStateMachine("idle", WithTrackHistory(100)) // 最多保留100条历史
//	sm := NewStateMachine("idle", WithTrackHistory(0))   // 不限制历史记录数量
func WithTrackHistory[S comparable](maxHistory int) StateMachineOption[S] {
	return func(sm *StateMachine[S]) {
		sm.trackHistory = true
		sm.maxHistory = maxHistory
		if maxHistory > 0 {
			sm.history = make([]StateTransition[S], 0, maxHistory)
		} else {
			sm.history = make([]StateTransition[S], 0)
		}
	}
}

// WithTimeFormat 设置历史记录的时间格式
// timeFormat 为时间格式字符串,默认使用 time.RFC3339
//
// 示例:
//
//	sm := NewStateMachine("idle", WithTrackHistory(100), WithTimeFormat("2006-01-02 15:04:05"))
func WithTimeFormat[S comparable](timeFormat string) StateMachineOption[S] {
	return func(sm *StateMachine[S]) {
		if timeFormat != "" {
			sm.timeFormat = timeFormat
		}
	}
}

// NewStateMachine 创建一个新的状态机
// initialState 为初始状态
//
// 示例:
//
//	sm := NewStateMachine("idle")
//	sm := NewStateMachine("idle", WithAllowAnyTransition()) // 允许任意转换
func NewStateMachine[S comparable](initialState S, opts ...StateMachineOption[S]) *StateMachine[S] {
	now := time.Now()
	sm := &StateMachine[S]{
		currentState:       initialState,
		transitions:        make(map[S]map[S]struct{}),
		onTransition:       make([]func(from, to S), 0),
		onEnter:            make(map[S][]func(from S)),
		onExit:             make(map[S][]func(to S)),
		allowAny:           false,
		trackHistory:       false,
		maxHistory:         0,
		timeFormat:         time.RFC3339,
		lastTransitionTime: now,
	}

	for _, opt := range opts {
		opt(sm)
	}

	return sm
}

// CurrentState 获取当前状态
// 线程安全
func (sm *StateMachine[S]) CurrentState() S {
	sm.mu.RLock()
	defer sm.mu.RUnlock()
	return sm.currentState
}

// AllowTransition 允许从状态 from 转换到状态 to
// 可以多次调用以配置状态转换规则
//
// 示例:
//
//	sm.AllowTransition("pending", "approved")
//	sm.AllowTransition("pending", "rejected")
func (sm *StateMachine[S]) AllowTransition(from, to S) {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	if sm.transitions[from] == nil {
		sm.transitions[from] = make(map[S]struct{})
	}
	sm.transitions[from][to] = struct{}{}
}

// AllowTransitions 批量允许状态转换
// 参数 from 为起始状态, toStates 为可转换到的所有目标状态
//
// 示例:
//
//	sm.AllowTransitions("pending", "approved", "rejected", "cancelled")
func (sm *StateMachine[S]) AllowTransitions(from S, toStates ...S) {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	if sm.transitions[from] == nil {
		sm.transitions[from] = make(map[S]struct{})
	}
	for _, to := range toStates {
		sm.transitions[from][to] = struct{}{}
	}
}

// DisallowTransition 禁止从状态 from 转换到状态 to
func (sm *StateMachine[S]) DisallowTransition(from, to S) {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	if sm.transitions[from] != nil {
		delete(sm.transitions[from], to)
	}
}

// CanTransitionTo 检查是否可以转换到目标状态
// 线程安全
//
// 示例:
//
//	if sm.CanTransitionTo("approved") {
//	    可以转换
//	}
func (sm *StateMachine[S]) CanTransitionTo(to S) bool {
	sm.mu.RLock()
	defer sm.mu.RUnlock()

	// 如果允许任意转换
	if sm.allowAny {
		return true
	}

	// 检查是否允许从当前状态转换到目标状态
	if sm.transitions[sm.currentState] == nil {
		return false
	}
	_, ok := sm.transitions[sm.currentState][to]
	return ok
}

// TransitionTo 转换到目标状态
// 如果转换不被允许,返回错误
// 转换成功时会触发所有已注册的回调函数
//
// 示例:
//
//	if err := sm.TransitionTo("approved"); err != nil {
//	    log.Printf("Transition failed: %v", err)
//	}
func (sm *StateMachine[S]) TransitionTo(to S) error {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	from := sm.currentState

	// 如果已经是目标状态,直接返回
	if from == to {
		return nil
	}

	// 检查转换是否被允许
	if !sm.allowAny {
		if sm.transitions[from] == nil {
			return fmt.Errorf("no transitions defined for state: %v", from)
		}
		if _, ok := sm.transitions[from][to]; !ok {
			return fmt.Errorf("transition from %v to %v is not allowed", from, to)
		}
	}

	// 触发离开回调
	if exitCallbacks, ok := sm.onExit[from]; ok {
		for _, callback := range exitCallbacks {
			callback(to)
		}
	}

	// 更新状态
	sm.currentState = to

	// 记录历史
	if sm.trackHistory {
		now := time.Now()
		duration := now.Sub(sm.lastTransitionTime)
		transition := StateTransition[S]{
			From:      from,
			To:        to,
			Timestamp: now,
			Duration:  duration,
		}
		sm.history = append(sm.history, transition)
		sm.lastTransitionTime = now

		// 如果设置了最大历史记录数量且超过限制,移除最旧的记录
		if sm.maxHistory > 0 && len(sm.history) > sm.maxHistory {
			sm.history = sm.history[len(sm.history)-sm.maxHistory:]
		}
	}

	// 触发进入回调
	if enterCallbacks, ok := sm.onEnter[to]; ok {
		for _, callback := range enterCallbacks {
			callback(from)
		}
	}

	// 触发转换回调
	for _, callback := range sm.onTransition {
		callback(from, to)
	}

	return nil
}

// GetLastTransitionTime 获取上次状态转换的时间
func (sm *StateMachine[S]) GetLastTransitionTime() time.Time {
	sm.mu.RLock()
	defer sm.mu.RUnlock()
	return sm.lastTransitionTime
}

// GetLastTransitionDuration 获取上次状态转换的持续时间（在上一个状态停留的时间）
// 如果没有历史记录，返回 0
func (sm *StateMachine[S]) GetLastTransitionDuration() time.Duration {
	sm.mu.RLock()
	defer sm.mu.RUnlock()

	if !sm.trackHistory || len(sm.history) == 0 {
		return 0
	}

	return sm.history[len(sm.history)-1].Duration
}

// GetLastTransition 获取最后一次状态转换记录
// 如果没有历史记录，返回零值和 false
func (sm *StateMachine[S]) GetLastTransition() (StateTransition[S], bool) {
	sm.mu.RLock()
	defer sm.mu.RUnlock()

	if !sm.trackHistory || len(sm.history) == 0 {
		return StateTransition[S]{}, false
	}

	return sm.history[len(sm.history)-1], true
}

// GetTotalDuration 获取从状态机创建到当前的总时长
func (sm *StateMachine[S]) GetTotalDuration() time.Duration {
	sm.mu.RLock()
	defer sm.mu.RUnlock()

	if !sm.trackHistory || len(sm.history) == 0 {
		// 如果没有历史记录，返回从创建到现在的时间
		return time.Since(sm.lastTransitionTime)
	}

	// 返回从第一次转换到现在的时间
	firstTransition := sm.history[0].Timestamp
	return time.Since(firstTransition)
}

// MustTransitionTo 强制转换到目标状态,如果失败则 panic
func (sm *StateMachine[S]) MustTransitionTo(to S) {
	if err := sm.TransitionTo(to); err != nil {
		panic(err)
	}
}

// OnTransition 注册状态转换回调
// 每次状态转换时都会调用所有已注册的回调函数
//
// 示例:
//
//	sm.OnTransition(func(from, to string) {
//	    log.Printf("State changed: %s -> %s", from, to)
//	})
func (sm *StateMachine[S]) OnTransition(handler func(from, to S)) {
	sm.mu.Lock()
	defer sm.mu.Unlock()
	sm.onTransition = append(sm.onTransition, handler)
}

// OnEnter 注册进入某个状态时的回调
// 参数 from 为转换前的状态
//
// 示例:
//
//	sm.OnEnter("connected", func(from string) {
//	    log.Printf("Entered connected state from %s", from)
//	})
func (sm *StateMachine[S]) OnEnter(state S, handler func(from S)) {
	sm.mu.Lock()
	defer sm.mu.Unlock()
	sm.onEnter[state] = append(sm.onEnter[state], handler)
}

// OnExit 注册离开某个状态时的回调
// 参数 to 为转换后的目标状态
//
// 示例:
//
//	sm.OnExit("connected", func(to string) {
//	    log.Printf("Exiting connected state to %s", to)
//	})
func (sm *StateMachine[S]) OnExit(state S, handler func(to S)) {
	sm.mu.Lock()
	defer sm.mu.Unlock()
	sm.onExit[state] = append(sm.onExit[state], handler)
}

// GetAllowedTransitions 获取当前状态允许转换到的所有状态
func (sm *StateMachine[S]) GetAllowedTransitions() []S {
	sm.mu.RLock()
	defer sm.mu.RUnlock()

	if sm.allowAny {
		return []S{} // 无法枚举所有可能的状态
	}

	transitions := sm.transitions[sm.currentState]
	if transitions == nil {
		return []S{}
	}

	result := make([]S, 0, len(transitions))
	for state := range transitions {
		result = append(result, state)
	}
	return result
}

// Reset 重置状态机到初始状态
func (sm *StateMachine[S]) Reset(initialState S) {
	sm.mu.Lock()
	defer sm.mu.Unlock()
	sm.currentState = initialState
}

// ClearCallbacks 清除所有回调函数
func (sm *StateMachine[S]) ClearCallbacks() {
	sm.mu.Lock()
	defer sm.mu.Unlock()
	sm.onTransition = make([]func(from, to S), 0)
	sm.onEnter = make(map[S][]func(from S))
	sm.onExit = make(map[S][]func(to S))
}

// GetHistory 获取状态转换历史记录
// 返回所有历史记录的副本,按时间顺序排列(最早的在前)
//
// 示例:
//
//	history := sm.GetHistory()
//	for _, trans := range history {
//	    fmt.Printf("%s: %v -> %v\n", trans.Timestamp.Format(time.RFC3339), trans.From, trans.To)
//	}
func (sm *StateMachine[S]) GetHistory() []StateTransition[S] {
	sm.mu.RLock()
	defer sm.mu.RUnlock()

	if !sm.trackHistory {
		return []StateTransition[S]{}
	}

	// 返回副本以避免外部修改
	history := make([]StateTransition[S], len(sm.history))
	copy(history, sm.history)
	return history
}

// GetRecentHistory 获取最近 n 条状态转换历史记录
// 如果 n 大于实际历史记录数量,返回所有历史记录
//
// 示例:
//
//	recentHistory := sm.GetRecentHistory(5) // 获取最近5条记录
func (sm *StateMachine[S]) GetRecentHistory(n int) []StateTransition[S] {
	sm.mu.RLock()
	defer sm.mu.RUnlock()

	if !sm.trackHistory || n <= 0 {
		return []StateTransition[S]{}
	}

	start := len(sm.history) - n
	if start < 0 {
		start = 0
	}

	history := make([]StateTransition[S], len(sm.history)-start)
	copy(history, sm.history[start:])
	return history
}

// ClearHistory 清除历史记录
func (sm *StateMachine[S]) ClearHistory() {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	if sm.trackHistory {
		if sm.maxHistory > 0 {
			sm.history = make([]StateTransition[S], 0, sm.maxHistory)
		} else {
			sm.history = make([]StateTransition[S], 0)
		}
	}
}

// GetHistoryCount 获取历史记录数量
func (sm *StateMachine[S]) GetHistoryCount() int {
	sm.mu.RLock()
	defer sm.mu.RUnlock()
	return len(sm.history)
}

// IsTrackingHistory 检查是否正在追踪历史记录
func (sm *StateMachine[S]) IsTrackingHistory() bool {
	sm.mu.RLock()
	defer sm.mu.RUnlock()
	return sm.trackHistory
}

// formatHistorySlice 将历史记录切片格式化为字符串
// 这是一个内部辅助方法,用于复用字符串格式化逻辑
func (sm *StateMachine[S]) formatHistorySlice(history []StateTransition[S]) string {
	if len(history) == 0 {
		return ""
	}

	var result string
	for i, trans := range history {
		if i > 0 {
			result += "\n"
		}
		result += fmt.Sprintf("%s: %v -> %v (duration: %v)",
			trans.Timestamp.Format(sm.timeFormat),
			trans.From,
			trans.To,
			trans.Duration)
	}
	return result
}

// GetHistoryString 获取完整执行流程的字符串表示
// 使用创建状态机时配置的时间格式
// 返回格式化后的历史记录字符串,每行一条记录
//
// 示例:
//
//	historyStr := sm.GetHistoryString()
//	fmt.Println(historyStr)
func (sm *StateMachine[S]) GetHistoryString() string {
	sm.mu.RLock()
	defer sm.mu.RUnlock()

	if !sm.trackHistory {
		return ""
	}

	return sm.formatHistorySlice(sm.history)
}

// GetRecentHistoryString 获取最近 n 条执行流程的字符串表示
// n 为记录数量，使用创建状态机时配置的时间格式
//
// 示例:
//
//	historyStr := sm.GetRecentHistoryString(5)
func (sm *StateMachine[S]) GetRecentHistoryString(n int) string {
	sm.mu.RLock()
	defer sm.mu.RUnlock()

	if !sm.trackHistory || n <= 0 {
		return ""
	}

	start := len(sm.history) - n
	if start < 0 {
		start = 0
	}

	return sm.formatHistorySlice(sm.history[start:])
}
