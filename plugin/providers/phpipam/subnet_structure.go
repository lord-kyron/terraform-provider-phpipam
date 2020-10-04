package phpipam

import (
	"errors"
	"regexp"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/pavel-z1/phpipam-sdk-go/controllers/subnets"
	"github.com/pavel-z1/phpipam-sdk-go/phpipam"
)

// resourceSubnetOptionalFields represents all the fields that are optional in
// the phpipam_subnet resource. These fields get flagged as Optional, with zero
// value defaults (the field is not set), in addition to being marked as
// Computed. Any field not listed here cannot be supplied by the resource and
// is solely computed.
var resourceSubnetOptionalFields = linearSearchSlice{
	"description",
	"linked_subnet_id",
	"vlan_id",
	"vrf_id",
	"master_subnet_id",
	"nameserver_id",
	"show_name",
	"create_ptr_records",
	"display_hostnames",
	"allow_ip_requests",
	"scan_agent_id",
	"include_in_ping",
	"host_discover_enabled",
	"is_full",
	"utilization_threshold",
}

// bareSubnetSchema returns a map[string]*schema.Schema with the schema used
// to represent a PHPIPAM subnet resource. This output should then be modified
// so that required and computed fields are set properly for both the data
// source and the resource.
func bareSubnetSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"subnet_id": &schema.Schema{
			Type: schema.TypeInt,
		},
		"subnet_address": &schema.Schema{
			Type: schema.TypeString,
		},
		"subnet_mask": &schema.Schema{
			Type: schema.TypeInt,
		},
		"description": &schema.Schema{
			Type: schema.TypeString,
		},
		"section_id": &schema.Schema{
			Type: schema.TypeInt,
		},
		"linked_subnet_id": &schema.Schema{
			Type: schema.TypeInt,
		},
		"vlan_id": &schema.Schema{
			Type: schema.TypeInt,
		},
		"vrf_id": &schema.Schema{
			Type: schema.TypeInt,
		},
		"master_subnet_id": &schema.Schema{
			Type: schema.TypeInt,
		},
		"nameserver_id": &schema.Schema{
			Type: schema.TypeInt,
		},
		"show_name": &schema.Schema{
			Type: schema.TypeBool,
		},
		"permissions": &schema.Schema{
			Type: schema.TypeString,
		},
		"create_ptr_records": &schema.Schema{
			Type: schema.TypeBool,
		},
		"display_hostnames": &schema.Schema{
			Type: schema.TypeBool,
		},
		"allow_ip_requests": &schema.Schema{
			Type: schema.TypeBool,
		},
		"scan_agent_id": &schema.Schema{
			Type: schema.TypeInt,
		},
		"include_in_ping": &schema.Schema{
			Type: schema.TypeBool,
		},
		"host_discovery_enabled": &schema.Schema{
			Type: schema.TypeBool,
		},
		"is_folder": &schema.Schema{
			Type: schema.TypeBool,
		},
		"is_full": &schema.Schema{
			Type: schema.TypeBool,
		},
		"utilization_threshold": &schema.Schema{
			Type: schema.TypeInt,
		},
		"location_id": &schema.Schema{
			Type: schema.TypeInt,
		},
		"edit_date": &schema.Schema{
			Type: schema.TypeString,
		},
		"gateway": &schema.Schema{
			Type: schema.TypeMap,
		},
		"gateway_id": &schema.Schema{
			Type: schema.TypeString,
		},
		"custom_fields": &schema.Schema{
			Type: schema.TypeMap,
		},
		"parent_subnet_id": &schema.Schema{
			Type: schema.TypeInt,
		},
	}
}

// resourceSubnetSchema returns the schema for the phpipam_subnet resource. It
// sets the required and optional fields, the latter defined in
// resourceSubnetRequiredFields, and ensures that all optional and
// non-configurable fields are computed as well.
func resourceSubnetSchema() map[string]*schema.Schema {
	schema := bareSubnetSchema()
	for k, v := range schema {
		switch {
		// Subnet Address and Mask are currently ForceNew
		case k == "subnet_address":
			v.Optional = true
			v.Computed = true
			v.ForceNew = true
			v.ConflictsWith = []string{"parent_subnet_id"}
		case k == "subnet_mask":
			v.Required = true
			v.ForceNew = true
		case k == "section_id":
			v.Required = true
		case k == "custom_fields":
			v.Optional = true
		case k == "parent_subnet_id":
			v.Optional = true
			v.ConflictsWith = []string{"subnet_id"}
		case resourceSubnetOptionalFields.Has(k):
			v.Optional = true
			v.Computed = true
		default:
			v.Computed = true
		}
	}
	return schema
}

