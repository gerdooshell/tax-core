package dataAccess

import (
	dataAccess "github.com/gerdooshell/tax-core/interactors/data_access"
)

type TaxDataService interface {
	dataAccess.FederalTaxData
	dataAccess.BCTaxData
	dataAccess.AlbertaTaxData
	dataAccess.TaxMargin
}
