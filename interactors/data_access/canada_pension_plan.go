package dataAccess

import (
	"context"
	"github.com/gerdooshell/tax-core/entities/canada/shared"
)

type CanadaPensionPlanData interface {
	GetCPP(ctx context.Context, year int) (<-chan shared.CanadaPensionPlan, <-chan error)
}
