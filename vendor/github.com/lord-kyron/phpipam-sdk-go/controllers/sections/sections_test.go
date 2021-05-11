package sections

import (
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"reflect"
	"testing"

	"github.com/davecgh/go-spew/spew"
	"github.com/pavel-z1/phpipam-sdk-go/controllers/subnets"
	"github.com/pavel-z1/phpipam-sdk-go/phpipam"
	"github.com/pavel-z1/phpipam-sdk-go/phpipam/session"
	"github.com/pavel-z1/phpipam-sdk-go/testacc"
)

var testListSectionsOutputExpected = []Section{
	Section{
		ID:          2,
		Name:        "IPv6",
		Description: "Section for IPv6 addresses",
		Permissions: "{\"3\":\"1\",\"2\":\"2\"}",
	},
	Section{
		ID:   3,
		Name: "foobar",
	},
	Section{
		ID:          1,
		Name:        "Customers",
		Description: "Section for customers",
		Permissions: "{\"3\":\"1\",\"2\":\"2\"}",
	},
}

const testListSectionsOutputJSON = `
{
  "code": 200,
  "success": true,
  "data": [
    {
      "id": "2",
      "name": "IPv6",
      "description": "Section for IPv6 addresses",
      "masterSection": "0",
      "permissions": "{\"3\":\"1\",\"2\":\"2\"}",
      "strictMode": "0",
      "subnetOrdering": null,
      "order": null,
      "editDate": null,
      "showVLAN": "0",
      "showVRF": "0",
      "DNS": null,
      "links": [
        {
          "rel": "self",
          "href": "/api/test/sections/2/"
        }
      ]
    },
    {
      "id": "3",
      "name": "foobar",
      "description": null,
      "masterSection": "0",
      "permissions": null,
      "strictMode": "0",
      "subnetOrdering": null,
      "order": null,
      "editDate": null,
      "showVLAN": "0",
      "showVRF": "0",
      "DNS": null,
      "links": [
        {
          "rel": "self",
          "href": "/api/test/sections/3/"
        }
      ]
    },
    {
      "id": "1",
      "name": "Customers",
      "description": "Section for customers",
      "masterSection": "0",
      "permissions": "{\"3\":\"1\",\"2\":\"2\"}",
      "strictMode": "0",
      "subnetOrdering": null,
      "order": null,
      "editDate": null,
      "showVLAN": "0",
      "showVRF": "0",
      "DNS": null,
      "links": [
        {
          "rel": "self",
          "href": "/api/test/sections/1/"
        }
      ]
    }
  ]
}
`

var testCreateSectionInput = Section{
	Name:        "foobar",
	StrictMode:  true,
	Permissions: "{\"3\":\"1\",\"2\":\"2\"}",
}

const testCreateSectionOutputExpected = `Section created`
const testCreateSectionOutputJSON = `
{
  "code": 201,
  "success": true,
  "data": "Section created"
}
`

var testGetSectionOutputExpected = Section{
	ID:          1,
	Name:        "Customers",
	Description: "Section for customers",
	Permissions: "{\"3\":\"1\",\"2\":\"2\"}",
}

const testGetSectionOutputJSON = `
{
  "code": 200,
  "success": true,
  "data": {
    "id": "1",
    "name": "Customers",
    "description": "Section for customers",
    "masterSection": "0",
    "permissions": "{\"3\":\"1\",\"2\":\"2\"}",
    "strictMode": "0",
    "subnetOrdering": null,
    "order": null,
    "editDate": null,
    "showVLAN": "0",
    "showVRF": "0",
    "DNS": null,
    "links": [
      {
        "rel": "self",
        "href": "/api/test/sections/1/",
        "methods": [
          "GET",
          "POST",
          "DELETE",
          "PATCH"
        ]
      },
      {
        "rel": "subnets",
        "href": "/api/test/sections/1/subnets/",
        "methods": [
          "GET"
        ]
      }
    ]
  }
}
`

