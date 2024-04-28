package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/wariusagi/assessment-tax/pkg/services"
)

type adminHandler struct {
	srv services.AdminService
}

func NewAdminHandler(srv services.AdminService) adminHandler {
	return adminHandler{srv: srv}
}

func (h adminHandler) SetDeduction(c echo.Context) error {
	req := services.AdminDeductionRequest{}
	err := c.Bind(&req)
	if err != nil {
		return c.JSON(http.StatusBadRequest, Err{Message: err.Error()})
	}

	res, err := h.srv.SetDeduction(req)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, Err{Message: err.Error()})
	}
	return c.JSON(http.StatusOK, res)
}

func (h adminHandler) SetKReceipt(c echo.Context) error {
	req := services.AdminKReceiptRequest{}
	err := c.Bind(&req)
	if err != nil {
		return c.JSON(http.StatusBadRequest, Err{Message: err.Error()})
	}

	res, err := h.srv.SetKReceipt(req)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, Err{Message: err.Error()})
	}
	return c.JSON(http.StatusOK, res)
}
