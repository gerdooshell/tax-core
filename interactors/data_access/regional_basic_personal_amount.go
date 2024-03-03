package dataAccess

import (
	"context"

	albertaCredits "github.com/gerdooshell/tax-core/entities/canada/alberta/credits"
	bcCredits "github.com/gerdooshell/tax-core/entities/canada/bc/credits"
)

type BritishColumbiaBPADataOut struct {
	BasicPersonalAmount bcCredits.BasicPersonalAmount
	Err                 error
}

type AlbertaBPADataOut struct {
	BasicPersonalAmount albertaCredits.BasicPersonalAmount
	Err                 error
}

type RegionalBPAData interface {
	GetAlbertaBPA(ctx context.Context, year int) <-chan AlbertaBPADataOut
	GetBritishColumbiaBPA(ctx context.Context, year int) <-chan BritishColumbiaBPADataOut
}
