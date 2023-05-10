package phpipam

import (
	"errors"
	"fmt"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

const testAccResourcePHPIPAMAddressName = "phpipam_address.address"
const testAccResourcePHPIPAMAddressCIDR = "10.10.1.10"
const testAccResourcePHPIPAMAddressConfig = `
resource "phpipam_section" "test" {
  name        = "tf-test"
  description = "Terraform test section"
}

resource "phpipam_subnet" "subnet" {
  section_id     = phpipam_section.test.section_id
  subnet_address = "10.10.1.0"
  subnet_mask    = 24
  depends_on = [phpipam_section.test]
}

resource "phpipam_address" "address" {
  subnet_id   = phpipam_subnet.subnet.subnet_id
  ip_address  = "10.10.1.10"
  description = "Terraform test address"
  hostname    = "tf-test.cust1.local"

  depends_on = [phpipam_subnet.subnet]
}
`

const testAccResourcePHPIPAMOptionalAddressConfig = `
resource "phpipam_section" "test" {
  name        = "tf-test"
  description = "Terraform test section"
}

resource "phpipam_subnet" "subnet" {
  section_id     = phpipam_section.test.section_id
  subnet_address = "10.10.1.0"
  subnet_mask    = 24
  depends_on = [phpipam_section.test]
}

resource "phpipam_address" "address" {
  subnet_id   = phpipam_subnet.subnet.subnet_id
  description = "Terraform test address"
  hostname    = "tf-test.cust1.local"
}
`

const testAccResourcePHPIPAMAddressCustomFieldConfig = `
resource "phpipam_section" "test" {
  name        = "tf-test"
  description = "Terraform test section"
}

resource "phpipam_subnet" "subnet" {
  section_id     = phpipam_section.test.section_id
  subnet_address = "10.10.1.0"
  subnet_mask    = 24
  depends_on = [phpipam_section.test]
}

resource "phpipam_address" "address" {
  subnet_id   = phpipam_subnet.subnet.subnet_id
  ip_address  = "10.10.1.10"
  description = "Terraform test address (custom fields)"
  hostname    = "tf-test.cust1.local"

  custom_fields = {
    custom_CustomTestAddresses = "terraform-test"
  }
}
`

const testAccResourcePHPIPAMAddressCustomFieldUpdateConfig = `
resource "phpipam_section" "test" {
  name        = "tf-test"
  description = "Terraform test section"
}

resource "phpipam_subnet" "subnet" {
  section_id     = phpipam_section.test.section_id
  subnet_address = "10.10.1.0"
  subnet_mask    = 24
  depends_on = [phpipam_section.test]
}

resource "phpipam_address" "address" {
  subnet_id   = phpipam_subnet.subnet.subnet_id
  ip_address  = "10.10.1.10"
  description = "Terraform test address (custom fields), step 2"
  hostname    = "tf-test.cust1.local"

  depends_on = [phpipam_subnet.subnet]
}
`

func TestAccResourcePHPIPAMAddress(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			sectionSweep("tf-test", t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckResourcePHPIPAMAddressDeleted,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccResourcePHPIPAMAddressConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckResourcePHPIPAMAddressCreated,
					resource.TestCheckResourceAttr("phpipam_address.address", "ip_address", "10.10.1.10"),
					resource.TestCheckResourceAttr("phpipam_address.address", "description", "Terraform test address"),
					resource.TestCheckResourceAttr("phpipam_address.address", "hostname", "tf-test.cust1.local"),
				),
			},
		},
	})
}

func TestAccResourcePHPIPAMOptionalAddress(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			sectionSweep("tf-test", t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckResourcePHPIPAMAddressDeleted,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccResourcePHPIPAMOptionalAddressConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckResourcePHPIPAMAddressCreated,
					resource.TestCheckResourceAttr("phpipam_address.address", "description", "Terraform test address"),
					resource.TestCheckResourceAttr("phpipam_address.address", "hostname", "tf-test.cust1.local"),
				),
			},
		},
	})
}

func TestAccResourcePHPIPAMAddress_CustomFields(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			sectionSweep("tf-test", t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckResourcePHPIPAMAddressDeleted,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccResourcePHPIPAMAddressCustomFieldConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckResourcePHPIPAMAddressCreated,
					resource.TestCheckResourceAttr("phpipam_address.address", "ip_address", "10.10.1.10"),
					resource.TestCheckResourceAttr("phpipam_address.address", "description", "Terraform test address (custom fields)"),
					resource.TestCheckResourceAttr("phpipam_address.address", "hostname", "tf-test.cust1.local"),
					resource.TestCheckResourceAttr("phpipam_address.address", "custom_fields.custom_CustomTestAddresses", "terraform-test"),
				),
			},
			resource.TestStep{
				Config: testAccResourcePHPIPAMAddressCustomFieldUpdateConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckResourcePHPIPAMAddressCreated,
					resource.TestCheckResourceAttr("phpipam_address.address", "ip_address", "10.10.1.10"),
					resource.TestCheckResourceAttr("phpipam_address.address", "description", "Terraform test address (custom fields), step 2"),
					resource.TestCheckResourceAttr("phpipam_address.address", "hostname", "tf-test.cust1.local"),
					resource.TestCheckNoResourceAttr("phpipam_address.address", "custom_fields.custom_CustomTestAddresses"),
				),
			},
		},
	})
}

func testAccCheckResourcePHPIPAMAddressCreated(s *terraform.State) error {
	r, ok := s.RootModule().Resources[testAccResourcePHPIPAMAddressName]
	if !ok {
		return fmt.Errorf("Resource name %s could not be found", testAccResourcePHPIPAMAddressName)
	}
	if r.Primary.ID == "" {
		return errors.New("No ID is set")
	}

	id, _ := strconv.Atoi(r.Primary.ID)

	c := testAccProvider.Meta().(*ProviderPHPIPAMClient).addressesController
	if _, err := c.GetAddressByID(id); err != nil {
		return err
	}
	return nil
}

func testAccCheckResourcePHPIPAMAddressDeleted(s *terraform.State) error {
	c := testAccProvider.Meta().(*ProviderPHPIPAMClient).addressesController
	_, err := c.GetAddressesByIP(testAccResourcePHPIPAMAddressCIDR)
	switch {
	case err == nil:
		return errors.New("Expected error, got none")
	case err != nil && err.Error() != "Error from API (404): Address not found":
		return fmt.Errorf("Expected 404, got %s", err)
	}

	return nil
}
