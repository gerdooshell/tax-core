package rrspCalculator

import (
	"context"
	dataProvider "github.com/gerdooshell/tax-core/data-access"
	dataAccess "github.com/gerdooshell/tax-core/interactors/data_access"
)

type RRSPInteractor interface {
	GetRRSPMaxContribution(ctx context.Context, year int, totalIncome float64) <-chan RRSPContributionOutput
}

type RRSPContributionOutput struct {
	MaxContribution float64
	Err             error
}

func NewRRSPInteractor() RRSPInteractor {
	return &rrspImpl{
		dataProvider: dataProvider.GetDataProviderServiceInstance(),
	}
}

type rrspImpl struct {
	dataProvider dataAccess.RRSPData
}

func (rrsp *rrspImpl) GetRRSPMaxContribution(ctx context.Context, year int, totalIncome float64) <-chan RRSPContributionOutput {
	out := make(chan RRSPContributionOutput, 1)
	go func() {
		defer close(out)
		rrspOutput := RRSPContributionOutput{}
		defer func() { out <- rrspOutput }()
		rrspChan, errChan := rrsp.dataProvider.GetRRSP(ctx, year)
		select {
		case rrspOutput.Err = <-errChan:
			return
		case rrspEntity := <-rrspChan:
			if rrspOutput.Err = rrspEntity.CalculateMaxContribution(totalIncome); rrspOutput.Err != nil {
				return
			}
			rrspOutput.MaxContribution = rrspEntity.GetContribution()
		}
	}()
	return out
}
