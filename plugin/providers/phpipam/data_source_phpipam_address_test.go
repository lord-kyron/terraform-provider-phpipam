package phpipam

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

const testAccDataSourcePHPIPAMAddressConfig = `
resource "phpipam_section" "test" {
  name        = "tf-test"
  description = "Terraform test section"
}

resource "phpipam_subnet" "subnet" {
  section_id     = phpipam_section.test.section_id
  subnet_address = "10.10.3.0"
  subnet_mask    = 24

  depends_on = [phpipam_section.test]
}

resource "phpipam_address" "address" {
  subnet_id   =  phpipam_subnet.subnet.subnet_id
  ip_address  = "10.10.3.245"
  description = "Server2"
  hostname    = "server1.cust1.local"

  depends_on = [phpipam_subnet.subnet]
}


data "phpipam_address" "address_by_address" {
  ip_address = "10.10.3.245"
  depends_on = [phpipam_address.address]
}

data "phpipam_address" "address_by_id" {
  address_id = data.phpipam_address.address_by_address.address_id
  depends_on = [data.phpipam_address.address_by_address]
}

data "phpipam_address" "address_by_hostname" {
  subnet_id = phpipam_subnet.subnet.subnet_id
  hostname  = "server1.cust1.local"
  depends_on     = [phpipam_address.address]
}

data "phpipam_address" "address_by_description" {
  subnet_id   = phpipam_subnet.subnet.subnet_id
  description = "Server2"
  depends_on     = [phpipam_address.address]
}
`

const testAccDataSourcePHPIPAMAddressCustomFieldConfig = `
resource "phpipam_section" "test" {
  name        = "tf-test"
  description = "Terraform test section"
}

resource "phpipam_subnet" "subnet" {
  section_id     = phpipam_section.test.section_id
  subnet_address = "10.10.3.0"
  subnet_mask    = 24

  depends_on = [phpipam_section.test]
}

resource "phpipam_address" "address" {
  subnet_id   =  phpipam_subnet.subnet.subnet_id
  ip_address  = "10.10.3.245"
  description = "Server2"
  hostname    = "server1.cust1.local"

  custom_fields = {
    custom_CustomTestAddresses  = "terraform-test"
    custom_CustomTestAddresses2 = "terraform2-test"
  }
  depends_on = [phpipam_subnet.subnet]
}


data "phpipam_address" "address_by_address" {
  ip_address = "10.10.3.245"
  depends_on = [phpipam_address.address]
}

data "phpipam_address" "address_by_id" {
  address_id = data.phpipam_address.address_by_address.address_id
  depends_on = [data.phpipam_address.address_by_address]
}

data "phpipam_address" "address_by_hostname" {
  subnet_id = phpipam_subnet.subnet.subnet_id
  hostname  = "server1.cust1.local"
  depends_on     = [phpipam_address.address]
}

data "phpipam_address" "address_by_description" {
  subnet_id   = phpipam_subnet.subnet.subnet_id
  description = "Server2"
  depends_on     = [phpipam_address.address]
}

data "phpipam_address" "custom_search" {
  subnet_id = phpipam_address.address.subnet_id

  custom_field_filter = {
    custom_CustomTestAddresses  = ".*terraform.*"
  }
}
`

func TestAccDataSourcePHPIPAMAddress(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			sectionSweep("tf-test", t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccDataSourcePHPIPAMAddressConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("phpipam_section.test", "name", "tf-test"),
					resource.TestCheckResourceAttrPair("data.phpipam_address.address_by_address", "address_id", "data.phpipam_address.address_by_id", "address_id"),
					resource.TestCheckResourceAttrPair("data.phpipam_address.address_by_address", "ip_address", "data.phpipam_address.address_by_id", "ip_address"),
					resource.TestCheckResourceAttr("data.phpipam_address.address_by_address", "description", "Server2"),
					resource.TestCheckResourceAttr("data.phpipam_address.address_by_hostname", "ip_address", "10.10.3.245"),
					resource.TestCheckResourceAttr("data.phpipam_address.address_by_description", "ip_address", "10.10.3.245"),
				),
			},
			resource.TestStep{
				Config: testAccDataSourcePHPIPAMAddressCustomFieldConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.phpipam_address.custom_search", "ip_address", "10.10.3.245"),
					resource.TestCheckResourceAttr("data.phpipam_address.custom_search", "description", "Server2"),
					resource.TestCheckResourceAttr("data.phpipam_address.custom_search", "hostname", "server1.cust1.local"),
					resource.TestCheckResourceAttr("data.phpipam_address.custom_search", "custom_fields.custom_CustomTestAddresses", "terraform-test"),
					resource.TestCheckResourceAttr("data.phpipam_address.custom_search", "custom_fields.custom_CustomTestAddresses2", "terraform2-test"),
				),
			},
		},
	})
}
