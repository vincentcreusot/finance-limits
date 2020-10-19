package logic

import (
	"encoding/json"
	"reflect"
	"testing"
	"time"
)

func Test_UnmarshalJSON(t *testing.T) {
	type args struct {
		jsonString string
	}
	type output struct {
		load     inputLoad
		hasError bool
	}
	tests := []struct {
		name string
		args args
		want output
	}{
		{
			name: "MarshallingOk",
			args: args{
				jsonString: `{"id": "1234","customer_id": "2345","load_amount": "$123.45","time": "2018-01-01T00:00:00Z"}`,
			},
			want: output{
				load: inputLoad{
					LoadID:     "1234",
					CustomerID: "2345",
					Amount:     loadAmount{Value: 123.45},
					Time:       time.Date(2018, time.January, 1, 0, 0, 0, 0, time.UTC),
				},
				hasError: false,
			},
		},
		{
			name: "MarshallingAmountFailed",
			args: args{
				jsonString: `{"id": "1234","customer_id": "2345","load_amount": "AAAAAAAAAA","time": "2018-01-01T00:00:00Z"}`,
			},
			want: output{
				load: inputLoad{
					LoadID:     "1234",
					CustomerID: "2345",
					Amount:     loadAmount{Value: 0},
					Time:       time.Date(0001, time.January, 1, 0, 0, 0, 0, time.UTC), // time not parsed because of error in amount
				},
				hasError: true,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := inputLoad{}
			err := json.Unmarshal([]byte(tt.args.jsonString), &l)
			if (err != nil) != tt.want.hasError || l != tt.want.load {
				t.Errorf("unmarshallJson = %v and %v, want %v", l, err, tt.want)
			}
		})
	}
}

