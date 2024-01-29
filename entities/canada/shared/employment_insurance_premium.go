package shared

import (
	"fmt"
	"github.com/gerdooshell/tax-core/library/mathHelper"
)

type EmploymentInsurancePremium struct {
	MaxInsurableEarning               float64
	Rate                              float64
	EmployerEmployeeContributionRatio float64
	eiEmployeeContribution            float64
	eiEmployerContribution            float64
}

// Calculate validates properties and calculates the ei premium contribution by employee and employer
func (eip *EmploymentInsurancePremium) Calculate(totalIncome float64) error {
	if err := eip.validateProperties(totalIncome); err != nil {
		return err
	}
	if totalIncome >= eip.MaxInsurableEarning {
		eip.eiEmployeeContribution = eip.MaxInsurableEarning * eip.Rate / 100
	} else {
		eip.eiEmployeeContribution = totalIncome * eip.Rate / 100
	}
	eip.eiEmployerContribution = eip.eiEmployeeContribution * eip.EmployerEmployeeContributionRatio
	eip.eiEmployeeContribution = mathHelper.RoundFloat64(eip.eiEmployeeContribution, 2)
	eip.eiEmployerContribution = mathHelper.RoundFloat64(eip.eiEmployerContribution, 2)
	return nil
}

func (eip *EmploymentInsurancePremium) GetEIEmployee() float64 {
	return eip.eiEmployeeContribution
}

func (eip *EmploymentInsurancePremium) GetEIEmployer() float64 {
	return eip.eiEmployerContribution
}

func (eip *EmploymentInsurancePremium) validateProperties(totalIncome float64) error {
	if eip.MaxInsurableEarning <= 0 {
		return fmt.Errorf("ei error: invalid max insurrable earcning: \"%v\"", eip.MaxInsurableEarning)
	}
	if eip.Rate < 0 || eip.Rate > 100 {
		return fmt.Errorf("ei error: invalid rate: \"%v\"", eip.Rate)
	}
	if eip.EmployerEmployeeContributionRatio < 0 {
		return fmt.Errorf("ei error: invalid employer contribution ratio: \"%v\"", eip.EmployerEmployeeContributionRatio)
	}
	if totalIncome < 0 {
		return fmt.Errorf("ei error: invalid totalIncome: \"%v\"", eip.EmployerEmployeeContributionRatio)
	}
	return nil
}
