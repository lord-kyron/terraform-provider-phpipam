package vlans

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

var testCreateVLANInput = VLAN{
	Name:   "foolan",
	Number: 1000,
}

const testCreateVLANOutputExpected = `Vlan created`
const testCreateVLANOutputJSON = `
{
  "code": 201,
  "success": true,
  "data": "Vlan created"
}
`

var testGetVLANByIDOutputExpected = VLAN{
	ID:       3,
	DomainID: 1,
	Name:     "foolan",
	Number:   1000,
}

const testGetVLANByIDOutputJSON = `
{
  "code": 200,
  "success": true,
  "data": {
    "id": "3",
    "domainId": "1",
    "name": "foolan",
    "number": "1000",
    "description": null,
    "editDate": null,
    "links": [
      {
        "rel": "self",
        "href": "/api/test/vlans/3/",
        "methods": [
          "GET",
          "POST",
          "DELETE",
          "PATCH"
        ]
      },
      {
        "rel": "subnets",
        "href": "/api/test/vlans/3/subnets/",
        "methods": [
          "GET"
        ]
      }
    ]
  }
}
`

var testGetVLANsByNumberOutputExpected = []VLAN{
	VLAN{
		ID:       3,
		DomainID: 1,
		Name:     "foolan",
		Number:   1000,
	},
}

const testGetVLANsByNumberOutputJSON = `
{
  "code": 200,
  "success": true,
  "data": [
    {
      "id": "3",
      "domainId": "1",
      "name": "foolan",
      "number": "1000",
      "description": null,
      "editDate": null,
      "links": [
        {
          "rel": "self",
          "href": "/api/test/vlans/3/"
        }
      ]
    }
  ]
}
`

var testGetVLANCustomFieldsSchemaExpected = map[string]phpipam.CustomField{
	"CustomTestVLANs": phpipam.CustomField{
		Name:    "CustomTestVLANs",
		Type:    "varchar(255)",
		Comment: "Test field for vlans controller",
		Null:    "YES",
		Default: "",
	},
	"CustomTestVLANs2": phpipam.CustomField{
		Name:    "CustomTestVLANs2",
		Type:    "varchar(255)",
		Comment: "Test field for vlans controller (second field)",
		Null:    "YES",
		Default: "",
	},
}

const testGetVLANCustomFieldsSchemaJSON = `
{
  "code": 200,
  "success": true,
  "data": {
    "CustomTestVLANs": {
      "name": "CustomTestVLANs",
      "type": "varchar(255)",
      "Comment": "Test field for vlans controller",
      "Null": "YES",
      "Default": ""
    },
    "CustomTestVLANs2": {
      "name": "CustomTestVLANs2",
      "type": "varchar(255)",
      "Comment": "Test field for vlans controller (second field)",
      "Null": "YES",
      "Default": ""
    }
  }
}
`

var testUpdateVLANInput = VLAN{
	ID:   3,
	Name: "bazlan",
}

const testUpdateVLANOutputExpected = `Vlan updated`
const testUpdateVLANOutputJSON = `
{
  "code": 200,
  "success": true,
  "data": "Vlan updated"
}
`

