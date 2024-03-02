package dataAccess

import (
	"context"
)

type EIPremiumData interface {
	GetEIPremium(ctx context.Context, year int) <-chan EIPremiumDataOut
}
