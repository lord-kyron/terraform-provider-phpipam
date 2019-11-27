package phpipam

import (
	"fmt"
	"reflect"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
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
