package phpipam

import (
	"errors"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/pavel-z1/phpipam-sdk-go/controllers/vlans"
)

func dataSourcePHPIPAMVLAN() *schema.Resource {
	return &schema.Resource{
		Read:   dataSourcePHPIPAMVLANRead,
		Schema: dataSourceVLANSchema(),
	}
}

func dataSourcePHPIPAMVLANRead(d *schema.ResourceData, meta interface{}) error {
	c := meta.(*ProviderPHPIPAMClient).vlansController
	var out vlans.VLAN
	// We need to determine how to get the vlan. An ID search takes priority,
	// and after that vlans.
	switch {
	case d.Get("vlan_id").(int) != 0:
		var err error
		out, err = c.GetVLANByID(d.Get("vlan_id").(int))
		if err != nil {
			return err
		}
	case d.Get("number").(int) != 0 && d.Get("l2_domain_id").(int) != 0:
		v, err := c.GetVLANsByNumberAndDomainID(d.Get("number").(int), d.Get("l2_domain_id").(int))
		if err != nil {
			return err
		}
		// Only one result should be returned by this search. Fail on multiples.
		if len(v) != 1 {
			return errors.New("VLAN search returned either zero or multiple results. Please correct your search and try again")
		}
		out = v[0]
	case d.Get("number").(int) != 0:
		v, err := c.GetVLANsByNumber(d.Get("number").(int))
		if err != nil {
			return err
		}
		// Only one result should be returned by this search. Fail on multiples.
		if len(v) != 1 {
			return errors.New("VLAN search returned either zero or multiple results. Please correct your search and try again")
		}
		out = v[0]
	default:
		return errors.New("vlan_id or number not defined, cannot proceed with reading data")
	}

	fields, err := c.GetVLANCustomFields(out.ID)
	switch {
	case err == nil:
		trimMap(fields)
		if err := d.Set("custom_fields", fields); err != nil {
			return err
		}
	case err != nil:
		return err
	}

	flattenVLAN(out, d)
	return nil
}
