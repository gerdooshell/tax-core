package credits

import (
	"fmt"
	"github.com/gerdooshell/tax-core/library/mathHelper"
)

type CanadaEmploymentAmount struct {
	Value float64
}

func (cea *CanadaEmploymentAmount) Calculate() error {
	if err := cea.validate(); err != nil {
		return err
	}
	cea.Value = mathHelper.RoundFloat64(cea.Value, 2)
	return nil
}

func (cea *CanadaEmploymentAmount) GetValue() float64 {
	return cea.Value
}

func (cea *CanadaEmploymentAmount) validate() error {
	if cea.Value < 0 {
		return fmt.Errorf("cea error: invalid cea: \"%v\"", cea.Value)
	}
	return nil
}
