package canadaTaxMarginCalculator

import (
	"context"
)

type CanadaTaxMarginCalculator interface {
	GetAllMarginalBrackets(ctx context.Context, input *Input) (out Output, err error)
}
