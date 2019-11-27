package phpipam

import (
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
)

const testAccDataSourcePHPIPAMAddressConfig = `
data "phpipam_address" "address_by_address" {
  ip_address = "10.10.1.245"
}

data "phpipam_address" "address_by_id" {
  address_id = "${data.phpipam_address.address_by_address.address_id}"
}

data "phpipam_address" "address_by_hostname" {
  subnet_id = 3
  hostname  = "server1.cust1.local"
}

data "phpipam_address" "address_by_description" {
  subnet_id   = 3
  description = "Server2"
}
`

const testAccDataSourcePHPIPAMAddressCustomFieldConfig = `
resource "phpipam_address" "address" {
  subnet_id   = 3
  ip_address  = "10.10.1.10"
  description = "Terraform test address (custom fields)"
  hostname    = "tf-test.cust1.local"

  custom_fields = {
    CustomTestAddresses  = "terraform-test"
    CustomTestAddresses2 = "terraform2-test"
  }
}

data "phpipam_address" "custom_search" {
  subnet_id = "${phpipam_address.address.subnet_id}"

  custom_field_filter = {
    CustomTestAddresses  = ".*terraform.*"
    CustomTestAddresses2 = ".*terraform2.*"
  }
}
`

func TestAccDataSourcePHPIPAMAddress(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccDataSourcePHPIPAMAddressConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair("data.phpipam_address.address_by_address", "address_id", "data.phpipam_address.address_by_id", "address_id"),
					resource.TestCheckResourceAttrPair("data.phpipam_address.address_by_address", "ip_address", "data.phpipam_address.address_by_id", "ip_address"),
					resource.TestCheckResourceAttr("data.phpipam_address.address_by_address", "description", "Gateway"),
					resource.TestCheckResourceAttr("data.phpipam_address.address_by_hostname", "ip_address", "10.10.1.3"),
					resource.TestCheckResourceAttr("data.phpipam_address.address_by_description", "ip_address", "10.10.1.4"),
				),
			},
		},
	})
}

func TestAccDataSourcePHPIPAMAddress_CustomField(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccDataSourcePHPIPAMAddressCustomFieldConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.phpipam_address.custom_search", "subnet_id", "3"),
					resource.TestCheckResourceAttr("data.phpipam_address.custom_search", "ip_address", "10.10.1.10"),
					resource.TestCheckResourceAttr("data.phpipam_address.custom_search", "description", "Terraform test address (custom fields)"),
					resource.TestCheckResourceAttr("data.phpipam_address.custom_search", "hostname", "tf-test.cust1.local"),
					resource.TestCheckResourceAttr("data.phpipam_address.custom_search", "custom_fields.CustomTestAddresses", "terraform-test"),
					resource.TestCheckResourceAttr("data.phpipam_address.custom_search", "custom_fields.CustomTestAddresses2", "terraform2-test"),
				),
			},
		},
	})
}
