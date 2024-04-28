package handlers

import (
	"encoding/csv"
	"fmt"
	"io"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/wariusagi/assessment-tax/pkg/services"
	"github.com/wariusagi/assessment-tax/pkg/utils"
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

	res, err := h.srv.CalculateTax(req)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, Err{Message: err.Error()})
	}
	return c.JSON(http.StatusOK, res)
}

func (h taxHandler) CalculateTaxFromCsv(c echo.Context) error {
	file, err := c.FormFile("taxFile")
	if err != nil {
		return c.JSON(http.StatusBadRequest, Err{Message: "get file: " + err.Error()})
	}

	src, err := file.Open()
	if err != nil {
		return c.JSON(http.StatusBadRequest, Err{Message: "open file" + err.Error()})
	}
	defer src.Close()

	reader := csv.NewReader(src)

	// read data csv
	req := []services.TaxCalculationUploadFileHeaderRequest{}
	isCheckHeader := false
	for {
		row, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		if len(row) != 3 {
			return c.JSON(http.StatusBadRequest, Err{Message: "column more than 3"})
		}

		if !isCheckHeader {
			if row[0] != "totalIncome" || row[1] != "wht" || row[2] != "donation" {
				return c.JSON(http.StatusBadRequest, Err{Message: "column name wrong"})
			}
			isCheckHeader = true
		} else {
			row0, err := utils.ParseFloat(row[0])
			if err != nil {
				return err
			}
			row1, err := utils.ParseFloat(row[1])
			if err != nil {
				return err
			}
			row2, err := utils.ParseFloat(row[2])
			if err != nil {
				return err
			}
			req = append(req, services.TaxCalculationUploadFileHeaderRequest{
				TotalIncome: row0,
				Wht:         row1,
				Donation:    row2,
			})
		}
	}
	fmt.Println("request file: ", req)
	// todo: calculate tax service
	return c.JSON(http.StatusOK, services.TaxCalculationUploadResponse{})
}
