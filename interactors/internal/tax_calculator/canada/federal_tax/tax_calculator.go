package federal_tax

import (
	"context"
	dataService "github.com/gerdooshell/tax-core/data-access"
	"github.com/gerdooshell/tax-core/entities/canada/federal/credits"
	sharedEntities "github.com/gerdooshell/tax-core/entities/canada/shared"
	"github.com/gerdooshell/tax-core/interactors/data_access"
	fedDS "github.com/gerdooshell/tax-core/interactors/internal/tax_calculator/canada/federal_tax/data_structures"
	"github.com/gerdooshell/tax-core/interactors/internal/tax_calculator/canada/shared"
	"github.com/gerdooshell/tax-core/library/mathHelper"
	"github.com/gerdooshell/tax-core/library/region/canada"
)

type FederalTaxCalculator interface {
	CalculateFederalTax(context.Context, *fedDS.Input) (fedDS.Output, error)
}

func NewTaxCalculator() FederalTaxCalculator {
	return &federalTaxCalculatorImpl{
		dataProvider: dataService.GetDataProviderServiceInstance(),
	}
}

type federalTaxCalculatorImpl struct {
	dataProvider     dataAccess.FederalTaxData
	bpa              credits.BasicPersonalAmount
	cea              credits.CanadaEmploymentAmount
	cpp              sharedEntities.CanadaPensionPlan
	ei               sharedEntities.EmploymentInsurancePremium
	brackets         []sharedEntities.TaxBracket
	tax              sharedEntities.Tax
	creditsReduction sharedEntities.Tax
}

func (fed *federalTaxCalculatorImpl) CalculateFederalTax(ctx context.Context, input *fedDS.Input) (out fedDS.Output, err error) {

	eiErr := fed.processEI(ctx, input)
	cppErr := fed.processCPP(ctx, input)
	bpaErr := fed.processBPA(ctx, input)
	ceaErr := fed.processCEA(ctx, input)
	bracketsErr := fed.processBrackets(ctx, input)

	if err = <-eiErr; err != nil {
		return
	}
	if err = <-cppErr; err != nil {
		return
	}
	if err = <-bpaErr; err != nil {
		return
	}
	if err = <-ceaErr; err != nil {
		return
	}
	if err = <-bracketsErr; err != nil {
		return
	}

	sharedDeductions := shared.AllCanadaTaxDeductions{
		CanadaPensionPlanAdditional:       fed.cpp.GetCPPFirstAdditionalEmployee(),
		CanadaPensionPlanSecondAdditional: fed.cpp.GetCPPSecondAdditionalEmployee(),
	}
	deductions := fedDS.FederalTaxDeductions{
		AllCanadaTaxDeductions: sharedDeductions,
	}
	taxableIncome := input.Salary - deductions.CanadaPensionPlanAdditional - deductions.CanadaPensionPlanSecondAdditional

	fed.tax = sharedEntities.Tax{
		TaxBrackets: fed.brackets,
	}
	if err = fed.tax.Calculate(taxableIncome, false); err != nil {
		return
	}

	fed.creditsReduction = sharedEntities.Tax{
		TaxBrackets: fed.brackets,
	}
	creditAmount := fed.ei.GetEIEmployee() + fed.cpp.GetCPPBasicEmployee() + fed.bpa.GetValue() + fed.cea.GetEmployeeValue()
	if err = fed.creditsReduction.Calculate(creditAmount, true); err != nil {
		return
	}
	out.Deductions = deductions
	out.Credits = fedDS.FederalTaxCredit{
		BasicPensionAmount:     fed.bpa.GetValue(),
		CanadaEmploymentAmount: fed.cea.GetEmployeeValue(),
		AllCanadaTaxCredits: shared.AllCanadaTaxCredits{
			CanadaPensionPlanBasic:     fed.cpp.GetCPPBasicEmployee(),
			EmploymentInsurancePremium: fed.ei.GetEIEmployee(),
		},
	}
	out.TotalTax = fed.tax.GetValue()
	out.PayableTax = mathHelper.RoundFloat64(fed.tax.GetValue()-fed.creditsReduction.GetValue(), 2)

	return
}

func (fed *federalTaxCalculatorImpl) processEI(ctx context.Context, input *fedDS.Input) <-chan error {
	out := make(chan error)
	eiChan, errChan := fed.dataProvider.GetEIPremium(ctx, input.Year)

	go func() {
		defer close(out)
		var err error
		defer func() { out <- err }()
		select {
		case err = <-errChan:
			return
		case fed.ei = <-eiChan:
			err = fed.ei.Calculate(input.Salary)
		}
	}()
	return out
}

func (fed *federalTaxCalculatorImpl) processCEA(ctx context.Context, input *fedDS.Input) <-chan error {
	out := make(chan error)
	ceaChan, errChan := fed.dataProvider.GetCEA(ctx, input.Year)

	go func() {
		defer close(out)
		var err error
		defer func() { out <- err }()
		select {
		case err = <-errChan:
			return
		case fed.cea = <-ceaChan:
			err = fed.cea.Calculate(input.Salary)
		}
	}()
	return out
}

func (fed *federalTaxCalculatorImpl) processCPP(ctx context.Context, input *fedDS.Input) <-chan error {
	out := make(chan error)
	cppChan, errChan := fed.dataProvider.GetCPP(ctx, input.Year)

	go func() {
		defer close(out)
		var err error
		defer func() { out <- err }()
		select {
		case err = <-errChan:
			return
		case fed.cpp = <-cppChan:
			err = fed.cpp.Calculate(input.Salary)
		}
	}()
	return out
}

func (fed *federalTaxCalculatorImpl) processBPA(ctx context.Context, input *fedDS.Input) <-chan error {
	out := make(chan error)
	bpaChan, errChan := fed.dataProvider.GetFederalBPA(ctx, input.Year)

	go func() {
		defer close(out)
		var err error
		defer func() { out <- err }()
		select {
		case err = <-errChan:
			return
		case fed.bpa = <-bpaChan:
			err = fed.bpa.Calculate(input.Salary)
		}
	}()
	return out
}

func (fed *federalTaxCalculatorImpl) processBrackets(ctx context.Context, input *fedDS.Input) <-chan error {
	out := make(chan error)
	bracketsChan, errChan := fed.dataProvider.GetTaxBrackets(ctx, input.Year, canada.Federal)

	go func() {
		defer close(out)
		var err error
		defer func() { out <- err }()
		select {
		case err = <-errChan:
			return
		case fed.brackets = <-bracketsChan:
		}
	}()
	return out
}
