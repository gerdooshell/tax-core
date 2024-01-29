package credits

import "testing"

func TestCanadaEmploymentAmount_Calculate(t *testing.T) {
	type fields struct {
		Value float64
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{name: "CAE for 2023", fields: fields{Value: 1368}, wantErr: false},
		{name: "negative CAE", fields: fields{Value: -2}, wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cea := &CanadaEmploymentAmount{
				Value: tt.fields.Value,
			}
			if err := cea.Calculate(); (err != nil) != tt.wantErr {
				t.Errorf("Calculate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestCanadaEmploymentAmount_GetValue(t *testing.T) {
	type fields struct {
		Value float64
	}
	tests := []struct {
		name   string
		fields fields
		want   float64
	}{
		{name: "get value", fields: fields{Value: 1368}, want: 1368},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cea := &CanadaEmploymentAmount{
				Value: tt.fields.Value,
			}
			if got := cea.GetValue(); got != tt.want {
				t.Errorf("GetValue() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCanadaEmploymentAmount_validate(t *testing.T) {
	type fields struct {
		Value float64
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{name: "CAE for 2023", fields: fields{Value: 1368}, wantErr: false},
		{name: "negative CAE", fields: fields{Value: -2}, wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cea := &CanadaEmploymentAmount{
				Value: tt.fields.Value,
			}
			if err := cea.validate(); (err != nil) != tt.wantErr {
				t.Errorf("validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
