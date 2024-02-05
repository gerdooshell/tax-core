package dataAccess

import (
	"context"
	"github.com/gerdooshell/tax-core/entities/canada/federal/credits"
)

type FederalTaxData interface {
	AllCanadaTaxData
	GetFederalBPA(ctx context.Context, year int) (<-chan credits.BasicPersonalAmount, <-chan error)
	GetCEA(ctx context.Context, year int) (<-chan credits.CanadaEmploymentAmount, <-chan error)
}
