package phpipam

import "github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

// resourcePHPIPAML2Domain returns the resource structure for the phpipam_l2domain
// resource.
//
// Note that we use the data source read function here to pull down data, as
// read workflow is identical for both the resource and the data source.
func resourcePHPIPAML2Domain() *schema.Resource {
	return &schema.Resource{
		Create: resourcePHPIPAML2DomainCreate,
		Read:   dataSourcePHPIPAML2DomainRead,
		Update: resourcePHPIPAML2DomainUpdate,
		Delete: resourcePHPIPAML2DomainDelete,
		Schema: resourceL2DomainSchema(),
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

func resourcePHPIPAML2DomainCreate(d *schema.ResourceData, meta interface{}) error {
	c := meta.(*ProviderPHPIPAMClient).l2domainsController
	in := expandL2Domain(d)

	// Assert the ID field here is empty. If this is not empty the request will fail.
	in.ID = 0

	if _, err := c.CreateL2Domain(in); err != nil {
		return err
	}

	return dataSourcePHPIPAML2DomainRead(d, meta)
}

func resourcePHPIPAML2DomainUpdate(d *schema.ResourceData, meta interface{}) error {
	c := meta.(*ProviderPHPIPAMClient).l2domainsController
	in := expandL2Domain(d)

	if err := c.UpdateL2Domain(in); err != nil {
		return err
	}

	return dataSourcePHPIPAML2DomainRead(d, meta)
}

func resourcePHPIPAML2DomainDelete(d *schema.ResourceData, meta interface{}) error {
	c := meta.(*ProviderPHPIPAMClient).l2domainsController
	in := expandL2Domain(d)

	if err := c.DeleteL2Domain(in.ID); err != nil {
		return err
	}
	d.SetId("")
	return nil
}
