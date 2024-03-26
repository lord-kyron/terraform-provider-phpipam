package phpipam

import (
	"errors"
	"log"
	"strings"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/pavel-z1/phpipam-sdk-go/controllers/sections"
)

func dataSourcePHPIPAMSection() *schema.Resource {
	return &schema.Resource{
		Read:   dataSourcePHPIPAMSectionRead,
		Schema: dataSourceSectionSchema(),
	}
}

func dataSourcePHPIPAMSectionRead(d *schema.ResourceData, meta interface{}) error {
	c := meta.(*ProviderPHPIPAMClient).sectionsController
	var out sections.Section
	var err error
	// We need to determine how to get the section. An ID search takes priority,
	// and after that section name.
	switch {
	case d.Get("section_id").(int) != 0:
		out, err = c.GetSectionByID(d.Get("section_id").(int))
		if err != nil {
			if strings.Contains(err.Error(), "Section does not exist") {
				log.Printf("Can't find section with id %d", d.Get("section_id").(int))
				return nil
			}
			return err
		}
	case d.Get("name").(string) != "":
		out, err = c.GetSectionByName(d.Get("name").(string))
		if err != nil {
			if strings.Contains(err.Error(), "Not Found") {
				log.Printf("Can't find section with name %s", d.Get("name").(string))
				return nil
			}
			return err
		}
	default:
		// We need to ensure imported resources are not recreated when terraform apply is ran
		// imported resources only have an Id which we need to map back to section_id
		id := d.Id()
		if len(id) > 0 {
			section_id, err := strconv.Atoi(id)
			if err != nil {
				return err
			}
			out, err = c.GetSectionByID(section_id)
			if err != nil {
				return err
			}
		} else {
			return errors.New("section_id or name not defined, cannot proceed with reading data")
		}
	}
	flattenSection(out, d)
	return nil
}
