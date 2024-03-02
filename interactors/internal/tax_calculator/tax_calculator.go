package taxCalculator

import (
	"context"
	"errors"
	"time"

	dataProvider "github.com/gerdooshell/tax-core/data-access"
	sharedEntities "github.com/gerdooshell/tax-core/entities/canada/shared"
	dataAccess "github.com/gerdooshell/tax-core/interactors/data_access"
	sharedCredits "github.com/gerdooshell/tax-core/interactors/internal/shared_credits"
	"github.com/gerdooshell/tax-core/library/mathHelper"
	"github.com/gerdooshell/tax-core/library/region/canada"
)

type TaxBracketsInteractor interface {
	GetTotalTax(ctx context.Context, year int, province canada.Province, taxableIncome float64) <-chan TaxOutput
	GetReducedTaxCredit(ctx context.Context, year int, province canada.Province, totalTaxCredit float64) <-chan TaxOutput
	GetPayableTaxGivenDeductions(ctx context.Context, year int, province canada.Province, totalIncome, totalDeductions float64) <-chan TaxOutput
}

type TaxOutput struct {
	Value float64
	Err   error
}

func NewTaxInteractor() TaxBracketsInteractor {
	return &totalTaxImpl{
		dataProvider:      dataProvider.GetDataProviderServiceInstance(),
		creditsCalculator: sharedCredits.NewTaxCreditInteractor(),
		timeout:           time.Second * 10,
	}
}

type totalTaxImpl struct {
	dataProvider      dataAccess.TaxBracketData
	creditsCalculator sharedCredits.TaxCreditInteractor
	timeout           time.Duration
}

func (t *totalTaxImpl) GetTotalTax(ctx context.Context, year int, province canada.Province, taxableIncome float64) <-chan TaxOutput {
	return t.applyTaxBrackets(ctx, year, province, taxableIncome, false)
}

func (t *totalTaxImpl) GetReducedTaxCredit(ctx context.Context, year int, province canada.Province, totalTaxCredit float64) <-chan TaxOutput {
	return t.applyTaxBrackets(ctx, year, province, totalTaxCredit, true)
}

func (t *totalTaxImpl) GetPayableTaxGivenDeductions(ctx context.Context, year int, province canada.Province, totalIncome, totalDeductions float64) <-chan TaxOutput {
	out := make(chan TaxOutput, 1)
	go func() {
		defer close(out)
		taxOutput := TaxOutput{}
		defer func() { out <- taxOutput }()
		taxableIncome := max(totalIncome-totalDeductions, 0)
		totalTaxRegionalChan := t.GetTotalTax(ctx, year, province, taxableIncome)
		totalTaxFederalChan := t.GetTotalTax(ctx, year, canada.Federal, taxableIncome)
		creditsRegionalChan := t.creditsCalculator.GetTaxCredits(ctx, year, province, totalIncome)
		taxCredits := <-creditsRegionalChan
		if taxCredits.Err != nil {
			taxOutput.Err = taxCredits.Err
			return
		}
		totalCreditsRegional := taxCredits.EIPremiumEmployee + taxCredits.CPPBasicEmployee + taxCredits.RegionalBPA
		totalCreditsFederal := taxCredits.CanadaEmploymentAmount + taxCredits.EIPremiumEmployee + taxCredits.CPPBasicEmployee + taxCredits.FederalBPA
		reducedCreditsRegionalChan := t.GetReducedTaxCredit(ctx, year, province, totalCreditsRegional)
		reducedCreditsFederalChan := t.GetReducedTaxCredit(ctx, year, canada.Federal, totalCreditsFederal)
		totalTaxRegional := <-totalTaxRegionalChan
		if totalTaxRegional.Err != nil {
			taxOutput.Err = totalTaxRegional.Err
			return
		}
		totalTaxFederal := <-totalTaxFederalChan
		if totalTaxFederal.Err != nil {
			taxOutput.Err = totalTaxFederal.Err
			return
		}
		reducedCreditsRegional := <-reducedCreditsRegionalChan
		if reducedCreditsRegional.Err != nil {
			taxOutput.Err = reducedCreditsRegional.Err
			return
		}
		reducedCreditsFederal := <-reducedCreditsFederalChan
		if reducedCreditsFederal.Err != nil {
			taxOutput.Err = reducedCreditsFederal.Err
			return
		}
		taxOutput.Value = max(mathHelper.RoundFloat64(totalTaxRegional.Value+totalTaxFederal.Value-reducedCreditsRegional.Value-reducedCreditsFederal.Value, 2), 0)
	}()
	return out
}

func (t *totalTaxImpl) applyTaxBrackets(ctx context.Context, year int, province canada.Province, amount float64, isCredit bool) <-chan TaxOutput {
	out := make(chan TaxOutput, 1)
	go func() {
		defer close(out)
		totalTaxOutput := TaxOutput{}
		defer func() { out <- totalTaxOutput }()
		taxBracketsChan := t.dataProvider.GetTaxBrackets(ctx, year, province)
		select {
		case taxBracketsDataOut := <-taxBracketsChan:
			if totalTaxOutput.Err = taxBracketsDataOut.Err; totalTaxOutput.Err != nil {
				return
			}
			taxEntity := sharedEntities.Tax{TaxBrackets: taxBracketsDataOut.TaxBrackets}
			if err := taxEntity.Calculate(amount, isCredit); err != nil {
				totalTaxOutput.Err = err
			}
			totalTaxOutput.Value = taxEntity.GetValue()
		case <-time.After(t.timeout):
			totalTaxOutput.Err = errors.New("get tax brackets data timed out")
		}
	}()
	return out
}
