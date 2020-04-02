package astvalidator

import (
	"bytes"
	"encoding/json"
	"reflect"
	"strings"
	"testing"
)

func TestGenerateConditionQueryStructure(t *testing.T) {
	type args struct {
		query string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "Normal case",
			args: args{
				query: `
                      (id = 1 
                      && member_id = 2 )
                      || (
                        division = engineering 
                        || division = finance
                      )
`,
			},
			want:    `{"conditions":[{"conditions":[{"attribute":{"name":"id","operator":"=","value":"1"}},{"operator":"AND","attribute":{"name":"member_id","operator":"=","value":"2"}}]},{"operator":"OR","conditions":[{"attribute":{"name":"division","operator":"=","value":"engineering"}},{"operator":"OR","attribute":{"name":"division","operator":"=","value":"finance"}}]}]}`,
			wantErr: false,
		},
		{
			name: "Normal case",
			args: args{
				query: `
                      (id = 1 
                      && member_id = 2 )
                      || (
                        division = engineering 
                        || division = finance
                      ) && user_id = 43
`,
			},
			want:    `{"conditions":[{"conditions":[{"attribute":{"name":"id","operator":"=","value":"1"}},{"operator":"AND","attribute":{"name":"member_id","operator":"=","value":"2"}}]},{"operator":"OR","conditions":[{"attribute":{"name":"division","operator":"=","value":"engineering"}},{"operator":"OR","attribute":{"name":"division","operator":"=","value":"finance"}}]},{"operator":"AND","attribute":{"name":"user_id","operator":"=","value":"43"}}]}`,
			wantErr: false,
		},
		{
			name: "Normal case",
			args: args{
				query: `
                      id = 1 
                      && member_id = 2 
                      && (
                        division = engineering 
                        || division = finance
                      )
`,
			},
			want:    `{"conditions":[{"attribute":{"name":"id","operator":"=","value":"1"}},{"operator":"AND","attribute":{"name":"member_id","operator":"=","value":"2"}},{"operator":"AND","conditions":[{"attribute":{"name":"division","operator":"=","value":"engineering"}},{"operator":"OR","attribute":{"name":"division","operator":"=","value":"finance"}}]}]}`,
			wantErr: false,
		},
		{
			name: "Normal case",
			args: args{
				query: `id=1 &&  member_id=2   &&   (division=engineering || division=finance)`,
			},
			want:    `{"conditions":[{"attribute":{"name":"id","operator":"=","value":"1"}},{"operator":"AND","attribute":{"name":"member_id","operator":"=","value":"2"}},{"operator":"AND","conditions":[{"attribute":{"name":"division","operator":"=","value":"engineering"}},{"operator":"OR","attribute":{"name":"division","operator":"=","value":"finance"}}]}]}`,
			wantErr: false,
		},
		{
			name: "Normal case",
			args: args{
				query: `
                  id = 1 
                  && member_id = 2 
                  && user_id = 3 
                  && (
                    province = jatim 
                    || city = mojokerto
                    || (
                      warehouse_id = 1 
                      && warehouse_detail_id = 2
                    )
                  )
`,
			},
			want:    `{"conditions":[{"attribute":{"name":"id","operator":"=","value":"1"}},{"operator":"AND","attribute":{"name":"member_id","operator":"=","value":"2"}},{"operator":"AND","attribute":{"name":"user_id","operator":"=","value":"3"}},{"operator":"AND","conditions":[{"attribute":{"name":"province","operator":"=","value":"jatim"}},{"operator":"OR","attribute":{"name":"city","operator":"=","value":"mojokerto"}},{"operator":"OR","conditions":[{"attribute":{"name":"warehouse_id","operator":"=","value":"1"}},{"operator":"AND","attribute":{"name":"warehouse_detail_id","operator":"=","value":"2"}}]}]}]}`,
			wantErr: false,
		},
		{
			name: "Normal case",
			args: args{
				query: `
                  id = 1
                  && member_id = 2
                  && user_id = 3
                  && (
                    province = jatim
                    || city = mojokerto
                    || (
                      warehouse_id = 1
                      && warehouse_detail_id = 2
                    )
                  )
				  && data_id = 54
`,
			},
			want:    `{"conditions":[{"attribute":{"name":"id","operator":"=","value":"1"}},{"operator":"AND","attribute":{"name":"member_id","operator":"=","value":"2"}},{"operator":"AND","attribute":{"name":"user_id","operator":"=","value":"3"}},{"operator":"AND","conditions":[{"attribute":{"name":"province","operator":"=","value":"jatim"}},{"operator":"OR","attribute":{"name":"city","operator":"=","value":"mojokerto"}},{"operator":"OR","conditions":[{"attribute":{"name":"warehouse_id","operator":"=","value":"1"}},{"operator":"AND","attribute":{"name":"warehouse_detail_id","operator":"=","value":"2"}}]}]},{"operator":"AND","attribute":{"name":"data_id","operator":"=","value":"54"}}]}`,
			wantErr: false,
		},
		{
			name: "Normal case",
			args: args{
				query: "((date<=2019-09-09 && date > 2019-08-08) || (p_date>=2019-01-01 && p_date<2019-02-02)) && (member_type=1||member_type=2)",
			},
			want:    `{"conditions":[{"conditions":[{"conditions":[{"attribute":{"name":"date","operator":"\u003c=","value":"2019-09-09"}},{"operator":"AND","attribute":{"name":"date","operator":"\u003e","value":"2019-08-08"}}]},{"operator":"OR","conditions":[{"attribute":{"name":"p_date","operator":"\u003e=","value":"2019-01-01"}},{"operator":"AND","attribute":{"name":"p_date","operator":"\u003c","value":"2019-02-02"}}]}]},{"operator":"AND","conditions":[{"attribute":{"name":"member_type","operator":"=","value":"1"}},{"operator":"OR","attribute":{"name":"member_type","operator":"=","value":"2"}}]}]}`,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GenerateCondition(tt.args.query)
			if (err != nil) != tt.wantErr {
				t.Errorf("GenerateCondition() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			byteBuf, _ := json.Marshal(got)
			if !strings.EqualFold(string(byteBuf), tt.want) {
				t.Errorf("GenerateCondition() = %v, want %v", string(byteBuf), tt.want)
			}
		})
	}
}

