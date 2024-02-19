package dataAccess

import (
	"context"

	albertaCredits "github.com/gerdooshell/tax-core/entities/canada/alberta/credits"
	bcCredits "github.com/gerdooshell/tax-core/entities/canada/bc/credits"
)

type RegionalBPAData interface {
	GetAlbertaBPA(ctx context.Context, year int) (<-chan albertaCredits.BasicPersonalAmount, <-chan error)
	GetBCBPA(ctx context.Context, year int) (<-chan bcCredits.BasicPersonalAmount, <-chan error)
}
