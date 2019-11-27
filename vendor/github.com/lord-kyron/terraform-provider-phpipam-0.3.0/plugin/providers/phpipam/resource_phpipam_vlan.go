package phpipam

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/hashicorp/terraform/helper/schema"
)

// resourcePHPIPAMVLAN returns the resource structure for the phpipam_vlan
// resource.
//
// Note that we use the data source read function here to pull down data, as
// read workflow is identical for both the resource and the data source.
func resourcePHPIPAMVLAN() *schema.Resource {
	return &schema.Resource{
		Create: resourcePHPIPAMVLANCreate,
		Read:   dataSourcePHPIPAMVLANRead,
		Update: resourcePHPIPAMVLANUpdate,
		Delete: resourcePHPIPAMVLANDelete,
		Schema: resourceVLANSchema(),
	}
}

func resourcePHPIPAMVLANCreate(d *schema.ResourceData, meta interface{}) error {
	c := meta.(*ProviderPHPIPAMClient).vlansController
	in := expandVLAN(d)

	// Assert the ID field here is empty. If this is not empty the request will fail.
	in.ID = 0

	if _, err := c.CreateVLAN(in); err != nil {
		return err
	}

	// If we have custom fields, set them now. We need to get the IP address's ID
	// beforehand.
	if customFields, ok := d.GetOk("custom_fields"); ok {
		vlans, err := c.GetVLANsByNumber(in.Number)
		if err != nil {
			return fmt.Errorf("Could not read VLAN after creating: %s", err)
		}

		if len(vlans) != 1 {
			return errors.New("VLAN either missing or multiple results returned by reading VLAN after creation")
		}

		d.SetId(strconv.Itoa(vlans[0].ID))

		if _, err := c.UpdateVLANCustomFields(vlans[0].ID, vlans[0].Name, customFields.(map[string]interface{})); err != nil {
			return err
		}
	}

	return dataSourcePHPIPAMVLANRead(d, meta)
}

func resourcePHPIPAMVLANUpdate(d *schema.ResourceData, meta interface{}) error {
	c := meta.(*ProviderPHPIPAMClient).vlansController
	in := expandVLAN(d)

	if _, err := c.UpdateVLAN(in); err != nil {
		return err
	}

	if err := updateCustomFields(d, c); err != nil {
		return err
	}

	return dataSourcePHPIPAMVLANRead(d, meta)
}

func resourcePHPIPAMVLANDelete(d *schema.ResourceData, meta interface{}) error {
	c := meta.(*ProviderPHPIPAMClient).vlansController
	in := expandVLAN(d)

	if _, err := c.DeleteVLAN(in.ID); err != nil {
		return err
	}
	d.SetId("")
	return nil
}
