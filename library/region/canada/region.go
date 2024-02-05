package canada

import "fmt"

type Province string

const (
	Federal         Province = "federal"
	Alberta         Province = "alberta"
	BritishColumbia Province = "british-columbia"
	Manitoba        Province = "manitoba"
	Ontario         Province = "ontario"
	NovaScotia      Province = "nova-scotia"
	NewBrunswick    Province = "new-brunswick"
	Quebec          Province = "quebec"
	Saskatchewan    Province = "saskatchewan"
	Yukon           Province = "yukon"
)

func GetProvinceFromString(province string) (pr Province, err error) {
	switch province {
	case "federal", "fed":
		pr = Federal
	case "alberta", "ab":
		pr = Alberta
	case "british-columbia", "bc":
		pr = BritishColumbia
	case "ontario", "on":
		pr = Ontario
	case "manitoba, mb":
		pr = Manitoba
	case "quebec", "qc", "pq":
		pr = Quebec
	case "saskatchewan", "sk":
		pr = Saskatchewan
	case "nova-scotia", "ns":
		pr = NovaScotia
	case "new-brunswick", "nb":
		pr = NewBrunswick
	case "yukon", "yt":
		pr = Yukon
	default:
		pr = Federal
		err = fmt.Errorf("invalid province: \"%v\"", province)
	}
	return
}
