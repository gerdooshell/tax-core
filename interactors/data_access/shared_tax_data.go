package dataAccess

import (
	"context"
	"github.com/gerdooshell/tax-core/entities/canada/shared"
	"github.com/gerdooshell/tax-core/library/region/canada"
)

type AllCanadaTaxData interface {
	GetTaxBrackets(ctx context.Context, year int, province canada.Province) (<-chan []shared.TaxBracket, <-chan error)
	GetCPP(ctx context.Context, year int) (<-chan shared.CanadaPensionPlan, <-chan error)
	GetEIPremium(ctx context.Context, year int) (<-chan shared.EmploymentInsurancePremium, <-chan error)
}
