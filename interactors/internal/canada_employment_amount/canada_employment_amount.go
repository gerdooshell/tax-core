package ceaCalculator

import (
	"context"

	dataProvider "github.com/gerdooshell/tax-core/data-access"
	federalEntities "github.com/gerdooshell/tax-core/entities/canada/federal/credits"
	dataAccess "github.com/gerdooshell/tax-core/interactors/data_access"
)

type CanadaEmploymentAmountInteractor interface {
	GetCEACredit(ctx context.Context, year int, totalIncome float64) <-chan CanadaEmploymentAmountOutput
}

type CanadaEmploymentAmountOutput struct {
	Value float64
	Err   error
}

func NewCanadaEmploymentAmountInteractor() CanadaEmploymentAmountInteractor {
	return &canadaEmploymentAmountImpl{
		dataProvider: dataProvider.GetDataProviderServiceInstance(),
	}
}

type canadaEmploymentAmountImpl struct {
	dataProvider dataAccess.CanadaEmploymentAmountData
	ceaEntity    federalEntities.CanadaEmploymentAmount
}

func (cea *canadaEmploymentAmountImpl) GetCEACredit(ctx context.Context, year int, totalIncome float64) <-chan CanadaEmploymentAmountOutput {
	out := make(chan CanadaEmploymentAmountOutput)
	ceaChan, errChan := cea.dataProvider.GetCEA(ctx, year)

	go func() {
		defer close(out)
		var ceaOutput CanadaEmploymentAmountOutput
		defer func() { out <- ceaOutput }()
		select {
		case ceaOutput.Err = <-errChan:
			return
		case cea.ceaEntity = <-ceaChan:
			if ceaOutput.Err = cea.ceaEntity.Calculate(totalIncome); ceaOutput.Err != nil {
				return
			}
			ceaOutput.Value = cea.ceaEntity.GetEmployeeValue()
		}
	}()
	return out
}
