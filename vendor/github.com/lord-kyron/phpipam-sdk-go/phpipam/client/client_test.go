package client

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/pavel-z1/phpipam-sdk-go/phpipam"
	"github.com/pavel-z1/phpipam-sdk-go/phpipam/session"
)

const testDateStamp = "2999-12-31 23:59:59"

const authErrorResponseText = `
{
  "code": 500,
  "success": false,
  "message": "Invalid username or password"
}
`

const sessionErrorResponseText = `
{
  "code": 403,
  "success": false,
  "message": "Invalid token"
}
`

var authOKResponseText = fmt.Sprintf(`
{
  "code": 200,
  "success": true,
  "data": {
    "token": "foobarbazboop",
    "expires": "%s"
  }
}
`, testDateStamp)

const subnetSearchOKResponseText = `
{
  "code": 200,
  "success": true,
  "data": [
    {
      "id": "3",
      "subnet": "10.10.1.0",
      "mask": "24",
      "sectionId": "1",
      "description": "Customer 1",
      "firewallAddressObject": null,
      "vrfId": "0",
      "masterSubnetId": "2",
      "allowRequests": "1",
      "vlanId": "0",
      "showName": "1",
      "device": "0",
      "permissions": "{\"3\":\"1\",\"2\":\"2\"}",
      "pingSubnet": "0",
      "discoverSubnet": "0",
      "DNSrecursive": "0",
      "DNSrecords": "0",
      "nameserverId": "0",
      "scanAgent": null,
      "isFolder": "0",
      "isFull": "0",
      "tag": "2",
      "editDate": null,
			"Projects": "bazboop",
      "links": [
        {
          "rel": "self",
          "href": "/api/test/subnets/3/"
        }
      ]
    }
  ]
}
`

const subnetGetOKResponseText = `
{
  "code": 200,
  "success": true,
  "data": {
    "id": "3",
    "subnet": "10.10.1.0",
    "mask": "24",
    "sectionId": "1",
    "description": "Customer 1",
    "firewallAddressObject": null,
    "vrfId": "0",
    "masterSubnetId": "2",
    "allowRequests": "1",
    "vlanId": "0",
    "showName": "1",
    "device": "0",
    "permissions": "{\"3\":\"1\",\"2\":\"2\"}",
    "pingSubnet": "0",
    "discoverSubnet": "0",
    "DNSrecursive": "0",
    "DNSrecords": "0",
    "nameserverId": "0",
    "scanAgent": null,
    "isFolder": "0",
    "isFull": "0",
    "tag": "2",
    "editDate": null,
    "Projects": "bazboop",
    "links": [
      {
        "rel": "self",
        "href": "/api/test/subnets/3/"
      }
    ]
  }
}
`

// testSubnetData represents a subnet object. This may match what ends up in
// the subnets controller. Some fields that are missing from the API
// documentation, or are ambiguous, are omitted.
type testSubnetData struct {
	ID             int `json:",string"`
	Subnet         string
	Mask           string
	SectionID      int `json:",string"`
	Description    string
	VrfID          int `json:",string"`
	MasterSubnetID int `json:",string"`
	AllowRequests  int `json:",string"`
	VlanID         int `json:",string"`
	ShowName       int `json:",string"`
	Device         int `json:",string"`
	Permissions    string
	PingSubnet     int `json:",string"`
	DiscoverSubnet int `json:",string"`
	DNSRecursive   int `json:",string"`
	DNSRecords     int `json:",string"`
	NameserverID   int `json:",string"`
	IsFolder       int `json:",string"`
	IsFull         int `json:",string"`
	EditDate       string
	TagID          int `json:"tag,string"`
}

type testSubnetDataResponse struct {
	Data []testSubnetData
}

const subnetSearchErrorResponseText = `
{
  "code": 404,
  "success": false,
  "message": "No subnets found"
}
`

const testCustomFieldsSchemaResponseText = `
{
  "code": 200,
  "success": true,
  "data": {
    "Projects": {
      "name": "Projects",
      "type": "varchar(255)",
      "Comment": "Projects assigned to subnet",
      "Null": "NO",
      "Default": "foobar"
    }
  }
}
`

