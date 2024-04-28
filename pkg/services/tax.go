package services

import "math"

type Allowance struct {
	AllowanceType string  `json:"allowanceType"`
	Amount        float64 `json:"amount"`
}

type TaxCalculationRequest struct {
	TotalIncome float64     `json:"totalIncome"`
	Wht         float64     `json:"wht"`
	Allowances  []Allowance `json:"allowances"`
}

type TaxCalculationResponse struct {
	Tax      float64       `json:"tax"`
	TaxLevel []TaxLevelRes `json:"taxLevel"`
}

type TaxLevelRes struct {
	Level string  `json:"level"`
	Tax   float64 `json:"tax"`
}

type TaxCalculatorLevel struct {
	Lower float64
	Upper float64
	Rate  float64
	Text  string
}

var TaxCalculatorLevels = []TaxCalculatorLevel{
	{0, 150000, 0, "0-150,000"}, // exempt
	{150000, 500000, 0.10, "150,001-500,000"},
	{500000, 1000000, 0.15, "500,001-1,000,000"},
	{1000000, 2000000, 0.20, "1,000,001-2,000,000"},
	{2000000, math.MaxFloat64, 0.35, "2,000,001 ขึ้นไป"},
}

type TaxService interface {
	CalculateTax(req TaxCalculationRequest) (TaxCalculationResponse, error)
}
