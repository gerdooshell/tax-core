package shared

import (
	canadaTaxCalculator "github.com/gerdooshell/tax-core/interactors/controller_access/canada_tax_calculator"
)

type AllCanadaTaxDeductions struct {
	CanadaPensionPlanAdditional       float64
	CanadaPensionPlanSecondAdditional float64
}

func FromModuleTaxDeductions(deducts canadaTaxCalculator.TaxDeductions) AllCanadaTaxDeductions {
	return AllCanadaTaxDeductions{
		CanadaPensionPlanAdditional:       deducts.CPPFirstAdditional,
		CanadaPensionPlanSecondAdditional: deducts.CPPSecondAdditional,
	}
}