const testDeleteVLANOutputExpected = `Vlan deleted`
const testDeleteVLANOutputJSON = `
{
  "code": 200,
  "success": true,
  "data": "Vlan deleted"
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

func TestCreateVLAN(t *testing.T) {
	ts := httpCreatedTestServer(testCreateVLANOutputJSON)
	defer ts.Close()
	sess := fullSessionConfig()
	sess.Config.Endpoint = ts.URL
	client := NewController(sess)

	in := testCreateVLANInput
	expected := testCreateVLANOutputExpected
	actual, err := client.CreateVLAN(in)
	if err != nil {
		t.Fatalf("Bad: %s", err)
	}

	if !reflect.DeepEqual(expected, actual) {
		t.Fatalf("Expected %#v, got %#v", expected, actual)
	}
}

func TestGetVLANByID(t *testing.T) {
	ts := httpOKTestServer(testGetVLANByIDOutputJSON)
	defer ts.Close()
	sess := fullSessionConfig()
	sess.Config.Endpoint = ts.URL
	client := NewController(sess)

	expected := testGetVLANByIDOutputExpected
	actual, err := client.GetVLANByID(3)
	if err != nil {
		t.Fatalf("Bad: %s", err)
	}

	if !reflect.DeepEqual(expected, actual) {
		t.Fatalf("Expected %#v, got %#v", expected, actual)
	}
}

func TestGetVLANsByNumber(t *testing.T) {
	ts := httpOKTestServer(testGetVLANsByNumberOutputJSON)
	defer ts.Close()
	sess := fullSessionConfig()
	sess.Config.Endpoint = ts.URL
	client := NewController(sess)

	expected := testGetVLANsByNumberOutputExpected
	actual, err := client.GetVLANsByNumber(1000)
	if err != nil {
		t.Fatalf("Bad: %s", err)
	}

	if !reflect.DeepEqual(expected, actual) {
		t.Fatalf("Expected %#v, got %#v", expected, actual)
	}
}

func TestGetVLANCustomFieldsSchema(t *testing.T) {
	ts := httpOKTestServer(testGetVLANCustomFieldsSchemaJSON)
	defer ts.Close()
	sess := fullSessionConfig()
	sess.Config.Endpoint = ts.URL
	client := NewController(sess)

	expected := testGetVLANCustomFieldsSchemaExpected
	actual, err := client.GetVLANCustomFieldsSchema()
	if err != nil {
		t.Fatalf("Bad: %s", err)
	}

	if !reflect.DeepEqual(expected, actual) {
		t.Fatalf("Expected %#v, got %#v", expected, actual)
	}
}

func TestUpdateVLAN(t *testing.T) {
	ts := httpOKTestServer(testUpdateVLANOutputJSON)
	defer ts.Close()
	sess := fullSessionConfig()
	sess.Config.Endpoint = ts.URL
	client := NewController(sess)

	in := testUpdateVLANInput
	expected := testUpdateVLANOutputExpected
	actual, err := client.UpdateVLAN(in)
	if err != nil {
		t.Fatalf("Bad: %s", err)
	}

	if !reflect.DeepEqual(expected, actual) {
		t.Fatalf("Expected %#v, got %#v", expected, actual)
	}
}

func TestDeleteVLAN(t *testing.T) {
	ts := httpOKTestServer(testDeleteVLANOutputJSON)
	defer ts.Close()
	sess := fullSessionConfig()
	sess.Config.Endpoint = ts.URL
	client := NewController(sess)

	expected := testDeleteVLANOutputExpected
	actual, err := client.DeleteVLAN(3)
	if err != nil {
		t.Fatalf("Bad: %s", err)
	}

	if !reflect.DeepEqual(expected, actual) {
		t.Fatalf("Expected %#v, got %#v", expected, actual)
	}
}

// testAccVLANCRUDCreate tests the creation part of the vlans controller
// CRUD acceptance test.
func testAccVLANCRUDCreate(t *testing.T, sess *session.Session, v VLAN) {
	c := NewController(sess)

	if _, err := c.CreateVLAN(v); err != nil {
		t.Fatalf("Create: Error creating vlan: %s", err)
	}
}

// testAccVLANCRUDReadByNumber tests the read part of the vlans controller
// acceptance test, by fetching the vlan by number. This is the first part of
// the 2-part read test, and also returns the ID of the vlan so that the
// test fixutre can be updated.
func testAccVLANCRUDReadByNumber(t *testing.T, sess *session.Session, v VLAN) int {
	c := NewController(sess)

	out, err := c.GetVLANsByNumber(v.Number)
	if err != nil {
		t.Fatalf("Can't get vlan by number: %s", err)
	}

	for _, val := range out {
		// We don't have an ID yet here, so set it.
		v.ID = val.ID
		if reflect.DeepEqual(v, val) {
			return val.ID
		}
	}

	t.Fatalf("ReadByNumber: Could not find vlan %#v in %#v", v, out)
	return 0
}

// testAccVLANCRUDReadByID tests the read part of the vlans controller
// acceptance test, by fetching the vlan by ID. This is the second part of
// the 2-part read test.
func testAccVLANCRUDReadByID(t *testing.T, sess *session.Session, v VLAN) {
	c := NewController(sess)

	out, err := c.GetVLANByID(v.ID)
	if err != nil {
		t.Fatalf("Can't find vlan by ID: %s", err)
	}

	if !reflect.DeepEqual(v, out) {
		t.Fatalf("ReadByID: Expected %#v, got %#v", v, out)
	}
}

// testAccVLANCRUDUpdate tests the update part of the vlans controller
// acceptance test.
func testAccVLANCRUDUpdate(t *testing.T, sess *session.Session, v VLAN) {
	c := NewController(sess)

	if _, err := c.UpdateVLAN(v); err != nil {
		t.Fatalf("Error updating vlan: %s", err)
	}

	// Assert update
	out, err := c.GetVLANByID(v.ID)

	if err != nil {
		t.Fatalf("Error fetching vlan after update: %s", err)
	}

	// Update updated date in original
	v.EditDate = out.EditDate

	if !reflect.DeepEqual(v, out) {
		t.Fatalf("Error after update: expected %#v, got %#v", v, out)
	}
}

// testAccVLANCRUDDelete tests the delete part of the vlans controller
// acceptance test.
func testAccVLANCRUDDelete(t *testing.T, sess *session.Session, v VLAN) {
	c := NewController(sess)

	if _, err := c.DeleteVLAN(v.ID); err != nil {
		t.Fatalf("Error deleting vlan: %s", err)
	}

	// check to see if vlan is actually gone
	if _, err := c.GetVLANByID(v.ID); err == nil {
		t.Fatalf("VLAN still present after delete")
	}
}

// TestAccVLANCRUD runs a full create-read-update-delete test for a PHPIPAM
// vlan.
func TestAccVLANCRUD(t *testing.T) {
	testacc.VetAccConditions(t)

	sess := session.NewSession()
	vlan := testCreateVLANInput
	if os.Getenv("TESTACC_CUSTOM_NESTED") != "" {
		vlan.CustomFields = map[string]interface{}{
			"CustomTestVLANs":  "foobar",
			"CustomTestVLANs2": nil,
		}
	} else {
		log.Println("Note: Not testing nested custom fields as TESTACC_CUSTOM_NESTED is not set")
	}
	testAccVLANCRUDCreate(t, sess, vlan)
	// Add the domain ID here as 1 is the default.
	vlan.DomainID = 1
	vlan.ID = testAccVLANCRUDReadByNumber(t, sess, vlan)
	testAccVLANCRUDReadByID(t, sess, vlan)
	vlan.Name = "bazlan"
	testAccVLANCRUDUpdate(t, sess, vlan)
	testAccVLANCRUDDelete(t, sess, vlan)
}

// TestAccGetVLANCustomFieldsSchema tests GetVLANCustomFieldsSchema against
// a live PHPIPAM instance.
func TestAccGetVLANCustomFieldsSchema(t *testing.T) {
	testacc.VetAccConditions(t)

	sess := session.NewSession()
	client := NewController(sess)

	expected := testGetVLANCustomFieldsSchemaExpected
	actual, err := client.GetVLANCustomFieldsSchema()
	if err != nil {
		t.Fatalf("Bad: %s", err)
	}

	if !reflect.DeepEqual(expected, actual) {
		t.Fatalf("Expected %#v, got %#v", expected, actual)
	}
}

// testAccVLANCustomFieldUpdate adds a custom field to an existing vlan
// entry, and verify the data changed by reading it back. Technically, this
// covers both UpdateVLANCustomFields and GetVLANCustomFields.
func testAccVLANCustomFieldUpdateRead(t *testing.T, sess *session.Session, id int, name string, fields map[string]interface{}) {
	c := NewController(sess)

	if _, err := c.UpdateVLANCustomFields(id, name, fields); err != nil {
		t.Fatalf("Error updating vlan custom fields: %s", err)
	}

	expected := fields
	actual, err := c.GetVLANCustomFields(id)
	if err != nil {
		t.Fatalf("Error fetching custom fields after update: %s", err)
	}

	if !reflect.DeepEqual(expected, actual) {
		t.Fatalf("Expected %#v, got %#v", expected, actual)
	}
}

// TestAccVLANCustomFieldUpdateRead runs acceptance tests for
// UpdateVLANCustomFields and GetVLANCustomFields, by setting a value and
// then reading it back to make sure it updated.
//
// We do this a few times to make sure that custom fields can be updated
// correctly.
func TestAccVLANCustomFieldUpdateRead(t *testing.T) {
	testacc.VetAccConditions(t)
	testacc.SkipIfCustomNested(t)

	sess := session.NewSession()
	fields := map[string]interface{}{
		"CustomTestVLANs":  "foobar",
		"CustomTestVLANs2": nil,
	}

	// We create a brand new vlan for this so we don't interfere with other
	// testing that works off of existing data.
	vlan := testCreateVLANInput
	testAccVLANCRUDCreate(t, sess, vlan)
	// Add the domain ID here as 1 is the default.
	vlan.DomainID = 1
	vlan.ID = testAccVLANCRUDReadByNumber(t, sess, vlan)

	testAccVLANCustomFieldUpdateRead(t, sess, vlan.ID, vlan.Name, fields)

	fields["CustomTestVLANs"] = "updated"
	testAccVLANCustomFieldUpdateRead(t, sess, vlan.ID, vlan.Name, fields)

	// Clearing out a optional field will render it as a null field in the JSON
	// response, so it needs to be nil here and not just an empty string.
	fields["CustomTestVLANs"] = nil
	testAccVLANCustomFieldUpdateRead(t, sess, vlan.ID, vlan.Name, fields)

	// clean up
	testAccVLANCRUDDelete(t, sess, vlan)
}
