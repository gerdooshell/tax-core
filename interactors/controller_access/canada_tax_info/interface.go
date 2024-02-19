package canadaTaxInfo

import "context"

type CanadaTaxInfo interface {
	CalculateLegacyTax(ctx context.Context, input *Input) (Output, error)
}
