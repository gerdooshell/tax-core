package optimalRrsp

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	logger "github.com/gerdooshell/tax-logger-client-go"
	"io"
	"net/http"
	"strconv"

	"github.com/gerdooshell/tax-core/controller/internal"
	"github.com/gerdooshell/tax-core/controller/internal/routes"
	restApi "github.com/gerdooshell/tax-core/controller/rest_api/handlers"
	"github.com/gerdooshell/tax-core/interactors/controller_access/registered_retirement_saving_plan"
	"github.com/gerdooshell/tax-core/library/region/canada"

	"github.com/gorilla/mux"
)

func NewOptimalRRSPController() restApi.Handler {
	return &optimalRRSPController{}
}

type state struct {
	input  *rrspInfo.OptimalInput
	apiKey string
}

type optimalRRSPController struct {
}

func (o optimalRRSPController) URL() string {
	return routes.OptimalRRSP
}

func (o optimalRRSPController) Methods() []string {
	return []string{http.MethodGet}
}

func (o optimalRRSPController) Authorize() error {
	return nil
}

func (o optimalRRSPController) ParseArgs(r *http.Request) (err error) {
	defer func() {
		if err != nil {
			logger.Error(err.Error())
		}
	}()
	reqState := state{}
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
		return err
	}

	input := rrspInfo.OptimalInput{
		Province:        province,
		Year:            year,
		TotalIncome:     income,
		ContributedRRSP: rrsp,
	}
	if err != nil {
		return err
	}
	reqState.apiKey = r.Header.Get(internal.APIKyeNameID)
	reqState.input = &input
	if err = validateInput(reqState); err != nil {
		return err
	}
	ctx := context.WithValue(r.Context(), "state", reqState)
	*r = *(r.WithContext(ctx))
	return nil
}

func (o optimalRRSPController) Process(r *http.Request) *http.Response {
	var err error
	defer func() {
		if err != nil {
			logger.Error(err.Error())
		}
	}()
	resp := new(http.Response)
	reqState := r.Context().Value("state").(state)
	ctx := r.Context()
	calculator := rrspInfo.NewRegisteredRetirementSavingPlan()
	out, err := calculator.GetOptimalRRSPContributions(ctx, reqState.input)
	if err != nil {
		resp.Body = io.NopCloser(bytes.NewReader([]byte(err.Error())))
		resp.StatusCode = http.StatusInternalServerError
		return resp
	}
	respBody := newResponseBodyFromModel(out)
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

type RRSPOptimalResp struct {
	OptimalRRSPSuggestions []RRSPOptimalModel `json:"optimalRRSPSuggestions"`
}

type RRSPOptimalModel struct {
	PayableTax    float64 `json:"payableTax"`
	TaxableIncome float64 `json:"taxableIncome"`
	TotalRRSP     float64 `json:"totalRRSP"`
	TaxReturn     float64 `json:"taxReturn"`
	LeftRRSPRoom  float64 `json:"leftRRSPRoom"`
}

func validateInput(reqState state) error {
	if reqState.input.Year <= 0 {
		return fmt.Errorf("invalid year \"%v\"", reqState.input.Year)
	}
	if reqState.input.TotalIncome <= 0 {
		return fmt.Errorf("invalid income \"%v\"", reqState.input.TotalIncome)
	}
	if reqState.input.Province == canada.Federal {
		return fmt.Errorf("invalid province \"%v\"", reqState.input.Province)
	}
	if reqState.input.ContributedRRSP < 0 {
		return fmt.Errorf("invalid rrsp contribution \"%v\"", reqState.input.Province)
	}
	return nil
}

func newResponseBodyFromModel(out []rrspInfo.OptimalOutput) RRSPOptimalResp {
	models := make([]RRSPOptimalModel, 0, len(out))
	for _, o := range out {
		models = append(models, RRSPOptimalModel{
			PayableTax:    o.PayableTax,
			TaxableIncome: o.TaxableIncome,
			TotalRRSP:     o.RRSP,
			TaxReturn:     o.TaxReturn,
			LeftRRSPRoom:  o.LeftRRSPRoom,
		})
	}
	return RRSPOptimalResp{OptimalRRSPSuggestions: models}
}
