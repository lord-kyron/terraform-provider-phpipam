package subnets

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"reflect"
	"testing"

	"github.com/pavel-z1/phpipam-sdk-go/controllers/addresses"
	"github.com/pavel-z1/phpipam-sdk-go/phpipam"
	"github.com/pavel-z1/phpipam-sdk-go/phpipam/session"
	"github.com/pavel-z1/phpipam-sdk-go/testacc"
)

var testCreateSubnetInput = Subnet{
	SubnetAddress:  "10.10.3.0",
	Mask:           24,
	SectionID:      1,
	MasterSubnetID: 2,
}

const testCreateSubnetOutputExpected = `Subnet created`
const testCreateSubnetOutputJSON = `
{
  "code": 201,
  "success": true,
  "data": "Subnet created"
}
`

var testCreateFirstFreeSubnetInput = Subnet{
	Description: "Subnet1",
}

const testCreateFirstFreeSubnetOutputExpected = "10.10.4.0/25"
const testCreateFirstFreeSubnetOutputJSON = `
{
  "code": 201,
  "success": true,
  "message": "Subnet created",
  "id": "10",
  "data": "10.10.4.0/25"
}
`

var testGetSubnetByIDOutputExpected = Subnet{
	ID:             8,
	SubnetAddress:  "10.10.3.0",
	Mask:           24,
	SectionID:      1,
	MasterSubnetID: 2,
}

const testGetSubnetByIDOutputJSON = `
{
  "code": 200,
  "success": true,
  "data": {
    "id": "8",
    "subnet": "10.10.3.0",
    "mask": "24",
    "sectionId": "1",
    "description": null,
    "firewallAddressObject": null,
    "vrfId": null,
    "masterSubnetId": "2",
    "allowRequests": "0",
    "vlanId": null,
    "showName": "0",
    "device": "0",
    "permissions": null,
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
    "links": [
      {
        "rel": "self",
        "href": "/api/test/subnets/8/",
        "methods": [
          "GET",
          "POST",
          "DELETE",
          "PATCH"
        ]
      },
      {
        "rel": "addresses",
        "href": "/api/test/subnets/8/addresses/",
        "methods": [
          "GET"
        ]
      },
      {
        "rel": "usage",
        "href": "/api/test/subnets/8/usage/",
        "methods": [
          "GET"
        ]
      },
      {
        "rel": "first_free",
        "href": "/api/test/subnets/8/first_free/",
        "methods": [
          "GET"
        ]
      },
      {
        "rel": "slaves",
        "href": "/api/test/subnets/8/slaves/",
        "methods": [
          "GET"
        ]
      },
      {
        "rel": "slaves_recursive",
        "href": "/api/test/subnets/8/slaves_recursive/",
        "methods": [
          "GET"
        ]
      },
      {
        "rel": "truncate",
        "href": "/api/test/subnets/8/truncate/",
        "methods": [
          "DELETE"
        ]
      },
      {
        "rel": "resize",
        "href": "/api/test/subnets/8/resize/",
        "methods": [
          "PATCH"
        ]
      },
      {
        "rel": "split",
        "href": "/api/test/subnets/8/split/",
        "methods": [
          "PATCH"
        ]
      }
    ]
  }
}
`

var testGetSubnetsByCIDROutputExpected = []Subnet{
	Subnet{
		ID:             8,
		SubnetAddress:  "10.10.3.0",
		Mask:           24,
		SectionID:      1,
		MasterSubnetID: 2,
	},
}

const testGetSubnetsByCIDROutputJSON = `
{
  "code": 200,
  "success": true,
  "data": [
    {
      "id": "8",
      "subnet": "10.10.3.0",
      "mask": "24",
      "sectionId": "1",
      "description": null,
      "firewallAddressObject": null,
      "vrfId": null,
      "masterSubnetId": "2",
      "allowRequests": "0",
      "vlanId": null,
      "showName": "0",
      "device": "0",
      "permissions": null,
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
      "links": [
        {
          "rel": "self",
          "href": "/api/test/subnets/8/"
        }
      ]
    }
  ]
}
`

const testGetFirstFreeSubnetOutputExpected = "10.10.4.0/25"
const testGetFirstFreeSubnetOutputJSON = `
{
  "code": 200,
  "success": true,
  "data": "10.10.4.0/25"
}
`

