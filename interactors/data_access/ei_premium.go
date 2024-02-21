package dataAccess

import (
	"context"

	"github.com/gerdooshell/tax-core/entities/canada/shared"
)

type EIPremiumData interface {
	GetEIPremium(ctx context.Context, year int) (<-chan shared.EmploymentInsurancePremium, <-chan error)
}
