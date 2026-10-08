package phpipam

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

const testAccDataSourcePHPIPAMFirstFreeAddressConfig = `
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

data "phpipam_first_free_address" "next" {
  subnet_id = data.phpipam_subnet.subnet_by_cidr.subnet_id
  depends_on = [data.phpipam_subnet.subnet_by_cidr]
}
`

const testAccDataSourcePHPIPAMFirstFreeAddressNoFreeConfig = `
resource "phpipam_section" "section" {
  name        = "tf-test"
  description = "Terraform test section"
}

resource "phpipam_subnet" "subnet" {
  section_id     = phpipam_section.section.section_id
  subnet_address = "10.10.3.0"
  subnet_mask    = 30
}

resource "phpipam_address" "address_1" {
  subnet_id  = phpipam_subnet.subnet.subnet_id
  ip_address = "10.10.3.1"
}

resource "phpipam_address" "address_2" {
  subnet_id  = phpipam_subnet.subnet.subnet_id
  ip_address = "10.10.3.2"
}

data "phpipam_first_free_address" "next" {
  subnet_id = phpipam_subnet.subnet.subnet_id

  depends_on = [
    "phpipam_address.address_1",
    "phpipam_address.address_2",
  ]
}
`

func TestAccDataSourcePHPIPAMFirstFreeAddress(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			sectionSweep("tf-test", t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccDataSourcePHPIPAMFirstFreeAddressConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.phpipam_first_free_address.next", "ip_address", "10.10.1.1"),
				),
			},
		},
	})
}

func TestAccDataSourcePHPIPAMFirstFreeAddressNoFree(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			sectionSweep("tf-test", t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config:      testAccDataSourcePHPIPAMFirstFreeAddressNoFreeConfig,
				ExpectError: regexp.MustCompile("No free addresses found"),
			},
		},
	})
}

const testAccDataSourcePHPIPAMFirstFreeAddressSubnetIDsConfig = `
resource "phpipam_section" "section" {
  name        = "tf-test"
  description = "Terraform test section"
}

resource "phpipam_subnet" "full" {
  section_id     = phpipam_section.section.section_id
  subnet_address = "10.10.4.0"
  subnet_mask    = 30
}

resource "phpipam_address" "full_1" {
  subnet_id  = phpipam_subnet.full.subnet_id
  ip_address = "10.10.4.1"
}

resource "phpipam_address" "full_2" {
  subnet_id  = phpipam_subnet.full.subnet_id
  ip_address = "10.10.4.2"
}

resource "phpipam_subnet" "free" {
  section_id     = phpipam_section.section.section_id
  subnet_address = "10.10.5.0"
  subnet_mask    = 24
}

data "phpipam_first_free_address" "next" {
  subnet_ids = [phpipam_subnet.full.subnet_id, phpipam_subnet.free.subnet_id]

  depends_on = [
    "phpipam_address.full_1",
    "phpipam_address.full_2",
    "phpipam_subnet.free",
  ]
}
`

const testAccDataSourcePHPIPAMFirstFreeAddressSubnetIDsNoFreeConfig = `
resource "phpipam_section" "section" {
  name        = "tf-test"
  description = "Terraform test section"
}

resource "phpipam_subnet" "full_1" {
  section_id     = phpipam_section.section.section_id
  subnet_address = "10.10.6.0"
  subnet_mask    = 30
}

resource "phpipam_address" "full_1_addr1" {
  subnet_id  = phpipam_subnet.full_1.subnet_id
  ip_address = "10.10.6.1"
}

resource "phpipam_address" "full_1_addr2" {
  subnet_id  = phpipam_subnet.full_1.subnet_id
  ip_address = "10.10.6.2"
}

resource "phpipam_subnet" "full_2" {
  section_id     = phpipam_section.section.section_id
  subnet_address = "10.10.7.0"
  subnet_mask    = 30
}

resource "phpipam_address" "full_2_addr1" {
  subnet_id  = phpipam_subnet.full_2.subnet_id
  ip_address = "10.10.7.1"
}

resource "phpipam_address" "full_2_addr2" {
  subnet_id  = phpipam_subnet.full_2.subnet_id
  ip_address = "10.10.7.2"
}

data "phpipam_first_free_address" "next" {
  subnet_ids = [phpipam_subnet.full_1.subnet_id, phpipam_subnet.full_2.subnet_id]

  depends_on = [
    "phpipam_address.full_1_addr1",
    "phpipam_address.full_1_addr2",
    "phpipam_address.full_2_addr1",
    "phpipam_address.full_2_addr2",
  ]
}
`

func TestAccDataSourcePHPIPAMFirstFreeAddressSubnetIDs(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			sectionSweep("tf-test", t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccDataSourcePHPIPAMFirstFreeAddressSubnetIDsConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.phpipam_first_free_address.next", "ip_address", "10.10.5.1"),
					resource.TestCheckResourceAttrPair("data.phpipam_first_free_address.next", "subnet_id", "phpipam_subnet.free", "subnet_id"),
				),
			},
		},
	})
}

func TestAccDataSourcePHPIPAMFirstFreeAddressSubnetIDsNoFree(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			sectionSweep("tf-test", t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config:      testAccDataSourcePHPIPAMFirstFreeAddressSubnetIDsNoFreeConfig,
				ExpectError: regexp.MustCompile("No free addresses found"),
			},
		},
	})
}
