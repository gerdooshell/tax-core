package taxCalculator

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/gerdooshell/tax-core/controller/internal"
	"github.com/gerdooshell/tax-core/controller/internal/routes"
	restApi "github.com/gerdooshell/tax-core/controller/rest_api/handlers"
	canadaTaxInfo "github.com/gerdooshell/tax-core/interactors/controller_access/canada_tax_info"
	"github.com/gerdooshell/tax-core/library/region/canada"

	"github.com/gorilla/mux"
)

type taxCalculator struct {
}

func NewTaxCalculatorController() restApi.Handler {
	return &taxCalculator{}
}

type State struct {
	input  *canadaTaxInfo.Input
	apiKey string
}

func (tc *taxCalculator) URL() string {
	return routes.TaxCalculator
}

func (tc *taxCalculator) Methods() []string {
	return []string{http.MethodGet}
}

func (tc *taxCalculator) Authorize() error {
	return nil
}

func (tc *taxCalculator) ParseArgs(r *http.Request) (*http.Request, error) {
	state := State{}
	pathVars := mux.Vars(r)
	provinceStr, ok := pathVars["province"]
	if !ok {
		return nil, fmt.Errorf("province is not provided")
	}
	var err error
	province, err := canada.GetProvinceFromString(provinceStr)
	if err != nil {
		return nil, err
	}
	yearStr, ok := pathVars["year"]
	if !ok {
		return nil, fmt.Errorf("year is not provided")
	}
	year, err := strconv.Atoi(yearStr)
	if err != nil {
		return nil, err
	}
	incomeStr, ok := pathVars["income"]
	if !ok {
		return nil, fmt.Errorf("income is not provided")
	}
	income, err := strconv.ParseFloat(incomeStr, 64)
	if err != nil {
		return nil, err
	}
	input := canadaTaxInfo.Input{
		Province:    province,
		Year:        year,
		TotalIncome: income,
	}
	if err != nil {
		return nil, err
	}
	state.apiKey = r.Header.Get(internal.APIKyeNameID)
	state.input = &input
	if err := validateInput(state); err != nil {
		return nil, err
	}
	ctx := context.WithValue(r.Context(), "state", state)
	r = r.WithContext(ctx)
	return r, nil
}

func (tc *taxCalculator) Process(r *http.Request) *http.Response {
	resp := new(http.Response)
	state := r.Context().Value("state").(State)
	ctx := r.Context()
	calculator := canadaTaxInfo.NewCanadaTaxInfo()
	out, err := calculator.CalculateLegacyTax(ctx, state.input)
	if err != nil {
		resp.Body = io.NopCloser(bytes.NewReader([]byte(err.Error())))
		resp.StatusCode = http.StatusInternalServerError
		return resp
	}
	respBody := NewResponseBodyModelFrom(out)
	respBodyByte, err := json.Marshal(respBody)
	if err != nil {
		resp.Body = io.NopCloser(bytes.NewReader([]byte(err.Error())))
		resp.StatusCode = http.StatusInternalServerError
		return resp
	}
	resp.Body = io.NopCloser(bytes.NewReader(respBodyByte))

	resp.StatusCode = http.StatusOK
	return resp
}

func validateInput(state State) error {
	if state.input.Year <= 0 {
		return fmt.Errorf("invalid year \"%v\"", state.input.Year)
	}
	if state.input.TotalIncome <= 0 {
		return fmt.Errorf("invalid income \"%v\"", state.input.TotalIncome)
	}
	if state.input.Province == canada.Federal {
		return fmt.Errorf("invalid province \"%v\"", state.input.Province)
	}
	return nil
}

type RequestBodyModel struct {
	Province string  `json:"province"`
	Income   float64 `json:"income"`
	Year     int     `json:"year"`
}

type ResponseBodyModel struct {
	FederalPayableTax  float64           `json:"federal_payable_tax"`
	FederalTotalTax    float64           `json:"federal_total_tax"`
	RegionalPayableTax float64           `json:"regional_payable_tax"`
	RegionalTotalTax   float64           `json:"regional_total_tax"`
	TaxCredits         TaxCreditModel    `json:"tax_credits"`
	TaxDeductions      TaxDeductionModel `json:"tax_deductions"`
}

func NewResponseBodyModelFrom(out canadaTaxInfo.Output) ResponseBodyModel {
	return ResponseBodyModel{
		FederalPayableTax:  out.FederalPayableTax,
		FederalTotalTax:    out.FederalTotalTax,
		RegionalPayableTax: out.RegionalPayableTax,
		RegionalTotalTax:   out.RegionalTotalTax,
		TaxCredits: TaxCreditModel{
			EIPremium:              out.TaxCredits.EIPremium,
			CPPBasic:               out.TaxCredits.CanadaPensionPlanBasic,
			BPAFederal:             out.TaxCredits.FederalBasicPensionAmount,
			BPARegional:            out.TaxCredits.RegionalBasicPensionAmount,
			CanadaEmploymentAmount: out.TaxCredits.CanadaEmploymentAmount,
		},
		TaxDeductions: TaxDeductionModel{
			CPPFirstAdditional:  out.TaxDeductions.CPPFirstAdditional,
			CPPSecondAdditional: out.TaxDeductions.CPPSecondAdditional,
		},
	}
}

type TaxCreditModel struct {
	BPAFederal             float64 `json:"basic_personal_amount_federal"`
	BPARegional            float64 `json:"basic_personal_amount_regional"`
	CanadaEmploymentAmount float64 `json:"canada_employment_amount"`
	EIPremium              float64 `json:"employment_insurance_premium"`
	CPPBasic               float64 `json:"canada_pension_plan_basic"`
}

type TaxDeductionModel struct {
	CPPFirstAdditional  float64 `json:"canada_pension_plan_first_additional"`
	CPPSecondAdditional float64 `json:"canada_pension_plan_second_additional"`
}
