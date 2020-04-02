package astvalidator

import (
	"strings"
	"time"
)

func (c *Condition) ValidateCondition(condition Condition) (isValid bool, err error) {
	referenceAttrMap := make(map[string]bool)
	inputAttrMap := make(map[string]bool)

	c.readAllAttributes(referenceAttrMap)
	condition.readAllAttributes(inputAttrMap)
	condition.setNonExistAttributeDefaultValue(referenceAttrMap, inputAttrMap)
	return c.validateConditionAttribute(condition)
}

func (c *Condition) readAllAttributes(attrMap map[string]bool) {
	if len(c.Conditions) > 0 {
		for _, condition := range c.Conditions {
			condition.readAllAttributes(attrMap)
		}
	} else {
		if c.Attribute != nil {
			if _, ok := attrMap[c.Attribute.Name]; !ok {
				attrMap[c.Attribute.Name] = true
			}
		}
	}
}

func (c *Condition) validateConditionAttribute(condition Condition) (isValid bool, err error) {
	if len(c.Conditions) > 0 {
		for i, subCondition := range c.Conditions {
			isSubValid, err := subCondition.validateConditionAttribute(condition)
			if err != nil {
				return false, err
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
		isValid, _, err = c.validateConditionValue("", condition)
	}
	return
}

func (c *Condition) validateConditionValue(prefix string, condition Condition) (isValid, isSkip bool, err error) {
	isValid = true
	if len(condition.Conditions) > 0 {
		for i, subCondition := range condition.Conditions {
			isSubValid, isSkip, err := c.validateConditionValue(prefix, *subCondition)
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
		if c.Attribute == nil || condition.Attribute == nil {
			return false, false, nil
		}
		if condition.Attribute.Name == c.Attribute.Name {
			operator := c.Attribute.Operator
			switch operator {
			case OperatorEqual:
				isValid = strings.EqualFold(condition.Attribute.Value, c.Attribute.Value)
			default:
				value := condition.Attribute.Value
				secondValue := c.Attribute.Value
				valueType := getValueType(c.Attribute.Value)

				switch valueType {
				case TypeTime:
					isValid = validateTime(stringToTime(value), operator, stringToTime(secondValue))
				default:
					isValid = validateNumeric(stringToFloat64(value), operator, stringToFloat64(secondValue))
				}
			}
		} else {
			return false, true, nil
		}
	}
	return
}

func (c *Condition) setNonExistAttributeDefaultValue(referenceAttrMap, inputAttrMap map[string]bool) {
	for attrName, _ := range referenceAttrMap {
		if _, ok := inputAttrMap[attrName]; !ok {
			c.Conditions = append(c.Conditions, &Condition{
				Operator: LogicalOperatorAnd,
				Attribute: &Attribute{
					Name:     attrName,
					Operator: "=",
					Value:    "",
				},
			})
		}
	}
}

func getValueType(value string) int {
	var varType, indexVal, dotCount int = TypeAlphanumeric, 0, 0
	for _, char := range value {
		if char == ',' {
			continue
		}
		if '0' <= char && char <= '9' {
			if indexVal == 0 || (indexVal > 0 && dotCount == 1) {
				varType = TypeNumeric
			}
		} else if char == '.' {
			if indexVal > 0 && varType == TypeNumeric {
				dotCount++
				varType = TypeAlphanumeric
			}
			if dotCount > 1 {
				varType = TypeAlphanumeric
				break
			}
		} else {
			varType = TypeAlphanumeric
			break
		}
		indexVal++
	}
	if varType == TypeAlphanumeric {
		if _, err := time.Parse(DateTimeFormat, value); err == nil {
			varType = TypeTime
		}
	}
	return varType
}
