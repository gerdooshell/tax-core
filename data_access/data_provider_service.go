package dataAccess

import (
	"context"
	"fmt"
	dataProvider "github.com/gerdooshell/tax-communication/src/data-provider"
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

func NewDataService(dataProviderUrl string) DataService {
	return &dataService{
		dataProviderUrl: dataProviderUrl,
		timeout:         time.Second * 3,
	}
}

type dataService struct {
	dataProviderUrl string
	timeout         time.Duration
}

func (ds *dataService) genDataServiceClient() (dataProvider.GRPCDataProviderClient, error) {
	connection, err := grpc.Dial(ds.dataProviderUrl, grpc.WithInsecure())
	if err != nil {
		return nil, fmt.Errorf("connection failed, error: \"%v\"", err)
	}
	client := dataProvider.NewGRPCDataProviderClient(connection)
	if err = connection.Close(); err != nil {
		return nil, fmt.Errorf("failed closing connection, error: %v\n", err)
	}
	return client, nil
}

func (ds *dataService) GetCPP(ctx context.Context, year int) (<-chan sharedEntities.CanadaPensionPlan, <-chan error) {
	//TODO implement me
	panic("implement me")
}

func (ds *dataService) GetEIPremium(ctx context.Context, year int) (<-chan sharedEntities.EmploymentInsurancePremium, <-chan error) {
	//TODO implement me
	panic("implement me")
}

func (ds *dataService) GetFederalBPA(ctx context.Context, year int) (<-chan federalEntities.BasicPersonalAmount, <-chan error) {
	out := make(chan federalEntities.BasicPersonalAmount)
	errChan := make(chan error)
	go func() {
		defer close(out)
		defer close(errChan)
		client, err := ds.genDataServiceClient()
		if err != nil {
			errChan <- err
			return
		}
		req := &dataProvider.FederalBPARequest{
			Year: int32(year),
		}
		ctx, cancel := context.WithTimeout(ctx, ds.timeout)
		defer cancel()
		resp, err := client.GetFederalBPA(ctx, req)
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
		client, err := ds.genDataServiceClient()
		if err != nil {
			errChan <- err
			return
		}
		req := &dataProvider.FederalTaxBracketsRequest{
			Year: int32(year),
		}
		ctx, cancel := context.WithTimeout(ctx, ds.timeout)
		defer cancel()
		resp, err := client.GetFederalTaxBrackets(ctx, req)
		if err != nil {
			errChan <- err
			return
		}
		rates := resp.GetRates()
		levels := resp.GetBrackets()
		brackets, err := sharedEntities.FromArray(rates, levels)
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
		client, err := ds.genDataServiceClient()
		if err != nil {
			errChan <- err
			return
		}
		req := &dataProvider.CanadaEmploymentAmountRequest{
			Year: int32(year),
		}
		ctx, cancel := context.WithTimeout(ctx, ds.timeout)
		defer cancel()
		resp, err := client.GetCEA(ctx, req)
		if err != nil {
			errChan <- err
			return
		}
		out <- federalEntities.CanadaEmploymentAmount{Value: resp.GetCeaValue()}
	}()
	return out, errChan
}

func (ds *dataService) GetBCBPA(ctx context.Context, year int) (<-chan bcCredits.BasicPersonalAmount, <-chan error) {
	//TODO implement me
	panic("implement me")
}

func (ds *dataService) GetBCTaxBrackets(ctx context.Context, year int) (<-chan []sharedEntities.TaxBracket, <-chan error) {
	//TODO implement me
	panic("implement me")
}

func (ds *dataService) GetAlbertaBPA(ctx context.Context, year int) (<-chan abCredits.BasicPersonalAmount, <-chan error) {
	//TODO implement me
	panic("implement me")
}

func (ds *dataService) GetAlbertaTaxBrackets(ctx context.Context, year int) (<-chan []sharedEntities.TaxBracket, <-chan error) {
	//TODO implement me
	panic("implement me")
}
