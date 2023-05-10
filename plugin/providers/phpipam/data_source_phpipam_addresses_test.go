package phpipam

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

const testAccDataSourcePHPIPAMAddressesConfigStage1 = `
variable "ip_addresses" {
  type = list

  default = [
    "10.10.3.10",
    "10.10.3.11",
    "10.10.3.12",
    "10.10.3.13",
    "10.10.3.14",
  ]
}

resource "phpipam_section" "section" {
  name        = "tf-test"
  description = "Terraform test section"
}

resource "phpipam_subnet" "subnet" {
  section_id     = phpipam_section.section.section_id
  subnet_address = "10.10.3.0"
  subnet_mask    = 24
  description    = "Terraform test subnet (multiple addresses data source)"
}

resource "phpipam_address" "addresses" {
  count       = length(var.ip_addresses)
  subnet_id   = phpipam_subnet.subnet.subnet_id
  ip_address  = var.ip_addresses[count.index]
  description = "Terraform test address (multiple addresses data source)"
  hostname    = "tf-addresses-test1.example.internal"

  custom_fields = {
    custom_CustomTestAddresses  = "terraform-test-multiple"
    custom_CustomTestAddresses2 = "Entry ${var.ip_addresses[count.index]}"
  }
}
`

const testAccDataSourcePHPIPAMAddressesConfigStage2 = testAccDataSourcePHPIPAMAddressesConfigStage1 + `
data "phpipam_addresses" "addresses_by_description" {
  subnet_id = phpipam_subnet.subnet.subnet_id
  description = "Terraform test address (multiple addresses data source)"
}

data "phpipam_addresses" "addresses_by_hostname" {
  subnet_id = phpipam_subnet.subnet.subnet_id
  hostname = "tf-addresses-test1.example.internal"
}

data "phpipam_addresses" "addresses_by_custom_fields" {
  subnet_id = phpipam_subnet.subnet.subnet_id

  custom_field_filter = {
    custom_CustomTestAddresses  = "terraform-test-multiple"
    custom_CustomTestAddresses2 = "^Entry [0-9]"
  }
}

output "expected_address_ids" {
  value = phpipam_address.addresses.*.address_id
}

output "actual_address_ids_description" {
  value = data.phpipam_addresses.addresses_by_description.address_ids
}

output "actual_address_ids_hostname" {
  value = data.phpipam_addresses.addresses_by_hostname.address_ids
}

output "actual_address_ids_custom_fields" {
  value = data.phpipam_addresses.addresses_by_custom_fields.address_ids
}
`

func TestAccDataSourcePHPIPAMAddresses(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			sectionSweep("tf-test", t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccDataSourcePHPIPAMAddressesConfigStage1,
			},
			resource.TestStep{
				Config: testAccDataSourcePHPIPAMAddressesConfigStage2,
				Check: resource.ComposeTestCheckFunc(
					testCheckOutputPair("expected_address_ids", "actual_address_ids_description"),
					testCheckOutputPair("expected_address_ids", "actual_address_ids_hostname"),
					testCheckOutputPair("expected_address_ids", "actual_address_ids_custom_fields"),
				),
			},
		},
	})
}
