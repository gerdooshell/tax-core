package dataAccess

import (
	"context"
	"fmt"
	dataProvider "github.com/gerdooshell/tax-communication/src/data-provider"
	abCredits "github.com/gerdooshell/tax-core/entities/canada/alberta/credits"
	bcCredits "github.com/gerdooshell/tax-core/entities/canada/bc/credits"
	"github.com/gerdooshell/tax-core/entities/canada/federal/credits"
	"github.com/gerdooshell/tax-core/entities/canada/shared"
	dataAccess "github.com/gerdooshell/tax-core/interactors/data_access"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"log"
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

func (ds *dataService) GetCPP(ctx context.Context, year int) (<-chan shared.CanadaPensionPlan, <-chan error) {
	//TODO implement me
	panic("implement me")
}

func (ds *dataService) GetEIPremium(ctx context.Context, year int) (<-chan shared.EmploymentInsurancePremium, <-chan error) {
	//TODO implement me
	panic("implement me")
}

func (ds *dataService) GetFederalBPA(ctx context.Context, year int) (<-chan credits.BasicPersonalAmount, <-chan error) {
	out := make(chan credits.BasicPersonalAmount)
	errChan := make(chan error)
	go func() {
		defer close(out)
		defer close(errChan)
		connection, err := grpc.Dial(ds.dataProviderUrl, grpc.WithInsecure())
		if err != nil {
			fmt.Println("connection failed, error:", err)
		}
		defer func() {
			if err = connection.Close(); err != nil {
				log.Printf("failed closing connection, error: %v\n", err)
			}
		}()
		ctx, cancel := context.WithTimeout(ctx, ds.timeout)
		defer cancel()
		client := dataProvider.NewGRPCDataProviderClient(connection)
		req := &dataProvider.FederalBPARequest{
			Year: int32(year),
		}
		resp, err := client.GetFederalBPA(ctx, req)
		if err != nil {
			statusError, ok := status.FromError(err)
			if ok {
				switch statusError.Code() {
				case codes.InvalidArgument:
					// handle the logic
				case codes.DeadlineExceeded:
					// handle the logic
				default:
					//handle the logic
				}
			} else {
				fmt.Printf("failed getting user info from server, error: %v\n", err)
			}
			errChan <- err
		}
		out <- credits.BasicPersonalAmount{
			MaxBPAAmount: resp.GetMaxBPAAmount(),
			MaxBPAIncome: resp.GetMaxBPAIncome(),
			MinBPAAmount: resp.GetMinBPAAmount(),
			MinBPAIncome: resp.GetMinBPAIncome(),
		}
	}()
	return out, errChan
}

func (ds *dataService) GetFederalTaxBrackets(ctx context.Context, year int) (<-chan []shared.TaxBracket, <-chan error) {
	//TODO implement me
	panic("implement me")
}

func (ds *dataService) GetCEA(ctx context.Context, year int) (<-chan credits.CanadaEmploymentAmount, <-chan error) {
	//TODO implement me
	panic("implement me")
}

func (ds *dataService) GetBCBPA(ctx context.Context, year int) (<-chan bcCredits.BasicPersonalAmount, <-chan error) {
	//TODO implement me
	panic("implement me")
}

func (ds *dataService) GetBCTaxBrackets(ctx context.Context, year int) (<-chan []shared.TaxBracket, <-chan error) {
	//TODO implement me
	panic("implement me")
}

func (ds *dataService) GetAlbertaBPA(ctx context.Context, year int) (<-chan abCredits.BasicPersonalAmount, <-chan error) {
	//TODO implement me
	panic("implement me")
}

func (ds *dataService) GetAlbertaTaxBrackets(ctx context.Context, year int) (<-chan []shared.TaxBracket, <-chan error) {
	//TODO implement me
	panic("implement me")
}
