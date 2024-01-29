package canadaTaxImplementation

import (
	"context"
	"fmt"
	"github.com/gerdooshell/tax-core/interactors/controller_access/canada_tax_calculator"
	fed "github.com/gerdooshell/tax-core/interactors/internal/tax_calculator/canada/federal_tax"
	fedDS "github.com/gerdooshell/tax-core/interactors/internal/tax_calculator/canada/federal_tax/data_structures"
	"github.com/gerdooshell/tax-core/interactors/internal/tax_calculator/canada/regional"
	regionalPlugin "github.com/gerdooshell/tax-core/interactors/internal/tax_calculator/canada/regional/module_plugin"
	regionalDS "github.com/gerdooshell/tax-core/interactors/internal/tax_calculator/canada/regional/module_plugin/data_structures"
	"github.com/gerdooshell/tax-core/library/region/canada"
)

func NewCanadaTaxCalculator() canadaTaxCalculator.CanadaTaxCalculator {
	return &taxCalculatorImpl{}
}

type taxCalculatorImpl struct {
	province              canada.Province
	regionalTaxCalculator regionalPlugin.TaxCalculator
	federalTaxCalculator  fed.FederalTaxCalculator
}

func (tax taxCalculatorImpl) Calculate(ctx context.Context, input *canadaTaxCalculator.Input) (out canadaTaxCalculator.Output, err error) {
	if input == nil {
		err = fmt.Errorf("canada tax calculator error: invalid input: nil ")
		return
	}
	tax.province = input.Province
	tax.regionalTaxCalculator = regional.NewRegionalTaxCalculator()
	regionalInput := mapInputToRegionalTaxInput(input)
	regionalTax, err := tax.regionalTaxCalculator.CalculateRegionalTax(ctx, regionalInput)
	if err != nil {
		return
	}
	tax.federalTaxCalculator = fed.NewTaxCalculator()
	federalTaxInput := mapInputToFederalTaxInput(input)
	federalTax, err := tax.federalTaxCalculator.CalculateFederalTax(ctx, federalTaxInput)
	if err != nil {
		return
	}
	out = mergeRegionalAndFederalTaxOutputs(regionalTax, federalTax)
	return out, nil
}

func mapInputToRegionalTaxInput(input *canadaTaxCalculator.Input) *regionalDS.Input {
	return &regionalDS.Input{
		Credits:    regionalDS.FromModuleTaxCredits(input.Credits),
		Province:   input.Province,
		Deductions: regionalDS.FromModuleTaxDeductions(input.Deductions),
		Salary:     input.Salary,
		Year:       input.Year,
	}
}

func mapInputToFederalTaxInput(input *canadaTaxCalculator.Input) *fedDS.Input {
	return &fedDS.Input{
		Credits:    fedDS.FromModuleTaxCredit(input.Credits),
		Deductions: fedDS.FromModuleTaxDeductions(input.Deductions),
		Salary:     input.Salary,
		Year:       input.Year,
	}
}

func mergeRegionalAndFederalTaxOutputs(regionalTax regionalDS.Output, federalTax fedDS.Output) canadaTaxCalculator.Output {
	return canadaTaxCalculator.Output{
		TaxCredits:         mergeRegionalAndFederalTaxCredits(regionalTax.Credits, federalTax.Credits),
		TaxDeductions:      mergeRegionalAndFederalTaxDeductions(regionalTax.Deductions, federalTax.Deductions),
		FederalPayableTax:  federalTax.PayableTax,
		FederalTotalTax:    federalTax.TotalTax,
		RegionalPayableTax: regionalTax.PayableTax,
		RegionalTotalTax:   regionalTax.TotalTax,
	}
}

func mergeRegionalAndFederalTaxCredits(regCredits regionalDS.RegionalTaxCredit, fedCredits fedDS.FederalTaxCredit) canadaTaxCalculator.TaxCredits {
	return canadaTaxCalculator.TaxCredits{
		FederalBasicPensionAmount:  fedCredits.BasicPensionAmount,
		EIPremium:                  fedCredits.EmploymentInsurancePremium,
		CanadaPensionPlanBasic:     fedCredits.CanadaPensionPlanBasic,
		RegionalBasicPensionAmount: regCredits.BasicPensionAmount,
		CanadaEmploymentAmount:     fedCredits.CanadaEmploymentAmount,
	}
}

func mergeRegionalAndFederalTaxDeductions(regDeducts regionalDS.RegionalTaxDeductions, fedDeducts fedDS.FederalTaxDeductions) canadaTaxCalculator.TaxDeductions {
	return canadaTaxCalculator.TaxDeductions{
		CPPFirstAdditional:  fedDeducts.CanadaPensionPlanAdditional,
		CPPSecondAdditional: fedDeducts.CanadaPensionPlanSecondAdditional,
	}
}
