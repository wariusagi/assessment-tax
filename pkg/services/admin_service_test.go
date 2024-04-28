package services_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/wariusagi/assessment-tax/pkg/database"
	"github.com/wariusagi/assessment-tax/pkg/services"
)

func (m *MockRepo) UpdateAmtPersonalDeductionDeduction(cycleYear int, amtPersonalDeduction float64) error {
	return m.err
}

func (m *MockRepo) UpdateAmtKReceiptDeduction(cycleYear int, amtKReceipt float64) error {
	return m.err
}

func callAdminService(req services.AdminDeductionRequest, mockRepo *MockRepo) (services.AdminDeductionPersonalResponse, error) {
	service := services.NewAdminService(mockRepo)
	return service.SetDeduction(req)

}

func TestAdminServiceSetDeduction_Success(t *testing.T) {
	// mock
	mockRepo := &MockRepo{
		data: database.MasterTaxDeduction{},
		err:  nil,
	}
	req := services.AdminDeductionRequest{
		Amount: 70000.0,
	}

	res, err := callAdminService(req, mockRepo)

	assert.NoError(t, err)
	resPersonalDeductionExpected := 70000.0
	assert.Equal(t, resPersonalDeductionExpected, res.PersonalDeduction)
}

func TestAdminServiceSetDeduction_SuccessWithAmountLimitConfig(t *testing.T) {
	// mock
	mockRepo := &MockRepo{
		data: database.MasterTaxDeduction{},
		err:  nil,
	}
	req := services.AdminDeductionRequest{
		Amount: 200000.0,
	}

	res, err := callAdminService(req, mockRepo)

	assert.NoError(t, err)
	resPersonalDeductionExpected := 100000.0
	assert.Equal(t, resPersonalDeductionExpected, res.PersonalDeduction)
}

func TestAdminServiceSetDeduction_ErrorAmountLessThanConfig(t *testing.T) {
	// mock
	mockRepo := &MockRepo{
		data: database.MasterTaxDeduction{},
		err:  fmt.Errorf("mock error"),
	}
	req := services.AdminDeductionRequest{
		Amount: 5000.0,
	}

	_, err := callAdminService(req, mockRepo)

	assert.Error(t, err)
}

func TestAdminServiceSetDeduction_ErrorUpdateDB(t *testing.T) {
	// mock
	mockRepo := &MockRepo{
		data: database.MasterTaxDeduction{},
		err:  fmt.Errorf("mock error"),
	}
	req := services.AdminDeductionRequest{
		Amount: 70000.0,
	}

	_, err := callAdminService(req, mockRepo)

	assert.Error(t, err)
}
