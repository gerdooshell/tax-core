package credits

import "testing"

func TestBasicPensionAmount_Calculate(t *testing.T) {
	type fields struct {
		Value float64
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{name: "zero bpabc", fields: fields{Value: 0}, wantErr: true},
		{name: "negative bpabc", fields: fields{Value: -1}, wantErr: true},
		{name: "positive bpabc", fields: fields{Value: 1}, wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bpa := &BasicPersonalAmount{
				Value: tt.fields.Value,
			}
			if err := bpa.Calculate(); (err != nil) != tt.wantErr {
				t.Errorf("Calculate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestBasicPensionAmount_GetValue(t *testing.T) {
	type fields struct {
		Value float64
	}
	tests := []struct {
		name   string
		fields fields
		want   float64
	}{
		{name: "bpabc 10", fields: fields{Value: 10}, want: 10},
		{name: "bpabc 100", fields: fields{Value: 100}, want: 100},
		{name: "bpabc 100000", fields: fields{Value: 100000}, want: 100000},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bpa := &BasicPersonalAmount{
				Value: tt.fields.Value,
			}
			if got := bpa.GetValue(); got != tt.want {
				t.Errorf("GetValue() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBasicPensionAmount_validate(t *testing.T) {
	type fields struct {
		Value float64
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{name: "zero bpabc", fields: fields{Value: 0}, wantErr: true},
		{name: "negative bpabc", fields: fields{Value: -1}, wantErr: true},
		{name: "positive bpabc", fields: fields{Value: 1}, wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bpa := &BasicPersonalAmount{
				Value: tt.fields.Value,
			}
			if err := bpa.validate(); (err != nil) != tt.wantErr {
				t.Errorf("validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
