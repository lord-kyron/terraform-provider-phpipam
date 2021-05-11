package addresses

import (
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"reflect"
	"testing"

	"github.com/pavel-z1/phpipam-sdk-go/phpipam"
	"github.com/pavel-z1/phpipam-sdk-go/phpipam/session"
	"github.com/pavel-z1/phpipam-sdk-go/testacc"
)

var testCreateAddressInput = Address{
	SubnetID:    3,
	IPAddress:   "10.10.1.10",
	Description: "foobar",
}

const testCreateAddressOutputExpected = `Address created`
const testCreateAddressOutputJSON = `
{
  "code": 201,
  "success": true,
  "data": "Address created"
}
`

var testGetAddressByIDOutputExpected = Address{
	ID:          11,
	SubnetID:    3,
	IPAddress:   "10.10.1.10",
	Description: "foobar",
	Tag:         2,
}

const testGetAddressByIDOutputJSON = `
{
  "code": 200,
  "success": true,
  "data": {
    "id": "11",
    "subnetId": "3",
    "ip": "10.10.1.10",
    "is_gateway": null,
    "description": "foobar",
    "hostname": null,
    "mac": null,
    "owner": null,
    "tag": "2",
    "deviceId": null,
    "port": null,
    "note": null,
    "lastSeen": null,
    "excludePing": null,
    "PTRignore": null,
    "PTR": "0",
    "firewallAddressObject": null,
    "editDate": null,
    "links": [
      {
        "rel": "self",
        "href": "/api/test/addresses/11/",
        "methods": [
          "GET",
          "POST",
          "DELETE",
          "PATCH"
        ]
      },
      {
        "rel": "ping",
        "href": "/api/test/addresses/11/ping/",
        "methods": [
          "GET"
        ]
      }
    ]
  }
}
`

var testGetAddressesByIPOutputExpected = []Address{
	Address{
		ID:          11,
		SubnetID:    3,
		IPAddress:   "10.10.1.10",
		Description: "foobar",
		Tag:         2,
	},
}

const testGetAddressesByIPOutputJSON = `
{
  "code": 200,
  "success": true,
  "data": [
    {
      "id": "11",
      "subnetId": "3",
      "ip": "10.10.1.10",
      "is_gateway": null,
      "description": "foobar",
      "hostname": null,
      "mac": null,
      "owner": null,
      "tag": "2",
      "deviceId": null,
      "port": null,
      "note": null,
      "lastSeen": null,
      "excludePing": null,
      "PTRignore": null,
      "PTR": "0",
      "firewallAddressObject": null,
      "editDate": null,
      "links": [
        {
          "rel": "self",
          "href": "/api/test/addresses/11/"
        }
      ]
    }
  ]
}
`

var testGetAddressCustomFieldsSchemaExpected = map[string]phpipam.CustomField{
	"CustomTestAddresses": phpipam.CustomField{
		Name:    "CustomTestAddresses",
		Type:    "varchar(255)",
		Comment: "Test field for addresses controller",
		Null:    "YES",
		Default: "",
	},
	"CustomTestAddresses2": phpipam.CustomField{
		Name:    "CustomTestAddresses2",
		Type:    "varchar(255)",
		Comment: "Test field for addresses controller (second field)",
		Null:    "YES",
		Default: "",
	},
}

const testGetAddressCustomFieldsSchemaJSON = `
{
  "code": 200,
  "success": true,
  "data": {
    "CustomTestAddresses": {
      "name": "CustomTestAddresses",
      "type": "varchar(255)",
      "Comment": "Test field for addresses controller",
      "Null": "YES",
      "Default": ""
    },
    "CustomTestAddresses2": {
      "name": "CustomTestAddresses2",
      "type": "varchar(255)",
      "Comment": "Test field for addresses controller (second field)",
      "Null": "YES",
      "Default": ""
    }
  }
}
`

var testUpdateAddressInput = Address{
	ID:          11,
	Description: "bazboop",
}

