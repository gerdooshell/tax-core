package marginCalculator

import (
	"context"
	"fmt"
	"github.com/gerdooshell/tax-core/interactors/internal/margin_calculator/data_structures"
	"github.com/gerdooshell/tax-core/library/region/canada"

	dataProvider "github.com/gerdooshell/tax-core/data-access"
	"github.com/gerdooshell/tax-core/entities/canada/shared"
	dataAccess "github.com/gerdooshell/tax-core/interactors/data_access"
)

type TaxMarginsInteractor interface {
	GetCombinedMarginalBrackets(ctx context.Context, input marginDS.Input) <-chan marginDS.Output
}

func NewTaxMarginCalculator() TaxMarginsInteractor {
	return &taxMarginsCa{
		dataProvider: dataProvider.GetDataProviderServiceInstance(),
	}
}

type taxMarginsCa struct {
	dataProvider        dataAccess.TaxMargin
	federalTaxBrackets  []shared.TaxBracket
	regionalTaxBrackets []shared.TaxBracket
	marginalTaxBrackets shared.TaxMarginalBracket
}

func (tm *taxMarginsCa) GetCombinedMarginalBrackets(ctx context.Context, input marginDS.Input) <-chan marginDS.Output {
	out := make(chan marginDS.Output, 1)
	go func() {
		defer close(out)
		marginOut := marginDS.Output{}
		defer func() { out <- marginOut }()
		bracketsChan, errChan := tm.dataProvider.GetCombinedMarginalBrackets(ctx, input.Year, input.Province)
		errRegChan := tm.getFederalBrackets(ctx, input)
		errFedChan := tm.getRegionalBrackets(ctx, input)
		select {
		case getError := <-errChan:
			fmt.Printf("failed getting combined marginal brackets: %v\n", getError)
		case brackets := <-bracketsChan:
			if len(brackets) == 0 {
				break
			}
			marginOut.Brackets = brackets
			return
		}
		if marginOut.Err = <-errFedChan; marginOut.Err != nil {
			return
		}
		if marginOut.Err = <-errRegChan; marginOut.Err != nil {
			return
		}
		tm.marginalTaxBrackets.RegionalTaxBrackets = tm.regionalTaxBrackets
		tm.marginalTaxBrackets.FederalBrackets = tm.federalTaxBrackets
		if marginOut.Err = tm.marginalTaxBrackets.CalcCombinedTaxMargins(); marginOut.Err != nil {
			return
		}
		brackets := tm.marginalTaxBrackets.GetMargins()
		saveChan, errSaveChan := tm.dataProvider.SaveMarginalTaxBrackets(ctx, input.Province, input.Year, brackets)
		select {
		case marginOut.Err = <-errSaveChan:
			return
		case _ = <-saveChan:
			fmt.Println("saved to database")
		}
		marginOut.Brackets = brackets
	}()
	return out
}

func (tm *taxMarginsCa) getFederalBrackets(ctx context.Context, input marginDS.Input) <-chan error {
	errChan := make(chan error, 1)
	go func() {
		defer close(errChan)
		var err error
		defer func() { errChan <- err }()
		out, errOut := tm.dataProvider.GetTaxBrackets(ctx, input.Year, canada.Federal)
		select {
		case tm.federalTaxBrackets = <-out:
		case err = <-errOut:
		}
	}()
	return errChan
}

func (tm *taxMarginsCa) getRegionalBrackets(ctx context.Context, input marginDS.Input) <-chan error {
	errChan := make(chan error, 1)
	go func() {
		defer close(errChan)
		var err error
		defer func() { errChan <- err }()
		out, errOut := tm.dataProvider.GetTaxBrackets(ctx, input.Year, input.Province)
		select {
		case tm.regionalTaxBrackets = <-out:
		case err = <-errOut:
		}
	}()
	return errChan
}
