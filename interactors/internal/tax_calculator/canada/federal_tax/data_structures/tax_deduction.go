package dataStructures

import (
	canadaTaxCalculator "github.com/gerdooshell/tax-core/interactors/controller_access/canada_tax_calculator"
	"github.com/gerdooshell/tax-core/interactors/internal/tax_calculator/canada/shared"
)

type FederalTaxDeductions struct {
	shared.AllCanadaTaxDeductions
	//TODO: implement me
}

func FromModuleTaxDeductions(deducts canadaTaxCalculator.TaxDeductions) FederalTaxDeductions {
	return FederalTaxDeductions{
		AllCanadaTaxDeductions: shared.FromModuleTaxDeductions(deducts),
	}
}
