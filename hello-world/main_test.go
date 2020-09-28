package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/aws/aws-lambda-go/events"
	"github.com/google/go-cmp/cmp"
)

func TestHandler(t *testing.T) {
	t.Run("Unable to get IP", func(t *testing.T) {
		DefaultHTTPGetAddress = "http://127.0.0.1:12345"

		_, err := handler(events.APIGatewayProxyRequest{})
		if err == nil {
			t.Fatal("Error failed to trigger with an invalid request")
		}
	})

	t.Run("Non 200 Response", func(t *testing.T) {
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(500)
		}))
		defer ts.Close()

		DefaultHTTPGetAddress = ts.URL

		_, err := handler(events.APIGatewayProxyRequest{})
		if err != nil && err.Error() != ErrNon200Response.Error() {
			t.Fatalf("Error failed to trigger with an invalid HTTP response: %v", err)
		}
	})

	t.Run("Unable decode IP", func(t *testing.T) {
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(500)
		}))
		defer ts.Close()

		DefaultHTTPGetAddress = ts.URL

		_, err := handler(events.APIGatewayProxyRequest{})
		if err == nil {
			t.Fatal("Error failed to trigger with an invalid HTTP response")
		}
	})

	t.Run("Successful Request", func(t *testing.T) {
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(200)
			fmt.Fprintf(w, "127.0.0.1")
		}))
		defer ts.Close()

		DefaultHTTPGetAddress = ts.URL

		os.Setenv("APP_NAME", "Hello")
		os.Setenv("ENV", "dev")

		want := events.APIGatewayProxyResponse{
			Body:       "Hello@dev",
			StatusCode: 200,
		}

		response, err := handler(events.APIGatewayProxyRequest{})

		if err != nil {
			t.Fatal("Everything should be ok")
		}
		if diff := cmp.Diff(want, response); diff != "" {
			t.Errorf("api response got differs: (-want +got)\n%s", diff)
		}

	})
}
