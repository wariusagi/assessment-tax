package services

type AdminDeductionRequest struct {
	Amount float64 `json:"amount"`
}

type AdminDeductionPersonalResponse struct {
	PersonalDeduction float64 `json:"personalDeduction"`
}

type AdminService interface {
	SetDeduction(req AdminDeductionRequest) (AdminDeductionPersonalResponse, error)
}
