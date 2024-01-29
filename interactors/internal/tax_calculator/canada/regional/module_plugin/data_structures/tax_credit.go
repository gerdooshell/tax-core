package dataStructures

import (
	canadaTaxCalculator "github.com/gerdooshell/tax-core/interactors/controller_access/canada_tax_calculator"
	"github.com/gerdooshell/tax-core/interactors/internal/tax_calculator/canada/shared"
)

type RegionalTaxCredit struct {
	shared.AllCanadaTaxCredits
	BasicPensionAmount float64
}

func FromModuleTaxCredits(cred canadaTaxCalculator.TaxCredits) RegionalTaxCredit {
	return RegionalTaxCredit{
		AllCanadaTaxCredits: shared.FromModuleTaxCredit(cred),
		BasicPensionAmount:  cred.RegionalBasicPensionAmount,
	}
}
