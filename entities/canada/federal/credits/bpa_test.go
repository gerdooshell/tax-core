package credits

import "testing"

func TestBasicPensionAmount_Calculate(t *testing.T) {
	type fields struct {
		MaxBPAIncome float64
		MinBPAIncome float64
		MinBPAAmount float64
		MaxBPAAmount float64
		value        float64
	}
	type args struct {
		income float64
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{name: "positive income", fields: fields{MaxBPAIncome: 50000, MaxBPAAmount: 15000, MinBPAIncome: 100000, MinBPAAmount: 10000}, args: args{income: 10000}, wantErr: false},
		{name: "zero income", fields: fields{MaxBPAIncome: 50000, MaxBPAAmount: 15000, MinBPAIncome: 100000, MinBPAAmount: 10000}, args: args{income: 0}, wantErr: false},
		{name: "negative income", fields: fields{MaxBPAIncome: 50000, MaxBPAAmount: 15000, MinBPAIncome: 100000, MinBPAAmount: 10000}, args: args{income: -10}, wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bpa := &BasicPersonalAmount{
				MaxBPAIncome: tt.fields.MaxBPAIncome,
				MinBPAIncome: tt.fields.MinBPAIncome,
				MinBPAAmount: tt.fields.MinBPAAmount,
				MaxBPAAmount: tt.fields.MaxBPAAmount,
				value:        tt.fields.value,
			}
			if err := bpa.Calculate(tt.args.income); (err != nil) != tt.wantErr {
				t.Errorf("Calculate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestBasicPensionAmount_GetValue(t *testing.T) {
	type fields struct {
		MaxBPAIncome float64
		MinBPAIncome float64
		MinBPAAmount float64
		MaxBPAAmount float64
		value        float64
	}
	type args struct {
		income float64
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   float64
	}{
		{name: "zero income", fields: fields{MaxBPAIncome: 50000, MaxBPAAmount: 15000, MinBPAIncome: 100000, MinBPAAmount: 10000}, args: args{income: 0}, want: 15000},
		{name: "low income", fields: fields{MaxBPAIncome: 50000, MaxBPAAmount: 15000, MinBPAIncome: 100000, MinBPAAmount: 10000}, args: args{income: 10000}, want: 15000},
		{name: "low edge income", fields: fields{MaxBPAIncome: 50000, MaxBPAAmount: 15000, MinBPAIncome: 100000, MinBPAAmount: 10000}, args: args{income: 50000}, want: 15000},
		{name: "middle income", fields: fields{MaxBPAIncome: 50000, MaxBPAAmount: 15000, MinBPAIncome: 100000, MinBPAAmount: 10000}, args: args{income: 80000}, want: 12000},
		{name: "high edge income", fields: fields{MaxBPAIncome: 50000, MaxBPAAmount: 15000, MinBPAIncome: 100000, MinBPAAmount: 10000}, args: args{income: 100000}, want: 10000},
		{name: "high income", fields: fields{MaxBPAIncome: 50000, MaxBPAAmount: 15000, MinBPAIncome: 100000, MinBPAAmount: 10000}, args: args{income: 150000}, want: 10000},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bpa := &BasicPersonalAmount{
				MaxBPAIncome: tt.fields.MaxBPAIncome,
				MinBPAIncome: tt.fields.MinBPAIncome,
				MinBPAAmount: tt.fields.MinBPAAmount,
				MaxBPAAmount: tt.fields.MaxBPAAmount,
				value:        tt.fields.value,
			}
			_ = bpa.Calculate(tt.args.income)
			if got := bpa.GetValue(); got != tt.want {
				t.Errorf("GetEmployeeValue() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBasicPensionAmount_validate(t *testing.T) {
	type fields struct {
		MaxBPAIncome float64
		MinBPAIncome float64
		MinBPAAmount float64
		MaxBPAAmount float64
		value        float64
	}
	type args struct {
		income float64
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{name: "positive income", fields: fields{MaxBPAIncome: 50000, MaxBPAAmount: 15000, MinBPAIncome: 100000, MinBPAAmount: 10000}, args: args{income: 10000}, wantErr: false},
		{name: "zero income", fields: fields{MaxBPAIncome: 50000, MaxBPAAmount: 15000, MinBPAIncome: 100000, MinBPAAmount: 10000}, args: args{income: 0}, wantErr: false},
		{name: "negative income", fields: fields{MaxBPAIncome: 50000, MaxBPAAmount: 15000, MinBPAIncome: 100000, MinBPAAmount: 10000}, args: args{income: -10}, wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bpa := &BasicPersonalAmount{
				MaxBPAIncome: tt.fields.MaxBPAIncome,
				MinBPAIncome: tt.fields.MinBPAIncome,
				MinBPAAmount: tt.fields.MinBPAAmount,
				MaxBPAAmount: tt.fields.MaxBPAAmount,
				value:        tt.fields.value,
			}
			if err := bpa.validate(tt.args.income); (err != nil) != tt.wantErr {
				t.Errorf("validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
