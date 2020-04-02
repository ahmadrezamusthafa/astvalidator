package astvalidator

import (
	"strconv"
	"time"
)

func stringToBool(value string) bool {
	switch value {
	case "t", "true":
		return true
	default:
		return false
	}
}

func stringToFloat64(value string) float64 {
	var floatValue float64
	floatValue, err := strconv.ParseFloat(value, 64)
	if err != nil {
		return 0
	}
	return floatValue
}

func stringToTime(value string) time.Time {
	var timeValue time.Time
	timeValue, err := time.Parse(DateTimeFormat, value)
	if err != nil {
		return time.Time{}
	}
	return timeValue
}

func interfacePtrToInt64(input interface{}) int64 {
	if val, ok := input.(*int64); ok {
		return *val
	}
	res := interfacePtrToInt(input)
	return int64(res)
}

func interfacePtrToInt(input interface{}) int {
	if val, ok := input.(*int); ok {
		return *val
	}
	return 0
}

func interfaceToInt64(input interface{}) int64 {
	if val, ok := input.(int64); ok {
		return val
	}
	res := interfaceToInt(input)
	return int64(res)
}

func interfaceToInt(input interface{}) int {
	if val, ok := input.(int); ok {
		return val
	}
	return 0
}

func interfacePtrToFloat64(input interface{}) float64 {
	if val, ok := input.(*float64); ok {
		return *val
	}
	res := interfacePtrToFloat32(input)
	return float64(res)
}

func interfacePtrToFloat32(input interface{}) float32 {
	if val, ok := input.(*float32); ok {
		return *val
	}
	return 0
}

func interfaceToFloat64(input interface{}) float64 {
	if val, ok := input.(float64); ok {
		return val
	}
	res := interfaceToFloat32(input)
	return float64(res)
}

func interfaceToFloat32(input interface{}) float32 {
	if val, ok := input.(float32); ok {
		return val
	}
	return 0
}
