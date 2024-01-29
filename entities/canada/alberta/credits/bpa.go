package credits

import (
	"fmt"
	"github.com/gerdooshell/tax-core/library/mathHelper"
)

type BasicPersonalAmount struct {
	Value float64
}

func (bpa *BasicPersonalAmount) Calculate() error {
	if err := bpa.validate(); err != nil {
		return err
	}
	bpa.Value = mathHelper.RoundFloat64(bpa.Value, 2)
	return nil
}

func (bpa *BasicPersonalAmount) GetValue() float64 {
	return bpa.Value
}

func (bpa *BasicPersonalAmount) validate() error {
	if bpa.Value <= 0 {
		return fmt.Errorf("bpaab error: invalid bpaab: \"%v\"", bpa.Value)
	}
	return nil
}
