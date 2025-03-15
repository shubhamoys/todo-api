package utils

import "github.com/gin-gonic/gin"

// Response Struct
type APIResponse struct {
	Status  bool        `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`  // Omits if nil
	Error   *APIError   `json:"error,omitempty"` // Omits if nil
}

// APIError Struct
type APIError struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Details interface{} `json:"details,omitempty"`
}

// Success Response
func SuccessResponse(c *gin.Context, statusCode int, message string, data interface{}) {
	c.JSON(statusCode, APIResponse{
		Status:  true,
		Message: message,
		Data:    data,
	})
}

// Error Response
func ErrorResponse(c *gin.Context, statusCode int, message string, err error, details interface{}) {
	c.JSON(statusCode, APIResponse{
		Status:  false,
		Message: message,
		Error: &APIError{
			Code:    statusCode,
			Message: err.Error(),
			Details: details,
		},
	})
}
