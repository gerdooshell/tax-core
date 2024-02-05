package taxCalculator

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"io"
	"net/http"
	"strconv"
	"time"

	"github.com/gerdooshell/tax-core/controller/internal"
	"github.com/gerdooshell/tax-core/controller/internal/routes"
	restApi "github.com/gerdooshell/tax-core/controller/rest_api/handlers"
	canadaTaxCalculator "github.com/gerdooshell/tax-core/interactors/controller_access/canada_tax_calculator"
	canadaTaxImplementation "github.com/gerdooshell/tax-core/interactors/controller_access/canada_tax_calculator/implementation"
	"github.com/gerdooshell/tax-core/library/region/canada"
)

type taxCalculator struct {
	state *State
}

func NewTaxCalculatorController() restApi.Handler {
	return &taxCalculator{}
}

type State struct {
	context context.Context
	input   *canadaTaxCalculator.Input
	apiKey  string
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
	tc.state = &State{
		context: context.Background(),
	}
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
	input := canadaTaxCalculator.Input{
		Province: province,
		Year:     year,
		Salary:   income,
	}
	if err != nil {
		return nil, err
	}
	tc.state.apiKey = r.Header.Get(internal.APIKyeNameID)
	tc.state.input = &input
	if err := tc.validateInput(); err != nil {
		return nil, err
	}
	return r, nil
}

func (tc *taxCalculator) Process(_ *http.Request) *http.Response {
	resp := new(http.Response)
	ctx, cancel := context.WithTimeout(tc.state.context, time.Second*5)
	defer cancel()
	calculator := canadaTaxImplementation.NewCanadaTaxCalculator()
	out, err := calculator.Calculate(ctx, tc.state.input)
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

func (tc *taxCalculator) validateInput() error {
	if tc.state.input.Year <= 0 {
		return fmt.Errorf("invalid year \"%v\"", tc.state.input.Year)
	}
	if tc.state.input.Salary <= 0 {
		return fmt.Errorf("invalid income \"%v\"", tc.state.input.Salary)
	}
	if tc.state.input.Province == canada.Federal {
		return fmt.Errorf("invalid province \"%v\"", tc.state.input.Province)
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

func NewResponseBodyModelFrom(out canadaTaxCalculator.Output) ResponseBodyModel {
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
