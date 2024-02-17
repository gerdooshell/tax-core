package canadaTaxMarginCalculator

import (
	"context"
	"errors"

	canadaTaxMrgins "github.com/gerdooshell/tax-core/interactors/internal/tax_margins/canada"
	marginDS "github.com/gerdooshell/tax-core/interactors/internal/tax_margins/canada/data_structures"
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
	margin := canadaTaxMrgins.NewTaxMarginCa()
	brackets, err := margin.GetCombinedMarginalBrackets(ctx, marginDS.Input{Year: input.Year, Province: input.Province})
	if err != nil {
		return
	}
	out.MarginalBrackets = brackets.Brackets
	return
}
