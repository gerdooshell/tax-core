package dataAccess

import (
	"context"
	sharedEntities "github.com/gerdooshell/tax-core/entities/canada/shared"
)

type RRSPData interface {
	GetRRSP(ctx context.Context, year int) (<-chan sharedEntities.RRSP, <-chan error)
}
