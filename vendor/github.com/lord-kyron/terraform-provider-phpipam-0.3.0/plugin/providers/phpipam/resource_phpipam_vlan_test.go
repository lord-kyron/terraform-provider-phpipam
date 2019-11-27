package phpipam

import (
	"errors"
	"fmt"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

const testAccResourcePHPIPAMVLANName = "phpipam_vlan.vlan"
const testAccResourcePHPIPAMVLANNumber = 1999
const testAccResourcePHPIPAMVLANConfig = `
resource "phpipam_vlan" "vlan" {
  name        = "terraform"
  number      = 1999
  description = "Terraform test vlan"
}
`

const testAccResourcePHPIPAMVLANCustomFieldConfig = `
resource "phpipam_vlan" "vlan" {
  name        = "terraform"
  number      = 1999
  description = "Terraform test vlan (custom field)"

  custom_fields = {
    CustomTestVLANs = "terraform-test"
  }
}
`

const testAccResourcePHPIPAMVLANCustomFieldUpdateConfig = `
resource "phpipam_vlan" "vlan" {
  name        = "terraform"
  number      = 1999
  description = "Terraform test vlan (custom field), step 2"
}
`

func TestAccResourcePHPIPAMVLAN(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckResourcePHPIPAMVLANDeleted,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccResourcePHPIPAMVLANConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckResourcePHPIPAMVLANCreated,
					resource.TestCheckResourceAttr("phpipam_vlan.vlan", "name", "terraform"),
					resource.TestCheckResourceAttr("phpipam_vlan.vlan", "number", "1999"),
					resource.TestCheckResourceAttr("phpipam_vlan.vlan", "description", "Terraform test vlan"),
				),
			},
		},
	})
}

func TestAccResourcePHPIPAMVLAN_CustomFields(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckResourcePHPIPAMVLANDeleted,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccResourcePHPIPAMVLANCustomFieldConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckResourcePHPIPAMVLANCreated,
					resource.TestCheckResourceAttr("phpipam_vlan.vlan", "name", "terraform"),
					resource.TestCheckResourceAttr("phpipam_vlan.vlan", "number", "1999"),
					resource.TestCheckResourceAttr("phpipam_vlan.vlan", "description", "Terraform test vlan (custom field)"),
					resource.TestCheckResourceAttr("phpipam_vlan.vlan", "custom_fields.CustomTestVLANs", "terraform-test"),
				),
			},
			resource.TestStep{
				Config: testAccResourcePHPIPAMVLANCustomFieldUpdateConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckResourcePHPIPAMVLANCreated,
					resource.TestCheckResourceAttr("phpipam_vlan.vlan", "name", "terraform"),
					resource.TestCheckResourceAttr("phpipam_vlan.vlan", "number", "1999"),
					resource.TestCheckResourceAttr("phpipam_vlan.vlan", "description", "Terraform test vlan (custom field), step 2"),
					resource.TestCheckNoResourceAttr("phpipam_vlan.vlan", "custom_fields.CustomTestVLANs"),
				),
			},
		},
	})
}

func testAccCheckResourcePHPIPAMVLANCreated(s *terraform.State) error {
	r, ok := s.RootModule().Resources[testAccResourcePHPIPAMVLANName]
	if !ok {
		return fmt.Errorf("Resource name %s could not be found", testAccResourcePHPIPAMVLANName)
	}
	if r.Primary.ID == "" {
		return errors.New("No ID is set")
	}

	id, _ := strconv.Atoi(r.Primary.ID)

	c := testAccProvider.Meta().(*ProviderPHPIPAMClient).vlansController
	if _, err := c.GetVLANByID(id); err != nil {
		return err
	}
	return nil
}

func testAccCheckResourcePHPIPAMVLANDeleted(s *terraform.State) error {
	c := testAccProvider.Meta().(*ProviderPHPIPAMClient).vlansController
	_, err := c.GetVLANsByNumber(testAccResourcePHPIPAMVLANNumber)
	switch {
	case err == nil:
		return errors.New("Expected error, got none")
	case err != nil && err.Error() != "Error from API (404): Vlans not found":
		return fmt.Errorf("Expected 404, got %s", err)
	}

	return nil
}
