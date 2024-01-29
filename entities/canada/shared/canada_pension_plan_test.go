package shared

import "testing"

func TestCanadaPensionPlan_Calculate(t *testing.T) {
	type fields struct {
		Year                            int
		BasicRateEmployee               float64
		BasicRateEmployer               float64
		FirstAdditionalRateEmployee     float64
		FirstAdditionalRateEmployer     float64
		SecondAdditionalRateEmployee    float64
		SecondAdditionalRateEmployer    float64
		BasicExemption                  float64
		MaxPensionableEarning           float64
		AdditionalMaxPensionableEarning float64
		cppBasicEmployee                float64
		cppBasicEmployer                float64
		cppBasicSelfEmployed            float64
		cppFirstAdditionalEmployee      float64
		cppFirstAdditionalEmployer      float64
		cppFirstAdditionalSelfEmployed  float64
		cppSecondAdditionalEmployee     float64
		cppSecondAdditionalEmployer     float64
		cppSecondAdditionalSelfEmployed float64
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
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cpp := &CanadaPensionPlan{
				Year:                            tt.fields.Year,
				BasicRateEmployee:               tt.fields.BasicRateEmployee,
				BasicRateEmployer:               tt.fields.BasicRateEmployer,
				FirstAdditionalRateEmployee:     tt.fields.FirstAdditionalRateEmployee,
				FirstAdditionalRateEmployer:     tt.fields.FirstAdditionalRateEmployer,
				SecondAdditionalRateEmployee:    tt.fields.SecondAdditionalRateEmployee,
				SecondAdditionalRateEmployer:    tt.fields.SecondAdditionalRateEmployer,
				BasicExemption:                  tt.fields.BasicExemption,
				MaxPensionableEarning:           tt.fields.MaxPensionableEarning,
				AdditionalMaxPensionableEarning: tt.fields.AdditionalMaxPensionableEarning,
				cppBasicEmployee:                tt.fields.cppBasicEmployee,
				cppBasicEmployer:                tt.fields.cppBasicEmployer,
				cppBasicSelfEmployed:            tt.fields.cppBasicSelfEmployed,
				cppFirstAdditionalEmployee:      tt.fields.cppFirstAdditionalEmployee,
				cppFirstAdditionalEmployer:      tt.fields.cppFirstAdditionalEmployer,
				cppFirstAdditionalSelfEmployed:  tt.fields.cppFirstAdditionalSelfEmployed,
				cppSecondAdditionalEmployee:     tt.fields.cppSecondAdditionalEmployee,
				cppSecondAdditionalEmployer:     tt.fields.cppSecondAdditionalEmployer,
				cppSecondAdditionalSelfEmployed: tt.fields.cppSecondAdditionalSelfEmployed,
			}
			if err := cpp.Calculate(tt.args.totalIncome); (err != nil) != tt.wantErr {
				t.Errorf("Calculate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestCanadaPensionPlan_GetCPPBasicEmployee(t *testing.T) {
	type fields struct {
		Year                            int
		BasicRateEmployee               float64
		BasicRateEmployer               float64
		FirstAdditionalRateEmployee     float64
		FirstAdditionalRateEmployer     float64
		SecondAdditionalRateEmployee    float64
		SecondAdditionalRateEmployer    float64
		BasicExemption                  float64
		MaxPensionableEarning           float64
		AdditionalMaxPensionableEarning float64
		cppBasicEmployee                float64
		cppBasicEmployer                float64
		cppBasicSelfEmployed            float64
		cppFirstAdditionalEmployee      float64
		cppFirstAdditionalEmployer      float64
		cppFirstAdditionalSelfEmployed  float64
		cppSecondAdditionalEmployee     float64
		cppSecondAdditionalEmployer     float64
		cppSecondAdditionalSelfEmployed float64
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
		{name: "zero income", fields: fields{Year: 2023, BasicExemption: 3500, BasicRateEmployee: 4.95, BasicRateEmployer: 4.95, FirstAdditionalRateEmployee: 1, FirstAdditionalRateEmployer: 1, SecondAdditionalRateEmployee: 0, SecondAdditionalRateEmployer: 0, MaxPensionableEarning: 66600}, args: args{totalIncome: 0}, want: 0},
		{name: "low income 1", fields: fields{Year: 2023, BasicExemption: 3500, BasicRateEmployee: 4.95, BasicRateEmployer: 4.95, FirstAdditionalRateEmployee: 1, FirstAdditionalRateEmployer: 1, SecondAdditionalRateEmployee: 0, SecondAdditionalRateEmployer: 0, MaxPensionableEarning: 66600}, args: args{totalIncome: 10000}, want: 321.75},
		{name: "low income 2", fields: fields{Year: 2023, BasicExemption: 3500, BasicRateEmployee: 4.95, BasicRateEmployer: 4.95, FirstAdditionalRateEmployee: 1, FirstAdditionalRateEmployer: 1, SecondAdditionalRateEmployee: 0, SecondAdditionalRateEmployer: 0, MaxPensionableEarning: 66600}, args: args{totalIncome: 60000}, want: 2796.75},
		{name: "edge income", fields: fields{Year: 2023, BasicExemption: 3500, BasicRateEmployee: 4.95, BasicRateEmployer: 4.95, FirstAdditionalRateEmployee: 1, FirstAdditionalRateEmployer: 1, SecondAdditionalRateEmployee: 0, SecondAdditionalRateEmployer: 0, MaxPensionableEarning: 66600}, args: args{totalIncome: 66600}, want: 3123.45},
		{name: "high income 1", fields: fields{Year: 2023, BasicExemption: 3500, BasicRateEmployee: 4.95, BasicRateEmployer: 4.95, FirstAdditionalRateEmployee: 1, FirstAdditionalRateEmployer: 1, SecondAdditionalRateEmployee: 0, SecondAdditionalRateEmployer: 0, MaxPensionableEarning: 66600}, args: args{totalIncome: 70000}, want: 3123.45},
		{name: "high income 2", fields: fields{Year: 2023, BasicExemption: 3500, BasicRateEmployee: 4.95, BasicRateEmployer: 4.95, FirstAdditionalRateEmployee: 1, FirstAdditionalRateEmployer: 1, SecondAdditionalRateEmployee: 0, SecondAdditionalRateEmployer: 0, MaxPensionableEarning: 66600}, args: args{totalIncome: 500000}, want: 3123.45},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cpp := &CanadaPensionPlan{
				Year:                            tt.fields.Year,
				BasicRateEmployee:               tt.fields.BasicRateEmployee,
				BasicRateEmployer:               tt.fields.BasicRateEmployer,
				FirstAdditionalRateEmployee:     tt.fields.FirstAdditionalRateEmployee,
				FirstAdditionalRateEmployer:     tt.fields.FirstAdditionalRateEmployer,
				SecondAdditionalRateEmployee:    tt.fields.SecondAdditionalRateEmployee,
				SecondAdditionalRateEmployer:    tt.fields.SecondAdditionalRateEmployer,
				BasicExemption:                  tt.fields.BasicExemption,
				MaxPensionableEarning:           tt.fields.MaxPensionableEarning,
				AdditionalMaxPensionableEarning: tt.fields.AdditionalMaxPensionableEarning,
				cppBasicEmployee:                tt.fields.cppBasicEmployee,
				cppBasicEmployer:                tt.fields.cppBasicEmployer,
				cppBasicSelfEmployed:            tt.fields.cppBasicSelfEmployed,
				cppFirstAdditionalEmployee:      tt.fields.cppFirstAdditionalEmployee,
				cppFirstAdditionalEmployer:      tt.fields.cppFirstAdditionalEmployer,
				cppFirstAdditionalSelfEmployed:  tt.fields.cppFirstAdditionalSelfEmployed,
				cppSecondAdditionalEmployee:     tt.fields.cppSecondAdditionalEmployee,
				cppSecondAdditionalEmployer:     tt.fields.cppSecondAdditionalEmployer,
				cppSecondAdditionalSelfEmployed: tt.fields.cppSecondAdditionalSelfEmployed,
			}
			_ = cpp.Calculate(tt.args.totalIncome)
			if got := cpp.GetCPPBasicEmployee(); got != tt.want {
				t.Errorf("GetCPPBasicEmployee() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCanadaPensionPlan_GetCPPBasicEmployer(t *testing.T) {
	type fields struct {
		Year                            int
		BasicRateEmployee               float64
		BasicRateEmployer               float64
		FirstAdditionalRateEmployee     float64
		FirstAdditionalRateEmployer     float64
		SecondAdditionalRateEmployee    float64
		SecondAdditionalRateEmployer    float64
		BasicExemption                  float64
		MaxPensionableEarning           float64
		AdditionalMaxPensionableEarning float64
		cppBasicEmployee                float64
		cppBasicEmployer                float64
		cppBasicSelfEmployed            float64
		cppFirstAdditionalEmployee      float64
		cppFirstAdditionalEmployer      float64
		cppFirstAdditionalSelfEmployed  float64
		cppSecondAdditionalEmployee     float64
		cppSecondAdditionalEmployer     float64
		cppSecondAdditionalSelfEmployed float64
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
		{name: "zero income", fields: fields{Year: 2023, BasicExemption: 3500, BasicRateEmployee: 4.95, BasicRateEmployer: 4.95, FirstAdditionalRateEmployee: 1, FirstAdditionalRateEmployer: 1, SecondAdditionalRateEmployee: 0, SecondAdditionalRateEmployer: 0, MaxPensionableEarning: 66600}, args: args{totalIncome: 0}, want: 0},
		{name: "low income 1", fields: fields{Year: 2023, BasicExemption: 3500, BasicRateEmployee: 4.95, BasicRateEmployer: 4.95, FirstAdditionalRateEmployee: 1, FirstAdditionalRateEmployer: 1, SecondAdditionalRateEmployee: 0, SecondAdditionalRateEmployer: 0, MaxPensionableEarning: 66600}, args: args{totalIncome: 10000}, want: 321.75},
		{name: "low income 2", fields: fields{Year: 2023, BasicExemption: 3500, BasicRateEmployee: 4.95, BasicRateEmployer: 4.95, FirstAdditionalRateEmployee: 1, FirstAdditionalRateEmployer: 1, SecondAdditionalRateEmployee: 0, SecondAdditionalRateEmployer: 0, MaxPensionableEarning: 66600}, args: args{totalIncome: 60000}, want: 2796.75},
		{name: "edge income", fields: fields{Year: 2023, BasicExemption: 3500, BasicRateEmployee: 4.95, BasicRateEmployer: 4.95, FirstAdditionalRateEmployee: 1, FirstAdditionalRateEmployer: 1, SecondAdditionalRateEmployee: 0, SecondAdditionalRateEmployer: 0, MaxPensionableEarning: 66600}, args: args{totalIncome: 66600}, want: 3123.45},
		{name: "high income 1", fields: fields{Year: 2023, BasicExemption: 3500, BasicRateEmployee: 4.95, BasicRateEmployer: 4.95, FirstAdditionalRateEmployee: 1, FirstAdditionalRateEmployer: 1, SecondAdditionalRateEmployee: 0, SecondAdditionalRateEmployer: 0, MaxPensionableEarning: 66600}, args: args{totalIncome: 70000}, want: 3123.45},
		{name: "high income 2", fields: fields{Year: 2023, BasicExemption: 3500, BasicRateEmployee: 4.95, BasicRateEmployer: 4.95, FirstAdditionalRateEmployee: 1, FirstAdditionalRateEmployer: 1, SecondAdditionalRateEmployee: 0, SecondAdditionalRateEmployer: 0, MaxPensionableEarning: 66600}, args: args{totalIncome: 500000}, want: 3123.45},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cpp := &CanadaPensionPlan{
				Year:                            tt.fields.Year,
				BasicRateEmployee:               tt.fields.BasicRateEmployee,
				BasicRateEmployer:               tt.fields.BasicRateEmployer,
				FirstAdditionalRateEmployee:     tt.fields.FirstAdditionalRateEmployee,
				FirstAdditionalRateEmployer:     tt.fields.FirstAdditionalRateEmployer,
				SecondAdditionalRateEmployee:    tt.fields.SecondAdditionalRateEmployee,
				SecondAdditionalRateEmployer:    tt.fields.SecondAdditionalRateEmployer,
				BasicExemption:                  tt.fields.BasicExemption,
				MaxPensionableEarning:           tt.fields.MaxPensionableEarning,
				AdditionalMaxPensionableEarning: tt.fields.AdditionalMaxPensionableEarning,
				cppBasicEmployee:                tt.fields.cppBasicEmployee,
				cppBasicEmployer:                tt.fields.cppBasicEmployer,
				cppBasicSelfEmployed:            tt.fields.cppBasicSelfEmployed,
				cppFirstAdditionalEmployee:      tt.fields.cppFirstAdditionalEmployee,
				cppFirstAdditionalEmployer:      tt.fields.cppFirstAdditionalEmployer,
				cppFirstAdditionalSelfEmployed:  tt.fields.cppFirstAdditionalSelfEmployed,
				cppSecondAdditionalEmployee:     tt.fields.cppSecondAdditionalEmployee,
				cppSecondAdditionalEmployer:     tt.fields.cppSecondAdditionalEmployer,
				cppSecondAdditionalSelfEmployed: tt.fields.cppSecondAdditionalSelfEmployed,
			}
			_ = cpp.Calculate(tt.args.totalIncome)
			if got := cpp.GetCPPBasicEmployer(); got != tt.want {
				t.Errorf("GetCPPBasicEmployer() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCanadaPensionPlan_GetCPPBasicSelfEmployed(t *testing.T) {
	type fields struct {
		Year                            int
		BasicRateEmployee               float64
		BasicRateEmployer               float64
		FirstAdditionalRateEmployee     float64
		FirstAdditionalRateEmployer     float64
		SecondAdditionalRateEmployee    float64
		SecondAdditionalRateEmployer    float64
		BasicExemption                  float64
		MaxPensionableEarning           float64
		AdditionalMaxPensionableEarning float64
		cppBasicEmployee                float64
		cppBasicEmployer                float64
		cppBasicSelfEmployed            float64
		cppFirstAdditionalEmployee      float64
		cppFirstAdditionalEmployer      float64
		cppFirstAdditionalSelfEmployed  float64
		cppSecondAdditionalEmployee     float64
		cppSecondAdditionalEmployer     float64
		cppSecondAdditionalSelfEmployed float64
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
		{name: "zero income", fields: fields{Year: 2023, BasicExemption: 3500, BasicRateEmployee: 4.95, BasicRateEmployer: 4.95, FirstAdditionalRateEmployee: 1, FirstAdditionalRateEmployer: 1, SecondAdditionalRateEmployee: 0, SecondAdditionalRateEmployer: 0, MaxPensionableEarning: 66600}, args: args{totalIncome: 0}, want: 0},
		{name: "low edge income", fields: fields{Year: 2023, BasicExemption: 3500, BasicRateEmployee: 4.95, BasicRateEmployer: 4.95, FirstAdditionalRateEmployee: 1, FirstAdditionalRateEmployer: 1, SecondAdditionalRateEmployee: 0, SecondAdditionalRateEmployer: 0, MaxPensionableEarning: 66600}, args: args{totalIncome: 3500}, want: 0},
		{name: "low income 1", fields: fields{Year: 2023, BasicExemption: 3500, BasicRateEmployee: 4.95, BasicRateEmployer: 4.95, FirstAdditionalRateEmployee: 1, FirstAdditionalRateEmployer: 1, SecondAdditionalRateEmployee: 0, SecondAdditionalRateEmployer: 0, MaxPensionableEarning: 66600}, args: args{totalIncome: 10000}, want: 321.75 * 2},
		{name: "low income 2", fields: fields{Year: 2023, BasicExemption: 3500, BasicRateEmployee: 4.95, BasicRateEmployer: 4.95, FirstAdditionalRateEmployee: 1, FirstAdditionalRateEmployer: 1, SecondAdditionalRateEmployee: 0, SecondAdditionalRateEmployer: 0, MaxPensionableEarning: 66600}, args: args{totalIncome: 60000}, want: 2796.75 * 2},
		{name: "edge income", fields: fields{Year: 2023, BasicExemption: 3500, BasicRateEmployee: 4.95, BasicRateEmployer: 4.95, FirstAdditionalRateEmployee: 1, FirstAdditionalRateEmployer: 1, SecondAdditionalRateEmployee: 0, SecondAdditionalRateEmployer: 0, MaxPensionableEarning: 66600}, args: args{totalIncome: 66600}, want: 3123.45 * 2},
		{name: "high income 1", fields: fields{Year: 2023, BasicExemption: 3500, BasicRateEmployee: 4.95, BasicRateEmployer: 4.95, FirstAdditionalRateEmployee: 1, FirstAdditionalRateEmployer: 1, SecondAdditionalRateEmployee: 0, SecondAdditionalRateEmployer: 0, MaxPensionableEarning: 66600}, args: args{totalIncome: 70000}, want: 3123.45 * 2},
		{name: "high income 2", fields: fields{Year: 2023, BasicExemption: 3500, BasicRateEmployee: 4.95, BasicRateEmployer: 4.95, FirstAdditionalRateEmployee: 1, FirstAdditionalRateEmployer: 1, SecondAdditionalRateEmployee: 0, SecondAdditionalRateEmployer: 0, MaxPensionableEarning: 66600}, args: args{totalIncome: 500000}, want: 3123.45 * 2},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cpp := &CanadaPensionPlan{
				Year:                            tt.fields.Year,
				BasicRateEmployee:               tt.fields.BasicRateEmployee,
				BasicRateEmployer:               tt.fields.BasicRateEmployer,
				FirstAdditionalRateEmployee:     tt.fields.FirstAdditionalRateEmployee,
				FirstAdditionalRateEmployer:     tt.fields.FirstAdditionalRateEmployer,
				SecondAdditionalRateEmployee:    tt.fields.SecondAdditionalRateEmployee,
				SecondAdditionalRateEmployer:    tt.fields.SecondAdditionalRateEmployer,
				BasicExemption:                  tt.fields.BasicExemption,
				MaxPensionableEarning:           tt.fields.MaxPensionableEarning,
				AdditionalMaxPensionableEarning: tt.fields.AdditionalMaxPensionableEarning,
				cppBasicEmployee:                tt.fields.cppBasicEmployee,
				cppBasicEmployer:                tt.fields.cppBasicEmployer,
				cppBasicSelfEmployed:            tt.fields.cppBasicSelfEmployed,
				cppFirstAdditionalEmployee:      tt.fields.cppFirstAdditionalEmployee,
				cppFirstAdditionalEmployer:      tt.fields.cppFirstAdditionalEmployer,
				cppFirstAdditionalSelfEmployed:  tt.fields.cppFirstAdditionalSelfEmployed,
				cppSecondAdditionalEmployee:     tt.fields.cppSecondAdditionalEmployee,
				cppSecondAdditionalEmployer:     tt.fields.cppSecondAdditionalEmployer,
				cppSecondAdditionalSelfEmployed: tt.fields.cppSecondAdditionalSelfEmployed,
			}
			_ = cpp.Calculate(tt.args.totalIncome)
			if got := cpp.GetCPPBasicSelfEmployed(); got != tt.want {
				t.Errorf("GetCPPBasicSelfEmployed() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCanadaPensionPlan_GetCPPFirstAdditionalEmployee(t *testing.T) {
	type fields struct {
		Year                            int
		BasicRateEmployee               float64
		BasicRateEmployer               float64
		FirstAdditionalRateEmployee     float64
		FirstAdditionalRateEmployer     float64
		SecondAdditionalRateEmployee    float64
		SecondAdditionalRateEmployer    float64
		BasicExemption                  float64
		MaxPensionableEarning           float64
		AdditionalMaxPensionableEarning float64
		cppBasicEmployee                float64
		cppBasicEmployer                float64
		cppBasicSelfEmployed            float64
		cppFirstAdditionalEmployee      float64
		cppFirstAdditionalEmployer      float64
		cppFirstAdditionalSelfEmployed  float64
		cppSecondAdditionalEmployee     float64
		cppSecondAdditionalEmployer     float64
		cppSecondAdditionalSelfEmployed float64
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
		{name: "zero income", fields: fields{Year: 2023, BasicExemption: 3500, BasicRateEmployee: 4.95, BasicRateEmployer: 4.95, FirstAdditionalRateEmployee: 1, FirstAdditionalRateEmployer: 1, SecondAdditionalRateEmployee: 0, SecondAdditionalRateEmployer: 0, MaxPensionableEarning: 66600}, args: args{totalIncome: 0}, want: 0},
		{name: "low edge income", fields: fields{Year: 2023, BasicExemption: 3500, BasicRateEmployee: 4.95, BasicRateEmployer: 4.95, FirstAdditionalRateEmployee: 1, FirstAdditionalRateEmployer: 1, SecondAdditionalRateEmployee: 0, SecondAdditionalRateEmployer: 0, MaxPensionableEarning: 66600}, args: args{totalIncome: 3500}, want: 0},
		{name: "low income 1", fields: fields{Year: 2023, BasicExemption: 3500, BasicRateEmployee: 4.95, BasicRateEmployer: 4.95, FirstAdditionalRateEmployee: 1, FirstAdditionalRateEmployer: 1, SecondAdditionalRateEmployee: 0, SecondAdditionalRateEmployer: 0, MaxPensionableEarning: 66600}, args: args{totalIncome: 10000}, want: 65},
		{name: "low income 2", fields: fields{Year: 2023, BasicExemption: 3500, BasicRateEmployee: 4.95, BasicRateEmployer: 4.95, FirstAdditionalRateEmployee: 1, FirstAdditionalRateEmployer: 1, SecondAdditionalRateEmployee: 0, SecondAdditionalRateEmployer: 0, MaxPensionableEarning: 66600}, args: args{totalIncome: 60000}, want: 565},
		{name: "edge income", fields: fields{Year: 2023, BasicExemption: 3500, BasicRateEmployee: 4.95, BasicRateEmployer: 4.95, FirstAdditionalRateEmployee: 1, FirstAdditionalRateEmployer: 1, SecondAdditionalRateEmployee: 0, SecondAdditionalRateEmployer: 0, MaxPensionableEarning: 66600}, args: args{totalIncome: 66600}, want: 631},
		{name: "high income 1", fields: fields{Year: 2023, BasicExemption: 3500, BasicRateEmployee: 4.95, BasicRateEmployer: 4.95, FirstAdditionalRateEmployee: 1, FirstAdditionalRateEmployer: 1, SecondAdditionalRateEmployee: 0, SecondAdditionalRateEmployer: 0, MaxPensionableEarning: 66600}, args: args{totalIncome: 70000}, want: 631},
		{name: "high income 2", fields: fields{Year: 2023, BasicExemption: 3500, BasicRateEmployee: 4.95, BasicRateEmployer: 4.95, FirstAdditionalRateEmployee: 1, FirstAdditionalRateEmployer: 1, SecondAdditionalRateEmployee: 0, SecondAdditionalRateEmployer: 0, MaxPensionableEarning: 66600}, args: args{totalIncome: 500000}, want: 631},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cpp := &CanadaPensionPlan{
				Year:                            tt.fields.Year,
				BasicRateEmployee:               tt.fields.BasicRateEmployee,
				BasicRateEmployer:               tt.fields.BasicRateEmployer,
				FirstAdditionalRateEmployee:     tt.fields.FirstAdditionalRateEmployee,
				FirstAdditionalRateEmployer:     tt.fields.FirstAdditionalRateEmployer,
				SecondAdditionalRateEmployee:    tt.fields.SecondAdditionalRateEmployee,
				SecondAdditionalRateEmployer:    tt.fields.SecondAdditionalRateEmployer,
				BasicExemption:                  tt.fields.BasicExemption,
				MaxPensionableEarning:           tt.fields.MaxPensionableEarning,
				AdditionalMaxPensionableEarning: tt.fields.AdditionalMaxPensionableEarning,
				cppBasicEmployee:                tt.fields.cppBasicEmployee,
				cppBasicEmployer:                tt.fields.cppBasicEmployer,
				cppBasicSelfEmployed:            tt.fields.cppBasicSelfEmployed,
				cppFirstAdditionalEmployee:      tt.fields.cppFirstAdditionalEmployee,
				cppFirstAdditionalEmployer:      tt.fields.cppFirstAdditionalEmployer,
				cppFirstAdditionalSelfEmployed:  tt.fields.cppFirstAdditionalSelfEmployed,
				cppSecondAdditionalEmployee:     tt.fields.cppSecondAdditionalEmployee,
				cppSecondAdditionalEmployer:     tt.fields.cppSecondAdditionalEmployer,
				cppSecondAdditionalSelfEmployed: tt.fields.cppSecondAdditionalSelfEmployed,
			}
			_ = cpp.Calculate(tt.args.totalIncome)
			if got := cpp.GetCPPFirstAdditionalEmployee(); got != tt.want {
				t.Errorf("GetCPPFirstAdditionalEmployee() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCanadaPensionPlan_GetCPPFirstAdditionalEmployer(t *testing.T) {
	type fields struct {
		Year                            int
		BasicRateEmployee               float64
		BasicRateEmployer               float64
		FirstAdditionalRateEmployee     float64
		FirstAdditionalRateEmployer     float64
		SecondAdditionalRateEmployee    float64
		SecondAdditionalRateEmployer    float64
		BasicExemption                  float64
		MaxPensionableEarning           float64
		AdditionalMaxPensionableEarning float64
		cppBasicEmployee                float64
		cppBasicEmployer                float64
		cppBasicSelfEmployed            float64
		cppFirstAdditionalEmployee      float64
		cppFirstAdditionalEmployer      float64
		cppFirstAdditionalSelfEmployed  float64
		cppSecondAdditionalEmployee     float64
		cppSecondAdditionalEmployer     float64
		cppSecondAdditionalSelfEmployed float64
	}
	tests := []struct {
		name   string
		fields fields
		want   float64
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cpp := &CanadaPensionPlan{
				Year:                            tt.fields.Year,
				BasicRateEmployee:               tt.fields.BasicRateEmployee,
				BasicRateEmployer:               tt.fields.BasicRateEmployer,
				FirstAdditionalRateEmployee:     tt.fields.FirstAdditionalRateEmployee,
				FirstAdditionalRateEmployer:     tt.fields.FirstAdditionalRateEmployer,
				SecondAdditionalRateEmployee:    tt.fields.SecondAdditionalRateEmployee,
				SecondAdditionalRateEmployer:    tt.fields.SecondAdditionalRateEmployer,
				BasicExemption:                  tt.fields.BasicExemption,
				MaxPensionableEarning:           tt.fields.MaxPensionableEarning,
				AdditionalMaxPensionableEarning: tt.fields.AdditionalMaxPensionableEarning,
				cppBasicEmployee:                tt.fields.cppBasicEmployee,
				cppBasicEmployer:                tt.fields.cppBasicEmployer,
				cppBasicSelfEmployed:            tt.fields.cppBasicSelfEmployed,
				cppFirstAdditionalEmployee:      tt.fields.cppFirstAdditionalEmployee,
				cppFirstAdditionalEmployer:      tt.fields.cppFirstAdditionalEmployer,
				cppFirstAdditionalSelfEmployed:  tt.fields.cppFirstAdditionalSelfEmployed,
				cppSecondAdditionalEmployee:     tt.fields.cppSecondAdditionalEmployee,
				cppSecondAdditionalEmployer:     tt.fields.cppSecondAdditionalEmployer,
				cppSecondAdditionalSelfEmployed: tt.fields.cppSecondAdditionalSelfEmployed,
			}
			if got := cpp.GetCPPFirstAdditionalEmployer(); got != tt.want {
				t.Errorf("GetCPPFirstAdditionalEmployer() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCanadaPensionPlan_GetCPPFirstAdditionalSelfEmployed(t *testing.T) {
	type fields struct {
		Year                            int
		BasicRateEmployee               float64
		BasicRateEmployer               float64
		FirstAdditionalRateEmployee     float64
		FirstAdditionalRateEmployer     float64
		SecondAdditionalRateEmployee    float64
		SecondAdditionalRateEmployer    float64
		BasicExemption                  float64
		MaxPensionableEarning           float64
		AdditionalMaxPensionableEarning float64
		cppBasicEmployee                float64
		cppBasicEmployer                float64
		cppBasicSelfEmployed            float64
		cppFirstAdditionalEmployee      float64
		cppFirstAdditionalEmployer      float64
		cppFirstAdditionalSelfEmployed  float64
		cppSecondAdditionalEmployee     float64
		cppSecondAdditionalEmployer     float64
		cppSecondAdditionalSelfEmployed float64
	}
	tests := []struct {
		name   string
		fields fields
		want   float64
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cpp := &CanadaPensionPlan{
				Year:                            tt.fields.Year,
				BasicRateEmployee:               tt.fields.BasicRateEmployee,
				BasicRateEmployer:               tt.fields.BasicRateEmployer,
				FirstAdditionalRateEmployee:     tt.fields.FirstAdditionalRateEmployee,
				FirstAdditionalRateEmployer:     tt.fields.FirstAdditionalRateEmployer,
				SecondAdditionalRateEmployee:    tt.fields.SecondAdditionalRateEmployee,
				SecondAdditionalRateEmployer:    tt.fields.SecondAdditionalRateEmployer,
				BasicExemption:                  tt.fields.BasicExemption,
				MaxPensionableEarning:           tt.fields.MaxPensionableEarning,
				AdditionalMaxPensionableEarning: tt.fields.AdditionalMaxPensionableEarning,
				cppBasicEmployee:                tt.fields.cppBasicEmployee,
				cppBasicEmployer:                tt.fields.cppBasicEmployer,
				cppBasicSelfEmployed:            tt.fields.cppBasicSelfEmployed,
				cppFirstAdditionalEmployee:      tt.fields.cppFirstAdditionalEmployee,
				cppFirstAdditionalEmployer:      tt.fields.cppFirstAdditionalEmployer,
				cppFirstAdditionalSelfEmployed:  tt.fields.cppFirstAdditionalSelfEmployed,
				cppSecondAdditionalEmployee:     tt.fields.cppSecondAdditionalEmployee,
				cppSecondAdditionalEmployer:     tt.fields.cppSecondAdditionalEmployer,
				cppSecondAdditionalSelfEmployed: tt.fields.cppSecondAdditionalSelfEmployed,
			}
			if got := cpp.GetCPPFirstAdditionalSelfEmployed(); got != tt.want {
				t.Errorf("GetCPPFirstAdditionalSelfEmployed() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCanadaPensionPlan_GetCPPSecondAdditionalEmployee(t *testing.T) {
	type fields struct {
		Year                            int
		BasicRateEmployee               float64
		BasicRateEmployer               float64
		FirstAdditionalRateEmployee     float64
		FirstAdditionalRateEmployer     float64
		SecondAdditionalRateEmployee    float64
		SecondAdditionalRateEmployer    float64
		BasicExemption                  float64
		MaxPensionableEarning           float64
		AdditionalMaxPensionableEarning float64
		cppBasicEmployee                float64
		cppBasicEmployer                float64
		cppBasicSelfEmployed            float64
		cppFirstAdditionalEmployee      float64
		cppFirstAdditionalEmployer      float64
		cppFirstAdditionalSelfEmployed  float64
		cppSecondAdditionalEmployee     float64
		cppSecondAdditionalEmployer     float64
		cppSecondAdditionalSelfEmployed float64
	}
	tests := []struct {
		name   string
		fields fields
		want   float64
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cpp := &CanadaPensionPlan{
				Year:                            tt.fields.Year,
				BasicRateEmployee:               tt.fields.BasicRateEmployee,
				BasicRateEmployer:               tt.fields.BasicRateEmployer,
				FirstAdditionalRateEmployee:     tt.fields.FirstAdditionalRateEmployee,
				FirstAdditionalRateEmployer:     tt.fields.FirstAdditionalRateEmployer,
				SecondAdditionalRateEmployee:    tt.fields.SecondAdditionalRateEmployee,
				SecondAdditionalRateEmployer:    tt.fields.SecondAdditionalRateEmployer,
				BasicExemption:                  tt.fields.BasicExemption,
				MaxPensionableEarning:           tt.fields.MaxPensionableEarning,
				AdditionalMaxPensionableEarning: tt.fields.AdditionalMaxPensionableEarning,
				cppBasicEmployee:                tt.fields.cppBasicEmployee,
				cppBasicEmployer:                tt.fields.cppBasicEmployer,
				cppBasicSelfEmployed:            tt.fields.cppBasicSelfEmployed,
				cppFirstAdditionalEmployee:      tt.fields.cppFirstAdditionalEmployee,
				cppFirstAdditionalEmployer:      tt.fields.cppFirstAdditionalEmployer,
				cppFirstAdditionalSelfEmployed:  tt.fields.cppFirstAdditionalSelfEmployed,
				cppSecondAdditionalEmployee:     tt.fields.cppSecondAdditionalEmployee,
				cppSecondAdditionalEmployer:     tt.fields.cppSecondAdditionalEmployer,
				cppSecondAdditionalSelfEmployed: tt.fields.cppSecondAdditionalSelfEmployed,
			}
			if got := cpp.GetCPPSecondAdditionalEmployee(); got != tt.want {
				t.Errorf("GetCPPSecondAdditionalEmployee() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCanadaPensionPlan_GetCPPSecondAdditionalEmployer(t *testing.T) {
	type fields struct {
		Year                            int
		BasicRateEmployee               float64
		BasicRateEmployer               float64
		FirstAdditionalRateEmployee     float64
		FirstAdditionalRateEmployer     float64
		SecondAdditionalRateEmployee    float64
		SecondAdditionalRateEmployer    float64
		BasicExemption                  float64
		MaxPensionableEarning           float64
		AdditionalMaxPensionableEarning float64
		cppBasicEmployee                float64
		cppBasicEmployer                float64
		cppBasicSelfEmployed            float64
		cppFirstAdditionalEmployee      float64
		cppFirstAdditionalEmployer      float64
		cppFirstAdditionalSelfEmployed  float64
		cppSecondAdditionalEmployee     float64
		cppSecondAdditionalEmployer     float64
		cppSecondAdditionalSelfEmployed float64
	}
	tests := []struct {
		name   string
		fields fields
		want   float64
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cpp := &CanadaPensionPlan{
				Year:                            tt.fields.Year,
				BasicRateEmployee:               tt.fields.BasicRateEmployee,
				BasicRateEmployer:               tt.fields.BasicRateEmployer,
				FirstAdditionalRateEmployee:     tt.fields.FirstAdditionalRateEmployee,
				FirstAdditionalRateEmployer:     tt.fields.FirstAdditionalRateEmployer,
				SecondAdditionalRateEmployee:    tt.fields.SecondAdditionalRateEmployee,
				SecondAdditionalRateEmployer:    tt.fields.SecondAdditionalRateEmployer,
				BasicExemption:                  tt.fields.BasicExemption,
				MaxPensionableEarning:           tt.fields.MaxPensionableEarning,
				AdditionalMaxPensionableEarning: tt.fields.AdditionalMaxPensionableEarning,
				cppBasicEmployee:                tt.fields.cppBasicEmployee,
				cppBasicEmployer:                tt.fields.cppBasicEmployer,
				cppBasicSelfEmployed:            tt.fields.cppBasicSelfEmployed,
				cppFirstAdditionalEmployee:      tt.fields.cppFirstAdditionalEmployee,
				cppFirstAdditionalEmployer:      tt.fields.cppFirstAdditionalEmployer,
				cppFirstAdditionalSelfEmployed:  tt.fields.cppFirstAdditionalSelfEmployed,
				cppSecondAdditionalEmployee:     tt.fields.cppSecondAdditionalEmployee,
				cppSecondAdditionalEmployer:     tt.fields.cppSecondAdditionalEmployer,
				cppSecondAdditionalSelfEmployed: tt.fields.cppSecondAdditionalSelfEmployed,
			}
			if got := cpp.GetCPPSecondAdditionalEmployer(); got != tt.want {
				t.Errorf("GetCPPSecondAdditionalEmployer() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCanadaPensionPlan_GetCPPSecondAdditionalSelfEmployed(t *testing.T) {
	type fields struct {
		Year                            int
		BasicRateEmployee               float64
		BasicRateEmployer               float64
		FirstAdditionalRateEmployee     float64
		FirstAdditionalRateEmployer     float64
		SecondAdditionalRateEmployee    float64
		SecondAdditionalRateEmployer    float64
		BasicExemption                  float64
		MaxPensionableEarning           float64
		AdditionalMaxPensionableEarning float64
		cppBasicEmployee                float64
		cppBasicEmployer                float64
		cppBasicSelfEmployed            float64
		cppFirstAdditionalEmployee      float64
		cppFirstAdditionalEmployer      float64
		cppFirstAdditionalSelfEmployed  float64
		cppSecondAdditionalEmployee     float64
		cppSecondAdditionalEmployer     float64
		cppSecondAdditionalSelfEmployed float64
	}
	tests := []struct {
		name   string
		fields fields
		want   float64
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cpp := &CanadaPensionPlan{
				Year:                            tt.fields.Year,
				BasicRateEmployee:               tt.fields.BasicRateEmployee,
				BasicRateEmployer:               tt.fields.BasicRateEmployer,
				FirstAdditionalRateEmployee:     tt.fields.FirstAdditionalRateEmployee,
				FirstAdditionalRateEmployer:     tt.fields.FirstAdditionalRateEmployer,
				SecondAdditionalRateEmployee:    tt.fields.SecondAdditionalRateEmployee,
				SecondAdditionalRateEmployer:    tt.fields.SecondAdditionalRateEmployer,
				BasicExemption:                  tt.fields.BasicExemption,
				MaxPensionableEarning:           tt.fields.MaxPensionableEarning,
				AdditionalMaxPensionableEarning: tt.fields.AdditionalMaxPensionableEarning,
				cppBasicEmployee:                tt.fields.cppBasicEmployee,
				cppBasicEmployer:                tt.fields.cppBasicEmployer,
				cppBasicSelfEmployed:            tt.fields.cppBasicSelfEmployed,
				cppFirstAdditionalEmployee:      tt.fields.cppFirstAdditionalEmployee,
				cppFirstAdditionalEmployer:      tt.fields.cppFirstAdditionalEmployer,
				cppFirstAdditionalSelfEmployed:  tt.fields.cppFirstAdditionalSelfEmployed,
				cppSecondAdditionalEmployee:     tt.fields.cppSecondAdditionalEmployee,
				cppSecondAdditionalEmployer:     tt.fields.cppSecondAdditionalEmployer,
				cppSecondAdditionalSelfEmployed: tt.fields.cppSecondAdditionalSelfEmployed,
			}
			if got := cpp.GetCPPSecondAdditionalSelfEmployed(); got != tt.want {
				t.Errorf("GetCPPSecondAdditionalSelfEmployed() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCanadaPensionPlan_calculateCppBasic(t *testing.T) {
	type fields struct {
		Year                            int
		BasicRateEmployee               float64
		BasicRateEmployer               float64
		FirstAdditionalRateEmployee     float64
		FirstAdditionalRateEmployer     float64
		SecondAdditionalRateEmployee    float64
		SecondAdditionalRateEmployer    float64
		BasicExemption                  float64
		MaxPensionableEarning           float64
		AdditionalMaxPensionableEarning float64
		cppBasicEmployee                float64
		cppBasicEmployer                float64
		cppBasicSelfEmployed            float64
		cppFirstAdditionalEmployee      float64
		cppFirstAdditionalEmployer      float64
		cppFirstAdditionalSelfEmployed  float64
		cppSecondAdditionalEmployee     float64
		cppSecondAdditionalEmployer     float64
		cppSecondAdditionalSelfEmployed float64
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
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cpp := &CanadaPensionPlan{
				Year:                            tt.fields.Year,
				BasicRateEmployee:               tt.fields.BasicRateEmployee,
				BasicRateEmployer:               tt.fields.BasicRateEmployer,
				FirstAdditionalRateEmployee:     tt.fields.FirstAdditionalRateEmployee,
				FirstAdditionalRateEmployer:     tt.fields.FirstAdditionalRateEmployer,
				SecondAdditionalRateEmployee:    tt.fields.SecondAdditionalRateEmployee,
				SecondAdditionalRateEmployer:    tt.fields.SecondAdditionalRateEmployer,
				BasicExemption:                  tt.fields.BasicExemption,
				MaxPensionableEarning:           tt.fields.MaxPensionableEarning,
				AdditionalMaxPensionableEarning: tt.fields.AdditionalMaxPensionableEarning,
				cppBasicEmployee:                tt.fields.cppBasicEmployee,
				cppBasicEmployer:                tt.fields.cppBasicEmployer,
				cppBasicSelfEmployed:            tt.fields.cppBasicSelfEmployed,
				cppFirstAdditionalEmployee:      tt.fields.cppFirstAdditionalEmployee,
				cppFirstAdditionalEmployer:      tt.fields.cppFirstAdditionalEmployer,
				cppFirstAdditionalSelfEmployed:  tt.fields.cppFirstAdditionalSelfEmployed,
				cppSecondAdditionalEmployee:     tt.fields.cppSecondAdditionalEmployee,
				cppSecondAdditionalEmployer:     tt.fields.cppSecondAdditionalEmployer,
				cppSecondAdditionalSelfEmployed: tt.fields.cppSecondAdditionalSelfEmployed,
			}
			if err := cpp.calculateCppBasic(tt.args.totalIncome); (err != nil) != tt.wantErr {
				t.Errorf("calculateCppBasic() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestCanadaPensionPlan_calculateCppFirst(t *testing.T) {
	type fields struct {
		Year                            int
		BasicRateEmployee               float64
		BasicRateEmployer               float64
		FirstAdditionalRateEmployee     float64
		FirstAdditionalRateEmployer     float64
		SecondAdditionalRateEmployee    float64
		SecondAdditionalRateEmployer    float64
		BasicExemption                  float64
		MaxPensionableEarning           float64
		AdditionalMaxPensionableEarning float64
		cppBasicEmployee                float64
		cppBasicEmployer                float64
		cppBasicSelfEmployed            float64
		cppFirstAdditionalEmployee      float64
		cppFirstAdditionalEmployer      float64
		cppFirstAdditionalSelfEmployed  float64
		cppSecondAdditionalEmployee     float64
		cppSecondAdditionalEmployer     float64
		cppSecondAdditionalSelfEmployed float64
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
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cpp := &CanadaPensionPlan{
				Year:                            tt.fields.Year,
				BasicRateEmployee:               tt.fields.BasicRateEmployee,
				BasicRateEmployer:               tt.fields.BasicRateEmployer,
				FirstAdditionalRateEmployee:     tt.fields.FirstAdditionalRateEmployee,
				FirstAdditionalRateEmployer:     tt.fields.FirstAdditionalRateEmployer,
				SecondAdditionalRateEmployee:    tt.fields.SecondAdditionalRateEmployee,
				SecondAdditionalRateEmployer:    tt.fields.SecondAdditionalRateEmployer,
				BasicExemption:                  tt.fields.BasicExemption,
				MaxPensionableEarning:           tt.fields.MaxPensionableEarning,
				AdditionalMaxPensionableEarning: tt.fields.AdditionalMaxPensionableEarning,
				cppBasicEmployee:                tt.fields.cppBasicEmployee,
				cppBasicEmployer:                tt.fields.cppBasicEmployer,
				cppBasicSelfEmployed:            tt.fields.cppBasicSelfEmployed,
				cppFirstAdditionalEmployee:      tt.fields.cppFirstAdditionalEmployee,
				cppFirstAdditionalEmployer:      tt.fields.cppFirstAdditionalEmployer,
				cppFirstAdditionalSelfEmployed:  tt.fields.cppFirstAdditionalSelfEmployed,
				cppSecondAdditionalEmployee:     tt.fields.cppSecondAdditionalEmployee,
				cppSecondAdditionalEmployer:     tt.fields.cppSecondAdditionalEmployer,
				cppSecondAdditionalSelfEmployed: tt.fields.cppSecondAdditionalSelfEmployed,
			}
			if err := cpp.calculateCppFirst(tt.args.totalIncome); (err != nil) != tt.wantErr {
				t.Errorf("calculateCppFirst() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestCanadaPensionPlan_calculateCppSecond(t *testing.T) {
	type fields struct {
		Year                            int
		BasicRateEmployee               float64
		BasicRateEmployer               float64
		FirstAdditionalRateEmployee     float64
		FirstAdditionalRateEmployer     float64
		SecondAdditionalRateEmployee    float64
		SecondAdditionalRateEmployer    float64
		BasicExemption                  float64
		MaxPensionableEarning           float64
		AdditionalMaxPensionableEarning float64
		cppBasicEmployee                float64
		cppBasicEmployer                float64
		cppBasicSelfEmployed            float64
		cppFirstAdditionalEmployee      float64
		cppFirstAdditionalEmployer      float64
		cppFirstAdditionalSelfEmployed  float64
		cppSecondAdditionalEmployee     float64
		cppSecondAdditionalEmployer     float64
		cppSecondAdditionalSelfEmployed float64
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
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cpp := &CanadaPensionPlan{
				Year:                            tt.fields.Year,
				BasicRateEmployee:               tt.fields.BasicRateEmployee,
				BasicRateEmployer:               tt.fields.BasicRateEmployer,
				FirstAdditionalRateEmployee:     tt.fields.FirstAdditionalRateEmployee,
				FirstAdditionalRateEmployer:     tt.fields.FirstAdditionalRateEmployer,
				SecondAdditionalRateEmployee:    tt.fields.SecondAdditionalRateEmployee,
				SecondAdditionalRateEmployer:    tt.fields.SecondAdditionalRateEmployer,
				BasicExemption:                  tt.fields.BasicExemption,
				MaxPensionableEarning:           tt.fields.MaxPensionableEarning,
				AdditionalMaxPensionableEarning: tt.fields.AdditionalMaxPensionableEarning,
				cppBasicEmployee:                tt.fields.cppBasicEmployee,
				cppBasicEmployer:                tt.fields.cppBasicEmployer,
				cppBasicSelfEmployed:            tt.fields.cppBasicSelfEmployed,
				cppFirstAdditionalEmployee:      tt.fields.cppFirstAdditionalEmployee,
				cppFirstAdditionalEmployer:      tt.fields.cppFirstAdditionalEmployer,
				cppFirstAdditionalSelfEmployed:  tt.fields.cppFirstAdditionalSelfEmployed,
				cppSecondAdditionalEmployee:     tt.fields.cppSecondAdditionalEmployee,
				cppSecondAdditionalEmployer:     tt.fields.cppSecondAdditionalEmployer,
				cppSecondAdditionalSelfEmployed: tt.fields.cppSecondAdditionalSelfEmployed,
			}
			if err := cpp.calculateCppSecond(tt.args.totalIncome); (err != nil) != tt.wantErr {
				t.Errorf("calculateCppSecond() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestCanadaPensionPlan_validateCppBasicInputs(t *testing.T) {
	type fields struct {
		Year                            int
		BasicRateEmployee               float64
		BasicRateEmployer               float64
		FirstAdditionalRateEmployee     float64
		FirstAdditionalRateEmployer     float64
		SecondAdditionalRateEmployee    float64
		SecondAdditionalRateEmployer    float64
		BasicExemption                  float64
		MaxPensionableEarning           float64
		AdditionalMaxPensionableEarning float64
		cppBasicEmployee                float64
		cppBasicEmployer                float64
		cppBasicSelfEmployed            float64
		cppFirstAdditionalEmployee      float64
		cppFirstAdditionalEmployer      float64
		cppFirstAdditionalSelfEmployed  float64
		cppSecondAdditionalEmployee     float64
		cppSecondAdditionalEmployer     float64
		cppSecondAdditionalSelfEmployed float64
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cpp := &CanadaPensionPlan{
				Year:                            tt.fields.Year,
				BasicRateEmployee:               tt.fields.BasicRateEmployee,
				BasicRateEmployer:               tt.fields.BasicRateEmployer,
				FirstAdditionalRateEmployee:     tt.fields.FirstAdditionalRateEmployee,
				FirstAdditionalRateEmployer:     tt.fields.FirstAdditionalRateEmployer,
				SecondAdditionalRateEmployee:    tt.fields.SecondAdditionalRateEmployee,
				SecondAdditionalRateEmployer:    tt.fields.SecondAdditionalRateEmployer,
				BasicExemption:                  tt.fields.BasicExemption,
				MaxPensionableEarning:           tt.fields.MaxPensionableEarning,
				AdditionalMaxPensionableEarning: tt.fields.AdditionalMaxPensionableEarning,
				cppBasicEmployee:                tt.fields.cppBasicEmployee,
				cppBasicEmployer:                tt.fields.cppBasicEmployer,
				cppBasicSelfEmployed:            tt.fields.cppBasicSelfEmployed,
				cppFirstAdditionalEmployee:      tt.fields.cppFirstAdditionalEmployee,
				cppFirstAdditionalEmployer:      tt.fields.cppFirstAdditionalEmployer,
				cppFirstAdditionalSelfEmployed:  tt.fields.cppFirstAdditionalSelfEmployed,
				cppSecondAdditionalEmployee:     tt.fields.cppSecondAdditionalEmployee,
				cppSecondAdditionalEmployer:     tt.fields.cppSecondAdditionalEmployer,
				cppSecondAdditionalSelfEmployed: tt.fields.cppSecondAdditionalSelfEmployed,
			}
			if err := cpp.validateCppBasicInputs(); (err != nil) != tt.wantErr {
				t.Errorf("validateCppBasicInputs() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestCanadaPensionPlan_validateCppFirstInputs(t *testing.T) {
	type fields struct {
		Year                            int
		BasicRateEmployee               float64
		BasicRateEmployer               float64
		FirstAdditionalRateEmployee     float64
		FirstAdditionalRateEmployer     float64
		SecondAdditionalRateEmployee    float64
		SecondAdditionalRateEmployer    float64
		BasicExemption                  float64
		MaxPensionableEarning           float64
		AdditionalMaxPensionableEarning float64
		cppBasicEmployee                float64
		cppBasicEmployer                float64
		cppBasicSelfEmployed            float64
		cppFirstAdditionalEmployee      float64
		cppFirstAdditionalEmployer      float64
		cppFirstAdditionalSelfEmployed  float64
		cppSecondAdditionalEmployee     float64
		cppSecondAdditionalEmployer     float64
		cppSecondAdditionalSelfEmployed float64
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cpp := &CanadaPensionPlan{
				Year:                            tt.fields.Year,
				BasicRateEmployee:               tt.fields.BasicRateEmployee,
				BasicRateEmployer:               tt.fields.BasicRateEmployer,
				FirstAdditionalRateEmployee:     tt.fields.FirstAdditionalRateEmployee,
				FirstAdditionalRateEmployer:     tt.fields.FirstAdditionalRateEmployer,
				SecondAdditionalRateEmployee:    tt.fields.SecondAdditionalRateEmployee,
				SecondAdditionalRateEmployer:    tt.fields.SecondAdditionalRateEmployer,
				BasicExemption:                  tt.fields.BasicExemption,
				MaxPensionableEarning:           tt.fields.MaxPensionableEarning,
				AdditionalMaxPensionableEarning: tt.fields.AdditionalMaxPensionableEarning,
				cppBasicEmployee:                tt.fields.cppBasicEmployee,
				cppBasicEmployer:                tt.fields.cppBasicEmployer,
				cppBasicSelfEmployed:            tt.fields.cppBasicSelfEmployed,
				cppFirstAdditionalEmployee:      tt.fields.cppFirstAdditionalEmployee,
				cppFirstAdditionalEmployer:      tt.fields.cppFirstAdditionalEmployer,
				cppFirstAdditionalSelfEmployed:  tt.fields.cppFirstAdditionalSelfEmployed,
				cppSecondAdditionalEmployee:     tt.fields.cppSecondAdditionalEmployee,
				cppSecondAdditionalEmployer:     tt.fields.cppSecondAdditionalEmployer,
				cppSecondAdditionalSelfEmployed: tt.fields.cppSecondAdditionalSelfEmployed,
			}
			if err := cpp.validateCppFirstInputs(); (err != nil) != tt.wantErr {
				t.Errorf("validateCppFirstInputs() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestCanadaPensionPlan_validateCppSecondInputs(t *testing.T) {
	type fields struct {
		Year                            int
		BasicRateEmployee               float64
		BasicRateEmployer               float64
		FirstAdditionalRateEmployee     float64
		FirstAdditionalRateEmployer     float64
		SecondAdditionalRateEmployee    float64
		SecondAdditionalRateEmployer    float64
		BasicExemption                  float64
		MaxPensionableEarning           float64
		AdditionalMaxPensionableEarning float64
		cppBasicEmployee                float64
		cppBasicEmployer                float64
		cppBasicSelfEmployed            float64
		cppFirstAdditionalEmployee      float64
		cppFirstAdditionalEmployer      float64
		cppFirstAdditionalSelfEmployed  float64
		cppSecondAdditionalEmployee     float64
		cppSecondAdditionalEmployer     float64
		cppSecondAdditionalSelfEmployed float64
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cpp := &CanadaPensionPlan{
				Year:                            tt.fields.Year,
				BasicRateEmployee:               tt.fields.BasicRateEmployee,
				BasicRateEmployer:               tt.fields.BasicRateEmployer,
				FirstAdditionalRateEmployee:     tt.fields.FirstAdditionalRateEmployee,
				FirstAdditionalRateEmployer:     tt.fields.FirstAdditionalRateEmployer,
				SecondAdditionalRateEmployee:    tt.fields.SecondAdditionalRateEmployee,
				SecondAdditionalRateEmployer:    tt.fields.SecondAdditionalRateEmployer,
				BasicExemption:                  tt.fields.BasicExemption,
				MaxPensionableEarning:           tt.fields.MaxPensionableEarning,
				AdditionalMaxPensionableEarning: tt.fields.AdditionalMaxPensionableEarning,
				cppBasicEmployee:                tt.fields.cppBasicEmployee,
				cppBasicEmployer:                tt.fields.cppBasicEmployer,
				cppBasicSelfEmployed:            tt.fields.cppBasicSelfEmployed,
				cppFirstAdditionalEmployee:      tt.fields.cppFirstAdditionalEmployee,
				cppFirstAdditionalEmployer:      tt.fields.cppFirstAdditionalEmployer,
				cppFirstAdditionalSelfEmployed:  tt.fields.cppFirstAdditionalSelfEmployed,
				cppSecondAdditionalEmployee:     tt.fields.cppSecondAdditionalEmployee,
				cppSecondAdditionalEmployer:     tt.fields.cppSecondAdditionalEmployer,
				cppSecondAdditionalSelfEmployed: tt.fields.cppSecondAdditionalSelfEmployed,
			}
			if err := cpp.validateCppSecondInputs(); (err != nil) != tt.wantErr {
				t.Errorf("validateCppSecondInputs() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
