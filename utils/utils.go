package utils

import (
	"billing_engine/model"
	"errors"
	"net/http"

	"github.com/go-playground/validator"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
)

type APIModel[T interface{}] struct {
	Data T `json:"data"`
}

const (
	USER_ID_HEADER    = "X-User-ID"
	REQUEST_ID_HEADER = "X-Req-ID"
)

const (
	INVALID_USER_ID = "invalid user id"
)

func GetUserID(c echo.Context) (string, error) {
	userID := c.Request().Header.Get(USER_ID_HEADER)
	// Check if userID is present
	if userID == "" {
		return "", errors.New(INVALID_USER_ID)
	}
	return userID, nil
}

func GetRequestID(c echo.Context) string {
	requestID := c.Request().Header.Get(REQUEST_ID_HEADER)
	return requestID
}

func InvalidRequest(err error) error {
	return echo.NewHTTPError(http.StatusBadRequest, err.Error())
}

func InternalServerError(reqID string, err error) error {
	logrus.WithFields(logrus.Fields{
		"requestID": reqID,
		"error":     err.Error(),
	}).Error(model.ISE_MESSAGE)
	return echo.NewHTTPError(http.StatusInternalServerError, model.ISE_MESSAGE)
}

func Response(data interface{}) APIModel[any] {
	return APIModel[any]{Data: data}
}

type CustomValidator struct {
	validator *validator.Validate
}

func (cv *CustomValidator) Validate(i interface{}) error {
	return cv.validator.Struct(i)
}

func AddValidator(e *echo.Echo) {
	e.Validator = &CustomValidator{validator: validator.New()}
}

func BindAndValidateGeneric[T any](c echo.Context) (T, error) {
	var input T
	if err := c.Bind(&input); err != nil {
		return input, err
	}
	if err := c.Validate(&input); err != nil {
		return input, err
	}
	return input, nil
}
