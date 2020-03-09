package phpipam

import (
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourcePHPIPAMSubnets() *schema.Resource {
	return &schema.Resource{
		Read: dataSourcePHPIPAMSubnetsRead,
		Schema: map[string]*schema.Schema{
			"section_id": &schema.Schema{
				Type:     schema.TypeInt,
				Required: true,
			},
			"description": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"description_match":   subnetDescriptionMatchSchema([]string{"description", "custom_field_filter"}),
			"custom_field_filter": customFieldFilterSchema([]string{"description", "description_match"}),
			"subnet_ids": &schema.Schema{
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeInt},
			},
		},
	}
}

func dataSourcePHPIPAMSubnetsRead(d *schema.ResourceData, meta interface{}) error {
	out, err := subnetSearchInSection(d, meta)
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
	err = d.Set("subnet_ids", ids)
	if err != nil {
		return err
	}

	return nil
}
