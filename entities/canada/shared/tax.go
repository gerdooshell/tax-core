package shared

import (
	"fmt"
	"github.com/gerdooshell/tax-core/library/mathHelper"
	"sort"
)

type Tax struct {
	TaxBrackets      []TaxBracket
	calculatedAmount float64
}

func (t *Tax) Calculate(taxableIncome float64) error {
	if err := t.validateParameters(taxableIncome); err != nil {
		return err
	}
	sort.Slice(t.TaxBrackets[:], func(i, j int) bool {
		return t.TaxBrackets[i].Low < t.TaxBrackets[j].Low
	})

	t.calculatedAmount = 0
	for _, element := range t.TaxBrackets {
		if taxableIncome >= element.Low && taxableIncome < element.High {
			t.calculatedAmount += (taxableIncome - element.Low) * element.Rate / 100
			break
		} else {
			t.calculatedAmount += (element.High - element.Low) * element.Rate / 100
		}
	}
	t.calculatedAmount = mathHelper.RoundFloat64(t.calculatedAmount, 2)
	return nil
}

func (t *Tax) GetValue() float64 {
	return t.calculatedAmount
}

func (t *Tax) validateParameters(taxableIncome float64) error {
	if taxableIncome < 0 {
		return fmt.Errorf("tax error: invalid taxable income: \"%v\"", taxableIncome)
	}
	if len(t.TaxBrackets) == 0 {
		return fmt.Errorf("invalid tax bracket length: \"%v\"", t.TaxBrackets)
	}
	return nil
}
