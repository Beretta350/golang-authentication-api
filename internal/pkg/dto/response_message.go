package dto

import "net/http"

type ResponseMessage struct {
	StatusCode int    `json:"-"`
	Status     string `json:"status"`
	Message    string `json:"message"`
	Data       any    `json:"data,omitempty"`
	Errors     any    `json:"errors,omitempty"`
}

func OkResponse(message string, data any) *ResponseMessage {
	return newResponseMessage(
		http.StatusOK,
		http.StatusText(http.StatusOK),
		message,
		data,
		nil,
	)
}

func CreatedResponse(message string, data any) *ResponseMessage {
	return newResponseMessage(
		http.StatusCreated,
		http.StatusText(http.StatusCreated),
		message,
		data,
		nil,
	)
}

func UnauthorizedResponse(message string) *ResponseMessage {
	return newResponseMessage(
		http.StatusUnauthorized,
		http.StatusText(http.StatusUnauthorized),
		message,
		nil,
		nil,
	)
}

func BadRequestResponse(message string, errs []string) *ResponseMessage {
	return newResponseMessage(
		http.StatusBadRequest,
		http.StatusText(http.StatusBadRequest),
		message,
		nil,
		errs,
	)
}

func InternalErrorResponse(message string, errs any) *ResponseMessage {
	return newResponseMessage(
		http.StatusInternalServerError,
		http.StatusText(http.StatusInternalServerError),
		message,
		nil,
		errs,
	)
}

func newResponseMessage(args ...any) *ResponseMessage {
	return &ResponseMessage{
		StatusCode: args[0].(int),
		Status:     args[1].(string),
		Message:    args[2].(string),
		Data:       args[3],
		Errors:     args[4],
	}
}
