// Package search provides types and methods for working with the search
// controller.
package search

import (
	"encoding/json"
	"fmt"
	"net"
	"strconv"

	"github.com/pavel-z1/phpipam-sdk-go/controllers/subnets"
	"github.com/pavel-z1/phpipam-sdk-go/phpipam/client"
	"github.com/pavel-z1/phpipam-sdk-go/phpipam/session"
)

// subnetSearchResult is used internally to parse subnet search responses.
type subnetSearchResult struct {
	Subnets struct {
		Code int             `json:"code"`
		Data json.RawMessage `json:"data"`
	} `json:"subnets"`
}

// Controller is the base client for the Search controller.
type Controller struct {
	client.Client
}

// NewController returns a new instance of the client for the Search controller.
func NewController(sess *session.Session) *Controller {
	c := &Controller{
		Client: *client.NewClient(sess),
	}
	return c
}

// SearchSubnets searches for subnets in the phpIPAM database.
//
// This requires PHPIPAM 1.6 or higher.
func (c *Controller) SearchSubnets(searchString string) (out []subnets.Subnet, err error) {
	var result subnetSearchResult
	err = c.SendRequest("GET", fmt.Sprintf("/search/%s/?addresses=0&subnets=1&vlans=0&vrfs=0", searchString), &struct{}{}, &result)
	if err != nil {
		return
	}

	if result.Subnets.Code != 200 {
		var errorMsg string
		if err := json.Unmarshal(result.Subnets.Data, &errorMsg); err == nil {
			return nil, fmt.Errorf("search API returned code %d: %s", result.Subnets.Code, errorMsg)
		}
		return nil, fmt.Errorf("search API returned non-200 code: %d", result.Subnets.Code)
	}

	if err = json.Unmarshal(result.Subnets.Data, &out); err != nil {
		return nil, fmt.Errorf("failed to parse subnet data: %w", err)
	}

	// For some reason, subnet addresses are returned as integers: https://github.com/phpipam/phpipam/issues/4159
	for i := range out {
		ip_integer, err := strconv.Atoi(out[i].SubnetAddress)
		if err != nil {
			return out, err
		}
		ip := net.IPv4(byte(ip_integer>>24), byte(ip_integer>>16), byte(ip_integer>>8), byte(ip_integer))
		out[i].SubnetAddress = ip.String()
	}

	return
}
