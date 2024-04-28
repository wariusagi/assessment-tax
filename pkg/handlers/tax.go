package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/wariusagi/assessment-tax/pkg/services"
)

type taxHandler struct {
	srv services.TaxService
}

type Err struct {
	Message string `json:"message"`
}

func NewTaxHandler(srv services.TaxService) taxHandler {
	return taxHandler{srv: srv}
}

func (h taxHandler) CalculateTax(c echo.Context) error {
	req := services.TaxCalculationRequest{}
	err := c.Bind(&req)
	if err != nil {
		return c.JSON(http.StatusBadRequest, Err{Message: err.Error()})
	}

	res, err := h.srv.NewTaxCalculation(req)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, Err{Message: err.Error()})
	}
	return c.JSON(http.StatusOK, res)
}
