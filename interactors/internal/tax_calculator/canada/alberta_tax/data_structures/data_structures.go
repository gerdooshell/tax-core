package dataStructures

type Input struct {
	Credits    AlbertaTaxCredit
	Deductions AlbertaTaxDeductions
	Salary     float64
	Year       int
}

type Output struct {
	Credits    AlbertaTaxCredit
	Deductions AlbertaTaxDeductions
	PayableTax float64
	TotalTax   float64
}
