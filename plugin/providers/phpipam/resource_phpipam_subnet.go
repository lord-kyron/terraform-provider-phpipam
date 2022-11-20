package phpipam

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// resourcePHPIPAMSubnet returns the resource structure for the phpipam_subnet
// resource.
//
// Note that we use the data source read function here to pull down data, as
// read workflow is identical for both the resource and the data source.
func resourcePHPIPAMSubnet() *schema.Resource {
	return &schema.Resource{
		Create: resourcePHPIPAMSubnetCreate,
		Read:   dataSourcePHPIPAMSubnetRead,
		Update: resourcePHPIPAMSubnetUpdate,
		Delete: resourcePHPIPAMSubnetDelete,
		Schema: resourceSubnetSchema(),
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

func resourcePHPIPAMSubnetCreate(d *schema.ResourceData, meta interface{}) error {
	c := meta.(*ProviderPHPIPAMClient).subnetsController
	in := expandSubnet(d)

	// Assert the ID field here is empty. If this is not empty the request will fail.
	in.ID = 0

	if _, ok := d.GetOk("subnet_address"); ok {
		if _, err := c.CreateSubnet(in); err != nil {
			return err
		}
	} else if parentSubnetId, ok := d.GetOk("parent_subnet_id"); ok {
		res, err := c.CreateFirstFreeSubnet(parentSubnetId.(int), d.Get("subnet_mask").(int), in)

		if err != nil {
			return err
		}

		netAndMask := strings.Split(res, "/")
		maskNum, _ := strconv.Atoi(netAndMask[1])

		d.Set("subnet_address", netAndMask[0])
		d.Set("subnet_mask", maskNum)
	} else {
		return errors.New("Unsupported scenario! One of 'subnet_address' or 'parent_subnet_id' must be set")
	}

	// If we have custom fields, set them now. We need to get the subnet's ID
	// beforehand.
	if customFields, ok := d.GetOk("custom_fields"); ok {
		subnets, err := c.GetSubnetsByCIDR(fmt.Sprintf("%s/%d", in.SubnetAddress, in.Mask))
		if err != nil {
			return fmt.Errorf("Could not read subnet after creating: %s", err)
		}

		if len(subnets) != 1 {
			return errors.New("Subnet either missing or multiple results returned by reading subnet after creation")
		}

		d.SetId(strconv.Itoa(subnets[0].ID))

		if _, err := c.UpdateSubnetCustomFields(subnets[0].ID, customFields.(map[string]interface{})); err != nil {
			return err
		}
	}

	return dataSourcePHPIPAMSubnetRead(d, meta)
}

func resourcePHPIPAMSubnetUpdate(d *schema.ResourceData, meta interface{}) error {
	c := meta.(*ProviderPHPIPAMClient).subnetsController
	in := expandSubnet(d)
	// Remove the CIDR fields from the request, as these fields being present
	// implies that the subnet will be either split or renamed, which is not
	// supported by UpdateSubnet. These are implemented in the API but not in the
	// SDK, so support may be added at a later time.
	in.SubnetAddress = ""
	in.Mask = 0
	if _, err := c.UpdateSubnet(in); err != nil {
		return err
	}

	if len(d.Get("custom_fields").(map[string]interface{})) != 0 {
		if err := updateCustomFields(d, c); err != nil {
			return err
		}
	}

	return dataSourcePHPIPAMSubnetRead(d, meta)
}

func resourcePHPIPAMSubnetDelete(d *schema.ResourceData, meta interface{}) error {
	c := meta.(*ProviderPHPIPAMClient).subnetsController
	in := expandSubnet(d)

	if _, err := c.DeleteSubnet(in.ID); err != nil {
		return err
	}
	d.SetId("")
	return nil
}
