/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2024-11-26 15:55:27
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2024-12-05 20:02:54
 * @FilePath: \go-toolbox\pkg\moment\timer.go
 * @Description: Cron表达式是一种用于指定定时任务的时间表达式，常用于设置任务的执行时间、频率和间隔。Cron表达式由6或7个字段组成，分别表示秒、分、时、日期、月份、星期和年份。
 * Cron表达式的语法规则：
 *
 * 字段      | 是否必需 | 取值范围
 * ----------|----------|----------------
 * 秒        | 是       | [0, 59]
 * 分        | 是       | [0, 59]
 * 时        | 是       | [0, 23]
 * 日期      | 是       | [1, 31]，需要考虑月的天数
 * 月份      | 是       | [1, 12]或[JAN, DEC]
 * 星期      | 是       | [0, 7]或[SUN, SAT]（0=SUN）
 * 年份      | 否       | [1970, 2099]
 *
 * 特殊字符：
 *
 * 字符 | 含义
 * -----|-----------------
 * *   | 匹配任意值
 * ,   | 列出枚举值
 * -   | 指定范围
 * /   | 指定数值的增量
 * ?   | 不指定值，仅用于日期和星期
 * L   | 表示最后一天，仅字段日期和星期支持
 * W   | 除周末以外的有效工作日
 * #   | 确定每个月的第几个星期几
 *
 * 常用的Cron表达式示例：
 *
 * 场景                                   | 表达式
 * --------------------------------------|-------------------------
 * 每天上午10:15执行任务                 | 0 15 10 ? * *
 * 每月15日上午10:15执行任务             | 0 15 10 15 * ?
 * 每个星期三中午12:00执行任务           | 0 0 12 ? * WED
 * 每月最后一天执行任务                   | 0 0 1 L * ?
 * Copyright (c) 2024 by kamalyes, All Rights Reserved.
 */
package moment

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/kamalyes/go-toolbox/pkg/stringx"
)

// TaskInfo 结构体表示任务的执行信息
type TaskInfo struct {
	methodName string // 方法名称
	callerInfo string // 调用者信息
	execCount  int    // 执行次数
}

// GetMethodName 获取方法名称
func (t *TaskInfo) GetMethodName() string {
	return t.methodName
}

// GetCallerInfo 获取调用者信息
func (t *TaskInfo) GetCallerInfo() string {
	return t.callerInfo
}

// GetExecCount 获取执行次数
func (t *TaskInfo) GetExecCount() int {
	return t.execCount
}

type CustomFunc func() error // 自定义函数类型

// Rule 结构体表示一个调度规则
type Rule struct {
	expression string       // Cron 表达式
	callback   func() error // 任务回调函数
	beforeFunc func()       // 执行前的函数
	afterFunc  func()       // 执行后的函数
	skipFunc   func() bool  // 跳过执行的条件函数
	customFunc CustomFunc   // 自定义函数
}

// GetExpression 获取 Cron 表达式
func (r *Rule) GetExpression() string {
	return r.expression
}

// GetCallback 获取任务回调函数
func (r *Rule) GetCallback() func() error {
	return r.callback
}

// GetBeforeFunc 获取执行前的函数
func (r *Rule) GetBeforeFunc() func() {
	return r.beforeFunc
}

// GetAfterFunc 获取执行后的函数
func (r *Rule) GetAfterFunc() func() {
	return r.afterFunc
}

// GetSkipFunc 获取跳过执行的条件函数
func (r *Rule) GetSkipFunc() func() bool {
	return r.skipFunc
}

// GetCustomFunc 获取自定义函数
func (r *Rule) GetCustomFunc() CustomFunc {
	return r.customFunc
}

// Timer 结构体表示一个定时器
type Timer struct {
	ctx              context.Context
	rules            []*Rule               // 存储调度规则
	cooldownDuration time.Duration         // 冷却时间
	sleepDuration    time.Duration         // 睡眠时间
	isBroken         bool                  // 标记定时器是否处于故障状态
	failureCount     int                   // 失败计数
	maxFailures      int                   // 最大失败次数
	tasks            map[string][]TaskInfo // 任务执行记录
	log              []string              // 日志记录
	mu               sync.RWMutex          // 读写锁
	logMu            sync.Mutex            // 添加一个新的锁来保护日志记录
	wg               sync.WaitGroup        // 用于等待任务完成
	taskChan         chan *Rule            // 任务通道
	workerCount      int                   // 工作池大小
	once             sync.Once             // 确保只关闭一次
	stopFlag         chan struct{}         // 用于通知停止
}

