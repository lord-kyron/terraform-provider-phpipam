// Package sections provides types and methods for working with the sections
// controller.
package sections

import (
	"fmt"

	"github.com/pavel-z1/phpipam-sdk-go/controllers/subnets"
	"github.com/pavel-z1/phpipam-sdk-go/phpipam"
	"github.com/pavel-z1/phpipam-sdk-go/phpipam/client"
	"github.com/pavel-z1/phpipam-sdk-go/phpipam/session"
)

// Section represents a PHPIPAM section.
type Section struct {
	// The section ID.
	ID int `json:"id,string,omitempty"`

	// The section's name.
	Name string `json:"name,omitempty"`

	// The section's description.
	Description string `json:"description,omitempty"`

	// The ID of the section's parent, if nested.
	MasterSection int `json:"masterSection,string,omitempty"`

	// A JSON object, stringified, that represents the permissions for this
	// section.
	Permissions string `json:"permissions,omitempty"`

	// Whether or not to check consistency for subnets and IP addresses.
	StrictMode phpipam.BoolIntString `json:"strictMode,omitempty"`

	// How to order subnets in this section when viewing.
	SubnetOrdering string `json:"subnetOrdering,omitempty"`

	// The order position of this section when displaying sections.
	Order int `json:"order,string,omitempty"`

	// The date of the last edit to this resource.
	EditDate string `json:"editDate,omitempty"`

	// Whether or not to show VLANs in the subnet listing of this section.
	ShowVLAN phpipam.BoolIntString `json:"showVLAN,omitempty"`

	// Whether or not to show VRF information in the subnet listing of this
	// section.
	ShowVRF phpipam.BoolIntString `json:"showVRF,omitempty"`

	// Whether or not to show only supernets in the subnet listing of this
	// section.
	ShowSupernetOnly phpipam.BoolIntString `json:"showSupernetOnly,omitempty"`

	// The ID of the DNS resolver to be used for this section.
	DNS int `json:"DNS,string,omitempty"`
}

// Controller is the base client for the Sections controller.
type Controller struct {
	client.Client
}

// NewController returns a new instance of the client for the Sections controller.
func NewController(sess *session.Session) *Controller {
	c := &Controller{
		Client: *client.NewClient(sess),
	}
	return c
}

// ListSections lists all sections.
func (c *Controller) ListSections() (out []Section, err error) {
	err = c.SendRequest("GET", "/sections/", &struct{}{}, &out)
	return
}

// CreateSection creates a section by sending a POST request.
func (c *Controller) CreateSection(in Section) (message string, err error) {
	err = c.SendRequest("POST", "/sections/", &in, &message)
	return
}

// GetSectionByID GETs a section via its ID.
func (c *Controller) GetSectionByID(id int) (out Section, err error) {
	err = c.SendRequest("GET", fmt.Sprintf("/sections/%d/", id), &struct{}{}, &out)
	return
}

// GetSectionByName GETs a section via its name.
func (c *Controller) GetSectionByName(name string) (out Section, err error) {
	err = c.SendRequest("GET", fmt.Sprintf("/sections/%s/", name), &struct{}{}, &out)
	return
}

// GetSubnetsInSection GETs the subnets in a section by section ID.
func (c *Controller) GetSubnetsInSection(id int) (out []subnets.Subnet, err error) {
	err = c.SendRequest("GET", fmt.Sprintf("/sections/%d/subnets/", id), &struct{}{}, &out)
	return
}

// UpdateSection updates a section by sending a PATCH request.
func (c *Controller) UpdateSection(in Section) (err error) {
	err = c.SendRequest("PATCH", "/sections/", &in, &struct{}{})
	return
}

// DeleteSection deletes a section by sending a DELETE request. All subnets and
// addresses in the section will be deleted as well.
func (c *Controller) DeleteSection(id int) (err error) {
	err = c.SendRequest("DELETE", fmt.Sprintf("/sections/%d/", id), &struct{}{}, &struct{}{})
	return
}
