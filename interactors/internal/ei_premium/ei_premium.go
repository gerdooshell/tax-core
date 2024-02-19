package eipCalculator

import (
	"context"
	"fmt"

	dataProvider "github.com/gerdooshell/tax-core/data-access"
	sharedEntities "github.com/gerdooshell/tax-core/entities/canada/shared"
	dataAccess "github.com/gerdooshell/tax-core/interactors/data_access"
)

type EIPremiumInteractor interface {
	GetEIContribution(ctx context.Context, year int, totalIncome float64) <-chan EIPremiumOutput
}

func NewEIPremiumInteractor() EIPremiumInteractor {
	return &eiPremiumImpl{
		dataProvider: dataProvider.GetDataProviderServiceInstance(),
	}
}

type eiPremiumImpl struct {
	dataProvider dataAccess.EIPremiumData
	eiEntity     sharedEntities.EmploymentInsurancePremium
}

type EIPremiumOutput struct {
	Employee float64
	Employer float64
	Err      error
}

func (eip *eiPremiumImpl) GetEIContribution(ctx context.Context, year int, totalIncome float64) <-chan EIPremiumOutput {
	out := make(chan EIPremiumOutput)
	eiChan, errChan := eip.dataProvider.GetEIPremium(ctx, year)

	go func() {
		defer close(out)
		var eipOutput EIPremiumOutput
		defer func() { out <- eipOutput }()
		select {
		case eipOutput.Err = <-errChan:
			return
		case eip.eiEntity = <-eiChan:
			if eipOutput.Err = eip.eiEntity.Calculate(totalIncome); eipOutput.Err != nil {
				return
			}
			eipOutput.Employee = eip.eiEntity.GetEIEmployee()
			eipOutput.Employer = eip.eiEntity.GetEIEmployer()
		case <-ctx.Done():
			eipOutput.Err = fmt.Errorf("processing ei premium canceled")
		}
	}()
	return out
}
