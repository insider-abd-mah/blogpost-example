package rest

// RestResponse
type RestResponse struct {
	Body       string
	StatusCode StatusCode
	Error      error
}

type StatusCode int

// IsSuccess
func (sc StatusCode) IsSuccess() bool {
	return sc >= 200 && sc < 300
}
