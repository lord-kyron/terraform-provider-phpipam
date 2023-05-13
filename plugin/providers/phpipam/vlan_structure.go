package phpipam

import (
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/pavel-z1/phpipam-sdk-go/controllers/vlans"
)

// resourceVLANOptionalFields represents all the fields that are optional in
// the phpipam_vlan resource. These fields get flagged as Optional, with zero
// value defaults (the field is not set), in addition to being marked as
// Computed. Any field not listed here cannot be supplied by the resource and
// is solely computed.
var resourceVLANOptionalFields = linearSearchSlice{
	"l2_domain_id",
	"description",
}

// bareVLANSchema returns a map[string]*schema.Schema with the schema used
// to represent a PHPIPAM VLAN resource. This output should then be modified
// so that required and computed fields are set properly for both the data
// source and the resource.
func bareVLANSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"vlan_id": &schema.Schema{
			Type: schema.TypeInt,
		},
		"l2_domain_id": &schema.Schema{
			Type: schema.TypeInt,
		},
		"name": &schema.Schema{
			Type: schema.TypeString,
		},
		"number": &schema.Schema{
			Type: schema.TypeInt,
		},
		"description": &schema.Schema{
			Type: schema.TypeString,
		},
		"edit_date": &schema.Schema{
			Type: schema.TypeString,
		},
		"custom_fields": &schema.Schema{
			Type: schema.TypeMap,
		},
	}
}

// resourceVLANSchema returns the schema for the phpipam_vlan resource.
// It sets the required and optional fields, the latter defined in
// resourceVLANRequiredFields, and ensures that all optional and
// non-configurable fields are computed as well.
func resourceVLANSchema() map[string]*schema.Schema {
	schema := bareVLANSchema()
	for k, v := range schema {
		switch {
		// VLAN name and number are required
		case k == "name" || k == "number":
			v.Required = true
		case k == "custom_fields":
			v.Optional = true
		case resourceVLANOptionalFields.Has(k):
			v.Optional = true
			v.Computed = true
		default:
			v.Computed = true
		}
	}
	return schema
}

// dataSourceVLANSchema returns the schema for the phpipam_vlan data source. It
// sets the searchable fields and sets up the attribute conflicts between VLAN
// entry ID and VLAN number. It also ensures that all fields are computed as
// well.
func dataSourceVLANSchema() map[string]*schema.Schema {
	s := bareVLANSchema()
	for k, v := range s {
		switch k {
		case "vlan_id":
			v.Optional = true
			v.Computed = true
			v.ConflictsWith = []string{"number"}
		case "number":
			v.Optional = true
			v.Computed = true
			v.ConflictsWith = []string{"vlan_id"}
		case "l2_domain_id":
			v.Optional = true
			v.Computed = true
			v.ConflictsWith = []string{"vlan_id"}
		default:
			v.Computed = true
		}
	}
	return s
}

// expandVLAN returns the vlans.VLAN structure for a
// phpiapm_vlan resource or data source. Depending on if we are dealing with
// the resource or data source, extra considerations may need to be taken.
func expandVLAN(d *schema.ResourceData) vlans.VLAN {
	v := vlans.VLAN{
		ID:          d.Get("vlan_id").(int),
		DomainID:    d.Get("l2_domain_id").(int),
		Name:        d.Get("name").(string),
		Number:      d.Get("number").(int),
		Description: d.Get("description").(string),
	}

	return v
}

// flattenVLAN sets fields in a *schema.ResourceData with fields supplied by
// the input vlans.VLAN. This is used in read operations.
func flattenVLAN(v vlans.VLAN, d *schema.ResourceData) {
	d.SetId(strconv.Itoa(v.ID))
	d.Set("vlan_id", v.ID)
	d.Set("l2_domain_id", v.DomainID)
	d.Set("name", v.Name)
	d.Set("number", v.Number)
	d.Set("description", v.Description)
	d.Set("edit_date", v.EditDate)
}
