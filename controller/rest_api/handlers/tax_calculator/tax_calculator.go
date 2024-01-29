package taxCalculator

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
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
	body    io.ReadCloser
	input   *canadaTaxCalculator.Input
	apiKey  string
}

func (tc *taxCalculator) URL() string {
	return routes.TaxCalculator
}

func (tc *taxCalculator) Methods() []string {
	return []string{http.MethodPost}
}

func (tc *taxCalculator) Authorize() error {
	if tc.state.apiKey != internal.APIKeyValue {
		return errors.New("not authorized: invalid api key")
	}
	return nil
}

func (tc *taxCalculator) ParseArgs(r *http.Request) (*http.Request, error) {
	tc.state = &State{
		context: context.Background(),
	}
	tc.state.body = r.Body
	input, err := tc.buildCalculatorInput()
	if err != nil {
		return nil, err
	}
	tc.state.apiKey = r.Header.Get(internal.APIKyeNameID)
	tc.state.input = input
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

func (tc *taxCalculator) buildCalculatorInput() (*canadaTaxCalculator.Input, error) {
	bodyInput := new(RequestBodyModel)
	if err := json.NewDecoder(tc.state.body).Decode(bodyInput); err != nil {
		return nil, err
	}
	input, err := bodyInput.toInteractorInput()
	if err != nil {
		return nil, err
	}
	return input, nil
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

func (reqBody RequestBodyModel) toInteractorInput() (*canadaTaxCalculator.Input, error) {
	input := new(canadaTaxCalculator.Input)
	input.Salary = reqBody.Income
	input.Year = reqBody.Year
	province, err := canada.GetProvinceFromString(reqBody.Province)
	if err != nil {
		return nil, err
	}
	input.Province = province
	return input, nil
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
