package dataStructures

import (
	regionalDS "github.com/gerdooshell/tax-core/interactors/internal/tax_calculator/canada/regional/module_plugin/data_structures"
	"github.com/gerdooshell/tax-core/interactors/internal/tax_calculator/canada/shared"
)

type BCTaxDeductions struct {
	shared.AllCanadaTaxDeductions
}

func (td *BCTaxDeductions) ToRegionalTaxDeductions() regionalDS.RegionalTaxDeductions {
	return regionalDS.RegionalTaxDeductions{
		AllCanadaTaxDeductions: td.AllCanadaTaxDeductions,
	}
}

func FromRegionalTaxDeductions(deducts regionalDS.RegionalTaxDeductions) BCTaxDeductions {
	return BCTaxDeductions{
		AllCanadaTaxDeductions: deducts.AllCanadaTaxDeductions,
	}
}
