package phpipam

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/pavel-z1/phpipam-sdk-go/controllers/subnets"
)

func dataSourcePHPIPAMSubnets() *schema.Resource {
	return &schema.Resource{
		Read: dataSourcePHPIPAMSubnetsRead,
		Schema: map[string]*schema.Schema{
			"section_id": &schema.Schema{
				Type:         schema.TypeInt,
				Optional:     true,
				AtLeastOneOf: []string{"section_id", "search"},
			},
			"search": &schema.Schema{
				Type:         schema.TypeString,
				Optional:     true,
				AtLeastOneOf: []string{"section_id", "search"},
				ValidateFunc: validateSubnetSearch,
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
	var out []subnets.Subnet
	var err error
	if d.Get("search").(string) != "" {
		out, err = subnetSearchAll(d, meta)
	} else {
		out, err = subnetSearchInSection(d, meta)
	}
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

// validateSubnetSearch rejects search strings that cannot be sent in the
// search controller's URL path.
func validateSubnetSearch(v interface{}, k string) (ws []string, errs []error) {
	if strings.ContainsAny(v.(string), "/?#") {
		errs = append(errs, fmt.Errorf("%s must not contain '/', '?' or '#'; to search for a subnet by CIDR, search for its network address", k))
	}
	return
}
