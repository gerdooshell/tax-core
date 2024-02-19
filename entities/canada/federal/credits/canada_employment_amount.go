package credits

import (
	"fmt"
	"github.com/gerdooshell/tax-core/library/mathHelper"
)

type CanadaEmploymentAmount struct {
	Value           float64
	calculatedValue float64
}

func (cea *CanadaEmploymentAmount) Calculate(totalIncome float64) error {
	if err := cea.validate(totalIncome); err != nil {
		return err
	}
	cea.calculatedValue = mathHelper.RoundFloat64(min(cea.Value, totalIncome), 2)
	return nil
}

func (cea *CanadaEmploymentAmount) GetEmployeeValue() float64 {
	return cea.calculatedValue
}

func (cea *CanadaEmploymentAmount) GetSelfEmployedValue() float64 {
	return 0
}

func (cea *CanadaEmploymentAmount) validate(totalIncome float64) error {
	if totalIncome < 0 {
		return fmt.Errorf("cea error: invlid total income: \"%v\"", totalIncome)
	}
	if cea.Value < 0 {
		return fmt.Errorf("cea error: invalid cea: \"%v\"", cea.Value)
	}
	return nil
}