func Test_validateLoad(t *testing.T) {
	type args struct {
		load         inputLoad
		historyLoads []inputLoad
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			"validateAcceptedNoHistory",
			args{
				load: inputLoad{
					LoadID:     "1",
					CustomerID: "1",
					Amount:     loadAmount{Value: 3000},
					Time:       time.Date(2000, time.Month(1), 1, 10, 0, 0, 0, time.UTC),
				},
				historyLoads: nil,
			},
			true,
		},
		{
			"validateRefusedMaxAmountNoHistory",
			args{
				load: inputLoad{
					LoadID:     "1",
					CustomerID: "1",
					Amount:     loadAmount{Value: 6000},
					Time:       time.Date(2000, time.Month(1), 1, 10, 0, 0, 0, time.UTC),
				},
				historyLoads: nil,
			},
			false,
		},
		{
			"validateMaxAmountDifferentDays",
			args{
				load: inputLoad{
					LoadID:     "2",
					CustomerID: "1",
					Amount:     loadAmount{Value: 3000},
					Time:       time.Date(2000, time.Month(1), 2, 10, 0, 0, 0, time.UTC),
				},
				historyLoads: []inputLoad{
					inputLoad{
						LoadID:     "1",
						CustomerID: "1",
						Amount:     loadAmount{Value: 3000},
						Time:       time.Date(2000, time.Month(1), 1, 10, 0, 0, 0, time.UTC),
					},
				},
			},
			true,
		},
		{
			"validateMaxAmountDay",
			args{
				load: inputLoad{
					LoadID:     "2",
					CustomerID: "1",
					Amount:     loadAmount{Value: 3000},
					Time:       time.Date(2000, time.Month(1), 1, 10, 0, 0, 0, time.UTC),
				},
				historyLoads: []inputLoad{
					inputLoad{
						LoadID:     "1",
						CustomerID: "1",
						Amount:     loadAmount{Value: 3000},
						Time:       time.Date(2000, time.Month(1), 1, 10, 0, 0, 0, time.UTC),
					},
				},
			},
			false,
		},
		{
			"validateExactMaxAmountDay",
			args{
				load: inputLoad{
					LoadID:     "2",
					CustomerID: "1",
					Amount:     loadAmount{Value: 3000},
					Time:       time.Date(2000, time.Month(1), 1, 10, 0, 0, 0, time.UTC),
				},
				historyLoads: []inputLoad{
					inputLoad{
						LoadID:     "1",
						CustomerID: "1",
						Amount:     loadAmount{Value: 2000},
						Time:       time.Date(2000, time.Month(1), 1, 10, 0, 0, 0, time.UTC),
					},
				},
			},
			true,
		},
		{
			"validateMaxCountDay",
			args{
				load: inputLoad{
					LoadID:     "4",
					CustomerID: "1",
					Amount:     loadAmount{Value: 1000},
					Time:       time.Date(2000, time.Month(1), 1, 15, 0, 0, 0, time.UTC),
				},
				historyLoads: []inputLoad{
					inputLoad{
						LoadID:     "1",
						CustomerID: "1",
						Amount:     loadAmount{Value: 1000},
						Time:       time.Date(2000, time.Month(1), 1, 10, 0, 0, 0, time.UTC),
					},
					inputLoad{
						LoadID:     "2",
						CustomerID: "1",
						Amount:     loadAmount{Value: 1000},
						Time:       time.Date(2000, time.Month(1), 1, 11, 0, 0, 0, time.UTC),
					},
					inputLoad{
						LoadID:     "3",
						CustomerID: "1",
						Amount:     loadAmount{Value: 1000},
						Time:       time.Date(2000, time.Month(1), 1, 12, 0, 0, 0, time.UTC),
					},
				},
			},
			false,
		},
		{
			"validateExactCountDay",
			args{
				load: inputLoad{
					LoadID:     "4",
					CustomerID: "1",
					Amount:     loadAmount{Value: 1000},
					Time:       time.Date(2000, time.Month(1), 1, 15, 0, 0, 0, time.UTC),
				},
				historyLoads: []inputLoad{
					inputLoad{
						LoadID:     "1",
						CustomerID: "1",
						Amount:     loadAmount{Value: 1000},
						Time:       time.Date(2000, time.Month(1), 1, 10, 0, 0, 0, time.UTC),
					},
					inputLoad{
						LoadID:     "2",
						CustomerID: "1",
						Amount:     loadAmount{Value: 1000},
						Time:       time.Date(2000, time.Month(1), 1, 11, 0, 0, 0, time.UTC),
					},
				},
			},
			true,
		},
		{
			"validateMaxCountOnDifferentDays",
			args{
				load: inputLoad{
					LoadID:     "4",
					CustomerID: "1",
					Amount:     loadAmount{Value: 1000},
					Time:       time.Date(2000, time.Month(1), 2, 15, 0, 0, 0, time.UTC),
				},
				historyLoads: []inputLoad{
					inputLoad{
						LoadID:     "1",
						CustomerID: "1",
						Amount:     loadAmount{Value: 1000},
						Time:       time.Date(2000, time.Month(1), 1, 10, 0, 0, 0, time.UTC),
					},
					inputLoad{
						LoadID:     "2",
						CustomerID: "1",
						Amount:     loadAmount{Value: 1000},
						Time:       time.Date(2000, time.Month(1), 1, 11, 0, 0, 0, time.UTC),
					},
					inputLoad{
						LoadID:     "3",
						CustomerID: "1",
						Amount:     loadAmount{Value: 1000},
						Time:       time.Date(2000, time.Month(1), 1, 12, 0, 0, 0, time.UTC),
					},
				},
			},
			true,
		},
		{
			"validateExactMaxOnWeek",
			args{
				load: inputLoad{
					LoadID:     "4",
					CustomerID: "1",
					Amount:     loadAmount{Value: 5000},
					Time:       time.Date(2020, time.Month(1), 9, 15, 0, 0, 0, time.UTC),
				},
				historyLoads: []inputLoad{
					inputLoad{
						LoadID:     "1",
						CustomerID: "1",
						Amount:     loadAmount{Value: 5000},
						Time:       time.Date(2020, time.Month(1), 6, 10, 0, 0, 0, time.UTC), // 6th Jan 2020 is a Monday
					},
					inputLoad{
						LoadID:     "2",
						CustomerID: "1",
						Amount:     loadAmount{Value: 5000},
						Time:       time.Date(2020, time.Month(1), 7, 11, 0, 0, 0, time.UTC),
					},
					inputLoad{
						LoadID:     "3",
						CustomerID: "1",
						Amount:     loadAmount{Value: 5000},
						Time:       time.Date(2020, time.Month(1), 8, 12, 0, 0, 0, time.UTC),
					},
				},
			},
			true,
		},
		{
			"validateMaxOnWeek",
			args{
				load: inputLoad{
					LoadID:     "5",
					CustomerID: "1",
					Amount:     loadAmount{Value: 4000},
					Time:       time.Date(2020, time.Month(1), 10, 15, 0, 0, 0, time.UTC),
				},
				historyLoads: []inputLoad{
					inputLoad{
						LoadID:     "1",
						CustomerID: "1",
						Amount:     loadAmount{Value: 5000},
						Time:       time.Date(2020, time.Month(1), 6, 10, 0, 0, 0, time.UTC), // 6 Jan is a Monday
					},
					inputLoad{
						LoadID:     "2",
						CustomerID: "1",
						Amount:     loadAmount{Value: 5000},
						Time:       time.Date(2020, time.Month(1), 7, 11, 0, 0, 0, time.UTC),
					},
					inputLoad{
						LoadID:     "3",
						CustomerID: "1",
						Amount:     loadAmount{Value: 5000},
						Time:       time.Date(2020, time.Month(1), 8, 12, 0, 0, 0, time.UTC),
					},
					inputLoad{
						LoadID:     "4",
						CustomerID: "1",
						Amount:     loadAmount{Value: 3000},
						Time:       time.Date(2020, time.Month(1), 9, 12, 0, 0, 0, time.UTC),
					},
				},
			},
			false,
		},
		{
			"validateMaxOnTwoWeeks",
			args{
				load: inputLoad{
					LoadID:     "5",
					CustomerID: "1",
					Amount:     loadAmount{Value: 4000},
					Time:       time.Date(2020, time.Month(1), 13, 15, 0, 0, 0, time.UTC),
				},
				historyLoads: []inputLoad{
					inputLoad{
						LoadID:     "1",
						CustomerID: "1",
						Amount:     loadAmount{Value: 5000},
						Time:       time.Date(2020, time.Month(1), 6, 10, 0, 0, 0, time.UTC), // 6 Jan is a Monday
					},
					inputLoad{
						LoadID:     "2",
						CustomerID: "1",
						Amount:     loadAmount{Value: 5000},
						Time:       time.Date(2020, time.Month(1), 7, 11, 0, 0, 0, time.UTC),
					},
					inputLoad{
						LoadID:     "3",
						CustomerID: "1",
						Amount:     loadAmount{Value: 5000},
						Time:       time.Date(2020, time.Month(1), 8, 12, 0, 0, 0, time.UTC),
					},
					inputLoad{
						LoadID:     "4",
						CustomerID: "1",
						Amount:     loadAmount{Value: 3000},
						Time:       time.Date(2020, time.Month(1), 9, 12, 0, 0, 0, time.UTC),
					},
				},
			},
			true,
		},
		{
			"validateMidnightOnWeeks",
			args{
				load: inputLoad{
					LoadID:     "5",
					CustomerID: "1",
					Amount:     loadAmount{Value: 4000},
					Time:       time.Date(2020, time.Month(1), 12, 23, 59, 59, 0, time.UTC).Add(time.Second),
				},
				historyLoads: []inputLoad{
					inputLoad{
						LoadID:     "1",
						CustomerID: "1",
						Amount:     loadAmount{Value: 5000},
						Time:       time.Date(2020, time.Month(1), 6, 10, 0, 0, 0, time.UTC), // 6 Jan is a Monday
					},
					inputLoad{
						LoadID:     "2",
						CustomerID: "1",
						Amount:     loadAmount{Value: 5000},
						Time:       time.Date(2020, time.Month(1), 7, 11, 0, 0, 0, time.UTC),
					},
					inputLoad{
						LoadID:     "3",
						CustomerID: "1",
						Amount:     loadAmount{Value: 5000},
						Time:       time.Date(2020, time.Month(1), 8, 12, 0, 0, 0, time.UTC),
					},
					inputLoad{
						LoadID:     "4",
						CustomerID: "1",
						Amount:     loadAmount{Value: 3000},
						Time:       time.Date(2020, time.Month(1), 9, 12, 0, 0, 0, time.UTC),
					},
				},
			},
			true,
		},
		{
			"validateMidnightOnTwoDays",
			args{
				load: inputLoad{
					LoadID:     "5",
					CustomerID: "1",
					Amount:     loadAmount{Value: 4000},
					Time:       time.Date(2020, time.Month(1), 7, 0, 0, 0, 0, time.UTC),
				},
				historyLoads: []inputLoad{
					inputLoad{
						LoadID:     "1",
						CustomerID: "1",
						Amount:     loadAmount{Value: 3000},
						Time:       time.Date(2020, time.Month(1), 6, 10, 0, 0, 0, time.UTC), // 6 Jan is a Monday
					},
				},
			},
			true,
		},
		{
			"validateMidnightMinusOneSecond",
			args{
				load: inputLoad{
					LoadID:     "5",
					CustomerID: "1",
					Amount:     loadAmount{Value: 4000},
					Time:       time.Date(2020, time.Month(1), 7, 0, 0, 0, 0, time.UTC).Add(-time.Second),
				},
				historyLoads: []inputLoad{
					inputLoad{
						LoadID:     "1",
						CustomerID: "1",
						Amount:     loadAmount{Value: 3000},
						Time:       time.Date(2020, time.Month(1), 6, 10, 0, 0, 0, time.UTC), // 6 Jan is a Monday
					},
				},
			},
			false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := validateLoad(tt.args.load, tt.args.historyLoads); got != tt.want {
				t.Errorf("validateLoad = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_validateLoadAndFillHistory(t *testing.T) {
	type args struct {
		load                 inputLoad
		customerHistoryLoads map[string][]inputLoad
	}
	type output struct {
		returnedValue        bool
		customerHistoryLoads map[string][]inputLoad
	}
	tests := []struct {
		name string
		args args
		want output
	}{
		{
			"validateEmptyAccepted",
			args{
				load: inputLoad{
					LoadID:     "1",
					CustomerID: "1",
					Amount:     loadAmount{Value: 3000},
					Time:       time.Date(2000, time.Month(1), 1, 10, 0, 0, 0, time.UTC),
				},
				customerHistoryLoads: make(map[string][]inputLoad),
			},
			output{
				returnedValue: true,
				customerHistoryLoads: map[string][]inputLoad{
					"1": {
						inputLoad{
							LoadID:     "1",
							CustomerID: "1",
							Amount:     loadAmount{Value: 3000},
							Time:       time.Date(2000, time.Month(1), 1, 10, 0, 0, 0, time.UTC),
						}},
				},
			},
		},
		{
			"validateEmptyNotAccepted",
			args{
				load: inputLoad{
					LoadID:     "1",
					CustomerID: "1",
					Amount:     loadAmount{Value: 6000},
					Time:       time.Date(2000, time.Month(1), 1, 10, 0, 0, 0, time.UTC),
				},
				customerHistoryLoads: make(map[string][]inputLoad),
			},
			output{
				returnedValue:        false,
				customerHistoryLoads: make(map[string][]inputLoad),
			},
		},
		{
			"validateNotEmptyAccepted",
			args{
				load: inputLoad{
					LoadID:     "1",
					CustomerID: "1",
					Amount:     loadAmount{Value: 3000},
					Time:       time.Date(2000, time.Month(1), 1, 10, 0, 0, 0, time.UTC),
				},
				customerHistoryLoads: map[string][]inputLoad{
					"1": {
						inputLoad{
							LoadID:     "2",
							CustomerID: "1",
							Amount:     loadAmount{Value: 1000},
							Time:       time.Date(2000, time.Month(1), 1, 5, 0, 0, 0, time.UTC),
						},
					},
				},
			},
			output{
				returnedValue: true,
				customerHistoryLoads: map[string][]inputLoad{
					"1": {
						inputLoad{
							LoadID:     "2",
							CustomerID: "1",
							Amount:     loadAmount{Value: 1000},
							Time:       time.Date(2000, time.Month(1), 1, 5, 0, 0, 0, time.UTC),
						},
						inputLoad{
							LoadID:     "1",
							CustomerID: "1",
							Amount:     loadAmount{Value: 3000},
							Time:       time.Date(2000, time.Month(1), 1, 10, 0, 0, 0, time.UTC),
						}},
				},
			},
		},
		{
			"validateNotEmptyNotAccepted",
			args{
				load: inputLoad{
					LoadID:     "1",
					CustomerID: "1",
					Amount:     loadAmount{Value: 3000},
					Time:       time.Date(2000, time.Month(1), 1, 10, 0, 0, 0, time.UTC),
				},
				customerHistoryLoads: map[string][]inputLoad{
					"1": {
						inputLoad{
							LoadID:     "2",
							CustomerID: "1",
							Amount:     loadAmount{Value: 3000},
							Time:       time.Date(2000, time.Month(1), 1, 5, 0, 0, 0, time.UTC),
						},
					},
				},
			},
			output{
				returnedValue: false,
				customerHistoryLoads: map[string][]inputLoad{
					"1": {
						inputLoad{
							LoadID:     "2",
							CustomerID: "1",
							Amount:     loadAmount{Value: 3000},
							Time:       time.Date(2000, time.Month(1), 1, 5, 0, 0, 0, time.UTC),
						}},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			loadParser := NewFinanceLogic()
			loadParser.CustomersLoads = tt.args.customerHistoryLoads
			if got := loadParser.validateLoadAndFillHistory(tt.args.load); got != tt.want.returnedValue || !reflect.DeepEqual(loadParser.CustomersLoads, tt.want.customerHistoryLoads) {
				t.Errorf("validateLoadAndFillHistory = %v and %v, want %v", got, loadParser.CustomersLoads, tt.want)
			}
		})
	}
}

func Test_addCustomerLoadToTreated(t *testing.T) {
	type args struct {
		load           inputLoad
		treatedLoadIds map[customerLoadID]interface{}
	}

	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "NewCustomerLoad",
			args: args{
				load: inputLoad{
					LoadID:     "1",
					CustomerID: "1",
					Amount:     loadAmount{},
					Time:       time.Time{},
				},
				treatedLoadIds: make(map[customerLoadID]interface{}),
			},
			want: true,
		},
		{
			name: "ExistingCustomerLoad",
			args: args{
				load: inputLoad{
					LoadID:     "1",
					CustomerID: "1",
					Amount:     loadAmount{},
					Time:       time.Time{},
				},
				treatedLoadIds: map[customerLoadID]interface{}{
					customerLoadID{
						LoadID:     "1",
						CustomerID: "1",
					}: nil,
				},
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			loadParser := NewFinanceLogic()
			loadParser.TreatedLoadIds = tt.args.treatedLoadIds
			if got := loadParser.addCustomerLoadToTreated(tt.args.load); got != tt.want {
				t.Errorf("addCustomerLoadToTreated = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_ParseLoads(t *testing.T) {
	type args struct {
		loadStrings []string
	}
	type output struct {
		loadResponses  []string
		numberOfErrors int
	}
	tests := []struct {
		name string
		args args
		want output
	}{
		{
			name: "oneLoad",
			args: args{
				loadStrings: []string{
					`{"id": "1234","customer_id": "2345","load_amount": "$123.45","time": "2018-01-01T00:00:00Z"}`,
				},
			},
			want: output{
				loadResponses:  []string{`{"id":"1234","customer_id":"2345","accepted":true}`},
				numberOfErrors: 0,
			}},
		{
			name: "emptyLoad",
			args: args{
				loadStrings: make([]string, 0),
			},
			want: output{
				loadResponses:  make([]string, 0),
				numberOfErrors: 0,
			}},
		{
			name: "jsonError",
			args: args{
				loadStrings: []string{`anerrorinjson`},
			},
			want: output{
				loadResponses:  make([]string, 0),
				numberOfErrors: 1,
			}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			loadParser := NewFinanceLogic()
			stringChannel := make(chan string)
			go func(stringChan chan string, stringsToLoad []string) {
				for _, line := range stringsToLoad {
					stringChan <- line
				}
				close(stringChan)
			}(stringChannel, tt.args.loadStrings)
			if gotLoads, gotErrors := loadParser.ParseLoads(stringChannel); !reflect.DeepEqual(gotLoads, tt.want.loadResponses) || len(gotErrors) != tt.want.numberOfErrors {
				t.Errorf("ParseLoads = %v,%v want %v", gotLoads, gotErrors, tt.want)
			}
		})
	}
}
