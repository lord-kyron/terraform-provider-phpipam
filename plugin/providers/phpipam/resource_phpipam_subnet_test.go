package phpipam

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

const testAccResourcePHPIPAMSubnetName = "phpipam_subnet.subnet"
const testAccResourcePHPIPAMSubnetCIDR = "10.10.3.0/24"
const testAccResourceSubnetPHPIPAMSectionName = "tf-test"

var testAccResourceSubnetPHPIPAMSectionID = 0

const testAccResourcePHPIPAMSubnetConfig = `
resource "phpipam_section" "section" {
  name        = "tf-test"
  description = "Terraform test section"
}

resource "phpipam_subnet" "subnet" {
	subnet_address = "10.10.3.0"
	subnet_mask = 24
	description = "Terraform test subnet"
	section_id = phpipam_section.section.section_id
}
`

const testAccResourcePHPIPAMSubnetCustomFieldConfig = `
resource "phpipam_section" "section" {
  name        = "tf-test"
  description = "Terraform test section"
}

resource "phpipam_subnet" "subnet" {
  subnet_address = "10.10.3.0"
  subnet_mask    = 24
  description    = "Terraform test subnet (custom fields)"
  section_id     = phpipam_section.section.section_id

  custom_fields = {
    custom_CustomTestSubnets = "terraform-test"
  }
}
`

const testAccResourcePHPIPAMSubnetCustomFieldUpdateConfig = `
resource "phpipam_section" "section" {
  name        = "tf-test"
  description = "Terraform test section"
}

resource "phpipam_subnet" "subnet" {
  subnet_address = "10.10.3.0"
  subnet_mask    = 24
  description    = "Terraform test subnet (custom fields), step 2"
  section_id     = phpipam_section.section.section_id

}

`

func TestAccResourcePHPIPAMSubnet(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			sectionSweep("tf-test", t)
		},
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
		PreCheck: func() {
			testAccPreCheck(t)
			sectionSweep("tf-test", t)
		},
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
					resource.TestCheckResourceAttr("phpipam_subnet.subnet", "custom_fields.custom_CustomTestSubnets", "terraform-test"),
				),
			},
			resource.TestStep{
				Config: testAccResourcePHPIPAMSubnetCustomFieldUpdateConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckResourcePHPIPAMSubnetCreated,
					resource.TestCheckResourceAttr("phpipam_subnet.subnet", "subnet_address", "10.10.3.0"),
					resource.TestCheckResourceAttr("phpipam_subnet.subnet", "subnet_mask", "24"),
					resource.TestCheckResourceAttr("phpipam_subnet.subnet", "description", "Terraform test subnet (custom fields), step 2"),
					resource.TestCheckNoResourceAttr("phpipam_subnet.subnet", "custom_fields.custom_CustomTestSubnets"),
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

	sectionController := testAccProvider.Meta().(*ProviderPHPIPAMClient).sectionsController
	section, err := sectionController.GetSectionByName(testAccResourceSubnetPHPIPAMSectionName)
	switch {
	case err != nil && strings.Contains(err.Error(), "Error from API (404): Not Found"):
		return fmt.Errorf("Expected 404, got %s", err)
	case err != nil:
		return fmt.Errorf("Error: %s", err)
	}

	testAccResourceSubnetPHPIPAMSectionID = section.ID

	return nil
}

func testAccCheckResourcePHPIPAMSubnetDeleted(s *terraform.State) error {
	subnetController := testAccProvider.Meta().(*ProviderPHPIPAMClient).subnetsController
	_, err := subnetController.GetSubnetsByCIDRAndSection(testAccResourcePHPIPAMSubnetCIDR, testAccResourceSubnetPHPIPAMSectionID)
	error_messages := linearSearchSlice{"Error from API (404): No results (filter applied)", "Error from API (404): No subnets found"}
	switch {
	case err == nil:
		return errors.New("Expected error, got none")
	case err != nil && !error_messages.Has(err.Error()):
		return fmt.Errorf("Expected 404, got %s", err)
	}

	return nil
}
