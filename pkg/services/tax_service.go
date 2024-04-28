package services

import (
	"fmt"
	"math"
	"strings"

	"github.com/wariusagi/assessment-tax/pkg/database"
)

type taxService struct {
	r database.Repository
}

const curYear = 2567

func NewTaxService(r database.Repository) TaxService {
	return taxService{r: r}
}

func (s taxService) CalculateTax(req TaxCalculationRequest) (TaxCalculationResponse, error) {
	data, err := s.r.GetMasterTaxDeduction(curYear)
	if err != nil {
		return TaxCalculationResponse{}, err
	}

	// discount: personal deduction
	netTotalIncome := req.TotalIncome - data.AmtPersonalDeductionMin

	// discount: allowances
	for _, a := range req.Allowances {
		switch strings.ToLower(a.AllowanceType) {
		case "donation":
			netTotalIncome -= math.Min(a.Amount, data.AmtDonationMax)
		default:
			return TaxCalculationResponse{}, fmt.Errorf("not allow type = %v", a.AllowanceType)
		}
	}

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
