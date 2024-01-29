package shared

import (
	"math"
	"reflect"
	"testing"
)

func TestFromArray(t *testing.T) {
	type args struct {
		brackets []float64
		rates    []float64
	}
	tests := []struct {
		name    string
		args    args
		want    []TaxBracket
		wantErr bool
	}{
		{name: "valid inputs", args: args{brackets: []float64{0, 50000, 100000}, rates: []float64{5, 10, 20}}, want: []TaxBracket{{Low: 0, High: 50000, Rate: 5}, {Low: 50000, High: 100000, Rate: 10}, {Low: 100000, High: math.MaxFloat64, Rate: 20}}, wantErr: false},
		{name: "invalid brackets length", args: args{brackets: []float64{0, 50000}, rates: []float64{5, 10, 20}}, want: []TaxBracket{}, wantErr: true},
		{name: "invalid rates length", args: args{brackets: []float64{0, 50000, 100000}, rates: []float64{5, 10}}, want: []TaxBracket{}, wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := FromArray(tt.args.brackets, tt.args.rates)
			if (err != nil) != tt.wantErr {
				t.Errorf("FromArray() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("FromArray() got = %v, want %v", got, tt.want)
			}
		})
	}
}