const testGetSubnetsInSectionOutputJSON = `
{
  "code": 200,
  "success": true,
  "data": [
    {
      "id": "5",
      "subnet": "0.0.0.0",
      "mask": "",
      "sectionId": "1",
      "description": "My folder",
      "firewallAddressObject": null,
      "vrfId": "0",
      "masterSubnetId": "0",
      "allowRequests": "0",
      "vlanId": "0",
      "showName": "0",
      "device": "0",
      "permissions": "{\"3\":\"1\",\"2\":\"2\"}",
      "pingSubnet": "0",
      "discoverSubnet": "0",
      "DNSrecursive": "0",
      "DNSrecords": "0",
      "nameserverId": "0",
      "scanAgent": null,
      "isFolder": "1",
      "isFull": "0",
      "tag": "2",
      "editDate": null,
      "links": [
        {
          "rel": "self",
          "href": "/api/test/subnets/5/"
        }
      ]
    },
    {
      "id": "2",
      "subnet": "10.10.0.0",
      "mask": "16",
      "sectionId": "1",
      "description": "Business customers",
      "firewallAddressObject": null,
      "vrfId": "0",
      "masterSubnetId": "0",
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
      "links": [
        {
          "rel": "self",
          "href": "/api/test/subnets/2/"
        }
      ]
    },
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
      "links": [
        {
          "rel": "self",
          "href": "/api/test/subnets/3/"
        }
      ]
    },
    {
      "id": "4",
      "subnet": "10.10.2.0",
      "mask": "24",
      "sectionId": "1",
      "description": "Customer 2",
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
      "links": [
        {
          "rel": "self",
          "href": "/api/test/subnets/4/"
        }
      ]
    },
    {
      "id": "6",
      "subnet": "10.65.22.0",
      "mask": "24",
      "sectionId": "1",
      "description": "DHCP range",
      "firewallAddressObject": null,
      "vrfId": "0",
      "masterSubnetId": "5",
      "allowRequests": "0",
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
      "links": [
        {
          "rel": "self",
          "href": "/api/test/subnets/6/"
        }
      ]
    }
  ]
}
`

var testGetSubnetsInSectionExpected = []subnets.Subnet{
	subnets.Subnet{
		ID:             5,
		SubnetAddress:  "0.0.0.0",
		Mask:           0,
		SectionID:      1,
		MasterSubnetID: 0,
		AllowRequests:  false,
		Description:    "My folder",
		ShowName:       false,
		Permissions:    "{\"3\":\"1\",\"2\":\"2\"}",
		IsFolder:       true,
	},
	subnets.Subnet{
		ID:             2,
		SubnetAddress:  "10.10.0.0",
		Mask:           16,
		SectionID:      1,
		MasterSubnetID: 0,
		AllowRequests:  true,
		Description:    "Business customers",
		ShowName:       true,
		Permissions:    "{\"3\":\"1\",\"2\":\"2\"}",
	},
	subnets.Subnet{
		ID:             3,
		SubnetAddress:  "10.10.1.0",
		Mask:           24,
		SectionID:      1,
		MasterSubnetID: 2,
		AllowRequests:  true,
		Description:    "Customer 1",
		ShowName:       true,
		Permissions:    "{\"3\":\"1\",\"2\":\"2\"}",
	},
	subnets.Subnet{
		ID:             4,
		SubnetAddress:  "10.10.2.0",
		Mask:           24,
		SectionID:      1,
		MasterSubnetID: 2,
		AllowRequests:  true,
		Description:    "Customer 2",
		ShowName:       true,
		Permissions:    "{\"3\":\"1\",\"2\":\"2\"}",
	},
	subnets.Subnet{
		ID:             6,
		SubnetAddress:  "10.65.22.0",
		Mask:           24,
		SectionID:      1,
		MasterSubnetID: 5,
		AllowRequests:  false,
		Description:    "DHCP range",
		ShowName:       true,
		Permissions:    "{\"3\":\"1\",\"2\":\"2\"}",
	},
}

var testUpdateSectionInput = Section{
	ID:   3,
	Name: "foobaz",
}

const testUpdateSectionOutputJSON = `
{
  "code": 200,
  "success": true
}
`

var testDeleteSectionInput = Section{
	ID: 3,
}

