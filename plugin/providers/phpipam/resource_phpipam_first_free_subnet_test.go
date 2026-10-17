package phpipam

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

const testAccResourcePHPIPAMFirstFreeSubnetConfig = `
resource "phpipam_section" "section" {
  name        = "tf-test"
  description = "Terraform test section"
}

resource "phpipam_subnet" "parent" {
  section_id     = phpipam_section.section.section_id
  subnet_address = "10.10.20.0"
  subnet_mask    = 24
}

resource "phpipam_first_free_subnet" "first" {
  parent_subnet_id = phpipam_subnet.parent.subnet_id
  subnet_mask      = 26
  description      = "Terraform test first free subnet"
}

resource "phpipam_first_free_subnet" "second" {
  parent_subnet_id = phpipam_subnet.parent.subnet_id
  subnet_mask      = 26
  description      = "Terraform test second free subnet"

  depends_on = [phpipam_first_free_subnet.first]
}
`

const testAccResourcePHPIPAMFirstFreeSubnetUpdateConfig = `
resource "phpipam_section" "section" {
  name        = "tf-test"
  description = "Terraform test section"
}

resource "phpipam_subnet" "parent" {
  section_id     = phpipam_section.section.section_id
  subnet_address = "10.10.20.0"
  subnet_mask    = 24
}

resource "phpipam_first_free_subnet" "first" {
  parent_subnet_id = phpipam_subnet.parent.subnet_id
  subnet_mask      = 26
  description      = "Terraform test first free subnet, updated"

  custom_fields = {
    custom_CustomTestSubnets = "terraform-test"
  }
}

resource "phpipam_first_free_subnet" "second" {
  parent_subnet_id = phpipam_subnet.parent.subnet_id
  subnet_mask      = 26
  description      = "Terraform test second free subnet"

  depends_on = [phpipam_first_free_subnet.first]
}
`

func TestAccResourcePHPIPAMFirstFreeSubnet(t *testing.T) {
	var id string

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			sectionSweep("tf-test", t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccResourcePHPIPAMFirstFreeSubnetConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("phpipam_first_free_subnet.first", "subnet_address", "10.10.20.0"),
					resource.TestCheckResourceAttr("phpipam_first_free_subnet.first", "subnet_mask", "26"),
					resource.TestCheckResourceAttrPair("phpipam_first_free_subnet.first", "section_id", "phpipam_section.section", "section_id"),
					resource.TestCheckResourceAttr("phpipam_first_free_subnet.second", "subnet_address", "10.10.20.64"),
					resource.TestCheckResourceAttr("phpipam_first_free_subnet.second", "subnet_mask", "26"),
					testCheckResourceIDSave("phpipam_first_free_subnet.first", &id),
				),
			},
			resource.TestStep{
				Config: testAccResourcePHPIPAMFirstFreeSubnetUpdateConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("phpipam_first_free_subnet.first", "subnet_address", "10.10.20.0"),
					resource.TestCheckResourceAttr("phpipam_first_free_subnet.first", "description", "Terraform test first free subnet, updated"),
					resource.TestCheckResourceAttr("phpipam_first_free_subnet.first", "custom_fields.custom_CustomTestSubnets", "terraform-test"),
					testCheckResourceIDUnchanged("phpipam_first_free_subnet.first", &id),
				),
			},
			resource.TestStep{
				ResourceName:      "phpipam_first_free_subnet.first",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
