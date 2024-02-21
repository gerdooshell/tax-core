package sharedDeductions

import (
	"context"
	"github.com/gerdooshell/tax-core/interactors/internal/canada_pension_plan"
)

type TaxDeductionInteractor interface {
	GetTaxDeductions(ctx context.Context, year int, totalIncome float64) <-chan TaxDeductionsOutput
}

type TaxDeductionsOutput struct {
	CPPFirstAdditionalEmployee      float64
	CPPFirstAdditionalEmployer      float64
	CPPFirstAdditionalSelfEmployed  float64
	CPPSecondAdditionalEmployee     float64
	CPPSecondAdditionalEmployer     float64
	CPPSecondAdditionalSelfEmployed float64
	Err                             error
}

func NewTaxDeductionInteractor() TaxDeductionInteractor {
	return &taxDeductionsImpl{
		cpp: cppCalculator.NewCanadaPensionPlanInteractor(),
	}
}

type taxDeductionsImpl struct {
	cpp cppCalculator.CanadaPensionPlanInteractor
}

func (t *taxDeductionsImpl) GetTaxDeductions(ctx context.Context, year int, totalIncome float64) <-chan TaxDeductionsOutput {
	out := make(chan TaxDeductionsOutput, 1)
	go func() {
		defer close(out)
		taxDeductionsOutput := TaxDeductionsOutput{}
		defer func() { out <- taxDeductionsOutput }()
		canadaPensionPlan := <-t.cpp.GetCPPContribution(ctx, year, totalIncome)
		if canadaPensionPlan.Err != nil {
			taxDeductionsOutput.Err = canadaPensionPlan.Err
			return
		}
		taxDeductionsOutput.CPPFirstAdditionalEmployee = canadaPensionPlan.EmployeeFirstAdditional
		taxDeductionsOutput.CPPSecondAdditionalEmployee = canadaPensionPlan.EmployeeSecondAdditional
		taxDeductionsOutput.CPPFirstAdditionalEmployer = canadaPensionPlan.EmployerFirstAdditional
		taxDeductionsOutput.CPPSecondAdditionalEmployer = canadaPensionPlan.EmployerSecondAdditional
		taxDeductionsOutput.CPPFirstAdditionalSelfEmployed = canadaPensionPlan.SelfEmployedFirstAdditional
		taxDeductionsOutput.CPPSecondAdditionalSelfEmployed = canadaPensionPlan.SelfEmployedSecondAdditional
	}()
	return out
}
