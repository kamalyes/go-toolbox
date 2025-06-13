/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2025-06-09 10:19:55
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2025-06-13 11:52:15
 * @FilePath: \go-toolbox\pkg\schedule\expression.go
 * @Description: 表达式
 *
 * Copyright (c) 2025 by kamalyes, All Rights Reserved.
 */
package schedule

import (
	"errors"
	"fmt"
	"sort"
	"strconv"
	"strings"
	"sync"

	"github.com/kamalyes/go-toolbox/pkg/mathx"
	"github.com/kamalyes/go-toolbox/pkg/random"
	"github.com/kamalyes/go-toolbox/pkg/syncx"
)

// parseFieldRange 将字符串转换为整数，并校验整数是否在指定范围内
// Params：
//
//	s - 待转换字符串
//	fr - 允许的整数范围
//
// Returns：
//
// val - 解析得到的整数值
// err - 解析或校验过程中产生的错误
func parseFieldRange(s string, fr FieldRange) (val int, err error) {
	// 尝试将字符串转换为整数
	if val, err = strconv.Atoi(s); err != nil {
		// 转换失败，返回格式错误
		return 0, fmt.Errorf(ErrValueInvalid, s)
	}

	// 校验整数是否在合法范围内
	if val < fr.Min || val > fr.Max {
		// 超出范围，返回范围错误
		return 0, fmt.Errorf(ErrValueOutOfRange, fr.Min, fr.Max, val)
	}

	// 返回合法的整数值
	return val, nil
}

// parseRangePart 解析单个范围字符串，支持以下格式：
// - 星号 "*"：表示整个范围
// - 范围 "a-b"：表示从a到b的整数区间
// - 单值 "n"：表示单个整数
// Params：
//
//	r - 待解析的范围字符串
//	fr - 允许的整数范围
//
// Returns：
//
// start - 范围起始值
// end - 范围结束值
// err - 解析过程中可能产生的错误
func parseRangePart(r string, fr FieldRange) (start, end int, err error) {
	// 如果是星号，表示整个范围，直接返回字段的最小值和最大值
	if r == StarSymbol {
		return fr.Min, fr.Max, nil
	}

	// 如果包含范围符号（如 "-"），则拆分为起始和结束两部分
	if strings.Contains(r, RangeSymbol) {
		parts := strings.SplitN(r, RangeSymbol, 2)

		// 解析起始值
		if start, err = parseFieldRange(parts[0], fr); err != nil {
			return
		}

		// 解析结束值
		if end, err = parseFieldRange(parts[1], fr); err != nil {
			return
		}

		// 校验起始值不能大于结束值
		if start > end {
			err = fmt.Errorf(ErrRangeStartGreater, r)
		}
		return
	}

	// 单值情况，起始和结束值相同
	start, err = parseFieldRange(r, fr)
	end = start
	return
}

// parseStepPart 解析步进表达式，如 "*/15" 或 "1-10/2"
// Params：
//
//	part - 需要解析的步长范围字符串，格式一般为 "start-end/step"
//	fr - 字段范围，用于辅助解析范围部分的函数或结构体
//
// Returns：
//
//	result - 解析得到的整数切片，包含从 start 到 end，步长为 step 的数字序列
//	err - 解析过程中可能产生的错误
func parseStepPart(part string, fr FieldRange) (result []int, err error) {
	var (
		bounds           = strings.SplitN(part, StepSymbol, 2)        // 按步长符号（"/"）拆分字符串，最多拆成两部分
		startRaw         = mathx.SafeGetIndexOrDefault(bounds, 0, "") // 获取拆分后第一部分，表示起始范围字符串
		endRaw           = mathx.SafeGetIndexOrDefault(bounds, 1, "") // 获取拆分后第二部分，表示步长字符串
		start, end, step int                                          // 定义起始、结束、步长变量
	)
	if len(bounds) != 2 {
		// 如果拆分结果不是两部分，说明格式不正确，返回格式错误
		return nil, fmt.Errorf(ErrStepFormat, part)
	}

	// 解析起始范围字符串，得到起始和结束整数
	if start, end, err = parseRangePart(startRaw, fr); err != nil {
		return nil, err
	}

	// 解析步长字符串，转换为整数且必须大于0，否则返回步长无效错误
	if step, err = strconv.Atoi(endRaw); err != nil || step <= 0 {
		return nil, fmt.Errorf(ErrStepValueInvalid, endRaw)
	}

	// 起始不能大于结束，否则返回起始大于结束错误
	if start > end {
		return nil, fmt.Errorf(ErrRangeStartGreater, part)
	}

	// 根据起始、结束和步长，生成数字序列
	result = random.RandNumerical(start, end, step)
	return result, nil
}

