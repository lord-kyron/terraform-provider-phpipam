package phpipam

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

const testAccDataSourcePHPIPAML2DomainConfig = `
resource "phpipam_l2domain" "l2domain" {
  name        = "tf-test-l2domain"
  description = "Terraform test l2domain"
}

data "phpipam_l2domain" "l2domain_by_name" {
  name = "tf-test-l2domain"
  depends_on = [phpipam_l2domain.l2domain]
}

data "phpipam_l2domain" "l2domain_by_id" {
  domain_id = data.phpipam_l2domain.l2domain_by_name.domain_id
  depends_on = [data.phpipam_l2domain.l2domain_by_name]
}
`

func TestAccDataSourcePHPIPAML2Domain(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			l2domainSweep("tf-test-l2domain", t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccDataSourcePHPIPAML2DomainConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair("data.phpipam_l2domain.l2domain_by_name", "domain_id", "data.phpipam_l2domain.l2domain_by_id", "domain_id"),
					resource.TestCheckResourceAttrPair("data.phpipam_l2domain.l2domain_by_name", "name", "data.phpipam_l2domain.l2domain_by_id", "name"),
					resource.TestCheckResourceAttr("data.phpipam_l2domain.l2domain_by_name", "description", "Terraform test l2domain"),
				),
			},
		},
	})
}
