package services

import (
	"fmt"

	"github.com/wariusagi/assessment-tax/pkg/database"
)

type adminService struct {
	r database.Repository
}

const (
	personalDeductionMin = 10000
	personalDeductionMax = 100000
)

func NewAdminService(r database.Repository) AdminService {
	return adminService{r: r}
}

func (s adminService) SetDeduction(req AdminDeductionRequest) (AdminDeductionPersonalResponse, error) {
	personalDeduction := req.Amount
	if req.Amount <= personalDeductionMin {
		return AdminDeductionPersonalResponse{}, fmt.Errorf("amount must be more than : %v", personalDeductionMin)
	} else if req.Amount > personalDeductionMax {
		personalDeduction = personalDeductionMax
	}

	err := s.r.UpdateAmtPersonalDeductionDeduction(curYear, personalDeduction)
	if err != nil {
		return AdminDeductionPersonalResponse{}, err
	}

	return AdminDeductionPersonalResponse{PersonalDeduction: personalDeduction}, nil
}
