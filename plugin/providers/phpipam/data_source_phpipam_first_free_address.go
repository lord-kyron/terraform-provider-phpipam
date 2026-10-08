package phpipam

import (
	"errors"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourcePHPIPAMFirstFreeAddress() *schema.Resource {
	return &schema.Resource{
		Read: dataSourcePHPIPAMFirstFreeAddressRead,
		Schema: map[string]*schema.Schema{
			"subnet_id": &schema.Schema{
				Type:         schema.TypeInt,
				Optional:     true,
				Computed:     true,
				ExactlyOneOf: []string{"subnet_id", "subnet_ids"},
			},
			"subnet_ids": &schema.Schema{
				Type:         schema.TypeList,
				Optional:     true,
				Elem:         &schema.Schema{Type: schema.TypeInt},
				ExactlyOneOf: []string{"subnet_id", "subnet_ids"},
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

	ids := subnetIDsFromResourceData(d)

	var lastErr error
	for _, id := range ids {
		out, err := c.GetFirstFreeAddress(id)
		switch {
		case err != nil && strings.Contains(err.Error(), "No free addresses found"):
			lastErr = err
			continue
		case err != nil:
			return err
		case out == "":
			continue
		}

		d.SetId(out)
		d.Set("ip_address", out)
		d.Set("subnet_id", id)
		return nil
	}

	if lastErr != nil {
		return lastErr
	}
	return errors.New("Subnet has no free IP addresses")
}
