package services

type taxService struct {
	repo string
}

func NewTax(repo string) TaxService {
	return taxService{repo: repo}
}

func (s taxService) NewTaxCalculation(req TaxCalculationRequest) (TaxCalculationResponse, error) {
	return TaxCalculationResponse{}, nil
}