const testGetFirstFreeAddressOutputExpected = "10.10.1.1"
const testGetFirstFreeAddressOutputJSON = `
{
  "code": 200,
  "success": true,
  "data": "10.10.1.1"
}
`

var testGetAddressesInSubnetExpected = []addresses.Address{
	addresses.Address{
		ID:          1,
		SubnetID:    3,
		IPAddress:   "10.10.1.3",
		IsGateway:   false,
		Description: "Server1",
		Hostname:    "server1.cust1.local",
		Tag:         2,
		LastSeen:    "1970-01-01 00:00:01",
	},
	addresses.Address{
		ID:          2,
		SubnetID:    3,
		IPAddress:   "10.10.1.4",
		IsGateway:   false,
		Description: "Server2",
		Hostname:    "server2.cust1.local",
		Tag:         2,
		LastSeen:    "1970-01-01 00:00:01",
	},
	addresses.Address{
		ID:          3,
		SubnetID:    3,
		IPAddress:   "10.10.1.5",
		IsGateway:   false,
		Description: "Server3",
		Hostname:    "server3.cust1.local",
		Tag:         3,
		LastSeen:    "1970-01-01 00:00:01",
	},
	addresses.Address{
		ID:          4,
		SubnetID:    3,
		IPAddress:   "10.10.1.6",
		IsGateway:   false,
		Description: "Server4",
		Hostname:    "server4.cust1.local",
		Tag:         3,
		LastSeen:    "1970-01-01 00:00:01",
	},
	addresses.Address{
		ID:          5,
		SubnetID:    3,
		IPAddress:   "10.10.1.245",
		IsGateway:   false,
		Description: "Gateway",
		Tag:         2,
		LastSeen:    "1970-01-01 00:00:01",
	},
}

const testGetAddressesInSubnetJSON = `
{
  "code": 200,
  "success": true,
  "data": [
    {
      "id": "1",
      "subnetId": "3",
      "ip": "10.10.1.3",
      "is_gateway": "0",
      "description": "Server1",
      "hostname": "server1.cust1.local",
      "mac": null,
      "owner": null,
      "tag": "2",
      "deviceId": null,
      "port": null,
      "note": null,
			"lastSeen": "1970-01-01 00:00:01",
      "excludePing": "0",
      "PTRignore": "0",
      "PTR": "0",
      "firewallAddressObject": null,
      "editDate": null,
      "links": [
        {
          "rel": "self",
          "href": "/api/test/addresses/1/"
        }
      ]
    },
    {
      "id": "2",
      "subnetId": "3",
      "ip": "10.10.1.4",
      "is_gateway": "0",
      "description": "Server2",
      "hostname": "server2.cust1.local",
      "mac": null,
      "owner": null,
      "tag": "2",
      "deviceId": null,
      "port": null,
      "note": null,
			"lastSeen": "1970-01-01 00:00:01",
      "excludePing": "0",
      "PTRignore": "0",
      "PTR": "0",
      "firewallAddressObject": null,
      "editDate": null,
      "links": [
        {
          "rel": "self",
          "href": "/api/test/addresses/2/"
        }
      ]
    },
    {
      "id": "3",
      "subnetId": "3",
      "ip": "10.10.1.5",
      "is_gateway": "0",
      "description": "Server3",
      "hostname": "server3.cust1.local",
      "mac": null,
      "owner": null,
      "tag": "3",
      "deviceId": null,
      "port": null,
      "note": null,
			"lastSeen": "1970-01-01 00:00:01",
      "excludePing": "0",
      "PTRignore": "0",
      "PTR": "0",
      "firewallAddressObject": null,
      "editDate": null,
      "links": [
        {
          "rel": "self",
          "href": "/api/test/addresses/3/"
        }
      ]
    },
    {
      "id": "4",
      "subnetId": "3",
      "ip": "10.10.1.6",
      "is_gateway": "0",
      "description": "Server4",
      "hostname": "server4.cust1.local",
      "mac": null,
      "owner": null,
      "tag": "3",
      "deviceId": null,
      "port": null,
      "note": null,
			"lastSeen": "1970-01-01 00:00:01",
      "excludePing": "0",
      "PTRignore": "0",
      "PTR": "0",
      "firewallAddressObject": null,
      "editDate": null,
      "links": [
        {
          "rel": "self",
          "href": "/api/test/addresses/4/"
        }
      ]
    },
    {
      "id": "5",
      "subnetId": "3",
      "ip": "10.10.1.245",
      "is_gateway": "0",
      "description": "Gateway",
      "hostname": null,
      "mac": null,
      "owner": null,
      "tag": "2",
      "deviceId": null,
      "port": null,
      "note": null,
			"lastSeen": "1970-01-01 00:00:01",
      "excludePing": "0",
      "PTRignore": "0",
      "PTR": "0",
      "firewallAddressObject": null,
      "editDate": null,
      "links": [
        {
          "rel": "self",
          "href": "/api/test/addresses/5/"
        }
      ]
    }
  ]
}
`

