package federalBPA

import (
	"context"
	"fmt"

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
	}
}

type federalBPAImpl struct {
	dataProvider dataAccess.FederalBPAData
}

func (f *federalBPAImpl) GetFederalBPA(ctx context.Context, year int, totalIncome float64) <-chan FederalBasicPersonalAmountOutput {
	out := make(chan FederalBasicPersonalAmountOutput)
	go func() {
		defer close(out)
		federalBPAOutput := FederalBasicPersonalAmountOutput{}
		defer func() { out <- federalBPAOutput }()
		bpaChan, errChan := f.dataProvider.GetFederalBPA(ctx, year)
		select {
		case federalBPAOutput.Err = <-errChan:
			return
		case federalBPA := <-bpaChan:
			if federalBPAOutput.Err = federalBPA.Calculate(totalIncome); federalBPAOutput.Err != nil {
				return
			}
			federalBPAOutput.Value = federalBPA.GetValue()
		case <-ctx.Done():
			federalBPAOutput.Err = fmt.Errorf("processing federal bpa canceled")
		}
	}()
	return out
}
