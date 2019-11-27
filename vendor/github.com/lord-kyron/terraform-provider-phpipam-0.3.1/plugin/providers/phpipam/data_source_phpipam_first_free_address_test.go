package phpipam

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
)

const testAccDataSourcePHPIPAMFirstFreeAddressConfig = `
data "phpipam_subnet" "subnet_by_cidr" {
	subnet_address = "10.10.1.0"
	subnet_mask = 24
}

data "phpipam_first_free_address" "next" {
	subnet_id = "${data.phpipam_subnet.subnet_by_cidr.subnet_id}"
}
`

const testAccDataSourcePHPIPAMFirstFreeAddressNoFreeConfig = `
resource "phpipam_subnet" "subnet" {
  section_id     = 1
  subnet_address = "10.10.3.0"
  subnet_mask    = 30
}

resource "phpipam_address" "address_1" {
  subnet_id  = "${phpipam_subnet.subnet.subnet_id}"
  ip_address = "10.10.3.1"
}

resource "phpipam_address" "address_2" {
  subnet_id  = "${phpipam_subnet.subnet.subnet_id}"
  ip_address = "10.10.3.2"
}

data "phpipam_first_free_address" "next" {
  subnet_id = "${phpipam_subnet.subnet.subnet_id}"

  depends_on = [
    "phpipam_address.address_1",
    "phpipam_address.address_2",
  ]
}
`

func TestAccDataSourcePHPIPAMFirstFreeAddress(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
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
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config:      testAccDataSourcePHPIPAMFirstFreeAddressNoFreeConfig,
				ExpectError: regexp.MustCompile("Subnet has no free IP addresses"),
			},
		},
	})
}
