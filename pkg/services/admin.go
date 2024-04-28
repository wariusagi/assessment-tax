package services

type AdminDeductionRequest struct {
	Amount float64 `json:"amount"`
}

type AdminKReceiptRequest struct {
	Amount float64 `json:"amount"`
}

type AdminDeductionPersonalResponse struct {
	PersonalDeduction float64 `json:"personalDeduction"`
}

type AdminKReceiptResponse struct {
	KReceipt float64 `json:"kReceipt"`
}

type AdminService interface {
	SetDeduction(req AdminDeductionRequest) (AdminDeductionPersonalResponse, error)
	SetKReceipt(req AdminKReceiptRequest) (AdminKReceiptResponse, error)
}
