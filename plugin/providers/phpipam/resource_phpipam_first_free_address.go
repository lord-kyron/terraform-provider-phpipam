package phpipam

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// resourcePHPIPAMAddress returns the resource structure for the phpipam_address
// resource.
//
// Note that we use the data source read function here to pull down data, as
// read workflow is identical for both the resource and the data source.
func resourcePHPIPAMFirstFreeAddress() *schema.Resource {
	return &schema.Resource{
		Create: resourcePHPIPAMFirstFreeAddressCreate,
		Read:   dataSourcePHPIPAMAddressRead,
		Update: resourcePHPIPAMFirstFreeAddressUpdate,
		Delete: resourcePHPIPAMFirstFreeAddressDelete,
		Schema: resourceFirstFreeAddressSchema(),
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
		// IP Address and Subnet ID are ForceNew
		case k == "subnet_id":
			v.Required = true
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
	return s
}

func resourcePHPIPAMFirstFreeAddressCreate(d *schema.ResourceData, meta interface{}) error {
	// Get first free IP from provided subnet_id
	subnet_id := d.Get("subnet_id").(int)
	d.Set("subnet_id", nil)

	// Get address controller and start address creation
	c := meta.(*ProviderPHPIPAMClient).addressesController

	in := expandAddress(d)

	out, err := c.CreateFirstFreeAddress(subnet_id, in)
	if err != nil {
		return err
	}
	d.Set("ip_address", out)

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
