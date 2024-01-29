package regional

import (
	"context"
	ab "github.com/gerdooshell/tax-core/interactors/internal/tax_calculator/canada/alberta_tax"
	albertaDS "github.com/gerdooshell/tax-core/interactors/internal/tax_calculator/canada/alberta_tax/data_structures"
	bc "github.com/gerdooshell/tax-core/interactors/internal/tax_calculator/canada/bc_tax"
	bcDS "github.com/gerdooshell/tax-core/interactors/internal/tax_calculator/canada/bc_tax/data_structures"
	modulePlugin "github.com/gerdooshell/tax-core/interactors/internal/tax_calculator/canada/regional/module_plugin"
	regionalDS "github.com/gerdooshell/tax-core/interactors/internal/tax_calculator/canada/regional/module_plugin/data_structures"
	"github.com/gerdooshell/tax-core/library/region/canada"
)

func NewRegionalTaxCalculator() modulePlugin.TaxCalculator {
	return &regionalTaxCalculatorImpl{}
}

type regionalTaxCalculatorImpl struct {
}

func (reg *regionalTaxCalculatorImpl) CalculateRegionalTax(ctx context.Context, input *regionalDS.Input) (out regionalDS.Output, err error) {
	region := input.Province
	switch region {
	case canada.Alberta:
		albertaTaxInput := mapInputToAlbertaInput(input)
		calculator := ab.NewAlbertaTaxCalculator()
		albertaOut, err := calculator.CalculateAlbertaTax(ctx, albertaTaxInput)
		if err != nil {
			return out, err
		}
		out = mapAlbertaOutToOutput(albertaOut)
	case canada.BritishColumbia:
		bcTaxInput := mapInputToBCInput(input)
		calculator := bc.NewBCTaxCalculator()
		bcOutput, err := calculator.CalculateBCTax(ctx, bcTaxInput)
		if err != nil {
			return out, err
		}
		out = mapBCOutToOutput(bcOutput)
	}
	return
}

func mapInputToAlbertaInput(input *regionalDS.Input) *albertaDS.Input {
	return &albertaDS.Input{
		Credits:    albertaDS.FromRegionalTaxCredit(input.Credits),
		Deductions: albertaDS.FromRegionalTaxDeductions(input.Deductions),
		Salary:     input.Salary,
		Year:       input.Year,
	}
}

func mapAlbertaOutToOutput(albertaOut albertaDS.Output) regionalDS.Output {
	return regionalDS.Output{
		Credits:    albertaOut.Credits.ToRegionalTaxCredit(),
		Deductions: albertaOut.Deductions.ToRegionalTaxDeductions(),
		PayableTax: albertaOut.PayableTax,
		TotalTax:   albertaOut.TotalTax,
	}
}

func mapInputToBCInput(input *regionalDS.Input) *bcDS.Input {
	return &bcDS.Input{
		Credits:    bcDS.FromRegionalTaxCredit(input.Credits),
		Deductions: bcDS.FromRegionalTaxDeductions(input.Deductions),
		Salary:     input.Salary,
		Year:       input.Year,
	}
}

func mapBCOutToOutput(bcOut bcDS.Output) regionalDS.Output {
	return regionalDS.Output{
		Credits:    bcOut.Credits.ToRegionalTaxCredit(),
		Deductions: bcOut.Deductions.ToRegionalTaxDeductions(),
		PayableTax: bcOut.PayableTax,
		TotalTax:   bcOut.TotalTax,
	}
}
