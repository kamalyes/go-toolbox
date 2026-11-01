/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2026-10-31 21:36:52
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2026-10-31 21:58:03
 * @FilePath: \go-toolbox\pkg\moment\cycle.go
 * @Description: 订阅周期：锚定起始日按月/季/年推进，不累积漂移
 *
 * Copyright (c) 2026 by kamalyes, All Rights Reserved.
 */

package moment

import (
	"errors"
	"fmt"
	"time"
)

// ErrInvalidCycleStep 每期跨度必须为正
//
// [EN] The step per period must be positive
var ErrInvalidCycleStep = errors.New("moment: 每期跨度必须大于 0")

// ErrInvalidCycleUnit 周期单位不支持
//
// [EN] Unsupported cycle unit
var ErrInvalidCycleUnit = errors.New("moment: 周期单位只支持 日/周/月/季/年")

// Cycle 订阅周期，锚定在订阅起始时刻
//
// **为什么需要这个类型，而不是反复调用 Add**
//
// 月末钳位对单次 Add 是对的（1 月 31 日加一个月得到 2 月 28 日），
// 但如果把上一期的结果再加一个月，钳位的结果就成了下一次计算的输入，误差会累积：
//
//	1-31 → 2-28 → 3-28 → 4-28 → 5-28    扣款日永久漂移到 28 号
//
// 正确做法是每期都从**锚点**算，钳位只作用一次：
//
//	1-31 → 2-28 → 3-31 → 4-30 → 5-31    漂移不累积，回到长月仍是 31 号
//
// Cycle 把锚点连同步长一起记住，所有方法内部都从锚点出发，
// 因此不存在"忘记从锚点算"这种误用
//
// [EN] Why this type exists instead of repeated Add calls. Month-end clamping is
// correct for a single Add, but feeding the clamped result back in makes the
// error accumulate: 31 January stepped month by month becomes 28 February and
// then stays on the 28th forever, silently moving the billing day. Computing
// every period from the anchor instead applies the clamp exactly once, so a
// short month does not poison the following ones. Cycle stores the anchor
// alongside the step and every method works from it, so the misuse is not
// expressible
//
// 本类型是值类型，可以自由复制，并发读也是安全的
// 它只有锚点与步长两个字段，不含任何可变状态
//
// [EN] Cycle is a value type: copy it freely, and concurrent reads are safe
// It holds only an anchor and a step, with no mutable state
type Cycle struct {
	anchor time.Time
	step   int
	unit   Unit
}

// NewCycle 创建订阅周期
//
// start 是订阅起始时刻，同时作为锚点；step 是每期跨几个 unit，必须为正
// unit 只接受 UnitDay、UnitWeek、UnitMonth、UnitQuarter、UnitYear
// 更小的单位（秒、分、时）对订阅没有意义，显式拒绝而不是静默接受
//
// [EN] start doubles as the anchor; step must be positive. Only day, week,
// month, quarter and year are accepted — sub-day units are meaningless for
// subscriptions and are rejected explicitly rather than silently honoured
func NewCycle(start Moment, step int, unit Unit) (Cycle, error) {
	if step <= 0 {
		return Cycle{}, ErrInvalidCycleStep
	}
	switch unit {
	case UnitDay, UnitWeek, UnitMonth, UnitQuarter, UnitYear:
	default:
		return Cycle{}, ErrInvalidCycleUnit
	}
	return Cycle{anchor: start.t, step: step, unit: unit}, nil
}

// MustNewCycle 同 NewCycle，参数非法时 panic
//
// 仅用于参数是写死的常量、不可能出错的场景（如 NewCycle(start, 1, UnitMonth)）
//
// [EN] Same as NewCycle but panics on invalid input; only for hard-coded
// arguments that cannot fail
func MustNewCycle(start Moment, step int, unit Unit) Cycle {
	c, err := NewCycle(start, step, unit)
	if err != nil {
		panic(err)
	}
	return c
}

// mustValid 拦住未初始化的零值 Cycle
//
// 零值的 step 是 0，会让 estimateIndex 里的除法除零，
// 报出的是 "integer divide by zero" 这种看不出原因的运行时错误
// 这里换成说清缘由的 panic：零值 Cycle 是编程错误，不该被静默当成某个默认周期，
// 但至少要让人一眼知道错在哪
//
// [EN] Guards against an uninitialised zero value, whose step of 0 makes the
// division in estimateIndex panic with an opaque "integer divide by zero". A
// zero Cycle is a programming error and must not be silently treated as some
// default period, but the panic should at least say why
func (c Cycle) mustValid() {
	if c.step <= 0 {
		panic("moment: Cycle 未初始化，请用 NewCycle 或 MustNewCycle 创建，不要直接声明零值")
	}
}

