// Package vlans provides types and methods for working with the VLAN
// controller.
package vlans

import (
	"fmt"

	"github.com/paybyphone/phpipam-sdk-go/phpipam"
	"github.com/paybyphone/phpipam-sdk-go/phpipam/client"
	"github.com/paybyphone/phpipam-sdk-go/phpipam/session"
)

// VLAN represents a PHPIPAM VLAN.
type VLAN struct {
	// The VLAN ID. This is the entry ID in the PHPIPAM database, and not the
	// VLAN number, which is represented by the Number field.
	ID int `json:"id,string,omitempty"`

	// The Layer 2 domain identifier of the VLAN.
	DomainID int `json:"domainId,string,omitempty"`

	// The VLAN name/label.
	Name string `json:"name,omitempty"`

	// The VLAN number.
	Number int `json:"number,string,omitempty"`

	// A detailed description of the VLAN.
	Description string `json:"description,omitempty"`

	// The date of the last edit to this resource.
	EditDate string `json:"editDate,omitempty"`

	// A map[string]interface{} of custom fields to set on the resource. Note
	// that this functionality requires PHPIPAM 1.3 or higher with the "Nest
	// custom fields" flag set on the specific API integration. If this is not
	// enabled, this map will be nil on GETs and POSTs and PATCHes with this
	// field set will fail. Use the explicit custom field functions instead.
	CustomFields map[string]interface{} `json:"custom_fields,omitempty"`
}

// Controller is the base client for the VLAN controller.
type Controller struct {
	client.Client
}

// NewController returns a new instance of the client for the VLAN controller.
func NewController(sess *session.Session) *Controller {
	c := &Controller{
		Client: *client.NewClient(sess),
	}
	return c
}

// CreateVLAN creates a VLAN by sending a POST request.
func (c *Controller) CreateVLAN(in VLAN) (message string, err error) {
	err = c.SendRequest("POST", "/vlans/", &in, &message)
	return
}

// GetVLANByID GETs a VLAN via its ID in the PHPIPAM database.
func (c *Controller) GetVLANByID(id int) (out VLAN, err error) {
	err = c.SendRequest("GET", fmt.Sprintf("/vlans/%d/", id), &struct{}{}, &out)
	return
}

// GetVLANsByNumber GETs a VLAN via its VLAN number.
//
// This function is a search, however it's not entirely clear from the API spec
// on how to enter a search term that would return multiple VLANs. Nontheless,
// the output from this method is an array of VLANs, so this function returns a
// slice.
func (c *Controller) GetVLANsByNumber(id int) (out []VLAN, err error) {
	err = c.SendRequest("GET", fmt.Sprintf("/vlans/search/%d/", id), &struct{}{}, &out)
	return
}

// GetVLANCustomFieldsSchema GETs the custom fields for the vlans controller via
// client.GetCustomFieldsSchema.
func (c *Controller) GetVLANCustomFieldsSchema() (out map[string]phpipam.CustomField, err error) {
	out, err = c.Client.GetCustomFieldsSchema("vlans")
	return
}

// GetVLANCustomFields GETs the custom fields for a subnet via
// client.GetCustomFields.
func (c *Controller) GetVLANCustomFields(id int) (out map[string]interface{}, err error) {
	out, err = c.Client.GetCustomFields(id, "vlans")
	return
}

// UpdateVLAN updates a VLAN by sending a PATCH request.
func (c *Controller) UpdateVLAN(in VLAN) (message string, err error) {
	err = c.SendRequest("PATCH", "/vlans/", &in, &message)
	return
}

// UpdateVLANCustomFields PATCHes the vlan's custom fields.
//
// This function differs from the custom field functions available in the
// addresses and subnets controller - while those two controllers do not
// require any other data outside of the ID to update the custom fields,
// updating a VLAN requires a name as well.
func (c *Controller) UpdateVLANCustomFields(id int, name string, in map[string]interface{}) (message string, err error) {
	// Verify that we are only updating fields that are custom fields.
	var schema map[string]phpipam.CustomField
	schema, err = c.GetVLANCustomFieldsSchema()
	if err != nil {
		return
	}
	for k := range in {
		for l := range schema {
			if k == l {
				goto customFieldFound
			}
		}
		// not found
		return "", fmt.Errorf("Custom field %s not found in schema for controller vlans", k)
		// found
	customFieldFound:
	}

	params := make(map[string]interface{})
	for k, v := range in {
		params[k] = v
	}

	params["id"] = id
	params["name"] = name
	err = c.SendRequest("PATCH", "/vlans/", &params, &message)
	return
}

// DeleteVLAN deletes a VLAN by its ID.
func (c *Controller) DeleteVLAN(id int) (message string, err error) {
	err = c.SendRequest("DELETE", fmt.Sprintf("/vlans/%d/", id), &struct{}{}, &message)
	return
}