// GetCooldownDuration 获取冷却时间
func (t *Timer) GetCooldownDuration() time.Duration {
	return t.cooldownDuration
}

// GetSleepDuration 获取睡眠时间
func (t *Timer) GetSleepDuration() time.Duration {
	return t.sleepDuration
}

// GetMaxFailures 获取最大失败次数
func (t *Timer) GetMaxFailures() int {
	return t.maxFailures
}

// GetFailureCount 获取当前失败次数
func (t *Timer) GetFailureCount() int {
	return t.failureCount
}

// GetTasksCopy 获取任务执行记录的副本
func (t *Timer) GetTasksCopy() map[string][]TaskInfo {
	tasksCopy := make(map[string][]TaskInfo)
	for k, v := range t.tasks {
		tasksCopy[k] = v
	}
	return tasksCopy
}

// GetTasks 获取任务执行记录
func (t *Timer) GetTasks() map[string][]TaskInfo {
	return t.tasks
}

// GetTask 获取特定表达式的任务执行记录
func (t *Timer) GetTask(expression string) ([]TaskInfo, bool) {
	taskInfo, exists := t.tasks[expression]
	return taskInfo, exists
}

// GetLogCopy 获取日志记录的副本
func (t *Timer) GetLogCopy() []string {
	logCopy := make([]string, len(t.log)) // 创建日志副本
	copy(logCopy, t.log)                  // 复制日志内容
	return logCopy
}

// GetLog 获取日志记录
func (t *Timer) GetLog() []string {
	return t.log
}

// GetRulesCopy 获取调度规则的副本
func (t *Timer) GetRulesCopy() []*Rule {
	rulesCopy := make([]*Rule, len(t.rules)) // 创建规则副本
	copy(rulesCopy, t.rules)                 // 复制规则内容
	return rulesCopy
}

// GetRules 获取调度规则
func (t *Timer) GetRules() []*Rule {
	return t.rules
}

// IsBroken 获取定时器是否处于熔断状态
func (t *Timer) IsBroken() bool {
	return t.isBroken
}

// GetWorkerCount 获取工作池大小
func (t *Timer) GetWorkerCount() int {
	return t.workerCount
}

// GetLogCount 获取日志条目数量
func (t *Timer) GetLogCount() int {
	return len(t.log)
}

// NewTimer 创建一个新的 Timer 实例
func NewTimer() *Timer {
	return NewTimerWithCtx(context.Background())
}

// NewCustomTimer 创建一个新的 自定义Timer 实例
func NewCustomTimer(timer *Timer) *Timer {
	return timer
}

// NewTimerWithCtx 创建一个新的 Timer Ctx实例
func NewTimerWithCtx(ctx context.Context) *Timer {
	return &Timer{
		ctx:              ctx,
		rules:            []*Rule{},
		cooldownDuration: 0,
		sleepDuration:    1 * time.Second, // 默认睡眠时间
		tasks:            make(map[string][]TaskInfo),
		maxFailures:      3,
		log:              []string{},
		taskChan:         make(chan *Rule, 100), // 创建带缓冲的通道
		workerCount:      5,                     // 默认工作池大小
	}
}

// 定义调度规则的类型
type ScheduleType int

// 调度规则的常量
const (
	EverySecond ScheduleType = iota
	EveryMinute
	EveryHalfMinute // 每半分钟
	EveryHour
	EveryHalfHour // 每半小时
	EveryDay
	EveryHalfDay // 每半天
	EveryWeek
	EveryMonth
	EveryYear
	WeekdaysOnly   // 仅工作日
	WeekendsOnly   // 仅非工作日
	PeakHours      // 高峰期
	OffPeakHours   // 低谷期
	NewYear        // 元旦
	SpringFestival // 春节
	NationalDay    // 国庆节
	Christmas      // 圣诞节
	Monday         // 星期一
	Tuesday        // 星期二
	Wednesday      // 星期三
	Thursday       // 星期四
	Friday         // 星期五
	Saturday       // 星期六
	Sunday         // 星期日
)

