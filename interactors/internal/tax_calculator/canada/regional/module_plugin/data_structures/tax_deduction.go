package dataStructures

import (
	canadaTaxCalculator "github.com/gerdooshell/tax-core/interactors/controller_access/canada_tax_calculator"
	"github.com/gerdooshell/tax-core/interactors/internal/tax_calculator/canada/shared"
)

type RegionalTaxDeductions struct {
	shared.AllCanadaTaxDeductions
}

func FromModuleTaxDeductions(deducts canadaTaxCalculator.TaxDeductions) RegionalTaxDeductions {
	return RegionalTaxDeductions{
		AllCanadaTaxDeductions: shared.FromModuleTaxDeductions(deducts),
	}
}
