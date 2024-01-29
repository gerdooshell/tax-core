package dataAccess

import (
	"context"
	"github.com/gerdooshell/tax-core/entities/canada/alberta/credits"
	"github.com/gerdooshell/tax-core/entities/canada/shared"
)

type AlbertaTaxData interface {
	AllCanadaTaxData
	GetAlbertaBPA(ctx context.Context, year int) (<-chan credits.BasicPersonalAmount, <-chan error)
	GetAlbertaTaxBrackets(ctx context.Context, year int) (<-chan []shared.TaxBracket, <-chan error)
}
