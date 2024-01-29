package shared

import (
	"fmt"
	"github.com/gerdooshell/tax-core/library/mathHelper"
)

const (
	fromYearCPPFirstAdditional  int = 2019
	fromYearCPPSecondAdditional int = 2024
)

type CanadaPensionPlan struct {
	Year                            int
	BasicRateEmployee               float64
	BasicRateEmployer               float64
	FirstAdditionalRateEmployee     float64
	FirstAdditionalRateEmployer     float64
	SecondAdditionalRateEmployee    float64
	SecondAdditionalRateEmployer    float64
	BasicExemption                  float64
	MaxPensionableEarning           float64
	AdditionalMaxPensionableEarning float64
	cppBasicEmployee                float64
	cppBasicEmployer                float64
	cppBasicSelfEmployed            float64
	cppFirstAdditionalEmployee      float64
	cppFirstAdditionalEmployer      float64
	cppFirstAdditionalSelfEmployed  float64
	cppSecondAdditionalEmployee     float64
	cppSecondAdditionalEmployer     float64
	cppSecondAdditionalSelfEmployed float64
}

func (cpp *CanadaPensionPlan) Calculate(totalIncome float64) error {
	if cpp.Year <= 0 {
		return fmt.Errorf("cpp error: invalid year: %v", cpp.Year)
	}
	if totalIncome < 0 {
		return fmt.Errorf("cpp error: invalid income: %v", totalIncome)
	}
	if err := cpp.calculateCppBasic(totalIncome); err != nil {
		return err
	}
	if err := cpp.calculateCppFirst(totalIncome); err != nil {
		return err
	}
	if err := cpp.calculateCppSecond(totalIncome); err != nil {
		return err
	}
	return nil
}

func (cpp *CanadaPensionPlan) GetCPPBasicEmployee() float64 {
	return cpp.cppBasicEmployee
}

func (cpp *CanadaPensionPlan) GetCPPBasicEmployer() float64 {
	return cpp.cppBasicEmployer
}

func (cpp *CanadaPensionPlan) GetCPPBasicSelfEmployed() float64 {
	return cpp.cppBasicSelfEmployed
}

func (cpp *CanadaPensionPlan) GetCPPFirstAdditionalEmployee() float64 {
	return cpp.cppFirstAdditionalEmployee
}

func (cpp *CanadaPensionPlan) GetCPPFirstAdditionalEmployer() float64 {
	return cpp.cppFirstAdditionalEmployer
}

func (cpp *CanadaPensionPlan) GetCPPFirstAdditionalSelfEmployed() float64 {
	return cpp.cppFirstAdditionalSelfEmployed
}

func (cpp *CanadaPensionPlan) GetCPPSecondAdditionalEmployee() float64 {
	return cpp.cppSecondAdditionalEmployee
}

func (cpp *CanadaPensionPlan) GetCPPSecondAdditionalEmployer() float64 {
	return cpp.cppSecondAdditionalEmployer
}

func (cpp *CanadaPensionPlan) GetCPPSecondAdditionalSelfEmployed() float64 {
	return cpp.cppSecondAdditionalSelfEmployed
}

func (cpp *CanadaPensionPlan) calculateCppBasic(totalIncome float64) error {
	if err := cpp.validateCppBasicInputs(); err != nil {
		return nil
	}
	higherValue := max(min(cpp.MaxPensionableEarning, totalIncome)-cpp.BasicExemption, 0)
	employee := higherValue * cpp.BasicRateEmployee / 100
	employer := higherValue * cpp.BasicRateEmployer / 100
	cpp.cppBasicEmployee = mathHelper.RoundFloat64(employee, 2)
	cpp.cppBasicEmployer = mathHelper.RoundFloat64(employer, 2)
	cpp.cppBasicSelfEmployed = mathHelper.RoundFloat64(employee+employer, 2)
	return nil
}