// Anchor 返回锚点，即订阅起始时刻
func (c Cycle) Anchor() Moment { return Moment{t: c.anchor} }

// Every 返回每期跨度
func (c Cycle) Every() (int, Unit) { return c.step, c.unit }

// NthStart 第 i 期的起点，i 从 0 开始，0 即锚点本身
//
// i 可以为负，表示锚点之前的周期，便于回推上一期
// 无论 i 多大，都是从锚点一次性推进 i*step 个单位，不存在逐期累积
//
// [EN] Start of period i, zero-based, where 0 is the anchor itself. Negative i
// addresses periods before the anchor. However large i is, it steps from the
// anchor once, so nothing accumulates
func (c Cycle) NthStart(i int) Moment {
	c.mustValid()

	// 溢出必须拦住，否则会静默给出荒谬的结果
	//
	// 实测未加检查时：锚点 2026-01-31 的月周期，NthStart(math.MaxInt) 返回 2025-12-28，
	// 比锚点还早；年周期的 NthStart(math.MinInt) 返回锚点本身，因为 12*MinInt 溢出成 0
	// 这两个都不报错，调用方拿到的是一个看起来正常的日期
	// 更糟的是 IndexAt 的校正循环依赖 NthStart 随 i 单调，溢出一旦发生单调性就断了，
	// 循环可能不再收敛
	//
	// [EN] Overflow has to be caught, or the result is silently absurd. Measured
	// without this check, a monthly cycle anchored at 2026-01-31 returned
	// 2025-12-28 for NthStart(math.MaxInt) — earlier than its own anchor — and a
	// yearly cycle returned the anchor itself for math.MinInt because 12*MinInt
	// wraps to 0. Neither reported an error, so the caller received a
	// plausible-looking date. Worse, IndexAt's correction loop relies on
	// NthStart being monotonic in i, and overflow breaks that, so the loop may
	// stop converging
	n := i * c.step
	if n/c.step != i {
		panic(fmt.Sprintf("moment: Cycle 期号 %d 乘以步长 %d 溢出", i, c.step))
	}

	out := Moment{t: c.anchor}.Add(n, c.unit)

	// 自检方向：i 非零时结果必须严格偏向对应一侧
	// Add 内部对季/年还要再乘 3 或 12，那一步也可能溢出，这里兜住
	//
	// 必须用「严格」而不是 Before/After 取反：年周期的 NthStart(math.MinInt) 里
	// 12*MinInt 恰好回绕成 0，结果等于锚点本身，既不早于也不晚于，
	// 只判 Before/After 会漏掉这一例（实测漏过）
	//
	// [EN] Direction self-check: for a non-zero i the result must fall strictly
	// on the matching side. Add multiplies by 3 or 12 again for quarters and
	// years, which can overflow in turn. Strictness matters: for a yearly cycle,
	// 12*math.MinInt wraps to exactly 0, so the result equals the anchor and is
	// neither before nor after it — testing Before/After alone missed this case,
	// as measured
	if (i > 0 && !out.t.After(c.anchor)) || (i < 0 && !out.t.Before(c.anchor)) {
		panic(fmt.Sprintf("moment: Cycle 期号 %d 超出可表示的时间范围（步长 %d）", i, c.step))
	}
	return out
}

// NthEnd 第 i 期的结束时刻，即下一期起点的前一纳秒
//
// 与 Moment.EndOf 一致采用闭区间右端；若下游需要半开区间 [start, next)，
// 请直接用 NthStart(i+1)
//
// [EN] End of period i, the nanosecond before the next period starts, matching
// Moment.EndOf. For a half-open range use NthStart(i+1) instead
func (c Cycle) NthEnd(i int) Moment {
	return Moment{t: c.NthStart(i + 1).t.Add(-time.Nanosecond)}
}

