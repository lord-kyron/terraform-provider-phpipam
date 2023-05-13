package phpipam

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

const testAccDataSourcePHPIPAMSubnetConfig = `
resource "phpipam_section" "section" {
  name        = "tf-test"
  description = "Terraform test section"
}

resource "phpipam_subnet" "subnet" {
  subnet_address = "10.10.3.0"
  subnet_mask    = 24
  description    = "Customer 1"
  section_id     = phpipam_section.section.section_id
  depends_on     = [phpipam_section.section]
}

resource "phpipam_subnet" "subnet2" {
  subnet_address = "10.10.2.0"
  subnet_mask    = 25
  description    = "Customer 2"
  section_id     = phpipam_section.section.section_id
  depends_on     = [phpipam_section.section]
}

data "phpipam_subnet" "subnet_by_cidr" {
  section_id     = phpipam_section.section.section_id
  subnet_address = "10.10.3.0"
  subnet_mask    = 24
  depends_on     = [phpipam_subnet.subnet] 
}

data "phpipam_subnet" "subnet_by_id" {
  subnet_id  = data.phpipam_subnet.subnet_by_cidr.subnet_id
  depends_on = [data.phpipam_subnet.subnet_by_cidr]
}

data "phpipam_subnet" "subnet_by_description" {
  section_id  = phpipam_section.section.section_id
  description = "Customer 1"
  depends_on  = [phpipam_subnet.subnet]
}

data "phpipam_subnet" "subnet_by_description_match" {
  section_id  = phpipam_section.section.section_id
  description_match = "[Cc]ustomer 2"
  depends_on  = [phpipam_subnet.subnet2]
}
`

const testAccDataSourcePHPIPAMSubnetCustomFieldConfig = `
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
    custom_CustomTestSubnets2 = "terraform2-test"
  }
}

data "phpipam_subnet" "custom_search" {
  section_id = phpipam_subnet.subnet.section_id

  custom_field_filter = {
    custom_CustomTestSubnets = ".*terraform.*"
    custom_CustomTestSubnets2 = ".*terraform2.*"
  }
}
`

func TestAccDataSourcePHPIPAMSubnet(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			sectionSweep("tf-test", t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccDataSourcePHPIPAMSubnetConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair("data.phpipam_subnet.subnet_by_cidr", "subnet_id", "data.phpipam_subnet.subnet_by_id", "subnet_id"),
					resource.TestCheckResourceAttrPair("data.phpipam_subnet.subnet_by_cidr", "subnet_address", "data.phpipam_subnet.subnet_by_id", "subnet_address"),
					resource.TestCheckResourceAttrPair("data.phpipam_subnet.subnet_by_cidr", "subnet_mask", "data.phpipam_subnet.subnet_by_id", "subnet_mask"),
					resource.TestCheckResourceAttr("data.phpipam_subnet.subnet_by_description", "subnet_address", "10.10.3.0"),
					resource.TestCheckResourceAttr("data.phpipam_subnet.subnet_by_description", "subnet_mask", "24"),
					resource.TestCheckResourceAttr("data.phpipam_subnet.subnet_by_description_match", "subnet_address", "10.10.2.0"),
				),
			},
		},
	})
}

func TestAccDataSourcePHPIPAMSubnet_CustomFields(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			sectionSweep("tf-test", t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccDataSourcePHPIPAMSubnetCustomFieldConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.phpipam_subnet.custom_search", "subnet_address", "10.10.3.0"),
					resource.TestCheckResourceAttr("data.phpipam_subnet.custom_search", "subnet_mask", "24"),
					resource.TestCheckResourceAttr("data.phpipam_subnet.custom_search", "description", "Terraform test subnet (custom fields)"),
					resource.TestCheckResourceAttr("data.phpipam_subnet.custom_search", "custom_fields.custom_CustomTestSubnets", "terraform-test"),
				),
			},
		},
	})
}