// dataSourceSubnetSchema returns the schema for the phpipam_subnet data
// source. It sets the searchable fields and sets up the attribute conflicts
// between subnet/mask and subnet ID. It also ensures that all fields are
// computed as well.
func dataSourceSubnetSchema() map[string]*schema.Schema {
	s := bareSubnetSchema()
	for k, v := range s {
		switch k {
		case "subnet_address", "subnet_mask":
			v.Optional = true
			v.Computed = true
			v.ConflictsWith = []string{"subnet_id", "section_id", "description", "custom_field_filter"}
		case "subnet_id":
			v.Optional = true
			v.Computed = true
			v.ConflictsWith = []string{"subnet_address", "subnet_mask", "section_id", "description", "custom_field_filter"}
		case "section_id":
			v.Optional = true
			v.Computed = true
			v.ConflictsWith = []string{"subnet_id", "subnet_address", "subnet_mask"}
		case "description":
			v.Optional = true
			v.Computed = true
			v.ConflictsWith = []string{"subnet_id", "subnet_address", "subnet_mask", "description_match", "custom_field_filter"}
		default:
			v.Computed = true
		}
	}
	// Add the description_match item to the schema. This is a meta-parameter
	// that is not part of the API resource and exists to instruct PHPIPAM to
	// do a regex search on the description field of the subnet. This conflicts
	// with "description" and the other fields that description would normally
	// conflict with.
	s["description_match"] = subnetDescriptionMatchSchema([]string{"subnet_id", "subnet_address", "subnet_mask", "description", "custom_field_filter"})

	// Add the custom_field_filter item to the schema. This is a meta-parameter
	// that allows searching for a custom field value in the data source.
	s["custom_field_filter"] = customFieldFilterSchema([]string{"subnet_id", "subnet_address", "subnet_mask", "description", "description_match"})

	return s
}

// dataSourceSubnetsSchema returns the sub-schema for the phpipam_subnets data
// source. All this function does is set all fields as computed.
func dataSourceSubnetsSchema() map[string]*schema.Schema {
	schema := bareSubnetSchema()
	for _, v := range schema {
		v.Computed = true
	}
	return schema
}

// expandSubnet returns the subnets.Subnet structure for a
// phpiapm_subnet resource or data source. Depending on if we are dealing with
// the resource or data source, extra considerations may need to be taken.
func expandSubnet(d *schema.ResourceData) subnets.Subnet {
	s := subnets.Subnet{
		ID:             d.Get("subnet_id").(int),
		SubnetAddress:  d.Get("subnet_address").(string),
		Mask:           phpipam.JSONIntString(d.Get("subnet_mask").(int)),
		Description:    d.Get("description").(string),
		SectionID:      d.Get("section_id").(int),
		LinkedSubnet:   d.Get("linked_subnet_id").(int),
		VLANID:         d.Get("vlan_id").(int),
		VRFID:          d.Get("vrf_id").(int),
		MasterSubnetID: d.Get("master_subnet_id").(int),
		NameserverID:   d.Get("nameserver_id").(int),
		ShowName:       phpipam.BoolIntString(d.Get("show_name").(bool)),
		Permissions:    d.Get("permissions").(string),
		DNSRecursive:   phpipam.BoolIntString(d.Get("create_ptr_records").(bool)),
		DNSRecords:     phpipam.BoolIntString(d.Get("display_hostnames").(bool)),
		AllowRequests:  phpipam.BoolIntString(d.Get("allow_ip_requests").(bool)),
		ScanAgent:      d.Get("scan_agent_id").(int),
		PingSubnet:     phpipam.BoolIntString(d.Get("include_in_ping").(bool)),
		DiscoverSubnet: phpipam.BoolIntString(d.Get("host_discovery_enabled").(bool)),
		IsFolder:       phpipam.BoolIntString(d.Get("is_folder").(bool)),
		IsFull:         phpipam.BoolIntString(d.Get("is_full").(bool)),
		Threshold:      d.Get("utilization_threshold").(int),
		Location:       d.Get("location_id").(int),
		Gateway:        d.Get("gateway").(map[string]interface{}),
		GatewayID:      d.Get("gateway_id").(string),
	}

	return s
}

