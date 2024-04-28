package services

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

type TaxService interface {
	NewTaxCalculation(req TaxCalculationRequest) (TaxCalculationResponse, error)
}
