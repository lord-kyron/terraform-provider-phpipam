package phpipam

import (
	"errors"
	"fmt"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

const testAccResourcePHPIPAML2DomainResourceName = "phpipam_l2domain.l2domain"
const testAccResourcePHPIPAML2DomainName = "tf-test-l2domain"
const testAccResourcePHPIPAML2DomainConfig = `
resource "phpipam_l2domain" "l2domain" {
  name        = "tf-test-l2domain"
  description = "Terraform test l2domain"
}
`

func TestAccResourcePHPIPAML2Domain(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			l2domainSweep("tf-test-l2domain", t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckResourcePHPIPAML2DomainDeleted,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccResourcePHPIPAML2DomainConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckResourcePHPIPAML2DomainCreated,
					resource.TestCheckResourceAttr("phpipam_l2domain.l2domain", "name", "tf-test-l2domain"),
					resource.TestCheckResourceAttr("phpipam_l2domain.l2domain", "description", "Terraform test l2domain"),
				),
			},
		},
	})
}

func testAccCheckResourcePHPIPAML2DomainCreated(s *terraform.State) error {
	r, ok := s.RootModule().Resources[testAccResourcePHPIPAML2DomainResourceName]
	if !ok {
		return fmt.Errorf("Resource name %s could not be found", testAccResourcePHPIPAML2DomainResourceName)
	}
	if r.Primary.ID == "" {
		return errors.New("No ID is set")
	}

	id, _ := strconv.Atoi(r.Primary.ID)

	c := testAccProvider.Meta().(*ProviderPHPIPAMClient).l2domainsController
	if _, err := c.GetL2DomainByID(id); err != nil {
		return err
	}
	return nil
}

func testAccCheckResourcePHPIPAML2DomainDeleted(s *terraform.State) error {
	c := testAccProvider.Meta().(*ProviderPHPIPAMClient).l2domainsController
	_, err := c.GetL2DomainByName(testAccResourcePHPIPAML2DomainName)
	switch {
	case err == nil:
		return errors.New("Expected error, got none")
	case err != nil && err.Error() != "Error from API (404): No results (filter applied)":
		return fmt.Errorf("Expected 404, got %s", err)
	}

	return nil
}
