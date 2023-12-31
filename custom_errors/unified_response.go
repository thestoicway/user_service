package customerrors

// UnifiedResponse represents the standard response format for all responses
type UnifiedResponse struct {
	// Status is either "success" or "error"
	Status string `json:"status"`
	// Data is the payload of the response
	Data interface{} `json:"data,omitempty"`
	// Metadata can be used to provide additional information about the response
	// For example, if the response is a list of items, then metadata can contain
	// the total number of items
	Metadata interface{} `json:"metadata,omitempty"`
	// Error is present only if status is "error"
	Error *ErrorInfo `json:"error,omitempty"`
}

// ErrorInfo represents detailed error information
type ErrorInfo struct {
	// Code is a machine-readable error code
	Code int `json:"code"`
	// Description is a human-readable description of the error
	Description string `json:"description"`
	// Details can be used to provide additional information about the error
	// For example, if user is banned from OTP, then details can contain
	// the date when the ban will be lifted
	Details interface{} `json:"details,omitempty"`
}

type newResponseParams struct {
	data     interface{}
	err      *ErrorInfo
	metadata interface{}
}

func newResponse(p newResponseParams) *UnifiedResponse {
	if p.err != nil {
		return &UnifiedResponse{
			Status: "error",
			Error:  p.err,
		}
	}

	return &UnifiedResponse{
		Status:   "success",
		Data:     p.data,
		Metadata: p.metadata,
	}
}

// NewSuccessResponse creates a new success response
func NewSuccessResponse(data interface{}) *UnifiedResponse {
	return newResponse(newResponseParams{
		data: data,
	})
}

// NewErrorResponse creates a new error response
func NewErrorResponse(err *CustomError) *UnifiedResponse {
	return newResponse(newResponseParams{
		err: &ErrorInfo{
			Code:        err.Code,
			Description: err.Message,
			Details:     err.Details,
		},
	})
}