func (cpp *CanadaPensionPlan) validateCppBasicInputs() error {
	if cpp.BasicRateEmployee <= 0 {
		return fmt.Errorf("cpp error: invalid basic employee rate \"%v\" for year \"%v\"", cpp.BasicRateEmployee, cpp.Year)
	}
	if cpp.BasicRateEmployer <= 0 {
		return fmt.Errorf("cpp error: invalid basic employer rate \"%v\" for year \"%v\"", cpp.BasicRateEmployer, cpp.Year)
	}
	return nil
}

func (cpp *CanadaPensionPlan) calculateCppSecond(totalIncome float64) error {
	if err := cpp.validateCppSecondInputs(); err != nil {
		return err
	}
	if totalIncome <= cpp.MaxPensionableEarning || cpp.Year < fromYearCPPSecondAdditional {
		cpp.cppSecondAdditionalEmployee = 0
		cpp.cppSecondAdditionalEmployer = 0
		cpp.cppSecondAdditionalSelfEmployed = 0
		return nil
	}
	higherValue := max(min(cpp.AdditionalMaxPensionableEarning, totalIncome)-cpp.BasicExemption, 0)
	employee := higherValue * cpp.SecondAdditionalRateEmployee / 100
	employer := higherValue * cpp.SecondAdditionalRateEmployer / 100
	cpp.cppSecondAdditionalEmployee = mathHelper.RoundFloat64(employee, 2)
	cpp.cppSecondAdditionalEmployer = mathHelper.RoundFloat64(employer, 2)
	cpp.cppSecondAdditionalSelfEmployed = mathHelper.RoundFloat64(employee+employer, 2)
	return nil
}

func (cpp *CanadaPensionPlan) validateCppSecondInputs() error {
	if cpp.SecondAdditionalRateEmployee <= 0 && cpp.Year >= fromYearCPPSecondAdditional {
		return fmt.Errorf("cpp error: invalid second additional employee rate \"%v\" for year \"%v\"", cpp.FirstAdditionalRateEmployee, cpp.Year)
	}
	if cpp.SecondAdditionalRateEmployer <= 0 && cpp.Year >= fromYearCPPSecondAdditional {
		return fmt.Errorf("cpp error: invalid second additional employer rate \"%v\" for year \"%v\"", cpp.FirstAdditionalRateEmployer, cpp.Year)
	}
	return nil
}

func (cpp *CanadaPensionPlan) calculateCppFirst(totalIncome float64) error {
	if err := cpp.validateCppFirstInputs(); err != nil {
		return err
	}
	if cpp.Year < fromYearCPPFirstAdditional {
		cpp.cppFirstAdditionalEmployee = 0
		cpp.cppFirstAdditionalEmployer = 0
		cpp.cppFirstAdditionalSelfEmployed = 0
		return nil
	}
	higherValue := max(min(cpp.MaxPensionableEarning, totalIncome)-cpp.BasicExemption, 0)
	employee := higherValue * cpp.FirstAdditionalRateEmployee / 100
	employer := higherValue * cpp.FirstAdditionalRateEmployer / 100
	cpp.cppFirstAdditionalEmployee = mathHelper.RoundFloat64(employee, 2)
	cpp.cppFirstAdditionalEmployer = mathHelper.RoundFloat64(employer, 2)
	cpp.cppFirstAdditionalSelfEmployed = mathHelper.RoundFloat64(employee+employer, 2)
	return nil
}

func (cpp *CanadaPensionPlan) validateCppFirstInputs() error {
	if cpp.FirstAdditionalRateEmployee <= 0 && cpp.Year >= fromYearCPPFirstAdditional {
		return fmt.Errorf("cpp error: invalid first additional employee rate \"%v\" for year \"%v\"", cpp.FirstAdditionalRateEmployee, cpp.Year)
	}
	if cpp.FirstAdditionalRateEmployer <= 0 && cpp.Year >= fromYearCPPFirstAdditional {
		return fmt.Errorf("cpp error: invalid first additional employer rate \"%v\" for year \"%v\"", cpp.FirstAdditionalRateEmployer, cpp.Year)
	}
	return nil
}
