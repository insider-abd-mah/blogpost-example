package rest

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// Client base rest client
type Client interface {
	Get(req RestRequest) RestResponse
	Post(req RestRequest) RestResponse
	Put(req RestRequest) RestResponse
	Delete(req RestRequest) RestResponse
}

// RestClient
type RestClient struct {
	BaseURL        string
	DefaultHeaders map[string]interface{}
}

// Get
func (c RestClient) Get(req RestRequest) RestResponse {
	req.Method = "GET"

	return c.call(req)
}

// Post
func (c RestClient) Post(req RestRequest) RestResponse {
	req.Method = "POST"

	return c.call(req)
}

// call
func (c RestClient) call(r RestRequest) RestResponse {
	res := RestResponse{}
	req, err := http.NewRequest(r.Method, c.getURL(r.Path, r.Query), strings.NewReader(string(r.Body)))

	if err != nil {
		res.Error = err

		return res
	}

	for key, val := range c.DefaultHeaders {
		req.Header.Set(key, fmt.Sprintf("%s", val))
	}

	for key, val := range r.Headers {
		req.Header.Set(key, fmt.Sprintf("%s", val))
	}

	req.Close = true

	cl := http.Client{Timeout: 60 * time.Second}
	resp, err := cl.Do(req)

	if err != nil {
		res.Error = err

		return res
	}

	respBody, err := ioutil.ReadAll(resp.Body)

	defer func() {
		_ = resp.Body.Close()
	}()

	res.StatusCode = StatusCode(resp.StatusCode)
	res.Body = string(respBody)
	res.Error = err

	return res
}

func (c RestClient) getURL(path string, query map[string]interface{}) string {
	params := url.Values{}

	for key, val := range query {
		params[key] = []string{fmt.Sprintf("%s", val)}
	}

	if len(params) > 0 {
		return fmt.Sprintf("%s%s?%s", c.BaseURL, path, params.Encode())
	}

	return fmt.Sprintf("%s%s", c.BaseURL, path)
}
