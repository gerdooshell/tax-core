package taxMargin

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"

	"github.com/gerdooshell/tax-core/controller/internal"
	"github.com/gerdooshell/tax-core/controller/internal/routes"
	restApi "github.com/gerdooshell/tax-core/controller/rest_api/handlers"
	"github.com/gerdooshell/tax-core/entities/canada/shared"
	canadaTaxMarginCalculator "github.com/gerdooshell/tax-core/interactors/controller_access/canada_tax_margin_calculator"
	"github.com/gerdooshell/tax-core/library/region/canada"

	"github.com/gorilla/mux"
)

type taxMargin struct {
	state *State
}

func NewTaxMarginController() restApi.Handler {
	return &taxMargin{}
}

type State struct {
	context context.Context
	input   *canadaTaxMarginCalculator.Input
	apiKey  string
}

func (tc *taxMargin) URL() string {
	return routes.TaxMargin
}

func (tc *taxMargin) Methods() []string {
	return []string{http.MethodGet}
}

func (tc *taxMargin) Authorize() error {
	return nil
}

func (tc *taxMargin) ParseArgs(r *http.Request) (*http.Request, error) {
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
	tc.state.input = &canadaTaxMarginCalculator.Input{Year: year, Province: province}
	tc.state.apiKey = r.Header.Get(internal.APIKyeNameID)
	return r, nil
}

func (tc *taxMargin) Process(_ *http.Request) *http.Response {
	resp := new(http.Response)
	ctx, cancel := context.WithTimeout(tc.state.context, time.Second*5)
	defer cancel()
	marginalTax := canadaTaxMarginCalculator.NewCanadaTaxMarginCalculator()
	out, err := marginalTax.GetAllMarginalBrackets(ctx, tc.state.input)
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

type ResponseBodyModel struct {
	MarginalBrackets []MarginalBracketModel `json:"marginal_brackets"`
}

func NewResponseBodyModelFrom(out canadaTaxMarginCalculator.Output) ResponseBodyModel {
	marginalBrackets := make([]MarginalBracketModel, 0, len(out.MarginalBrackets))
	for _, bracket := range out.MarginalBrackets {
		marginalBrackets = append(marginalBrackets, NewMarginalBracketModelFromTaxBracket(bracket))
	}
	return ResponseBodyModel{
		MarginalBrackets: marginalBrackets,
	}
}

type MarginalBracketModel struct {
	Low  float64 `json:"low"`
	High float64 `json:"high"`
	Rate float64 `json:"rate"`
}

func NewMarginalBracketModelFromTaxBracket(bracket shared.TaxBracket) MarginalBracketModel {
	return MarginalBracketModel{
		Low:  bracket.Low,
		High: bracket.High,
		Rate: bracket.Rate,
	}
}
