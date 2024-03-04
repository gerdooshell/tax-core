package dataAccess

import (
	"context"
	federalCredits "github.com/gerdooshell/tax-core/entities/canada/federal/credits"
)

type CEADataOut struct {
	CanadaEmploymentAmount federalCredits.CanadaEmploymentAmount
	Err                    error
}

type CanadaEmploymentAmountData interface {
	GetCEA(ctx context.Context, year int) <-chan CEADataOut
}
