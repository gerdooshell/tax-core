package shared

import (
	"fmt"
	"github.com/gerdooshell/tax-core/library/mathHelper"
)

type RRSP struct {
	Year                  int
	Rate                  float64
	MaxContribution       float64
	OverContributionLimit float64
	OverContributionRate  float64
	contribution          float64
}

func (rrsp *RRSP) CalculateMaxContribution(totalIncome float64) error {
	if err := rrsp.validate(totalIncome); err != nil {
		return err
	}
	rrsp.contribution = mathHelper.RoundFloat64(min(totalIncome*rrsp.Rate/100, rrsp.MaxContribution), 2)
	return nil
}

func (rrsp *RRSP) GetContribution() float64 {
	return rrsp.contribution
}

func (rrsp *RRSP) validate(totalIncome float64) (err error) {
	if totalIncome < 0 {
		return fmt.Errorf("rrsp: invalid income: \"%v\"", totalIncome)
	}
	if rrsp.Year <= 0 {
		return fmt.Errorf("rrsp: invalid year: \"%v\"", rrsp.Year)
	}
	if rrsp.OverContributionLimit < 0 {
		return fmt.Errorf("rrsp: invalid over contribution limit: \"%v\"", rrsp.OverContributionLimit)
	}
	if rrsp.Year <= 0 {
		return fmt.Errorf("rrsp: invalid over contribution rate: \"%v\"", rrsp.OverContributionRate)
	}
	return nil
}