// monthsPerPeriod 每期折合多少个月，仅月族单位有效
//
// **这个函数只服务于估算，正确性不依赖它。** IndexAt 的校正循环才是权威：
// 实测把这里的 3 * c.step 改成 c.step（季度算错 3 倍），
// 包括 -8 到 8 期的季度往返在内的全部测试依然通过，只是校正要多走几步
// 所以改坏它的代价是慢，不是错——反过来说，也别指望测试能发现这里的错
//
// [EN] This feeds the estimate only; correctness does not depend on it, because
// IndexAt's correction loop is authoritative. Making this wrong by a factor of
// three was measured to keep every test passing, including the quarter
// round-trip over periods -8 to 8, at the cost of extra correction steps. So
// breaking it costs speed rather than correctness — and equally, do not expect
// the tests to catch a mistake here
func (c Cycle) monthsPerPeriod() int {
	switch c.unit {
	case UnitMonth:
		return c.step
	case UnitQuarter:
		return 3 * c.step
	case UnitYear:
		return 12 * c.step
	default:
		return 0
	}
}

// IndexAt 返回 t 落在第几期
//
// 期的区间是左闭右开：恰好等于某期起点时归入该期
// t 早于锚点时返回负数
//
// 实现是"先估算再校正"：月族按日历月差估，日/周按时长估，
// 然后用 NthStart 逐步校正到准确值
// 必须校正，因为钳位会让"第 i 期起点"与按月差直接推算的结果不完全对应
//
// [EN] Which period contains t, left-closed and right-open, negative before the
// anchor. Estimated from the calendar-month difference (month family) or elapsed
// duration (day and week), then corrected against NthStart — the correction is
// required because clamping breaks the direct correspondence between a month
// difference and a period index
func (c Cycle) IndexAt(t Moment) int {
	c.mustValid()
	return c.correctIndex(c.estimateIndex(t), t)
}

// correctIndex 把任意估算值 i 校正为准确期号
//
// 抽成独立函数有两个用处：一是把「估算只需接近、校正才是权威」这件事落成代码结构，
// 二是测试可以故意喂偏大或偏小的估算，把两个方向都走到
// 当前两种估算都是真值的**上界**（见 estimateIndex），所以实际只会走到向下校正那支，
// 向上那支在正常路径上不可达——留着是为了估算以后若被改动仍然正确，
// 不留则一旦估算变成下界就会静默算错
//
// [EN] Corrects an arbitrary estimate i into the exact period index. Extracting
// it serves two purposes: it encodes structurally that the estimate only has to
// be close while the correction is authoritative, and it lets tests feed a
// deliberately high or low estimate to exercise both directions. Both current
// estimates are upper bounds on the true index (see estimateIndex), so only the
// downward branch is reachable on the normal path; the upward one is kept so
// that a future change to the estimate stays correct, since dropping it would
// silently break the moment an estimate became a lower bound
func (c Cycle) correctIndex(i int, t Moment) int {
	// 校正：区间左闭右开
	//
	// 两个循环一定会停：step 已由 NewCycle 保证为正，因此 NthStart(i) 随 i 严格单调递增，
	// 而估算值与真值相差不超过一两期，所以每个循环只走几步
	// 单调性是终止的前提——这也是 mustValid 必须拦住零值的另一个原因，
	// step 为 0 时 NthStart 恒等于锚点，下面的循环会无限自减
	//
	// [EN] Both loops terminate: NewCycle guarantees a positive step, so
	// NthStart is strictly increasing in i, and the estimate is off by at most a
	// period or two. That monotonicity is what makes termination hold, which is
	// the other reason mustValid has to reject the zero value: with a step of 0,
	// NthStart would be constant and the loop below would decrement forever
	for c.NthStart(i).t.After(t.t) {
		i--
	}
	for !c.NthStart(i + 1).t.After(t.t) {
		i++
	}
	return i
}

