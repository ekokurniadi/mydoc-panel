package helper

import "github.com/go-playground/validator/v10"

type Response struct {
	Meta Meta        `json:"meta"`
	Data interface{} `json:"data"`
}
type Meta struct {
	Message string `json:"message"`
	Code    int    `json:"code"`
	Status  string `json:"status"`
}

type ServerSideResponse struct {
	TotalItems   int         `json:"total_items"`
	Data         interface{} `json:"data"`
	TotalPages   int         `json:"total_pages"`
	CurrentPages int         `json:"current_pages"`
}

func ServerSideResponses(totalItems int, totalPages int, currentPages int, data interface{}) ServerSideResponse {
	serverSideResponse := ServerSideResponse{
		TotalItems:   totalItems,
		TotalPages:   totalPages,
		CurrentPages: currentPages,
		Data:         data,
	}
	jsonResponse := serverSideResponse
	return jsonResponse
}

func ApiResponse(message string, code int, status string, data interface{}) Response {
	meta := Meta{
		Message: message,
		Code:    code,
		Status:  status,
	}
	jsonResponse := Response{
		Meta: meta,
		Data: data,
	}
	return jsonResponse
}
func FormatValidationError(err error) []string {
	var errors []string
	for _, e := range err.(validator.ValidationErrors) {
		errors = append(errors, e.Error())
	}
	return errors
}
