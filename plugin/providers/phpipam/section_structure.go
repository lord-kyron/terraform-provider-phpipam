package phpipam

import (
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/pavel-z1/phpipam-sdk-go/controllers/sections"
	"github.com/pavel-z1/phpipam-sdk-go/phpipam"
)

// resourceSectionOptionalFields represents all the fields that are optional in
// the phpipam_section resource. These fields get flagged as Optional, with zero
// value defaults (the field is not set), in addition to being marked as
// Computed. Any field not listed here cannot be supplied by the resource and
// is solely computed.
var resourceSectionOptionalFields = linearSearchSlice{
	"description",
	"master_section_id",
	"strict_mode",
	"subnet_ordering",
	"display_order",
	"show_vlan_in_subnet_listing",
	"show_vrf_in_subnet_listing",
	"show_supernet_only",
	"dns_resolver_id",
}

// bareSectionSchema returns a map[string]*schema.Schema with the schema used
// to represent a PHPIPAM Section resource. This output should then be modified
// so that required and computed fields are set properly for both the data
// source and the resource.
func bareSectionSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"section_id": &schema.Schema{
			Type: schema.TypeInt,
		},
		"name": &schema.Schema{
			Type: schema.TypeString,
		},
		"description": &schema.Schema{
			Type: schema.TypeString,
		},
		"master_section_id": &schema.Schema{
			Type: schema.TypeInt,
		},
		"permissions": &schema.Schema{
			Type: schema.TypeString,
		},
		"strict_mode": &schema.Schema{
			Type: schema.TypeBool,
		},
		"subnet_ordering": &schema.Schema{
			Type: schema.TypeString,
		},
		"display_order": &schema.Schema{
			Type: schema.TypeInt,
		},
		"edit_date": &schema.Schema{
			Type: schema.TypeString,
		},
		"show_vlan_in_subnet_listing": &schema.Schema{
			Type: schema.TypeBool,
		},
		"show_vrf_in_subnet_listing": &schema.Schema{
			Type: schema.TypeBool,
		},
		"show_supernet_only": &schema.Schema{
			Type: schema.TypeBool,
		},
		"dns_resolver_id": &schema.Schema{
			Type: schema.TypeInt,
		},
	}
}

// resourceSectionSchema returns the schema for the phpipam_section resource.
// It sets the required and optional fields, the latter defined in
// resourceSectionRequiredFields, and ensures that all optional and
// non-configurable fields are computed as well.
func resourceSectionSchema() map[string]*schema.Schema {
	schema := bareSectionSchema()
	for k, v := range schema {
		switch {
		// Section name is required
		case k == "name":
			v.Required = true
		case resourceSectionOptionalFields.Has(k):
			v.Optional = true
			v.Computed = true
		default:
			v.Computed = true
		}
	}
	return schema
}

// dataSourceSectionSchema returns the schema for the phpipam_section data source. It
// sets the searchable fields and sets up the attribute conflicts between Section
// entry ID and Section name. It also ensures that all fields are computed as
// well.
func dataSourceSectionSchema() map[string]*schema.Schema {
	schema := bareSectionSchema()
	for k, v := range schema {
		switch k {
		case "section_id":
			v.Optional = true
			v.Computed = true
			v.ConflictsWith = []string{"name"}
		case "name":
			v.Optional = true
			v.Computed = true
			v.ConflictsWith = []string{"section_id"}
		default:
			v.Computed = true
		}
	}
	return schema
}

// expandSection returns the sections.Section structure for a
// phpiapm_section resource or data source. Depending on if we are dealing with
// the resource or data source, extra considerations may need to be taken.
func expandSection(d *schema.ResourceData) sections.Section {
	s := sections.Section{
		ID:               d.Get("section_id").(int),
		Name:             d.Get("name").(string),
		Description:      d.Get("description").(string),
		MasterSection:    d.Get("master_section_id").(int),
		Permissions:      d.Get("permissions").(string),
		StrictMode:       phpipam.BoolIntString(d.Get("strict_mode").(bool)),
		SubnetOrdering:   d.Get("subnet_ordering").(string),
		Order:            d.Get("display_order").(int),
		EditDate:         d.Get("edit_date").(string),
		ShowVLAN:         phpipam.BoolIntString(d.Get("show_vlan_in_subnet_listing").(bool)),
		ShowVRF:          phpipam.BoolIntString(d.Get("show_vrf_in_subnet_listing").(bool)),
		ShowSupernetOnly: phpipam.BoolIntString(d.Get("show_supernet_only").(bool)),
		DNS:              d.Get("dns_resolver_id").(int),
	}

	return s
}

// flattenSection sets fields in a *schema.ResourceData with fields supplied by
// the input sections.Section. This is used in read operations.
func flattenSection(s sections.Section, d *schema.ResourceData) {
	d.SetId(strconv.Itoa(s.ID))
	d.Set("section_id", s.ID)
	d.Set("name", s.Name)
	d.Set("description", s.Description)
	d.Set("master_section_id", s.MasterSection)
	d.Set("permissions", s.Permissions)
	d.Set("strict_mode", s.StrictMode)
	d.Set("subnet_ordering", s.SubnetOrdering)
	d.Set("display_order", s.Order)
	d.Set("edit_date", s.EditDate)
	d.Set("show_vlan_in_subnet_listing", s.ShowVLAN)
	d.Set("show_vrf_in_subnet_listing", s.ShowVRF)
	d.Set("show_supernet_only", s.ShowSupernetOnly)
	d.Set("dns_resolver_id", s.DNS)
}