// parseSimplePart 解析不含逗号的单个字段，支持星号、范围、步进和单值
// Params：
//
//	part - 字段字符串
//	fr - 允许的整数范围
//
// Returns：
//
//	[]int - 解析得到的整数切片
//	error - 解析过程中可能产生的错误
func parseSimplePart(part string, fr FieldRange) ([]int, error) {
	// 如果字符串中包含步长符号（如 "/"），则调用解析步长范围的函数
	if strings.Contains(part, StepSymbol) {
		return parseStepPart(part, fr)
	}

	// 否则，直接解析普通范围或单个数字
	start, end, err := parseRangePart(part, fr)
	if err != nil {
		return nil, err
	}

	// 生成从 start 到 end，步长为 1 的数字序列
	vals := random.RandNumerical(start, end, 1)
	return vals, nil
}

// parseCommaSeparated 解析逗号分隔的字段，支持多值组合
// Params：
//
//	field - 字段字符串，可能包含多个用逗号分隔的部分
//	fr - 允许的整数范围
//
// Returns：
//
//	合并去重排序后的整数切片，或错误
func parseCommaSeparated(field string, fr FieldRange) ([]int, error) {
	// 按逗号分割字符串，得到多个部分
	parts := strings.Split(field, CommaSymbol)

	var allVals []int
	// 逐个解析每个部分
	for _, p := range parts {
		vals, err := parseSimplePart(p, fr)
		if err != nil {
			// 解析出错则返回错误
			return nil, err
		}
		// 将解析结果追加到总结果切片中
		allVals = append(allVals, vals...)
	}

	// 对所有结果进行去重，防止重复数字
	return mathx.SliceUniq(allVals), nil
}

// ------------------- 星期字段结构体及解析 -------------------

// WeekdayField 表示星期字段，支持普通值、多值、#扩展语法和问号(?)
type WeekdayField struct {
	vals       []int       // 匹配的星期数字，0-6，0表示周日
	nthWeekday map[int]int // #扩展语法映射，key为星期数字，value为第几个
	isQuestion bool        // 是否为问号(?)，表示不限制星期
}

// Vals 返回匹配的星期数字列表
func (wf *WeekdayField) Vals() []int {
	return wf.vals
}

// NthWeekday 返回 # 语法映射
func (wf *WeekdayField) NthWeekday() map[int]int {
	return wf.nthWeekday
}

// IsQuestion 判断是否为问号(?)，表示不限制星期
func (wf *WeekdayField) IsQuestion() bool {
	return wf.isQuestion
}

// normalizeWeekday 将输入的星期数字根据配置转换为0-6范围（0=周日）
// Params：
//
//	wd - 输入的星期数字
//	numbering - 星期数字体系（0-6或1-7）
//
// Returns：
//
//	标准化后的星期数字，或错误
func normalizeWeekday(wd int, numbering WeekdayNumbering) (int, error) {
	switch numbering {
	case WeekdayZeroBased:
		if wd < WeekdayZeroBasedMin || wd > WeekdayZeroBasedMax {
			return 0, fmt.Errorf(ErrWeekdayValueRangeZero, wd)
		}
		return wd, nil
	case WeekdayOneBased:
		if wd < WeekdayOneBasedMin || wd > WeekdayOneBasedMax {
			return 0, fmt.Errorf(ErrWeekdayValueRangeOne, wd)
		}
		// 特殊转换 1->0(周日), 7->6(周六), 其他数字减1
		return wd - WeekdayOneBasedMin, nil
	default:
		return 0, errors.New(ErrWeekdayUnknownNumbering)
	}
}

// parseAndNormalizeWeekday 将字符串类型的星期数字解析成int，
// 并调用normalizeWeekday进行标准化处理。
func parseAndNormalizeWeekday(wdStr string, numbering WeekdayNumbering) (int, error) {
	wd, err := strconv.Atoi(wdStr)
	if err != nil {
		return 0, fmt.Errorf(ErrValueInvalid, err)
	}
	return normalizeWeekday(wd, numbering)
}

