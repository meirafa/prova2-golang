package web

import "github.com/gin-gonic/gin"

type errorResponse struct {
	StatusCode int    `json:"status_code"`
	Status     string `json:"status"`
	Message    string `json:"message"`
}

type response struct {
	Data interface{} `json:"data"`
}

// ResponseOK escreve uma mensagem de êxito
func ResponseOK(ctx *gin.Context, statusCode int, data interface{}) {
	ctx.JSON(statusCode, response{data})
}

// BadResponse escreve uma mensagem indicando que a operação não foi bem sucedida
func BadResponse(ctx *gin.Context, statusCode int, status, message string) {
	ctx.JSON(status, errorResponse{
		StatusCode: statusCode,
		Status:     status,
		Message:    message,
	})
}

func DeleteResponse(ctx *gin.Context, statusCode int, message string) {
	ctx.JSON(statusCode, errorResponse{
		StatusCode: statusCode,
		Status:     "success",
		Message:    message,
	})
}
