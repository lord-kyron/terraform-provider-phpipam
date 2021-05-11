package request

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"regexp"
	"strings"
	"testing"

	"github.com/pavel-z1/phpipam-sdk-go/phpipam"
	"github.com/pavel-z1/phpipam-sdk-go/phpipam/session"
)

const errorResponseText = `
{
  "code": 500,
  "success": false,
  "message": "Invalid username or password"
}
`

const errorResponseNonJSONText = "<html><head><title>Service Unavailable</title></head><body><b>Service Unavailable</b></body></html>"

const okResponseText = `
{
  "code": 200,
  "success": true,
  "data": {
    "token": "foobarbazboop",
    "expires": "2017-03-03 00:56:34"
  }
}
`

type okAuthResponseData struct {
	Expires string
	Token   string
}

func okResponse() okAuthResponseData {
	return okAuthResponseData{
		Expires: "2017-03-03 00:56:34",
		Token:   "foobarbazboop",
	}
}

const errorResponse = "Error from API (500): Invalid username or password"

func errorResponseNonJSON() string {
	return fmt.Sprintf("Non-API error (503 Service Unavailable): %s", errorResponseNonJSONText)
}

func newHTTPTestServer(f func(w http.ResponseWriter, r *http.Request)) *httptest.Server {
	ts := httptest.NewServer(http.HandlerFunc(f))
	return ts
}

func httpErrorTestServer() *httptest.Server {
	return newHTTPTestServer(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/json")
		http.Error(w, errorResponseText, http.StatusInternalServerError)
	})
}

func httpNonJSONErrorTestServer() *httptest.Server {
	return newHTTPTestServer(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "text/html")
		http.Error(w, errorResponseNonJSONText, http.StatusServiceUnavailable)
	})
}

func httpOKTestServer() *httptest.Server {
	return newHTTPTestServer(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/json")
		http.Error(w, okResponseText, http.StatusOK)
	})
}

func phpipamConfig() phpipam.Config {
	return phpipam.Config{
		AppID:    "0123456789abcdefgh",
		Password: "changeit",
		Username: "nobody",
	}
}

func testRequest(c phpipam.Config, in interface{}, out interface{}) *Request {
	s := &session.Session{
		Config: c,
	}
	r := NewRequest(s)
	r.Method = "GET"
	r.URI = "/api/test/users/"
	r.Input = in
	r.Output = out
	return r
}

func TestRequestSendSuccess(t *testing.T) {
	ts := httpOKTestServer()
	defer ts.Close()
	cfg := phpipamConfig()
	cfg.Endpoint = ts.URL
	in := struct{}{}
	out := okAuthResponseData{}
	r := testRequest(cfg, &in, &out)
	err := r.Send()

	if err != nil {
		t.Fatalf("Unexpected request error: %s", err)
	}

	expected := okResponse()

	if reflect.DeepEqual(expected, out) == false {
		t.Fatalf("expected %v, got %v", expected, out)
	}
}

func TestRequestSendError(t *testing.T) {
	ts := httpErrorTestServer()
	defer ts.Close()
	cfg := phpipamConfig()
	cfg.Endpoint = ts.URL
	in := struct{}{}
	out := okAuthResponseData{}
	r := testRequest(cfg, &in, &out)
	err := r.Send()

	if err == nil {
		t.Fatalf("Expected error, got success")
	}

	expected := errorResponse

	if err.Error() != expected {
		t.Fatalf("expected %s, got %s", expected, err)
	}
}

func TestRequestSendNonJSONError(t *testing.T) {
	ts := httpNonJSONErrorTestServer()
	defer ts.Close()
	cfg := phpipamConfig()
	cfg.Endpoint = ts.URL
	in := struct{}{}
	out := okAuthResponseData{}
	r := testRequest(cfg, &in, &out)
	err := r.Send()

	if err == nil {
		t.Fatalf("Expected error, got success")
	}

	expected := errorResponseNonJSON()

	// HTTP server gives a bunch of whitespace after for some reason
	if strings.TrimSpace(err.Error()) != expected {
		t.Fatalf("expected %s (%T), got %s (%T)", expected, expected, err.Error(), err.Error())
	}
}

func TestRequestSendProtocolError(t *testing.T) {
	ts := httpOKTestServer()
	cfg := phpipamConfig()
	cfg.Endpoint = ts.URL
	in := struct{}{}
	out := okAuthResponseData{}
	r := testRequest(cfg, &in, &out)
	ts.Close()
	err := r.Send()

	if err == nil {
		t.Fatalf("Expected error, got success")
	}

	expected := "^HTTP protocol error"

	if ok, _ := regexp.MatchString(expected, err.Error()); ok == false {
		t.Fatalf("expected error to match %s, got %s", expected, err)
	}
}
