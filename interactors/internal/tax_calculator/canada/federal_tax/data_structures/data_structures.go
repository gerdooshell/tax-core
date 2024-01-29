package dataStructures

type Input struct {
	Credits    FederalTaxCredit
	Deductions FederalTaxDeductions
	Salary     float64
	Year       int
}

type Output struct {
	Credits    FederalTaxCredit
	Deductions FederalTaxDeductions
	PayableTax float64
	TotalTax   float64
}
