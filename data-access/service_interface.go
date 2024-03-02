package dataAccess

import (
	dataAccessInteractor "github.com/gerdooshell/tax-core/interactors/data_access"
)

type DataProviderService interface {
	dataAccessInteractor.TaxMargin
	dataAccessInteractor.FederalBPAData
	dataAccessInteractor.EIPremiumData
	dataAccessInteractor.CanadaEmploymentAmountData
	dataAccessInteractor.CanadaPensionPlanData
	dataAccessInteractor.RegionalBPAData
	dataAccessInteractor.TaxBracketData
	dataAccessInteractor.RRSPData
}
