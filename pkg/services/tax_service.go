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
		case "k-receipt":
			netTotalIncome -= math.Min(a.Amount, data.AmtKReceiptMax)
		default:
			return TaxCalculationResponse{}, fmt.Errorf("not allow type = %v", a.AllowanceType)
		}
	}

	if netTotalIncome < 0 {
		netTotalIncome = 0
	}

	// calculate tax
	var taxTotal float64
	taxLevel := []TaxLevelRes{}
	taxLevelDefault := []TaxLevelRes{}
	idxFinalTax := 0
	for i, lv := range TaxCalculatorLevels {
		taxLevel = append(taxLevel, TaxLevelRes{Level: lv.Text})
		taxLevelDefault = append(taxLevelDefault, TaxLevelRes{Level: lv.Text})
		if netTotalIncome >= lv.Upper {
			tax := (lv.Upper - lv.Lower) * lv.Rate
			if tax > 0 {
				taxLevel[i].Tax = tax
				taxTotal += tax
			}
		} else {
			tax := (netTotalIncome - lv.Lower) * lv.Rate
			taxLevel[i].Tax = tax
			taxTotal += tax
			idxFinalTax = i
			break
		}
	}

	for i := idxFinalTax + 1; i < len(TaxCalculatorLevels); i++ {
		taxLevel = append(taxLevel, TaxLevelRes{Level: TaxCalculatorLevels[i].Text})
		taxLevelDefault = append(taxLevelDefault, TaxLevelRes{Level: TaxCalculatorLevels[i].Text})
	}

	// discount: wht
	var taxRefund float64
	if req.Wht > 0 {
		if req.Wht > taxTotal {
			taxRefund = req.Wht - taxTotal
			taxLevel = taxLevelDefault
			taxTotal = 0

		} else {
			taxTotal -= req.Wht
		}
	}
	return TaxCalculationResponse{Tax: taxTotal, TaxLevel: taxLevel, TaxRefund: taxRefund}, nil
}
