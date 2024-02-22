package routes

const (
	BaseURL       = ":8185"
	ServerPath    = "/api"
	TaxCalculator = ServerPath + "/tax/year/{year}/province/{province}/income/{income}"
	TaxMargin     = ServerPath + "/margin/year/{year}/province/{province}"
	OptimalRRSP   = ServerPath + "/optimal-rrsp/year/{year}/province/{province}/income/{income}/contributed-rrsp/{contributed-rrsp}"
)
