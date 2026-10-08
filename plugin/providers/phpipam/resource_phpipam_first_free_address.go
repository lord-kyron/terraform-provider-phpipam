package phpipam

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// resourcePHPIPAMAddress returns the resource structure for the phpipam_address
// resource.
//
// Note that we use the data source read function here to pull down data, as
// read workflow is identical for both the resource and the data source.
func resourcePHPIPAMFirstFreeAddress() *schema.Resource {
	return &schema.Resource{
		Create:        resourcePHPIPAMFirstFreeAddressCreate,
		Read:          dataSourcePHPIPAMAddressRead,
		Update:        resourcePHPIPAMFirstFreeAddressUpdate,
		Delete:        resourcePHPIPAMFirstFreeAddressDelete,
		CustomizeDiff: resourcePHPIPAMFirstFreeAddressCustomizeDiff,
		Schema:        resourceFirstFreeAddressSchema(),
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

// bareAddressSchema returns a map[string]*schema.Schema with the schema used
// to represent a PHPIPAM address resource. This output should then be modified
// so that required and computed fields are set properly for both the data
// source and the resource.
func bareFirstFreeAddressSchema() map[string]*schema.Schema {
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
// resourceAddressRequiredFields, and ensures that all optional and
// non-configurable fields are computed as well.
func resourceFirstFreeAddressSchema() map[string]*schema.Schema {
	s := bareAddressSchema()
	for k, v := range s {
		switch {
		// Subnet ID is ForceNew, and is mutually exclusive with subnet_ids.
		// It stays Computed so it reflects the subnet actually used once an
		// address has been allocated from subnet_ids.
		case k == "subnet_id":
			v.Optional = true
			v.Computed = true
			v.ForceNew = true
			v.ExactlyOneOf = []string{"subnet_id", "subnet_ids"}
		case k == "custom_fields":
			v.Optional = true
		case resourceAddressOptionalFields.Has(k):
			v.Optional = true
			v.Computed = true
		default:
			v.Computed = true
		}
	}
	// subnet_ids allows picking the first subnet, in order, that still has a
	// free address, instead of a single mandatory subnet_id. It is Computed so
	// resourcePHPIPAMFirstFreeAddressCustomizeDiff can clear its diff when
	// migrating from subnet_id without forcing a replacement.
	s["subnet_ids"] = &schema.Schema{
		Type:         schema.TypeList,
		Optional:     true,
		Computed:     true,
		ForceNew:     true,
		Elem:         &schema.Schema{Type: schema.TypeInt},
		ExactlyOneOf: []string{"subnet_id", "subnet_ids"},
	}
	return s
}

// resourcePHPIPAMFirstFreeAddressCustomizeDiff prevents a spurious replacement
// when a config moves from a single subnet_id to a subnet_ids list that still
// contains the subnet the address already lives in - the address itself does
// not need to move, so there is nothing to force a new resource for.
func resourcePHPIPAMFirstFreeAddressCustomizeDiff(ctx context.Context, d *schema.ResourceDiff, meta interface{}) error {
	oldSubnetIDRaw, _ := d.GetChange("subnet_id")
	currentSubnetID := oldSubnetIDRaw.(int)
	if currentSubnetID == 0 {
		return nil
	}

	// subnet_ids is Optional+Computed, so GetOk alone would happily fall back
	// to its last known state value even when the config only sets subnet_id
	// (e.g. when genuinely moving to a different subnet). Only reconcile when
	// subnet_ids is actually present in the raw config.
	rawConfig := d.GetRawConfig()
	if rawConfig.IsNull() || !rawConfig.IsKnown() {
		return nil
	}
	subnetIDsInConfig := rawConfig.GetAttr("subnet_ids")
	if subnetIDsInConfig.IsNull() || !subnetIDsInConfig.IsKnown() {
		return nil
	}

	rawIDs, ok := d.GetOk("subnet_ids")
	if !ok {
		return nil
	}

	for _, v := range rawIDs.([]interface{}) {
		if v.(int) == currentSubnetID {
			if err := d.Clear("subnet_id"); err != nil {
				return err
			}
			return d.Clear("subnet_ids")
		}
	}

	return nil
}

func resourcePHPIPAMFirstFreeAddressCreate(d *schema.ResourceData, meta interface{}) error {
	// Try subnet_ids (or the single subnet_id) in order until one has a free address.
	ids := subnetIDsFromResourceData(d)
	d.Set("subnet_id", nil)

	// Get address controller and start address creation
	c := meta.(*ProviderPHPIPAMClient).addressesController

	in := expandAddress(d)

	var out string
	var err error
	found := false
	for _, id := range ids {
		out, err = c.CreateFirstFreeAddress(id, in)
		switch {
		case err != nil && strings.Contains(err.Error(), "No free addresses found"):
			continue
		case err != nil:
			return err
		}
		found = true
		break
	}
	if !found {
		return fmt.Errorf("No free IP addresses found in any of the provided subnets: %s", err)
	}
	d.Set("ip_address", out)
	// Persist the candidate list so future plans see subnet_ids as already
	// known and don't treat it as still-to-be-computed (which, combined with
	// ForceNew, would force a replacement on every plan).
	d.Set("subnet_ids", ids)

	// If we have custom fields, set them now. We need to get the IP address's ID
	// beforehand.
	if customFields, ok := d.GetOk("custom_fields"); ok {
		addrs, err := c.GetAddressesByIP(out)
		if err != nil {
			return fmt.Errorf("Could not read IP address after creating: %s", err)
		}
		//addrs := d.Get("ip_address")

		if len(addrs) != 1 {
			return errors.New("IP address either missing or multiple results returned by reading IP after creation")
		}

		d.SetId(strconv.Itoa(addrs[0].ID))

		if _, err := c.UpdateAddressCustomFields(addrs[0].ID, customFields.(map[string]interface{})); err != nil {
			return err
		}
	}

	return dataSourcePHPIPAMAddressRead(d, meta)
}

func resourcePHPIPAMFirstFreeAddressUpdate(d *schema.ResourceData, meta interface{}) error {
	c := meta.(*ProviderPHPIPAMClient).addressesController
	in := expandAddress(d)

	// IPAddress and SubnetID need to be removed for update requests.
	in.IPAddress = ""
	in.SubnetID = 0
	if _, err := c.UpdateAddress(in); err != nil {
		return err
	}

	if err := updateCustomFields(d, c); err != nil {
		return err
	}

	return dataSourcePHPIPAMAddressRead(d, meta)
}

func resourcePHPIPAMFirstFreeAddressDelete(d *schema.ResourceData, meta interface{}) error {
	c := meta.(*ProviderPHPIPAMClient).addressesController
	in := expandAddress(d)

	//	if _, err := c.DeleteAddress(in.ID, phpipam.BoolIntString(d.Get("remove_dns_on_delete").(bool))); err != nil {
	if _, err := c.DeleteAddress(in.ID, false); err != nil {
		return err
	}
	d.SetId("")
	return nil
}
