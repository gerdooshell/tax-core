package dataAccess

import (
	"context"
	"github.com/gerdooshell/tax-core/entities/canada/shared"
	"github.com/gerdooshell/tax-core/library/region/canada"
)

type TaxMargin interface {
	SaveMarginalTaxBrackets(ctx context.Context, province canada.Province, year int, brackets []shared.TaxBracket) (<-chan bool, <-chan error)
	GetCombinedMarginalBrackets(ctx context.Context, year int, province canada.Province) (<-chan []shared.TaxBracket, <-chan error)
	GetTaxBrackets(ctx context.Context, year int, province canada.Province) (<-chan []shared.TaxBracket, <-chan error)
}
