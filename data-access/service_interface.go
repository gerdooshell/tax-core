package dataAccess

import (
	dataAccess "github.com/gerdooshell/tax-core/interactors/data_access"
)

type DataProviderService interface {
	dataAccess.FederalTaxData
	dataAccess.BCTaxData
	dataAccess.AlbertaTaxData
	dataAccess.TaxMargin
	dataAccess.FederalBPAData
	dataAccess.EIPremiumData
	dataAccess.CanadaEmploymentAmountData
	dataAccess.CanadaPensionPlanData
	dataAccess.RegionalBPAData
	dataAccess.TaxBracketData
}
