package phpipam

import (
	"errors"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/pavel-z1/phpipam-sdk-go/controllers/l2domains"
)

func dataSourcePHPIPAML2Domain() *schema.Resource {
	return &schema.Resource{
		Read:   dataSourcePHPIPAML2DomainRead,
		Schema: dataSourceL2DomainSchema(),
	}
}

func dataSourcePHPIPAML2DomainRead(d *schema.ResourceData, meta interface{}) error {
	c := meta.(*ProviderPHPIPAMClient).l2domainsController
	var out l2domains.L2Domain
	var err error
	// We need to determine how to get the l2domain. An ID search takes priority,
	// and after that l2domain name.
	switch {
	case d.Get("domain_id").(int) != 0:
		out, err = c.GetL2DomainByID(d.Get("domain_id").(int))
		if err != nil {
			if strings.Contains(err.Error(), "Invalid domain id") {
				log.Printf("Can't find l2domain with id %d", d.Get("domain_id").(int))
				return nil
			}
			return err
		}
	case d.Get("name").(string) != "":
		list_out, err := c.GetL2DomainByName(d.Get("name").(string))
		if err != nil {
			if strings.Contains(err.Error(), "No results (filter applied)") {
				log.Printf("Can't find l2domain with name %s", d.Get("name").(string))
				return nil
			}
			return err
		}
		if len(list_out) != 1 {
			return errors.New("L2 Domain either missing or multiple results returned by reading l2domains")
		}
		out = list_out[0]
	default:
		return errors.New("domain_id or name not defined, cannot proceed with reading data")
	}
	flattenL2Domain(out, d)
	return nil
}
