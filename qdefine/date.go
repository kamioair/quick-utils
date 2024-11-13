package qdefine

import (
	"fmt"
	"github.com/kamioair/quick-utils/qconfig"
	"github.com/kamioair/quick-utils/qconvert"
	"strconv"
	"strings"
	"time"
)

var (
	dateFormat = "" // 日期掩码
)

// NewDate
//
//	@Description: 创建日期
//	@param t 时间
//	@return Date
func NewDate(t time.Time) Date {
	t = t.Local()
	s := fmt.Sprintf("%04d%02d%02d", t.Year(), t.Month(), t.Day())
	v, _ := strconv.ParseUint(s, 10, 32)
	return Date(v)
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

// ToString
//
//	@Description: 根据全局format格式化输出
//	@return string
//
//goland:noinspection GoMixedReceiverTypes
func (d Date) ToString() string {
	if dateFormat == "" {
		dateFormat = qconfig.Get("", "com.dateFormat", "yyyy-MM-dd")
	}
	return qconvert.DateTime.ToString(d.ToTime(), "yyyy-MM-dd")
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
	v, err := qconvert.DateTime.ToTime(string(data))
	if err == nil {
		s := fmt.Sprintf("%04d%02d%02d", v.Year(), v.Month(), v.Day())
		t, _ := strconv.ParseUint(s, 10, 64)
		*d = Date(t)
	}
	return err
}
