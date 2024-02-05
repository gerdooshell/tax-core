package alberta_tax

import (
	"context"
	dataService "github.com/gerdooshell/tax-core/data-access"
	"github.com/gerdooshell/tax-core/entities/canada/alberta/credits"
	"github.com/gerdooshell/tax-core/entities/canada/shared"
	"github.com/gerdooshell/tax-core/interactors/data_access"
	abDS "github.com/gerdooshell/tax-core/interactors/internal/tax_calculator/canada/alberta_tax/data_structures"
	shared2 "github.com/gerdooshell/tax-core/interactors/internal/tax_calculator/canada/shared"
	"github.com/gerdooshell/tax-core/library/mathHelper"
	"github.com/gerdooshell/tax-core/library/region/canada"
)

type AlbertaTaxCalculator interface {
	CalculateAlbertaTax(context.Context, *abDS.Input) (abDS.Output, error)
}

func NewAlbertaTaxCalculator() AlbertaTaxCalculator {
	return &taxCalculatorImpl{
		dataProvider: dataService.NewDataProviderService(),
	}
}

type taxCalculatorImpl struct {
	ei               shared.EmploymentInsurancePremium
	bpa              credits.BasicPersonalAmount
	brackets         []shared.TaxBracket
	dataProvider     dataAccess.AlbertaTaxData
	cpp              shared.CanadaPensionPlan
	tax              shared.Tax
	creditsReduction shared.Tax
}

func (ab *taxCalculatorImpl) CalculateAlbertaTax(ctx context.Context, input *abDS.Input) (out abDS.Output, err error) {
	eiErr := ab.processEI(ctx, input)
	bpaErr := ab.processBPA(ctx, input)
	cppErr := ab.processCPP(ctx, input)
	bracketErr := ab.processTaxBrackets(ctx, input)
	if err = <-eiErr; err != nil {
		return
	}
	if err = <-bpaErr; err != nil {
		return
	}
	if err = <-cppErr; err != nil {
		return
	}
	if err = <-bracketErr; err != nil {
		return
	}
	taxableIncome := input.Salary - ab.cpp.GetCPPFirstAdditionalEmployee() - ab.cpp.GetCPPSecondAdditionalEmployee()
	ab.tax = shared.Tax{
		TaxBrackets: ab.brackets,
	}
	if err = ab.tax.Calculate(taxableIncome); err != nil {
		return
	}

	ab.creditsReduction = shared.Tax{
		TaxBrackets: ab.brackets,
	}
	reductionAmount := ab.bpa.GetValue() + ab.ei.GetEIEmployee() + ab.cpp.GetCPPBasicEmployee()
	if err = ab.creditsReduction.Calculate(reductionAmount); err != nil {
		return
	}
	out.Credits = abDS.AlbertaTaxCredit{
		BasicPensionAmount: ab.bpa.GetValue(),
		AllCanadaTaxCredits: shared2.AllCanadaTaxCredits{
			CanadaPensionPlanBasic:     ab.cpp.GetCPPBasicEmployee(),
			EmploymentInsurancePremium: ab.ei.GetEIEmployee(),
		},
	}
	out.Deductions = abDS.AlbertaTaxDeductions{
		AllCanadaTaxDeductions: shared2.AllCanadaTaxDeductions{
			CanadaPensionPlanAdditional:       ab.cpp.GetCPPFirstAdditionalEmployee(),
			CanadaPensionPlanSecondAdditional: ab.cpp.GetCPPSecondAdditionalEmployee(),
		},
	}
	out.TotalTax = ab.tax.GetValue()
	out.PayableTax = mathHelper.RoundFloat64(ab.tax.GetValue()-ab.creditsReduction.GetValue(), 2)
	return out, nil
}

func (ab *taxCalculatorImpl) processEI(ctx context.Context, input *abDS.Input) <-chan error {
	out := make(chan error)
	eiChan, errChan := ab.dataProvider.GetEIPremium(ctx, input.Year)

	go func() {
		defer close(out)
		var err error
		defer func() { out <- err }()
		select {
		case err = <-errChan:
			return
		case ab.ei = <-eiChan:
			err = ab.ei.Calculate(input.Salary)
		}
	}()
	return out
}

func (ab *taxCalculatorImpl) processBPA(ctx context.Context, input *abDS.Input) <-chan error {
	out := make(chan error)
	bpaChan, bpaErr := ab.dataProvider.GetAlbertaBPA(ctx, input.Year)

	go func() {
		defer close(out)
		var err error
		defer func() { out <- err }()
		select {
		case err = <-bpaErr:
			return
		case ab.bpa = <-bpaChan:
			err = ab.bpa.Calculate()
		}
	}()
	return out
}

func (ab *taxCalculatorImpl) processTaxBrackets(ctx context.Context, input *abDS.Input) <-chan error {
	out := make(chan error)
	bracketsChan, bracketsErr := ab.dataProvider.GetTaxBrackets(ctx, input.Year, canada.Alberta)
	go func() {
		defer close(out)
		var err error
		defer func() { out <- err }()
		select {
		case err = <-bracketsErr:
			return
		case ab.brackets = <-bracketsChan:
		}
	}()
	return out
}

func (ab *taxCalculatorImpl) processCPP(ctx context.Context, input *abDS.Input) <-chan error {
	out := make(chan error)
	cppChan, errChan := ab.dataProvider.GetCPP(ctx, input.Year)

	go func() {
		defer close(out)
		var err error
		defer func() { out <- err }()
		select {
		case err = <-errChan:
			return
		case ab.cpp = <-cppChan:
			err = ab.cpp.Calculate(input.Salary)
		}
	}()
	return out
}
