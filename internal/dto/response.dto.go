package dto

type Response struct {
	Success bool `example:"true"`
	Data    any
	Message string `example:"success"`
}

type ErrorResponse struct {
	Success bool   `example:"false"`
	Message string `example:"error"`
}
