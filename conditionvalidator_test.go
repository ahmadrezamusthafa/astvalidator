package astvalidator

import "testing"

func TestCondition_ValidateCondition(t *testing.T) {
	tests := []struct {
		name           string
		referenceQuery string
		input          string
		wantIsValid    bool
		wantErr        bool
	}{
		{
			name:           "Normal case",
			referenceQuery: "id=1 && member_id=45",
			input:          "id=2 || member_id=45",
			wantIsValid:    false,
			wantErr:        false,
		},
		{
			name:           "Normal case - ignore case",
			referenceQuery: "name=Budi && brand=Arava && member_id=45",
			input:          "name=budi && brand=arava && member_id=45",
			wantIsValid:    true,
			wantErr:        false,
		},
		{
			name:           "Normal case",
			referenceQuery: "id=1 && member_id=45",
			input:          "(id=2||id=1) && member_id=45",
			wantIsValid:    true,
			wantErr:        false,
		},
		{
			name:           "Normal case",
			referenceQuery: "id=1 && member_id=45",
			input:          "id=1",
			wantIsValid:    false,
			wantErr:        false,
		},
		{
			name:           "Normal case",
			referenceQuery: "id=1 || member_id=45",
			input:          "id=1",
			wantIsValid:    true,
			wantErr:        false,
		},
		{
			name:           "Normal case",
			referenceQuery: "id=1 || member_id=45",
			input:          "member_id=45",
			wantIsValid:    true,
			wantErr:        false,
		},
		{
			name:           "Normal case",
			referenceQuery: "id=1 && member_id=45",
			input:          "id=1 && member_id=45",
			wantIsValid:    true,
			wantErr:        false,
		},
		{
			name:           "Normal case",
			referenceQuery: "id=1 && member_id=45",
			input:          "id=1 && (member_id=23||member_id=35)",
			wantIsValid:    false,
			wantErr:        false,
		},
		{
			name:           "Normal case",
			referenceQuery: "id=1 && member_id=45",
			input:          "id=1 && (member_id=23||member_id=45)",
			wantIsValid:    true,
			wantErr:        false,
		},
		{
			name:           "Normal case",
			referenceQuery: "id=1 && member_id=45",
			input:          "id=1 && member_id=44",
			wantIsValid:    false,
			wantErr:        false,
		},
		{
			name:           "Normal case",
			referenceQuery: "id=1 || member_id=45",
			input:          "id=1 && member_id=22",
			wantIsValid:    true,
			wantErr:        false,
		},
		{
			name:           "Normal case - condition group",
			referenceQuery: "(id=1 || id=2) && member_id=45",
			input:          "id=1 && member_id=45",
			wantIsValid:    true,
			wantErr:        false,
		},
		{
			name:           "Normal case - multi condition group",
			referenceQuery: "id=1 &&  member_id=3  && ((division=engineering || division=finance || division=people)&&(member_id=2||id=1))",
			input:          "id=1 && member_id=3 && division=engineering",
			wantIsValid:    true,
			wantErr:        false,
		},
		{
			name:           "Normal case - multi condition group",
			referenceQuery: "id=1 &&  member_id=3  && ((division=engineering || division=finance || division=people)&&(member_id=2||id=1))",
			input:          "(id=1 && member_id=3) && (division=tech&&division=finance)",
			wantIsValid:    false,
			wantErr:        false,
		},
		{
			name:           "Normal case - multi condition group",
			referenceQuery: "id=1 &&  member_id=3  && ((division=engineering || division=finance || division=people)&&(member_id=2||id=1))",
			input:          "((id=1 && member_id=3) && (division=tech||division=finance))",
			wantIsValid:    true,
			wantErr:        false,
		},
		{
			name:           "Normal case - multi condition group",
			referenceQuery: "(id=1 || id=2) && (member_id=45||member_id=10)",
			input:          "id=1 && (member_id=10)",
			wantIsValid:    true,
			wantErr:        false,
		},
		{
			name:           "Normal case - multi condition group",
			referenceQuery: "(id=1 || id=2) && (member_id=45||member_id=10)",
			input:          "id=3 && member_id=10 || id=14",
			wantIsValid:    false,
			wantErr:        false,
		},
		{
			name:           "Normal case - multi condition group",
			referenceQuery: "(id=1 || id=2) && (member_id=45||member_id=10) && (segment=trial||segment=free)",
			input:          "id=1 && member_id=10 && segment=free",
			wantIsValid:    true,
			wantErr:        false,
		},
		{
			name:           "Normal case - match only one input",
			referenceQuery: "id=1",
			input:          "id=1 && member_id=10 && segment=free",
			wantIsValid:    true,
			wantErr:        false,
		},
		{
			name:           "Normal case - string condition",
			referenceQuery: "deviceType=mobile || memberId=xxx",
			input:          "deviceType=mobile",
			wantIsValid:    true,
			wantErr:        false,
		},
		{
			name:           "Normal case - string condition",
			referenceQuery: "deviceType=mobile && memberId=xxx",
			input:          "deviceType=mobile",
			wantIsValid:    false,
			wantErr:        false,
		},
		{
			name:           "Normal case - rule engine case",
			referenceQuery: "deviceType=mobile && ABTest=xxx ",
			input:          "deviceType=mobile && ABTest=yyy",
			wantIsValid:    false,
			wantErr:        false,
		},
		{
			name:           "Normal case - rule engine case with group condition",
			referenceQuery: "deviceType=mobile && ABTest=xxx ",
			input:          "deviceType=mobile && (ABTest=yyy||ABTest=xxx)",
			wantIsValid:    true,
			wantErr:        false,
		},
		{
			name:           "Normal case - greater than operator - integer",
			referenceQuery: "(id=1 || id=2) && member_id>100",
			input:          "id=1 && member_id=111",
			wantIsValid:    true,
			wantErr:        false,
		},
		{
			name:           "Normal case - greater than operator - integer",
			referenceQuery: "(id=1 || id=2) && member_id>100",
			input:          "id=1 && member_id=111",
			wantIsValid:    true,
			wantErr:        false,
		},
		{
			name:           "Normal case - greater than equal operator - datetime",
			referenceQuery: `(id=1 || id=2) && create_date>="2020-02-02 12:12:12"`,
			input:          `id=1 && create_date="2020-02-02 12:12:12"`,
			wantIsValid:    true,
			wantErr:        false,
		},
		{
			name:           "Normal case - greater than operator - float",
			referenceQuery: "(id=1 || id=2) && price>1200.50 && (segment=hijaber||segment=girl||segment=cantik) && poin>100",
			input:          "id=1 && price=1200.51 && ((segment=cantik&&poin=58)||(segment=girl&&poin=518))",
			wantIsValid:    true,
			wantErr:        false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			referenceCondition, err := GenerateCondition(tt.referenceQuery)
			if err != nil {
				t.Errorf("Condition.ValidateCondition() referenceQuery error = %v", err)
				return
			}
			inputCondition, err := GenerateCondition(tt.input)
			if err != nil {
				t.Errorf("Condition.ValidateCondition() input error = %v", err)
				return
			}
			gotIsValid, err := referenceCondition.ValidateCondition(inputCondition)
			if (err != nil) != tt.wantErr {
				t.Errorf("Condition.ValidateCondition() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotIsValid != tt.wantIsValid {
				t.Errorf("Condition.ValidateCondition() = %v, want %v", gotIsValid, tt.wantIsValid)
			}
		})
	}
}

func Test_getValueType(t *testing.T) {
	type args struct {
		value string
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{
			name: "Normal case",
			args: args{
				value: "ajkshd-kjhad-ueydsi-asdk9238",
			},
			want: TypeAlphanumeric,
		},
		{
			name: "Normal case",
			args: args{
				value: "10.01.200.01",
			},
			want: TypeAlphanumeric,
		},
		{
			name: "Normal case",
			args: args{
				value: "19.123",
			},
			want: TypeNumeric,
		},
		{
			name: "Normal case",
			args: args{
				value: "123",
			},
			want: TypeNumeric,
		},
		{
			name: "Normal case",
			args: args{
				value: "2020-02-02 12:00:21",
			},
			want: TypeTime,
		},
		{
			name: "Normal case",
			args: args{
				value: "1245.",
			},
			want: TypeAlphanumeric,
		},
		{
			name: "Normal case",
			args: args{
				value: ".123.",
			},
			want: TypeAlphanumeric,
		},
		{
			name: "Normal case",
			args: args{
				value: "free-member",
			},
			want: TypeAlphanumeric,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getValueType(tt.args.value); got != tt.want {
				t.Errorf("getValueType() = %v, want %v", got, tt.want)
			}
		})
	}
}
