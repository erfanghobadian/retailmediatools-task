package utils

type FieldError struct {
	Field  string `json:"field"`
	Reason string `json:"reason"`
}

type ErrorResponse struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Details interface{} `json:"details,omitempty"`
}
