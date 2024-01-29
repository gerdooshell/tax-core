package dataStructures

import (
	"github.com/gerdooshell/tax-core/library/region/canada"
)

type Input struct {
	Province   canada.Province
	Credits    RegionalTaxCredit
	Deductions RegionalTaxDeductions
	Salary     float64
	Year       int
}

type Output struct {
	Credits    RegionalTaxCredit
	Deductions RegionalTaxDeductions
	PayableTax float64
	TotalTax   float64
}
