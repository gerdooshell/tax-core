package dataAccess

import (
	"context"
	"github.com/gerdooshell/tax-core/entities/canada/shared"
	"github.com/gerdooshell/tax-core/library/region/canada"
)

type TaxMargin interface {
	GetAlbertaTaxBrackets(ctx context.Context, year int) (<-chan []shared.TaxBracket, <-chan error)
	GetFederalTaxBrackets(ctx context.Context, year int) (<-chan []shared.TaxBracket, <-chan error)
	GetBCTaxBrackets(ctx context.Context, year int) (<-chan []shared.TaxBracket, <-chan error)
	SaveMarginalTaxBrackets(ctx context.Context, province canada.Province, year int, brackets []shared.TaxBracket) (<-chan bool, <-chan error)
}
