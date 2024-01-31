package dataAccess

import (
	"context"
	"fmt"
	dataProvider "github.com/gerdooshell/tax-communication/src/data_provider"
	abCredits "github.com/gerdooshell/tax-core/entities/canada/alberta/credits"
	bcCredits "github.com/gerdooshell/tax-core/entities/canada/bc/credits"
	federalEntities "github.com/gerdooshell/tax-core/entities/canada/federal/credits"
	sharedEntities "github.com/gerdooshell/tax-core/entities/canada/shared"
	dataAccess "github.com/gerdooshell/tax-core/interactors/data_access"
	"google.golang.org/grpc"
	"time"
)

type DataService interface {
	dataAccess.FederalTaxData
	dataAccess.BCTaxData
	dataAccess.AlbertaTaxData
	dataAccess.AllCanadaTaxData
}

func NewDataProviderService() DataService {
	return &dataService{
		dataProviderUrl: "localhost:45432",
		timeout:         time.Second * 3,
	}
}

type dataService struct {
	dataProviderUrl string
	grpcClient      dataProvider.GRPCDataProviderClient
	timeout         time.Duration
}

func (ds *dataService) generateDataServiceClient() error {
	if ds.grpcClient != nil {
		return nil
	}
	connection, err := grpc.Dial(ds.dataProviderUrl, grpc.WithInsecure())
	if err != nil {
		return fmt.Errorf("connection failed, error: \"%v\"", err)
	}
	ds.grpcClient = dataProvider.NewGRPCDataProviderClient(connection)
	//if err = connection.Close(); err != nil {
	//	return nil, fmt.Errorf("failed closing connection, error: %v\n", err)
	//}
	return nil
}

func (ds *dataService) GetCPP(ctx context.Context, year int) (<-chan sharedEntities.CanadaPensionPlan, <-chan error) {
	out := make(chan sharedEntities.CanadaPensionPlan)
	errChan := make(chan error)
	go func() {
		defer close(out)
		defer close(errChan)
		err := ds.generateDataServiceClient()
		if err != nil {
			errChan <- err
			return
		}
		req := &dataProvider.CanadaPensionPlanRequest{
			Year: int32(year),
		}
		ctx, cancel := context.WithTimeout(ctx, ds.timeout)
		defer cancel()
		resp, err := ds.grpcClient.GetCanadaPensionPlan(ctx, req)
		if err != nil {
			errChan <- err
			return
		}
		out <- sharedEntities.CanadaPensionPlan{
			Year:                            year,
			BasicExemption:                  resp.GetBasicExemption(),
			BasicRateEmployee:               resp.GetBasicRateEmployee(),
			BasicRateEmployer:               resp.GetBasicRateEmployer(),
			FirstAdditionalRateEmployee:     resp.GetFirstAdditionalRateEmployee(),
			FirstAdditionalRateEmployer:     resp.GetFirstAdditionalRateEmployer(),
			SecondAdditionalRateEmployee:    resp.GetSecondAdditionalRateEmployee(),
			SecondAdditionalRateEmployer:    resp.GetSecondAdditionalRateEmployer(),
			MaxPensionableEarning:           resp.GetMaxPensionableEarning(),
			AdditionalMaxPensionableEarning: resp.GetAdditionalMaxPensionableEarning(),
		}
	}()
	return out, errChan
}

func (ds *dataService) GetEIPremium(ctx context.Context, year int) (<-chan sharedEntities.EmploymentInsurancePremium, <-chan error) {
	out := make(chan sharedEntities.EmploymentInsurancePremium)
	errChan := make(chan error)
	go func() {
		defer close(out)
		defer close(errChan)
		err := ds.generateDataServiceClient()
		if err != nil {
			errChan <- err
			return
		}
		req := &dataProvider.EmploymentInsurancePremiumRequest{
			Year: int32(year),
		}
		ctx, cancel := context.WithTimeout(ctx, ds.timeout)
		defer cancel()
		resp, err := ds.grpcClient.GetEmploymentInsurancePremium(ctx, req)
		if err != nil {
			errChan <- err
			return
		}
		out <- sharedEntities.EmploymentInsurancePremium{
			Rate:                              resp.GetRate(),
			MaxInsurableEarning:               resp.GetMaxInsurableEarning(),
			EmployerEmployeeContributionRatio: resp.GetEmployerEmployeeContributionRatio(),
		}
	}()
	return out, errChan
}

func (ds *dataService) GetFederalBPA(ctx context.Context, year int) (<-chan federalEntities.BasicPersonalAmount, <-chan error) {
	out := make(chan federalEntities.BasicPersonalAmount)
	errChan := make(chan error)
	go func() {
		defer close(out)
		defer close(errChan)
		err := ds.generateDataServiceClient()
		if err != nil {
			errChan <- err
			return
		}
		req := &dataProvider.FederalBPARequest{
			Year: int32(year),
		}
		ctx, cancel := context.WithTimeout(ctx, ds.timeout)
		defer cancel()
		resp, err := ds.grpcClient.GetFederalBPA(ctx, req)
		if err != nil {
			errChan <- err
			return
		}
		out <- federalEntities.BasicPersonalAmount{
			MaxBPAAmount: resp.GetMaxBPAAmount(),
			MaxBPAIncome: resp.GetMaxBPAIncome(),
			MinBPAAmount: resp.GetMinBPAAmount(),
			MinBPAIncome: resp.GetMinBPAIncome(),
		}
	}()
	return out, errChan
}

