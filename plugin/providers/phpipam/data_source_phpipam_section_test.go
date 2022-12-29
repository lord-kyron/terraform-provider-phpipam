package phpipam

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

const testAccDataSourcePHPIPAMSectionConfig = `
data "phpipam_section" "section_by_name" {
	name = "Customers"
}

data "phpipam_section" "section_by_id" {
	section_id = "${data.phpipam_section.section_by_name.section_id}"
}
`

func TestAccDataSourcePHPIPAMSection(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccDataSourcePHPIPAMSectionConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair("data.phpipam_section.section_by_name", "section_id", "data.phpipam_section.section_by_id", "section_id"),
					resource.TestCheckResourceAttrPair("data.phpipam_section.section_by_name", "name", "data.phpipam_section.section_by_id", "name"),
					resource.TestCheckResourceAttr("data.phpipam_section.section_by_name", "description", "Section for customers"),
				),
			},
		},
	})
}
