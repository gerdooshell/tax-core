package dataAccess

import (
	"context"
	"github.com/gerdooshell/tax-core/entities/canada/alberta/credits"
)

type AlbertaTaxData interface {
	AllCanadaTaxData
	GetAlbertaBPA(ctx context.Context, year int) (<-chan credits.BasicPersonalAmount, <-chan error)
}