// parseWeekdayField 解析星期字段，支持逗号分隔、多值，#扩展符号，?符号
// Params：
//
//	field - 星期字段字符串
//	numbering - 星期数字体系配置
//
// Returns：
//
//	解析后的WeekdayField结构体，或错误
func parseWeekdayField(field string, numbering WeekdayNumbering) (WeekdayField, error) {
	wf := WeekdayField{nthWeekday: make(map[int]int)}
	if field == QuestionSymbol {
		wf.isQuestion = true
		return wf, nil
	}

	parts := strings.Split(field, CommaSymbol)
	valSet := make(map[int]struct{})

	for _, part := range parts {
		part = strings.TrimSpace(part)

		// 处理 # 语法，如 "2#3" 表示第三个星期一
		if strings.Contains(part, HashSymbol) {
			sub := strings.Split(part, HashSymbol)
			if len(sub) != 2 {
				return wf, fmt.Errorf(ErrWeekdayHashFormat, part)
			}
			nth, err := strconv.Atoi(sub[1])
			if err != nil {
				return wf, fmt.Errorf(ErrWeekdayHashValue, part)
			}
			wd, err := parseAndNormalizeWeekday(sub[0], numbering)
			if err != nil {
				return wf, err
			}
			wf.nthWeekday[wd] = nth
			valSet[wd] = struct{}{}
			continue
		}

		// 处理星号 "*"，表示所有星期
		if part == StarSymbol {
			for i := WeekdayZeroBasedMin; i <= WeekdayZeroBasedMax; i++ {
				valSet[i] = struct{}{}
			}
			continue
		}

		// 处理区间表达式 "a-b"
		if strings.Contains(part, RangeSymbol) {
			bounds := strings.SplitN(part, RangeSymbol, 2)
			if len(bounds) != 2 {
				return wf, fmt.Errorf(ErrRangeFormat, part)
			}
			start, err := parseAndNormalizeWeekday(bounds[0], numbering)
			if err != nil {
				return wf, err
			}
			end, err := parseAndNormalizeWeekday(bounds[1], numbering)
			if err != nil {
				return wf, err
			}
			if start > end {
				return wf, fmt.Errorf(ErrRangeStartGreater, part)
			}
			for i := start; i <= end; i++ {
				valSet[i] = struct{}{}
			}
			continue
		}

		// 处理单个星期数字
		wd, err := parseAndNormalizeWeekday(part, numbering)
		if err != nil {
			return wf, err
		}
		valSet[wd] = struct{}{}
	}

	// 转换成切片并排序
	for v := range valSet {
		wf.vals = append(wf.vals, v)
	}
	sort.Ints(wf.vals)

	return wf, nil
}

// Contains 判断给定星期数字是否匹配当前星期字段
// 参数 wd - 星期数字，0-6
// 返回 true 表示匹配，false 表示不匹配
func (wf *WeekdayField) Contains(wd int) bool {
	// ? 表示不限制，匹配所有
	return mathx.IF(wf.isQuestion, true, mathx.SliceContains(wf.vals, wd))
}

// ------------------- 日字段结构体及解析 -------------------

// DayField 表示日字段，支持普通值、?符号、L和W扩展符号
type DayField struct {
	vals       []int // 普通匹配的日期值
	isQuestion bool  // 是否为问号(?)，表示不限制日期
	isLastDay  bool  // 是否为L，表示当月最后一天
	nearestWD  bool  // 是否为W，表示最近的工作日
}

// IsQuestion 判断是否为问号(?)
func (df *DayField) IsQuestion() bool {
	return df.isQuestion
}

// IsLastDay 判断是否为L，表示当月最后一天
func (df *DayField) IsLastDay() bool {
	return df.isLastDay
}

// IsNearestWD 判断是否为W，表示最近的工作日
func (df *DayField) IsNearestWD() bool {
	return df.nearestWD
}

// Vals 返回普通匹配的日期值
func (df *DayField) Vals() []int {
	return df.vals
}

