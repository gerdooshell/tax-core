package dataAccess

import (
	federalCredits "github.com/gerdooshell/tax-core/entities/canada/federal/credits"
	sharedEntities "github.com/gerdooshell/tax-core/entities/canada/shared"
)

type FederalBPADataOut struct {
	BasicPersonalAmount federalCredits.BasicPersonalAmount
	Err                 error
}

type TaxBracketsDataOut struct {
	TaxBrackets []sharedEntities.TaxBracket
	Err         error
}

type EIPremiumDataOut struct {
	EmploymentInsurancePremium sharedEntities.EmploymentInsurancePremium
	Err                        error
}
