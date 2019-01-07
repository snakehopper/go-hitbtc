package spiral

type errorResponse struct {
	ErrorCode int64  `json:"error_code"`
	Message   string `json:"message"`
}
