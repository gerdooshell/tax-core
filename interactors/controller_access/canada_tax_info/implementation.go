package canadaTaxInfo

import (
	"context"
	"fmt"
	"sync"

	sharedCredits "github.com/gerdooshell/tax-core/interactors/internal/shared_credits"
	sharedDeductions "github.com/gerdooshell/tax-core/interactors/internal/shared_deductions"
	taxCalculator "github.com/gerdooshell/tax-core/interactors/internal/tax_calculator"
	"github.com/gerdooshell/tax-core/library/mathHelper"
	"github.com/gerdooshell/tax-core/library/region/canada"
)

func NewCanadaTaxInfo() CanadaTaxInfo {
	return &taxInfoImpl{
		taxCalculator:        taxCalculator.NewTaxInteractor(),
		creditsCalculator:    sharedCredits.NewTaxCreditInteractor(),
		deductionsCalculator: sharedDeductions.NewTaxDeductionInteractor(),
	}
}

type taxInfoImpl struct {
	taxCalculator        taxCalculator.TaxBracketsInteractor
	creditsCalculator    sharedCredits.TaxCreditInteractor
	deductionsCalculator sharedDeductions.TaxDeductionInteractor
}

func (t *taxInfoImpl) CalculateLegacyTax(ctx context.Context, input *Input) (out Output, err error) {
	if input == nil {
		err = fmt.Errorf("CalculateLegacyTax error: nil input")
		return
	}
	var totalFederalCredits, totalRegionalCredits float64
	wg := sync.WaitGroup{}
	wg.Add(2)
	go t.calculateLegacyDeductionAndTax(ctx, input, &wg, &out, &err)

	go t.calculateLegacyCredits(ctx, input, &wg, &out, &err, &totalFederalCredits, &totalRegionalCredits)
	wg.Wait()
	out.FederalPayableTax = max(mathHelper.RoundFloat64(out.FederalTotalTax-totalFederalCredits, 2), 0)
	out.RegionalPayableTax = max(mathHelper.RoundFloat64(out.RegionalTotalTax-totalRegionalCredits, 2), 0)
	return
}

func (t *taxInfoImpl) calculateLegacyDeductionAndTax(ctx context.Context, input *Input, wg *sync.WaitGroup, out *Output, err *error) {
	defer wg.Done()
	deductionsChan := t.deductionsCalculator.GetTaxDeductions(ctx, input.Year, input.TotalIncome)
	deductions := <-deductionsChan
	if deductions.Err != nil {
		*err = deductions.Err
		return
	}
	taxableIncome := input.TotalIncome - deductions.CPPFirstAdditionalEmployee - deductions.CPPSecondAdditionalEmployee
	taxRegionalChan := t.taxCalculator.GetTotalTax(ctx, input.Year, input.Province, taxableIncome)
	taxFederalChan := t.taxCalculator.GetTotalTax(ctx, input.Year, canada.Federal, taxableIncome)
	out.TaxDeductions.CPPFirstAdditional = deductions.CPPFirstAdditionalEmployee
	out.TaxDeductions.CPPSecondAdditional = deductions.CPPSecondAdditionalEmployee
	regionalTax := <-taxRegionalChan
	federalTax := <-taxFederalChan
	if regionalTax.Err != nil {
		*err = regionalTax.Err
		return
	}
	if federalTax.Err != nil {
		*err = federalTax.Err
		return
	}
	out.FederalTotalTax = federalTax.Value
	out.RegionalTotalTax = regionalTax.Value
}

func (t *taxInfoImpl) calculateLegacyCredits(ctx context.Context, input *Input, wg *sync.WaitGroup, out *Output, err *error, totalFederalCredits, totalRegionalCredits *float64) {
	defer wg.Done()
	creditsChan := t.creditsCalculator.GetTaxCredits(ctx, input.Year, input.Province, input.TotalIncome)
	creditsInfo := <-creditsChan
	if creditsInfo.Err != nil {
		*err = creditsInfo.Err
		return
	}
	creditsSumFederal := creditsInfo.EIPremiumEmployee + creditsInfo.CPPBasicEmployee + creditsInfo.FederalBPA + creditsInfo.CanadaEmploymentAmount
	creditsSumRegional := creditsInfo.EIPremiumEmployee + creditsInfo.CPPBasicEmployee + creditsInfo.RegionalBPA
	reducedTaxCreditRegional := t.taxCalculator.GetReducedTaxCredit(ctx, input.Year, input.Province, creditsSumRegional)
	reducedTaxCreditFederal := t.taxCalculator.GetReducedTaxCredit(ctx, input.Year, canada.Federal, creditsSumFederal)
	out.TaxCredits.EIPremium = creditsInfo.EIPremiumEmployee
	out.TaxCredits.CanadaPensionPlanBasic = creditsInfo.CPPBasicEmployee
	out.TaxCredits.FederalBasicPensionAmount = creditsInfo.FederalBPA
	out.TaxCredits.RegionalBasicPensionAmount = creditsInfo.RegionalBPA
	out.TaxCredits.CanadaEmploymentAmount = creditsInfo.CanadaEmploymentAmount
	totalFederalCreditsResp := <-reducedTaxCreditFederal
	if totalFederalCreditsResp.Err != nil {
		*err = totalFederalCreditsResp.Err
		return
	}
	totalRegionalCreditsResp := <-reducedTaxCreditRegional
	if totalRegionalCreditsResp.Err != nil {
		*err = totalRegionalCreditsResp.Err
		return
	}
	*totalFederalCredits = mathHelper.RoundFloat64(totalFederalCreditsResp.Value, 2)
	*totalRegionalCredits = mathHelper.RoundFloat64(totalRegionalCreditsResp.Value, 2)
}
