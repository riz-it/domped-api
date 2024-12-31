package domain

import (
	"fmt"
	"strings"

	"riz.it/domped/internal/dto"
)

type ValidationError struct {
	FailedField string
	Tag         string
	Param       *string
	Value       *interface{}
}

func NewError(code int, opts ...interface{}) *dto.ApiResponse[any] {
	message := statusMessages[code]
	if len(opts) > 0 {
		if customMessage, ok := opts[0].(string); ok {
			message = customMessage
		}
	}

	errors := extractValidationErrors(opts)

	return &dto.ApiResponse[any]{
		Status:  false,
		Message: message,
		Data:    nil,
		Errors:  errors,
	}
}

func extractValidationErrors(opts []interface{}) map[string][]string {
	if len(opts) > 1 {
		if validationErrors, ok := opts[1].([]ValidationError); ok {
			return mapValidationErrors(validationErrors)
		}
	}
	return nil
}

func mapValidationErrors(validationErrors []ValidationError) map[string][]string {
	errors := make(map[string][]string)
	for _, err := range validationErrors {
		field := formatFieldName(err.FailedField)
		errors[field] = append(errors[field], translateValidationError(err))
	}
	return errors
}

func formatFieldName(field string) string {
	return strings.ReplaceAll(strings.ToLower(field), " ", "_")
}

func translateValidationError(err ValidationError) string {
	messages := map[string]string{
		"required": "%s is required",
		"email":    "%s must be a valid email address",
		"min":      "%s must be at least %s",
		"max":      "%s must be at most %s",
		"eqfield":  "%s must be equal to %s",
		"numeric":  "%s must be a number",
		"number":   "%s must be a number",
		"exists":   "%s does not exist",
		"unique":   "%s has already been taken",
		"oneof":    "%s must be one of %s",
		"url":      "%s must be a valid URL",
	}

	if msg, exists := messages[err.Tag]; exists {
		if strings.Contains(msg, "%s") && err.Param != nil {
			return fmt.Sprintf(msg, err.FailedField)
		}
		return fmt.Sprintf(msg, err.FailedField)
	}
	return err.Tag
}

var statusMessages = map[int]string{
	200: "OK",
	201: "Created",
	202: "Accepted",
	203: "Non-Authoritative Information",
	204: "No Content",
	205: "Reset Content",
	206: "Partial Content",
	207: "Multi-Status",
	208: "Already Reported",
	226: "IM Used",
	400: "Bad Request",
	401: "Unauthorized",
	402: "Payment Required",
	403: "Forbidden",
	404: "Not Found",
	405: "Method Not Allowed",
	406: "Not Acceptable",
	407: "Proxy Authentication Required",
	408: "Request Timeout",
	409: "Conflict",
	410: "Gone",
	411: "Length Required",
	412: "Precondition Failed",
	413: "Request Entity Too Large",
	414: "Request URI Too Long",
	415: "Unsupported Media Type",
	416: "Requested Range Not Satisfiable",
	417: "Expectation Failed",
	418: "I'm a teapot",
	421: "Misdirected Request",
	422: "Unprocessable Entity",
	423: "Locked",
	424: "Failed Dependency",
	425: "Too Early",
	426: "Upgrade Required",
	428: "Precondition Required",
	429: "Too Many Requests",
	431: "Request Header Fields Too Large",
	451: "Unavailable For Legal Reasons",
	500: "Internal Server Error",
	501: "Not Implemented",
	502: "Bad Gateway",
	503: "Service Unavailable",
	504: "Gateway Timeout",
	505: "HTTP Version Not Supported",
	506: "Variant Also Negotiates",
	507: "Insufficient Storage",
	508: "Loop Detected",
	510: "Not Extended",
	511: "Network Authentication Required",
}
