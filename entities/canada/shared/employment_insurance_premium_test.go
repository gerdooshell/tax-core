package shared

import "testing"

func TestEmploymentInsurancePremium_Calculate(t *testing.T) {
	type fields struct {
		MaxInsurableEarning               float64
		Rate                              float64
		EmployerEmployeeContributionRatio float64
		eiEmployeeContribution            float64
		eiEmployerContribution            float64
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
		{name: "valid inputs", fields: fields{MaxInsurableEarning: 50000, Rate: 5, EmployerEmployeeContributionRatio: 1.4}, args: args{totalIncome: 100000}, wantErr: false},
		{name: "invalid MIE", fields: fields{MaxInsurableEarning: 0, Rate: 5, EmployerEmployeeContributionRatio: 1.4}, args: args{totalIncome: 100000}, wantErr: true},
		{name: "1. edge Rate", fields: fields{MaxInsurableEarning: 50000, Rate: 0, EmployerEmployeeContributionRatio: 1.4}, args: args{totalIncome: 100000}, wantErr: false},
		{name: "2. edge Rate", fields: fields{MaxInsurableEarning: 50000, Rate: 100, EmployerEmployeeContributionRatio: 1.4}, args: args{totalIncome: 100000}, wantErr: false},
		{name: "invalid Rate", fields: fields{MaxInsurableEarning: 50000, Rate: 101, EmployerEmployeeContributionRatio: 1.4}, args: args{totalIncome: 100000}, wantErr: true},
		{name: "invalid EECR", fields: fields{MaxInsurableEarning: 50000, Rate: 5, EmployerEmployeeContributionRatio: -1}, args: args{totalIncome: 100000}, wantErr: true},
		{name: "1. invalid Salary", fields: fields{MaxInsurableEarning: 50000, Rate: 5, EmployerEmployeeContributionRatio: 1.4}, args: args{totalIncome: 0}, wantErr: false},
		{name: "2. invalid Salary", fields: fields{MaxInsurableEarning: 50000, Rate: 5, EmployerEmployeeContributionRatio: 1.4}, args: args{totalIncome: -100}, wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			eip := &EmploymentInsurancePremium{
				MaxInsurableEarning:               tt.fields.MaxInsurableEarning,
				Rate:                              tt.fields.Rate,
				EmployerEmployeeContributionRatio: tt.fields.EmployerEmployeeContributionRatio,
				eiEmployeeContribution:            tt.fields.eiEmployeeContribution,
				eiEmployerContribution:            tt.fields.eiEmployerContribution,
			}
			if err := eip.Calculate(tt.args.totalIncome); (err != nil) != tt.wantErr {
				t.Errorf("Calculate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestEmploymentInsurancePremium_GetEIEmployee(t *testing.T) {
	type fields struct {
		MaxInsurableEarning               float64
		Rate                              float64
		EmployerEmployeeContributionRatio float64
		eiEmployeeContribution            float64
		eiEmployerContribution            float64
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
		{name: "over max salary", fields: fields{MaxInsurableEarning: 50000, Rate: 1.63, EmployerEmployeeContributionRatio: 1.4}, args: args{totalIncome: 100000}, want: 815},
		{name: "high edge salary", fields: fields{MaxInsurableEarning: 50000, Rate: 1.63, EmployerEmployeeContributionRatio: 1.4}, args: args{totalIncome: 50000}, want: 815},
		{name: "1. calculate for salary", fields: fields{MaxInsurableEarning: 50000, Rate: 1.63, EmployerEmployeeContributionRatio: 1.4}, args: args{totalIncome: 40000}, want: 652},
		{name: "2. calculate for salary", fields: fields{MaxInsurableEarning: 50000, Rate: 1.63, EmployerEmployeeContributionRatio: 1.4}, args: args{totalIncome: 4000}, want: 65.2},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			eip := &EmploymentInsurancePremium{
				MaxInsurableEarning:               tt.fields.MaxInsurableEarning,
				Rate:                              tt.fields.Rate,
				EmployerEmployeeContributionRatio: tt.fields.EmployerEmployeeContributionRatio,
				eiEmployeeContribution:            tt.fields.eiEmployeeContribution,
				eiEmployerContribution:            tt.fields.eiEmployerContribution,
			}
			_ = eip.Calculate(tt.args.totalIncome)
			if got := eip.GetEIEmployee(); got != tt.want {
				t.Errorf("GetEIEmployee() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestEmploymentInsurancePremium_GetEIEmployer(t *testing.T) {
	type fields struct {
		MaxInsurableEarning               float64
		Rate                              float64
		EmployerEmployeeContributionRatio float64
		eiEmployeeContribution            float64
		eiEmployerContribution            float64
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
		{name: "over max salary", fields: fields{MaxInsurableEarning: 50000, Rate: 1.63, EmployerEmployeeContributionRatio: 1.4}, args: args{totalIncome: 100000}, want: 815 * 1.4},
		{name: "high edge salary", fields: fields{MaxInsurableEarning: 50000, Rate: 1.63, EmployerEmployeeContributionRatio: 1.4}, args: args{totalIncome: 50000}, want: 815 * 1.4},
		{name: "1. calculate for salary", fields: fields{MaxInsurableEarning: 50000, Rate: 1.63, EmployerEmployeeContributionRatio: 1.4}, args: args{totalIncome: 40000}, want: 652 * 1.4},
		{name: "2. calculate for salary", fields: fields{MaxInsurableEarning: 50000, Rate: 1.63, EmployerEmployeeContributionRatio: 1.4}, args: args{totalIncome: 4000}, want: 65.2 * 1.4},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			eip := &EmploymentInsurancePremium{
				MaxInsurableEarning:               tt.fields.MaxInsurableEarning,
				Rate:                              tt.fields.Rate,
				EmployerEmployeeContributionRatio: tt.fields.EmployerEmployeeContributionRatio,
				eiEmployeeContribution:            tt.fields.eiEmployeeContribution,
				eiEmployerContribution:            tt.fields.eiEmployerContribution,
			}
			_ = eip.Calculate(tt.args.totalIncome)
			if got := eip.GetEIEmployer(); got != tt.want {
				t.Errorf("GetEIEmployer() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestEmploymentInsurancePremium_validateProperties(t *testing.T) {
	type fields struct {
		MaxInsurableEarning               float64
		Rate                              float64
		EmployerEmployeeContributionRatio float64
		eiEmployeeContribution            float64
		eiEmployerContribution            float64
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
		{name: "valid inputs", fields: fields{MaxInsurableEarning: 50000, Rate: 5, EmployerEmployeeContributionRatio: 1.4}, args: args{totalIncome: 100000}, wantErr: false},
		{name: "invalid MIE", fields: fields{MaxInsurableEarning: 0, Rate: 5, EmployerEmployeeContributionRatio: 1.4}, args: args{totalIncome: 100000}, wantErr: true},
		{name: "1. edge Rate", fields: fields{MaxInsurableEarning: 50000, Rate: 0, EmployerEmployeeContributionRatio: 1.4}, args: args{totalIncome: 100000}, wantErr: false},
		{name: "2. edge Rate", fields: fields{MaxInsurableEarning: 50000, Rate: 100, EmployerEmployeeContributionRatio: 1.4}, args: args{totalIncome: 100000}, wantErr: false},
		{name: "invalid Rate", fields: fields{MaxInsurableEarning: 50000, Rate: 101, EmployerEmployeeContributionRatio: 1.4}, args: args{totalIncome: 100000}, wantErr: true},
		{name: "invalid EECR", fields: fields{MaxInsurableEarning: 50000, Rate: 5, EmployerEmployeeContributionRatio: -1}, args: args{totalIncome: 100000}, wantErr: true},
		{name: "1. invalid Salary", fields: fields{MaxInsurableEarning: 50000, Rate: 5, EmployerEmployeeContributionRatio: 1.4}, args: args{totalIncome: 0}, wantErr: false},
		{name: "2. invalid Salary", fields: fields{MaxInsurableEarning: 50000, Rate: 5, EmployerEmployeeContributionRatio: 1.4}, args: args{totalIncome: -100}, wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			eip := &EmploymentInsurancePremium{
				MaxInsurableEarning:               tt.fields.MaxInsurableEarning,
				Rate:                              tt.fields.Rate,
				EmployerEmployeeContributionRatio: tt.fields.EmployerEmployeeContributionRatio,
				eiEmployeeContribution:            tt.fields.eiEmployeeContribution,
				eiEmployerContribution:            tt.fields.eiEmployerContribution,
			}
			if err := eip.validateProperties(tt.args.totalIncome); (err != nil) != tt.wantErr {
				t.Errorf("validateProperties() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