// parseDayField 解析日字段，支持 ?、L、W 和多值、范围、步进
// Params：
//
//	field - 日字段字符串
//
// Returns：
//
//	解析后的DayField结构体，或错误
func parseDayField(field string) (DayField, error) {
	df := DayField{}
	// 处理 ? 表示不限制

	if field == QuestionSymbol {
		df.isQuestion = true
		return df, nil
	}

	// 处理 L 表示当月最后一天
	if strings.Contains(field, LSymbol) {
		if field == LSymbol {
			df.isLastDay = true
			return df, nil
		}
		// 复杂L格式暂不支持
		return df, fmt.Errorf(ErrDayLComplexFormat, field)
	}

	// 处理 W 表示最近工作日，只支持类似 15W 格式
	if strings.Contains(field, WSymbol) {
		if strings.HasSuffix(field, WSymbol) {
			dayPart := strings.TrimSuffix(field, WSymbol)
			day, err := strconv.Atoi(dayPart)
			if err != nil {
				return df, fmt.Errorf(ErrDayWFormat, field)
			}
			if day < dayRange.Min || day > dayRange.Max {
				return df, fmt.Errorf(ErrDayWOutOfRange, day)
			}
			df.nearestWD = true
			df.vals = []int{day}
			return df, nil
		}
		return df, fmt.Errorf(ErrDayWFormat, field)
	}

	// 普通多值、范围、步进解析
	vals, err := parseCommaSeparated(field, dayRange)
	if err != nil {
		return df, err
	}
	df.vals = vals

	return df, nil
}

// Contains 判断给定日期是否匹配日字段
// 注意：L和W的具体计算需要结合年月，当前未实现，暂返回false
func (df *DayField) Contains(day int) bool {
	if df.isQuestion {
		return true
	}
	if df.isLastDay || df.nearestWD {
		return false
	}
	return mathx.SliceContains(df.vals, day)
}

// ------------------- 简单字段结构体及解析 -------------------

// SimpleField 用于表示秒、分、时、月、年字段的匹配结果
type SimpleField struct {
	vals []int // 匹配的整数集合，例如秒字段可能是 [0,15,30,55]
}

// Vals 返回匹配的整数切片
func (sf *SimpleField) Vals() []int {
	return sf.vals
}

// parseSimpleField 解析秒、分、时、月、年字段，支持多值、范围、步进
// Params：
//
//	field - 字段字符串
//	fr - 允许的整数范围
//
// Returns：
//
//	解析后的SimpleField结构体，或错误
func parseSimpleField(field string, fr FieldRange) (SimpleField, error) {
	vals, err := parseCommaSeparated(field, fr)
	return mathx.ReturnIfErr(SimpleField{vals: vals}, err)
}

// Contains 判断给定整数是否匹配
func (sf *SimpleField) Contains(v int) bool {
	return mathx.SliceContains(sf.vals, v)
}

// ------------------- Cron表达式结构体及解析 -------------------

// CronExpr 表示完整的7字段Cron表达式解析结果
type CronExpr struct {
	Second  SimpleField  // 秒字段，0-59
	Minute  SimpleField  // 分字段，0-59
	Hour    SimpleField  // 时字段，0-23
	Day     DayField     // 日字段，支持扩展符号
	Month   SimpleField  // 月字段，1-12
	Weekday WeekdayField // 星期字段，0-6或1-7，支持扩展符号
	Year    SimpleField  // 年字段，1970~(nowYear+99)
}

// CronParser 负责解析Cron表达式
type CronParser struct {
	expr          string           // 原始表达式字符串
	weekNumbering WeekdayNumbering // 星期数字体系配置
	mu            sync.RWMutex     // 读写锁
}

// NewCronParser 创建CronParser实例，默认星期数字体系为0-6，支持1-7方式
func NewCronParser() *CronParser {
	return &CronParser{
		weekNumbering: WeekdayZeroBased,
	}
}

// SetExpression 设置待解析的Cron表达式字符串，支持链式调用
func (c *CronParser) SetExpression(expr string) *CronParser {
	return syncx.WithLockReturnValue(&c.mu, func() *CronParser {
		c.expr = expr
		return c
	})
}

// SetWeekNumbering 设置星期数字体系，支持链式调用
func (c *CronParser) SetWeekNumbering(num WeekdayNumbering) *CronParser {
	return syncx.WithLockReturnValue(&c.mu, func() *CronParser {
		c.weekNumbering = num
		return c
	})
}

