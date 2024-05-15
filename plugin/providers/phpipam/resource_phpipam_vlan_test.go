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

const testAccResourcePHPIPAMVLANName = "phpipam_vlan.vlan"
const testAccResourcePHPIPAMVLANNumber = 1999
const testAccResourceVlanPHPIPAML2DomainName = "tf-test-l2domain"

var testAccResourceVlanPHPIPAML2DomainID int

const testAccResourcePHPIPAMVLANConfig = `
resource "phpipam_l2domain" "tf-test-l2domain" {
  name        = "tf-test-l2domain"
  description = "Terraform test l2domain"
}

resource "phpipam_vlan" "vlan" {
  name         = "terraform"
  number       = 1999
  description  = "Terraform test vlan"
  l2_domain_id = phpipam_l2domain.tf-test-l2domain.domain_id
}
`

const testAccResourcePHPIPAMVLANCustomFieldConfig = `
resource "phpipam_l2domain" "tf-test-l2domain" {
  name        = "tf-test-l2domain"
  description = "Terraform test l2domain"
}

resource "phpipam_vlan" "vlan" {
  name        = "terraform"
  number      = 1999
  description = "Terraform test vlan (custom field)"
  l2_domain_id = phpipam_l2domain.tf-test-l2domain.domain_id

  custom_fields = {
    custom_CustomTestVLANs = "terraform-test"
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
		PreCheck: func() {
			testAccPreCheck(t)
			l2domainSweep("tf-test-l2domain", t)
		},
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

/*
	func TestAccResourcePHPIPAMVLAN_CustomFields(t *testing.T) {
		resource.Test(t, resource.TestCase{
	                PreCheck:  func() {
	                        testAccPreCheck(t)
	                        l2domainSweep("tf-test-l2domain", t)
	                        },
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
						resource.TestCheckResourceAttr("phpipam_vlan.vlan", "custom_fields.custom_CustomTestVLANs", "terraform-test"),
					),
				},
				resource.TestStep{
					Config: testAccResourcePHPIPAMVLANCustomFieldUpdateConfig,
					Check: resource.ComposeTestCheckFunc(
						testAccCheckResourcePHPIPAMVLANCreated,
						resource.TestCheckResourceAttr("phpipam_vlan.vlan", "name", "terraform"),
						resource.TestCheckResourceAttr("phpipam_vlan.vlan", "number", "1999"),
						resource.TestCheckResourceAttr("phpipam_vlan.vlan", "description", "Terraform test vlan (custom field), step 2"),
						resource.TestCheckNoResourceAttr("phpipam_vlan.vlan", "custom_fields.custom_CustomTestVLANs"),
					),
				},
			},
		})
	}
*/
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

	l2domainController := testAccProvider.Meta().(*ProviderPHPIPAMClient).l2domainsController
	l2domains, err := l2domainController.GetL2DomainByName(testAccResourceVlanPHPIPAML2DomainName)
	switch {
	case err != nil && strings.Contains(err.Error(), "Error from API (404): Not Found"):
		return fmt.Errorf("Expected 404, got %s", err)
	case err != nil:
		return fmt.Errorf("Error: %s", err)
	case err == nil && len(l2domains) != 1:
		return fmt.Errorf("Multiple l2 domains with the same name: %s", testAccResourceVlanPHPIPAML2DomainName)
	}

	testAccResourceVlanPHPIPAML2DomainID = l2domains[0].ID

	return nil
}

func testAccCheckResourcePHPIPAMVLANDeleted(s *terraform.State) error {
	c := testAccProvider.Meta().(*ProviderPHPIPAMClient).vlansController
	_, err := c.GetVLANsByNumberAndDomainID(testAccResourcePHPIPAMVLANNumber, testAccResourceVlanPHPIPAML2DomainID)
	error_messages := linearSearchSlice{"Error from API (404): No results (filter applied)", "Error from API (404): Vlans not found"}
	switch {
	case err == nil:
		return errors.New("Expected error, got none")
	case err != nil && !error_messages.Has(err.Error()):
		return fmt.Errorf("Expected 404, got %s", err)
	}

	return nil
}
