package eipCalculator

import (
	"context"
	"errors"
	"time"

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
		timeout:      time.Second * 10,
	}
}

type eiPremiumImpl struct {
	dataProvider dataAccess.EIPremiumData
	eiEntity     sharedEntities.EmploymentInsurancePremium
	timeout      time.Duration
}

type EIPremiumOutput struct {
	Employee float64
	Employer float64
	Err      error
}

func (eip *eiPremiumImpl) GetEIContribution(ctx context.Context, year int, totalIncome float64) <-chan EIPremiumOutput {
	out := make(chan EIPremiumOutput, 1)

	go func() {
		defer close(out)
		var eipOutput EIPremiumOutput
		defer func() { out <- eipOutput }()
		eiChan := eip.dataProvider.GetEIPremium(ctx, year)
		select {
		case eiPremiumDataOut := <-eiChan:
			if eipOutput.Err = eiPremiumDataOut.Err; eipOutput.Err != nil {
				return
			}
			eip.eiEntity = eiPremiumDataOut.EmploymentInsurancePremium
			if eipOutput.Err = eip.eiEntity.Calculate(totalIncome); eipOutput.Err != nil {
				return
			}
			eipOutput.Employee = eip.eiEntity.GetEIEmployee()
			eipOutput.Employer = eip.eiEntity.GetEIEmployer()
		case <-time.After(eip.timeout):
			eipOutput.Err = errors.New("get ei premium contribution data timed out")
		}
	}()
	return out
}
