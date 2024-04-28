package services_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/wariusagi/assessment-tax/pkg/database"
	"github.com/wariusagi/assessment-tax/pkg/services"
)

type MockRepo struct {
	data database.MasterTaxDeduction
	err  error
}

func (m *MockRepo) GetMasterTaxDeduction(year int) (database.MasterTaxDeduction, error) {
	return m.data, m.err
}

func callService(req services.TaxCalculationRequest, mockRepo *MockRepo) (services.TaxCalculationResponse, error) {
	service := services.NewTaxService(mockRepo)
	return service.CalculateTax(req)

}

func TestServiceCalculateTax_Success(t *testing.T) {
	// mock
	mockRepo := &MockRepo{
		data: database.MasterTaxDeduction{
			AmtPersonalDeductionMin: 60000,
		},
		err: nil,
	}
	req := services.TaxCalculationRequest{
		TotalIncome: 500000.0,
	}

	res, err := callService(req, mockRepo)

	assert.NoError(t, err)
	resTaxExpected := 29000.0
	assert.Equal(t, resTaxExpected, res.Tax)
}

func TestServiceCalculateTax_SuccessWithDiscountWht(t *testing.T) {
	// mock
	mockRepo := &MockRepo{
		data: database.MasterTaxDeduction{
			AmtPersonalDeductionMin: 60000,
		},
		err: nil,
	}
	req := services.TaxCalculationRequest{
		TotalIncome: 500000.0,
		Wht:         25000.0,
	}

	res, err := callService(req, mockRepo)

	assert.NoError(t, err)
	resTaxExpected := 4000.0
	assert.Equal(t, resTaxExpected, res.Tax)
}

func TestServiceCalculateTax_SuccessWithDiscountDonation(t *testing.T) {
	// mock
	mockRepo := &MockRepo{
		data: database.MasterTaxDeduction{
			AmtPersonalDeductionMin: 60000,
			AmtDonationMax:          100000,
		},
		err: nil,
	}
	req := services.TaxCalculationRequest{
		TotalIncome: 500000.0,
		Wht:         0.0,
		Allowances: []services.Allowance{
			{
				AllowanceType: "donation",
				Amount:        200000.0,
			},
		},
	}

	res, err := callService(req, mockRepo)

	assert.NoError(t, err)
	resTaxExpected := 19000.0
	assert.Equal(t, resTaxExpected, res.Tax)
}

func TestServiceCalculateTax_ErrorGetDB(t *testing.T) {
	// mock
	mockRepo := &MockRepo{
		data: database.MasterTaxDeduction{},
		err:  fmt.Errorf("mock error"),
	}
	req := services.TaxCalculationRequest{
		TotalIncome: 500000.0,
	}

	_, err := callService(req, mockRepo)

	assert.Error(t, err)
}

func TestServiceCalculateTax_ErrorAllowanceTypeNotAllow(t *testing.T) {
	// mock
	mockRepo := &MockRepo{
		data: database.MasterTaxDeduction{
			AmtPersonalDeductionMin: 60000,
			AmtDonationMax:          100000,
		},
		err: nil,
	}
	req := services.TaxCalculationRequest{
		TotalIncome: 500000.0,
		Wht:         0.0,
		Allowances: []services.Allowance{
			{
				AllowanceType: "test",
				Amount:        200000.0,
			},
		},
	}

	_, err := callService(req, mockRepo)

	assert.Error(t, err)
}
