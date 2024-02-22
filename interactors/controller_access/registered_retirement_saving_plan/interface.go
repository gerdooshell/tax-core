package rrspInfo

import "context"

type RegisteredRetirementSavingPlan interface {
	GetOptimalRRSPContributions(ctx context.Context, input *OptimalInput) ([]OptimalOutput, error)
}