var testCustomFieldsSchemaExpected = map[string]phpipam.CustomField{
	"Projects": phpipam.CustomField{
		Name:    "Projects",
		Type:    "varchar(255)",
		Comment: "Projects assigned to subnet",
		Null:    "NO",
		Default: "foobar",
	},
}

var testGetCustomFieldsRequestExpected = map[string]interface{}{
	"Projects": "bazboop",
}

const testUpdateCustomFieldsRequestResponseText = `
{
  "code": 200,
  "success": true,
  "data": "subnet updated"
}
`

const testUpdateCustomFieldsRequestExpected = "subnet updated"

const authErrorExpectedResponse = "Error from API (500): Invalid username or password"
const sessionErrorExpectedResponse = "Error from API (403): Invalid token"
const subnetsErrorExpectedResponse = "Error from API (404): No subnets found"
const updateCustomFieldsErrorExpectedResponse = "Custom field Description not found in schema for controller subnets"

func newHTTPTestServer(f func(w http.ResponseWriter, r *http.Request)) *httptest.Server {
	ts := httptest.NewServer(http.HandlerFunc(f))
	return ts
}

func httpAuthErrorTestServer() *httptest.Server {
	return newHTTPTestServer(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/json")
		http.Error(w, authErrorResponseText, http.StatusInternalServerError)
	})
}

func httpSessionErrorTestServer() *httptest.Server {
	return newHTTPTestServer(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/json")
		http.Error(w, sessionErrorResponseText, http.StatusForbidden)
	})
}

func httpSubnetSearchErrorTestServer() *httptest.Server {
	return newHTTPTestServer(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/json")
		http.Error(w, subnetSearchErrorResponseText, http.StatusNotFound)
	})
}

func httpAuthOKTestServer() *httptest.Server {
	return newHTTPTestServer(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/json")
		http.Error(w, authOKResponseText, http.StatusOK)
	})
}

func httpSubnetSearchOKTestServer() *httptest.Server {
	return newHTTPTestServer(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/json")
		http.Error(w, subnetSearchOKResponseText, http.StatusOK)
	})
}

func httpSubnetGetOKTestServer() *httptest.Server {
	return newHTTPTestServer(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/json")
		http.Error(w, subnetGetOKResponseText, http.StatusOK)
	})
}

func httpCustomFieldsSchemaTestServer() *httptest.Server {
	return newHTTPTestServer(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/json")
		http.Error(w, testCustomFieldsSchemaResponseText, http.StatusOK)
	})
}

func httpUpdateCustomFieldsRequestTestServer() *httptest.Server {
	return newHTTPTestServer(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/json")
		http.Error(w, testUpdateCustomFieldsRequestResponseText, http.StatusOK)
	})
}

func phpipamConfig() phpipam.Config {
	return phpipam.Config{
		AppID:    "0123456789abcdefgh",
		Password: "changeit",
		Username: "nobody",
	}
}

func fullSessionConfig() *session.Session {
	return &session.Session{
		Config: phpipamConfig(),
		Token: session.Token{
			String: "foobarbazboop",
		},
	}
}

func TestNewClient(t *testing.T) {
	sess := session.NewSession(phpipamConfig())

	expected := &Client{
		Session: sess,
	}

	actual := NewClient(sess)

	if !reflect.DeepEqual(expected, actual) {
		t.Fatalf("Expected session to be %#v, got %#v", expected, actual)
	}
}

func TestLoginSessionSuccess(t *testing.T) {
	ts := httpAuthOKTestServer()
	defer ts.Close()
	cfg := phpipamConfig()
	cfg.Endpoint = ts.URL
	sess := session.NewSession(cfg)
	client := NewClient(sess)
	if err := loginSession(client.Session); err != nil {
		t.Fatalf("Unexpected error: %#v", err)
	}

	expected := session.Token{
		String: "foobarbazboop",
	}
	actual := client.Session.Token

	if !reflect.DeepEqual(expected, actual) {
		t.Fatalf("Expected session token to be %#v, got %#v", expected, actual)
	}
}