// 调度规则映射
var ScheduleRules = map[ScheduleType]string{
	EverySecond:     "0 * * * * *",       // 每秒
	EveryMinute:     "0 */1 * * * *",     // 每分钟
	EveryHalfMinute: "*/30 * * * * *",    // 每半分钟
	EveryHour:       "0 0 */1 * * *",     // 每小时
	EveryHalfHour:   "0 */30 * * * *",    // 每半小时
	EveryDay:        "0 0 1 * * *",       // 每天
	EveryHalfDay:    "0 0,12 * * *",      // 每半天（每天的 0:00 和 12:00）
	EveryWeek:       "0 0 * * 0",         // 每周
	EveryMonth:      "0 0 1 * *",         // 每月
	EveryYear:       "0 0 1 1 *",         // 每年
	WeekdaysOnly:    "0 0 * * 1-5",       // 仅工作日（周一到周五）
	WeekendsOnly:    "0 0 * * 6,0",       // 仅非工作日（周六和周日）
	PeakHours:       "0 8-9,17-18 * * *", // 高峰期（例如 8:00 - 9:00 和 17:00 - 18:00）
	OffPeakHours:    "0 9-17 * * *",      // 低谷期（例如 9:00 - 17:00）
	Monday:          "0 0 * * 1",         // 星期一
	Tuesday:         "0 0 * * 2",         // 星期二
	Wednesday:       "0 0 * * 3",         // 星期三
	Thursday:        "0 0 * * 4",         // 星期四
	Friday:          "0 0 * * 5",         // 星期五
	Saturday:        "0 0 * * 6",         // 星期六
	Sunday:          "0 0 * * 0",         // 星期日
}

// String 返回调度规则的 Cron 表达式字符串
func (s ScheduleType) String() (string, error) {
	if rule, exists := ScheduleRules[s]; exists {
		return rule, nil
	}
	return "", fmt.Errorf("invalid schedule type: %d", s)
}

// SetDefaultScheduleRule 设置指定调度类型的 cron 表达式
func (t *Timer) SetDefaultScheduleRule(scheduleType ScheduleType, callback func() error) *Timer {
	t.AddRule(ScheduleRules[scheduleType], callback)
	return t
}

// SetCooldownDuration 设置冷却时间
func (t *Timer) SetCooldownDuration(duration time.Duration) *Timer {
	t.cooldownDuration = duration
	return t
}

// SetSleepDuration 设置睡眠时间
func (t *Timer) SetSleepDuration(duration time.Duration) *Timer {
	t.sleepDuration = duration
	return t
}

// SetTaskChanCapacity 设置任务通道的容量
func (t *Timer) SetTaskChanCapacity(capacity int) *Timer {
	t.mu.Lock() // 锁定以确保安全
	defer t.mu.Unlock()

	// 关闭旧通道，创建新通道
	close(t.taskChan)
	t.taskChan = make(chan *Rule, capacity)
	return t
}

// SetWorkerCount 设置工作池的大小
func (t *Timer) SetWorkerCount(count int) *Timer {
	t.workerCount = count
	return t
}

// SetMaxFailures 设置最大失败次数
func (t *Timer) SetMaxFailures(maxFailures int) *Timer {
	t.maxFailures = maxFailures
	return t
}

// SetCustomFunc 设置自定义函数
func (t *Timer) SetCustomFunc(customFunc CustomFunc) *Timer {
	if len(t.rules) > 0 {
		t.rules[len(t.rules)-1].customFunc = customFunc
	}
	return t
}

// AddRule 添加调度规则
func (t *Timer) AddRule(expression string, callback func() error) *Timer {
	rule := &Rule{expression: expression, callback: callback}
	t.rules = append(t.rules, rule)
	return t
}

// setDecorator 设置装饰器
func (t *Timer) setDecorator(decoratorFunc func(rule *Rule)) *Timer {
	if len(t.rules) > 0 {
		decoratorFunc(t.rules[len(t.rules)-1])
	}
	return t
}

// Before 添加 Before 装饰器
func (t *Timer) Before(beforeFunc func()) *Timer {
	return t.setDecorator(func(rule *Rule) {
		rule.beforeFunc = beforeFunc
	})
}

// After 添加 After 装饰器
func (t *Timer) After(afterFunc func()) *Timer {
	return t.setDecorator(func(rule *Rule) {
		rule.afterFunc = afterFunc
	})
}

// Skip 添加 Skip 装饰器
func (t *Timer) Skip(skipFunc func() bool) *Timer {
	return t.setDecorator(func(rule *Rule) {
		rule.skipFunc = skipFunc
	})
}

