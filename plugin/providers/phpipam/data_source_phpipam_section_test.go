package phpipam

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

const testAccDataSourcePHPIPAMSectionConfig = `
resource "phpipam_section" "section" {
  name        = "tf-test"
  description = "Terraform test section"
}

data "phpipam_section" "section_by_name" {
  name = "tf-test"
  depends_on = [phpipam_section.section]
}

data "phpipam_section" "section_by_id" {
  section_id = data.phpipam_section.section_by_name.section_id
  depends_on = [data.phpipam_section.section_by_name]
}
`

func TestAccDataSourcePHPIPAMSection(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			sectionSweep("tf-test", t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccDataSourcePHPIPAMSectionConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair("data.phpipam_section.section_by_name", "section_id", "data.phpipam_section.section_by_id", "section_id"),
					resource.TestCheckResourceAttrPair("data.phpipam_section.section_by_name", "name", "data.phpipam_section.section_by_id", "name"),
					resource.TestCheckResourceAttr("data.phpipam_section.section_by_name", "description", "Terraform test section"),
				),
			},
		},
	})
}
