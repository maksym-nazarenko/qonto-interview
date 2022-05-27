package api

type Error string

func (e Error) Error() string {
	return string(e)
}

const (
	ErrMalformedInput = Error("malformed input data")
)

type errorResponse struct {
	Error string `json:"error,omitempty"`
}

func wrapError(err error) *errorResponse {
	return &errorResponse{
		Error: err.Error(),
	}
}
