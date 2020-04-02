package astvalidator

import (
	"reflect"
	"testing"
	"time"
)

func TestCondition_Validate(t *testing.T) {
	type args struct {
		query  string
		object interface{}
	}
	tests := []struct {
		name        string
		args        args
		wantIsValid bool
		wantErr     bool
	}{
		{
			name: "Normal case - struct validation",
			args: args{
				query: `(id=1 && (member_id=12||member_id=2))  &&   (division=engineering || division=finance)`,
				object: struct {
					ID       string `json:"id"`
					MemberID string `json:"member_id"`
					Division string `json:"division"`
				}{
					ID:       "1",
					MemberID: "2",
					Division: "finance",
				},
			},
			wantIsValid: true,
			wantErr:     false,
		},
		{
			name: "Normal case - struct validation",
			args: args{
				query: `(id=1 &&  member_id=2  &&   (division=engineering || division=finance))||(member_id=3)`,
				object: struct {
					ID       int    `json:"id"`
					MemberID int    `json:"member_id"`
					Division string `json:"division"`
				}{
					ID:       1,
					MemberID: 3,
					Division: "finance",
				},
			},
			wantIsValid: true,
			wantErr:     false,
		},
		{
			name: "Normal case - struct validation - brand attribute is not exist",
			args: args{
				query: `(id=1 &&  member_id=2  &&   (division=engineering || division=finance))||(member_id=3&&brand=abc)`,
				object: struct {
					ID       int    `json:"id"`
					MemberID int    `json:"member_id"`
					Division string `json:"division"`
				}{
					ID:       1,
					MemberID: 3,
					Division: "finance",
				},
			},
			wantIsValid: false,
			wantErr:     false,
		},
		{
			name: "Normal case - struct validation",
			args: args{
				query: `id=1 &&  member_id=3  && ((division=engineering || division=finance || division=people)&&(member_id=2||id=1))`,
				object: struct {
					ID       int    `json:"id"`
					MemberID int    `json:"member_id"`
					Division string `json:"division"`
				}{
					ID:       1,
					MemberID: 3,
					Division: "people",
				},
			},
			wantIsValid: true,
			wantErr:     false,
		},
		{
			name: "Normal case - struct validation",
			args: args{
				query: `ID=1 &&  MemberID=2  &&   (Division=engineering || Division=finance)`,
				object: struct {
					ID       int
					MemberID int
					Division string
				}{
					ID:       1,
					MemberID: 2,
					Division: "engineering",
				},
			},
			wantIsValid: true,
			wantErr:     false,
		},
		{
			name: "Normal case - struct validation - Brand attribute is not exist",
			args: args{
				query: `ID=1 &&  MemberID=2  &&   (Division=engineering || Division=finance) && Brand=Adidas`,
				object: struct {
					ID       int
					MemberID int
					Division string
				}{
					ID:       1,
					MemberID: 2,
					Division: "engineering",
				},
			},
			wantIsValid: false,
			wantErr:     false,
		},
		{
			name: "Normal case - struct validation - skip not exist attribute because using OR condition",
			args: args{
				query: `ID=1 &&  MemberID=2  &&   (Division=engineering || Division=finance) && (Category=Bawahan || ID=1 || Brand=nike)`,
				object: struct {
					ID       int
					MemberID int
					Division string
				}{
					ID:       1,
					MemberID: 2,
					Division: "engineering",
				},
			},
			wantIsValid: true,
			wantErr:     false,
		},
		{
			name: "Error case",
			args: args{
				query:  `ID=1 &&  MemberID=2  &&   (Division=engineering || Division=finance)`,
				object: nil,
			},
			wantIsValid: false,
			wantErr:     true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			condition, _ := GenerateCondition(tt.args.query)
			gotIsValid, err := condition.Validate(tt.args.object)
			if (err != nil) != tt.wantErr {
				t.Errorf("Condition.Validate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotIsValid != tt.wantIsValid {
				t.Errorf("Condition.Validate() = %v, want %v", gotIsValid, tt.wantIsValid)
			}
		})
	}
}

