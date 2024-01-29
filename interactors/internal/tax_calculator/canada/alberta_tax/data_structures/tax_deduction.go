package dataStructures

import (
	regionalDS "github.com/gerdooshell/tax-core/interactors/internal/tax_calculator/canada/regional/module_plugin/data_structures"
	"github.com/gerdooshell/tax-core/interactors/internal/tax_calculator/canada/shared"
)

type AlbertaTaxDeductions struct {
	shared.AllCanadaTaxDeductions
}

func (td *AlbertaTaxDeductions) ToRegionalTaxDeductions() regionalDS.RegionalTaxDeductions {
	return regionalDS.RegionalTaxDeductions{
		AllCanadaTaxDeductions: td.AllCanadaTaxDeductions,
	}
}

func FromRegionalTaxDeductions(deducts regionalDS.RegionalTaxDeductions) AlbertaTaxDeductions {
	return AlbertaTaxDeductions{
		AllCanadaTaxDeductions: deducts.AllCanadaTaxDeductions,
	}
}
