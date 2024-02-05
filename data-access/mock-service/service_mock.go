package mockService

import (
	"context"
	"github.com/gerdooshell/tax-core/data-access"
	abCredits "github.com/gerdooshell/tax-core/entities/canada/alberta/credits"
	bcCredits "github.com/gerdooshell/tax-core/entities/canada/bc/credits"
	fedCredits "github.com/gerdooshell/tax-core/entities/canada/federal/credits"
	"github.com/gerdooshell/tax-core/entities/canada/shared"
	"github.com/gerdooshell/tax-core/library/region/canada"
	"math"
)

func NewPostgresServiceMock() dataAccess.DataProviderService {
	return &postgresServiceMock{}
}

type postgresServiceMock struct {
}

func (pg *postgresServiceMock) GetTaxBrackets(ctx context.Context, year int, province canada.Province) (<-chan []shared.TaxBracket, <-chan error) {
	//TODO implement me
	panic("implement me")
}

func (pg *postgresServiceMock) GetCombinedMarginalBrackets(ctx context.Context, year int, province canada.Province) (<-chan []shared.TaxBracket, <-chan error) {
	//TODO implement me
	panic("implement me")
}

func (pg *postgresServiceMock) SaveMarginalTaxBrackets(ctx context.Context, province canada.Province, year int, brackets []shared.TaxBracket) (<-chan bool, <-chan error) {
	//TODO implement me
	panic("implement me")
}

func (pg *postgresServiceMock) GetCPP(ctx context.Context, year int) (<-chan shared.CanadaPensionPlan, <-chan error) {
	data := make(chan shared.CanadaPensionPlan)
	go func() {
		data <- shared.CanadaPensionPlan{
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
		}
	}()
	return data, nil
}

func (pg *postgresServiceMock) GetFederalBPA(_ context.Context, _ int) (<-chan fedCredits.BasicPersonalAmount, <-chan error) {
	data := make(chan fedCredits.BasicPersonalAmount)
	go func() {
		data <- fedCredits.BasicPersonalAmount{
			MaxBPAIncome: 165430,
			MinBPAIncome: 235675,
			MaxBPAAmount: 15000,
			MinBPAAmount: 13521,
		}
	}()
	return data, nil
}

func (pg *postgresServiceMock) GetEIPremium(_ context.Context, _ int) (<-chan shared.EmploymentInsurancePremium, <-chan error) {
	data := make(chan shared.EmploymentInsurancePremium)
	go func() {
		data <- shared.EmploymentInsurancePremium{
			MaxInsurableEarning:               61500,
			Rate:                              1.63,
			EmployerEmployeeContributionRatio: 1.4}
	}()
	return data, nil
}

func (pg *postgresServiceMock) GetFederalTaxBrackets(_ context.Context, _ int) (<-chan []shared.TaxBracket, <-chan error) {
	data := make(chan []shared.TaxBracket)
	go func() {
		data <- []shared.TaxBracket{
			{High: 53359, Low: 0, Rate: 15},
			{High: 106717, Low: 53359, Rate: 20.5},
			{High: 165430, Low: 106717, Rate: 26},
			{High: 235675, Low: 165430, Rate: 29},
			{High: math.MaxFloat64, Low: 235675, Rate: 33},
		}
	}()
	return data, nil
}

func (pg *postgresServiceMock) GetBCBPA(_ context.Context, _ int) (<-chan bcCredits.BasicPersonalAmount, <-chan error) {
	data := make(chan bcCredits.BasicPersonalAmount)
	go func() { data <- bcCredits.BasicPersonalAmount{Value: 11981} }()
	return data, nil
}

func (pg *postgresServiceMock) GetBCTaxBrackets(_ context.Context, _ int) (<-chan []shared.TaxBracket, <-chan error) {
	data := make(chan []shared.TaxBracket)
	go func() {
		data <- []shared.TaxBracket{
			{High: 45654, Low: 0, Rate: 5.06},
			{High: 91310, Low: 45654, Rate: 7.7},
			{High: 104835, Low: 91310, Rate: 10.5},
			{High: 127299, Low: 104835, Rate: 12.29},
			{High: 172602, Low: 127299, Rate: 14.7},
			{High: 240716, Low: 172602, Rate: 16.8},
			{High: math.MaxFloat64, Low: 240716, Rate: 20.5},
		}
	}()
	return data, nil
}

func (pg *postgresServiceMock) GetCEA(_ context.Context, _ int) (<-chan fedCredits.CanadaEmploymentAmount, <-chan error) {
	data := make(chan fedCredits.CanadaEmploymentAmount)
	go func() {
		data <- fedCredits.CanadaEmploymentAmount{
			Value: 1368,
		}
	}()
	return data, nil
}

func (pg *postgresServiceMock) GetAlbertaBPA(_ context.Context, _ int) (<-chan abCredits.BasicPersonalAmount, <-chan error) {
	return nil, nil
}

func (pg *postgresServiceMock) GetAlbertaTaxBrackets(_ context.Context, _ int) (<-chan []shared.TaxBracket, <-chan error) {
	return nil, nil
}
