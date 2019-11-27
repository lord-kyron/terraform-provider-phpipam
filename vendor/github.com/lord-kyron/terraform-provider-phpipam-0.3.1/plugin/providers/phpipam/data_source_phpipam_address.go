package phpipam

import (
	"errors"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/lord-kyron/phpipam-sdk-go/controllers/addresses"
)

func dataSourcePHPIPAMAddress() *schema.Resource {
	return &schema.Resource{
		Read:   dataSourcePHPIPAMAddressRead,
		Schema: dataSourceAddressSchema(),
	}
}

func dataSourcePHPIPAMAddressRead(d *schema.ResourceData, meta interface{}) error {
	c := meta.(*ProviderPHPIPAMClient).addressesController
	out := make([]addresses.Address, 1)
	var err error
	// We need to determine how to get the address. An ID search takes priority,
	// and after that addresss.
	switch {
	case d.Get("address_id").(int) != 0:
		out[0], err = c.GetAddressByID(d.Get("address_id").(int))
		if err != nil {
			return err
		}
	case d.Get("ip_address").(string) != "":
		out, err = c.GetAddressesByIP(d.Get("ip_address").(string))
		if err != nil {
			return err
		}
	case d.Get("subnet_id").(int) != 0 && (d.Get("description").(string) != "" || d.Get("hostname").(string) != "" || len(d.Get("custom_field_filter").(map[string]interface{})) > 0):
		out, err = addressSearchInSubnet(d, meta)
		if err != nil {
			return err
		}
	default:
		return errors.New("No valid combination of parameters found - need one of address_id, ip_address, or subnet_id and (description|hostname|custom_field_filter)")
	}
	if len(out) != 1 {
		return errors.New("Your search returned zero or multiple results. Please correct your search and try again")
	}
	flattenAddress(out[0], d)
	fields, err := c.GetAddressCustomFields(out[0].ID)
	if err != nil {
		return err
	}
	trimMap(fields)
	if err := d.Set("custom_fields", fields); err != nil {
		return err
	}
	return nil
}
