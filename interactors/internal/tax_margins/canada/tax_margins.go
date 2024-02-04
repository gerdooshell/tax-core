package canadaTaxMrgins

import (
	"context"
	dataProvider "github.com/gerdooshell/tax-core/data-access"
	"github.com/gerdooshell/tax-core/entities/canada/shared"
	dataAccess "github.com/gerdooshell/tax-core/interactors/data_access"
	marginDS "github.com/gerdooshell/tax-core/interactors/internal/tax_margins/canada/data_structures"
	"github.com/gerdooshell/tax-core/library/region/canada"
)

type TaxMarginsCa interface {
	GetAllMarginalBrackets(ctx context.Context, input marginDS.Input) (out marginDS.Output, err error)
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

func (tm *taxMarginsCa) GetAllMarginalBrackets(ctx context.Context, input marginDS.Input) (out marginDS.Output, err error) {
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
	out.Brackets = tm.marginalTaxBrackets.GetMargins()
	return
}

func (tm *taxMarginsCa) getFederalBrackets(ctx context.Context, input marginDS.Input) <-chan error {
	errChan := make(chan error)
	go func() {
		defer close(errChan)
		var err error
		defer func() { errChan <- err }()
		out, errOut := tm.dataProvider.GetFederalTaxBrackets(ctx, input.Year)
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
		var out <-chan []shared.TaxBracket
		var errOut <-chan error
		switch input.Province {
		case canada.BritishColumbia:
			out, errOut = tm.dataProvider.GetBCTaxBrackets(ctx, input.Year)
		case canada.Alberta:
			out, errOut = tm.dataProvider.GetAlbertaTaxBrackets(ctx, input.Year)
		}
		select {
		case tm.regionalTaxBrackets = <-out:
		case err = <-errOut:
		}
	}()
	return errChan
}
