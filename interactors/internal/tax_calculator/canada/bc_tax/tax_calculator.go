package bc_tax

import (
	"context"
	dataService "github.com/gerdooshell/tax-core/data-access"
	"github.com/gerdooshell/tax-core/entities/canada/bc/credits"
	sharedEntities "github.com/gerdooshell/tax-core/entities/canada/shared"
	"github.com/gerdooshell/tax-core/interactors/data_access"
	bcDS "github.com/gerdooshell/tax-core/interactors/internal/tax_calculator/canada/bc_tax/data_structures"
	"github.com/gerdooshell/tax-core/interactors/internal/tax_calculator/canada/shared"
	"github.com/gerdooshell/tax-core/library/mathHelper"
)

type BritishColumbiaTaxCalculator interface {
	CalculateBCTax(ctx context.Context, input *bcDS.Input) (bcDS.Output, error)
}

func NewBCTaxCalculator() BritishColumbiaTaxCalculator {
	return &taxCalculatorImpl{
		dataProvider: dataService.NewDataProviderService("localhost:45432"),
	}
}

type taxCalculatorImpl struct {
	tax              sharedEntities.Tax
	brackets         []sharedEntities.TaxBracket
	creditsReduction sharedEntities.Tax
	bpa              credits.BasicPersonalAmount
	cpp              sharedEntities.CanadaPensionPlan
	dataProvider     dataAccess.BCTaxData
	ei               sharedEntities.EmploymentInsurancePremium
}

func (bc *taxCalculatorImpl) CalculateBCTax(ctx context.Context, input *bcDS.Input) (out bcDS.Output, err error) {
	eiErr := bc.processEI(ctx, input)
	cppErr := bc.processCPP(ctx, input)
	bpaErr := bc.processBPA(ctx, input)
	bracketErr := bc.processTaxBrackets(ctx, input)
	if err = <-bpaErr; err != nil {
		return
	}
	if err = <-bracketErr; err != nil {
		return
	}
	if err = <-eiErr; err != nil {
		return
	}
	if err = <-cppErr; err != nil {
		return
	}
	sharedDeductions := shared.AllCanadaTaxDeductions{
		CanadaPensionPlanAdditional:       bc.cpp.GetCPPFirstAdditionalEmployee(),
		CanadaPensionPlanSecondAdditional: bc.cpp.GetCPPSecondAdditionalEmployee(),
	}
	deductions := bcDS.BCTaxDeductions{
		AllCanadaTaxDeductions: sharedDeductions,
	}
	deductionAmount := deductions.CanadaPensionPlanAdditional + deductions.CanadaPensionPlanSecondAdditional
	bc.tax = sharedEntities.Tax{
		TaxBrackets: bc.brackets,
	}
	err = bc.tax.Calculate(input.Salary - deductionAmount)
	creditAmount := bc.bpa.GetValue() + bc.cpp.GetCPPBasicEmployee() + bc.ei.GetEIEmployee()
	bc.creditsReduction = sharedEntities.Tax{
		TaxBrackets: bc.brackets,
	}
	if err = bc.creditsReduction.Calculate(creditAmount); err != nil {
		return
	}
	out.Credits = bcDS.BCTaxCredits{
		BasicPensionAmount: bc.bpa.GetValue(),
		AllCanadaTaxCredits: shared.AllCanadaTaxCredits{
			CanadaPensionPlanBasic:     bc.cpp.GetCPPBasicEmployee(),
			EmploymentInsurancePremium: bc.ei.GetEIEmployee(),
		},
	}
	out.Deductions = deductions
	out.TotalTax = bc.tax.GetValue()
	out.PayableTax = mathHelper.RoundFloat64(bc.tax.GetValue()-bc.creditsReduction.GetValue(), 2)
	return out, nil
}

func (bc *taxCalculatorImpl) processEI(ctx context.Context, input *bcDS.Input) <-chan error {
	out := make(chan error)
	eiChan, errChan := bc.dataProvider.GetEIPremium(ctx, input.Year)

	go func() {
		defer close(out)
		var err error
		defer func() { out <- err }()
		select {
		case err = <-errChan:
			return
		case bc.ei = <-eiChan:
			err = bc.ei.Calculate(input.Salary)
		}
	}()
	return out
}

func (bc *taxCalculatorImpl) processCPP(ctx context.Context, input *bcDS.Input) <-chan error {
	out := make(chan error)
	cppChan, errChan := bc.dataProvider.GetCPP(ctx, input.Year)

	go func() {
		defer close(out)
		var err error
		defer func() { out <- err }()
		select {
		case err = <-errChan:
			return
		case bc.cpp = <-cppChan:
			err = bc.cpp.Calculate(input.Salary)
		}
	}()
	return out
}

func (bc *taxCalculatorImpl) processBPA(ctx context.Context, input *bcDS.Input) <-chan error {
	out := make(chan error)
	bpaChan, bpaErr := bc.dataProvider.GetBCBPA(ctx, input.Year)

	go func() {
		defer close(out)
		var err error
		defer func() { out <- err }()
		select {
		case err = <-bpaErr:
			return
		case bc.bpa = <-bpaChan:
			err = bc.bpa.Calculate()
		}
	}()
	return out
}

func (bc *taxCalculatorImpl) processTaxBrackets(ctx context.Context, input *bcDS.Input) <-chan error {
	out := make(chan error)
	bracketsChan, bracketsErr := bc.dataProvider.GetBCTaxBrackets(ctx, input.Year)
	go func() {
		defer close(out)
		var err error
		defer func() { out <- err }()
		select {
		case err = <-bracketsErr:
			return
		case bc.brackets = <-bracketsChan:
		}
	}()
	return out
}
