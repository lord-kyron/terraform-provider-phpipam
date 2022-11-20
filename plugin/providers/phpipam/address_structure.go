package phpipam

import (
	"errors"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/pavel-z1/phpipam-sdk-go/controllers/addresses"
	"github.com/pavel-z1/phpipam-sdk-go/phpipam"
)

// resourceAddressOptionalFields represents all the fields that are optional in
// the phpipam_address resource. These fields get flagged as Optional, with zero
// value defaults (the field is not set), in addition to being marked as
// Computed. Any field not listed here cannot be supplied by the resource and
// is solely computed.
var resourceAddressOptionalFields = linearSearchSlice{
	"is_gateway",
	"description",
	"hostname",
	"mac_address",
	"owner",
	"state_tag_id",
	"skip_ptr_record",
	"ptr_record_id",
	"device_id",
	"switch_port_label",
	"note",
	"exclude_ping",
}

// bareAddressSchema returns a map[string]*schema.Schema with the schema used
// to represent a PHPIPAM address resource. This output should then be modified
// so that required and computed fields are set properly for both the data
// source and the resource.
func bareAddressSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"address_id": &schema.Schema{
			Type: schema.TypeInt,
		},
		"subnet_id": &schema.Schema{
			Type: schema.TypeInt,
		},
		"ip_address": &schema.Schema{
			Type: schema.TypeString,
		},
		"is_gateway": &schema.Schema{
			Type: schema.TypeBool,
		},
		"description": &schema.Schema{
			Type: schema.TypeString,
		},
		"hostname": &schema.Schema{
			Type: schema.TypeString,
		},
		"mac_address": &schema.Schema{
			Type: schema.TypeString,
		},
		"owner": &schema.Schema{
			Type: schema.TypeString,
		},
		"state_tag_id": &schema.Schema{
			Type: schema.TypeInt,
		},
		"skip_ptr_record": &schema.Schema{
			Type: schema.TypeBool,
		},
		"ptr_record_id": &schema.Schema{
			Type: schema.TypeInt,
		},
		"device_id": &schema.Schema{
			Type: schema.TypeInt,
		},
		"switch_port_label": &schema.Schema{
			Type: schema.TypeString,
		},
		"note": &schema.Schema{
			Type: schema.TypeString,
		},
		"last_seen": &schema.Schema{
			Type: schema.TypeString,
		},
		"exclude_ping": &schema.Schema{
			Type: schema.TypeBool,
		},
		"edit_date": &schema.Schema{
			Type: schema.TypeString,
		},
		"custom_fields": &schema.Schema{
			Type: schema.TypeMap,
		},
	}
}

// resourceAddressSchema returns the schema for the phpipam_address resource.
// It sets the required and optional fields, the latter defined in
// resourceAddressOptionalFields, and ensures that all optional and
// non-configurable fields are computed as well.
func resourceAddressSchema() map[string]*schema.Schema {
	s := bareAddressSchema()
	for k, v := range s {
		switch {
		// IP Address and Subnet ID are ForceNew
		case k == "subnet_id":
			v.Required = true
			v.ForceNew = true
		case k == "ip_address":
			v.Optional = true
			v.Computed = true
			v.ForceNew = true
		case k == "custom_fields":
			v.Optional = true
		case resourceAddressOptionalFields.Has(k):
			v.Optional = true
			v.Computed = true
		default:
			v.Computed = true
		}
	}
	// Add the remove_dns_on_delete item to the schema. This is a meta-parameter
	// that is not part of the API resource and exists to instruct PHPIPAM to
	// gracefully remove the address from its DNS integrations as well when it is
	// removed. The default on this option is true.
	s["remove_dns_on_delete"] = &schema.Schema{
		Type:     schema.TypeBool,
		Optional: true,
		Default:  true,
	}
	return s
}

// dataSourceAddressSchema returns the schema for the phpipam_address data
// source. It sets the searchable fields and sets up the attribute conflicts
// between IP address and address ID. It also ensures that all fields are
// computed as well.
func dataSourceAddressSchema() map[string]*schema.Schema {
	s := bareAddressSchema()
	for k, v := range s {
		switch k {
		case "address_id":
			v.Optional = true
			v.Computed = true
			v.ConflictsWith = []string{"ip_address", "subnet_id", "description", "hostname", "custom_field_filter"}
		case "ip_address":
			v.Optional = true
			v.Computed = true
			v.ConflictsWith = []string{"address_id", "subnet_id", "description", "hostname", "custom_field_filter"}
		case "subnet_id":
			v.Optional = true
			v.Computed = true
			v.ConflictsWith = []string{"ip_address", "address_id"}
		case "description":
			v.Optional = true
			v.Computed = true
			v.ConflictsWith = []string{"ip_address", "address_id", "hostname", "custom_field_filter"}
		case "hostname":
			v.Optional = true
			v.Computed = true
			v.ConflictsWith = []string{"ip_address", "address_id", "description", "custom_field_filter"}
		default:
			v.Computed = true
		}
	}
	// Add the custom_field_filter item to the schema. This is a meta-parameter
	// that allows searching for a custom field value in the data source.
	s["custom_field_filter"] = customFieldFilterSchema([]string{"ip_address", "address_id", "hostname", "description"})

	return s
}

