package phpipam

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

const testAccDataSourcePHPIPAMFirstFreeSubnetConfig = `
resource "phpipam_section" "section" {
  name        = "tf-test"
  description = "Terraform test section"
}

resource "phpipam_subnet" "subnet" {
  section_id     = phpipam_section.section.section_id
  subnet_address = "10.10.1.0"
  subnet_mask    = 24
}

data "phpipam_subnet" "subnet_by_cidr" {
  section_id     = phpipam_section.section.section_id
  subnet_address = "10.10.1.0"
  subnet_mask = 24
  depends_on = [phpipam_subnet.subnet]
}

data "phpipam_first_free_subnet" "next" {
  subnet_id = data.phpipam_subnet.subnet_by_cidr.subnet_id
  subnet_mask = 26
  depends_on = [data.phpipam_subnet.subnet_by_cidr]
}
`

const testAccDataSourcePHPIPAMFirstFreeSubnetNoFreeConfig = `
resource "phpipam_section" "section" {
  name        = "tf-test"
  description = "Terraform test section"
}

resource "phpipam_subnet" "subnet1"  {
  section_id     = phpipam_section.section.section_id
  subnet_address = "10.10.3.0"
  subnet_mask    = 31
}

data "phpipam_subnet" "subnet_by_cidr" {
  section_id     = phpipam_section.section.section_id
  subnet_address = "10.10.3.0"
  subnet_mask    = 31
  depends_on     = [phpipam_subnet.subnet1]
}

data "phpipam_first_free_subnet" "next" {
  subnet_id   = data.phpipam_subnet.subnet_by_cidr.subnet_id
  subnet_mask = 30

  depends_on  = [phpipam_subnet.subnet1]
}

resource "phpipam_subnet" "subnet2"  {
  section_id     = phpipam_section.section.section_id
  subnet_address = cidrhost(data.phpipam_first_free_subnet.next.ip_address, data.phpipam_first_free_subnet.next.subnet_mask)
  subnet_mask    = data.phpipam_first_free_subnet.next.subnet_mask
}
`

func TestAccDataSourcePHPIPAMFirstFreeSubnet(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			sectionSweep("tf-test", t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccDataSourcePHPIPAMFirstFreeSubnetConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.phpipam_first_free_subnet.next", "ip_address", "10.10.1.0/26"),
				),
			},
		},
	})
}

func TestAccDataSourcePHPIPAMFirstFreeSubnetNoFree(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			sectionSweep("tf-test", t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config:      testAccDataSourcePHPIPAMFirstFreeSubnetNoFreeConfig,
				ExpectError: regexp.MustCompile("No subnets found"),
			},
		},
	})
}
