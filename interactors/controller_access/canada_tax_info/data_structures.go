package canadaTaxInfo

import "github.com/gerdooshell/tax-core/library/region/canada"

type Input struct {
	Province    canada.Province
	TotalIncome float64
	Year        int
}

type Output struct {
	TaxCredits         TaxCredits
	TaxDeductions      TaxDeductions
	FederalPayableTax  float64
	FederalTotalTax    float64
	RegionalPayableTax float64
	RegionalTotalTax   float64
}

type TaxCredits struct {
	FederalBasicPensionAmount  float64
	CanadaEmploymentAmount     float64
	RegionalBasicPensionAmount float64
	EIPremium                  float64
	CanadaPensionPlanBasic     float64
}

type TaxDeductions struct {
	CPPFirstAdditional  float64
	CPPSecondAdditional float64
}
