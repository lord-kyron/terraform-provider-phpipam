package phpipam

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/pavel-z1/phpipam-sdk-go/controllers/subnets"
)

// resourcePHPIPAMSubnet returns the resource structure for the phpipam_subnet
// resource.
//
// Note that we use the data source read function here to pull down data, as
// read workflow is identical for both the resource and the data source.
func resourcePHPIPAMSubnet() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourcePHPIPAMSubnetCreate,
		ReadContext:   dataSourcePHPIPAMSubnetRead,
		UpdateContext: resourcePHPIPAMSubnetUpdate,
		DeleteContext: resourcePHPIPAMSubnetDelete,
		Schema:        resourceSubnetSchema(),
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

func resourcePHPIPAMSubnetCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*ProviderPHPIPAMClient).subnetsController
	in := expandSubnet(d, meta.(*ProviderPHPIPAMClient).NestCustomFields)

	// Assert the ID field here is empty. If this is not empty the request will fail.
	in.ID = 0

	if _, ok := d.GetOk("subnet_address"); ok {
		if _, err := c.CreateSubnet(in); err != nil {
			return diag.FromErr(err)
		}
	} else if parentSubnetId, ok := d.GetOk("parent_subnet_id"); ok {
		res, err := c.CreateFirstFreeSubnet(parentSubnetId.(int), d.Get("subnet_mask").(int), in)

		if err != nil {
			return diag.FromErr(err)
		}

		netAndMask := strings.Split(res, "/")
		maskNum, _ := strconv.Atoi(netAndMask[1])

		d.Set("subnet_address", netAndMask[0])
		d.Set("subnet_mask", maskNum)
	} else {
		return diag.FromErr(errors.New("Unsupported scenario! One of 'subnet_address' or 'parent_subnet_id' must be set"))
	}

	if !meta.(*ProviderPHPIPAMClient).NestCustomFields {
		// If we have custom fields, set them now. We need to get the subnet's ID
		// beforehand.
		if customFields, ok := d.GetOk("custom_fields"); ok {
			var subnets []subnets.Subnet
			var err error
			switch {
			case in.SectionID != 0:
				subnets, err = c.GetSubnetsByCIDRAndSection(fmt.Sprintf("%s/%d", in.SubnetAddress, in.Mask), in.SectionID)
			default:
				subnets, err = c.GetSubnetsByCIDR(fmt.Sprintf("%s/%d", in.SubnetAddress, in.Mask))
			}
			if err != nil {
				return diag.FromErr(fmt.Errorf("Could not read subnet after creating: %s", err))
			}

			if len(subnets) != 1 {
				return diag.FromErr(errors.New("Subnet either missing or multiple results returned by reading subnet after creation"))
			}

			d.SetId(strconv.Itoa(subnets[0].ID))

			if _, err := c.UpdateSubnetCustomFields(subnets[0].ID, customFields.(map[string]interface{})); err != nil {
				return diag.FromErr(err)
			}
		}
	}

	return dataSourcePHPIPAMSubnetRead(ctx, d, meta)
}

func resourcePHPIPAMSubnetUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*ProviderPHPIPAMClient).subnetsController
	in := expandSubnet(d, meta.(*ProviderPHPIPAMClient).NestCustomFields)
	// Remove the CIDR fields from the request, as these fields being present
	// implies that the subnet will be either split or renamed, which is not
	// supported by UpdateSubnet. These are implemented in the API but not in the
	// SDK, so support may be added at a later time.
	in.SubnetAddress = ""
	in.Mask = 0
	if _, err := c.UpdateSubnet(in); err != nil {
		return diag.FromErr(err)
	}

	if !meta.(*ProviderPHPIPAMClient).NestCustomFields {
		if err := updateCustomFields(d, c); err != nil {
			return diag.FromErr(err)
		}
	}

	return dataSourcePHPIPAMSubnetRead(ctx, d, meta)
}

func resourcePHPIPAMSubnetDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*ProviderPHPIPAMClient).subnetsController
	in := expandSubnet(d, meta.(*ProviderPHPIPAMClient).NestCustomFields)

	if _, err := c.DeleteSubnet(in.ID); err != nil {
		return diag.FromErr(err)
	}
	d.SetId("")
	return nil
}
