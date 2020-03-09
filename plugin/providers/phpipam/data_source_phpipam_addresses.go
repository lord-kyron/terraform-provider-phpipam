package phpipam

import (
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourcePHPIPAMAddresses() *schema.Resource {
	return &schema.Resource{
		Read: dataSourcePHPIPAMAddressesRead,
		Schema: map[string]*schema.Schema{
			"subnet_id": &schema.Schema{
				Type:     schema.TypeInt,
				Required: true,
			},
			"description": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"hostname": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"custom_field_filter": customFieldFilterSchema([]string{"description", "hostname"}),
			"address_ids": &schema.Schema{
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeInt},
			},
		},
	}
}

func dataSourcePHPIPAMAddressesRead(d *schema.ResourceData, meta interface{}) error {
	out, err := addressSearchInSubnet(d, meta)
	if err != nil {
		return err
	}
	var sum int
	ids := make([]int, 0)
	for _, v := range out {
		sum += v.ID
		ids = append(ids, v.ID)
	}

	d.SetId(strconv.Itoa(sum))
	err = d.Set("address_ids", ids)
	if err != nil {
		return err
	}

	return nil
}
