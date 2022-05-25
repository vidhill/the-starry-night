package customerrors

import "fmt"

type UnexpectedResponseError struct {
	statusCode int
}

func (e UnexpectedResponseError) Error() string {
	return fmt.Sprintf("non success response code received from api, received %v status", e.statusCode)
}

func NewUnexpectedResponseError(i int) UnexpectedResponseError {
	return UnexpectedResponseError{
		statusCode: i,
	}
}
