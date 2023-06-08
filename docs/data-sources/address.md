# phpipam_address

The `phpipam_address` data source allows one to get information about a specific
IP address within PHPIPAM. Use this address to get general information about a
specific IP address such as its host name, description and more.

Lookups for IP addresses can only happen at this time via its entry in the
database, or the IP address itself. Future versions of this resource, when such
features become generally available in the PHPIPAM API, will allow lookup based
on host name, allowing for better ability for this resource to discover IP
addresses that have been pre-assigned for a specific resource.

**Example:**

```hcl
data "phpipam_address" "address" {
  ip_address = "10.10.1.1"
}

output "address_description" {
  value = data.phpipam_address.address.description
}
```

**Example With `subnet_id` when multiple subnets share the same ip ranges :**

```hcl
data "phpipam_address" "address" {
  ip_address = "10.10.1.1"
  subnet_id         = 3
}

output "address_description" {
  value = data.phpipam_address.address.description
}
```

**Example With `description`:**

```hcl
data "phpipam_address" "address" {
  subnet_id         = 3
  description = "Customer 1"
}

output "address_description" {
  value = data.phpipam_address.address.description
}
```

**Example With `custom_field_filter`:**

```hcl
data "phpipam_address" "address" {
  subnet_id = 3

  custom_field_filter {
    custom_CustomTestAddresses = ".*terraform.*"
  }
}

output "address_description" {
  value = data.phpipam_address.address.description
}
```

## Argument Reference

The data source takes the following parameters:

- `address_id` - The ID of the IP address in the PHPIPAM database.
- `ip_address` - The actual IP address in PHPIPAM.
- `subnet_id` - The ID of the subnet that the address resides in. This is
  required to search on the `description` or `hostname` field. Optional if
  multiple subnets have the same ip ranges ( multiple subnets behind NAT )
- `description` - The description of the IP address. `subnet_id` is required
  when using this field.
- `hostname` - The host name of the IP address. `subnet_id` is required when
  using this field.
- `custom_field_filter` - A map of custom fields to search for. The filter
  values are regular expressions that follow the RE2 syntax for which you can
  find documentation [here](https://github.com/google/re2/wiki/Syntax). All
  fields need to match for the match to succeed.

⚠️ **NOTE:** `description`, `hostname`, and `custom_field_filter` fields return
the first match found without any warnings. If you are looking to return
multiple addresses, combine this data source with the
[`phpipam_addresses`](./addresses.md) data source.

⚠️ **NOTE:** An empty or unspecified `custom_field_filter` value is the
equivalent to a regular expression that matches everything, and hence will
return the first address it sees in the subnet. Custom fileds must contain mandatory
prefix `custom_`.

Arguments are processed in the following order of precedence:

- `address_id`
- `ip_address`
- `subnet_id`, and either one of `description`, `hostname`, or
   `custom_field_filter`

## Attribute Reference

The following attributes are exported:

- `address_id` - The ID of the IP address in the PHPIPAM database.
- `ip_address` - the IP address.
- `subnet_id` - The database ID of the subnet this IP address belongs to.
- `is_gateway` - `true` if this IP address has been designated as a gateway.
- `description` - The description provided to this IP address.
- `hostname` - The hostname supplied to this IP address.
- `owner` - The owner name provided to this IP address.
- `mac_address` - The MAC address provided to this IP address.
- `state_tag_id` - The tag ID in the database for the IP address' specific
   state. **NOTE:** This is currently represented as an integer but may change
   to the specific string representation at a later time.
- `skip_ptr_record` - `true` if PTR records are not being created for this IP
   address.
- `ptr_record_id` - The ID of the associated PTR record in the PHPIPAM
   database.
- `device_id` - The ID of the associated device in the PHPIPAM database.
- `switch_port_label` - A string port label that is associated with this
   address.
- `note` - The note supplied to this IP address.
- `last_seen` - The last time this IP address answered ping probes.
- `exclude_ping` - `true` if this address is excluded from ping probes.
- `edit_date` - The last time this resource was modified.
- `custom_fields` - A key/value map of custom fields for this address.