const testDeleteSectionOutputJSON = `
{
  "code": 200,
  "success": true
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

func TestListSections(t *testing.T) {
	ts := httpOKTestServer(testListSectionsOutputJSON)
	defer ts.Close()
	sess := fullSessionConfig()
	sess.Config.Endpoint = ts.URL
	client := NewController(sess)

	expected := testListSectionsOutputExpected
	actual, err := client.ListSections()
	if err != nil {
		t.Fatalf("Bad: %s", err)
	}

	if !reflect.DeepEqual(expected, actual) {
		t.Fatalf("Expected %#v, got %#v", expected, actual)
	}
}

func TestCreateSection(t *testing.T) {
	ts := httpCreatedTestServer(testCreateSectionOutputJSON)
	defer ts.Close()
	sess := fullSessionConfig()
	sess.Config.Endpoint = ts.URL
	client := NewController(sess)

	in := testCreateSectionInput
	expected := testCreateSectionOutputExpected
	actual, err := client.CreateSection(in)
	if err != nil {
		t.Fatalf("Bad: %s", err)
	}

	if !reflect.DeepEqual(expected, actual) {
		t.Fatalf("Expected %#v, got %#v", expected, actual)
	}
}

func TestGetSectionByID(t *testing.T) {
	ts := httpOKTestServer(testGetSectionOutputJSON)
	defer ts.Close()
	sess := fullSessionConfig()
	sess.Config.Endpoint = ts.URL
	client := NewController(sess)

	expected := testGetSectionOutputExpected
	actual, err := client.GetSectionByID(1)
	if err != nil {
		t.Fatalf("Bad: %s", err)
	}

	if !reflect.DeepEqual(expected, actual) {
		t.Fatalf("Expected %#v, got %#v", expected, actual)
	}
}

func TestGetSectionByName(t *testing.T) {
	ts := httpOKTestServer(testGetSectionOutputJSON)
	defer ts.Close()
	sess := fullSessionConfig()
	sess.Config.Endpoint = ts.URL
	client := NewController(sess)

	expected := testGetSectionOutputExpected
	actual, err := client.GetSectionByName("Customers")
	if err != nil {
		t.Fatalf("Bad: %s", err)
	}

	if !reflect.DeepEqual(expected, actual) {
		t.Fatalf("Expected %#v, got %#v", expected, actual)
	}
}

func TestGetSubnetsInSection(t *testing.T) {
	ts := httpOKTestServer(testGetSubnetsInSectionOutputJSON)
	defer ts.Close()
	sess := fullSessionConfig()
	sess.Config.Endpoint = ts.URL
	client := NewController(sess)

	expected := testGetSubnetsInSectionExpected
	actual, err := client.GetSubnetsInSection(1)
	if err != nil {
		t.Fatalf("Bad: %s", err)
	}

	if !reflect.DeepEqual(expected, actual) {
		t.Fatalf("Expected %#v, got %#v", expected, actual)
	}
}

func TestUpdateSection(t *testing.T) {
	ts := httpOKTestServer(testUpdateSectionOutputJSON)
	defer ts.Close()
	sess := fullSessionConfig()
	sess.Config.Endpoint = ts.URL
	client := NewController(sess)

	in := testUpdateSectionInput
	err := client.UpdateSection(in)
	if err != nil {
		t.Fatalf("Bad: %s", err)
	}
}
func TestDeleteSection(t *testing.T) {
	ts := httpOKTestServer(testUpdateSectionOutputJSON)
	defer ts.Close()
	sess := fullSessionConfig()
	sess.Config.Endpoint = ts.URL
	client := NewController(sess)

	err := client.DeleteSection(3)
	if err != nil {
		t.Fatalf("Bad: %s", err)
	}
}

// testAccSectionsCRUDCreate tests the creation part of the sections controller
// CRUD acceptance test.
func testAccSectionsCRUDCreate(t *testing.T, sess *session.Session, s Section) {
	c := NewController(sess)

	if _, err := c.CreateSection(s); err != nil {
		t.Fatalf("Create: Error creating section: %s", err)
	}
}

// testAccSectionsCRUDReadByName tests the read part of the sections controller
// acceptance test, by fetching the section by name. This is the first part of
// the 3-part read test, and also returns the ID of the section so that the
// test fixutre can be updated.
func testAccSectionsCRUDReadByName(t *testing.T, sess *session.Session, s Section) int {
	c := NewController(sess)

	out, err := c.GetSectionByName(s.Name)
	if err != nil {
		t.Fatalf("Can't get section by name: %s", err)
	}
	// We don't have an ID yet here, so set it.
	s.ID = out.ID

	if !reflect.DeepEqual(s, out) {
		t.Fatalf("ReadByName: Expected %s, got %s", spew.Sdump(s), spew.Sdump(out))
	}
	return out.ID
}

// testAccSectionsCRUDReadByID tests the read part of the sections controller
// acceptance test, by fetching the section by ID. This is the second part of
// the 3-part read test.
func testAccSectionsCRUDReadByID(t *testing.T, sess *session.Session, s Section) {
	c := NewController(sess)

	out, err := c.GetSectionByID(s.ID)
	if err != nil {
		t.Fatalf("Can't find section by ID: %s", err)
	}

	if !reflect.DeepEqual(s, out) {
		t.Fatalf("ReadByID: Expected %#v, got %#v", s, out)
	}
}

// testAccSectionsCRUDReadByList tests the read part of the sections controller
// acceptance test, by fetching the section by searching for it in the sections
// listing. This is the third part of the 3-part read test.
func testAccSectionsCRUDReadByList(t *testing.T, sess *session.Session, s Section) {
	c := NewController(sess)

	out, err := c.ListSections()
	if err != nil {
		t.Fatalf("Can't list sections: %s", err)
	}
	for _, v := range out {
		if reflect.DeepEqual(s, v) {
			return
		}
	}

	t.Fatalf("ReadByList: Could not find section %#v in %#v", s, out)
}

// testAccSectionsCRUDUpdate tests the update part of the sections controller
// acceptance test.
func testAccSectionsCRUDUpdate(t *testing.T, sess *session.Session, s Section) {
	c := NewController(sess)

	if err := c.UpdateSection(s); err != nil {
		t.Fatalf("Error updating section: %s", err)
	}

	// Assert update
	out, err := c.GetSectionByID(s.ID)

	if err != nil {
		t.Fatalf("Error fetching section after update: %s", err)
	}

	// Update updated date in original
	s.EditDate = out.EditDate

	if !reflect.DeepEqual(s, out) {
		t.Fatalf("Error after update: expected %#v, got %#v", s, out)
	}
}

// testAccSectionsCRUDDelete tests the delete part of the sections controller
// acceptance test.
func testAccSectionsCRUDDelete(t *testing.T, sess *session.Session, s Section) {
	c := NewController(sess)

	if err := c.DeleteSection(s.ID); err != nil {
		t.Fatalf("Error deleting section: %s", err)
	}

	// check to see if section is actually gone
	if _, err := c.GetSectionByID(s.ID); err == nil {
		t.Fatalf("Section still present after delete")
	}
}

// TestAccSectionsCRUD runs a full create-read-update-delete test for a PHPIPAM
// section.
func TestAccSectionsCRUD(t *testing.T) {
	testacc.VetAccConditions(t)

	sess := session.NewSession()
	section := testCreateSectionInput
	testAccSectionsCRUDCreate(t, sess, section)
	section.ID = testAccSectionsCRUDReadByName(t, sess, section)
	testAccSectionsCRUDReadByID(t, sess, section)
	testAccSectionsCRUDReadByList(t, sess, section)
	section.Name = "bazboop"
	testAccSectionsCRUDUpdate(t, sess, section)
	testAccSectionsCRUDDelete(t, sess, section)
}

// TestAccGetSubnetsInSection tests GetSubnetsInSection against a live PHPIPAM
// instance.
func TestAccGetSubnetsInSection(t *testing.T) {
	testacc.VetAccConditions(t)

	sess := session.NewSession()
	client := NewController(sess)

	expected := testGetSubnetsInSectionExpected
	if os.Getenv("TESTACC_CUSTOM_NESTED") != "" {
		for n := range expected {
			expected[n].CustomFields = map[string]interface{}{
				"CustomTestSubnets":  nil,
				"CustomTestSubnets2": nil,
			}
		}
	} else {
		log.Println("Note: Not testing nested custom fields as TESTACC_CUSTOM_NESTED is not set")
	}
	actual, err := client.GetSubnetsInSection(1)
	if err != nil {
		t.Fatalf("Bad: %s", err)
	}

	if !reflect.DeepEqual(expected, actual) {
		t.Fatalf("Expected %s, got %s", spew.Sdump(expected), spew.Sdump(actual))
	}
}