// expandAddress returns the addresses.Address structure for a
// phpipam_address resource or data source. Depending on if we are dealing with
// the resource or data source, extra considerations may need to be taken.
func expandAddress(d *schema.ResourceData) addresses.Address {
	s := addresses.Address{
		ID:          d.Get("address_id").(int),
		SubnetID:    d.Get("subnet_id").(int),
		IPAddress:   d.Get("ip_address").(string),
		IsGateway:   phpipam.BoolIntString(d.Get("is_gateway").(bool)),
		Description: d.Get("description").(string),
		Hostname:    d.Get("hostname").(string),
		MACAddress:  d.Get("mac_address").(string),
		Owner:       d.Get("owner").(string),
		Tag:         d.Get("state_tag_id").(int),
		PTRIgnore:   phpipam.BoolIntString(d.Get("skip_ptr_record").(bool)),
		PTRRecordID: d.Get("ptr_record_id").(int),
		DeviceID:    d.Get("device_id").(int),
		Port:        d.Get("switch_port_label").(string),
		Note:        d.Get("note").(string),
		LastSeen:    d.Get("last_seen").(string),
		ExcludePing: phpipam.BoolIntString(d.Get("exclude_ping").(bool)),
	}

	return s
}

// flattenAddress sets fields in a *schema.ResourceData with fields supplied by
// the input addresses.Address. This is used in read operations.
func flattenAddress(a addresses.Address, d *schema.ResourceData) {
	d.SetId(strconv.Itoa(a.ID))
	d.Set("address_id", a.ID)
	d.Set("subnet_id", a.SubnetID)
	d.Set("ip_address", a.IPAddress)
	d.Set("is_gateway", a.IsGateway)
	d.Set("description", a.Description)
	d.Set("hostname", a.Hostname)
	d.Set("mac_address", a.MACAddress)
	d.Set("owner", a.Owner)
	d.Set("state_tag_id", a.Tag)
	d.Set("skip_ptr_record", a.PTRIgnore)
	d.Set("ptr_record_id", a.PTRRecordID)
	d.Set("device_id", a.DeviceID)
	d.Set("switch_port_label", a.Port)
	d.Set("note", a.Note)
	d.Set("last_seen", a.LastSeen)
	d.Set("exclude_ping", a.ExcludePing)
	d.Set("edit_date", a.EditDate)
}

// addressSearchInSubnet provides the address search functionality for both the
// phpipam_address and phpipam_addresses data sources, returning an
// []addresses.Address to the particular data source that is calling the
// function. From here it's up to the specific data source to determine what
// they want to do with the results (ie: reject it on matching nothing or more
// than one for the singular data source, or extracting the IDs for the plural
// one).
func addressSearchInSubnet(d *schema.ResourceData, meta interface{}) ([]addresses.Address, error) {
	c := meta.(*ProviderPHPIPAMClient).addressesController
	s := meta.(*ProviderPHPIPAMClient).subnetsController
	result := make([]addresses.Address, 0)
	v, err := s.GetAddressesInSubnet(d.Get("subnet_id").(int))
	if err != nil {
		return result, err
	}
	if len(v) == 0 {
		return result, errors.New("No addresses were found in the supplied subnet")
	}
	for _, r := range v {
		switch {
		// Double-assert that we don't have empty strings in the conditionals
		// to ensure there there is no edge cases with matching zero values.
		case d.Get("description").(string) != "" && r.Description == d.Get("description").(string):
			result = append(result, r)
		case d.Get("hostname").(string) != "" && r.Hostname == d.Get("hostname").(string):
			result = append(result, r)
		case len(d.Get("custom_field_filter").(map[string]interface{})) > 0:
			fields, err := c.GetAddressCustomFields(r.ID)
			if err != nil {
				return result, err
			}
			search := d.Get("custom_field_filter").(map[string]interface{})
			matched, err := customFieldFilter(fields, search)
			if err != nil {
				return result, err
			}
			if matched {
				result = append(result, r)
			}
		}
	}
	return result, nil
}
