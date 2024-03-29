package astvalidator

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"time"
)

func (c *Condition) Validate(data interface{}) (isValid bool, err error) {
	if data == nil {
		return false, fmt.Errorf(ErrorMessageInvalidData, "nil")
	}
	rType := reflect.TypeOf(data)
	switch rType.Kind() {
	case reflect.Struct, reflect.Map:
		isValid, _, err := c.validateAttribute(rType, data)
		return isValid, err
	default:
		return false, fmt.Errorf(ErrorMessageInvalidType, "struct")
	}
}

func (c *Condition) ValidateObjects(data ... interface{}) (isValid bool, err error) {
	if data == nil {
		return false, fmt.Errorf(ErrorMessageInvalidData, "nil")
	}
	rType := reflect.TypeOf(data)
	switch rType.Kind() {
	case reflect.Slice:
		preparedData, err := c.prepareDataFromSlice(data)
		if err != nil {
			return false, err
		}
		return c.Validate(preparedData)
	default:
		return false, fmt.Errorf(ErrorMessageInvalidType, "slice")
	}
}

func (c *Condition) FilterSlice(data interface{}) (result interface{}, err error) {
	if data == nil {
		return result, fmt.Errorf(ErrorMessageInvalidData, "nil")
	}
	rType := reflect.TypeOf(data)
	switch rType.Kind() {
	case reflect.Slice:
		rValue := reflect.ValueOf(data)
		rSlice := reflect.MakeSlice(rType, 0, 1)
		for i := 0; i < rValue.Len(); i++ {
			obj := rValue.Index(i).Interface()
			isValid, err := c.Validate(obj)
			if err != nil {
				return rSlice, err
			}
			if isValid {
				rSlice = reflect.Append(rSlice, rValue.Index(i))
			}
		}
		result = rSlice.Interface()
		return
	default:
		return result, fmt.Errorf(ErrorMessageInvalidType, "slice")
	}
}

func (c *Condition) prepareDataFromSlice(data interface{}) (interface{}, error) {
	var preparedData interface{}
	rValue := reflect.ValueOf(data)
	if rValue.Type().Kind() != reflect.Slice {
		return false, fmt.Errorf(ErrorMessageInvalidType, "slice")
	}
	if rValue.Len() == 0 {
		return false, fmt.Errorf(ErrorMessageInvalidData, "empty slice")
	}

	firstValue := rValue.Index(0).Interface()
	rFirstValue := reflect.ValueOf(firstValue)
	if firstValue == nil {
		return false, fmt.Errorf(ErrorMessageInvalidData, "nil")
	}
	switch rFirstValue.Type().Kind() {
	case reflect.Struct:
		preparedData = firstValue
	default:
		length := rFirstValue.Len()
		switch length {
		case 0:
			return false, fmt.Errorf(ErrorMessageInvalidData, "empty slice")
		case 1:
			preparedData = rFirstValue.Index(0).Interface()
		default:
			mapObj := make(map[string]interface{})
			mapValue := reflect.MakeMap(reflect.TypeOf(mapObj))
			for i := 0; i < length; i++ {
				rDetailValue := reflect.ValueOf(rFirstValue.Index(i).Interface())
				mapValue.SetMapIndex(reflect.ValueOf(rDetailValue.Type().Name()), rDetailValue)
			}
			preparedData = mapValue.Interface()
		}
	}
	return preparedData, nil
}

func (c *Condition) validateAttribute(rType reflect.Type, data interface{}) (isValid, isSkip bool, err error) {
	if len(c.Conditions) > 0 {
		for i, subCondition := range c.Conditions {
			isSubValid, isSkip, err := subCondition.validateAttribute(rType, data)
			if err != nil {
				return false, false, err
			}
			if isSkip {
				continue
			}
			if i == 0 {
				isValid = isSubValid
			} else {
				if subCondition.Operator == LogicalOperatorOr {
					isValid = isValid || isSubValid
				} else {
					isValid = isValid && isSubValid
				}
			}
		}
	} else {
		switch rType.Kind() {
		case reflect.Map:
			if value, ok := data.(map[string]interface{}); ok {
				isValid, isSkip, err = c.validateMapValue(value)
			} else {
				return false, false, errors.New(ErrorMessageUnableToCastObject)
			}
		default:
			isValid, err = c.validateStructValue("", data)
		}
	}
	return
}