const testUpdateAddressOutputExpected = `Address updated`
const testUpdateAddressOutputJSON = `
{
  "code": 200,
  "success": true,
  "data": "Address updated"
}
`

const testDeleteAddressOutputExpected = `Address deleted`
const testDeleteAddressOutputJSON = `
{
  "code": 200,
  "success": true,
  "data": "Address deleted"
}
`

func newHTTPTestServer(f func(w http.ResponseWriter, r *http.Request)) *httptest.Server {
	ts := httptest.NewServer(http.HandlerFunc(f))
	return ts
}

func httpOKTestServer(output string) *httptest.Server {
	return newHTTPTestServer(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/json")
		http.Error(w, output, http.StatusOK)
	})
}

func httpCreatedTestServer(output string) *httptest.Server {
	return newHTTPTestServer(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/json")
		http.Error(w, output, http.StatusCreated)
	})
}

func fullSessionConfig() *session.Session {
	return &session.Session{
		Config: phpipam.Config{
			AppID:    "0123456789abcdefgh",
			Password: "changeit",
			Username: "nobody",
		},
		Token: session.Token{
			String: "foobarbazboop",
		},
	}
}

func TestCreateAddress(t *testing.T) {
	ts := httpCreatedTestServer(testCreateAddressOutputJSON)
	defer ts.Close()
	sess := fullSessionConfig()
	sess.Config.Endpoint = ts.URL
	client := NewController(sess)

	in := testCreateAddressInput
	expected := testCreateAddressOutputExpected
	actual, err := client.CreateAddress(in)
	if err != nil {
		t.Fatalf("Bad: %s", err)
	}

	if !reflect.DeepEqual(expected, actual) {
		t.Fatalf("Expected %#v, got %#v", expected, actual)
	}
}

func TestGetAddressByID(t *testing.T) {
	ts := httpOKTestServer(testGetAddressByIDOutputJSON)
	defer ts.Close()
	sess := fullSessionConfig()
	sess.Config.Endpoint = ts.URL
	client := NewController(sess)

	expected := testGetAddressByIDOutputExpected
	actual, err := client.GetAddressByID(11)
	if err != nil {
		t.Fatalf("Bad: %s", err)
	}

	if !reflect.DeepEqual(expected, actual) {
		t.Fatalf("Expected %#v, got %#v", expected, actual)
	}
}

func TestGetAddressesByIP(t *testing.T) {
	ts := httpOKTestServer(testGetAddressesByIPOutputJSON)
	defer ts.Close()
	sess := fullSessionConfig()
	sess.Config.Endpoint = ts.URL
	client := NewController(sess)

	expected := testGetAddressesByIPOutputExpected
	actual, err := client.GetAddressesByIP("10.10.1.10/24")
	if err != nil {
		t.Fatalf("Bad: %s", err)
	}

	if !reflect.DeepEqual(expected, actual) {
		t.Fatalf("Expected %#v, got %#v", expected, actual)
	}
}

func TestGetAddressCustomFieldsSchema(t *testing.T) {
	ts := httpOKTestServer(testGetAddressCustomFieldsSchemaJSON)
	defer ts.Close()
	sess := fullSessionConfig()
	sess.Config.Endpoint = ts.URL
	client := NewController(sess)

	expected := testGetAddressCustomFieldsSchemaExpected
	actual, err := client.GetAddressCustomFieldsSchema()
	if err != nil {
		t.Fatalf("Bad: %s", err)
	}

	if !reflect.DeepEqual(expected, actual) {
		t.Fatalf("Expected %#v, got %#v", expected, actual)
	}
}

func TestUpdateAddress(t *testing.T) {
	ts := httpOKTestServer(testUpdateAddressOutputJSON)
	defer ts.Close()
	sess := fullSessionConfig()
	sess.Config.Endpoint = ts.URL
	client := NewController(sess)

	in := testUpdateAddressInput
	expected := testUpdateAddressOutputExpected
	actual, err := client.UpdateAddress(in)
	if err != nil {
		t.Fatalf("Bad: %s", err)
	}

	if !reflect.DeepEqual(expected, actual) {
		t.Fatalf("Expected %#v, got %#v", expected, actual)
	}
}

