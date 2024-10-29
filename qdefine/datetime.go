package qdefine

import (
	"fmt"
	"github.com/liaozhibinair/quick-utils/qconvert"
	"strconv"
	"strings"
	"time"
)

var (
	dateFormat     = "yyyy-MM-dd"          // 日期掩码
	dateTimeFormat = "yyyy-MM-dd HH:mm:ss" // 日期时间掩码
)

type Date uint32

// FromTime
//
//	@Description: 通过原生的time赋值
//	@param time
//
//goland:noinspection GoMixedReceiverTypes
func (d *Date) FromTime(time time.Time) {
	t := time.Local()
	s := fmt.Sprintf("%04d%02d%02d", t.Year(), t.Month(), t.Day())
	v, _ := strconv.ParseUint(s, 10, 32)
	*d = Date(v)
}

// FromTime
//
//	@Description: 通过原生的time赋值
//	@param time
func FromTime(time time.Time) (d Date) {
	t := time.Local()
	s := fmt.Sprintf("%04d%02d%02d", t.Year(), t.Month(), t.Day())
	v, _ := strconv.ParseUint(s, 10, 32)
	return Date(v)
}

// ToString
//
//	@Description: 根据全局format格式化输出
//	@return string
//
//goland:noinspection GoMixedReceiverTypes
func (d Date) ToString() string {
	return qconvert.TimeToString(d.ToTime(), dateFormat)
}

// ToTime
//
//	@Description: 转为原生时间对象
//	@return time.Time
//
//goland:noinspection GoMixedReceiverTypes
func (d Date) ToTime() time.Time {
	if d == 0 {
		return time.Time{}
	}
	str := fmt.Sprintf("%d", d)
	if len(str) != 8 {
		str = str + strings.Repeat("0", 8-len(str))
	}
	year, _ := strconv.Atoi(str[0:4])
	month, _ := strconv.Atoi(str[4:6])
	day, _ := strconv.Atoi(str[6:8])
	return time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.Local)
}

// MarshalJSON
//
//	@Description: 复写json转换
//	@return []byte
//	@return error
//
//goland:noinspection GoMixedReceiverTypes
func (d Date) MarshalJSON() ([]byte, error) {
	str := fmt.Sprintf("\"%s\"", d.ToString())
	return []byte(str), nil
}

// UnmarshalJSON
//
//	@Description: 复写json转换
//	@param data
//	@return error
//
//goland:noinspection GoMixedReceiverTypes
func (d *Date) UnmarshalJSON(data []byte) error {
	v, err := qconvert.StringToTime(string(data))
	if err == nil {
		s := fmt.Sprintf("%04d%02d%02d", v.Year(), v.Month(), v.Day())
		t, _ := strconv.ParseUint(s, 10, 64)
		*d = Date(t)
	}
	return err
}

type DateTime uint64

// NowTime
//
//	@Description: 当前系统时间
//	@return DateTime
func NowTime() DateTime {
	var now DateTime
	now.FromTime(time.Now().Local())
	return now
}

// FromTime
//
//	@Description: 通过原生的time赋值
//	@param time
//
//goland:noinspection GoMixedReceiverTypes
func (d *DateTime) FromTime(time time.Time) {
	t := time.Local()
	s := fmt.Sprintf("%04d%02d%02d%02d%02d%02d", t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), t.Second())
	v, _ := strconv.ParseUint(s, 10, 64)
	*d = DateTime(v)
}

// DayFirst
//
//	@Description: 返回日期的首秒，即 年月日000000
//	@return DateTime
func (d DateTime) DayFirstSecond() DateTime {
	if d == 0 {
		return 0
	}
	str := fmt.Sprintf("%d", d)
	if len(str) != 14 {
		str = str + strings.Repeat("0", 14-len(str))
	}
	year, _ := strconv.Atoi(str[0:4])
	month, _ := strconv.Atoi(str[4:6])
	day, _ := strconv.Atoi(str[6:8])
	full, _ := strconv.ParseUint(fmt.Sprintf("%02d%02d%02d000000", year, month, day), 10, 64)
	return DateTime(full)
}

// DayLastSecond
//
//	@Description: 返回日期的最后一秒，即 年月日235959
//	@return DateTime
func (d DateTime) DayLastSecond() DateTime {
	if d == 0 {
		return 0
	}
	str := fmt.Sprintf("%d", d)
	if len(str) != 14 {
		str = str + strings.Repeat("0", 14-len(str))
	}
	year, _ := strconv.Atoi(str[0:4])
	month, _ := strconv.Atoi(str[4:6])
	day, _ := strconv.Atoi(str[6:8])
	full, _ := strconv.ParseUint(fmt.Sprintf("%02d%02d%02d235959", year, month, day), 10, 64)
	return DateTime(full)
}

// Add
//
//	@Description: 添加时间
//	@param duration
//	@return DateTime
func (d DateTime) Add(duration time.Duration) DateTime {
	t := d.ToTime()
	t = t.Add(duration)
	var dt DateTime
	dt.FromTime(t)
	return dt
}

// Date
//
//	@Description: 转为日期
//	@return Date
//
//goland:noinspection GoMixedReceiverTypes
func (d DateTime) Date() Date {
	var date Date
	date.FromTime(d.ToTime())
	return date
}

// ToString
//
//	@Description: 根据全局format格式化输出
//	@return string
//
//goland:noinspection GoMixedReceiverTypes
func (d DateTime) ToString() string {
	return qconvert.TimeToString(d.ToTime(), dateTimeFormat)
}

// ToTime
//
//	@Description: 转为原生时间对象
//	@return time.Time
//
//goland:noinspection GoMixedReceiverTypes
func (d DateTime) ToTime() time.Time {
	if d == 0 {
		return time.Time{}
	}
	str := fmt.Sprintf("%d", d)
	if len(str) != 14 {
		str = str + strings.Repeat("0", 14-len(str))
	}
	year, _ := strconv.Atoi(str[0:4])
	month, _ := strconv.Atoi(str[4:6])
	day, _ := strconv.Atoi(str[6:8])
	hour, _ := strconv.Atoi(str[8:10])
	minute, _ := strconv.Atoi(str[10:12])
	second, _ := strconv.Atoi(str[12:14])
	return time.Date(year, time.Month(month), day, hour, minute, second, 0, time.Local)
}

// MarshalJSON
//
//	@Description: 复写json转换
//	@return []byte
//	@return error
func (d DateTime) MarshalJSON() ([]byte, error) {
	str := fmt.Sprintf("\"%s\"", d.ToString())
	return []byte(str), nil
}

// UnmarshalJSON
//
//	@Description: 复写json转换
//	@param data
//	@return error
func (d *DateTime) UnmarshalJSON(data []byte) error {
	v, err := qconvert.StringToTime(string(data))
	if err == nil {
		s := fmt.Sprintf("%04d%02d%02d%02d%02d%02d", v.Year(), v.Month(), v.Day(), v.Hour(), v.Minute(), v.Second())
		t, _ := strconv.ParseUint(s, 10, 64)
		*d = DateTime(t)
	}
	return err
}