// flattenSubnet sets fields in a *schema.ResourceData with fields supplied by
// the input subnets.Subnet. This is used in read operations.
func flattenSubnet(s subnets.Subnet, d *schema.ResourceData) {
	d.SetId(strconv.Itoa(s.ID))
	d.Set("subnet_id", s.ID)
	d.Set("subnet_address", s.SubnetAddress)
	d.Set("subnet_mask", s.Mask)
	d.Set("description", s.Description)
	d.Set("section_id", s.SectionID)
	d.Set("linked_subnet_id", s.LinkedSubnet)
	d.Set("vlan_id", s.VLANID)
	d.Set("vrf_id", s.VRFID)
	d.Set("master_subnet_id", s.MasterSubnetID)
	d.Set("nameserver_id", s.NameserverID)
	d.Set("show_name", s.ShowName)
	d.Set("permissions", s.Permissions)
	d.Set("create_ptr_records", s.DNSRecursive)
	d.Set("display_hostnames", s.DNSRecords)
	d.Set("allow_ip_requests", s.AllowRequests)
	d.Set("scan_agent_id", s.ScanAgent)
	d.Set("include_in_ping", s.PingSubnet)
	d.Set("host_discovery_enabled", s.DiscoverSubnet)
	d.Set("is_folder", s.IsFolder)
	d.Set("is_full", s.IsFull)
	d.Set("utilization_threshold", s.Threshold)
	d.Set("location_id", s.Location)
	d.Set("edit_date", s.EditDate)
	d.Set("gateway", s.Gateway)
	d.Set("gateway_id", s.GatewayID)
}

// subnetDescriptionMatchSchema returns a *schema.Schema for description
// matching for subnet-related resources. The conflicting keys are populated by
// the passed in string slice.
func subnetDescriptionMatchSchema(conflicts []string) *schema.Schema {
	return &schema.Schema{
		Type:          schema.TypeString,
		Optional:      true,
		ConflictsWith: conflicts,
		ValidateFunc: func(v interface{}, k string) (ws []string, errors []error) {
			_, err := regexp.Compile(v.(string))
			if err != nil {
				errors = append(errors, err)
			}
			return
		},
	}
}

// subnetSearchInSection provides the subnet search functionality for both the
// phpipam_subnet and phpipam_subnets data sources, returning a
// []subnets.Subnet to the particular data source that is calling the function.
// From here it's up to the specific data source to determine what they want to
// do with the results (ie: reject it on matching nothing or more than one for
// the singular data source, or extracting the IDs for the plural one).
func subnetSearchInSection(d *schema.ResourceData, meta interface{}) ([]subnets.Subnet, error) {
	c := meta.(*ProviderPHPIPAMClient).subnetsController
	s := meta.(*ProviderPHPIPAMClient).sectionsController
	result := make([]subnets.Subnet, 0)

	v, err := s.GetSubnetsInSection(d.Get("section_id").(int))
	if err != nil {
		return result, err
	}
	if len(v) == 0 {
		return result, errors.New("No subnets were found in the supplied section")
	}
	for _, r := range v {
		switch {
		// Double-assert that we don't have empty strings in the conditionals
		// to ensure there there is no edge cases with matching zero values.
		case d.Get("description_match").(string) != "":
			// Don't trap error here because we should have already validated the regex via the ValidateFunc.
			if matched, _ := regexp.MatchString(d.Get("description_match").(string), r.Description); matched {
				result = append(result, r)
			}
		case d.Get("description").(string) != "" && r.Description == d.Get("description").(string):
			result = append(result, r)
		case len(d.Get("custom_field_filter").(map[string]interface{})) > 0:
			// Skip folders for now as there is issues pulling them down in the API.
			if r.IsFolder {
				continue
			}
			fields, err := c.GetSubnetCustomFields(r.ID)
			if err != nil {
				return result, err
			}
			search := d.Get("custom_field_filter").(map[string]interface{})
			if err != nil {
				return result, err
			}
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
