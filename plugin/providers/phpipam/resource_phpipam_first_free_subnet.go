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
func resourcePHPIPAMFirstFreeSubnet() *schema.Resource {
	return &schema.Resource{
		Create: resourcePHPIPAMFirstFreeSubnetCreate,
		Read:   dataSourcePHPIPAMSubnetRead,
		Update: resourcePHPIPAMFirstFreeSubnetUpdate,
		Delete: resourcePHPIPAMFirstFreeSubnetDelete,
		Schema: resourceFirstFreeSubnetSchema(),
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

func resourcePHPIPAMFirstFreeSubnetCreate(d *schema.ResourceData, meta interface{}) error {
	// Get first free subnet from provided subnet_id
	subnet_id := d.Get("parent_subnet_id").(int)
	d.Set("subnet_id", nil)
	subnet_mask := d.Get("subnet_mask").(int)
	// Get address controller and start address creation
	c := meta.(*ProviderPHPIPAMClient).subnetsController

	in := expandSubnet(d)

	out, err := c.CreateFirstFreeSubnet(subnet_id, subnet_mask, in)
	if err != nil {
		return err
	}
	d.Set("subnet_address", out)

	// If we have custom fields, set them now. We need to get the IP address's ID
	// beforehand.
	if customFields, ok := d.GetOk("custom_fields"); ok {
		addrs, err := c.GetSubnetsByCIDR(fmt.Sprintf("%s/%d", out, in.Mask))
		if err != nil {
			return fmt.Errorf("Could not read IP address after creating: %s", err)
		}

		if len(addrs) != 1 {
			return errors.New("IP address either missing or multiple results returned by reading IP after creation")
		}

		d.SetId(strconv.Itoa(addrs[0].ID))

		if _, err := c.UpdateSubnetCustomFields(addrs[0].ID, customFields.(map[string]interface{})); err != nil {
			return err
		}
	}

	return dataSourcePHPIPAMSubnetRead(d, meta)
}

func resourcePHPIPAMFirstFreeSubnetUpdate(d *schema.ResourceData, meta interface{}) error {
	c := meta.(*ProviderPHPIPAMClient).subnetsController
	in := expandSubnet(d)

	// SubnetAddress and mask need to be removed for update requests.
	in.SubnetAddress = ""
	in.Mask = 0
	if _, err := c.UpdateSubnet(in); err != nil {
		return err
	}

	if err := updateCustomFields(d, c); err != nil {
		return err
	}

	return dataSourcePHPIPAMSubnetRead(d, meta)
}

func resourcePHPIPAMFirstFreeSubnetDelete(d *schema.ResourceData, meta interface{}) error {
	c := meta.(*ProviderPHPIPAMClient).subnetsController
	in := expandSubnet(d)

	if _, err := c.DeleteSubnet(in.ID); err != nil {
		return err
	}
	d.SetId("")
	return nil
}
