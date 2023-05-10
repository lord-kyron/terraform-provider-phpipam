package phpipam

import (
	"fmt"
	"log"
	"reflect"
	"regexp"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/pavel-z1/phpipam-sdk-go/controllers/addresses"
	"github.com/pavel-z1/phpipam-sdk-go/controllers/subnets"
	"github.com/pavel-z1/phpipam-sdk-go/controllers/vlans"
)

// customFieldFilterSchema returns a *schema.Schema for the custom_field_filter
// attribute on select data sources. The conflict keys are populated by the
// supplied string slice.
func customFieldFilterSchema(conflicts []string) *schema.Schema {
	return &schema.Schema{
		Type:          schema.TypeMap,
		Optional:      true,
		ConflictsWith: conflicts,
		ValidateFunc: func(m interface{}, k string) (ws []string, errors []error) {
			for _, v := range m.(map[string]interface{}) {
				_, err := regexp.Compile(v.(string))
				if err != nil {
					errors = append(errors, err)
				}
			}
			return
		},
	}
}

// customFieldFilter takes two maps - one with custom field data, and one with
// key/value search items. The function returns a true match only if all of the
// search fields match. The data is matched against value as a regex. For exact
// matching, ensure your match is enclosed in the ^ (start of line) and the $
// (end of line) anchors.
//
// PHPIPAM currently stringifies most, if not all, values coming out of the
// API. As such, we don't attempt to cast here - anything that is not a string
// is an error. If the need arises for this to be changed at some point in time
// this function will be updated.
func customFieldFilter(data, search map[string]interface{}) (bool, error) {
	// zero-length or nil map is a panic. This should never happen
	if search == nil || len(search) == 0 {
		panic("Zero length or nil map passed as search terms to customFieldFilter")
	}

	for k, expr := range search {
		if w, ok := data[k]; ok {
			switch v := w.(type) {
			case nil:
				// no field value, not a match
				return false, nil
			case string:
				if match, _ := regexp.MatchString(expr.(string), v); !match {
					// not a match if one of the search values do not match
					return false, nil
				}
			default:
				return false, fmt.Errorf("Key %s's value is not a string or stringified value, which we currently do not support (%#v)", k, v)
			}
		} else {
			// not a match if one of the search keys is not present at all
			return false, nil
		}
	}
	// All keys matched, we have a winner
	return true, nil
}

// trimMap goes thru a map[string]interface{}, and removes keys that
// have zero or nil values.
func trimMap(in map[string]interface{}) {
	for k, v := range in {
		switch {
		case v == nil:
			fallthrough
		case reflect.ValueOf(v).Interface() == reflect.Zero(reflect.TypeOf(v)).Interface():
			delete(in, k)
		}
	}
}

// updateCustomFields performs an update of custom fields on a resource, with
// the following stipulations:
//   - If we have custom fields, we need to do a diff on what is set versus
//     what isn't set, and ensure that we clear out the keys that aren't set.
//     Since our SDK does not currently support NOT NULL custom fields in
//     PHPIPAM, we can safely set these to nil.
//   - If we don't have a value for
//     custom_fields at all, set all keys to nil and update so that all custom
//     fields get blown away.
func updateCustomFields(d *schema.ResourceData, client interface{}) error {
	log.Printf("Start Update custom fields ...............")
	customFields := make(map[string]interface{})
	log.Printf("Defined custom fields ...............%s", customFields)

	if m, ok := d.GetOk("custom_fields"); ok {
		customFields = m.(map[string]interface{})
	}
	var old map[string]interface{}
	var err error
	switch c := client.(type) {
	case *addresses.Controller:
		old, err = c.GetAddressCustomFields(d.Get("address_id").(int))
	case *subnets.Controller:
		old, err = c.GetSubnetCustomFields(d.Get("subnet_id").(int))
	case *vlans.Controller:
		old, err = c.GetVLANCustomFields(d.Get("vlan_id").(int))
	default:
		panic(fmt.Errorf("Invalid client type passed %#v - this is a bug", client))
	}
	if err != nil {
		return fmt.Errorf("Error getting custom fields for updating: %s", err)
	}
nextKey:
	for k := range old {
		for l, v := range customFields {
			if k == l {
				customFields[l] = v
				continue nextKey
			}
		}
		customFields[k] = nil
	}

	switch c := client.(type) {
	case *addresses.Controller:
		_, err = c.UpdateAddressCustomFields(d.Get("address_id").(int), customFields)
	case *subnets.Controller:
		_, err = c.UpdateSubnetCustomFields(d.Get("subnet_id").(int), customFields)
	case *vlans.Controller:
		_, err = c.UpdateVLANCustomFields(d.Get("vlan_id").(int), d.Get("name").(string), customFields)
	default:
		panic(fmt.Errorf("Invalid client type passed %#v - this is a bug", client))
	}
	return err
}
