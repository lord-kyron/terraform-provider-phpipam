package phpipam

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

const testAccResourcePHPIPAMIPv6Config = `
resource "phpipam_section" "section" {
  name        = "tf-test"
  description = "Terraform test section"
}

resource "phpipam_subnet" "subnet" {
  section_id     = phpipam_section.section.section_id
  subnet_address = "2001:db8:10::"
  subnet_mask    = 64
  description    = "Terraform test IPv6 subnet"
}

resource "phpipam_address" "address" {
  subnet_id   = phpipam_subnet.subnet.subnet_id
  ip_address  = "2001:db8:10::10"
  hostname    = "tf-test-v6.cust1.local"
  description = "Terraform test IPv6 address"
}

resource "phpipam_first_free_address" "next" {
  subnet_id   = phpipam_subnet.subnet.subnet_id
  hostname    = "tf-test-v6-next.cust1.local"
  description = "Terraform test IPv6 first free address"

  depends_on = [phpipam_address.address]
}

resource "phpipam_first_free_subnet" "child" {
  parent_subnet_id = phpipam_subnet.subnet.subnet_id
  subnet_mask      = 80
  description      = "Terraform test IPv6 first free subnet"

  depends_on = [phpipam_first_free_address.next]
}

data "phpipam_subnet" "by_cidr" {
  subnet_address = phpipam_subnet.subnet.subnet_address
  subnet_mask    = phpipam_subnet.subnet.subnet_mask
  section_id     = phpipam_section.section.section_id
}

data "phpipam_address" "by_ip" {
  ip_address = phpipam_address.address.ip_address
}
`

func TestAccResourcePHPIPAMIPv6(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			sectionSweep("tf-test", t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccResourcePHPIPAMIPv6Config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("phpipam_subnet.subnet", "subnet_address", "2001:db8:10::"),
					resource.TestCheckResourceAttr("phpipam_subnet.subnet", "subnet_mask", "64"),
					resource.TestCheckResourceAttr("phpipam_address.address", "ip_address", "2001:db8:10::10"),
					resource.TestCheckResourceAttr("phpipam_address.address", "hostname", "tf-test-v6.cust1.local"),
					// Unlike IPv4, phpIPAM hands out the first address of an IPv6 subnet.
					resource.TestCheckResourceAttr("phpipam_first_free_address.next", "ip_address", "2001:db8:10::"),
					resource.TestCheckResourceAttrPair("phpipam_first_free_address.next", "subnet_id", "phpipam_subnet.subnet", "subnet_id"),
					resource.TestCheckResourceAttr("phpipam_first_free_subnet.child", "subnet_address", "2001:db8:10::"),
					resource.TestCheckResourceAttr("phpipam_first_free_subnet.child", "subnet_mask", "80"),
					resource.TestCheckResourceAttrPair("data.phpipam_subnet.by_cidr", "subnet_id", "phpipam_subnet.subnet", "subnet_id"),
					resource.TestCheckResourceAttrPair("data.phpipam_address.by_ip", "address_id", "phpipam_address.address", "address_id"),
				),
			},
			resource.TestStep{
				ResourceName:      "phpipam_subnet.subnet",
				ImportState:       true,
				ImportStateVerify: true,
			},
			resource.TestStep{
				ResourceName:      "phpipam_address.address",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
