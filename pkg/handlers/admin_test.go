package handlers_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/wariusagi/assessment-tax/pkg/handlers"
	"github.com/wariusagi/assessment-tax/pkg/services"
)

type MockAdminService struct {
	res services.AdminDeductionPersonalResponse
	err error
}

func (m *MockAdminService) SetDeduction(req services.AdminDeductionRequest) (services.AdminDeductionPersonalResponse, error) {
	return m.res, m.err
}

func (m *MockAdminService) SetKReceipt(req services.AdminKReceiptRequest) (services.AdminKReceiptResponse, error) {
	return services.AdminKReceiptResponse{}, m.err
}

func callAdminHandler(body io.Reader, mockService *MockAdminService) (*httptest.ResponseRecorder, error) {
	// req
	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/tax/calculations", body)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// handler
	handler := handlers.NewAdminHandler(mockService)
	err := handler.SetDeduction(c)
	return rec, err
}

func TestAdminHandlerSetDeduction_Success(t *testing.T) {
	reqBody := services.AdminDeductionRequest{}
	reqBodyJson, err := json.Marshal(reqBody)
	if err != nil {
		t.Fatalf("Marshalling request body failed: %v", err)
	}

	// mock
	mockService := &MockAdminService{
		res: services.AdminDeductionPersonalResponse{PersonalDeduction: 70000.0},
		err: nil,
	}

	rec, err := callAdminHandler(bytes.NewReader(reqBodyJson), mockService)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)

	resBodyExpected := services.AdminDeductionPersonalResponse{
		PersonalDeduction: 70000.0,
	}
	resBody := services.AdminDeductionPersonalResponse{}
	if err := json.Unmarshal(rec.Body.Bytes(), &resBody); err != nil {
		t.Fatalf("Parsing response body failed: %v", err)
	}
	assert.Equal(t, resBodyExpected, resBody)
}

func TestAdminHandlerSetDeduction_ErrorBindReq(t *testing.T) {
	// req
	reqBodyStr := `{"amount":"mock"}`

	// mock
	mockService := &MockAdminService{
		res: services.AdminDeductionPersonalResponse{},
		err: nil,
	}

	rec, err := callAdminHandler(bytes.NewBufferString(reqBodyStr), mockService)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestAdminHandlerSetDeduction_ErrorSetDeductionService(t *testing.T) {
	// req
	reqBody := services.AdminDeductionRequest{Amount: 70000.0}
	reqBodyJson, err := json.Marshal(reqBody)
	if err != nil {
		t.Fatalf("Marshalling request body failed: %v", err)
	}

	// mock
	mockService := &MockAdminService{
		res: services.AdminDeductionPersonalResponse{},
		err: fmt.Errorf("mock error"),
	}

	rec, err := callAdminHandler(bytes.NewReader(reqBodyJson), mockService)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusInternalServerError, rec.Code)
}
