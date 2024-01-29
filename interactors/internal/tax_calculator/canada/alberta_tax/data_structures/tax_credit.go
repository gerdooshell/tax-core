package dataStructures

import (
	modulePlugin "github.com/gerdooshell/tax-core/interactors/internal/tax_calculator/canada/regional/module_plugin/data_structures"
	"github.com/gerdooshell/tax-core/interactors/internal/tax_calculator/canada/shared"
)

type AlbertaTaxCredit struct {
	shared.AllCanadaTaxCredits
	BasicPensionAmount float64
}

func (tc *AlbertaTaxCredit) ToRegionalTaxCredit() modulePlugin.RegionalTaxCredit {
	return modulePlugin.RegionalTaxCredit{
		BasicPensionAmount:  tc.BasicPensionAmount,
		AllCanadaTaxCredits: tc.AllCanadaTaxCredits,
	}
}

func FromRegionalTaxCredit(cred modulePlugin.RegionalTaxCredit) AlbertaTaxCredit {
	return AlbertaTaxCredit{
		BasicPensionAmount:  cred.BasicPensionAmount,
		AllCanadaTaxCredits: cred.AllCanadaTaxCredits,
	}
}
