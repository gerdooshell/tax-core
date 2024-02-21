package credits

import "testing"

func TestCanadaEmploymentAmount_Calculate(t *testing.T) {
	type fields struct {
		Value float64
	}
	type args struct {
		totalIncome float64
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{name: "CAE for 2023", fields: fields{Value: 1368}, args: args{totalIncome: 1}, wantErr: false},
		{name: "negative CAE", fields: fields{Value: -2}, args: args{totalIncome: 1}, wantErr: true},
		{name: "negative income", fields: fields{Value: 1368}, args: args{totalIncome: -1}, wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cea := &CanadaEmploymentAmount{
				Value: tt.fields.Value,
			}
			if err := cea.Calculate(tt.args.totalIncome); (err != nil) != tt.wantErr {
				t.Errorf("Calculate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestCanadaEmploymentAmount_GetValue(t *testing.T) {
	type fields struct {
		Value float64
	}
	type args struct {
		totalIncome float64
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   float64
	}{
		{name: "get value", fields: fields{Value: 1368}, args: args{totalIncome: 2000}, want: 1368},
		{name: "get value", fields: fields{Value: 1368}, args: args{totalIncome: 200}, want: 200},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cea := &CanadaEmploymentAmount{
				Value: tt.fields.Value,
			}
			_ = cea.Calculate(tt.args.totalIncome)
			if got := cea.GetEmployeeValue(); got != tt.want {
				t.Errorf("GetEmployeeValue() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCanadaEmploymentAmount_validate(t *testing.T) {
	type fields struct {
		Value float64
	}
	type args struct {
		totalIncome float64
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{name: "CAE for 2023", fields: fields{Value: 1368}, args: args{totalIncome: 2000}, wantErr: false},
		{name: "negative CAE", fields: fields{Value: -2}, args: args{totalIncome: 2000}, wantErr: true},
		{name: "negative income", fields: fields{Value: 1368}, args: args{totalIncome: -2000}, wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cea := &CanadaEmploymentAmount{
				Value: tt.fields.Value,
			}
			if err := cea.validate(tt.args.totalIncome); (err != nil) != tt.wantErr {
				t.Errorf("validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
