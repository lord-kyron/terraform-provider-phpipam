package phpipam

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

const testAccDataSourcePHPIPAMVLANConfig = `
resource "phpipam_l2domain" "tf-test-l2domain" {
  name        = "tf-test-l2domain"
  description = "Terraform test l2domain"
}

resource "phpipam_vlan" "vlan" {
  name        = "tf-test-vlan"
  number      = 2001
  l2_domain_id = phpipam_l2domain.tf-test-l2domain.domain_id
  description = "Terraform test vlan"
  depends_on  =[phpipam_l2domain.tf-test-l2domain]
}

data "phpipam_vlan" "vlan_by_number" {
  number       = 2001
  l2_domain_id = phpipam_l2domain.tf-test-l2domain.domain_id
  depends_on   = [phpipam_vlan.vlan]
}

data "phpipam_vlan" "vlan_by_id" {
   vlan_id = data.phpipam_vlan.vlan_by_number.vlan_id
   depends_on = [data.phpipam_vlan.vlan_by_number]
}
`

func TestAccDataSourcePHPIPAMVLAN(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			l2domainSweep("tf-test-l2domain", t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccDataSourcePHPIPAMVLANConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair("data.phpipam_vlan.vlan_by_number", "vlan_id", "data.phpipam_vlan.vlan_by_id", "vlan_id"),
					resource.TestCheckResourceAttrPair("data.phpipam_vlan.vlan_by_number", "number", "data.phpipam_vlan.vlan_by_id", "number"),
					resource.TestCheckResourceAttr("data.phpipam_vlan.vlan_by_number", "name", "tf-test-vlan"),
					resource.TestCheckResourceAttr("data.phpipam_vlan.vlan_by_number", "description", "Terraform test vlan"),
				),
			},
		},
	})
}
