package routes

const (
	BaseURL       = ":8185"
	ServerPath    = "/api"
	TaxCalculator = ServerPath + "/tax"
	TaxMargin     = ServerPath + "/margin/province/{province}/year/{year}"
)
