package schema

type ErrorResponse struct {
	ErrorCode    int64  `json:"error_code"`
	ErrorMessage string `json:"error_message"`
}

var ErrNotSupportChain = ErrorResponse{
	ErrorCode:    400001,
	ErrorMessage: "the chain is not support yet",
}

var ErrInternal = ErrorResponse{
	ErrorCode:    500001,
	ErrorMessage: "internal server error",
}
