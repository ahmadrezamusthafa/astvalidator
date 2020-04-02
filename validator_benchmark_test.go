package astvalidator

import (
	"testing"
	"time"
)

//BENCHMARK Validate
//Improvement history:
//------------------------------------
//	attempt	   |  time per loop
//------------------------------------
//  484363	      2316 ns/op (now)
//------------------------------------
func BenchmarkValidate(b *testing.B) {
	object := struct {
		ID       string `json:"id"`
		MemberID string `json:"member_id"`
		Division string `json:"division"`
	}{
		ID:       "1",
		MemberID: "2",
		Division: "finance",
	}

	query := "(id=1 && (member_id=12||member_id=2))  &&   (division=engineering || division=finance)"
	condition, _ := GenerateCondition(query)
	for n := 0; n < b.N; n++ {
		condition.Validate(object)
	}
}

//BENCHMARK ValidateObjects
//Improvement history:
//------------------------------------
//	attempt	   |  time per loop
//------------------------------------
//  359617	      2829 ns/op (now)
//------------------------------------
func BenchmarkValidateObjects(b *testing.B) {
	object := struct {
		ID       string `json:"id"`
		MemberID string `json:"member_id"`
		Division string `json:"division"`
	}{
		ID:       "1",
		MemberID: "2",
		Division: "finance",
	}

	query := "(id=1 && (member_id=12||member_id=2))  &&   (division=engineering || division=finance)"
	condition, _ := GenerateCondition(query)
	for n := 0; n < b.N; n++ {
		condition.ValidateObjects(object)
	}
}

//BENCHMARK ValidateComplexOperator
//Improvement history:
//------------------------------------
//	attempt	   |  time per loop
//------------------------------------
//  172010	      7053 ns/op (now)
//------------------------------------
func BenchmarkValidateComplexOperator(b *testing.B) {
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
	for n := 0; n < b.N; n++ {
		condition.Validate(object)
	}
}
