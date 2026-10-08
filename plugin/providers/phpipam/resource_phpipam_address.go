package phpipam

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/pavel-z1/phpipam-sdk-go/phpipam"
)

// resourcePHPIPAMAddress returns the resource structure for the phpipam_address
// resource.
//
// Note that we use the data source read function here to pull down data, as
// read workflow is identical for both the resource and the data source.
func resourcePHPIPAMAddress() *schema.Resource {
	return &schema.Resource{
		Create: resourcePHPIPAMAddressCreate,
		Read:   resourcePHPIPAMAddressRead,
		Update: resourcePHPIPAMAddressUpdate,
		Delete: resourcePHPIPAMAddressDelete,
		Schema: resourceAddressSchema(),
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

// resourcePHPIPAMAddressRead restores address_id from Terraform's resource ID
// for states created by provider versions that did not persist address_id.
// Reading by ID avoids the IP-and-subnet endpoint, which no longer works with
func resourcePHPIPAMAddressRead(d *schema.ResourceData, meta interface{}) error {
	if d.Get("address_id").(int) == 0 && d.Id() != "" {
		addressID, err := strconv.Atoi(d.Id())
		if err != nil {
			return fmt.Errorf("invalid address resource ID %q: %w", d.Id(), err)
		}
		if err := d.Set("address_id", addressID); err != nil {
			return err
		}
	}

	return dataSourcePHPIPAMAddressRead(d, meta)
}

func resourcePHPIPAMAddressCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ProviderPHPIPAMClient)
	c := client.addressesController

	// Get the next free address if no IP is specified.
	if d.Get("ip_address").(string) == "" {
		// By default Terraform runs operations in parallel. Protect the
		// GetFirstFreeAddress and CreateAddress operations with a lock so they are
		// not run concurrently.
		client.addressAllocationLock.Lock()
		defer client.addressAllocationLock.Unlock()

		subnet_c := client.subnetsController
		out, err := subnet_c.GetFirstFreeAddress(d.Get("subnet_id").(int))
		if err != nil {
			return err
		}
		if out == "" {
			return errors.New("Subnet has no free IP addresses")
		}

		d.Set("ip_address", out)
	}

	in := expandAddress(d)

	// Assert the ID field here is empty. If this is not empty the request will fail.
	in.ID = 0

	if _, err := c.CreateAddress(in); err != nil {
		return err
	}

	// phpIPAM's create endpoint does not return the new address ID. Resolve it
	// before the first read so that the resource is read by ID, rather than via
	// the IP-and-subnet endpoint, whose behavior changed in phpIPAM 1.8.3.
	addrs, err := c.GetAddressesByIP(in.IPAddress)
	if err != nil {
		return fmt.Errorf("could not read IP address after creating: %w", err)
	}

	var addressID int
	for _, addr := range addrs {
		if addr.SubnetID == in.SubnetID {
			addressID = addr.ID
			break
		}
	}
	if addressID == 0 {
		return errors.New("created IP address was not found in the requested subnet")
	}

	d.Set("address_id", addressID)
	d.SetId(strconv.Itoa(addressID))

	if customFields, ok := d.GetOk("custom_fields"); ok {
		if _, err := c.UpdateAddressCustomFields(addressID, customFields.(map[string]interface{})); err != nil {
			return err
		}
	}

	return dataSourcePHPIPAMAddressRead(d, meta)
}

func resourcePHPIPAMAddressUpdate(d *schema.ResourceData, meta interface{}) error {
	c := meta.(*ProviderPHPIPAMClient).addressesController
	in := expandAddress(d)

	// IPAddress and SubnetID need to be removed for update requests.
	in.IPAddress = ""
	in.SubnetID = 0
	// The server-managed timestamps need to be removed as well. They are
	// populated in state by the read, and sending them back on an update can
	// fail with "Invalid request key". Both are computed-only in the schema, so
	// they can never have been supplied by the practitioner.
	in.LastSeen = ""
	in.EditDate = ""
	if _, err := c.UpdateAddress(in); err != nil {
		return err
	}

	if err := updateCustomFields(d, c); err != nil {
		return err
	}

	return dataSourcePHPIPAMAddressRead(d, meta)
}

func resourcePHPIPAMAddressDelete(d *schema.ResourceData, meta interface{}) error {
	c := meta.(*ProviderPHPIPAMClient).addressesController
	in := expandAddress(d)

	if _, err := c.DeleteAddress(in.ID, phpipam.BoolIntString(d.Get("remove_dns_on_delete").(bool))); err != nil {
		return err
	}
	d.SetId("")
	return nil
}