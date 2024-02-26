package taxCalculator

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/gerdooshell/tax-core/controller/internal"
	"github.com/gerdooshell/tax-core/controller/internal/routes"
	restApi "github.com/gerdooshell/tax-core/controller/rest_api/handlers"
	canadaTaxInfo "github.com/gerdooshell/tax-core/interactors/controller_access/canada_tax_info"
	"github.com/gerdooshell/tax-core/library/region/canada"

	logger "github.com/gerdooshell/tax-logger-client-go"
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

func (tc *taxCalculator) ParseArgs(r *http.Request) (err error) {
	defer func() {
		if err != nil {
			logger.Error(err.Error())
		}
	}()
	state := State{}
	pathVars := mux.Vars(r)
	provinceStr, ok := pathVars["province"]
	if !ok {
		return errors.New("province is not provided")
	}
	province, err := canada.GetProvinceFromString(provinceStr)
	if err != nil {
		return err
	}
	yearStr, ok := pathVars["year"]
	if !ok {
		return errors.New("year is not provided")
	}
	year, err := strconv.Atoi(yearStr)
	if err != nil {
		return err
	}
	incomeStr, ok := pathVars["income"]
	if !ok {
		return errors.New("income is not provided")
	}
	income, err := strconv.ParseFloat(incomeStr, 64)
	if err != nil {
		return err
	}
	rrspStr, ok := pathVars["rrsp"]
	if !ok {
		return errors.New("rrsp is not provided")
	}
	rrsp, err := strconv.ParseFloat(rrspStr, 64)
	if err != nil {
		logger.Error(err.Error())
		return err
	}
	input := canadaTaxInfo.Input{
		Province:    province,
		Year:        year,
		TotalIncome: income,
		RRSP:        rrsp,
	}
	if err != nil {
		return err
	}
	state.apiKey = r.Header.Get(internal.APIKyeNameID)
	state.input = &input
	if err = validateInput(state); err != nil {
		return err
	}
	ctx := context.WithValue(r.Context(), "state", state)
	*r = *(r.WithContext(ctx))
	return nil
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
		logger.Error(err.Error())
		return resp
	}
	respBody := NewResponseBodyModelFrom(out)
	respBodyByte, err := json.Marshal(respBody)
	if err != nil {
		resp.Body = io.NopCloser(bytes.NewReader([]byte(err.Error())))
		resp.StatusCode = http.StatusInternalServerError
		logger.Error(err.Error())
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
	FederalPayableTax  float64           `json:"federalPayableTax"`
	FederalTotalTax    float64           `json:"federalTotalTax"`
	RegionalPayableTax float64           `json:"regionalPayableTax"`
	RegionalTotalTax   float64           `json:"regionalTotalTax"`
	LeftRRSPRoom       float64           `json:"leftRRSPRoom"`
	TaxCredits         TaxCreditModel    `json:"taxCredits"`
	TaxDeductions      TaxDeductionModel `json:"taxDeductions"`
}

func NewResponseBodyModelFrom(out canadaTaxInfo.Output) ResponseBodyModel {
	return ResponseBodyModel{
		FederalPayableTax:  out.FederalPayableTax,
		FederalTotalTax:    out.FederalTotalTax,
		RegionalPayableTax: out.RegionalPayableTax,
		RegionalTotalTax:   out.RegionalTotalTax,
		LeftRRSPRoom:       out.LeftRRSPRoom,
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
	BPAFederal             float64 `json:"basicPersonalAmountFederal"`
	BPARegional            float64 `json:"basicPersonalAmountRegional"`
	CanadaEmploymentAmount float64 `json:"canadaEmploymentAmount"`
	EIPremium              float64 `json:"employmentInsurancePremium"`
	CPPBasic               float64 `json:"canadaPensionPlanBasic"`
}

type TaxDeductionModel struct {
	CPPFirstAdditional  float64 `json:"canadaPensionPlanFirstAdditional"`
	CPPSecondAdditional float64 `json:"canadaPensionPlanSecondAdditional"`
}
