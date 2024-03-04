package regionalBPA

import (
	"context"
	"errors"
	"time"

	dataProvider "github.com/gerdooshell/tax-core/data-access"
	dataAccess "github.com/gerdooshell/tax-core/interactors/data_access"
	"github.com/gerdooshell/tax-core/library/region/canada"
)

type RegionalBasicPersonalAmountInteractor interface {
	GetRegionalBPA(ctx context.Context, year int, totalIncome float64, province canada.Province) <-chan RegionalBasicPersonalAmountOutput
}

type RegionalBasicPersonalAmountOutput struct {
	Value float64
	Err   error
}

func NewRegionalBasicPersonalAmountInteractor() RegionalBasicPersonalAmountInteractor {
	return &regionalBPAImpl{
		dataProvider: dataProvider.GetDataProviderServiceInstance(),
		timeout:      time.Second * 10,
	}
}

type regionalBPAImpl struct {
	dataProvider dataAccess.RegionalBPAData
	timeout      time.Duration
}

func (r *regionalBPAImpl) GetRegionalBPA(ctx context.Context, year int, totalIncome float64, province canada.Province) <-chan RegionalBasicPersonalAmountOutput {
	out := make(chan RegionalBasicPersonalAmountOutput, 1)
	go func() {
		defer close(out)
		regionalBPAOutput := RegionalBasicPersonalAmountOutput{}
		defer func() { out <- regionalBPAOutput }()
		switch province {
		case canada.Alberta:
			regionalBPAOutput = r.getAlbertaBPA(ctx, year)
		case canada.BritishColumbia:
			regionalBPAOutput = r.getBritishColumbiaBPA(ctx, year)
		}
	}()
	return out
}

func (r *regionalBPAImpl) getBritishColumbiaBPA(ctx context.Context, year int) (regionalBPAOutput RegionalBasicPersonalAmountOutput) {
	bcChan := r.dataProvider.GetBritishColumbiaBPA(ctx, year)
	select {
	case bcBPADataOut := <-bcChan:
		if bcBPADataOut.Err != nil {
			regionalBPAOutput.Err = bcBPADataOut.Err
			return
		}
		if regionalBPAOutput.Err = bcBPADataOut.BasicPersonalAmount.Calculate(); regionalBPAOutput.Err != nil {
			return
		}
		regionalBPAOutput.Value = bcBPADataOut.BasicPersonalAmount.GetValue()
		return
	case <-time.After(r.timeout):
		regionalBPAOutput.Err = errors.New("get bc bpa timed out")
	}
	return
}

func (r *regionalBPAImpl) getAlbertaBPA(ctx context.Context, year int) (regionalBPAOutput RegionalBasicPersonalAmountOutput) {
	abChan := r.dataProvider.GetAlbertaBPA(ctx, year)
	select {
	case abBPADataOut := <-abChan:
		if abBPADataOut.Err != nil {
			regionalBPAOutput.Err = abBPADataOut.Err
			return
		}
		if regionalBPAOutput.Err = abBPADataOut.BasicPersonalAmount.Calculate(); regionalBPAOutput.Err != nil {
			return
		}
		regionalBPAOutput.Value = abBPADataOut.BasicPersonalAmount.GetValue()
	case <-time.After(r.timeout):
		regionalBPAOutput.Err = errors.New("get alberta bpa timed out")
	}
	return
}
