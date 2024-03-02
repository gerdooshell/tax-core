package dataAccess

import (
	"context"
)

type FederalBPAData interface {
	GetFederalBPA(ctx context.Context, year int) <-chan FederalBPADataOut
}
