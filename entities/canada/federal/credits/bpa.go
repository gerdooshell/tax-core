package credits

import (
	"fmt"
	"github.com/gerdooshell/tax-core/library/mathHelper"
)

type BasicPersonalAmount struct {
	MaxBPAIncome float64
	MinBPAIncome float64
	MinBPAAmount float64
	MaxBPAAmount float64
	value        float64
}

func (bpa *BasicPersonalAmount) Calculate(income float64) error {
	err := bpa.validate(income)
	if err != nil {
		return err
	}

	if income >= bpa.MinBPAIncome {
		bpa.value = bpa.MinBPAAmount
	} else if income > bpa.MaxBPAIncome && income < bpa.MinBPAIncome {
		bpa.value = bpa.MaxBPAAmount -
			((income - bpa.MaxBPAIncome) *
				((bpa.MaxBPAAmount - bpa.MinBPAAmount) / (bpa.MinBPAIncome - bpa.MaxBPAIncome)))
	} else {
		bpa.value = bpa.MaxBPAAmount
	}
	bpa.value = mathHelper.RoundFloat64(bpa.value, 2)
	return nil
}

func (bpa *BasicPersonalAmount) GetValue() float64 {
	return bpa.value
}

func (bpa *BasicPersonalAmount) validate(income float64) error {
	if income < 0 {
		return fmt.Errorf("bpaf error: invalid income: \"%v\"", income)
	}
	if bpa.MaxBPAIncome <= 0 {
		return fmt.Errorf("bpaf error: invalid MaxBPAIncome: \"%v\"", bpa.MaxBPAIncome)
	}
	if bpa.MinBPAIncome <= 0 {
		return fmt.Errorf("bpaf error: invalid MinBPAIncome: \"%v\"", bpa.MinBPAIncome)
	}
	if bpa.MinBPAAmount <= 0 {
		return fmt.Errorf("bpaf error: invalid MinBPAAmount: \"%v\"", bpa.MinBPAAmount)
	}
	if bpa.MaxBPAAmount <= 0 {
		return fmt.Errorf("bpaf error: invalid MaxBPAAmount: \"%v\"", bpa.MaxBPAAmount)
	}
	return nil
}
