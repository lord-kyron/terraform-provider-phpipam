// Package l2domains provides types and methods for working with the l2domains
// controller.
package l2domains

import (
	"fmt"

	"github.com/pavel-z1/phpipam-sdk-go/controllers/vlans"
	//"github.com/pavel-z1/phpipam-sdk-go/phpipam"
	"github.com/pavel-z1/phpipam-sdk-go/phpipam/client"
	"github.com/pavel-z1/phpipam-sdk-go/phpipam/session"
)

// L2Domain represents a PHPIPAM l2domain.
type L2Domain struct {
	// The L2 domain ID.
	ID int `json:"id,string,omitempty"`

	// The L2 domains name.
	Name string `json:"name,omitempty"`

	// The l2domain's description.
	Description string `json:"description,omitempty"`

	// The ID of the section's parent, if nested.
	Sections string `json:"sections,omitempty"`
}

// Controller is the base client for the L2Domains controller.
type Controller struct {
	client.Client
}

// NewController returns a new instance of the client for the L2Domains controller.
func NewController(sess *session.Session) *Controller {
	c := &Controller{
		Client: *client.NewClient(sess),
	}
	return c
}

// ListL2Domains lists all l2domains.
func (c *Controller) ListL2Domains() (out []L2Domain, err error) {
	err = c.SendRequest("GET", "/l2domains/", &struct{}{}, &out)
	return
}

// CreateL2Domain creates a l2domain by sending a POST request.
func (c *Controller) CreateL2Domain(in L2Domain) (message string, err error) {
	err = c.SendRequest("POST", "/l2domains/", &in, &message)
	return
}

// GetL2DomainByID GETs a l2domain via its ID.
func (c *Controller) GetL2DomainByID(id int) (out L2Domain, err error) {
	err = c.SendRequest("GET", fmt.Sprintf("/l2domains/%d/", id), &struct{}{}, &out)
	return
}

// GetL2DomainByName GETs a l2domain via its name.
func (c *Controller) GetL2DomainByName(name string) (out []L2Domain, err error) {
	err = c.SendRequest("GET", fmt.Sprintf("/l2domains/?filter_by=name&filter_value=%s", name), &struct{}{}, &out)
	return
}

// GetVlansInL2Domain GETs the vlans in a l2domains by l2domain ID.
func (c *Controller) GetVlansInl2Domain(id int) (out []vlans.VLAN, err error) {
	err = c.SendRequest("GET", fmt.Sprintf("/l2domains/%d/vlans/", id), &struct{}{}, &out)
	return
}

// UpdateL2Domain updates a l2domain by sending a PATCH request.
func (c *Controller) UpdateL2Domain(in L2Domain) (err error) {
	err = c.SendRequest("PATCH", "/l2domains/", &in, &struct{}{})
	return
}

// DeleteL2Domain deletes a l2domain by sending a DELETE request. All subnets and
func (c *Controller) DeleteL2Domain(id int) (err error) {
	err = c.SendRequest("DELETE", fmt.Sprintf("/l2domains/%d/", id), &struct{}{}, &struct{}{})
	return
}
