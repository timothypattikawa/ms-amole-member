package exception

import (
	"errors"
	"github.com/labstack/echo/v4"
	"github.com/timothypattikawa/ms-kamoro-costumer/internal/dto"
	"net/http"
)

func sendErrorResponse(c echo.Context, statusCode int, errorMessage string) bool {
	err := c.JSON(statusCode, dto.BaseResponse{
		Error: errorMessage,
	})

	return err == nil
}

func CostumeEchoError(err error, c echo.Context) {
	if notFoundError(err, c) {
		return
	} else if badRequestError(err, c) {
		return
	} else if unauthorizedError(err, c) {
		return
	} else {
		internalError(err, c)
		return
	}
}

func badRequestError(err error, c echo.Context) bool {
	var response *BadReqeustError
	ok := errors.As(err, &response)
	if ok {
		errorResponse := sendErrorResponse(c, http.StatusBadRequest, err.Error())
		return errorResponse
	}
	return false
}

func notFoundError(err error, c echo.Context) bool {
	var response *NotFoundError
	ok := errors.As(err, &response)
	if ok {
		errorResponse := sendErrorResponse(c, http.StatusBadRequest, err.Error())
		return errorResponse
	}
	return false
}

func internalError(err error, c echo.Context) bool {
	var response *InternalServerError
	ok := errors.As(err, &response)
	if ok {
		errorResponse := sendErrorResponse(c, http.StatusInternalServerError, err.Error())
		return errorResponse
	}
	return false
}

func unauthorizedError(err error, c echo.Context) bool {
	var response *Unauthorized
	ok := errors.As(err, &response)
	if ok {
		errorResponse := sendErrorResponse(c, http.StatusBadRequest, err.Error())
		return errorResponse
	}
	return false
}
