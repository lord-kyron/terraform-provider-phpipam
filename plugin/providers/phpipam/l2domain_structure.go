package phpipam

import (
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/pavel-z1/phpipam-sdk-go/controllers/l2domains"
)

// resourceL2DomainOptionalFields represents all the fields that are optional in
// the phpipam_l2domain resource. These fields get flagged as Optional, with zero
// value defaults (the field is not set), in addition to being marked as
// Computed. Any field not listed here cannot be supplied by the resource and
// is solely computed.
var resourceL2DomainOptionalFields = linearSearchSlice{
	"description",
	"sections",
}

// bareL2DomainSchema returns a map[string]*schema.Schema with the schema used
// to represent a PHPIPAM L2Domain resource. This output should then be modified
// so that required and computed fields are set properly for both the data
// source and the resource.
func bareL2DomainSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"domain_id": &schema.Schema{
			Type: schema.TypeInt,
		},
		"name": &schema.Schema{
			Type: schema.TypeString,
		},
		"description": &schema.Schema{
			Type: schema.TypeString,
		},
		"sections": &schema.Schema{
			Type: schema.TypeString,
		},
	}
}

// resourceL2DomainSchema returns the schema for the phpipam_l2domain resource.
// It sets the required and optional fields, the latter defined in
// resourceL2DomanRequiredFields, and ensures that all optional and
// non-configurable fields are computed as well.
func resourceL2DomainSchema() map[string]*schema.Schema {
	schema := bareL2DomainSchema()
	for k, v := range schema {
		switch {
		// Section name is required
		case k == "name":
			v.Required = true
		case resourceL2DomainOptionalFields.Has(k):
			v.Optional = true
			v.Computed = true
		default:
			v.Computed = true
		}
	}
	return schema
}

// dataSourceL2DomainSchema returns the schema for the phpipam_l2domain data source. It
// sets the searchable fields and sets up the attribute conflicts between L2Domain
// entry ID and L2Domain name. It also ensures that all fields are computed as
// well.
func dataSourceL2DomainSchema() map[string]*schema.Schema {
	schema := bareL2DomainSchema()
	for k, v := range schema {
		switch k {
		case "domain_id":
			v.Optional = true
			v.Computed = true
			v.ConflictsWith = []string{"name"}
		case "name":
			v.Optional = true
			v.Computed = true
			v.ConflictsWith = []string{"domain_id"}
		default:
			v.Computed = true
		}
	}
	return schema
}

// expandL2Domain returns the l2domains.L2Domain structure for a
// phpiapm_l2domain resource or data source. Depending on if we are dealing with
// the resource or data source, extra considerations may need to be taken.
func expandL2Domain(d *schema.ResourceData) l2domains.L2Domain {
	l := l2domains.L2Domain{
		ID:          d.Get("domain_id").(int),
		Name:        d.Get("name").(string),
		Description: d.Get("description").(string),
		Sections:    d.Get("sections").(string),
	}

	return l
}

// flattenL2Domain sets fields in a *schema.ResourceData with fields supplied by
// the input l2domains.L2Domain. This is used in read operations.
func flattenL2Domain(l l2domains.L2Domain, d *schema.ResourceData) {
	d.SetId(strconv.Itoa(l.ID))
	d.Set("domain_id", l.ID)
	d.Set("name", l.Name)
	d.Set("description", l.Description)
	d.Set("sections", l.Sections)
}
