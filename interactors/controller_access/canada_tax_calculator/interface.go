package canadaTaxCalculator

import "context"

type CanadaTaxCalculator interface {
	Calculate(context.Context, *Input) (Output, error)
}
