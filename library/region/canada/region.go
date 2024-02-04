package canada

import "fmt"

type Province string

const (
	Federal         Province = "Federal"
	Alberta         Province = "Alberta"
	BritishColumbia Province = "British Columbia"
	Manitoba        Province = "Manitoba"
	Ontario         Province = "Ontario"
	NovaScotia      Province = "Nova Scotia"
	NewBrunswick    Province = "New Brunswick"
	Quebec          Province = "Quebec"
	Saskatchewan    Province = "Saskatchewan"
	Yukon           Province = "Yukon"
)

func GetProvinceFromString(province string) (pr Province, err error) {
	switch province {
	case "federal", "fed":
		pr = Federal
	case "alberta", "ab":
		pr = Alberta
	case "british_columbia", "bc":
		pr = BritishColumbia
	case "ontario", "on":
		pr = Ontario
	case "manitoba, mb":
		pr = Manitoba
	case "quebec", "qc", "pq":
		pr = Quebec
	case "saskatchewan", "sk":
		pr = Saskatchewan
	case "nova_scotia", "ns":
		pr = NovaScotia
	case "new_brunswick", "nb":
		pr = NewBrunswick
	case "yukon", "yk":
		pr = Yukon
	default:
		pr = Federal
		err = fmt.Errorf("invalid province: \"%v\"", province)
	}
	return
}
