package phpipam

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

const testAccDataSourcePHPIPAMSubnetsConfigStage1 = `
variable "subnet_addresses" {
  type = "list"

  default = [
    "10.10.3.0/24",
    "10.10.4.0/24",
    "10.10.5.0/24",
  ]
}

resource "phpipam_subnet" "subnets" {
  count          = "${length(var.subnet_addresses)}"
  section_id     = 1
  subnet_address = "${element(split("/", element(var.subnet_addresses, count.index)), 0)}"
  subnet_mask    = "${element(split("/", element(var.subnet_addresses, count.index)), 1)}"
  description    = "Terraform test subnet (multiple subnet data source)"

  custom_fields = {
    CustomTestSubnets  = "terraform-test-multiple"
    CustomTestSubnets2 = "Entry ${var.subnet_addresses[count.index]}"
  }
}
`

const testAccDataSourcePHPIPAMSubnetsConfigStage2 = testAccDataSourcePHPIPAMSubnetsConfigStage1 + `
data "phpipam_subnets" "subnets_by_description" {
  section_id  = 1
  description = "Terraform test subnet (multiple subnet data source)"
}

data "phpipam_subnets" "subnets_by_description_match" {
  section_id        = 1
  description_match = ".*(multiple subnet data source).*"
}

data "phpipam_subnets" "subnets_by_custom_fields" {
  section_id = 1

  custom_field_filter = {
    CustomTestSubnets  = "terraform-test-multiple"
    CustomTestSubnets2 = "^Entry [0-9]"
  }
}

output "expected_subnet_ids" {
  value = ["${phpipam_subnet.subnets.*.id}"]
}

output "actual_subnet_ids_description" {
  value = "${data.phpipam_subnets.subnets_by_description.subnet_ids}"
}

output "actual_subnet_ids_description_match" {
  value = "${data.phpipam_subnets.subnets_by_description_match.subnet_ids}"
}

output "actual_subnet_ids_custom_fields" {
  value = "${data.phpipam_subnets.subnets_by_custom_fields.subnet_ids}"
}
`

func TestAccDataSourcePHPIPAMSubnets(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccDataSourcePHPIPAMSubnetsConfigStage1,
			},
			resource.TestStep{
				Config: testAccDataSourcePHPIPAMSubnetsConfigStage2,
				Check: resource.ComposeTestCheckFunc(
					testCheckOutputPair("expected_subnet_ids", "actual_subnet_ids_description"),
					testCheckOutputPair("expected_subnet_ids", "actual_subnet_ids_description_match"),
					testCheckOutputPair("expected_subnet_ids", "actual_subnet_ids_custom_fields"),
				),
			},
		},
	})
}
