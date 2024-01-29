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
	case "federal":
		pr = Federal
	case "alberta":
		pr = Alberta
	case "british_columbia":
		pr = BritishColumbia
	case "ontario":
		pr = Ontario
	case "manitoba":
		pr = Manitoba
	case "quebec":
		pr = Quebec
	case "saskatchewan":
		pr = Saskatchewan
	case "nova_scotia":
		pr = NovaScotia
	case "new_brunswick":
		pr = NewBrunswick
	case "yukon":
		pr = Yukon
	default:
		pr = Federal
		err = fmt.Errorf("invalid province: \"%v\"", province)
	}
	return
}