var testGetSubnetCustomFieldsSchemaExpected = map[string]phpipam.CustomField{
	"CustomTestSubnets": phpipam.CustomField{
		Name:    "CustomTestSubnets",
		Type:    "varchar(255)",
		Comment: "Test field for subnets controller",
		Null:    "YES",
		Default: "",
	},
	"CustomTestSubnets2": phpipam.CustomField{
		Name:    "CustomTestSubnets2",
		Type:    "varchar(255)",
		Comment: "Test field for subnets controller (second field)",
		Null:    "YES",
		Default: "",
	},
}

const testGetSubnetCustomFieldsSchemaJSON = `
{
  "code": 200,
  "success": true,
  "data": {
    "CustomTestSubnets": {
      "name": "CustomTestSubnets",
      "type": "varchar(255)",
      "Comment": "Test field for subnets controller",
      "Null": "YES",
      "Default": null
    },
    "CustomTestSubnets2": {
      "name": "CustomTestSubnets2",
      "type": "varchar(255)",
      "Comment": "Test field for subnets controller (second field)",
      "Null": "YES",
      "Default": null
    }
  }
}
`

var testUpdateSubnetInput = Subnet{
	ID:          8,
	Description: "foobat",
}

const testUpdateSubnetOutputExpected = `Subnet updated`
const testUpdateSubnetOutputJSON = `
{
  "code": 200,
  "success": true,
  "data": "Subnet updated"
}
`

