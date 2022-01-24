package response

import (
	"fmt"
	"net/http"
	"regexp"
	"strings"

	"github.com/go-playground/validator"
)

type Result struct {
	Code    int      `json:"code"`
	Message string   `json:"message"`
	Errors  []string `json:"errors"`
}

type Response struct {
	Signature string      `json:"signature"`
	Result    Result      `json:"result"`
	Data      interface{} `json:"data"`
}

func NewResponse(signature string, result Result, data interface{}) Response {
	return Response{
		Signature: signature,
		Result:    result,
		Data:      data,
	}
}

var Map = map[string]Result{
	"ok":                  {Code: 200000, Message: strings.ToLower(http.StatusText(http.StatusOK))},
	"deleted":             {Code: 200001, Message: strings.ToLower(http.StatusText(http.StatusOK))},
	"created":             {Code: 201000, Message: strings.ToLower(http.StatusText(http.StatusCreated))},
	"badRequest":          {Code: 400000, Message: strings.ToLower(http.StatusText(http.StatusBadRequest))},
	"unauthorized":        {Code: 401000, Message: strings.ToLower(http.StatusText(http.StatusBadRequest))},
	"internalServerError": {Code: 500000, Message: strings.ToLower(http.StatusText(http.StatusBadRequest))},
}

func BuildErrorBodyRequestValidator(err error) []string {
	var errors []string
	for _, err := range err.(validator.ValidationErrors) {
		message := ErrorBodyRequestValidatorString(err.Tag(), err.Field(), err.Param())
		errors = append(errors, message)
	}
	return errors
}

func ErrorBodyRequestValidatorString(tag string, field string, param string) string {
	error := ""
	switch tag {
	case "required":
		error = fmt.Sprintf("field %s is required", SnakeToCamel(field))
	}
	return error
}

var matchFirstCap = regexp.MustCompile("(.)([A-Z][a-z]+)")
var matchAllCap = regexp.MustCompile("([a-z0-9])([A-Z])")

func SnakeToCamel(str string) string {
	snake := matchFirstCap.ReplaceAllString(str, "${1}_${2}")
	snake = matchAllCap.ReplaceAllString(snake, "${1}_${2}")
	return strings.ToLower(snake)
}