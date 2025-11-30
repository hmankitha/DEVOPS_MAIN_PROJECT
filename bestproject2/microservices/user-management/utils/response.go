package utils

import "github.com/golang-jwt/jwt/v5"

type Claims struct {
	UserID   string `json:"user_id"`
	Email    string `json:"email"`
	Username string `json:"username"`
	Role     string `json:"role"`
	jwt.RegisteredClaims
}

type Response struct {
	Success bool        `json:"success"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}

type PaginatedResponse struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data"`
	Page    int         `json:"page"`
	Limit   int         `json:"limit"`
	Total   int         `json:"total"`
}

func SuccessResponse(data interface{}, message string) *Response {
	return &Response{
		Success: true,
		Message: message,
		Data:    data,
	}
}

func ErrorResponse(err string) *Response {
	return &Response{
		Success: false,
		Error:   err,
	}
}
