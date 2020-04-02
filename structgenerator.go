package astvalidator

import (
	"bytes"
)

var (
	mapOperator = map[string]interface{}{
		OperatorEqual:            nil,
		OperatorLessThan:         nil,
		OperatorGreaterThan:      nil,
		OperatorLessThanEqual:    nil,
		OperatorGreaterThanEqual: nil,
	}

	mapLogicalOperator = map[string]string{
		LogicalOperatorAndSyntax: LogicalOperatorAnd,
		LogicalOperatorOrSyntax:  LogicalOperatorOr,
	}
)

func GenerateCondition(query string) (Condition, error) {
	tokenAttributes := getTokenAttributes(query)
	if len(tokenAttributes) == 0 {
		return Condition{Attribute: &Attribute{}}, nil
	}
	_, condition := buildCondition(Condition{}, tokenAttributes)
	return condition, nil
}

func buildCondition(condition Condition, attrs []*TokenAttribute) (int, Condition) {
	var (
		conditionItem *Condition
		lastPos       int
		operator      string
	)
	for i := 0; i < len(attrs); i++ {
		lastPos = i
		attr := attrs[i]
		if attr.hasCalled {
			continue
		}
		attr.hasCalled = true
		if attr.value == ")" {
			break
		}
		if attr.value == "(" {
			newCondition := Condition{
				Operator: operator,
			}
			lastPos, resp := buildCondition(newCondition, attrs[i+1:])
			condition.Conditions = append(condition.Conditions, &resp)
			i = i + lastPos + 1
			continue
		}

		if val, ok := mapLogicalOperator[attr.value]; ok {
			operator = val
			conditionItem = nil
		} else if _, ok := mapOperator[attr.value]; ok {
			conditionItem.Attribute.Operator = attr.value
		} else {
			if conditionItem == nil {
				conditionItem = &Condition{
					Attribute: &Attribute{
						Name: attr.value,
					},
				}
				conditionItem.Attribute = &Attribute{
					Name: attr.value,
				}
			} else {
				conditionItem.Attribute.Value = attr.value
				if condition.Conditions == nil {
					condition.Conditions = []*Condition{}
				}
				conditionItem.Operator = operator
				condition.Conditions = append(condition.Conditions, conditionItem)
			}
		}
	}
	return lastPos, condition
}

func getTokenAttributes(query string) []*TokenAttribute {
	tokenAttributes := []*TokenAttribute{}
	buffer := &bytes.Buffer{}
	isOpenQuote := false
	for _, char := range query {
		switch char {
		case ' ', '\n', '\'':
			if !isOpenQuote {
				continue
			} else {
				buffer.WriteRune(char)
			}
		case '|', '&', '<', '>':
			if buffer.Len() > 0 {
				bufBytes := buffer.Bytes()
				switch bufBytes[0] {
				case ByteVerticalBar:
					tokenAttributes = appendAttribute(tokenAttributes, buffer, LogicalOperatorOrSyntax)
				case ByteAmpersand:
					tokenAttributes = appendAttribute(tokenAttributes, buffer, LogicalOperatorAndSyntax)
				default:
					tokenAttributes = appendAttribute(tokenAttributes, buffer, string(bufBytes))
					buffer.WriteRune(char)
				}
			} else {
				buffer.WriteRune(char)
			}
		case '=', '(', ')':
			if buffer.Len() > 0 {
				bufBytes := buffer.Bytes()
				switch bufBytes[0] {
				case ByteLessThan, ByteGreaterThan:
					tokenAttributes = appendAttribute(tokenAttributes, buffer, string(bufBytes)+string(char))
					continue
				default:
					tokenAttributes = appendAttribute(tokenAttributes, buffer, string(bufBytes))
				}
			}
			tokenAttributes = append(tokenAttributes, &TokenAttribute{
				value: string(char),
			})
		case '"':
			isOpenQuote = !isOpenQuote
		default:
			if buffer.Len() > 0 {
				bufByte := buffer.Bytes()[0]
				if bufByte == ByteLessThan || bufByte == ByteGreaterThan {
					tokenAttributes = appendAttribute(tokenAttributes, buffer, string(bufByte))
				}
			}
			buffer.WriteRune(char)
		}
	}
	if buffer.Len() > 0 {
		tokenAttributes = appendAttribute(tokenAttributes, buffer, buffer.String())
	}
	return tokenAttributes
}

func appendAttribute(tokenAttributes []*TokenAttribute, buffer *bytes.Buffer, value string) []*TokenAttribute {
	tokenAttributes = append(tokenAttributes, &TokenAttribute{
		value: value,
	})
	buffer.Reset()
	return tokenAttributes
}