// Validate 检查 Timer 的配置是否有效
func (t *Timer) Validate() error {
	if t.sleepDuration <= 0 {
		return errors.New("sleep duration must be greater than zero")
	}
	if len(t.rules) == 0 {
		return errors.New("at least one rule must be defined")
	}
	return nil
}

// Start 启动定时器
func (t *Timer) Start() error {
	if err := t.Validate(); err != nil {
		return err
	}

	t.stopFlag = make(chan struct{}) // 初始化停止标志通道

	// 启动工作池
	for i := 0; i < t.workerCount; i++ {
		go t.worker()
	}

	go t.run() // 启动定时器运行逻辑
	return nil
}

// run 包含定时器的主要运行逻辑
func (t *Timer) run() {
	for {
		select {
		case <-t.ctx.Done():
			t.Stop() // 在上下文取消时停止定时器
			return
		case <-t.stopFlag: // 检查停止标志
			return
		default:
			t.processRules() // 处理规则
		}
	}
}

// processRules 处理所有规则
func (t *Timer) processRules() {
	if t.isBroken {
		t.recoverTimer()
	} else {
		for _, rule := range t.rules {
			if !t.sendTask(rule) {
				t.logAction("Task channel is full, skipping task.", "System")
			}
		}
		time.Sleep(t.sleepDuration) // 使用配置的睡眠时间
	}
}

// sendTask 尝试将任务放入通道
func (t *Timer) sendTask(rule *Rule) bool {
	select {
	case t.taskChan <- rule: // 将任务放入通道
		return true
	default:
		return false // 通道已满
	}
}

// recoverTimer 恢复定时器状态
func (t *Timer) recoverTimer() {
	time.AfterFunc(t.cooldownDuration, func() {
		t.logAction("Timer is recovering...", "System")
		t.isBroken = false
		t.failureCount = 0
	})
}

// worker 执行任务的工作协程
func (t *Timer) worker() {
	for {
		select {
		case <-t.ctx.Done(): // 响应上下文取消
			return
		case rule, ok := <-t.taskChan:
			if !ok {
				return // 通道已关闭，退出
			}
			// 直接执行规则，而不是使用 goroutine
			t.executeRule(rule)
		case <-t.stopFlag: // 检查停止标志
			return
		}
	}
}

// Stop 停止定时器
func (t *Timer) Stop() {
	t.once.Do(func() {
		close(t.stopFlag) // 关闭停止标志
		close(t.taskChan) // 关闭任务通道
		t.wg.Wait()       // 等待所有任务完成
		t.logAction("Timer stopped.", "System")
	})
}

// Break 熔断定时器
func (t *Timer) Break() {
	t.mu.Lock()         // 写锁
	defer t.mu.Unlock() // 确保解锁
	t.isBroken = true
	t.logAction("Timer is broken!", "System")
}

// executeRule 执行单个规则
func (t *Timer) executeRule(rule *Rule) {
	now := time.Now()
	t.logAction(fmt.Sprintf("Executing rule: %s at %s", rule.expression, now), "System")

	if err := t.executeCustomFunction(rule); err != nil {
		t.logAction(fmt.Sprintf("Error executing custom function for rule '%s': %v", rule.expression, err), "System")
		return
	}

	nextRunTime, err := t.getNextRunTime(now, rule.expression)
	if err != nil {
		t.logAction(fmt.Sprintf("Error getting next run time for rule '%s': %v", rule.expression, err), "System")
		return
	}

	if nextRunTime != nil && nextRunTime.Before(now) {
		if t.shouldSkipRule(rule) {
			return
		}

		t.executeOptionalFunc(rule.beforeFunc) // 执行前函数
		if err := t.executeCallback(rule); err != nil {
			t.handleCallbackError(rule, err)
			return
		}

		t.recordTask(rule.expression, "User     ")
		t.executeOptionalFunc(rule.afterFunc) // 执行后函数
	}
}

// executeOptionalFunc 执行可选函数
func (t *Timer) executeOptionalFunc(f func()) {
	if f != nil {
		f()
	}
}

// executeCustomFunction 执行自定义函数
func (t *Timer) executeCustomFunction(rule *Rule) error {
	if rule.customFunc != nil {
		return rule.customFunc()
	}
	return nil
}

// shouldSkipRule 判断是否跳过规则
func (t *Timer) shouldSkipRule(rule *Rule) bool {
	if rule.skipFunc != nil && rule.skipFunc() {
		t.logAction(fmt.Sprintf("Skipping rule with expression '%s'", rule.expression), "System")
		return true
	}
	return false
}

