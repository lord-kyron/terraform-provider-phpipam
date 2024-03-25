package phpipam

import (
	"context"
	"errors"
	"fmt"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/pavel-z1/phpipam-sdk-go/controllers/subnets"
)

func dataSourcePHPIPAMSubnet() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourcePHPIPAMSubnetRead,
		Schema:      dataSourceSubnetSchema(),
	}
}

func dataSourcePHPIPAMSubnetRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*ProviderPHPIPAMClient).subnetsController
	out := make([]subnets.Subnet, 1)
	var err error
	// We need to determine how to get the subnet. An ID search takes priority,
	// and after that subnets.
	switch {
	case d.Get("subnet_id").(int) != 0:
		out[0], err = c.GetSubnetByID(d.Get("subnet_id").(int))
		if err != nil {
			return diag.FromErr(err)
		}
	case d.Get("subnet_address").(string) != "" && d.Get("subnet_mask").(int) != 0 && d.Get("section_id").(int) == 0:
		out, err = c.GetSubnetsByCIDR(fmt.Sprintf("%s/%d", d.Get("subnet_address"), d.Get("subnet_mask")))
		if err != nil {
			return diag.FromErr(err)
		}
	case d.Get("subnet_address").(string) != "" && d.Get("subnet_mask").(int) != 0 && d.Get("section_id").(int) != 0:
		out, err = c.GetSubnetsByCIDRAndSection(fmt.Sprintf("%s/%d", d.Get("subnet_address"), d.Get("subnet_mask")), d.Get("section_id").(int))
		if err != nil {
			return diag.FromErr(err)
		}
	case d.Get("section_id").(int) != 0 && (d.Get("description").(string) != "" || d.Get("description_match").(string) != "" || len(d.Get("custom_field_filter").(map[string]interface{})) > 0):
		out, err = subnetSearchInSection(d, meta)
		if err != nil {
			return diag.FromErr(err)
		}
	default:
		// We need to ensure imported resources are not recreated when terraform apply is ran
		// imported resources only have an Id which we need to map back to subnet_id
		id := d.Id()
		if len(id) > 0 {
			subnet_id, err := strconv.Atoi(id)
			if err != nil {
				return diag.FromErr(err)
			}
			out[0], err = c.GetSubnetByID(subnet_id)
			if err != nil {
				return diag.FromErr(err)
			}
		} else {
			return diag.FromErr(errors.New("No valid combination of parameters found - need one of subnet_id, subnet_address and subnet_mask, or section_id and (description|description_match|custom_field_filter)"))
		}
	}
	if len(out) != 1 {
		return diag.FromErr(errors.New("Your search returned zero or multiple results. Please correct your search and try again"))
	}

	if !meta.(*ProviderPHPIPAMClient).NestCustomFields {
		if checkSubnetsCustomFiledsExists(d, c) {
			fields, err := c.GetSubnetCustomFields(out[0].ID)
			switch {
			case err == nil:
				trimMap(fields)
				if err := d.Set("custom_fields", fields); err != nil {
					return diag.FromErr(err)
				}
			case err != nil:
				return diag.FromErr(err)
			}
		}
	}

	flattenSubnet(out[0], d)

	if out[0].CustomFields != nil && !meta.(*ProviderPHPIPAMClient).NestCustomFields {

		var diags diag.Diagnostics
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  "Nest custom fields is enabled on the API",
			Detail:   "This API has enabled nested custom fields. Please set nest_custom_fields to true in the provider configuration.",
		})
		return diags
	}

	return nil
}

func checkSubnetsCustomFiledsExists(d *schema.ResourceData, client *subnets.Controller) bool {
	if _, ok := d.GetOk("custom_field_filter"); ok {
		return true
	} else if _, err := client.GetSubnetCustomFieldsSchema(); err == nil {
		return true
	} else {
		return false
	}
}
