package sharedCredits

import (
	"context"

	"github.com/gerdooshell/tax-core/interactors/internal/canada_employment_amount"
	"github.com/gerdooshell/tax-core/interactors/internal/canada_pension_plan"
	"github.com/gerdooshell/tax-core/interactors/internal/ei_premium"
	"github.com/gerdooshell/tax-core/interactors/internal/federal_basic_personal_amount"
	"github.com/gerdooshell/tax-core/interactors/internal/regional_basic_personal_amount"
	"github.com/gerdooshell/tax-core/library/region/canada"
)

type TaxCreditInteractor interface {
	GetTaxCredits(ctx context.Context, year int, province canada.Province, totalIncome float64) <-chan TaxCreditsOutput
}

func NewTaxCreditInteractor() TaxCreditInteractor {
	return &taxCreditImpl{
		ceaCalculator:         ceaCalculator.NewCanadaEmploymentAmountInteractor(),
		eipCalculator:         eipCalculator.NewEIPremiumInteractor(),
		cppCalculator:         cppCalculator.NewCanadaPensionPlanInteractor(),
		fedBPACalculator:      federalBPA.NewFederalBasicPersonalAmountInteractor(),
		regionalBPACalculator: regionalBPA.NewRegionalBasicPersonalAmountInteractor(),
	}
}

type TaxCreditsOutput struct {
	CPPBasicEmployee       float64
	CPPBasicEmployer       float64
	CPPBasicSelfEmployed   float64
	EIPremiumEmployee      float64
	EIPremiumEmployer      float64
	CanadaEmploymentAmount float64
	FederalBPA             float64
	RegionalBPA            float64
	Err                    error
}

type taxCreditImpl struct {
	ceaCalculator         ceaCalculator.CanadaEmploymentAmountInteractor
	eipCalculator         eipCalculator.EIPremiumInteractor
	cppCalculator         cppCalculator.CanadaPensionPlanInteractor
	fedBPACalculator      federalBPA.FederalBasicPersonalAmountInteractor
	regionalBPACalculator regionalBPA.RegionalBasicPersonalAmountInteractor
}

func (t *taxCreditImpl) GetTaxCredits(ctx context.Context, year int, province canada.Province, totalIncome float64) <-chan TaxCreditsOutput {
	out := make(chan TaxCreditsOutput)
	go func() {
		defer close(out)
		taxCreditsOutput := TaxCreditsOutput{}
		defer func() { out <- taxCreditsOutput }()
		ceaChan := t.ceaCalculator.GetCEACredit(ctx, year, totalIncome)
		eipChan := t.eipCalculator.GetEIContribution(ctx, year, totalIncome)
		cppChan := t.cppCalculator.GetCPPContribution(ctx, year, totalIncome)
		fedBPAChan := t.fedBPACalculator.GetFederalBPA(ctx, year, totalIncome)
		regionalBPAChan := t.regionalBPACalculator.GetRegionalBPA(ctx, year, totalIncome, province)
		canadaEmploymentAmount := <-ceaChan
		if canadaEmploymentAmount.Err != nil {
			taxCreditsOutput.Err = canadaEmploymentAmount.Err
			return
		}
		eiPremium := <-eipChan
		if eiPremium.Err != nil {
			taxCreditsOutput.Err = eiPremium.Err
			return
		}

		canadaPensionPlan := <-cppChan
		if canadaPensionPlan.Err != nil {
			taxCreditsOutput.Err = canadaPensionPlan.Err
			return
		}

		federalBasicPensionAmount := <-fedBPAChan
		if federalBasicPensionAmount.Err != nil {
			taxCreditsOutput.Err = federalBasicPensionAmount.Err
			return
		}
		regionalBasicPersonalAmount := <-regionalBPAChan
		if regionalBasicPersonalAmount.Err != nil {
			taxCreditsOutput.Err = regionalBasicPersonalAmount.Err
			return
		}
		taxCreditsOutput.CanadaEmploymentAmount = canadaEmploymentAmount.Value
		taxCreditsOutput.EIPremiumEmployee = eiPremium.Employee
		taxCreditsOutput.EIPremiumEmployer = eiPremium.Employer
		taxCreditsOutput.CPPBasicEmployee = canadaPensionPlan.EmployeeBasic
		taxCreditsOutput.CPPBasicEmployer = canadaPensionPlan.EmployerBasic
		taxCreditsOutput.CPPBasicSelfEmployed = canadaPensionPlan.SelfEmployedBasic
		taxCreditsOutput.FederalBPA = federalBasicPensionAmount.Value
		taxCreditsOutput.RegionalBPA = regionalBasicPersonalAmount.Value
	}()
	return out
}
