package ceaCalculator

import (
	"context"
	"errors"
	dataProvider "github.com/gerdooshell/tax-core/data-access"
	federalEntities "github.com/gerdooshell/tax-core/entities/canada/federal/credits"
	dataAccess "github.com/gerdooshell/tax-core/interactors/data_access"
	"time"
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
		timeout:      time.Second * 10,
	}
}

type canadaEmploymentAmountImpl struct {
	dataProvider dataAccess.CanadaEmploymentAmountData
	ceaEntity    federalEntities.CanadaEmploymentAmount
	timeout      time.Duration
}

func (cea *canadaEmploymentAmountImpl) GetCEACredit(ctx context.Context, year int, totalIncome float64) <-chan CanadaEmploymentAmountOutput {
	out := make(chan CanadaEmploymentAmountOutput, 1)
	go func() {
		defer close(out)
		var ceaOutput CanadaEmploymentAmountOutput
		defer func() { out <- ceaOutput }()
		ceaChan := cea.dataProvider.GetCEA(ctx, year)
		select {

		case ceaDataOut := <-ceaChan:
			if ceaDataOut.Err != nil {
				ceaOutput.Err = ceaDataOut.Err
				return
			}
			cea.ceaEntity = ceaDataOut.CanadaEmploymentAmount
			if ceaOutput.Err = cea.ceaEntity.Calculate(totalIncome); ceaOutput.Err != nil {
				return
			}
			ceaOutput.Value = cea.ceaEntity.GetEmployeeValue()
		case <-time.After(cea.timeout):
			ceaOutput.Err = errors.New("get cea data timed out")
		}
	}()
	return out
}
