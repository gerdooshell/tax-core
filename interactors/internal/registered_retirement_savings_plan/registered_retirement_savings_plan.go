package rrspCalculator

import (
	"context"
	"errors"
	"time"

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
		timeout:      time.Second * 10,
	}
}

type rrspImpl struct {
	dataProvider dataAccess.RRSPData
	timeout      time.Duration
}

func (rrsp *rrspImpl) GetRRSPMaxContribution(ctx context.Context, year int, totalIncome float64) <-chan RRSPContributionOutput {
	out := make(chan RRSPContributionOutput, 1)
	go func() {
		defer close(out)
		rrspOutput := RRSPContributionOutput{}
		defer func() { out <- rrspOutput }()
		rrspChan := rrsp.dataProvider.GetRRSP(ctx, year)
		select {
		case rrspDataOut := <-rrspChan:
			if rrspDataOut.Err != nil {
				rrspOutput.Err = rrspDataOut.Err
				return
			}
			if rrspOutput.Err = rrspDataOut.RRSP.CalculateMaxContribution(totalIncome); rrspOutput.Err != nil {
				return
			}
			rrspOutput.MaxContribution = rrspDataOut.RRSP.GetContribution()
		case <-time.After(rrsp.timeout):
			rrspOutput.Err = errors.New("get rrsp data timed out")
		}
	}()
	return out
}
