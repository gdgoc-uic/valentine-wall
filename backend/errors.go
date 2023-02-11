package main

import (
	"net/http"

	"github.com/pocketbase/pocketbase/apis"
)

type ResponseError struct {
	StatusCode int    `json:"-"`
	WError     error  `json:"-"`
	Message    string `json:"error_message"`
}

func (re *ResponseError) Error() string {
	if re.WError != nil {
		return re.WError.Error()
	} else if len(re.Message) == 0 {
		return http.StatusText(re.StatusCode)
	} else {
		return re.Message
	}
}

func (re *ResponseError) ToApiError() error {
	return apis.NewApiError(re.StatusCode, re.Message, re.WError)
}