func TestDeleteAddress(t *testing.T) {
	ts := httpOKTestServer(testDeleteAddressOutputJSON)
	defer ts.Close()
	sess := fullSessionConfig()
	sess.Config.Endpoint = ts.URL
	client := NewController(sess)

	expected := testDeleteAddressOutputExpected
	actual, err := client.DeleteAddress(11, false)
	if err != nil {
		t.Fatalf("Bad: %s", err)
	}

	if !reflect.DeepEqual(expected, actual) {
		t.Fatalf("Expected %#v, got %#v", expected, actual)
	}
}

// testAccAddressCRUDCreate tests the creation part of the addresss controller
// CRUD acceptance test.
func testAccAddressCRUDCreate(t *testing.T, sess *session.Session, a Address) {
	c := NewController(sess)

	if _, err := c.CreateAddress(a); err != nil {
		t.Fatalf("Create: Error creating address: %s", err)
	}
}

// testAccAddressCRUDReadByIP tests the read part of the addresss controller
// acceptance test, by fetching the address by IP. This is the first part of
// the 2-part read test, and also returns the ID of the address so that the
// test fixutre can be updated.
func testAccAddressCRUDReadByIP(t *testing.T, sess *session.Session, a Address) int {
	c := NewController(sess)

	out, err := c.GetAddressesByIP(a.IPAddress)
	if err != nil {
		t.Fatalf("Can't get address by IP: %s", err)
	}

	for _, v := range out {
		// We don't have an ID yet here, so set it.
		a.ID = v.ID
		if reflect.DeepEqual(a, v) {
			return v.ID
		}
	}

	t.Fatalf("ReadByIP: Could not find address %#v in %#v", a, out)
	return 0
}

// testAccAddressCRUDReadByID tests the read part of the addresss controller
// acceptance test, by fetching the address by ID. This is the second part of
// the 2-part read test.
func testAccAddressCRUDReadByID(t *testing.T, sess *session.Session, a Address) {
	c := NewController(sess)

	out, err := c.GetAddressByID(a.ID)
	if err != nil {
		t.Fatalf("Can't find address by ID: %s", err)
	}

	if !reflect.DeepEqual(a, out) {
		t.Fatalf("ReadByID: Expected %#v, got %#v", a, out)
	}
}

// testAccAddressCRUDUpdate tests the update part of the addresss controller
// acceptance test.
func testAccAddressCRUDUpdate(t *testing.T, sess *session.Session, a Address) {
	c := NewController(sess)

	// IP and subnetID can't be in request
	params := a
	params.IPAddress = ""
	params.SubnetID = 0

	if _, err := c.UpdateAddress(params); err != nil {
		t.Fatalf("Error updating address: %s", err)
	}

	// Assert update
	out, err := c.GetAddressByID(a.ID)

	if err != nil {
		t.Fatalf("Error fetching address after update: %s", err)
	}

	// Update updated date in original
	a.EditDate = out.EditDate

	if !reflect.DeepEqual(a, out) {
		t.Fatalf("Error after update: expected %#v, got %#v", a, out)
	}
}

// testAccAddressCRUDDelete tests the delete part of the addresss controller
// acceptance test.
func testAccAddressCRUDDelete(t *testing.T, sess *session.Session, a Address) {
	c := NewController(sess)

	if _, err := c.DeleteAddress(a.ID, false); err != nil {
		t.Fatalf("Error deleting address: %s", err)
	}

	// check to see if address is actually gone
	if _, err := c.GetAddressByID(a.ID); err == nil {
		t.Fatalf("Address still present after delete")
	}
}