// executeCallback 执行回调并处理错误
func (t *Timer) executeCallback(rule *Rule) error {
	if rule.callback != nil {
		return rule.callback()
	}
	return nil
}

// handleCallbackError 处理回调错误
func (t *Timer) handleCallbackError(rule *Rule, err error) {
	t.logAction(fmt.Sprintf("Error executing rule '%s': %v", rule.expression, err), "System")
	t.recordTask(rule.expression, "System") // 记录任务
	t.failureCount++                        // 增加失败计数
	if t.failureCount >= t.maxFailures {
		t.Break() // 达到最大失败次数，熔断
	}
}

// recordTask 记录任务执行次数和调用者信息
func (t *Timer) recordTask(expression string, caller string) {
	t.mu.Lock()         // 写锁
	defer t.mu.Unlock() // 确保解锁

	// 检查该表达式是否已经存在记录
	if _, exists := t.tasks[expression]; !exists {
		t.tasks[expression] = []TaskInfo{}
	}

	// 查找是否已经有该记录
	for i, taskInfo := range t.tasks[expression] {
		if taskInfo.methodName == "Callback" && taskInfo.callerInfo == caller {
			// 如果找到，增加执行次数
			t.tasks[expression][i].execCount++
			// 记录日志
			t.logAction(fmt.Sprintf("Task with expression '%s' executed by %s, new count: %d", expression, caller, t.tasks[expression][i].execCount), caller)
			return
		}
	}

	// 如果没有找到，添加新的记录
	t.tasks[expression] = append(t.tasks[expression], TaskInfo{
		methodName: "Callback",
		callerInfo: caller,
		execCount:  1,
	})

	// 记录日志
	t.logAction(fmt.Sprintf("Task with expression '%s' executed by %s", expression, caller), caller)
}

// logAction 记录日志
func (t *Timer) logAction(action string, caller string) {
	timestamp := time.Now().Format(time.RFC3339)
	logEntry := fmt.Sprintf("[%s] [%s] %s", timestamp, caller, action)

	// 使用写锁保护日志
	t.logMu.Lock()         // 使用写锁
	defer t.logMu.Unlock() // 确保解锁
	t.log = append(t.log, logEntry)
	fmt.Println(logEntry) // 打印日志
}

// getNextRunTime 计算下一个运行时间
func (t *Timer) getNextRunTime(now time.Time, expression string) (*time.Time, error) {
	cronParts := strings.Fields(expression)
	if len(cronParts) < 6 || len(cronParts) > 7 {
		return nil, errors.New("invalid cron expression")
	}

	var second, minute, hour, day, month, week string
	if len(cronParts) == 6 {
		second, minute, hour, day, month, week = cronParts[0], cronParts[1], cronParts[2], cronParts[3], cronParts[4], cronParts[5]
	} else {
		second, minute, hour, day, month, week = cronParts[1], cronParts[2], cronParts[3], cronParts[4], cronParts[5], cronParts[6]
	}

	// 处理 '?' 字符
	if week == "?" {
		week = "*" // 将 '?' 转换为 '*' 以避免错误
	}

	nextTime := time.Date(now.Year(), now.Month(), now.Day(), now.Hour(), now.Minute(), now.Second(), 0, now.Location())

	// 验证并计算下一个时间字段
	if err := t.validateAndCalculateField(&nextTime, now, second, minute, hour, day, month, week); err != nil {
		return nil, err
	}

	return &nextTime, nil
}

// validateAndCalculateField 验证并计算下一个时间字段
func (t *Timer) validateAndCalculateField(nextTime *time.Time, now time.Time, second, minute, hour, day, month, week string) error {
	// 计算下一个时间字段
	nextSecond := t.getNextField(now.Second(), second, 60)
	nextMinute := t.getNextField(now.Minute(), minute, 60)
	nextHour := t.getNextField(now.Hour(), hour, 24)
	nextDay, err := t.getNextDay(now, day, nextHour)
	if err != nil {
		return err
	}
	nextMonthNum, err := t.getNextMonth(int(now.Month()), month)
	if err != nil {
		return err
	}
	// 更新 nextTime
	*nextTime = nextTime.Add(time.Duration(nextSecond-now.Second()) * time.Second)
	*nextTime = nextTime.Add(time.Duration(nextMinute-now.Minute()) * time.Minute)
	*nextTime = nextTime.Add(time.Duration(nextHour-now.Hour()) * time.Hour)
	*nextTime = nextTime.AddDate(0, nextMonthNum-int(now.Month()), nextDay-now.Day())

	// 处理星期字段
	if week != "*" {
		weekDay := int(now.Weekday())
		nextWeekDay, err := t.getNextWeek(weekDay, week)
		if err != nil {
			return err
		}
		if nextWeekDay != weekDay {
			*nextTime = nextTime.AddDate(0, 0, nextWeekDay-weekDay)
		}
	}

	if nextTime.Before(now) {
		*nextTime = nextTime.AddDate(0, 0, 1)
	}

	return nil
}