// Parse 解析Cron表达式，返回解析结果或错误
// 支持7字段格式（秒 分 时 日 月 星期 年）
// 校验字段数量和字段合法性，支持各种扩展符号
func (c *CronParser) Parse() (cronExpr *CronExpr, err error) {
	return syncx.WithRLockReturn(&c.mu, func() (*CronExpr, error) {

		// 将表达式按空白字符拆分成字段数组
		fields := strings.Fields(c.expr)

		// 校验字段数量是否在允许范围内（cronRange.Min到cronRange.Max）
		if len(fields) < cronRange.Min || len(fields) > cronRange.Max {
			return nil, errors.New(ErrFieldCount) // 字段数量错误，返回错误
		}

		// 先将各字段赋值给局部变量，方便后续使用和阅读
		var (
			hour, minute SimpleField                                                 // 时分匹配值
			year, month  SimpleField                                                 // 年月匹配值
			second       SimpleField                                                 // 秒匹配值
			weekday      WeekdayField                                                // 星期匹配值
			day          DayField                                                    // 日字段匹配值
			secondField  = mathx.SafeGetIndexOrDefaultNoSpace(fields, 0, StarSymbol) // 秒字段字符串
			minuteField  = mathx.SafeGetIndexOrDefaultNoSpace(fields, 1, StarSymbol) // 分字段字符串
			hourField    = mathx.SafeGetIndexOrDefaultNoSpace(fields, 2, StarSymbol) // 时字段字符串
			dayField     = mathx.SafeGetIndexOrDefaultNoSpace(fields, 3, StarSymbol) // 日字段字符串
			monthField   = mathx.SafeGetIndexOrDefaultNoSpace(fields, 4, StarSymbol) // 月字段字符串
			weekdayField = mathx.SafeGetIndexOrDefaultNoSpace(fields, 5, StarSymbol) // 星期字段字符串
			yearField    = mathx.SafeGetIndexOrDefaultNoSpace(fields, 6, StarSymbol) // 年字段字符串
		)

		// 解析秒字段，使用秒的合法范围secondRange
		if second, err = parseSimpleField(secondField, secondRange); err != nil {
			// 返回带上下文的错误信息，方便定位秒字段解析错误
			return nil, fmt.Errorf("%s, %v", ErrParseSecond, err)
		}

		// 解析分字段，使用分的合法范围minuteRange
		if minute, err = parseSimpleField(minuteField, minuteRange); err != nil {
			return nil, fmt.Errorf("%s, %v", ErrParseMinute, err)
		}

		// 解析时字段，使用时的合法范围hourRange
		if hour, err = parseSimpleField(hourField, hourRange); err != nil {
			return nil, fmt.Errorf("%s, %v", ErrParseHour, err)
		}

		// 解析日字段，日字段可能包含特殊符号，需要专门解析函数
		if day, err = parseDayField(dayField); err != nil {
			return nil, fmt.Errorf("%s, %v", ErrParseDay, err)
		}

		// 解析月字段，使用月的合法范围monthRange
		if month, err = parseSimpleField(monthField, monthRange); err != nil {
			return nil, fmt.Errorf("%s, %v", ErrParseMonth, err)
		}

		// 解析星期字段，传入weekNuming控制数字范围
		if weekday, err = parseWeekdayField(weekdayField, c.weekNumbering); err != nil {
			return nil, fmt.Errorf("%s, %v", ErrParseWeekday, err)
		}

		// 解析年字段，使用年的合法范围yearRange
		if year, err = parseSimpleField(yearField, yearRange); err != nil {
			return nil, fmt.Errorf("%s, %v", ErrParseYear, err)
		}
		// 校验“日”和“星期”字段的互斥规则：
		// 如果日字段和星期字段都为问号，则表达式无效，返回错误
		if day.isQuestion && weekday.isQuestion {
			return nil, errors.New(ErrWeekdayQuestionMutual)
		}

		// 返回完整解析结果
		cronExpr = &CronExpr{
			Second:  second,
			Minute:  minute,
			Hour:    hour,
			Day:     day,
			Month:   month,
			Weekday: weekday,
			Year:    year,
		}
		return cronExpr, err
	})
}
