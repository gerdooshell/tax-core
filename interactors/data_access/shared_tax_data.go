package dataAccess

import (
	"context"
	"github.com/gerdooshell/tax-core/entities/canada/shared"
)

type AllCanadaTaxData interface {
	GetCPP(ctx context.Context, year int) (<-chan shared.CanadaPensionPlan, <-chan error)
	GetEIPremium(ctx context.Context, year int) (<-chan shared.EmploymentInsurancePremium, <-chan error)
}
