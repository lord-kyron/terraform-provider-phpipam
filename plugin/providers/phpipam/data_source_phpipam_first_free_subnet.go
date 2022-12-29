package phpipam

import (
	"errors"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourcePHPIPAMFirstFreeSubnet() *schema.Resource {
	return &schema.Resource{
		Read: dataSourcePHPIPAMFirstFreeSubnetRead,
		Schema: map[string]*schema.Schema{
			"subnet_id": &schema.Schema{
				Type:     schema.TypeInt,
				Required: true,
			},
			"subnet_mask": &schema.Schema{
				Type:     schema.TypeInt,
				Required: true,
			},
			"ip_address": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourcePHPIPAMFirstFreeSubnetRead(d *schema.ResourceData, meta interface{}) error {
	c := meta.(*ProviderPHPIPAMClient).subnetsController
	out, err := c.GetFirstFreeSubnet(d.Get("subnet_id").(int), d.Get("subnet_mask").(int))
	if err != nil {
		return err
	}
	if out == "" {
		return errors.New("Subnet has no free IP addresses")
	}

	d.SetId(out)
	d.Set("ip_address", out)

	return nil
}