func Test_getToken(t *testing.T) {
	type args struct {
		value string
	}
	tests := []struct {
		name string
		args args
		want []*TokenAttribute
	}{
		{
			name: "Normal case",
			args: args{
				value: "id=1 &&  member_id=2   &&   (division=engineering      || division=finance)",
			},
			want: []*TokenAttribute{
				{
					value: "id",
				},
				{
					value: "=",
				},
				{
					value: "1",
				},
				{
					value: "&&",
				},
				{
					value: "member_id",
				},
				{
					value: "=",
				},
				{
					value: "2",
				},
				{
					value: "&&",
				},
				{
					value: "(",
				},
				{
					value: "division",
				},
				{
					value: "=",
				},
				{
					value: "engineering",
				},
				{
					value: "||",
				},
				{
					value: "division",
				},
				{
					value: "=",
				},
				{
					value: "finance",
				},
				{
					value: ")",
				},
			},
		},
		{
			name: "Normal case",
			args: args{
				value: "id>1 &&  member_id>=2 && (test_id<10 || pr_id<=28)",
			},
			want: []*TokenAttribute{
				{
					value: "id",
				},
				{
					value: ">",
				},
				{
					value: "1",
				},
				{
					value: "&&",
				},
				{
					value: "member_id",
				},
				{
					value: ">=",
				},
				{
					value: "2",
				},
				{
					value: "&&",
				},
				{
					value: "(",
				},
				{
					value: "test_id",
				},
				{
					value: "<",
				},
				{
					value: "10",
				},
				{
					value: "||",
				},
				{
					value: "pr_id",
				},
				{
					value: "<=",
				},
				{
					value: "28",
				},
				{
					value: ")",
				},
			},
		},
		{
			name: "Normal case",
			args: args{
				value: "date>2019-09-01 && date<=2019-10-10 && (segment_id=12||segment_id=13)",
			},
			want: []*TokenAttribute{
				{
					value: "date",
				},
				{
					value: ">",
				},
				{
					value: "2019-09-01",
				},
				{
					value: "&&",
				},
				{
					value: "date",
				},
				{
					value: "<=",
				},
				{
					value: "2019-10-10",
				},
				{
					value: "&&",
				},
				{
					value: "(",
				},
				{
					value: "segment_id",
				},
				{
					value: "=",
				},
				{
					value: "12",
				},
				{
					value: "||",
				},
				{
					value: "segment_id",
				},
				{
					value: "=",
				},
				{
					value: "13",
				},
				{
					value: ")",
				},
			},
		},
		{
			name: "Normal case",
			args: args{
				value: `date>"2019-09-01 00:10:00" && date<=2019-10-10 && (segment_id="12"||segment_id=13)`,
			},
			want: []*TokenAttribute{
				{
					value: "date",
				},
				{
					value: ">",
				},
				{
					value: "2019-09-01 00:10:00",
				},
				{
					value: "&&",
				},
				{
					value: "date",
				},
				{
					value: "<=",
				},
				{
					value: "2019-10-10",
				},
				{
					value: "&&",
				},
				{
					value: "(",
				},
				{
					value: "segment_id",
				},
				{
					value: "=",
				},
				{
					value: "12",
				},
				{
					value: "||",
				},
				{
					value: "segment_id",
				},
				{
					value: "=",
				},
				{
					value: "13",
				},
				{
					value: ")",
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getTokenAttributes(tt.args.value); !reflect.DeepEqual(got, tt.want) {
				strbGot := bytes.Buffer{}
				for _, g := range got {
					strbGot.WriteString("\"" + g.value + "\" ")
				}
				strbWant := bytes.Buffer{}
				for _, g := range tt.want {
					strbWant.WriteString("\"" + g.value + "\" ")
				}
				t.Errorf("getTokenAttributes() = %v, want %v", strbGot.String(), strbWant.String())
			}
		})
	}
}
