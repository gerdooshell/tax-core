package modulePlugin

import (
	"context"
	"github.com/gerdooshell/tax-core/interactors/internal/tax_calculator/canada/regional/module_plugin/data_structures"
)

type TaxCalculator interface {
	CalculateRegionalTax(context.Context, *dataStructures.Input) (dataStructures.Output, error)
}
