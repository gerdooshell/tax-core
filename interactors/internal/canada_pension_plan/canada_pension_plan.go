package cppCalculator

import (
	"context"
	"errors"
	"time"

	dataProvider "github.com/gerdooshell/tax-core/data-access"
	sharedEntities "github.com/gerdooshell/tax-core/entities/canada/shared"
	dataAccess "github.com/gerdooshell/tax-core/interactors/data_access"
)

type CanadaPensionPlanInteractor interface {
	GetCPPContribution(ctx context.Context, year int, totalIncome float64) <-chan CanadaPensionPlanOutput
}

type CanadaPensionPlanOutput struct {
	EmployeeBasic                float64
	EmployerBasic                float64
	SelfEmployedBasic            float64
	EmployeeFirstAdditional      float64
	EmployerFirstAdditional      float64
	SelfEmployedFirstAdditional  float64
	EmployeeSecondAdditional     float64
	EmployerSecondAdditional     float64
	SelfEmployedSecondAdditional float64
	Err                          error
}

func NewCanadaPensionPlanInteractor() CanadaPensionPlanInteractor {
	return &canadaPensionPlanImpl{
		dataProvider: dataProvider.GetDataProviderServiceInstance(),
		timeout:      time.Second * 10,
	}
}

type canadaPensionPlanImpl struct {
	dataProvider dataAccess.CanadaPensionPlanData
	cppEntity    sharedEntities.CanadaPensionPlan
	timeout      time.Duration
}

func (c *canadaPensionPlanImpl) GetCPPContribution(ctx context.Context, year int, totalIncome float64) <-chan CanadaPensionPlanOutput {
	out := make(chan CanadaPensionPlanOutput, 1)

	go func() {
		defer close(out)
		var cppOutput CanadaPensionPlanOutput
		defer func() { out <- cppOutput }()
		cppChan := c.dataProvider.GetCPP(ctx, year)
		select {
		case cppDataOut := <-cppChan:
			if cppDataOut.Err != nil {
				cppOutput.Err = cppDataOut.Err
				return
			}
			c.cppEntity = cppDataOut.CanadaPensionPlan
			if cppOutput.Err = c.cppEntity.Calculate(totalIncome); cppOutput.Err != nil {
				return
			}
			cppOutput.EmployeeBasic = c.cppEntity.GetCPPBasicEmployee()
			cppOutput.EmployerBasic = c.cppEntity.GetCPPBasicEmployer()
			cppOutput.SelfEmployedBasic = c.cppEntity.GetCPPBasicSelfEmployed()
			cppOutput.EmployeeFirstAdditional = c.cppEntity.GetCPPFirstAdditionalEmployee()
			cppOutput.EmployerFirstAdditional = c.cppEntity.GetCPPFirstAdditionalEmployer()
			cppOutput.SelfEmployedFirstAdditional = c.cppEntity.GetCPPFirstAdditionalSelfEmployed()
			cppOutput.EmployeeSecondAdditional = c.cppEntity.GetCPPSecondAdditionalEmployee()
			cppOutput.EmployerSecondAdditional = c.cppEntity.GetCPPSecondAdditionalEmployer()
			cppOutput.SelfEmployedSecondAdditional = c.cppEntity.GetCPPSecondAdditionalSelfEmployed()
		case <-time.After(c.timeout):
			cppOutput.Err = errors.New("get cpp data timed out")
		}
	}()
	return out
}