func (ds *dataService) GetFederalTaxBrackets(ctx context.Context, year int) (<-chan []sharedEntities.TaxBracket, <-chan error) {
	out := make(chan []sharedEntities.TaxBracket)
	errChan := make(chan error)
	go func() {
		defer close(out)
		defer close(errChan)
		err := ds.generateDataServiceClient()
		if err != nil {
			errChan <- err
			return
		}
		req := &dataProvider.FederalTaxBracketsRequest{
			Year: int32(year),
		}
		ctx, cancel := context.WithTimeout(ctx, ds.timeout)
		defer cancel()
		resp, err := ds.grpcClient.GetFederalTaxBrackets(ctx, req)
		if err != nil {
			errChan <- err
			return
		}
		rates := resp.GetRates()
		levels := resp.GetLevels()
		brackets, err := sharedEntities.FromArray(levels, rates)
		if err != nil {
			errChan <- err
			return
		}
		out <- brackets
	}()
	return out, errChan
}

func (ds *dataService) GetCEA(ctx context.Context, year int) (<-chan federalEntities.CanadaEmploymentAmount, <-chan error) {
	out := make(chan federalEntities.CanadaEmploymentAmount)
	errChan := make(chan error)
	go func() {
		defer close(out)
		defer close(errChan)
		err := ds.generateDataServiceClient()
		if err != nil {
			errChan <- err
			return
		}
		req := &dataProvider.CanadaEmploymentAmountRequest{
			Year: int32(year),
		}
		ctx, cancel := context.WithTimeout(ctx, ds.timeout)
		defer cancel()
		resp, err := ds.grpcClient.GetCEA(ctx, req)
		if err != nil {
			errChan <- err
			return
		}
		out <- federalEntities.CanadaEmploymentAmount{Value: resp.GetCeaValue()}
	}()
	return out, errChan
}

func (ds *dataService) GetBCBPA(ctx context.Context, year int) (<-chan bcCredits.BasicPersonalAmount, <-chan error) {
	out := make(chan bcCredits.BasicPersonalAmount)
	errChan := make(chan error)
	go func() {
		defer close(out)
		defer close(errChan)
		err := ds.generateDataServiceClient()
		if err != nil {
			errChan <- err
			return
		}
		req := &dataProvider.BritishColumbiaBPARequest{
			Year: int32(year),
		}
		ctx, cancel := context.WithTimeout(ctx, ds.timeout)
		defer cancel()
		resp, err := ds.grpcClient.GetBritishColumbiaBPA(ctx, req)
		if err != nil {
			errChan <- err
			return
		}
		out <- bcCredits.BasicPersonalAmount{Value: resp.GetBpaValue()}
	}()
	return out, errChan
}

func (ds *dataService) GetBCTaxBrackets(ctx context.Context, year int) (<-chan []sharedEntities.TaxBracket, <-chan error) {
	out := make(chan []sharedEntities.TaxBracket)
	errChan := make(chan error)
	go func() {
		defer close(out)
		defer close(errChan)
		err := ds.generateDataServiceClient()
		if err != nil {
			errChan <- err
			return
		}
		req := &dataProvider.BritishColumbiaTaxBracketsRequest{
			Year: int32(year),
		}
		ctx, cancel := context.WithTimeout(ctx, ds.timeout)
		defer cancel()
		resp, err := ds.grpcClient.GetBritishColumbiaTaxBrackets(ctx, req)
		if err != nil {
			errChan <- err
			return
		}
		rates := resp.GetRates()
		levels := resp.GetLevels()
		brackets, err := sharedEntities.FromArray(levels, rates)
		if err != nil {
			errChan <- err
			return
		}
		out <- brackets
	}()
	return out, errChan
}

func (ds *dataService) GetAlbertaBPA(ctx context.Context, year int) (<-chan abCredits.BasicPersonalAmount, <-chan error) {
	out := make(chan abCredits.BasicPersonalAmount)
	errChan := make(chan error)
	go func() {
		defer close(out)
		defer close(errChan)
		err := ds.generateDataServiceClient()
		if err != nil {
			errChan <- err
			return
		}
		req := &dataProvider.AlbertaBPARequest{
			Year: int32(year),
		}
		ctx, cancel := context.WithTimeout(ctx, ds.timeout)
		defer cancel()
		resp, err := ds.grpcClient.GetAlbertaBPA(ctx, req)
		if err != nil {
			errChan <- err
			return
		}
		out <- abCredits.BasicPersonalAmount{Value: resp.GetBpaValue()}
	}()
	return out, errChan
}

func (ds *dataService) GetAlbertaTaxBrackets(ctx context.Context, year int) (<-chan []sharedEntities.TaxBracket, <-chan error) {
	out := make(chan []sharedEntities.TaxBracket)
	errChan := make(chan error)
	go func() {
		defer close(out)
		defer close(errChan)
		err := ds.generateDataServiceClient()
		if err != nil {
			errChan <- err
			return
		}
		req := &dataProvider.AlbertaTaxBracketsRequest{
			Year: int32(year),
		}
		ctx, cancel := context.WithTimeout(ctx, ds.timeout)
		defer cancel()
		resp, err := ds.grpcClient.GetAlbertaTaxBrackets(ctx, req)
		if err != nil {
			errChan <- err
			return
		}
		rates := resp.GetRates()
		levels := resp.GetLevels()
		brackets, err := sharedEntities.FromArray(levels, rates)
		if err != nil {
			errChan <- err
			return
		}
		out <- brackets
	}()
	return out, errChan
}
