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
	Tax float64 `json:"tax"`
}

type TaxLevel struct {
	Lower float64
	Upper float64
	Rate  float64
}

var TaxLevels = []TaxLevel{
	{0, 150000, 0}, // exempt
	{150000, 500000, 0.10},
	{500000, 1000000, 0.15},
	{1000000, 2000000, 0.20},
	{2000000, math.MaxFloat64, 0.35},
}

type TaxService interface {
	CalculateTax(req TaxCalculationRequest) (TaxCalculationResponse, error)
}
