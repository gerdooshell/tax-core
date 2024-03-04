package dataAccess

import (
	"context"
	sharedEntities "github.com/gerdooshell/tax-core/entities/canada/shared"
)

type EIPremiumDataOut struct {
	EmploymentInsurancePremium sharedEntities.EmploymentInsurancePremium
	Err                        error
}

type EIPremiumData interface {
	GetEIPremium(ctx context.Context, year int) <-chan EIPremiumDataOut
}
