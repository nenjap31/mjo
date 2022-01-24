package merchant

import (
	"database/sql"
	"mjo/controller/merchant/defined"
	"mjo/controller/util/queryparams"
	"mjo/controller/util/response"
	"mjo/service/merchant"
	serviceDefined "mjo/service/merchant/defined"
	"net/http"

	"github.com/labstack/echo/v4"
)

type Controller struct {
	service merchant.IService
}

func NewController(service merchant.IService) *Controller {
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
	data := serviceDefined.Merchant{
		UserId:      bodyRequest.UserId,
		MerchantName:       sql.NullString{String: bodyRequest.MerchantName, Valid:true},
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
	data := serviceDefined.Merchant{
		UserId:      bodyRequest.UserId,
		MerchantName:       sql.NullString{String: bodyRequest.MerchantName, Valid:true},
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