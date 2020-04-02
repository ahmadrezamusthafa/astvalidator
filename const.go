package astvalidator

const (
	LogicalOperatorAnd = "AND"
	LogicalOperatorOr  = "OR"

	LogicalOperatorAndSyntax = "&&"
	LogicalOperatorOrSyntax  = "||"
)

const (
	OperatorEqual            = "="
	OperatorLessThan         = "<"
	OperatorLessThanEqual    = "<="
	OperatorGreaterThan      = ">"
	OperatorGreaterThanEqual = ">="
)

const (
	ByteAmpersand   = 38
	ByteLessThan    = 60
	ByteGreaterThan = 62
	ByteVerticalBar = 124
)

const (
	TypeTime         = 1
	TypeNumeric      = 2
	TypeAlphanumeric = 3
)

const DateTimeFormat = "2006-01-02 15:04:05"

const (
	ErrorMessageInvalidData        = "data can't be %s"
	ErrorMessageInvalidParameter   = "invalid parameter, %s is required"
	ErrorMessageInvalidType        = "invalid type, %s is required"
	ErrorMessageUnableToCastObject = "unable to cast object"
)
