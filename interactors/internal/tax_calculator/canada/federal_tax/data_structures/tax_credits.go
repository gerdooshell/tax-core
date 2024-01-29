package dataStructures

import (
	canadaTaxCalculator "github.com/gerdooshell/tax-core/interactors/controller_access/canada_tax_calculator"
	"github.com/gerdooshell/tax-core/interactors/internal/tax_calculator/canada/shared"
)

type FederalTaxCredit struct {
	BasicPensionAmount     float64
	CanadaEmploymentAmount float64
	shared.AllCanadaTaxCredits
}

func FromModuleTaxCredit(cred canadaTaxCalculator.TaxCredits) FederalTaxCredit {
	return FederalTaxCredit{
		BasicPensionAmount:     cred.FederalBasicPensionAmount,
		AllCanadaTaxCredits:    shared.FromModuleTaxCredit(cred),
		CanadaEmploymentAmount: cred.CanadaEmploymentAmount,
	}
}
