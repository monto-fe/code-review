package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Response 标准响应结构
type Response struct {
	RetCode int         `json:"ret_code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
	Error   interface{} `json:"error,omitempty"`
}

// Success 成功响应
func Success(c *gin.Context, data interface{}, message string, statusCode int) {
	if message == "" {
		message = "Success"
	}
	c.JSON(http.StatusOK, Response{
		RetCode: statusCode,
		Message: message,
		Data:    data,
	})
}

// Error 错误响应
func Error(c *gin.Context, err error, message string, statusCode int) {
	if message == "" {
		message = "An error occurred"
	}
	if statusCode == 0 {
		statusCode = 10000
	}
	c.JSON(http.StatusOK, Response{
		RetCode: statusCode,
		Message: message,
		Error:   err,
	})
}
