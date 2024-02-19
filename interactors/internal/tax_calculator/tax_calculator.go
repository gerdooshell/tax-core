package taxCalculator

import (
	"context"
	"fmt"

	dataProvider "github.com/gerdooshell/tax-core/data-access"
	sharedEntities "github.com/gerdooshell/tax-core/entities/canada/shared"
	dataAccess "github.com/gerdooshell/tax-core/interactors/data_access"
	"github.com/gerdooshell/tax-core/library/region/canada"
)

type TaxBracketsInteractor interface {
	GetTotalTax(ctx context.Context, year int, province canada.Province, taxableIncome float64) <-chan TotalTaxOutput
	GetReducedTaxCredit(ctx context.Context, year int, province canada.Province, totalTaxCredit float64) <-chan TotalTaxOutput
}

type TotalTaxOutput struct {
	Value float64
	Err   error
}

func NewTaxInteractor() TaxBracketsInteractor {
	return &totalTaxImpl{
		dataProvider: dataProvider.GetDataProviderServiceInstance(),
	}
}

type totalTaxImpl struct {
	dataProvider dataAccess.TaxBracketData
}

func (t *totalTaxImpl) GetTotalTax(ctx context.Context, year int, province canada.Province, taxableIncome float64) <-chan TotalTaxOutput {
	return t.applyTaxBrackets(ctx, year, province, taxableIncome, false)
}

func (t *totalTaxImpl) GetReducedTaxCredit(ctx context.Context, year int, province canada.Province, totalTaxCredit float64) <-chan TotalTaxOutput {
	return t.applyTaxBrackets(ctx, year, province, totalTaxCredit, true)
}

func (t *totalTaxImpl) applyTaxBrackets(ctx context.Context, year int, province canada.Province, amount float64, isCredit bool) <-chan TotalTaxOutput {
	out := make(chan TotalTaxOutput)
	go func() {
		defer close(out)
		totalTaxOutput := TotalTaxOutput{}
		defer func() { out <- totalTaxOutput }()
		taxBracketsChan, errChan := t.dataProvider.GetTaxBrackets(ctx, year, province)
		select {
		case totalTaxOutput.Err = <-errChan:
			return
		case taxBrackets := <-taxBracketsChan:
			taxEntity := sharedEntities.Tax{TaxBrackets: taxBrackets}
			if err := taxEntity.Calculate(amount, isCredit); err != nil {
				totalTaxOutput.Err = err
			}
			totalTaxOutput.Value = taxEntity.GetValue()
		case <-ctx.Done():
			totalTaxOutput.Err = fmt.Errorf("processing tax brackets canceled")
			return
		}
	}()
	return out
}
