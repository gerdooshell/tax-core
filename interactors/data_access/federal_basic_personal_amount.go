package dataAccess

import (
	"context"

	"github.com/gerdooshell/tax-core/entities/canada/federal/credits"
)

type FederalBPAData interface {
	GetFederalBPA(ctx context.Context, year int) (<-chan credits.BasicPersonalAmount, <-chan error)
}
