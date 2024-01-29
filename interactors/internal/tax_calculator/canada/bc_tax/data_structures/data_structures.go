package dataStructures

type Input struct {
	Credits    BCTaxCredits
	Deductions BCTaxDeductions
	Salary     float64
	Year       int
}

type Output struct {
	Credits    BCTaxCredits
	Deductions BCTaxDeductions
	PayableTax float64
	TotalTax   float64
}