func (c *Condition) validateStructValue(prefix string, data interface{}) (isValid bool, err error) {
	rValue := reflect.ValueOf(data)
	if rValue.Type().Kind() != reflect.Struct {
		return false, fmt.Errorf(ErrorMessageInvalidType, "struct")
	}
	for i := 0; i < rValue.NumField(); i++ {
		field := rValue.Field(i)
		typeField := rValue.Type().Field(i)
		tag := typeField.Name
		jsonTag, ok := typeField.Tag.Lookup("json")
		if ok && jsonTag != "" {
			tag = jsonTag
		}
		tag = prefix + tag

		if tag == c.Attribute.Name {
			var conditionValue interface{}
			validationType := TypeNumeric
			value := field.Interface()
			operator := c.Attribute.Operator

			if field.Kind() == reflect.Ptr && field.IsNil() {
				return false, nil
			}

			switch value.(type) {
			case int, int64:
				value = interfaceToInt64(value)
				conditionValue, err = strconv.ParseInt(c.Attribute.Value, 10, 64)
			case *int, *int64:
				value = interfacePtrToInt64(value)
				conditionValue, err = strconv.ParseInt(c.Attribute.Value, 10, 64)
			case float32, float64:
				value = interfaceToFloat64(value)
				conditionValue, err = strconv.ParseFloat(c.Attribute.Value, 64)
			case *float32, *float64:
				value = interfacePtrToFloat64(value)
				conditionValue, err = strconv.ParseInt(c.Attribute.Value, 10, 64)
			case time.Time:
				validationType = TypeTime
				conditionValue, err = time.Parse(DateTimeFormat, c.Attribute.Value)
			case *time.Time:
				validationType = TypeTime
				res, ok := value.(*time.Time)
				if ok {
					value = *res
				}
				conditionValue, err = time.Parse(DateTimeFormat, c.Attribute.Value)
			case bool:
				validationType = TypeAlphanumeric
				conditionValue = stringToBool(c.Attribute.Value)
			default:
				validationType = TypeAlphanumeric
				conditionValue = c.Attribute.Value
			}
			if err != nil {
				return false, err
			}

			switch operator {
			case OperatorEqual:
				isValid = value == conditionValue
			default:
				switch validationType {
				case TypeTime:
					isValid = validateTime(value, operator, conditionValue)
				default:
					isValid = validateNumeric(value, operator, conditionValue)
				}
			}
		}
	}
	return
}

func (c *Condition) validateMapValue(data map[string]interface{}) (isValid, isSkip bool, err error) {
	isSkip = true
	for key, value := range data {
		if len(key) > 0 && !strings.HasPrefix(c.Attribute.Name, key) {
			continue
		}
		isSkip = false
		isValid, err = c.validateStructValue(key+".", value)
		if err != nil {
			return false, false, err
		}
		if !isValid {
			break
		}
	}
	return
}

func validateTime(firstVal interface{}, operator string, secondVal interface{}) bool {
	firstTime, ok := firstVal.(time.Time)
	if !ok {
		return false
	}
	secondTime, ok := secondVal.(time.Time)
	if !ok {
		return false
	}

	switch operator {
	case OperatorGreaterThan:
		return firstTime.After(secondTime)
	case OperatorLessThan:
		return firstTime.Before(secondTime)
	case OperatorGreaterThanEqual:
		return firstTime.After(secondTime) || firstTime.Equal(secondTime)
	default:
		return firstTime.Before(secondTime) || firstTime.Equal(secondTime)
	}
}

func validateNumeric(firstVal interface{}, operator string, secondVal interface{}) bool {
	firstFloat, ok := firstVal.(float64)
	if !ok {
		firstInt, ok := firstVal.(int64)
		if !ok {
			return false
		}
		firstFloat = float64(firstInt)
	}
	secondFloat, ok := secondVal.(float64)
	if !ok {
		secondInt, ok := secondVal.(int64)
		if !ok {
			return false
		}
		secondFloat = float64(secondInt)
	}

	switch operator {
	case OperatorGreaterThan:
		return firstFloat > secondFloat
	case OperatorLessThan:
		return firstFloat < secondFloat
	case OperatorGreaterThanEqual:
		return firstFloat >= secondFloat
	default:
		return firstFloat <= secondFloat
	}
}
