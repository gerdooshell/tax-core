package mockService

import (
	"context"
	"github.com/gerdooshell/tax-core/data-access"
	bcCredits "github.com/gerdooshell/tax-core/entities/canada/bc/credits"
	fedCredits "github.com/gerdooshell/tax-core/entities/canada/federal/credits"
	"github.com/gerdooshell/tax-core/entities/canada/shared"
	dataAccessInteractor "github.com/gerdooshell/tax-core/interactors/data_access"
	"github.com/gerdooshell/tax-core/library/region/canada"
)

func NewPostgresServiceMock() dataAccess.DataProviderService {
	return &postgresServiceMock{}
}

type postgresServiceMock struct {
}

func (pg *postgresServiceMock) GetRRSP(ctx context.Context, year int) <-chan dataAccessInteractor.RRSPDataOut {
	//TODO implement me
	panic("implement me")
}

func (pg *postgresServiceMock) GetTaxBrackets(ctx context.Context, year int, province canada.Province) <-chan dataAccessInteractor.TaxBracketsDataOut {
	//TODO implement me
	panic("implement me")
}

func (pg *postgresServiceMock) GetCombinedMarginalBrackets(ctx context.Context, year int, province canada.Province) <-chan dataAccessInteractor.TaxBracketsDataOut {
	//TODO implement me
	panic("implement me")
}

func (pg *postgresServiceMock) SaveMarginalTaxBrackets(ctx context.Context, province canada.Province, year int, brackets []shared.TaxBracket) <-chan error {
	//TODO implement me
	panic("implement me")
}

func (pg *postgresServiceMock) GetCPP(ctx context.Context, year int) <-chan dataAccessInteractor.CPPDataOut {
	data := make(chan dataAccessInteractor.CPPDataOut, 1)
	go func() {
		defer close(data)
		data <- dataAccessInteractor.CPPDataOut{
			CanadaPensionPlan: shared.CanadaPensionPlan{
				Year:                            year,
				BasicExemption:                  3500,
				BasicRateEmployee:               4.950,
				BasicRateEmployer:               4.950,
				FirstAdditionalRateEmployee:     1,
				FirstAdditionalRateEmployer:     1,
				SecondAdditionalRateEmployee:    0,
				SecondAdditionalRateEmployer:    0,
				MaxPensionableEarning:           66600,
				AdditionalMaxPensionableEarning: 0,
			},
		}
	}()
	return data
}

func (pg *postgresServiceMock) GetFederalBPA(ctx context.Context, year int) <-chan dataAccessInteractor.FederalBPADataOut {
	data := make(chan dataAccessInteractor.FederalBPADataOut, 1)
	go func() {
		defer close(data)
		data <- dataAccessInteractor.FederalBPADataOut{
			BasicPersonalAmount: fedCredits.BasicPersonalAmount{
				MaxBPAIncome: 165430,
				MinBPAIncome: 235675,
				MaxBPAAmount: 15000,
				MinBPAAmount: 13521,
			},
		}
	}()
	return data
}

func (pg *postgresServiceMock) GetEIPremium(_ context.Context, _ int) <-chan dataAccessInteractor.EIPremiumDataOut {
	data := make(chan dataAccessInteractor.EIPremiumDataOut, 1)
	go func() {
		defer close(data)
		data <- dataAccessInteractor.EIPremiumDataOut{
			EmploymentInsurancePremium: shared.EmploymentInsurancePremium{
				MaxInsurableEarning:               61500,
				Rate:                              1.63,
				EmployerEmployeeContributionRatio: 1.4},
		}
	}()
	return data
}

func (pg *postgresServiceMock) GetBritishColumbiaBPA(_ context.Context, _ int) <-chan dataAccessInteractor.BritishColumbiaBPADataOut {
	data := make(chan dataAccessInteractor.BritishColumbiaBPADataOut, 1)
	go func() {
		defer close(data)
		data <- dataAccessInteractor.BritishColumbiaBPADataOut{
			BasicPersonalAmount: bcCredits.BasicPersonalAmount{Value: 11981},
		}
	}()
	return data
}

func (pg *postgresServiceMock) GetCEA(ctx context.Context, year int) <-chan dataAccessInteractor.CEADataOut {
	data := make(chan dataAccessInteractor.CEADataOut, 1)
	go func() {
		defer close(data)
		data <- dataAccessInteractor.CEADataOut{
			CanadaEmploymentAmount: fedCredits.CanadaEmploymentAmount{
				Value: 1368,
			},
		}
	}()
	return data
}

func (pg *postgresServiceMock) GetAlbertaBPA(_ context.Context, _ int) <-chan dataAccessInteractor.AlbertaBPADataOut {
	return nil
}