func TestLoginSessionError(t *testing.T) {
	ts := httpAuthErrorTestServer()
	defer ts.Close()
	cfg := phpipamConfig()
	cfg.Endpoint = ts.URL
	sess := session.NewSession(cfg)
	client := NewClient(sess)
	err := loginSession(client.Session)

	if err == nil {
		t.Fatalf("Expected error, got none")
	}

	expected := authErrorExpectedResponse
	actual := err.Error()

	if expected != actual {
		t.Fatalf("Expected error to be %s, got %s", expected, actual)
	}
}

func TestSendRequestSuccess(t *testing.T) {
	ts := httpSubnetSearchOKTestServer()
	defer ts.Close()
	sess := fullSessionConfig()
	sess.Config.Endpoint = ts.URL
	client := NewClient(sess)

	var parsed testSubnetDataResponse
	actual := make([]testSubnetData, 0)
	if err := json.Unmarshal([]byte(subnetSearchOKResponseText), &parsed); err != nil {
		t.Fatalf("Bad: %#v", err)
	}
	expected := parsed.Data

	if err := client.SendRequest("GET", "/subnets/cidr/10.10.1.0/24/", struct{}{}, &actual); err != nil {
		t.Fatalf("Unexpected error: %#v", err)
	}

	if !reflect.DeepEqual(expected, actual) {
		t.Fatalf("Expected output to be token to be %#v, got %#v", expected, actual)
	}
}

func TestSendRequestError(t *testing.T) {
	ts := httpSubnetSearchErrorTestServer()
	defer ts.Close()
	sess := fullSessionConfig()
	sess.Config.Endpoint = ts.URL
	client := NewClient(sess)

	tmp := make([]testSubnetData, 0)
	err := client.SendRequest("GET", "/subnets/cidr/10.10.1.0/24/", struct{}{}, &tmp)

	if err == nil {
		t.Fatalf("Expected error, got none")
	}

	expected := subnetsErrorExpectedResponse
	actual := err.Error()

	if expected != actual {
		t.Fatalf("Expected error to be %s, got %s", expected, actual)
	}
}

func TestGetCustomFieldsSchema(t *testing.T) {
	ts := httpCustomFieldsSchemaTestServer()
	defer ts.Close()
	sess := fullSessionConfig()
	sess.Config.Endpoint = ts.URL
	client := NewClient(sess)

	expected := testCustomFieldsSchemaExpected
	actual, err := client.GetCustomFieldsSchema("subnets")
	if err != nil {
		t.Fatalf("Bad: %s", err)
	}

	if !reflect.DeepEqual(expected, actual) {
		t.Fatalf("Expected %#v, got %#v", expected, actual)
	}
}

func TestUpdateCustomFieldsRequest(t *testing.T) {
	ts := httpUpdateCustomFieldsRequestTestServer()
	defer ts.Close()
	sess := fullSessionConfig()
	sess.Config.Endpoint = ts.URL
	client := NewClient(sess)

	in := map[string]interface{}{
		"Projects": "updated",
	}

	expected := testUpdateCustomFieldsRequestExpected
	actual, err := client.updateCustomFieldsRequest(3, in, "subnets", testCustomFieldsSchemaExpected)
	if err != nil {
		t.Fatalf("Bad: %s", err)
	}

	if !reflect.DeepEqual(expected, actual) {
		t.Fatalf("Expected %#v, got %#v", expected, actual)
	}
}

func TestUpdateCustomFieldsRequestIllegalField(t *testing.T) {
	ts := httpUpdateCustomFieldsRequestTestServer()
	defer ts.Close()
	sess := fullSessionConfig()
	sess.Config.Endpoint = ts.URL
	client := NewClient(sess)

	in := map[string]interface{}{
		"Description": "sneaky",
	}

	_, err := client.updateCustomFieldsRequest(3, in, "subnets", testCustomFieldsSchemaExpected)
	if err == nil {
		t.Fatalf("Expected error, got none")
	}

	if err.Error() != updateCustomFieldsErrorExpectedResponse {
		t.Fatalf("Expected %q, got %q", updateCustomFieldsErrorExpectedResponse, err.Error())
	}
}
