package dataAccess

import (
	"context"
	sharedEntities "github.com/gerdooshell/tax-core/entities/canada/shared"
)

type RRSPDataOut struct {
	RRSP sharedEntities.RRSP
	Err  error
}

type RRSPData interface {
	GetRRSP(ctx context.Context, year int) <-chan RRSPDataOut
}
