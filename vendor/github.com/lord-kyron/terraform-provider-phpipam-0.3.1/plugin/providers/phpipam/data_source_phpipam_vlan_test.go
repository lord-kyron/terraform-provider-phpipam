package phpipam

import (
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
)

const testAccDataSourcePHPIPAMVLANConfig = `
data "phpipam_vlan" "vlan_by_number" {
	number = 2001
}

data "phpipam_vlan" "vlan_by_id" {
	vlan_id = "${data.phpipam_vlan.vlan_by_number.vlan_id}"
}
`

func TestAccDataSourcePHPIPAMVLAN(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccDataSourcePHPIPAMVLANConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair("data.phpipam_vlan.vlan_by_number", "vlan_id", "data.phpipam_vlan.vlan_by_id", "vlan_id"),
					resource.TestCheckResourceAttrPair("data.phpipam_vlan.vlan_by_number", "number", "data.phpipam_vlan.vlan_by_id", "number"),
					resource.TestCheckResourceAttr("data.phpipam_vlan.vlan_by_number", "name", "IPv6 private 1"),
					resource.TestCheckResourceAttr("data.phpipam_vlan.vlan_by_number", "description", "IPv6 private 1 subnets"),
				),
			},
		},
	})
}
