package handler

// JSONResponse defines the standardized API response structure
type JSONResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
	Error   interface{} `json:"error,omitempty"`
}

// SuccessResponse creates a successful JSON response
func SuccessResponse(message string, data interface{}) JSONResponse {
	return JSONResponse{
		Success: true,
		Message: message,
		Data:    data,
	}
}

// ErrorResponse creates an error JSON response
func ErrorResponse(message string, err interface{}) JSONResponse {
	return JSONResponse{
		Success: false,
		Message: message,
		Error:   err,
	}
}
