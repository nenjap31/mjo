package transaction

import (
	"mjo/controller/transaction/defined"
	"mjo/controller/util/queryparams"
	"mjo/controller/util/response"
	"mjo/service/transaction"
	serviceDefined "mjo/service/transaction/defined"
	"net/http"
	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo/v4"
)

type Controller struct {
	service transaction.IService
}

func NewController(service transaction.IService) *Controller {
	return &Controller{
		service: service,
	}
}

func (controller Controller) List(c echo.Context) error {
	queryParams := c.QueryParams()
	cleanQueryParams := queryparams.QueryParamsCleaner(queryParams)
	result, err := controller.service.List(cleanQueryParams.QueryParams, cleanQueryParams.PerPage, cleanQueryParams.Offset)
	if err != nil {
		result := response.Map["badRequest"]
		result.Errors = append(result.Errors, err.Error())
		return c.JSON(http.StatusBadRequest, response.NewResponse("", result, nil))
	}
	var results []defined.DefaultResponse
	for _, data := range result {
		results = append(results, *defined.NewDefaultResponse(data))
	}
	return c.JSON(http.StatusOK, response.NewResponse("", response.Map["ok"], results))
}

func (controller Controller) Create(c echo.Context) error {
	bodyRequest := new(defined.InsertRequest)
	if err := c.Bind(bodyRequest); err != nil {
		result := response.Map["badRequest"]
		return c.JSON(http.StatusBadRequest, response.NewResponse("", result, nil))
	}
	if err := c.Validate(bodyRequest); err != nil {
		errors := response.BuildErrorBodyRequestValidator(err)
		result := response.Map["badRequest"]
		result.Errors = errors
		return c.JSON(http.StatusBadRequest, response.NewResponse("", result, nil))
	}
	data := serviceDefined.Transaction{
		MerchantId:      bodyRequest.MerchantId,
		OutletId:       bodyRequest.OutletId,
		BillTotal:       bodyRequest.BillTotal,
		CreatedBy:    bodyRequest.CreatedBy,
	}
	result, err := controller.service.Create(data)
	if err != nil {
		result := response.Map["badRequest"]
		result.Errors = append(result.Errors, err.Error())
		return c.JSON(http.StatusBadRequest, response.NewResponse("", result, nil))
	}
	return c.JSON(http.StatusCreated, response.NewResponse("", response.Map["created"], defined.NewDefaultResponse(result)))
}

func (controller Controller) FindById(c echo.Context) error {
	id := c.Param("id")
	result, err := controller.service.FindById(id)
	if err != nil {
		result := response.Map["badRequest"]
		result.Errors = append(result.Errors, err.Error())
		return c.JSON(http.StatusBadRequest, response.NewResponse("", result, nil))
	}
	return c.JSON(http.StatusOK, response.NewResponse("", response.Map["ok"], defined.NewDefaultResponse(result)))
}

func (controller Controller) UpdateById(c echo.Context) error {
	id := c.Param("id")
	bodyRequest := new(defined.UpdateByIdRequest)
	if err := c.Bind(bodyRequest); err != nil {
		result := response.Map["badRequest"]
		return c.JSON(http.StatusBadRequest, response.NewResponse("", result, nil))
	}
	if err := c.Validate(bodyRequest); err != nil {
		errors := response.BuildErrorBodyRequestValidator(err)
		result := response.Map["badRequest"]
		result.Errors = errors
		return c.JSON(http.StatusBadRequest, response.NewResponse("", result, nil))
	}
	data := serviceDefined.Transaction{
		MerchantId:      bodyRequest.MerchantId,
		OutletId:       bodyRequest.OutletId,
		BillTotal:       bodyRequest.BillTotal,
		UpdatedBy:    bodyRequest.UpdatedBy,
	}
	result, err := controller.service.UpdateById(id, data)
	if err != nil {
		result := response.Map["badRequest"]
		result.Errors = []string{err.Error()}
		return c.JSON(http.StatusBadRequest, response.NewResponse("", result, nil))
	}
	return c.JSON(http.StatusOK, response.NewResponse("", response.Map["ok"], defined.NewDefaultResponse(result)))
}

func (controller Controller) DeleteById(c echo.Context) error {
	id := c.Param("id")
	err := controller.service.DeleteById(id)
	if err != nil {
		result := response.Map["badRequest"]
		result.Errors = []string{err.Error()}
		return c.JSON(http.StatusBadRequest, response.NewResponse("", result, nil))
	}
	return c.JSON(http.StatusOK, response.NewResponse("", response.Map["deleted"], nil))
}

func (controller Controller) MonthlyReport(c echo.Context) error {
	token := c.Get("user").(*jwt.Token)
	bodyRequest := new(defined.ReportMerchantRequest)
	if err := c.Bind(bodyRequest); err != nil {
		result := response.Map["badRequest"]
		return c.JSON(http.StatusBadRequest, response.NewResponse("", result, nil))
	}
	if err := c.Validate(bodyRequest); err != nil {
		errors := response.BuildErrorBodyRequestValidator(err)
		result := response.Map["badRequest"]
		result.Errors = errors
		return c.JSON(http.StatusBadRequest, response.NewResponse("", result, nil))
	}
	result, err := controller.service.MonthlyReport(token, bodyRequest.StartDate, bodyRequest.EndDate, bodyRequest.Limit, bodyRequest.Offset)
	if err != nil {
		result := response.Map["badRequest"]
		result.Errors = append(result.Errors, err.Error())
		return c.JSON(http.StatusBadRequest, response.NewResponse("", result, nil))
	}
	var results []defined.ReportMonthlyResponse
	for _, data := range result {
		results = append(results, *defined.NewReportMonthlyResponse(data))
	}
	return c.JSON(http.StatusOK, response.NewResponse("", response.Map["ok"], results))
}

func (controller Controller) MonthlyOutletReport(c echo.Context) error {
	token := c.Get("user").(*jwt.Token)
	bodyRequest := new(defined.ReportMerchantRequest)
	if err := c.Bind(bodyRequest); err != nil {
		result := response.Map["badRequest"]
		return c.JSON(http.StatusBadRequest, response.NewResponse("", result, nil))
	}
	if err := c.Validate(bodyRequest); err != nil {
		errors := response.BuildErrorBodyRequestValidator(err)
		result := response.Map["badRequest"]
		result.Errors = errors
		return c.JSON(http.StatusBadRequest, response.NewResponse("", result, nil))
	}
	result, err := controller.service.MonthlyOutletReport(token, bodyRequest.StartDate, bodyRequest.EndDate, bodyRequest.Limit, bodyRequest.Offset)
	if err != nil {
		result := response.Map["badRequest"]
		result.Errors = append(result.Errors, err.Error())
		return c.JSON(http.StatusBadRequest, response.NewResponse("", result, nil))
	}
	var results []defined.ReportOutletMonthlyResponse
	for _, data := range result {
		results = append(results, *defined.NewReportOutletMonthlyResponse(data))
	}
	return c.JSON(http.StatusOK, response.NewResponse("", response.Map["ok"], results))
}