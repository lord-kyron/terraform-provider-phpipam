// Package request provides the HTTP request functionality for the PHPIPAM API.
package request

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/apex/log"
	"github.com/apex/log/handlers/logfmt"
	"github.com/pavel-z1/phpipam-sdk-go/phpipam/session"
)

// APIResponse represents a PHPIPAM response body. Both successful and
// unsuccessful requests share the same response format.
type APIResponse struct {
	// The HTTP result code.
	Code int

	// The response data. This is further unmarshaled into the data type set by
	// Request.Output.
	Data json.RawMessage

	// The error message, if the request failed.
	Message string

	// Whether or not the API request was successful.
	Success bool
}

// Request represents the API request.
type Request struct {
	// The API session.
	Session *session.Session

	// The request method.
	Method string

	// The request URI.
	URI string

	// The request data.
	Input interface{}

	// The output of the request. This corresponds to the "data" field in a
	// response.
	Output interface{}
}

// requestResponse is an unexported struct that encompasses status codes
// and request body in a fashion that can be read after the request
// is closed.
type requestResponse struct {
	// Status code.
	StatusCode int

	// Status code with short-form message.
	Status string

	// Response body.
	Body []byte
}

// BodyString converts requestResponse.Body to string.
func (r *requestResponse) BodyString() string {
	buf := bytes.NewBuffer(r.Body)
	return buf.String()
}

// readResponseJSON reads a "successful" response body as JSON into variable
// pointed to by v.
//
// First the main HTTP response is unmarshalled. If the request at that point
// failed according to the success field, the request is handed off to
// handleError and the resulting error message is returned. Otherwise, the
// request is successful and the response data is unmarshalled.
func (r *requestResponse) ReadResponseJSON(v interface{}) error {
	var resp APIResponse
	if err := json.Unmarshal(r.Body, &resp); err != nil {
		return fmt.Errorf("JSON parsing error: %s - Response body: %s", err, r.Body)
	}

	if !resp.Success {
		return r.handleError()
	}

	if string(resp.Data) != "" {
		if err := json.Unmarshal(resp.Data, v); err != nil {
			return fmt.Errorf("JSON parsing error: %s - Response data: %s", err, string(resp.Data))
		}
	}
	return nil
}

// handleError handles a PHPIPAM API error response.
func (r *requestResponse) handleError() error {
	var resp APIResponse
	if err := json.Unmarshal(r.Body, &resp); err != nil {
		// more than likely not JSON, just pull together the body and return it as
		// the error message
		return fmt.Errorf("Non-API error (%s): %s", r.Status, r.BodyString())
	}

	// Return a properly formatted error from the appropraite fields.
	return fmt.Errorf("Error from API (%d): %s", resp.Code, resp.Message)
}

// newRequestResponse creates a new requestResponse instance off a HTTP
// response. Warning: This also closes the Body.
func newRequestResponse(r *http.Response) *requestResponse {
	rr := &requestResponse{
		StatusCode: r.StatusCode,
		Status:     r.Status,
	}
	defer r.Body.Close()
	body, err := ioutil.ReadAll(r.Body)
	log.Debugf("Response Body Debug ................... %s", body)
	if err != nil {
		panic(err)
	}
	rr.Body = body
	return rr
}

// Send sends a request to the API endpoint, and parsees the response.
//
// Note that by design, Send does not handle redirects - if you get a 302 error
// or some other sort of 300 error from the SDK, please check your API
// endpoints.
func (r *Request) Send() error {
	var req *http.Request
	var err error
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: r.Session.Config.Insecure},
		Proxy:           http.ProxyFromEnvironment,
	}
	client := &http.Client{
		Transport: tr,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}

	switch r.Method {
	case "OPTIONS", "GET", "POST", "PUT", "PATCH", "DELETE":
		bs, err := json.Marshal(r.Input)
		log.Debugf("Request Body Debug ................... %s", bs)
		if err != nil {
			return fmt.Errorf("Error preparing request data: %s", err)
		}
		buf := bytes.NewBuffer(bs)
		log.Debugf("Request URL Debug ...................Method: %s, UR: %s/%s%s", r.Method, r.Session.Config.Endpoint, r.Session.Config.AppID, r.URI)
		req, err = http.NewRequest(r.Method, fmt.Sprintf("%s/%s%s", r.Session.Config.Endpoint, r.Session.Config.AppID, r.URI), buf)
		req.Header.Add("Content-Type", "application/json")
	default:
		return fmt.Errorf("API request method %s not supported by PHPIPAM", r.Method)
	}

	if err != nil {
		panic(err)
	}

	// Add session token if it exists, otherwise append username/password from the config.
	// Note that according to the PHPIPAM docs, Basic Auth does not work on
	// anything else other than the user controller. Falling back to basic auth
	// should only be used for setting up the session only.
	if r.Session.Token.String != "" {
		req.Header.Add("phpipam-token", r.Session.Token.String)
	} else {
		req.SetBasicAuth(r.Session.Config.Username, r.Session.Config.Password)
	}

	re, err := client.Do(req)

	if err != nil {
		return fmt.Errorf("HTTP protocol error: %s", err)
	}

	resp := newRequestResponse(re)

	// A response code of 300 or higher is an error. We do not handle redirects.
	if resp.StatusCode >= 300 {
		return resp.handleError()
	}

	// Unmarshal response into Output. The service is responsible for
	// this being functional past JSON parsing.
	if err := resp.ReadResponseJSON(r.Output); err != nil {
		return err
	}

	return nil
}

// NewRequest creates a new request instance with configuration set.
func NewRequest(s *session.Session) *Request {
	log.SetLevel(log.InfoLevel)
	log.SetHandler(logfmt.New(os.Stderr))

	env_loglevel := os.Getenv("PHPIPAMSDK_LOGLEVEL")
	if env_loglevel != "" {
		loglevel, err := log.ParseLevel(env_loglevel)
		if err == nil {
			log.SetLevel(loglevel)
		} else {
			log.Warnf("Invalid log level, defaulting to info: %s", err)
		}
	}

	r := &Request{
		Session: s,
	}
	return r
}

// change logger level, default is info
func SetLevel(level log.Level) {
	log.SetLevel(level)
}
