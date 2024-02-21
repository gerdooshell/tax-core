package dataAccess

import (
	"context"

	"github.com/gerdooshell/tax-core/entities/canada/federal/credits"
)

type CanadaEmploymentAmountData interface {
	GetCEA(ctx context.Context, year int) (<-chan credits.CanadaEmploymentAmount, <-chan error)
}
