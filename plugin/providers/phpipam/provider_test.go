package phpipam

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var testAccProvider *schema.Provider
var testAccProviders map[string]*schema.Provider

const envErrMsg = `PHPIPAM_APP_ID, PHPIPAM_ENDPOINT_ADDR, PHPIPAM_PASSWORD, and PHPIPAM_USER_NAME must be set for acceptance tests`

func init() {
	testAccProvider = Provider()
	testAccProviders = map[string]*schema.Provider{
		"phpipam": testAccProvider,
	}
}

func TestProvider(t *testing.T) {
	if err := Provider().InternalValidate(); err != nil {
		t.Fatalf("err: %s", err)
	}
}

func TestProvider_impl(t *testing.T) {
	var _ *schema.Provider = Provider()
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

func testAccProviderMeta(t *testing.T) (interface{}, error) {
	t.Helper()
	d := schema.TestResourceDataRaw(t, testAccProvider.Schema, make(map[string]interface{}))
	return providerConfigure(d)
}

func sectionSweep(sectionName string, t *testing.T) error {
	meta, err := testAccProviderMeta(t)
	if err != nil {
		t.Fatalf("bad: %s", err)
	}

	c := meta.(*ProviderPHPIPAMClient).sectionsController
	section, err := c.GetSectionByName(sectionName)
	switch {
	case err != nil && err.Error() == "Error from API (404): Not Found":
		return nil
	case err != nil:
		t.Fatalf("bad: %s", err)
	}

	if err := c.DeleteSection(section.ID); err != nil {
		t.Fatalf("bad: %s", err)
	}

	return nil
}

func l2domainSweep(l2domainName string, t *testing.T) error {
	meta, err := testAccProviderMeta(t)
	if err != nil {
		t.Fatalf("bad: %s", err)
	}

	c := meta.(*ProviderPHPIPAMClient).l2domainsController
	l2domains, err := c.GetL2DomainByName(l2domainName)
	switch {
	case err != nil && err.Error() == "Error from API (404): No results (filter applied)":
		return nil
	case err != nil:
		t.Fatalf("bad: %s", err)
	case err == nil && len(l2domains) != 1:
		t.Fatalf("Multiple l2 domains with the same name: %s", l2domainName)
	}

	if err := c.DeleteL2Domain(l2domains[0].ID); err != nil {
		t.Fatalf("bad: %s", err)
	}

	return nil
}
