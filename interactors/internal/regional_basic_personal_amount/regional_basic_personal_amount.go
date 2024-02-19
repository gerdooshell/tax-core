package regionalBPA

import (
	"context"
	"fmt"

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
	}
}

type regionalBPAImpl struct {
	dataProvider dataAccess.RegionalBPAData
}

func (r *regionalBPAImpl) GetRegionalBPA(ctx context.Context, year int, totalIncome float64, province canada.Province) <-chan RegionalBasicPersonalAmountOutput {
	out := make(chan RegionalBasicPersonalAmountOutput)
	go func() {
		defer close(out)
		regionalBPAOutput := RegionalBasicPersonalAmountOutput{}
		defer func() { out <- regionalBPAOutput }()
		switch province {
		case canada.Alberta:
			abChan, errChan := r.dataProvider.GetAlbertaBPA(ctx, year)
			select {
			case regionalBPAOutput.Err = <-errChan:
				return
			case abBPA := <-abChan:
				if regionalBPAOutput.Err = abBPA.Calculate(); regionalBPAOutput.Err != nil {
					return
				}
				regionalBPAOutput.Value = abBPA.GetValue()
				return
			case <-ctx.Done():
				regionalBPAOutput.Err = fmt.Errorf("processing regional bpa canceled")
			}
		case canada.BritishColumbia:
			bcChan, errChan := r.dataProvider.GetBCBPA(ctx, year)
			select {
			case regionalBPAOutput.Err = <-errChan:
				return
			case bcBPA := <-bcChan:
				if regionalBPAOutput.Err = bcBPA.Calculate(); regionalBPAOutput.Err != nil {
					return
				}
				regionalBPAOutput.Value = bcBPA.GetValue()
				return
			case <-ctx.Done():
				regionalBPAOutput.Err = fmt.Errorf("processing regional bpa canceled")
			}
		}
	}()
	return out
}
