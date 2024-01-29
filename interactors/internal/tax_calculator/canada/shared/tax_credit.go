package shared

import (
	canadaTaxCalculator "github.com/gerdooshell/tax-core/interactors/controller_access/canada_tax_calculator"
)

type AllCanadaTaxCredits struct {
	CanadaPensionPlanBasic     float64
	EmploymentInsurancePremium float64
}

func FromModuleTaxCredit(cred canadaTaxCalculator.TaxCredits) AllCanadaTaxCredits {
	return AllCanadaTaxCredits{
		CanadaPensionPlanBasic:     cred.CanadaPensionPlanBasic,
		EmploymentInsurancePremium: cred.EIPremium,
	}
}
