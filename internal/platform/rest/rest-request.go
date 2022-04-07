package rest

// RestRequest
type RestRequest struct {
	Path    string
	Query   map[string]interface{}
	Body    []byte
	Headers map[string]interface{}
	Method  string
}
