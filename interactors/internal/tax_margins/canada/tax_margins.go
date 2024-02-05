package canadaTaxMrgins

import (
	"context"

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
		dataProvider: dataProvider.NewDataProviderService(),
	}
}

type taxMarginsCa struct {
	dataProvider        dataAccess.TaxMargin
	federalTaxBrackets  []shared.TaxBracket
	regionalTaxBrackets []shared.TaxBracket
	marginalTaxBrackets shared.TaxMarginalBracket
}

func (tm *taxMarginsCa) GetCombinedMarginalBrackets(ctx context.Context, input marginDS.Input) (out marginDS.Output, err error) {
	// TODO: try to read from database first
	errRegChan := tm.getFederalBrackets(ctx, input)
	errFedChan := tm.getRegionalBrackets(ctx, input)
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
	_, errSaveChan := tm.dataProvider.SaveMarginalTaxBrackets(ctx, input.Province, input.Year, brackets)
	if err = <-errSaveChan; err != nil {
		return
	}
	out.Brackets = brackets
	return
}

func (tm *taxMarginsCa) getFederalBrackets(ctx context.Context, input marginDS.Input) <-chan error {
	errChan := make(chan error)
	go func() {
		defer close(errChan)
		var err error
		defer func() { errChan <- err }()
		out, errOut := tm.dataProvider.GetTaxBrackets(ctx, input.Year, input.Province)
		select {
		case tm.federalTaxBrackets = <-out:
		case err = <-errOut:
		}
	}()
	return errChan
}

func (tm *taxMarginsCa) getRegionalBrackets(ctx context.Context, input marginDS.Input) <-chan error {
	errChan := make(chan error)
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
