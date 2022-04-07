package main

import (
	"blog-example/internal/platform/database"
	"blog-example/internal/platform/rest"
	"blog-example/test"
	"bytes"
	"fmt"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func runTestServer() *httptest.Server {
	database.Init()

	return httptest.NewServer(setupServer())
}

func Test_post_api_integration_tests_store_endpoint(t *testing.T) {
	ts := runTestServer()
	defer ts.Close()

	t.Run("it should return 200 when health is ok", func(t *testing.T) {
		resp, err := http.Get(fmt.Sprintf("%s/health", ts.URL))

		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}

		assert.Equal(t, 200, resp.StatusCode)
	})

	t.Run("it should return validation error while store when request miss required parameters", func(t *testing.T) {
		_, teardown := test.NewIntegration(t)
		defer teardown()

		resp := (rest.RestClient{}).Post(
			rest.RestRequest{
				Path:   fmt.Sprintf("%s/posts/v1/store", ts.URL),
				Body:   []byte(fmt.Sprintf(`{"description": "%v"}`, "test")),
				Method: "POST",
			},
		)
		respBody := strings.TrimSuffix(fmt.Sprintf("%s", resp.Body), "\n")

		assert.Equal(t, 400, int(resp.StatusCode))
		assert.Equal(t, `{"statusMessage":"validation error"}`, respBody)
	})

	t.Run("it should return ok when insert new post successfully", func(t *testing.T) {
		buf := new(bytes.Buffer)
		_, teardown := test.NewIntegration(t)
		defer teardown()

		req := test.ReadStubFile("posts/store-request.json")
		io := bytes.NewBuffer(req)
		resp, _ := http.Post(fmt.Sprintf("%s/posts/v1/store", ts.URL), "application/json", io)

		_, _ = buf.ReadFrom(resp.Body)

		assert.Equal(t, 200, resp.StatusCode)
		assert.Equal(t, `{"statusMessage":"ok"}`, buf.String())
	})
}

func Test_amplitude_cohorts_api_integration_tests_get_all_endpoint(t *testing.T) {
	ts := runTestServer()
	defer ts.Close()

	t.Run("it should return the list of existed posts", func(t *testing.T) {
		db, teardown := test.NewIntegration(t)
		defer teardown()

		_, _ = db.Exec((`INSERT INTO posts (title, description) VALUES ("test title", "test description")`))
		resp := (rest.RestClient{}).Get(
			rest.RestRequest{
				Path:   fmt.Sprintf("%s/posts/v1/all", ts.URL),
				Method: "GET",
			},
		)

		assert.Equal(t, 200, int(resp.StatusCode))
		assert.Equal(t, `{"data":[{"title":"test title","description":"test description"}],"statusMessage":"ok"}`, resp.Body)
	})
}
