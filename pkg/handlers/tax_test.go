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

type MockTaxService struct {
	res services.TaxCalculationResponse
	err error
}

func (m *MockTaxService) CalculateTax(req services.TaxCalculationRequest) (services.TaxCalculationResponse, error) {
	return m.res, m.err
}

func callTaxHandler(body io.Reader, mockService *MockTaxService) (*httptest.ResponseRecorder, error) {
	// req
	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/tax/calculations", body)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// handler
	handler := handlers.NewTaxHandler(mockService)
	err := handler.CalculateTax(c)
	return rec, err
}

func TestTaxHandlerCalculateTax_Success(t *testing.T) {
	reqBody := services.TaxCalculationRequest{}
	reqBodyJson, err := json.Marshal(reqBody)
	if err != nil {
		t.Fatalf("Marshalling request body failed: %v", err)
	}

	// mock
	mockService := &MockTaxService{
		res: services.TaxCalculationResponse{},
		err: nil,
	}

	rec, err := callTaxHandler(bytes.NewReader(reqBodyJson), mockService)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)

	resBodyExpected := services.TaxCalculationResponse{}
	resBody := services.TaxCalculationResponse{}
	if err := json.Unmarshal(rec.Body.Bytes(), &resBody); err != nil {
		t.Fatalf("Parsing response body failed: %v", err)
	}
	assert.Equal(t, resBodyExpected, resBody)
}

func TestTaxHandlerCalculateTax_ErrorBindReq(t *testing.T) {
	// req
	reqBodyStr := `{"wht":"mock"}`

	// mock
	mockService := &MockTaxService{
		res: services.TaxCalculationResponse{},
		err: nil,
	}

	rec, err := callTaxHandler(bytes.NewBufferString(reqBodyStr), mockService)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestTaxHandlerCalculateTax_ErrorCalculateTaxService(t *testing.T) {
	// req
	reqBody := services.TaxCalculationRequest{}
	reqBodyJson, err := json.Marshal(reqBody)
	if err != nil {
		t.Fatalf("Marshalling request body failed: %v", err)
	}

	// mock
	mockService := &MockTaxService{
		res: services.TaxCalculationResponse{},
		err: fmt.Errorf("mock error"),
	}

	rec, err := callTaxHandler(bytes.NewReader(reqBodyJson), mockService)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusInternalServerError, rec.Code)
}
