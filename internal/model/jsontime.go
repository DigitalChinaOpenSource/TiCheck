package model

import (
	"database/sql/driver"
	"errors"
	"fmt"
	"time"
)

const (
	DateFormat = "2006-01-02"
	TimeFormat = "2006-01-02 15:04:05"
)

type JsonTime time.Time

func (t JsonTime) String() string {
	return time.Time(t).Format(TimeFormat)
}

//MarshalJSON 实现它的json序列化方法
func (t JsonTime) MarshalJSON() ([]byte, error) {
	var stamp = fmt.Sprintf("\"%s\"", time.Time(t).Format(TimeFormat))
	return []byte(stamp), nil
}

func (t *JsonTime) UnmarshalJSON(data []byte) (err error) {
	now, err := time.ParseInLocation(`"`+TimeFormat+`"`, string(data), time.Local)
	*t = JsonTime(now)
	return
}

func (t JsonTime) Value() (driver.Value, error) {
	// MyTime 转换成 time.Time 类型
	tTime := time.Time(t)
	return tTime.Format(TimeFormat), nil
}

func (t *JsonTime) Scan(v interface{}) error {
	switch vt := v.(type) {
	case time.Time:
		// 字符串转成 time.Time 类型
		*t = JsonTime(vt)
	default:
		return errors.New("类型处理错误")
	}
	return nil
}
