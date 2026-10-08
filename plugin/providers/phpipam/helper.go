package phpipam

import "github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

// linearSearchSlice provides a []string with a helper search function.
type linearSearchSlice []string

// Has checks the linearSearchSlice for the string provided by x and returns
// true if it finds a match.
func (s *linearSearchSlice) Has(x string) bool {
	for _, v := range *s {
		if v == x {
			return true
		}
	}
	return false
}

// subnetIDsFromResourceData returns the ordered list of subnet IDs to try,
// supporting both the legacy singular "subnet_id" and the "subnet_ids" list.
func subnetIDsFromResourceData(d *schema.ResourceData) []int {
	if raw, ok := d.GetOk("subnet_ids"); ok {
		list := raw.([]interface{})
		ids := make([]int, len(list))
		for i, v := range list {
			ids[i] = v.(int)
		}
		return ids
	}
	return []int{d.Get("subnet_id").(int)}
}
