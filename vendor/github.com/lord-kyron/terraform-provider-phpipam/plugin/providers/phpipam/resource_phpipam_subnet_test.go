package phpipam

import (
	"errors"
	"fmt"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

const testAccResourcePHPIPAMSubnetName = "phpipam_subnet.subnet"
const testAccResourcePHPIPAMSubnetCIDR = "10.10.3.0/24"
const testAccResourcePHPIPAMSubnetConfig = `
resource "phpipam_subnet" "subnet" {
	subnet_address = "10.10.3.0"
	subnet_mask = 24
	description = "Terraform test subnet"
	section_id = 1
}
`

const testAccResourcePHPIPAMSubnetCustomFieldConfig = `
resource "phpipam_subnet" "subnet" {
  subnet_address = "10.10.3.0"
  subnet_mask    = 24
  description    = "Terraform test subnet (custom fields)"
  section_id     = 1

  custom_fields = {
    CustomTestSubnets = "terraform-test"
  }
}
`

const testAccResourcePHPIPAMSubnetCustomFieldUpdateConfig = `
resource "phpipam_subnet" "subnet" {
  subnet_address = "10.10.3.0"
  subnet_mask    = 24
  description    = "Terraform test subnet (custom fields), step 2"
  section_id     = 1
}

`

func TestAccResourcePHPIPAMSubnet(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckResourcePHPIPAMSubnetDeleted,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccResourcePHPIPAMSubnetConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckResourcePHPIPAMSubnetCreated,
					resource.TestCheckResourceAttr("phpipam_subnet.subnet", "subnet_address", "10.10.3.0"),
					resource.TestCheckResourceAttr("phpipam_subnet.subnet", "subnet_mask", "24"),
					resource.TestCheckResourceAttr("phpipam_subnet.subnet", "description", "Terraform test subnet"),
				),
			},
		},
	})
}

func TestAccResourcePHPIPAMSubnet_CustomFields(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckResourcePHPIPAMSubnetDeleted,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccResourcePHPIPAMSubnetCustomFieldConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckResourcePHPIPAMSubnetCreated,
					resource.TestCheckResourceAttr("phpipam_subnet.subnet", "subnet_address", "10.10.3.0"),
					resource.TestCheckResourceAttr("phpipam_subnet.subnet", "subnet_mask", "24"),
					resource.TestCheckResourceAttr("phpipam_subnet.subnet", "description", "Terraform test subnet (custom fields)"),
					resource.TestCheckResourceAttr("phpipam_subnet.subnet", "custom_fields.CustomTestSubnets", "terraform-test"),
				),
			},
			resource.TestStep{
				Config: testAccResourcePHPIPAMSubnetCustomFieldUpdateConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckResourcePHPIPAMSubnetCreated,
					resource.TestCheckResourceAttr("phpipam_subnet.subnet", "subnet_address", "10.10.3.0"),
					resource.TestCheckResourceAttr("phpipam_subnet.subnet", "subnet_mask", "24"),
					resource.TestCheckResourceAttr("phpipam_subnet.subnet", "description", "Terraform test subnet (custom fields), step 2"),
					resource.TestCheckNoResourceAttr("phpipam_subnet.subnet", "custom_fields.CustomTestSubnets"),
				),
			},
		},
	})
}

func testAccCheckResourcePHPIPAMSubnetCreated(s *terraform.State) error {
	r, ok := s.RootModule().Resources[testAccResourcePHPIPAMSubnetName]
	if !ok {
		return fmt.Errorf("Resource name %s could not be found", testAccResourcePHPIPAMSubnetName)
	}
	if r.Primary.ID == "" {
		return errors.New("No ID is set")
	}

	id, _ := strconv.Atoi(r.Primary.ID)

	c := testAccProvider.Meta().(*ProviderPHPIPAMClient).subnetsController
	if _, err := c.GetSubnetByID(id); err != nil {
		return err
	}
	return nil
}

func testAccCheckResourcePHPIPAMSubnetDeleted(s *terraform.State) error {
	c := testAccProvider.Meta().(*ProviderPHPIPAMClient).subnetsController
	_, err := c.GetSubnetsByCIDR(testAccResourcePHPIPAMSubnetCIDR)
	switch {
	case err == nil:
		return errors.New("Expected error, got none")
	case err != nil && err.Error() != "Error from API (404): No subnets found":
		return fmt.Errorf("Expected 404, got %s", err)
	}

	return nil
}
