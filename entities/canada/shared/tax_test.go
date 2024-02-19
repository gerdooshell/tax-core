package shared

import (
	"math"
	"testing"
)

func TestTax_Calculate(t1 *testing.T) {
	type fields struct {
		TaxBrackets      []TaxBracket
		calculatedAmount float64
	}
	type args struct {
		taxableIncome float64
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{name: "1. valid input", fields: fields{TaxBrackets: []TaxBracket{{High: 50000, Low: 0, Rate: 5}, {High: 100000, Low: 50000, Rate: 10}}}, args: args{taxableIncome: 0}, wantErr: false},
		{name: "2. valid input", fields: fields{TaxBrackets: []TaxBracket{{High: 50000, Low: 0, Rate: 5}, {High: 100000, Low: 50000, Rate: 10}}}, args: args{taxableIncome: 100000000000}, wantErr: false},
		{name: "invalid amount < 0", fields: fields{TaxBrackets: []TaxBracket{{High: 50000, Low: 0, Rate: 5}, {High: 100000, Low: 50000, Rate: 10}}}, args: args{taxableIncome: -1}, wantErr: true},
		{name: "invalid brackets", fields: fields{TaxBrackets: []TaxBracket{}}, args: args{taxableIncome: -1}, wantErr: true},
	}
	for _, tt := range tests {
		t1.Run(tt.name, func(t1 *testing.T) {
			t := &Tax{
				TaxBrackets:      tt.fields.TaxBrackets,
				calculatedAmount: tt.fields.calculatedAmount,
			}
			if err := t.Calculate(tt.args.taxableIncome, false); (err != nil) != tt.wantErr {
				t1.Errorf("Calculate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestTax_GetValue(t1 *testing.T) {
	type fields struct {
		TaxBrackets      []TaxBracket
		calculatedAmount float64
	}
	type args struct {
		taxableIncome float64
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   float64
	}{
		{name: "zero income", fields: fields{TaxBrackets: []TaxBracket{{High: 50000, Low: 0, Rate: 5}, {High: 100000, Low: 50000, Rate: 10}}}, args: args{taxableIncome: 0}, want: 0},
		{name: "low income", fields: fields{TaxBrackets: []TaxBracket{{High: 50000, Low: 0, Rate: 5}, {High: 100000, Low: 50000, Rate: 10}}}, args: args{taxableIncome: 10000}, want: 500},
		{name: "low edge income", fields: fields{TaxBrackets: []TaxBracket{{High: 50000, Low: 0, Rate: 5}, {High: 100000, Low: 50000, Rate: 10}}}, args: args{taxableIncome: 50000}, want: 2500},
		{name: "between income", fields: fields{TaxBrackets: []TaxBracket{{High: 50000, Low: 0, Rate: 5}, {High: 100000, Low: 50000, Rate: 10}}}, args: args{taxableIncome: 70000}, want: 4500},
		{name: "high edge income", fields: fields{TaxBrackets: []TaxBracket{{High: 50000, Low: 0, Rate: 5}, {High: 100000, Low: 50000, Rate: 10}}}, args: args{taxableIncome: 100000}, want: 7500},
		{name: "1. high income", fields: fields{TaxBrackets: []TaxBracket{{High: 50000, Low: 0, Rate: 5}, {High: 100000, Low: 50000, Rate: 10}, {High: math.MaxFloat64, Low: 100000, Rate: 20}}}, args: args{taxableIncome: 120000}, want: 11500},
		{name: "2. high income", fields: fields{TaxBrackets: []TaxBracket{{High: 50000, Low: 0, Rate: 5}, {High: 100000, Low: 50000, Rate: 10}, {High: math.MaxFloat64, Low: 100000, Rate: 20}}}, args: args{taxableIncome: 500000}, want: 87500},
	}
	for _, tt := range tests {
		t1.Run(tt.name, func(t1 *testing.T) {
			t := &Tax{
				TaxBrackets:      tt.fields.TaxBrackets,
				calculatedAmount: tt.fields.calculatedAmount,
			}
			_ = t.Calculate(tt.args.taxableIncome, false)
			if got := t.GetValue(); got != tt.want {
				t1.Errorf("GetEmployeeValue() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTax_validateParameters(t1 *testing.T) {
	type fields struct {
		TaxBrackets      []TaxBracket
		calculatedAmount float64
	}
	type args struct {
		taxableIncome float64
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{name: "1. valid input", fields: fields{TaxBrackets: []TaxBracket{{High: 50000, Low: 0, Rate: 5}, {High: 100000, Low: 50000, Rate: 10}}}, args: args{taxableIncome: 0}, wantErr: false},
		{name: "2. valid input", fields: fields{TaxBrackets: []TaxBracket{{High: 50000, Low: 0, Rate: 5}, {High: 100000, Low: 50000, Rate: 10}}}, args: args{taxableIncome: 100000000000}, wantErr: false},
		{name: "invalid amount < 0", fields: fields{TaxBrackets: []TaxBracket{{High: 50000, Low: 0, Rate: 5}, {High: 100000, Low: 50000, Rate: 10}}}, args: args{taxableIncome: -1}, wantErr: true},
		{name: "invalid brackets", fields: fields{TaxBrackets: []TaxBracket{}}, args: args{taxableIncome: -1}, wantErr: true},
	}
	for _, tt := range tests {
		t1.Run(tt.name, func(t1 *testing.T) {
			t := &Tax{
				TaxBrackets:      tt.fields.TaxBrackets,
				calculatedAmount: tt.fields.calculatedAmount,
			}
			if err := t.validateParameters(tt.args.taxableIncome); (err != nil) != tt.wantErr {
				t1.Errorf("validateParameters() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
