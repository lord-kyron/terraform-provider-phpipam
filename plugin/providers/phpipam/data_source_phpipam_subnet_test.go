package phpipam

import (
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
)

const testAccDataSourcePHPIPAMSubnetConfig = `
data "phpipam_subnet" "subnet_by_cidr" {
  subnet_address = "10.10.2.0"
  subnet_mask    = 24
}

data "phpipam_subnet" "subnet_by_id" {
  subnet_id = "${data.phpipam_subnet.subnet_by_cidr.subnet_id}"
}

data "phpipam_subnet" "subnet_by_description" {
  section_id  = 1
  description = "Customer 1"
}

data "phpipam_subnet" "subnet_by_description_match" {
  section_id  = 1
  description_match = "ustomer 2"
}
`

const testAccDataSourcePHPIPAMSubnetCustomFieldConfig = `
resource "phpipam_subnet" "subnet" {
  subnet_address = "10.10.3.0"
  subnet_mask    = 24
  description    = "Terraform test subnet (custom fields)"
  section_id     = 1

  custom_fields = {
    CustomTestSubnets = "terraform-test"
    CustomTestSubnets2 = "terraform2-test"
  }
}

data "phpipam_subnet" "custom_search" {
  section_id = "${phpipam_subnet.subnet.section_id}"

  custom_field_filter = {
    CustomTestSubnets = ".*terraform.*"
    CustomTestSubnets2 = ".*terraform2.*"
  }
}
`

func TestAccDataSourcePHPIPAMSubnet(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccDataSourcePHPIPAMSubnetConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair("data.phpipam_subnet.subnet_by_cidr", "subnet_id", "data.phpipam_subnet.subnet_by_id", "subnet_id"),
					resource.TestCheckResourceAttrPair("data.phpipam_subnet.subnet_by_cidr", "subnet_address", "data.phpipam_subnet.subnet_by_id", "subnet_address"),
					resource.TestCheckResourceAttrPair("data.phpipam_subnet.subnet_by_cidr", "subnet_mask", "data.phpipam_subnet.subnet_by_id", "subnet_mask"),
					resource.TestCheckResourceAttr("data.phpipam_subnet.subnet_by_description", "subnet_address", "10.10.1.0"),
					resource.TestCheckResourceAttr("data.phpipam_subnet.subnet_by_description", "subnet_mask", "24"),
					resource.TestCheckResourceAttr("data.phpipam_subnet.subnet_by_description_match", "subnet_address", "10.10.2.0"),
				),
			},
		},
	})
}

func TestAccDataSourcePHPIPAMSubnet_CustomFields(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccDataSourcePHPIPAMSubnetCustomFieldConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.phpipam_subnet.custom_search", "subnet_address", "10.10.3.0"),
					resource.TestCheckResourceAttr("data.phpipam_subnet.custom_search", "subnet_mask", "24"),
					resource.TestCheckResourceAttr("data.phpipam_subnet.custom_search", "description", "Terraform test subnet (custom fields)"),
					resource.TestCheckResourceAttr("data.phpipam_subnet.custom_search", "custom_fields.CustomTestSubnets", "terraform-test"),
				),
			},
		},
	})
}
