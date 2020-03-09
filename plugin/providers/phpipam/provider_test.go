package phpipam

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

var testAccProvider *schema.Provider
var testAccProviders map[string]terraform.ResourceProvider

const envErrMsg = `PHPIPAM_APP_ID, PHPIPAM_ENDPOINT_ADDR, PHPIPAM_PASSWORD, and PHPIPAM_USER_NAME must be set for acceptance tests`

func init() {
	testAccProvider = Provider().(*schema.Provider)
	testAccProviders = map[string]terraform.ResourceProvider{
		"phpipam": testAccProvider,
	}
}

func TestProvider(t *testing.T) {
	if err := Provider().(*schema.Provider).InternalValidate(); err != nil {
		t.Fatalf("err: %s", err)
	}
}

func TestProvider_impl(t *testing.T) {
	var _ terraform.ResourceProvider = Provider()
}

func testAccPreCheck(t *testing.T) {
	switch {
	case os.Getenv("PHPIPAM_APP_ID") == "":
		t.Fatal(envErrMsg)
	case os.Getenv("PHPIPAM_ENDPOINT_ADDR") == "":
		t.Fatal(envErrMsg)
	case os.Getenv("PHPIPAM_PASSWORD") == "":
		t.Fatal(envErrMsg)
	case os.Getenv("PHPIPAM_USER_NAME") == "":
		t.Fatal(envErrMsg)
	}
}
