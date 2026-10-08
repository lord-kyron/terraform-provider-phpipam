package phpipam

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

const testAccResourcePHPIPAMFirstFreeAddressConfig = `
resource "phpipam_section" "section" {
  name        = "tf-test"
  description = "Terraform test section"
}

resource "phpipam_subnet" "subnet" {
  section_id     = phpipam_section.section.section_id
  subnet_address = "10.10.8.0"
  subnet_mask    = 24
}

resource "phpipam_first_free_address" "next" {
  subnet_id   = phpipam_subnet.subnet.subnet_id
  hostname    = "tf-test.cust1.local"
  description = "Terraform test address"

  depends_on = [phpipam_subnet.subnet]
}
`

const testAccResourcePHPIPAMFirstFreeAddressSubnetIDsConfig = `
resource "phpipam_section" "section" {
  name        = "tf-test"
  description = "Terraform test section"
}

resource "phpipam_subnet" "full" {
  section_id     = phpipam_section.section.section_id
  subnet_address = "10.10.9.0"
  subnet_mask    = 30
}

resource "phpipam_address" "full_1" {
  subnet_id  = phpipam_subnet.full.subnet_id
  ip_address = "10.10.9.1"
}

resource "phpipam_address" "full_2" {
  subnet_id  = phpipam_subnet.full.subnet_id
  ip_address = "10.10.9.2"
}

resource "phpipam_subnet" "free" {
  section_id     = phpipam_section.section.section_id
  subnet_address = "10.10.10.0"
  subnet_mask    = 24
}

resource "phpipam_first_free_address" "next" {
  subnet_ids  = [phpipam_subnet.full.subnet_id, phpipam_subnet.free.subnet_id]
  hostname    = "tf-test.cust1.local"
  description = "Terraform test address"

  depends_on = [
    "phpipam_address.full_1",
    "phpipam_address.full_2",
    "phpipam_subnet.free",
  ]
}
`

func TestAccResourcePHPIPAMFirstFreeAddress(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			sectionSweep("tf-test", t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccResourcePHPIPAMFirstFreeAddressConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("phpipam_first_free_address.next", "ip_address", "10.10.8.1"),
					resource.TestCheckResourceAttrPair("phpipam_first_free_address.next", "subnet_id", "phpipam_subnet.subnet", "subnet_id"),
				),
			},
		},
	})
}

func TestAccResourcePHPIPAMFirstFreeAddressSubnetIDs(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			sectionSweep("tf-test", t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccResourcePHPIPAMFirstFreeAddressSubnetIDsConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("phpipam_first_free_address.next", "ip_address", "10.10.10.1"),
					resource.TestCheckResourceAttrPair("phpipam_first_free_address.next", "subnet_id", "phpipam_subnet.free", "subnet_id"),
				),
			},
		},
	})
}

const testAccResourcePHPIPAMFirstFreeAddressMigrateStep1Config = `
resource "phpipam_section" "section" {
  name        = "tf-test"
  description = "Terraform test section"
}

resource "phpipam_subnet" "a" {
  section_id     = phpipam_section.section.section_id
  subnet_address = "10.10.15.0"
  subnet_mask    = 24
}

// Created up-front (but unused by the address yet) so that its subnet_id is
// already known by the time step 2 switches to subnet_ids - otherwise the
// still-to-be-created subnet keeps the whole subnet_ids list "unknown" at
// plan time and the CustomizeDiff can't tell it's a safe, no-op migration.
resource "phpipam_subnet" "b" {
  section_id     = phpipam_section.section.section_id
  subnet_address = "10.10.16.0"
  subnet_mask    = 24
}

resource "phpipam_first_free_address" "next" {
  subnet_id   = phpipam_subnet.a.subnet_id
  hostname    = "tf-test.cust1.local"
  description = "Terraform test address"

  depends_on = [phpipam_subnet.a, phpipam_subnet.b]
}
`

const testAccResourcePHPIPAMFirstFreeAddressMigrateStep2Config = `
resource "phpipam_section" "section" {
  name        = "tf-test"
  description = "Terraform test section"
}

resource "phpipam_subnet" "a" {
  section_id     = phpipam_section.section.section_id
  subnet_address = "10.10.15.0"
  subnet_mask    = 24
}

resource "phpipam_subnet" "b" {
  section_id     = phpipam_section.section.section_id
  subnet_address = "10.10.16.0"
  subnet_mask    = 24
}

resource "phpipam_first_free_address" "next" {
  subnet_ids  = [phpipam_subnet.a.subnet_id, phpipam_subnet.b.subnet_id]
  hostname    = "tf-test.cust1.local"
  description = "Terraform test address"

  depends_on = [phpipam_subnet.a, phpipam_subnet.b]
}
`

// TestAccResourcePHPIPAMFirstFreeAddressMigrateToSubnetIDs ensures that
// switching a resource's config from subnet_id to a subnet_ids list that
// still contains the current subnet is a no-op (no replacement), since the
// address does not need to move. See resourcePHPIPAMFirstFreeAddressCustomizeDiff.
func TestAccResourcePHPIPAMFirstFreeAddressMigrateToSubnetIDs(t *testing.T) {
	var firstID string

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			sectionSweep("tf-test", t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccResourcePHPIPAMFirstFreeAddressMigrateStep1Config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("phpipam_first_free_address.next", "ip_address", "10.10.15.1"),
					func(s *terraform.State) error {
						r, ok := s.RootModule().Resources["phpipam_first_free_address.next"]
						if !ok {
							return fmt.Errorf("resource not found in state")
						}
						firstID = r.Primary.ID
						return nil
					},
				),
			},
			resource.TestStep{
				Config: testAccResourcePHPIPAMFirstFreeAddressMigrateStep2Config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("phpipam_first_free_address.next", "ip_address", "10.10.15.1"),
					resource.TestCheckResourceAttrPair("phpipam_first_free_address.next", "subnet_id", "phpipam_subnet.a", "subnet_id"),
					func(s *terraform.State) error {
						r, ok := s.RootModule().Resources["phpipam_first_free_address.next"]
						if !ok {
							return fmt.Errorf("resource not found in state")
						}
						if r.Primary.ID != firstID {
							return fmt.Errorf("expected address to be preserved (id %s), got a new one (%s) - resource was replaced", firstID, r.Primary.ID)
						}
						return nil
					},
				),
			},
		},
	})
}

