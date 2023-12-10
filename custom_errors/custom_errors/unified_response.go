package customerrors

// UnifiedResponse represents the standard format for all responses
type UnifiedResponse struct {
	Status   string      `json:"status"`
	Data     interface{} `json:"data,omitempty"`
	Error    *ErrorInfo  `json:"error,omitempty"`
	Metadata interface{} `json:"metadata,omitempty"`
}

func NewSuccessResponse(data interface{}) *UnifiedResponse {
	return &UnifiedResponse{
		Status: "success",
		Data:   data,
	}
}

func NewErrorResponse(err *CustomError) *UnifiedResponse {
	return &UnifiedResponse{
		Status: "error",
		Error: &ErrorInfo{
			Code:        err.Code,
			Description: err.Message,
		},
	}
}

func NewInternalServerException(err error) *UnifiedResponse {
	return &UnifiedResponse{
		Status: "error",
		Error: &ErrorInfo{
			Code:        ErrUnknown,
			Description: err.Error(),
		},
	}
}

// ErrorInfo represents detailed error information
type ErrorInfo struct {
	Code        int    `json:"code"`
	Description string `json:"description"`
}