func TestCondition_FilterSlice(t *testing.T) {
	type Account struct {
		ID        int        `json:"id"`
		MemberID  int        `json:"member_id"`
		Division  string     `json:"division"`
		Score     *int       `json:"score"`
		Point     *int64     `json:"point"`
		Wallet    *float32   `json:"wallet"`
		Money     *float64   `json:"money"`
		JoinDate  time.Time  `json:"join_date"`
		LeaveDate *time.Time `json:"leave_date"`
	}

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
	fTime := func(t time.Time) *time.Time {
		return &t
	}

	testData := []Account{
		{
			ID:        1,
			MemberID:  21,
			Division:  "people",
			Score:     fInt(90),
			Point:     fInt64(12000),
			Wallet:    fFloat(100000),
			Money:     fFloat64(10000),
			JoinDate:  time.Date(2020, 3, 9, 0, 0, 0, 0, time.UTC),
			LeaveDate: fTime(time.Date(2020, 12, 9, 0, 0, 0, 0, time.UTC)),
		},
		{
			ID:        2,
			MemberID:  22,
			Division:  "finance",
			Score:     fInt(40),
			Point:     fInt64(1000),
			Wallet:    fFloat(1000),
			Money:     fFloat64(50000),
			JoinDate:  time.Date(2014, 1, 9, 0, 0, 0, 0, time.UTC),
			LeaveDate: fTime(time.Date(2015, 12, 9, 0, 0, 0, 0, time.UTC)),
		},
		{
			ID:        3,
			MemberID:  23,
			Division:  "business",
			Score:     fInt(60),
			Point:     fInt64(5000),
			Wallet:    fFloat(5000),
			Money:     fFloat64(80000),
			JoinDate:  time.Date(2016, 12, 9, 0, 0, 0, 0, time.UTC),
			LeaveDate: fTime(time.Date(2017, 12, 9, 0, 0, 0, 0, time.UTC)),
		},
		{
			ID:        4,
			MemberID:  24,
			Division:  "managerial",
			Score:     fInt(70),
			Point:     fInt64(20000),
			Wallet:    fFloat(4000),
			Money:     fFloat64(900000),
			JoinDate:  time.Date(2018, 4, 9, 0, 0, 0, 0, time.UTC),
			LeaveDate: fTime(time.Date(2019, 12, 9, 0, 0, 0, 0, time.UTC)),
		},
		{
			ID:        5,
			MemberID:  25,
			Division:  "engineering",
			Score:     fInt(100),
			Point:     fInt64(3000),
			Wallet:    fFloat(100),
			Money:     fFloat64(1500000),
			JoinDate:  time.Date(2015, 10, 9, 0, 0, 0, 0, time.UTC),
			LeaveDate: nil,
		},
		{
			ID:        5,
			MemberID:  25,
			Division:  "engineering",
			Score:     nil,
			Point:     nil,
			Wallet:    nil,
			Money:     nil,
			JoinDate:  time.Date(2015, 7, 9, 0, 0, 0, 0, time.UTC),
			LeaveDate: fTime(time.Date(2016, 12, 9, 0, 0, 0, 0, time.UTC)),
		},
	}

	type args struct {
		query   string
		objects interface{}
	}
	tests := []struct {
		name        string
		args        args
		wantResults interface{}
		wantErr     bool
	}{
		{
			name: "Normal case",
			args: args{
				query:   "id=1||id=2",
				objects: testData,
			},
			wantResults: []Account{
				{
					ID:        1,
					MemberID:  21,
					Division:  "people",
					Score:     fInt(90),
					Point:     fInt64(12000),
					Wallet:    fFloat(100000),
					Money:     fFloat64(10000),
					JoinDate:  time.Date(2020, 3, 9, 0, 0, 0, 0, time.UTC),
					LeaveDate: fTime(time.Date(2020, 12, 9, 0, 0, 0, 0, time.UTC)),
				},
				{
					ID:        2,
					MemberID:  22,
					Division:  "finance",
					Score:     fInt(40),
					Point:     fInt64(1000),
					Wallet:    fFloat(1000),
					Money:     fFloat64(50000),
					JoinDate:  time.Date(2014, 1, 9, 0, 0, 0, 0, time.UTC),
					LeaveDate: fTime(time.Date(2015, 12, 9, 0, 0, 0, 0, time.UTC)),
				},
			},
			wantErr: false,
		},
		{
			name: "Normal case",
			args: args{
				query:   "(member_id=23)||((id=1 && member_id=21)&&(division=people||division=managerial))",
				objects: testData,
			},
			wantResults: []Account{
				{
					ID:        1,
					MemberID:  21,
					Division:  "people",
					Score:     fInt(90),
					Point:     fInt64(12000),
					Wallet:    fFloat(100000),
					Money:     fFloat64(10000),
					JoinDate:  time.Date(2020, 3, 9, 0, 0, 0, 0, time.UTC),
					LeaveDate: fTime(time.Date(2020, 12, 9, 0, 0, 0, 0, time.UTC)),
				},
				{
					ID:        3,
					MemberID:  23,
					Division:  "business",
					Score:     fInt(60),
					Point:     fInt64(5000),
					Wallet:    fFloat(5000),
					Money:     fFloat64(80000),
					JoinDate:  time.Date(2016, 12, 9, 0, 0, 0, 0, time.UTC),
					LeaveDate: fTime(time.Date(2017, 12, 9, 0, 0, 0, 0, time.UTC)),
				},
			},
			wantErr: false,
		},
		{
			name: "Normal case",
			args: args{
				query:   "member_id=21 && score=90",
				objects: testData,
			},
			wantResults: []Account{
				{
					ID:        1,
					MemberID:  21,
					Division:  "people",
					Score:     fInt(90),
					Point:     fInt64(12000),
					Wallet:    fFloat(100000),
					Money:     fFloat64(10000),
					JoinDate:  time.Date(2020, 3, 9, 0, 0, 0, 0, time.UTC),
					LeaveDate: fTime(time.Date(2020, 12, 9, 0, 0, 0, 0, time.UTC)),
				},
			},
			wantErr: false,
		},
		{
			name: "Normal case",
			args: args{
				query:   "member_id=25 && point>=3000",
				objects: testData,
			},
			wantResults: []Account{
				{
					ID:        5,
					MemberID:  25,
					Division:  "engineering",
					Score:     fInt(100),
					Point:     fInt64(3000),
					Wallet:    fFloat(100),
					Money:     fFloat64(1500000),
					JoinDate:  time.Date(2015, 10, 9, 0, 0, 0, 0, time.UTC),
					LeaveDate: nil,
				},
			},
			wantErr: false,
		},
		{
			name: "Normal case",
			args: args{
				query:   `join_date>"2015-01-01 00:00:00" && join_date<="2016-01-01 00:00:00" && score>80 && point<4000 && wallet>90 && money>=1500000`,
				objects: testData,
			},
			wantResults: []Account{
				{
					ID:        5,
					MemberID:  25,
					Division:  "engineering",
					Score:     fInt(100),
					Point:     fInt64(3000),
					Wallet:    fFloat(100),
					Money:     fFloat64(1500000),
					JoinDate:  time.Date(2015, 10, 9, 0, 0, 0, 0, time.UTC),
					LeaveDate: nil,
				},
			},
			wantErr: false,
		},
		{
			name: "Normal case - empty",
			args: args{
				query:   "(member_id=25)||((id=1 && member_id=21)&&(division=people||division=managerial))",
				objects: []Account{},
			},
			wantResults: []Account{},
			wantErr:     false,
		},
		{
			name: "Error case - nil object",
			args: args{
				query:   "(member_id=25)||((id=1 && member_id=21)&&(division=people||division=managerial))",
				objects: nil,
			},
			wantResults: nil,
			wantErr:     true,
		},
		{
			name: "Error case - invalid type",
			args: args{
				query:   "(member_id=25)||((id=1 && member_id=21)&&(division=people||division=managerial))",
				objects: Account{},
			},
			wantResults: nil,
			wantErr:     true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			condition, _ := GenerateCondition(tt.args.query)
			gotResults, err := condition.FilterSlice(tt.args.objects)
			if (err != nil) != tt.wantErr {
				t.Errorf("Condition.FilterSlice() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotResults, tt.wantResults) {
				t.Errorf("Condition.FilterSlice() = %v, want %v", gotResults, tt.wantResults)
			}
		})
	}
}

