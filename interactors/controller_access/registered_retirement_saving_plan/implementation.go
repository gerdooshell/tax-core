package rrspInfo

import (
	"context"
	"errors"
	marginCalculator "github.com/gerdooshell/tax-core/interactors/internal/margin_calculator"
	marginDS "github.com/gerdooshell/tax-core/interactors/internal/margin_calculator/data_structures"
	"github.com/gerdooshell/tax-core/interactors/internal/registered_retirement_savings_plan"
	sharedCredits "github.com/gerdooshell/tax-core/interactors/internal/shared_credits"
	sharedDeductions "github.com/gerdooshell/tax-core/interactors/internal/shared_deductions"
	taxCalculator "github.com/gerdooshell/tax-core/interactors/internal/tax_calculator"
	"github.com/gerdooshell/tax-core/library/mathHelper"
)

func NewRegisteredRetirementSavingPlan() RegisteredRetirementSavingPlan {
	return &rrspInfo{
		deductionsInteractor: sharedDeductions.NewTaxDeductionInteractor(),
		creditsInteractor:    sharedCredits.NewTaxCreditInteractor(),
		taxInteractor:        taxCalculator.NewTaxInteractor(),
		marginsInteractor:    marginCalculator.NewTaxMarginCalculator(),
		rrspInteractor:       rrspCalculator.NewRRSPInteractor(),
	}
}

type rrspInfo struct {
	deductionsInteractor sharedDeductions.TaxDeductionInteractor
	creditsInteractor    sharedCredits.TaxCreditInteractor
	taxInteractor        taxCalculator.TaxBracketsInteractor
	marginsInteractor    marginCalculator.TaxMarginsInteractor
	rrspInteractor       rrspCalculator.RRSPInteractor
}

func (r *rrspInfo) GetOptimalRRSPContributions(ctx context.Context, input *OptimalInput) ([]OptimalOutput, error) {
	outs := make([]OptimalOutput, 0, 10)
	if input == nil {
		return outs, errors.New("optimal rrsp error: nil input is passed")
	}
	rrspLimitsChan := r.rrspInteractor.GetRRSPMaxContribution(ctx, input.Year, input.TotalIncome)
	deductionsChan := r.deductionsInteractor.GetTaxDeductions(ctx, input.Year, input.TotalIncome)
	marginsChan := r.marginsInteractor.GetCombinedMarginalBrackets(ctx, marginDS.Input{Year: input.Year, Province: input.Province})
	deductions := <-deductionsChan
	if deductions.Err != nil {
		return outs, deductions.Err
	}
	calculatedDeductions := deductions.CPPFirstAdditionalEmployee + deductions.CPPSecondAdditionalEmployee
	preContribTaxChan := r.taxInteractor.GetPayableTaxGivenDeductions(ctx, input.Year, input.Province, input.TotalIncome, calculatedDeductions)
	totalDeductions := calculatedDeductions + input.ContributedRRSP
	margins := <-marginsChan
	if margins.Err != nil {
		return outs, margins.Err
	}
	taxableIncome := max(input.TotalIncome-totalDeductions, 0)
	taxableIncomes := make([]float64, 0, len(margins.Brackets))
	rrspLimits := <-rrspLimitsChan
	if rrspLimits.Err != nil {
		return outs, rrspLimits.Err
	}
	maxRRSPContribTaxableIncome := max(input.TotalIncome-rrspLimits.MaxContribution-calculatedDeductions, 0)
	isMaxContribAdded := false
	for _, br := range margins.Brackets {
		if br.Low <= 0 || br.Low >= taxableIncome {
			continue
		}
		if br.Low > maxRRSPContribTaxableIncome && !isMaxContribAdded {
			isMaxContribAdded = true
			taxableIncomes = append(taxableIncomes, maxRRSPContribTaxableIncome)
		}
		taxableIncomes = append(taxableIncomes, br.Low)
	}
	if !isMaxContribAdded {
		isMaxContribAdded = true
		taxableIncomes = append(taxableIncomes, maxRRSPContribTaxableIncome)
	}
	taxableIncomes = append(taxableIncomes, taxableIncome)
	var deduction float64
	taxChannels := make([]<-chan taxCalculator.TaxOutput, len(taxableIncomes))
	for i, tIncome := range taxableIncomes {
		deduction = input.TotalIncome - tIncome
		taxChannels[i] = r.taxInteractor.GetPayableTaxGivenDeductions(ctx, input.Year, input.Province, input.TotalIncome, deduction)
	}
	preContribTax := <-preContribTaxChan
	if preContribTax.Err != nil {
		return outs, preContribTax.Err
	}
	for i, taxChan := range taxChannels {
		tax := <-taxChan
		if tax.Err != nil {
			return outs, tax.Err
		}
		totalRRSP := mathHelper.RoundFloat64(input.TotalIncome-taxableIncomes[i]-calculatedDeductions, 2)
		outs = append(outs, OptimalOutput{
			TaxableIncome: mathHelper.RoundFloat64(taxableIncomes[i], 2),
			RRSP:          totalRRSP,
			PayableTax:    tax.Value,
			TaxReturn:     mathHelper.RoundFloat64(preContribTax.Value-tax.Value, 2),
			LeftRRSPRoom:  rrspLimits.MaxContribution - totalRRSP,
		})
	}
	return outs, nil
}