// TestAccAddressCRUD runs a full create-read-update-delete test for a PHPIPAM
// address.
func TestAccAddressCRUD(t *testing.T) {
	testacc.VetAccConditions(t)

	sess := session.NewSession()
	address := testCreateAddressInput
	if os.Getenv("TESTACC_CUSTOM_NESTED") != "" {
		address.CustomFields = map[string]interface{}{
			"CustomTestAddresses":  "foobar",
			"CustomTestAddresses2": nil,
		}
	} else {
		log.Println("Note: Not testing nested custom fields as TESTACC_CUSTOM_NESTED is not set")
	}
	testAccAddressCRUDCreate(t, sess, address)
	// tag goes to used (default ID 2) when an IP is created
	address.Tag = 2
	address.ID = testAccAddressCRUDReadByIP(t, sess, address)
	testAccAddressCRUDReadByID(t, sess, address)
	address.Description = "foobaz"
	if os.Getenv("TESTACC_CUSTOM_NESTED") != "" {
		address.CustomFields["CustomTestAddresses"] = "bazboop"
	}
	testAccAddressCRUDUpdate(t, sess, address)
	testAccAddressCRUDDelete(t, sess, address)
}

// TestAccGetAddressCustomFieldsSchema tests GetAddressCustomFieldsSchema against
// a live PHPIPAM instance.
func TestAccGetAddressCustomFieldsSchema(t *testing.T) {
	testacc.VetAccConditions(t)

	sess := session.NewSession()
	client := NewController(sess)

	expected := testGetAddressCustomFieldsSchemaExpected
	actual, err := client.GetAddressCustomFieldsSchema()
	if err != nil {
		t.Fatalf("Bad: %s", err)
	}

	if !reflect.DeepEqual(expected, actual) {
		t.Fatalf("Expected %#v, got %#v", expected, actual)
	}
}

// testAccAddressCustomFieldUpdate adds a custom field to an existing subnet
// entry, and verify the data changed by reading it back. Technically, this
// covers both UpdateAddressCustomFields and GetAddressCustomFields.
func testAccAddressCustomFieldUpdateRead(t *testing.T, sess *session.Session, id int, fields map[string]interface{}) {
	c := NewController(sess)

	if _, err := c.UpdateAddressCustomFields(id, fields); err != nil {
		t.Fatalf("Error updating subnet custom fields: %s", err)
	}

	expected := fields
	actual, err := c.GetAddressCustomFields(id)
	if err != nil {
		t.Fatalf("Error fetching custom fields after update: %s", err)
	}

	if !reflect.DeepEqual(expected, actual) {
		t.Fatalf("Expected %#v, got %#v", expected, actual)
	}
}

// TestAccAddressCustomFieldUpdateRead runs acceptance tests for
// UpdateAddressCustomFields and GetAddressCustomFields, by setting a value and
// then reading it back to make sure it updated. This test is skipped if
// TESTACC_CUSTOM_NESTED is set.
//
// We do this a few times to make sure that custom fields can be updated
// correctly.
func TestAccAddressCustomFieldUpdateRead(t *testing.T) {
	testacc.VetAccConditions(t)
	testacc.SkipIfCustomNested(t)

	sess := session.NewSession()
	fields := map[string]interface{}{
		"CustomTestAddresses":  "foobar",
		"CustomTestAddresses2": nil,
	}

	// We create a brand new address for this so we don't interfere with other
	// testing that works off of existing data.
	address := testCreateAddressInput
	testAccAddressCRUDCreate(t, sess, address)
	// tag goes to used (default ID 2) when an IP is created
	address.Tag = 2
	address.ID = testAccAddressCRUDReadByIP(t, sess, address)

	testAccAddressCustomFieldUpdateRead(t, sess, address.ID, fields)

	fields["CustomTestAddresses"] = "updated"
	testAccAddressCustomFieldUpdateRead(t, sess, address.ID, fields)

	// Clearing out a optional field will render it as a null field in the JSON
	// response, so it needs to be nil here and not just an empty string.
	fields["CustomTestAddresses"] = nil
	testAccAddressCustomFieldUpdateRead(t, sess, address.ID, fields)

	// clean up
	testAccAddressCRUDDelete(t, sess, address)
}