func TestCondition_ValidateObjects(t *testing.T) {
	type fields struct {
		Operator   string
		Attribute  *Attribute
		Conditions []*Condition
	}
	type args struct {
		query string
		data  interface{}
	}
	type firstStruct struct {
		ID       string `json:"id"`
		MemberID string `json:"member_id"`
		Division string `json:"division"`
	}
	type secondStruct struct {
		Name string `json:"name"`
	}
	type thirdStruct struct {
		Type    string `json:"type"`
		Segment string `json:"segment"`
	}

	thirdData := thirdStruct{
		Type:    "ABC",
		Segment: "new-member",
	}

	tests := []struct {
		name        string
		fields      fields
		args        args
		wantIsValid bool
		wantErr     bool
	}{
		{
			name: "Normal case - one struct validation",
			args: args{
				query: `member_id=345`,
				data: []interface{}{
					firstStruct{
						ID:       "123",
						MemberID: "345",
						Division: "engineering",
					},
				},
			},
			wantIsValid: true,
			wantErr:     false,
		},
		{
			name: "Normal case - one struct validation - attribute brand not exist",
			args: args{
				query: `member_id=345 && brand=adidas`,
				data: []interface{}{
					firstStruct{
						ID:       "123",
						MemberID: "345",
						Division: "engineering",
					},
				},
			},
			wantIsValid: false,
			wantErr:     false,
		},
		{
			name: "Normal case - multi struct validation - all attributes exist",
			args: args{
				query: `thirdStruct.type=ABC && secondStruct.name=Test`,
				data: []interface{}{
					thirdData,
					secondStruct{
						Name: "Test",
					},
				},
			},
			wantIsValid: true,
			wantErr:     false,
		},
		{
			name: "Normal case - multi struct validation - attribute memberId in secondStruct not exist",
			args: args{
				query: `secondStruct.name=Test && secondStruct.memberId=1010101`,
				data: []interface{}{
					thirdData,
					secondStruct{
						Name: "Test",
					},
				},
			},
			wantIsValid: false,
			wantErr:     false,
		},
		{
			name: "Normal case - multi struct validation",
			args: args{
				query: `firstStruct.id=123 && secondStruct.name=Test`,
				data: []interface{}{
					firstStruct{
						ID:       "123",
						MemberID: "345",
						Division: "engineering",
					},
					secondStruct{
						Name: "Test",
					},
				},
			},
			wantIsValid: true,
			wantErr:     false,
		},
		{
			name: "Normal case - struct validation",
			args: args{
				query: `id=123`,
				data: firstStruct{
					ID:       "123",
					MemberID: "345",
					Division: "engineering",
				},
			},
			wantIsValid: true,
			wantErr:     false,
		},
		{
			name: "Normal case - multi struct validation",
			args: args{
				query: `(firstStruct.id=1234 || secondStruct=Test || thirdStruct.segment=new-member) && (firstStruct.member_id=345 && secondStruct.name=Test) && thirdStruct.type=ABC`,
				data: []interface{}{
					firstStruct{
						ID:       "123",
						MemberID: "345",
						Division: "engineering",
					},
					secondStruct{
						Name: "Test",
					},
					thirdData,
				},
			},
			wantIsValid: true,
			wantErr:     false,
		},
		{
			name: "Error case",
			args: args{
				query: `id=1`,
				data:  []interface{}{},
			},
			wantIsValid: false,
			wantErr:     true,
		},
		{
			name: "Error case",
			args: args{
				query: `id=1`,
				data:  nil,
			},
			wantIsValid: false,
			wantErr:     true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			condition, _ := GenerateCondition(tt.args.query)
			gotIsValid, err := condition.ValidateObjects(tt.args.data)
			if (err != nil) != tt.wantErr {
				t.Errorf("Condition.ValidateObjects() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotIsValid != tt.wantIsValid {
				t.Errorf("Condition.ValidateObjects() = %v, want %v", gotIsValid, tt.wantIsValid)
			}
		})
	}
}
