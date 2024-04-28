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

	// discount: personal deduction
	netTotalIncome := req.TotalIncome - data.AmtPersonalDeductionMin

	var tax float64
	for _, lv := range TaxLevels {
		if netTotalIncome >= lv.Upper {
			tax += (lv.Upper - lv.Lower) * lv.Rate
		} else {
			tax += (netTotalIncome - lv.Lower) * lv.Rate
			break
		}
	}

	// discount: wht
	if req.Wht > 0 {
		tax -= req.Wht
	}

	return TaxCalculationResponse{Tax: tax}, nil
}
