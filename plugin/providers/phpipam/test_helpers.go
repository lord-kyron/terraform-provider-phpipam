package phpipam

import (
	"fmt"
	"reflect"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

// testCheckOutputPair compares two outputs, like TestCheckResourceAttrPair for
// resources.
func testCheckOutputPair(nameFirst, nameSecond string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		ms := s.RootModule()
		rsFirst, ok := ms.Outputs[nameFirst]
		if !ok {
			return fmt.Errorf("Output not found: %s", nameFirst)
		}
		rsSecond, ok := ms.Outputs[nameSecond]
		if !ok {
			return fmt.Errorf("Output not found: %s", nameSecond)
		}

		if !reflect.DeepEqual(rsFirst.Value, rsSecond.Value) {
			return fmt.Errorf(
				"Output %q (value %#v) did not match output %q (value %#v)",
				nameFirst,
				rsFirst.Value,
				nameSecond,
				rsSecond.Value)
		}

		return nil
	}
}

// testCheckResourceIDSave stores the ID of a resource, so that a later step can
// check with testCheckResourceIDUnchanged that it was updated in place.
func testCheckResourceIDSave(name string, id *string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		r, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("Resource not found: %s", name)
		}
		*id = r.Primary.ID
		return nil
	}
}

// testCheckResourceIDUnchanged fails if a resource got a new ID, i.e. it was
// replaced instead of updated in place.
func testCheckResourceIDUnchanged(name string, id *string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		r, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("Resource not found: %s", name)
		}
		if r.Primary.ID != *id {
			return fmt.Errorf("%s was replaced: ID changed from %s to %s", name, *id, r.Primary.ID)
		}
		return nil
	}
}