// getNextDay 计算下一个有效的日期
func (t *Timer) getNextDay(now time.Time, day string, nextHour int) (int, error) {
	// 处理日期字段
	nextDay := t.getNextField(now.Day(), day, 31)

	// 检查下一个月份的天数
	nextMonth := int(now.Month())
	if nextHour == 0 && nextDay == 1 {
		nextMonth = (nextMonth + 1) % 12
		if nextMonth == 0 {
			nextMonth = 1
		}
	}
	daysInNextMonth := DaysInMonth(nextMonth, now.Year())
	if nextDay > daysInNextMonth {
		nextDay = daysInNextMonth
	}

	return nextDay, nil
}

// getNextMonth 计算下一个有效的月份
func (t *Timer) getNextMonth(currentMonth int, month string) (int, error) {
	if month == "*" {
		return (currentMonth % 12) + 1, nil
	}
	if strings.Contains(month, ",") {
		options := strings.Split(month, ",")
		for _, option := range options {
			if val, err := stringx.ParseMonth(option); err == nil && val > currentMonth {
				return val, nil
			}
		}
		return (currentMonth % 12) + 1, nil
	}
	return stringx.ParseMonth(month)
}

// getNextWeek 计算下一个有效的星期
func (t *Timer) getNextWeek(currentWeek int, week string) (int, error) {
	if week == "*" {
		return (currentWeek % 7) + 1, nil // 1=MON, 7=SUN
	}
	if strings.Contains(week, ",") {
		options := strings.Split(week, ",")
		for _, option := range options {
			if val, err := stringx.ParseWeek(option); err == nil && val > currentWeek {
				return val, nil
			}
		}
		return (currentWeek % 7) + 1, nil // 处理完所有选项后返回下一个星期
	}
	return stringx.ParseWeek(week)
}

// getNextField 计算下一个字段的值
func (t *Timer) getNextField(current int, field string, max int) int {
	switch {
	case field == "*":
		return (current + 1) % max
	case strings.Contains(field, ","):
		return t.handleCommaSeparated(current, field)
	case strings.Contains(field, "-"):
		return t.handleRange(current, field)
	case strings.HasPrefix(field, "*/"):
		return t.handleStep(current, field, max)
	default:
		return t.handleValue(current, field, max)
	}
}

// handleValue 处理普通值字段
func (t *Timer) handleValue(current int, field string, max int) int {
	if val, err := strconv.Atoi(field); err == nil && val >= 0 && val < max {
		if val > current {
			return val
		}
	}
	return (current + 1) % max // 默认返回下一个值
}

// handleCommaSeparated 处理逗号分隔的字段
func (t *Timer) handleCommaSeparated(current int, field string) int {
	options := strings.Split(field, ",")
	for _, option := range options {
		if val, err := strconv.Atoi(option); err == nil && val > current {
			return val
		}
	}
	return (current + 1) % 60 // 默认返回下一个值
}

// handleRange 处理范围字段
func (t *Timer) handleRange(current int, field string) int {
	rangeParts := strings.Split(field, "-")
	if len(rangeParts) != 2 {
		return (current + 1) % 60 // 默认返回下一个值
	}

	start, err1 := strconv.Atoi(rangeParts[0])
	end, err2 := strconv.Atoi(rangeParts[1])
	if err1 == nil && err2 == nil {
		if current < start {
			return start
		} else if current >= end {
			return start
		}
		return current + 1
	}
	return (current + 1) % 60 // 默认返回下一个值
}

// handleStep 处理步进字段
func (t *Timer) handleStep(current int, field string, max int) int {
	step, err := strconv.Atoi(strings.TrimPrefix(field, "*/"))
	if err == nil {
		return (current + step) % max
	}
	return (current + 1) % max // 默认返回下一个值
}
