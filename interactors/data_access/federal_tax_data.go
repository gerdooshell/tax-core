package dataAccess

import (
	"context"
	"github.com/gerdooshell/tax-core/entities/canada/federal/credits"
	"github.com/gerdooshell/tax-core/entities/canada/shared"
)

type FederalTaxData interface {
	AllCanadaTaxData
	GetFederalBPA(ctx context.Context, year int) (<-chan credits.BasicPersonalAmount, <-chan error)
	GetFederalTaxBrackets(ctx context.Context, year int) (<-chan []shared.TaxBracket, <-chan error)
	GetCEA(ctx context.Context, year int) (<-chan credits.CanadaEmploymentAmount, <-chan error)
}
