package phpipam

import "github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

// resourcePHPIPAMSection returns the resource structure for the phpipam_section
// resource.
//
// Note that we use the data source read function here to pull down data, as
// read workflow is identical for both the resource and the data source.
func resourcePHPIPAMSection() *schema.Resource {
	return &schema.Resource{
		Create: resourcePHPIPAMSectionCreate,
		Read:   dataSourcePHPIPAMSectionRead,
		Update: resourcePHPIPAMSectionUpdate,
		Delete: resourcePHPIPAMSectionDelete,
		Schema: resourceSectionSchema(),
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

func resourcePHPIPAMSectionCreate(d *schema.ResourceData, meta interface{}) error {
	c := meta.(*ProviderPHPIPAMClient).sectionsController
	in := expandSection(d)

	// Assert the ID field here is empty. If this is not empty the request will fail.
	in.ID = 0

	if _, err := c.CreateSection(in); err != nil {
		return err
	}

	return dataSourcePHPIPAMSectionRead(d, meta)
}

func resourcePHPIPAMSectionUpdate(d *schema.ResourceData, meta interface{}) error {
	c := meta.(*ProviderPHPIPAMClient).sectionsController
	in := expandSection(d)

	if err := c.UpdateSection(in); err != nil {
		return err
	}

	return dataSourcePHPIPAMSectionRead(d, meta)
}

func resourcePHPIPAMSectionDelete(d *schema.ResourceData, meta interface{}) error {
	c := meta.(*ProviderPHPIPAMClient).sectionsController
	in := expandSection(d)

	if err := c.DeleteSection(in.ID); err != nil {
		return err
	}
	d.SetId("")
	return nil
}
