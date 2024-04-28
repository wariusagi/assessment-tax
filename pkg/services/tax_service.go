package services

import (
	"github.com/wariusagi/assessment-tax/pkg/database"
)

type taxService struct {
	r database.Repository
}

func NewTaxService(r database.Repository) TaxService {
	return taxService{r: r}
}

func (s taxService) CalculateTax(req TaxCalculationRequest) (TaxCalculationResponse, error) {
	curYear := 2567
	data, err := s.r.GetMasterTaxDeduction(curYear)
	if err != nil {
		return TaxCalculationResponse{}, err
	}

	taxTotalIncome := req.TotalIncome - data.AmtPersonalDeductionMin

	var tax float64
	for _, lv := range TaxLevels {
		if taxTotalIncome >= lv.Upper {
			tax += (lv.Upper - lv.Lower) * lv.Rate
		} else {
			tax += (taxTotalIncome - lv.Lower) * lv.Rate
			break
		}
	}

	return TaxCalculationResponse{Tax: tax}, nil
}
