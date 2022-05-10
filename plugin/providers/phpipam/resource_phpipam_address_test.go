package phpipam

import (
	"errors"
	"fmt"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

const testAccResourcePHPIPAMAddressName = "phpipam_address.address"
const testAccResourcePHPIPAMAddressCIDR = "10.10.1.10"
const testAccResourcePHPIPAMAddressConfig = `
resource "phpipam_address" "address" {
  subnet_id   = 3
  ip_address  = "10.10.1.10"
  description = "Terraform test address"
  hostname    = "tf-test.cust1.local"
}
`

const testAccResourcePHPIPAMOptionalAddressConfig = `
resource "phpipam_address" "address" {
  subnet_id   = 3
  description = "Terraform test address"
  hostname    = "tf-test.cust1.local"
}
`

const testAccResourcePHPIPAMAddressCustomFieldConfig = `
resource "phpipam_address" "address" {
  subnet_id   = 3
  ip_address  = "10.10.1.10"
  description = "Terraform test address (custom fields)"
  hostname    = "tf-test.cust1.local"

  custom_fields = {
    CustomTestAddresses = "terraform-test"
  }
}
`

const testAccResourcePHPIPAMAddressCustomFieldUpdateConfig = `
resource "phpipam_address" "address" {
  subnet_id   = 3
  ip_address  = "10.10.1.10"
  description = "Terraform test address (custom fields), step 2"
  hostname    = "tf-test.cust1.local"
}
`

func TestAccResourcePHPIPAMAddress(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckResourcePHPIPAMAddressDeleted,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccResourcePHPIPAMAddressConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckResourcePHPIPAMAddressCreated,
					resource.TestCheckResourceAttr("phpipam_address.address", "subnet_id", "3"),
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
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckResourcePHPIPAMAddressDeleted,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccResourcePHPIPAMOptionalAddressConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckResourcePHPIPAMAddressCreated,
					resource.TestCheckResourceAttr("phpipam_address.address", "subnet_id", "3"),
					resource.TestCheckNoResourceAttr("phpipam_address.address", "ip_address"),
					resource.TestCheckResourceAttr("phpipam_address.address", "description", "Terraform test address"),
					resource.TestCheckResourceAttr("phpipam_address.address", "hostname", "tf-test.cust1.local"),
				),
			},
		},
	})
}

func TestAccResourcePHPIPAMAddress_CustomFields(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckResourcePHPIPAMAddressDeleted,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccResourcePHPIPAMAddressCustomFieldConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckResourcePHPIPAMAddressCreated,
					resource.TestCheckResourceAttr("phpipam_address.address", "subnet_id", "3"),
					resource.TestCheckResourceAttr("phpipam_address.address", "ip_address", "10.10.1.10"),
					resource.TestCheckResourceAttr("phpipam_address.address", "description", "Terraform test address (custom fields)"),
					resource.TestCheckResourceAttr("phpipam_address.address", "hostname", "tf-test.cust1.local"),
					resource.TestCheckResourceAttr("phpipam_address.address", "custom_fields.CustomTestAddresses", "terraform-test"),
				),
			},
			resource.TestStep{
				Config: testAccResourcePHPIPAMAddressCustomFieldUpdateConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckResourcePHPIPAMAddressCreated,
					resource.TestCheckResourceAttr("phpipam_address.address", "subnet_id", "3"),
					resource.TestCheckResourceAttr("phpipam_address.address", "ip_address", "10.10.1.10"),
					resource.TestCheckResourceAttr("phpipam_address.address", "description", "Terraform test address (custom fields), step 2"),
					resource.TestCheckResourceAttr("phpipam_address.address", "hostname", "tf-test.cust1.local"),
					resource.TestCheckNoResourceAttr("phpipam_address.address", "custom_fields.CustomTestAddresses"),
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
