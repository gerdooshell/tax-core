package canadaTaxMarginCalculator

import (
	"context"
	canadaTaxMrgins "github.com/gerdooshell/tax-core/interactors/internal/tax_margins/canada"
	marginDS "github.com/gerdooshell/tax-core/interactors/internal/tax_margins/canada/data_structures"
)

type canadaTaxMarginCalculator struct {
}

func NewCanadaTaxMarginCalculator() CanadaTaxMarginCalculator {
	return &canadaTaxMarginCalculator{}
}

func (c canadaTaxMarginCalculator) GetAllMarginalBrackets(ctx context.Context, input *Input) (out Output, err error) {
	margin := canadaTaxMrgins.NewTaxMarginCa()
	brackets, err := margin.GetAllMarginalBrackets(ctx, marginDS.Input{Year: input.Year, Province: input.Province})
	if err != nil {
		return
	}
	out.MarginalBrackets = brackets.Brackets
	return
}
