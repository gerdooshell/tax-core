package rrspCalculator

import (
	"context"
	dataProvider "github.com/gerdooshell/tax-core/data-access"
	dataAccess "github.com/gerdooshell/tax-core/interactors/data_access"
)

type RRSPInteractor interface {
}

type RRSPOutput struct {
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

func (rrsp *rrspImpl) GetRRSPLimits(ctx context.Context, year int, totalIncome float64) <-chan RRSPOutput {
	out := make(chan RRSPOutput)
	go func() {
		defer close(out)
		rrspOutput := RRSPOutput{}
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
