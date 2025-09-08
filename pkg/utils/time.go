package utils

import (
	"database/sql/driver"
	"fmt"
	"strings"
	"time"
)

// ParseTimeRange 解析时间范围
func ParseTimeRange(startTimeStr, endTimeStr string) (*time.Time, *time.Time, error) {
	var startTime, endTime *time.Time
	var err error

	// 支持的时间格式
	timeFormats := []string{
		"2006-01-02 15:04:05",
		"2006-01-02T15:04:05Z",
		"2006-01-02T15:04:05.000Z",
		"2006-01-02",
	}

	// 解析开始时间
	if startTimeStr != "" {
		var parsedTime time.Time
		for _, format := range timeFormats {
			parsedTime, err = time.Parse(format, startTimeStr)
			if err == nil {
				break
			}
		}
		if err != nil {
			return nil, nil, fmt.Errorf("开始时间格式错误: %s", startTimeStr)
		}
		startTime = &parsedTime
	}

	// 解析结束时间
	if endTimeStr != "" {
		var parsedTime time.Time
		for _, format := range timeFormats {
			parsedTime, err = time.Parse(format, endTimeStr)
			if err == nil {
				break
			}
		}
		if err != nil {
			return nil, nil, fmt.Errorf("结束时间格式错误: %s", endTimeStr)
		}

		// 如果只有日期，设置为当天的23:59:59
		if len(endTimeStr) == 10 { // YYYY-MM-DD 格式
			parsedTime = parsedTime.Add(23*time.Hour + 59*time.Minute + 59*time.Second)
		}
		endTime = &parsedTime
	}

	// 验证时间范围
	if startTime != nil && endTime != nil && startTime.After(*endTime) {
		return nil, nil, fmt.Errorf("开始时间不能晚于结束时间")
	}

	return startTime, endTime, nil
}

// FormatTime 格式化时间
func FormatTime(t time.Time) string {
	return t.Format("2006-01-02 15:04:05")
}

// FormatTimePtr 格式化时间指针
func FormatTimePtr(t *time.Time) string {
	if t == nil {
		return ""
	}
	return t.Format("2006-01-02 15:04:05")
}

// DateTime 自定义时间类型，用于统一JSON时间格式
type DateTime struct {
	time.Time
}

// NewDateTime 创建新的DateTime
func NewDateTime(t time.Time) DateTime {
	return DateTime{Time: t}
}

// NewDateTimeNow 创建当前时间的DateTime
func NewDateTimeNow() DateTime {
	return DateTime{Time: time.Now()}
}

// NewDateTimePtr 创建DateTime指针
func NewDateTimePtr(t time.Time) *DateTime {
	if t.IsZero() {
		return nil
	}
	return &DateTime{Time: t}
}

// MarshalJSON 自定义JSON序列化
func (dt DateTime) MarshalJSON() ([]byte, error) {
	if dt.IsZero() {
		return []byte("null"), nil
	}
	formatted := dt.Format("2006-01-02 15:04:05")
	return []byte(fmt.Sprintf(`"%s"`, formatted)), nil
}

// UnmarshalJSON 自定义JSON反序列化
func (dt *DateTime) UnmarshalJSON(data []byte) error {
	str := strings.Trim(string(data), `"`)
	if str == "null" || str == "" {
		dt.Time = time.Time{}
		return nil
	}

	// 支持多种时间格式解析
	timeFormats := []string{
		"2006-01-02 15:04:05",
		"2006-01-02T15:04:05Z",
		"2006-01-02T15:04:05.000Z",
		"2006-01-02T15:04:05+08:00",
		"2006-01-02",
	}

	var err error
	for _, format := range timeFormats {
		dt.Time, err = time.Parse(format, str)
		if err == nil {
			return nil
		}
	}

	return fmt.Errorf("无法解析时间格式: %s", str)
}

// Value 实现 driver.Valuer 接口，用于数据库存储
func (dt DateTime) Value() (driver.Value, error) {
	if dt.IsZero() {
		return nil, nil
	}
	return dt.Time, nil
}

// Scan 实现 sql.Scanner 接口，用于数据库读取
func (dt *DateTime) Scan(value interface{}) error {
	if value == nil {
		dt.Time = time.Time{}
		return nil
	}

	switch v := value.(type) {
	case time.Time:
		dt.Time = v
		return nil
	case string:
		return dt.UnmarshalJSON([]byte(`"` + v + `"`))
	default:
		return fmt.Errorf("无法将 %T 转换为 DateTime", value)
	}
}

// String 返回格式化的时间字符串
func (dt DateTime) String() string {
	if dt.IsZero() {
		return ""
	}
	return dt.Format("2006-01-02 15:04:05")
}

// IsValid 检查时间是否有效
func (dt DateTime) IsValid() bool {
	return !dt.IsZero()
}
