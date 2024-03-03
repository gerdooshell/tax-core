package dataAccess

import (
	"context"
	federalCredits "github.com/gerdooshell/tax-core/entities/canada/federal/credits"
)

type FederalBPADataOut struct {
	BasicPersonalAmount federalCredits.BasicPersonalAmount
	Err                 error
}

type FederalBPAData interface {
	GetFederalBPA(ctx context.Context, year int) <-chan FederalBPADataOut
}
