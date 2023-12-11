package customerrors

const (
	// ErrUnknown is an unknown error.
	// This should be used only as a fallback error when something unexpected happens.
	// For example, if the error is not handled in the application, like 500 Internal Server Error.
	ErrUnknown = iota + 1
	// ErrWrongInput is returned when wrong input is provided.
	// For example, if body of the request is not valid JSON or required fields are missing.
	ErrWrongInput
	// ErrWrongCredentials is returned when the user provides wrong credentials.
	// For example, if the user provides wrong email or password during login
	ErrWrongCredentials
)

// CustomError represents a custom error.
// It can be used to provide additional information about the error.
// The main objective of this struct is to provide a unified error response
// for all the errors that can occur in the application.
//
// This enables client applications to handle errors in a unified way
// and display a user-friendly error message to the user.
type CustomError struct {
	// Code is a machine-readable error code.
	// It can be used by client applications to handle errors in a unified way.
	// For example, if the error code is ErrWrongCredentials,
	// then the client application can display a message like
	// "Wrong credentials, either email or password is wrong".
	Code int
	// StatusCode is the HTTP status code that should be returned to the client.
	StatusCode int
	// Message is a human-readable description of the error.
	// Included mainly for debugging purposes and a clarification of the error.
	Message string
}

// Error returns the error message.
// It is needed to implement the error interface.
func (e *CustomError) Error() string {
	return e.Message
}

// NewWrongCredentialsError returns a new error with the given message.
// This error should be used when the user provides wrong credentials.
func NewWrongCredentialsError() error {
	return &CustomError{
		Code:       ErrWrongCredentials,
		Message:    "Wrong credentials, either email or password is wrong",
		StatusCode: 401,
	}
}

// NewWrongInputError returns a new error with the given message.
// This error should be used when wrong input is provided.
func NewWrongInputError(message string) error {
	return &CustomError{
		Code:       ErrWrongInput,
		Message:    message,
		StatusCode: 400,
	}
}
