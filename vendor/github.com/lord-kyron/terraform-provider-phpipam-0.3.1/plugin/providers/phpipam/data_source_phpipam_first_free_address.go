package phpipam

import (
	"errors"

	"github.com/hashicorp/terraform/helper/schema"
)

func dataSourcePHPIPAMFirstFreeAddress() *schema.Resource {
	return &schema.Resource{
		Read: dataSourcePHPIPAMFirstFreeAddressRead,
		Schema: map[string]*schema.Schema{
			"subnet_id": &schema.Schema{
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

func dataSourcePHPIPAMFirstFreeAddressRead(d *schema.ResourceData, meta interface{}) error {
	c := meta.(*ProviderPHPIPAMClient).subnetsController
	out, err := c.GetFirstFreeAddress(d.Get("subnet_id").(int))
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
