package federalBPA

import (
	"context"
	"errors"
	"time"

	dataProvider "github.com/gerdooshell/tax-core/data-access"
	dataAccess "github.com/gerdooshell/tax-core/interactors/data_access"
)

type FederalBasicPersonalAmountInteractor interface {
	GetFederalBPA(ctx context.Context, year int, totalIncome float64) <-chan FederalBasicPersonalAmountOutput
}

type FederalBasicPersonalAmountOutput struct {
	Value float64
	Err   error
}

func NewFederalBasicPersonalAmountInteractor() FederalBasicPersonalAmountInteractor {
	return &federalBPAImpl{
		dataProvider: dataProvider.GetDataProviderServiceInstance(),
		timeout:      time.Second * 10,
	}
}

type federalBPAImpl struct {
	dataProvider dataAccess.FederalBPAData
	timeout      time.Duration
}

func (f *federalBPAImpl) GetFederalBPA(ctx context.Context, year int, totalIncome float64) <-chan FederalBasicPersonalAmountOutput {
	out := make(chan FederalBasicPersonalAmountOutput, 1)
	go func() {
		defer close(out)
		federalBPAOutput := FederalBasicPersonalAmountOutput{}
		defer func() { out <- federalBPAOutput }()
		bpaChan := f.dataProvider.GetFederalBPA(ctx, year)
		select {
		case <-time.After(f.timeout):
			federalBPAOutput.Err = errors.New("get federal bpa timed out")
			return
		case federalBPA := <-bpaChan:
			if federalBPAOutput.Err = federalBPA.Err; federalBPAOutput.Err != nil {
				return
			}
			if federalBPAOutput.Err = federalBPA.BasicPersonalAmount.Calculate(totalIncome); federalBPAOutput.Err != nil {
				return
			}
			federalBPAOutput.Value = federalBPA.BasicPersonalAmount.GetValue()
		}
	}()
	return out
}
