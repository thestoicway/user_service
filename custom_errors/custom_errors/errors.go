package customerrors

const (
	ErrUnknown = iota + 1
	ErrWrongCredentials
)

type CustomError struct {
	Code       int
	StatusCode int
	Message    string
}

func (e *CustomError) Error() string {
	return e.Message
}

// NewWrongCredentialsError returns a new error with the given message.
func NewWrongCredentialsError() error {
	return &CustomError{
		Code:       ErrWrongCredentials,
		Message:    "Wrong credentials, either email or password is wrong",
		StatusCode: 401,
	}
}
