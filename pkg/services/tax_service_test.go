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

func callTaxService(req services.TaxCalculationRequest, mockRepo *MockRepo) (services.TaxCalculationResponse, error) {
	service := services.NewTaxService(mockRepo)
	return service.CalculateTax(req)

}

func TestTaxServiceCalculateTax_Success(t *testing.T) {
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

	res, err := callTaxService(req, mockRepo)

	assert.NoError(t, err)
	resTaxExpected := 29000.0
	assert.Equal(t, resTaxExpected, res.Tax)
}

func TestTaxServiceCalculateTax_SuccessWithDiscountWht(t *testing.T) {
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

	res, err := callTaxService(req, mockRepo)

	assert.NoError(t, err)
	resTaxExpected := 4000.0
	assert.Equal(t, resTaxExpected, res.Tax)
}

func TestTaxServiceCalculateTax_SuccessWithDiscountDonation(t *testing.T) {
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

	res, err := callTaxService(req, mockRepo)

	assert.NoError(t, err)
	resTaxExpected := 19000.0
	assert.Equal(t, resTaxExpected, res.Tax)
}

func TestTaxServiceCalculateTax_SuccessWithResponseTaxLevel(t *testing.T) {
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

	res, err := callTaxService(req, mockRepo)

	assert.NoError(t, err)
	resTaxExpected := 19000.0
	assert.Equal(t, resTaxExpected, res.Tax)
	assert.Equal(t, 5, len(res.TaxLevel))
	assert.Equal(t, 0.0, res.TaxLevel[0].Tax)
	assert.Equal(t, 19000.0, res.TaxLevel[1].Tax)
	assert.Equal(t, 0.0, res.TaxLevel[2].Tax)
	assert.Equal(t, 0.0, res.TaxLevel[3].Tax)
	assert.Equal(t, 0.0, res.TaxLevel[4].Tax)
}

func TestTaxServiceCalculateTax_SuccessWithResponseTaxLevelHaveTaxInSeveralLevel(t *testing.T) {
	// mock
	mockRepo := &MockRepo{
		data: database.MasterTaxDeduction{
			AmtPersonalDeductionMin: 60000,
			AmtDonationMax:          100000,
		},
		err: nil,
	}
	req := services.TaxCalculationRequest{
		TotalIncome: 2120000.0,
		Wht:         0.0,
		Allowances: []services.Allowance{
			{
				AllowanceType: "donation",
				Amount:        50000.0,
			},
		},
	}

	res, err := callTaxService(req, mockRepo)

	assert.NoError(t, err)
	resTaxExpected := 313500.0
	assert.Equal(t, resTaxExpected, res.Tax)
	assert.Equal(t, 5, len(res.TaxLevel))
	assert.Equal(t, 0.0, res.TaxLevel[0].Tax)
	assert.Equal(t, 35000.0, res.TaxLevel[1].Tax)
	assert.Equal(t, 75000.0, res.TaxLevel[2].Tax)
	assert.Equal(t, 200000.0, res.TaxLevel[3].Tax)
	assert.Equal(t, 3500.0, res.TaxLevel[4].Tax)
}

func TestTaxServiceCalculateTax_SuccessWithResponseTaxLevelHaveTaxZero(t *testing.T) {
	// mock
	mockRepo := &MockRepo{
		data: database.MasterTaxDeduction{
			AmtPersonalDeductionMin: 60000,
			AmtDonationMax:          100000,
		},
		err: nil,
	}
	req := services.TaxCalculationRequest{
		TotalIncome: 155000.0,
		Wht:         0.0,
		Allowances: []services.Allowance{
			{
				AllowanceType: "donation",
				Amount:        97000.0,
			},
		},
	}

	res, err := callTaxService(req, mockRepo)

	assert.NoError(t, err)
	resTaxExpected := 0.0
	assert.Equal(t, resTaxExpected, res.Tax)
	assert.Equal(t, 5, len(res.TaxLevel))
	assert.Equal(t, 0.0, res.TaxLevel[0].Tax)
	assert.Equal(t, 0.0, res.TaxLevel[1].Tax)
	assert.Equal(t, 0.0, res.TaxLevel[2].Tax)
	assert.Equal(t, 0.0, res.TaxLevel[3].Tax)
	assert.Equal(t, 0.0, res.TaxLevel[4].Tax)
}

func TestTaxServiceCalculateTax_ErrorGetDB(t *testing.T) {
	// mock
	mockRepo := &MockRepo{
		data: database.MasterTaxDeduction{},
		err:  fmt.Errorf("mock error"),
	}
	req := services.TaxCalculationRequest{
		TotalIncome: 500000.0,
	}

	_, err := callTaxService(req, mockRepo)

	assert.Error(t, err)
}

func TestTaxServiceCalculateTax_ErrorAllowanceTypeNotAllow(t *testing.T) {
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

	_, err := callTaxService(req, mockRepo)

	assert.Error(t, err)
}
