package rrspInfo

import "github.com/gerdooshell/tax-core/library/region/canada"

type OptimalInput struct {
	TotalIncome     float64
	ContributedRRSP float64
	Year            int
	Province        canada.Province
}

type OptimalOutput struct {
	TaxableIncome      float64
	RRSP               float64
	PayableTax         float64
	SuggestedTaxReturn float64
	TaxReturn          float64
	LeftRRSPRoom       float64
}