const testDeleteSubnetOutputExpected = `Subnet deleted`
const testDeleteSubnetOutputJSON = `
{
  "code": 200,
  "success": true,
  "data": "Subnet deleted"
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

func TestCreateSubnet(t *testing.T) {
	ts := httpCreatedTestServer(testCreateSubnetOutputJSON)
	defer ts.Close()
	sess := fullSessionConfig()
	sess.Config.Endpoint = ts.URL
	client := NewController(sess)

	in := testCreateSubnetInput
	expected := testCreateSubnetOutputExpected
	actual, err := client.CreateSubnet(in)
	if err != nil {
		t.Fatalf("Bad: %s", err)
	}

	if !reflect.DeepEqual(expected, actual) {
		t.Fatalf("Expected %#v, got %#v", expected, actual)
	}
}

func TestCreateFirstFreeSubnet(t *testing.T){
	ts := httpCreatedTestServer(testCreateFirstFreeSubnetOutputJSON)
	defer ts.Close()
	sess := fullSessionConfig()
	sess.Config.Endpoint = ts.URL
	client := NewController(sess)

	in := testCreateFirstFreeSubnetInput
	mask := 25
	id := 2
	expected := testCreateFirstFreeSubnetOutputExpected
	actual, err := client.CreateFirstFreeSubnet(id, mask, in)
	if err != nil {
		t.Fatalf("Bad: %s", err)
	}

	if !reflect.DeepEqual(expected, actual) {
		t.Fatalf("Expected %#v, got %#v", expected, actual)
	}
}

func TestGetSubnetByID(t *testing.T) {
	ts := httpOKTestServer(testGetSubnetByIDOutputJSON)
	defer ts.Close()
	sess := fullSessionConfig()
	sess.Config.Endpoint = ts.URL
	client := NewController(sess)

	expected := testGetSubnetByIDOutputExpected
	actual, err := client.GetSubnetByID(8)
	if err != nil {
		t.Fatalf("Bad: %s", err)
	}

	if !reflect.DeepEqual(expected, actual) {
		t.Fatalf("Expected %#v, got %#v", expected, actual)
	}
}

func TestGetSubnetsByCIDR(t *testing.T) {
	ts := httpOKTestServer(testGetSubnetsByCIDROutputJSON)
	defer ts.Close()
	sess := fullSessionConfig()
	sess.Config.Endpoint = ts.URL
	client := NewController(sess)

	expected := testGetSubnetsByCIDROutputExpected
	actual, err := client.GetSubnetsByCIDR("10.10.3.0/24")
	if err != nil {
		t.Fatalf("Bad: %s", err)
	}

	if !reflect.DeepEqual(expected, actual) {
		t.Fatalf("Expected %#v, got %#v", expected, actual)
	}
}

func TestGetFirstFreeSubnet(t *testing.T) {
	ts := httpOKTestServer(testGetFirstFreeSubnetOutputJSON)
	defer ts.Close()
	sess := fullSessionConfig()
	sess.Config.Endpoint = ts.URL
	client := NewController(sess)

	id := 2
	mask := 25
	expected := testGetFirstFreeSubnetOutputExpected
	actual, err := client.GetFirstFreeSubnet(id, mask)
	if err != nil {
		t.Fatalf("Bad: %s", err)
	}

	if !reflect.DeepEqual(expected, actual) {
		t.Fatalf("Expected %#v, got %#v", expected, actual)
	}
}

func TestGetFirstFreeAddress(t *testing.T) {
	ts := httpOKTestServer(testGetFirstFreeAddressOutputJSON)
	defer ts.Close()
	sess := fullSessionConfig()
	sess.Config.Endpoint = ts.URL
	client := NewController(sess)

	expected := testGetFirstFreeAddressOutputExpected
	actual, err := client.GetFirstFreeAddress(3)
	if err != nil {
		t.Fatalf("Bad: %s", err)
	}

	if expected != actual {
		t.Fatalf("Expected %#v, got %#v", expected, actual)
	}
}

func TestGetAddressesInSubnet(t *testing.T) {
	ts := httpOKTestServer(testGetAddressesInSubnetJSON)
	defer ts.Close()
	sess := fullSessionConfig()
	sess.Config.Endpoint = ts.URL
	client := NewController(sess)

	expected := testGetAddressesInSubnetExpected
	actual, err := client.GetAddressesInSubnet(3)
	if err != nil {
		t.Fatalf("Bad: %s", err)
	}

	if !reflect.DeepEqual(expected, actual) {
		t.Fatalf("Expected %#v, got %#v", expected, actual)
	}
}

func TestGetSubnetCustomFieldsSchema(t *testing.T) {
	ts := httpOKTestServer(testGetSubnetCustomFieldsSchemaJSON)
	defer ts.Close()
	sess := fullSessionConfig()
	sess.Config.Endpoint = ts.URL
	client := NewController(sess)

	expected := testGetSubnetCustomFieldsSchemaExpected
	actual, err := client.GetSubnetCustomFieldsSchema()
	if err != nil {
		t.Fatalf("Bad: %s", err)
	}

	if !reflect.DeepEqual(expected, actual) {
		t.Fatalf("Expected %#v, got %#v", expected, actual)
	}
}

func TestUpdateSubnet(t *testing.T) {
	ts := httpOKTestServer(testUpdateSubnetOutputJSON)
	defer ts.Close()
	sess := fullSessionConfig()
	sess.Config.Endpoint = ts.URL
	client := NewController(sess)

	in := testUpdateSubnetInput
	expected := testUpdateSubnetOutputExpected
	actual, err := client.UpdateSubnet(in)
	if err != nil {
		t.Fatalf("Bad: %s", err)
	}

	if !reflect.DeepEqual(expected, actual) {
		t.Fatalf("Expected %#v, got %#v", expected, actual)
	}
}

func TestDeleteSubnet(t *testing.T) {
	ts := httpOKTestServer(testDeleteSubnetOutputJSON)
	defer ts.Close()
	sess := fullSessionConfig()
	sess.Config.Endpoint = ts.URL
	client := NewController(sess)

	expected := testDeleteSubnetOutputExpected
	actual, err := client.DeleteSubnet(8)
	if err != nil {
		t.Fatalf("Bad: %s", err)
	}

	if !reflect.DeepEqual(expected, actual) {
		t.Fatalf("Expected %#v, got %#v", expected, actual)
	}
}

// testAccSubnetCRUDCreate tests the creation part of the subnets controller
// CRUD acceptance test.
func testAccSubnetCRUDCreate(t *testing.T, sess *session.Session, s Subnet) {
	c := NewController(sess)

	if _, err := c.CreateSubnet(s); err != nil {
		t.Fatalf("Create: Error creating subnet: %s", err)
	}
}

// testAccSubnetCRUDReadByCIDR tests the read part of the subnets controller
// acceptance test, by fetching the subnet by CIDR. This is the first part of
// the 3-part read test, and also returns the ID of the subnet so that the
// test fixutre can be updated.
func testAccSubnetCRUDReadByCIDR(t *testing.T, sess *session.Session, s Subnet) int {
	c := NewController(sess)

	out, err := c.GetSubnetsByCIDR(fmt.Sprintf("%s/%d", s.SubnetAddress, s.Mask))
	if err != nil {
		t.Fatalf("Can't get subnet by CIDR: %s", err)
	}

	for _, v := range out {
		// We don't have an ID yet here, so set it.
		s.ID = v.ID
		if reflect.DeepEqual(s, v) {
			return v.ID
		}
	}

	t.Fatalf("ReadByCIDR: Could not find subnet %#v in %#v", s, out)
	return 0
}

// testAccSubnetCRUDReadFirstFreeAddress tests the read part of the subnets
// controller acceptance test, by fetching the first available address in the
// created subnet. This is the second part of the 3-part read test, and also
// returns the ID of the subnet so that the test fixutre can be updated.
func testAccSubnetCRUDReadFirstFreeAddress(t *testing.T, sess *session.Session, s Subnet) {
	c := NewController(sess)

	out, err := c.GetFirstFreeAddress(s.ID)
	if err != nil {
		t.Fatalf("Can't read first free address: %s", err)
	}

	if out != "10.10.3.1" {
		t.Fatalf("Expected first free address to be 10.10.3.1, got %s", out)
	}
}

// testAccSubnetCRUDReadByID tests the read part of the subnets controller
// acceptance test, by fetching the subnet by ID. This is the third part of
// the 3-part read test.
func testAccSubnetCRUDReadByID(t *testing.T, sess *session.Session, s Subnet) {
	c := NewController(sess)

	out, err := c.GetSubnetByID(s.ID)
	if err != nil {
		t.Fatalf("Can't find subnet by ID: %s", err)
	}

	if !reflect.DeepEqual(s, out) {
		t.Fatalf("ReadByID: Expected %#v, got %#v", s, out)
	}
}

// testAccSubnetCRUDUpdate tests the update part of the subnets controller
// acceptance test.
func testAccSubnetCRUDUpdate(t *testing.T, sess *session.Session, s Subnet) {
	c := NewController(sess)

	// Address or mask can't be in an update request.
	params := s
	params.SubnetAddress = ""
	params.Mask = 0

	if _, err := c.UpdateSubnet(params); err != nil {
		t.Fatalf("Error updating subnet: %s", err)
	}

	// Assert update
	out, err := c.GetSubnetByID(s.ID)

	if err != nil {
		t.Fatalf("Error fetching subnet after update: %s", err)
	}

	// Update updated date in original
	s.EditDate = out.EditDate

	if !reflect.DeepEqual(s, out) {
		t.Fatalf("Error after update: expected %#v, got %#v", s, out)
	}
}

// testAccSubnetCRUDDelete tests the delete part of the subnets controller
// acceptance test.
func testAccSubnetCRUDDelete(t *testing.T, sess *session.Session, s Subnet) {
	c := NewController(sess)

	if _, err := c.DeleteSubnet(s.ID); err != nil {
		t.Fatalf("Error deleting subnet: %s", err)
	}

	// check to see if subnet is actually gone
	if _, err := c.GetSubnetByID(s.ID); err == nil {
		t.Fatalf("Subnet still present after delete")
	}
}

// TestAccSubnetCRUD runs a full create-read-update-delete test for a PHPIPAM
// subnet.
func TestAccSubnetCRUD(t *testing.T) {
	testacc.VetAccConditions(t)

	sess := session.NewSession()
	subnet := testCreateSubnetInput
	// Permissions get added even though they are optional
	subnet.Permissions = "{\"3\":\"1\",\"2\":\"2\"}"
	if os.Getenv("TESTACC_CUSTOM_NESTED") != "" {
		subnet.CustomFields = map[string]interface{}{
			"CustomTestSubnets":  "foobar",
			"CustomTestSubnets2": nil,
		}
	} else {
		log.Println("Note: Not testing nested custom fields as TESTACC_CUSTOM_NESTED is not set")
	}
	testAccSubnetCRUDCreate(t, sess, subnet)
	subnet.ID = testAccSubnetCRUDReadByCIDR(t, sess, subnet)
	testAccSubnetCRUDReadByID(t, sess, subnet)
	subnet.Description = "Updating subnet!"
	if os.Getenv("TESTACC_CUSTOM_NESTED") != "" {
		subnet.CustomFields["CustomTestSubnets"] = "bazboop"
	}
	testAccSubnetCRUDUpdate(t, sess, subnet)
	testAccSubnetCRUDDelete(t, sess, subnet)
}

// TestAccGetAddressesInSubnet tests GetAddressesInSubnet against a live PHPIPAM
// instance.
func TestAccGetAddressesInSubnet(t *testing.T) {
	testacc.VetAccConditions(t)

	sess := session.NewSession()
	client := NewController(sess)

	expected := testGetAddressesInSubnetExpected
	if os.Getenv("TESTACC_CUSTOM_NESTED") != "" {
		for n := range expected {
			expected[n].CustomFields = map[string]interface{}{
				"CustomTestAddresses":  nil,
				"CustomTestAddresses2": nil,
			}
		}
	} else {
		log.Println("Note: Not testing nested custom fields as TESTACC_CUSTOM_NESTED is not set")
	}
	actual, err := client.GetAddressesInSubnet(3)
	if err != nil {
		t.Fatalf("Bad: %s", err)
	}

	if !reflect.DeepEqual(expected, actual) {
		t.Fatalf("Expected %#v, got %#v", expected, actual)
	}
}

// TestAccGetSubnetCustomFieldsSchema tests GetSubnetCustomFieldsSchema against
// a live PHPIPAM instance.
func TestAccGetSubnetCustomFieldsSchema(t *testing.T) {
	testacc.VetAccConditions(t)

	sess := session.NewSession()
	client := NewController(sess)

	expected := testGetSubnetCustomFieldsSchemaExpected
	actual, err := client.GetSubnetCustomFieldsSchema()
	if err != nil {
		t.Fatalf("Bad: %s", err)
	}

	if !reflect.DeepEqual(expected, actual) {
		t.Fatalf("Expected %#v, got %#v", expected, actual)
	}
}

// testAccSubnetCustomFieldUpdate adds a custom field to an existing subnet
// entry, and verify the data changed by reading it back. Technically, this
// covers both UpdateSubnetCustomFields and GetSubnetCustomFields.
func testAccSubnetCustomFieldUpdateRead(t *testing.T, sess *session.Session, id int, fields map[string]interface{}) {
	c := NewController(sess)

	if _, err := c.UpdateSubnetCustomFields(id, fields); err != nil {
		t.Fatalf("Error updating subnet custom fields: %s", err)
	}

	expected := fields
	actual, err := c.GetSubnetCustomFields(id)
	if err != nil {
		t.Fatalf("Error fetching custom fields after update: %s", err)
	}

	if !reflect.DeepEqual(expected, actual) {
		t.Fatalf("Expected %#v, got %#v", expected, actual)
	}
}

// TestAccSubnetCustomFieldUpdateRead runs acceptance tests for
// UpdateSubnetCustomFields and GetSubnetCustomFields, by setting a value and
// then reading it back to make sure it updated.
//
// We do this a few times to make sure that custom fields can be updated
// correctly.
func TestAccSubnetCustomFieldUpdateRead(t *testing.T) {
	testacc.VetAccConditions(t)
	testacc.SkipIfCustomNested(t)

	sess := session.NewSession()
	fields := map[string]interface{}{
		"CustomTestSubnets":  "foobar",
		"CustomTestSubnets2": nil,
	}

	// We create a brand new subnet for this so we don't interfere with other
	// testing that works off of existing data.
	subnet := testCreateSubnetInput
	// Permissions get added even though they are optional
	subnet.Permissions = "{\"3\":\"1\",\"2\":\"2\"}"
	testAccSubnetCRUDCreate(t, sess, subnet)
	subnet.ID = testAccSubnetCRUDReadByCIDR(t, sess, subnet)

	testAccSubnetCustomFieldUpdateRead(t, sess, subnet.ID, fields)

	fields["CustomTestSubnets"] = "updated"
	testAccSubnetCustomFieldUpdateRead(t, sess, subnet.ID, fields)

	// Clearing out a optional field will render it as a null field in the JSON
	// response, so it needs to be nil here and not just an empty string.
	fields["CustomTestSubnets"] = nil
	testAccSubnetCustomFieldUpdateRead(t, sess, subnet.ID, fields)

	// clean up
	testAccSubnetCRUDDelete(t, sess, subnet)
}
