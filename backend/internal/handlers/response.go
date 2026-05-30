package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Response[T any] struct {
	Data    *T         `json:"data,omitempty"`
	Error   *ErrorBody `json:"error,omitempty"`
	Message string     `json:"message,omitempty"`
}

type ErrorBody struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Details any    `json:"details,omitempty"`
}

func respondOK(c *gin.Context, data any, message string) {
	var responseData *any
	if data != nil {
		responseData = &data
	}

	c.JSON(http.StatusOK, Response[any]{
		Data:    responseData,
		Message: message,
	})
}

func respondError(c *gin.Context, status int, code, message string) {
	c.JSON(status, Response[any]{
		Error: &ErrorBody{
			Code:    code,
			Message: message,
		},
	})
}
