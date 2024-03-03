package dataAccess

import (
	"context"
	sharedEntities "github.com/gerdooshell/tax-core/entities/canada/shared"
)

type CPPDataOut struct {
	CanadaPensionPlan sharedEntities.CanadaPensionPlan
	Err               error
}

type CanadaPensionPlanData interface {
	GetCPP(ctx context.Context, year int) <-chan CPPDataOut
}
