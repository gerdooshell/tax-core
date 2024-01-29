package dataStructures

import (
	"github.com/gerdooshell/tax-core/interactors/internal/tax_calculator/canada/regional/module_plugin/data_structures"
	"github.com/gerdooshell/tax-core/interactors/internal/tax_calculator/canada/shared"
)

type BCTaxCredits struct {
	shared.AllCanadaTaxCredits
	BasicPensionAmount float64
}

func (tc *BCTaxCredits) ToRegionalTaxCredit() dataStructures.RegionalTaxCredit {
	return dataStructures.RegionalTaxCredit{
		AllCanadaTaxCredits: tc.AllCanadaTaxCredits,
		BasicPensionAmount:  tc.BasicPensionAmount,
	}
}

func FromRegionalTaxCredit(cred dataStructures.RegionalTaxCredit) BCTaxCredits {
	return BCTaxCredits{
		AllCanadaTaxCredits: cred.AllCanadaTaxCredits,
		BasicPensionAmount:  cred.BasicPensionAmount,
	}
}
