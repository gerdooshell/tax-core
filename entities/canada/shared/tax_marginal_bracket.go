package shared

import (
	"errors"
	"math"
	"sort"
)

type TaxMarginalBracket struct {
	FederalBrackets     []TaxBracket
	RegionalTaxBrackets []TaxBracket
	margins             []TaxBracket
}

func (tm *TaxMarginalBracket) getMarginalRateForTaxableIncome(income float64, isRegional bool) float64 {
	brackets := tm.FederalBrackets
	if isRegional {
		brackets = tm.RegionalTaxBrackets
	}
	var rate float64
	for _, bracket := range brackets {
		if income >= bracket.Low && income < bracket.High {
			rate = bracket.Rate
			break
		}
	}
	return rate
}

func (tm *TaxMarginalBracket) getLowValues(isRegional bool) []float64 {
	brackets := tm.FederalBrackets
	if isRegional {
		brackets = tm.RegionalTaxBrackets
	}
	lows := make([]float64, 0)
	for _, bracket := range brackets {
		lows = append(lows, bracket.Low)
	}
	return lows
}

func (tm *TaxMarginalBracket) CalcCombinedTaxMargins() error {
	if err := tm.validateInputParameters(); err != nil {
		return err
	}
	sort.Slice(tm.FederalBrackets, func(i, j int) bool {
		return tm.FederalBrackets[i].Low < tm.FederalBrackets[j].Low
	})
	sort.Slice(tm.RegionalTaxBrackets, func(i, j int) bool {
		return tm.RegionalTaxBrackets[i].Low < tm.RegionalTaxBrackets[j].Low
	})
	federalLows := tm.getLowValues(false)
	regionalLows := tm.getLowValues(true)
	lows := append(federalLows, regionalLows...)
	sort.Slice(lows, func(i, j int) bool { return lows[i] < lows[j] })
	mergedBrackets := make([]TaxBracket, 0, len(lows))
	for i, low := range lows {
		if (i > 1 && lows[i-1] == low) || (i < len(lows)-1) && low == lows[i+1] {
			continue
		}
		fedRate := tm.getMarginalRateForTaxableIncome(low, false)
		regRate := tm.getMarginalRateForTaxableIncome(low, true)
		mergedBrackets = append(mergedBrackets, TaxBracket{High: 0, Low: low, Rate: regRate + fedRate})
	}
	var high float64
	for i := range mergedBrackets {
		if i < len(mergedBrackets)-1 {
			high = mergedBrackets[i+1].Low
		} else {
			high = math.MaxFloat64
		}
		mergedBrackets[i].High = high
	}
	tm.margins = mergedBrackets
	return nil
}

func (tm *TaxMarginalBracket) GetMargins() []TaxBracket {
	return tm.margins
}

func (tm *TaxMarginalBracket) validateInputParameters() error {
	if len(tm.FederalBrackets) == 0 {
		return errors.New("marginal tax: invalid federal brackets")
	}
	if len(tm.RegionalTaxBrackets) == 0 {
		return errors.New("marginal tax: invalid regional brackets")
	}
	return nil
}
