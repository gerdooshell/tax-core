package dataAccess

import (
	dataAccess "github.com/gerdooshell/tax-core/interactors/data_access"
)

type DataProviderService interface {
	dataAccess.FederalTaxData
	dataAccess.BCTaxData
	dataAccess.AlbertaTaxData
	dataAccess.TaxMargin
}
