package canadaTaxMrgins

import (
	"context"
	"fmt"
	"github.com/gerdooshell/tax-core/library/region/canada"

	dataProvider "github.com/gerdooshell/tax-core/data-access"
	"github.com/gerdooshell/tax-core/entities/canada/shared"
	dataAccess "github.com/gerdooshell/tax-core/interactors/data_access"
	marginDS "github.com/gerdooshell/tax-core/interactors/internal/tax_margins/canada/data_structures"
)

type TaxMarginsCa interface {
	GetCombinedMarginalBrackets(ctx context.Context, input marginDS.Input) (out marginDS.Output, err error)
}

func NewTaxMarginCa() TaxMarginsCa {
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

func (tm *taxMarginsCa) GetCombinedMarginalBrackets(ctx context.Context, input marginDS.Input) (out marginDS.Output, err error) {
	bracketsCahn, errChan := tm.dataProvider.GetCombinedMarginalBrackets(ctx, input.Year, input.Province)
	errRegChan := tm.getFederalBrackets(ctx, input)
	errFedChan := tm.getRegionalBrackets(ctx, input)
	select {
	case getError := <-errChan:
		fmt.Printf("failed getting combined marginal brackets: %v\n", getError)
	case brackets := <-bracketsCahn:
		if len(brackets) == 0 {
			break
		}
		out.Brackets = brackets
		return
	}

	if err = <-errFedChan; err != nil {
		return
	}
	if err = <-errRegChan; err != nil {
		return
	}
	tm.marginalTaxBrackets.RegionalTaxBrackets = tm.regionalTaxBrackets
	tm.marginalTaxBrackets.FederalBrackets = tm.federalTaxBrackets
	if err = tm.marginalTaxBrackets.CalcCombinedTaxMargins(); err != nil {
		return
	}
	brackets := tm.marginalTaxBrackets.GetMargins()
	saveChan, errSaveChan := tm.dataProvider.SaveMarginalTaxBrackets(ctx, input.Province, input.Year, brackets)
	select {
	case err = <-errSaveChan:
		return
	case _ = <-saveChan:
		fmt.Println("saved to database")
	}
	out.Brackets = brackets
	return
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
