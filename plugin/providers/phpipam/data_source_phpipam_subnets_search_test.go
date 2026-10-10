package phpipam

import (
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/pavel-z1/phpipam-sdk-go/controllers/search"
	"github.com/pavel-z1/phpipam-sdk-go/controllers/subnets"
	"github.com/pavel-z1/phpipam-sdk-go/phpipam"
	"github.com/pavel-z1/phpipam-sdk-go/phpipam/session"
)

// testSubnetSearchResponse is a search controller response with three
// subnets in two sections. Subnet addresses are returned as integers, as
// phpIPAM does for this endpoint.
const testSubnetSearchResponse = `{
  "code": 200,
  "success": true,
  "data": {
    "subnets": {
      "code": 200,
      "data": [
        {"id": "3", "subnet": "168099840", "mask": "24", "sectionId": "1", "description": "Customer 1"},
        {"id": "4", "subnet": "168100096", "mask": "24", "sectionId": "1", "description": "Customer 2"},
        {"id": "7", "subnet": "168100352", "mask": "24", "sectionId": "2", "description": "Customer 1"}
      ]
    }
  }
}`

const testSubnetSearchEmptyResponse = `{
  "code": 200,
  "success": true,
  "data": {"subnets": {"code": 404, "data": "No subnets found"}}
}`

// testSubnetSearchMeta returns provider client meta pointing at a test server
// that answers every request with body, and records the request URIs.
func testSubnetSearchMeta(t *testing.T, body string, uris *[]string) *ProviderPHPIPAMClient {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		*uris = append(*uris, r.URL.RequestURI())
		w.Header().Set("Content-Type", "application/json")
		http.Error(w, body, http.StatusOK)
	}))
	t.Cleanup(ts.Close)

	sess := session.NewSession(phpipam.Config{
		AppID:    "tf",
		Endpoint: ts.URL,
		Username: "nobody",
		Password: "changeit",
	})
	sess.Token.String = "token"
	return &ProviderPHPIPAMClient{
		searchController:  search.NewController(sess),
		subnetsController: subnets.NewController(sess),
	}
}

func TestDataSourcePHPIPAMSubnetsSearch(t *testing.T) {
	cases := []struct {
		name     string
		body     string
		raw      map[string]interface{}
		expected []interface{}
	}{
		{
			name:     "search only",
			body:     testSubnetSearchResponse,
			raw:      map[string]interface{}{"search": "Customer"},
			expected: []interface{}{3, 4, 7},
		},
		{
			name:     "search in section",
			body:     testSubnetSearchResponse,
			raw:      map[string]interface{}{"search": "Customer", "section_id": 1},
			expected: []interface{}{3, 4},
		},
		{
			name:     "search with description",
			body:     testSubnetSearchResponse,
			raw:      map[string]interface{}{"search": "Customer", "description": "Customer 1"},
			expected: []interface{}{3, 7},
		},
		{
			name:     "search with description_match in section",
			body:     testSubnetSearchResponse,
			raw:      map[string]interface{}{"search": "Customer", "section_id": 1, "description_match": "2$"},
			expected: []interface{}{4},
		},
		{
			name:     "no results",
			body:     testSubnetSearchEmptyResponse,
			raw:      map[string]interface{}{"search": "nothing"},
			expected: []interface{}{},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			var uris []string
			meta := testSubnetSearchMeta(t, tc.body, &uris)
			d := schema.TestResourceDataRaw(t, dataSourcePHPIPAMSubnets().Schema, tc.raw)

			if err := dataSourcePHPIPAMSubnetsRead(d, meta); err != nil {
				t.Fatalf("Bad: %s", err)
			}
			if actual := d.Get("subnet_ids").([]interface{}); !reflect.DeepEqual(tc.expected, actual) {
				t.Fatalf("Expected subnet_ids %#v, got %#v", tc.expected, actual)
			}
			if len(uris) != 1 || !strings.HasPrefix(uris[0], "/tf/search/") {
				t.Fatalf("Expected a single request to the search controller, got %#v", uris)
			}
		})
	}
}

func TestValidateSubnetSearch(t *testing.T) {
	for _, v := range []string{"Customer 1", "10.10.3.0", "2001:db8::"} {
		if _, errs := validateSubnetSearch(v, "search"); len(errs) != 0 {
			t.Fatalf("%q: unexpected errors: %v", v, errs)
		}
	}
	for _, v := range []string{"10.10.3.0/24", "a?b", "a#b"} {
		if _, errs := validateSubnetSearch(v, "search"); len(errs) == 0 {
			t.Fatalf("%q: expected an error, got none", v)
		}
	}
}

func TestDataSourcePHPIPAMSubnetsRequiresSectionOrSearch(t *testing.T) {
	r := dataSourcePHPIPAMSubnets()
	diags := r.Validate(terraform.NewResourceConfigRaw(map[string]interface{}{"description": "x"}))
	if !diags.HasError() {
		t.Fatalf("Expected an error when neither section_id nor search is set")
	}
}
