package phpipam

import (
	"errors"
	"fmt"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

const testAccResourcePHPIPAMSectionResourceName = "phpipam_section.section"
const testAccResourcePHPIPAMSectionName = "tf-test"
const testAccResourcePHPIPAMSectionConfig = `
resource "phpipam_section" "section" {
  name        = "tf-test"
  description = "Terraform test section"
}
`

func TestAccResourcePHPIPAMSection(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckResourcePHPIPAMSectionDeleted,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccResourcePHPIPAMSectionConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckResourcePHPIPAMSectionCreated,
					resource.TestCheckResourceAttr("phpipam_section.section", "name", "tf-test"),
					resource.TestCheckResourceAttr("phpipam_section.section", "description", "Terraform test section"),
				),
			},
		},
	})
}

func testAccCheckResourcePHPIPAMSectionCreated(s *terraform.State) error {
	r, ok := s.RootModule().Resources[testAccResourcePHPIPAMSectionResourceName]
	if !ok {
		return fmt.Errorf("Resource name %s could not be found", testAccResourcePHPIPAMSectionResourceName)
	}
	if r.Primary.ID == "" {
		return errors.New("No ID is set")
	}

	id, _ := strconv.Atoi(r.Primary.ID)

	c := testAccProvider.Meta().(*ProviderPHPIPAMClient).sectionsController
	if _, err := c.GetSectionByID(id); err != nil {
		return err
	}
	return nil
}

func testAccCheckResourcePHPIPAMSectionDeleted(s *terraform.State) error {
	c := testAccProvider.Meta().(*ProviderPHPIPAMClient).sectionsController
	_, err := c.GetSectionByName(testAccResourcePHPIPAMSectionName)
	switch {
	case err == nil:
		return errors.New("Expected error, got none")
	case err != nil && err.Error() != "Error from API (404): Not Found":
		return fmt.Errorf("Expected 404, got %s", err)
	}

	return nil
}
