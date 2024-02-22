package canadaTaxMarginCalculator

import (
	"context"
	"errors"
	canadaTaxMrgins "github.com/gerdooshell/tax-core/interactors/internal/margin_calculator"
	"github.com/gerdooshell/tax-core/interactors/internal/margin_calculator/data_structures"
)

type canadaTaxMarginCalculator struct {
}

func NewCanadaTaxMarginCalculator() CanadaTaxMarginCalculator {
	return &canadaTaxMarginCalculator{}
}

func (c canadaTaxMarginCalculator) GetAllMarginalBrackets(ctx context.Context, input *Input) (out Output, err error) {
	if input == nil {
		err = errors.New("null marginal input")
		return
	}
	margin := canadaTaxMrgins.NewTaxMarginCalculator()
	brackets := <-margin.GetCombinedMarginalBrackets(ctx, marginDS.Input{Year: input.Year, Province: input.Province})
	if brackets.Err != nil {
		return
	}
	out.MarginalBrackets = brackets.Brackets
	return
}
