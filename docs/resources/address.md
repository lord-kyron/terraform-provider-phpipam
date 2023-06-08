# phpipam_address

The `phpipam_address` resource manages an IP address in PHPIPAM. You can use it
to create IP address reservations for IP addresses that have been created by
other Terraform resources, or supplied by the
[`phpipam_first_free_address`](../data-sources/first_free_address.md) data
source. An example usage is below. Custom fields must contain the mandatory
prefix `custom_`.

⚠️  **NOTE:** If you are using the
[`phpipam_first_free_address`](../data-sources/first_free_address.md) to get
the first free IP address in a specific subnet, make sure you set `subnet_id`
and `ip_address` as ignored attributes with the `ignore_changes` lifecycle
attribute. This will prevent Terraform from perpetually deleting and
re-allocating the address when it sees a different available IP address in the
[`phpipam_first_free_address` data source](../data-sources/first_free_address.md).

**Example:**

```hcl
// Look up the subnet
data "phpipam_subnet" "subnet" {
  subnet_address = "10.10.2.0"
  subnet_mask    = 24
}

// Get the first available address
data "phpipam_first_free_address" "next_address" {
  subnet_id = data.phpipam_subnet.subnet.subnet_id
}

// Reserve the address. Note that we use ignore_changes here to ensure that we
// don't end up re-allocating this address on future Terraform runs.
resource "phpipam_address" "newip" {
  subnet_id   = data.phpipam_subnet.subnet.subnet_id
  ip_address  = data.phpipam_first_free_address.next_address.ip_address
  hostname    = "tf-test-host.example.internal"
  description = "Managed by Terraform"

  custom_fields = {
    custom_CustomTestAddresses = "terraform-test"
  }

  lifecycle {
    ignore_changes = [
      subnet_id,
      ip_address,
    ]
  }
}
```

## Argument Reference

The resource takes the following parameters:

- `subnet_id` (Required) - The database ID of the subnet this IP address
   belongs to.
- `ip_address` (Required) - The IP address to reserve.
- `is_gateway` (Optional) - `true` if this IP address has been designated as a
   gateway.
- `description` (Optional) - The description provided to this IP address.
- `hostname` (Optional) - The hostname supplied to this IP address.
- `owner` (Optional) - The owner name provided to this IP address.
- `mac_address` (Optional) - The MAC address provided to this IP address.
- `state_tag_id` (Optional) - The tag ID in the database for the IP address'
   specific state. **NOTE:** This is currently represented as an integer but may
   change to the specific string representation at a later time.
- `skip_ptr_record` (Optional) - `true` if PTR records are not being created
   for this IP address.
- `ptr_record_id` (Optional) - The ID of the associated PTR record in the
   PHPIPAM database.
- `device_id` (Optional) - The ID of the associated device in the PHPIPAM
   database.
- `switch_port_label` (Optional) - A string port label that is associated with
   this address.
- `note` (Optional) - The note supplied to this IP address.
- `exclude_ping` (Optional) - `true` if this address is excluded from ping
   probes.
- `remove_dns_on_delete` (Optional) - Removes DNS records created by PHPIPAM
   when the address is deleted from Terraform. Defaults to `true`.
- `custom_fields` (Optional) -  A key/value map of custom fields for this address.

⚠️  **NOTE on custom fields:** PHPIPAM installations with custom fields must have
all fields set to optional when using this plugin. For more info see
[here](https://github.com/phpipam/phpipam/issues/1073). Further to this, either
ensure that your fields also do not have default values, or ensure the default
is set in your TF configuration. Diff loops may happen otherwise!
Custom fileds must contain mandatory prefix `custom_`.

## Attribute Reference

The following attributes are exported:

- `address_id` - The ID of the IP address in the PHPIPAM database.
- `last_seen` - The last time this IP address answered ping probes.
- `edit_date` - The last time this resource was modified.