// estimateIndex 估算期号，不保证准确，由 IndexAt 校正
func (c Cycle) estimateIndex(t Moment) int {
	if mpp := c.monthsPerPeriod(); mpp > 0 {
		am := c.anchor.Year()*12 + int(c.anchor.Month())
		tm := t.t.Year()*12 + int(t.t.Month())
		return floorDiv(tm-am, mpp)
	}

	days := c.step
	if c.unit == UnitWeek {
		days = 7 * c.step
	}

	// 用公历日序号相减，而不是 t.Sub(anchor)
	//
	// time.Duration 是 int64 纳秒，上限约 292 年，超出就饱和
	// 原先用 Sub 的写法在相距很远时估算会被钉在 292 年处，校正循环只能一步一步爬：
	// 实测日周期相距 2000 年时单次 IndexAt 要 21.7 毫秒，相距 500 年要 2.0 毫秒，
	// 而零值 time.Time（公元 1 年）当锚点要 14.5 毫秒
	// 未初始化的时间字段在真实数据里很常见，这个开销还完全不显形
	// 换成公历日序号后不存在饱和，估算与真值最多差一天，校正循环恒定只走一两步
	//
	// [EN] Civil day numbers are subtracted instead of calling t.Sub(anchor),
	// because time.Duration is an int64 of nanoseconds and saturates at roughly
	// 292 years. With Sub, a distant t pinned the estimate at that ceiling and
	// the correction loop had to crawl the rest: a daily cycle 2000 years out
	// measured 21.7ms for a single IndexAt, 500 years out 2.0ms, and a
	// zero-value time.Time as the anchor 14.5ms — and uninitialised time fields
	// are common in real data, with none of that cost being visible. Civil day
	// numbers cannot saturate, so the estimate is within a day and the
	// correction is a constant step or two
	ay, am, ad := c.anchor.Date()
	ty, tm, td := t.t.Date()
	return floorDiv(daysFromCivil(ty, int(tm), td)-daysFromCivil(ay, int(am), ad), days)
}

// daysFromCivil 把公历年月日折成相对 1970-01-01 的天序号
//
// 取自 Howard Hinnant 的 chrono 算法：把 3 月当作年首，闰年规则就退化成
// 4/100/400 的整除判断，于是整个换算没有分支也没有查表
// 只用 int 运算，不碰 time.Duration，因此不存在溢出或饱和
//
// [EN] Converts a civil date to a day number relative to 1970-01-01, using
// Howard Hinnant's algorithm: treating March as the first month reduces the leap
// rule to divisibility by 4, 100 and 400, so the conversion needs no branches or
// tables. It stays in int arithmetic and never touches time.Duration, so it
// cannot overflow or saturate
func daysFromCivil(y, m, d int) int {
	if m <= 2 {
		y--
	}
	era := y / 400
	if y < 0 {
		era = (y - 399) / 400
	}
	yoe := y - era*400 // [0, 399]
	mp := m - 3        // 以 3 月为年首
	if m <= 2 {
		mp = m + 9
	}
	doy := (153*mp+2)/5 + d - 1 // [0, 365]
	doe := yoe*365 + yoe/4 - yoe/100 + doy
	return era*146097 + doe - 719468
}

// floorDiv 向下取整的整数除法
//
// Go 的 / 对负数是向零截断（-1/3 == 0），而期号需要向下取整（-1/3 == -1），
// 否则锚点之前的时刻会被归错期
//
// [EN] Floor division. Go truncates towards zero for negative operands, but a
// period index must round down, otherwise instants before the anchor land in the
// wrong period
func floorDiv(a, b int) int {
	q := a / b
	if (a%b != 0) && ((a < 0) != (b < 0)) {
		q--
	}
	return q
}

// PeriodAt 返回 t 所在期的期号与起止时刻
//
// [EN] The index and bounds of the period containing t
func (c Cycle) PeriodAt(t Moment) (index int, start, end Moment) {
	index = c.IndexAt(t)
	return index, c.NthStart(index), c.NthEnd(index)
}

// Next 返回 after 之后的首个期起点，即下一次扣款时刻
//
// 严格晚于 after：若 after 恰好是某期起点，返回的是再下一期
// 这正是扣款语义——今天扣过了，下一次应在下个周期
//
// [EN] The first period start strictly after `after`, i.e. the next charge. If
// `after` is itself a period start the following one is returned, which matches
// billing semantics: today is already charged
func (c Cycle) Next(after Moment) Moment {
	return c.NthStart(c.IndexAt(after) + 1)
}

// Schedule 返回从第 from 期起、连续 count 期的起点
//
// count <= 0 返回空切片
//
// [EN] The starts of count consecutive periods beginning at period `from`;
// a non-positive count yields an empty slice
func (c Cycle) Schedule(from, count int) []Moment {
	if count <= 0 {
		return nil
	}
	out := make([]Moment, count)
	for k := 0; k < count; k++ {
		out[k] = c.NthStart(from + k)
	}
	return out
}
