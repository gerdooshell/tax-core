package dataAccess

import (
	"context"
	"fmt"
	"math"
	"sync"

	abCredits "github.com/gerdooshell/tax-core/entities/canada/alberta/credits"
	bcCredits "github.com/gerdooshell/tax-core/entities/canada/bc/credits"
	federalEntities "github.com/gerdooshell/tax-core/entities/canada/federal/credits"
	sharedEntities "github.com/gerdooshell/tax-core/entities/canada/shared"
	"github.com/gerdooshell/tax-core/environment"
	dataAccessInteractor "github.com/gerdooshell/tax-core/interactors/data_access"
	"github.com/gerdooshell/tax-core/library/cache/lrucache"
	"github.com/gerdooshell/tax-core/library/region/canada"

	dataProvider "github.com/gerdooshell/tax-communication/src/data_provider"
	logger "github.com/gerdooshell/tax-logger-client-go"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var singletonInstance *dataService

type cacheKey struct {
	region   canada.Province
	year     int
	resource string
}

var instanceMu sync.Mutex

func GetDataProviderServiceInstance() DataProviderService {
	instanceMu.Lock()
	defer instanceMu.Unlock()
	if singletonInstance != nil {
		return singletonInstance
	}
	singletonInstance = &dataService{
		dataProviderUrl: getDataProviderUrl(),
		cache:           lrucache.NewLRUCache[cacheKey](1000),
		mu:              make(map[string]*sync.Mutex),
	}
	if err := singletonInstance.generateDataServiceClient(); err != nil {
		logger.Error(err.Error())
	}
	return singletonInstance
}

func getDataProviderUrl() string {
	if environment.GetEnvironment() == environment.Dev {
		return "localhost:45432"
	}
	return "data-provider:45432"
}

type dataService struct {
	dataProviderUrl string
	grpcClient      dataProvider.GRPCDataProviderClient
	cache           lrucache.LRUCache[cacheKey]
	mu              map[string]*sync.Mutex
}

var dialMu sync.Mutex

func (ds *dataService) generateDataServiceClient() error {
	dialMu.Lock()
	defer dialMu.Unlock()
	if ds.grpcClient != nil {
		return nil
	}
	connection, err := grpc.Dial(ds.dataProviderUrl, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		err = fmt.Errorf("connection failed to data provider")
	}
	ds.grpcClient = dataProvider.NewGRPCDataProviderClient(connection)
	//if err = connection.Close(); err != nil {
	//	return nil, fmt.Errorf("failed closing connection, error: %v\n", err)
	//}
	return err
}

var mapMu sync.Mutex

func (ds *dataService) registerToMutex(funcName string) {
	mapMu.Lock()
	defer mapMu.Unlock()
	if _, ok := ds.mu[funcName]; ok {
		return
	}
	var mu sync.Mutex
	ds.mu[funcName] = &mu
}

func (ds *dataService) readFromMutex(funcName string) (*sync.Mutex, bool) {
	mapMu.Lock()
	defer mapMu.Unlock()
	mu, ok := ds.mu[funcName]
	return mu, ok
}

func (ds *dataService) SaveMarginalTaxBrackets(ctx context.Context, province canada.Province, year int, brackets []sharedEntities.TaxBracket) (<-chan bool, <-chan error) {
	funcName := "SaveMarginalTaxBrackets"
	ds.registerToMutex(funcName)
	out := make(chan bool, 1)
	errChan := make(chan error, 1)
	mu, _ := ds.readFromMutex(funcName)
	go func() {
		defer close(out)
		defer close(errChan)
		mu.Lock()
		defer mu.Unlock()
		reqBrackets := make([]*dataProvider.Bracket, 0, len(brackets))
		for _, bracket := range brackets {
			reqBrackets = append(reqBrackets, &dataProvider.Bracket{Low: bracket.Low, Rate: bracket.Rate})
		}
		req := &dataProvider.SaveCombinedMarginalBracketsRequest{
			Year:     int32(year),
			Province: string(province),
			Brackets: reqBrackets,
		}
		resp, err := ds.grpcClient.SaveCombinedMarginalBrackets(ctx, req)
		if err != nil {
			errChan <- err
			return
		}
		if resp == nil {
			errChan <- fmt.Errorf("nil response: PostCombinedMarginalBrackets")
		}
		out <- resp.Success
	}()
	return out, errChan
}

func (ds *dataService) GetCPP(ctx context.Context, year int) (<-chan sharedEntities.CanadaPensionPlan, <-chan error) {
	funcName := "GetCPP"
	ds.registerToMutex(funcName)
	out := make(chan sharedEntities.CanadaPensionPlan, 1)
	errChan := make(chan error, 1)
	mu, _ := ds.readFromMutex(funcName)
	go func() {
		defer close(out)
		defer close(errChan)
		mu.Lock()
		defer mu.Unlock()
		cppCacheKey := cacheKey{region: canada.Federal, year: year, resource: funcName}
		if value, cacheErr := ds.cache.Read(cppCacheKey); cacheErr == nil {
			out <- value.(sharedEntities.CanadaPensionPlan)
			return
		}
		req := &dataProvider.GetCanadaPensionPlanRequest{
			Year: int32(year),
		}
		resp, err := ds.grpcClient.GetCanadaPensionPlan(ctx, req)
		if err != nil {
			errChan <- err
			return
		}
		value := sharedEntities.CanadaPensionPlan{
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
		if _, err = ds.cache.Add(cppCacheKey, value); err != nil {
			errChan <- err
			return
		}
		out <- value
	}()
	return out, errChan
}

func (ds *dataService) GetEIPremium(ctx context.Context, year int) <-chan dataAccessInteractor.EIPremiumDataOut {
	funcName := "GetEIPremium"
	ds.registerToMutex(funcName)
	outChan := make(chan dataAccessInteractor.EIPremiumDataOut, 1)
	mu, _ := ds.readFromMutex(funcName)
	go func() {
		defer close(outChan)
		out := dataAccessInteractor.EIPremiumDataOut{}
		defer func() { outChan <- out }()
		mu.Lock()
		defer mu.Unlock()
		eipCacheKey := cacheKey{region: canada.Federal, year: year, resource: funcName}
		if value, cacheErr := ds.cache.Read(eipCacheKey); cacheErr == nil {
			out.EmploymentInsurancePremium = value.(sharedEntities.EmploymentInsurancePremium)
			return
		}
		req := &dataProvider.GetEmploymentInsurancePremiumRequest{
			Year: int32(year),
		}
		resp, err := ds.grpcClient.GetEmploymentInsurancePremium(ctx, req)
		if err != nil {
			out.Err = err
			return
		}
		value := sharedEntities.EmploymentInsurancePremium{
			Rate:                              resp.GetRate(),
			MaxInsurableEarning:               resp.GetMaxInsurableEarning(),
			EmployerEmployeeContributionRatio: resp.GetEmployerEmployeeContributionRatio(),
		}
		if _, err = ds.cache.Add(eipCacheKey, value); err != nil {
			out.Err = err
			return
		}
		out.EmploymentInsurancePremium = value
	}()
	return outChan
}

func (ds *dataService) GetFederalBPA(ctx context.Context, year int) <-chan dataAccessInteractor.FederalBPADataOut {
	funcName := "GetFederalBPA"
	ds.registerToMutex(funcName)
	out := make(chan dataAccessInteractor.FederalBPADataOut, 1)
	mu, _ := ds.readFromMutex(funcName)
	go func() {
		defer close(out)
		bpaOut := dataAccessInteractor.FederalBPADataOut{}
		defer func() { out <- bpaOut }()
		mu.Lock()
		defer mu.Unlock()
		bpaCacheKey := cacheKey{region: canada.Federal, year: year, resource: funcName}
		if value, cacheErr := ds.cache.Read(bpaCacheKey); cacheErr == nil {
			bpaOut.BasicPersonalAmount = value.(federalEntities.BasicPersonalAmount)
			return
		}
		req := &dataProvider.GetFederalBPARequest{
			Year: int32(year),
		}
		resp, err := ds.grpcClient.GetFederalBPA(ctx, req)
		if err != nil {
			bpaOut.Err = err
			return
		}
		value := federalEntities.BasicPersonalAmount{
			MaxBPAAmount: resp.GetMaxBPAAmount(),
			MaxBPAIncome: resp.GetMaxBPAIncome(),
			MinBPAAmount: resp.GetMinBPAAmount(),
			MinBPAIncome: resp.GetMinBPAIncome(),
		}
		if _, err = ds.cache.Add(bpaCacheKey, value); err != nil {
			bpaOut.Err = err
			return
		}
		bpaOut.BasicPersonalAmount = value
	}()
	return out
}

func (ds *dataService) GetTaxBrackets(ctx context.Context, year int, province canada.Province) <-chan dataAccessInteractor.TaxBracketsDataOut {
	funcName := "GetTaxBrackets"
	ds.registerToMutex(funcName)
	outChan := make(chan dataAccessInteractor.TaxBracketsDataOut, 1)
	mu, _ := ds.readFromMutex(funcName)
	go func() {
		defer close(outChan)
		out := dataAccessInteractor.TaxBracketsDataOut{}
		defer func() { outChan <- out }()
		mu.Lock()
		defer mu.Unlock()
		bracketsCacheKey := cacheKey{region: province, year: year, resource: funcName}
		if value, cacheErr := ds.cache.Read(bracketsCacheKey); cacheErr == nil {
			out.TaxBrackets = value.([]sharedEntities.TaxBracket)
			return
		}
		req := &dataProvider.GetTaxBracketsRequest{
			Year:     int32(year),
			Province: string(province),
		}
		resp, err := ds.grpcClient.GetTaxBrackets(ctx, req)
		if err != nil {
			out.Err = err
			return
		}
		respBrackets := resp.GetBrackets()
		brackets := make([]sharedEntities.TaxBracket, 0, len(respBrackets))
		for i, rb := range respBrackets {
			bracket := sharedEntities.TaxBracket{Rate: rb.GetRate(), Low: rb.GetLow()}
			if i < len(respBrackets)-1 {
				bracket.High = respBrackets[i+1].Low
			} else {
				bracket.High = math.MaxFloat64
			}
			brackets = append(brackets, bracket)
		}
		if _, err = ds.cache.Add(bracketsCacheKey, brackets); err != nil {
			out.Err = err
			return
		}
		out.TaxBrackets = brackets
	}()
	return outChan
}

func (ds *dataService) GetCombinedMarginalBrackets(ctx context.Context, year int, province canada.Province) (<-chan []sharedEntities.TaxBracket, <-chan error) {
	funcName := "GetCombinedMarginalBrackets"
	ds.registerToMutex(funcName)
	out := make(chan []sharedEntities.TaxBracket, 1)
	errChan := make(chan error, 1)
	mu, _ := ds.readFromMutex(funcName)
	go func() {
		defer close(out)
		defer close(errChan)
		mu.Lock()
		defer mu.Unlock()
		marginalCacheKey := cacheKey{region: province, year: year, resource: funcName}
		if value, cacheErr := ds.cache.Read(marginalCacheKey); cacheErr == nil {
			out <- value.([]sharedEntities.TaxBracket)
			return
		}
		req := &dataProvider.GetCombinedMarginalBracketsRequest{
			Year:     int32(year),
			Province: string(province),
		}
		resp, err := ds.grpcClient.GetCombinedMarginalBrackets(ctx, req)
		if err != nil {
			errChan <- err
			return
		}
		respBrackets := resp.GetBrackets()
		brackets := make([]sharedEntities.TaxBracket, 0, len(respBrackets))
		for i, rb := range respBrackets {
			bracket := sharedEntities.TaxBracket{Rate: rb.GetRate(), Low: rb.GetLow()}
			if i < len(respBrackets)-1 {
				bracket.High = respBrackets[i+1].Low
			} else {
				bracket.High = math.MaxFloat64
			}
			brackets = append(brackets, bracket)
		}
		if _, err = ds.cache.Add(marginalCacheKey, brackets); err != nil {
			errChan <- err
			return
		}
		out <- brackets
	}()
	return out, errChan
}

func (ds *dataService) GetCEA(ctx context.Context, year int) (<-chan federalEntities.CanadaEmploymentAmount, <-chan error) {
	funcName := "GetCEA"
	ds.registerToMutex(funcName)
	out := make(chan federalEntities.CanadaEmploymentAmount, 1)
	errChan := make(chan error, 1)
	mu, _ := ds.readFromMutex(funcName)
	go func() {
		defer close(out)
		defer close(errChan)
		mu.Lock()
		defer mu.Unlock()
		ceaCacheKey := cacheKey{region: canada.Federal, year: year, resource: funcName}
		if value, cacheErr := ds.cache.Read(ceaCacheKey); cacheErr == nil {
			out <- value.(federalEntities.CanadaEmploymentAmount)
			return
		}
		req := &dataProvider.GetCanadaEmploymentAmountRequest{
			Year: int32(year),
		}
		resp, err := ds.grpcClient.GetCEA(ctx, req)
		if err != nil {
			errChan <- err
			return
		}
		value := federalEntities.CanadaEmploymentAmount{Value: resp.GetCeaValue()}
		if _, err = ds.cache.Add(ceaCacheKey, value); err != nil {
			errChan <- err
			return
		}
		out <- value
	}()
	return out, errChan
}

func (ds *dataService) GetBCBPA(ctx context.Context, year int) (<-chan bcCredits.BasicPersonalAmount, <-chan error) {
	funcName := "GetBCBPA"
	ds.registerToMutex(funcName)
	out := make(chan bcCredits.BasicPersonalAmount, 1)
	errChan := make(chan error, 1)
	mu, _ := ds.readFromMutex(funcName)
	go func() {
		defer close(out)
		defer close(errChan)
		mu.Lock()
		defer mu.Unlock()
		bpaCacheKey := cacheKey{region: canada.BritishColumbia, year: year, resource: funcName}
		if value, cacheErr := ds.cache.Read(bpaCacheKey); cacheErr == nil {
			out <- value.(bcCredits.BasicPersonalAmount)
			return
		}
		req := &dataProvider.GetBritishColumbiaBPARequest{
			Year: int32(year),
		}
		resp, err := ds.grpcClient.GetBritishColumbiaBPA(ctx, req)
		if err != nil {
			errChan <- err
			return
		}
		value := bcCredits.BasicPersonalAmount{Value: resp.GetBpaValue()}
		if _, err = ds.cache.Add(bpaCacheKey, value); err != nil {
			errChan <- err
			return
		}
		out <- value
	}()
	return out, errChan
}

func (ds *dataService) GetAlbertaBPA(ctx context.Context, year int) (<-chan abCredits.BasicPersonalAmount, <-chan error) {
	funcName := "GetAlbertaBPA"
	ds.registerToMutex(funcName)
	out := make(chan abCredits.BasicPersonalAmount, 1)
	errChan := make(chan error, 1)
	mu, _ := ds.readFromMutex(funcName)
	go func() {
		defer close(out)
		defer close(errChan)
		mu.Lock()
		defer mu.Unlock()
		bpaCacheKey := cacheKey{region: canada.Alberta, year: year, resource: funcName}
		if value, cacheErr := ds.cache.Read(bpaCacheKey); cacheErr == nil {
			out <- value.(abCredits.BasicPersonalAmount)
			return
		}
		req := &dataProvider.GetAlbertaBPARequest{
			Year: int32(year),
		}
		resp, err := ds.grpcClient.GetAlbertaBPA(ctx, req)
		if err != nil {
			errChan <- err
			return
		}
		value := abCredits.BasicPersonalAmount{Value: resp.GetBpaValue()}
		if _, err = ds.cache.Add(bpaCacheKey, value); err != nil {
			errChan <- err
			return
		}
		out <- value
	}()
	return out, errChan

}

func (ds *dataService) GetRRSP(ctx context.Context, year int) (<-chan sharedEntities.RRSP, <-chan error) {
	funcName := "GetRRSP"
	ds.registerToMutex(funcName)
	out := make(chan sharedEntities.RRSP, 1)
	errChan := make(chan error, 1)
	mu, _ := ds.readFromMutex(funcName)
	go func() {
		defer close(out)
		defer close(errChan)
		mu.Lock()
		defer mu.Unlock()
		rrspCacheKey := cacheKey{region: canada.Federal, year: year, resource: funcName}
		if value, cacheErr := ds.cache.Read(rrspCacheKey); cacheErr == nil {
			out <- value.(sharedEntities.RRSP)
			return
		}
		req := &dataProvider.GetRegisteredRetirementSavingsPlanRequest{
			Year: int32(year),
		}
		resp, err := ds.grpcClient.GetRegisteredRetirementSavingsPlan(ctx, req)
		if err != nil {
			errChan <- err
			return
		}
		value := sharedEntities.RRSP{
			Year:                  year,
			Rate:                  resp.GetRate(),
			MaxContribution:       resp.GetMaxContribution(),
			OverContributionRate:  resp.GetOverContributionRate(),
			OverContributionLimit: resp.GetOverContributionLimit(),
		}
		if _, err = ds.cache.Add(rrspCacheKey, value); err != nil {
			errChan <- err
			return
		}
		out <- value
	}()
	return out, errChan
}
