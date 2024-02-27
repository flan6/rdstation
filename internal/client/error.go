package client

import "fmt"

type RDError struct {
	Errors Errors `json:"errors"`
}

type Errors struct {
	StatusCode int
	Type       string `json:"error_type"`
	Message    string `json:"error_message"`
}

func (e Errors) Error() string {
	return fmt.Sprintf("%v: %s - %s", e.StatusCode, e.Type, e.Message)
}

func (e RDError) Error() string {
	return fmt.Sprintf("%s", e.Errors)
}
