# AST Validator

> AST Validator is a library to validate object or parameter based on query string

## Function
> **GenerateCondition :**
 generate condition object based on query string
 ```
func GenerateCondition(query string) (Condition, error) {...}
  ```
> **Validate :**
 validate object or parameter using generated condition
 ```
func (c *Condition) Validate(data interface{}) (isValid bool, err error) {...}
```
> **ValidateObjects :**
 validate multi objects or parameters using generated condition
 ```
func (c *Condition) ValidateObjects(data ... interface{}) (isValid bool, err error) {...}
```
> **ValidateCondition :**
 validate custom condition using generated condition
 ```
func (c *Condition) ValidateCondition(condition Condition) (isValid bool, err error) {...}
```

## Support
#### Operator
> Equal

> Less than

> Less than equal

> Greater than

> Greater than equal

#### Value Type
> Numeric

> Alphanumeric

> Time

---

## Sample
```cgo
func validateSample1() bool{
	object := struct {
		ID       string `json:"id"`
		MemberID string `json:"member_id"`
		Division string `json:"division"`
	}{
		ID:       "1",
		MemberID: "2",
		Division: "finance",
	}

	query := "(id=1 && (member_id=12||member_id=2)) && (division=engineering || division=finance)"
	condition, _ := GenerateCondition(query)
	isValid, err := condition.Validate(object)
	if err != nil{
	    return false
	}
	return isValid
}
```

```cgo
func validateSample2() bool {
	fInt := func(i int) *int {
		return &i
	}
	fInt64 := func(i int64) *int64 {
		return &i
	}
	fFloat := func(f float32) *float32 {
		return &f
	}
	fFloat64 := func(f float64) *float64 {
		return &f
	}

	object := struct {
		ID        int        `json:"id"`
		MemberID  int        `json:"member_id"`
		Division  string     `json:"division"`
		Score     *int       `json:"score"`
		Point     *int64     `json:"point"`
		Wallet    *float32   `json:"wallet"`
		Money     *float64   `json:"money"`
		JoinDate  time.Time  `json:"join_date"`
		LeaveDate *time.Time `json:"leave_date"`
	}{
		ID:        5,
		MemberID:  25,
		Division:  "engineering",
		Score:     fInt(100),
		Point:     fInt64(3000),
		Wallet:    fFloat(100),
		Money:     fFloat64(1500000),
		JoinDate:  time.Date(2015, 10, 9, 0, 0, 0, 0, time.UTC),
		LeaveDate: nil,
	}

	query := `join_date>"2015-01-01 00:00:00" && join_date<="2016-01-01 00:00:00" && score>80 && point<4000 && wallet=100`
	condition, _ := GenerateCondition(query)
	isValid, err := condition.Validate(object)
	if err != nil{
        return false
    }
    return isValid
}
```
```cgo
func validateSample3() {
	referenceQuery := `id=1 && ( division = engineering || division = finance )`
	input := `id=1 && division = engineering`
	referenceCondition, _ := GenerateCondition(referenceQuery)
	inputCondition, _ := GenerateCondition(input)

	isValid, err := referenceCondition.ValidateCondition(inputCondition)	
	if err != nil{
        return false
    }
    return isValid
}
```

## Benchmark
```
BenchmarkGenerateCondition-12             572319              2144 ns/op
BenchmarkValidate-12                      631578              1913 ns/op
BenchmarkValidateObjects-12               570372              2079 ns/op
BenchmarkValidateCondition-12            4878042               243 ns/op
```

## Future Development
> Improve case coverages

> Nested object validation