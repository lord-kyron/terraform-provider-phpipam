// Package client contains generic client structs and methods that are
// designed to be used by specific PHPIPAM services and resources.
package client

import (
	"fmt"
	"os"

	"github.com/apex/log"
	"github.com/apex/log/handlers/logfmt"
	"github.com/pavel-z1/phpipam-sdk-go/phpipam"
	"github.com/pavel-z1/phpipam-sdk-go/phpipam/request"
	"github.com/pavel-z1/phpipam-sdk-go/phpipam/session"
)

// Client encompasses a generic client object that is further extended by
// services. Any common configuration and functionality goes here.
type Client struct {
	// The session for this client.
	Session *session.Session
}

// NewClient creates a new client.
func NewClient(s *session.Session) *Client {
	log.SetLevel(log.InfoLevel)
	log.SetHandler(logfmt.New(os.Stderr))

	env_loglevel := os.Getenv("PHPIPAMSDK_LOGLEVEL")
	if env_loglevel != "" {
		loglevel, err := log.ParseLevel(env_loglevel)
		if err == nil {
			log.SetLevel(loglevel)
		} else {
			log.Warnf("Invalid log level, defaulting to info: %s", err)
		}
	}

	c := &Client{
		Session: s,
	}
	return c
}

// change logger level, default is info
func SetLevel(level log.Level) {
	log.SetLevel(level)
}

// loginSession logs in a session via the user controller. This is the only
// valid operation if the session does not have a token yet.
func loginSession(s *session.Session) error {
	if s.Config.Username == "" {
		s.Token.String = s.Config.Password
	} else {
		var out session.Token
		r := request.NewRequest(s)
		r.Method = "POST"
		r.URI = "/user/"
		r.Input = &struct{}{}
		r.Output = &out
		if err := r.Send(); err != nil {
			return err
		}
		s.Token = out
	}
	return nil
}

// SendRequest sends a request to a request.Request object.  It's expected that
// references to specific data types are passed - no checking is done to make
// sure that references are passed.
//
// This function also wraps session management into the workflow, logging in
// and refreshing session tokens as needed.
func (c *Client) SendRequest(method, uri string, in, out interface{}) error {
	// Check to make sure our session is ok first.
	if c.Session.Token.String == "" {
		if err := loginSession(c.Session); err != nil {
			return fmt.Errorf("Error logging into PHPIPAM: %s", err)
		}
	}

	r := request.NewRequest(c.Session)
	r.Method = method
	r.URI = uri
	r.Input = in
	r.Output = out
	err := r.Send()
	switch {
	case err == nil:
		return nil
	case err.Error() == "Error from API (403): Token expired":
		if err := loginSession(c.Session); err != nil {
			return fmt.Errorf("Error refreshing expired PHPIPAM session token: %s", err)
		}
		return r.Send()
	}
	return err
}

// GetCustomFieldsSchema GETs the custom fields for the supplied controller
// name and returns them as a map[string]phpipam.CustomField.
//
// This function is called out to in a controller to implement this
// functionality in a specific pacakge.
func (c *Client) GetCustomFieldsSchema(controller string) (out map[string]phpipam.CustomField, err error) {
	err = c.SendRequest("GET", fmt.Sprintf("/%s/custom_fields/", controller), &struct{}{}, &out)
	return
}

// GetCustomFields GETs the custom fields for a resource, and returns them
// as a map[string]interface{}. A call out to GetCustomFields is performed
// first, and then a GET is performed on the subnet resource with only the
// custom fields returned.
//
// Note that due to how PHPIPAM stringifies most output, this will, in most
// cases, mean that attribute values will be strings and will need to be
// convereted externally. This function does not explicitly lock to
// map[string]string to allow for possible cases where this is not the case,
// and to also allow for future de-stringification of the JSON.
//
// This function is called out to in a controller to implement this
// functionality in a specific pacakge.
func (c *Client) GetCustomFields(id int, controller string) (out map[string]interface{}, err error) {
	var schema map[string]phpipam.CustomField
	schema, err = c.GetCustomFieldsSchema(controller)
	switch {
	case err != nil:
		log.Warnf("Error getting custom Fields: %s", err)
		return
	}

	out, err = c.getCustomFieldsRequest(id, controller, schema)
	return
}

// getCustomFieldsRequest performs the actual work for GetCustomFields. This is
// separated off to make testing easier.
func (c *Client) getCustomFieldsRequest(id int, controller string, schema map[string]phpipam.CustomField) (out map[string]interface{}, err error) {
	err = c.SendRequest("GET", fmt.Sprintf("/%s/%d/", controller, id), &struct{}{}, &out)
	if err != nil {
		return
	}
	for k := range out {
		for l := range schema {
			if k == l {
				goto customFieldFound
			}
		}
		// not found
		delete(out, k)
		// found
	customFieldFound:
	}
	return
}

// UpdateCustomFields uses PATCH on a resource controller to update a specific
// resoruce ID with the custom fields provided in the key/value map defined by
// in.
//
// Internal validation is preformed first to ensure that this field is not
// setting a custom field that is *not* defined in the schema. This is to
// prevent abuse - if this was not in place, this function could technically be
// used to update *any* field, as PHPIPAM does not maintain a separate subtype
// for custom fields.
//
// This function is called out to in a controller to implement this
// functionality in a specific pacakge.
func (c *Client) UpdateCustomFields(id int, in map[string]interface{}, controller string) (message string, err error) {
	var schema map[string]phpipam.CustomField
	schema, err = c.GetCustomFieldsSchema(controller)
	switch {
	// Ignore this error if the caller is not setting any fields.
	case len(in) == 0 && err.Error() == "Error from API (200): No custom fields defined":
		err = nil
		return
	case err != nil:
		return
	}
	message, err = c.updateCustomFieldsRequest(id, in, controller, schema)
	return
}

// updateCustomFieldsRequest performs the actual validation and request work
// for UpdateCustomFields. This is separated off to make testing easier.
func (c *Client) updateCustomFieldsRequest(id int, in map[string]interface{}, controller string, schema map[string]phpipam.CustomField) (message string, err error) {
	for k := range in {
		for l := range schema {
			if k == l {
				goto customFieldFound
			}
		}
		// not found
		return "", fmt.Errorf("Custom field %s not found in schema for controller %s", k, controller)
		// found
	customFieldFound:
	}

	params := make(map[string]interface{})
	for k, v := range in {
		params[k] = v
	}

	params["id"] = id
	err = c.SendRequest("PATCH", fmt.Sprintf("/%s/", controller), &params, &message)
	return
}
