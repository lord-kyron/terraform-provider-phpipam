# phpipam_addresses

The `phpipam_addresses` data source allows you to search for IP addresses, much
in the same way as you can in the single-form [`phpipam_address`](./address.md)
data source. However, multiple addresses are returned from this data source as
a single list of address IDs as they are found in the PHPIPAM database. You can
then use the single-form [`phpipam_address`](./address.md) data source to
extract the IP data for each matched address in the database.

**Example:**

⚠️  **NOTE:** The below example requires Terraform v0.12.0 or later!

```hcl
data "phpipam_addresses" "address_search" {
  subnet_id = 3

  custom_field_filter {
    custom_CustomTestAddresses = ".*terraform.*"
  }
}

data "phpipam_address" "addresses" {
  count      = length(data.phpipam_addresses.address_search.address_ids)
  address_id = element(data.phpipam_addresses.address_search.address_ids, count.index)
}

output "ip_addresses" {
  value = [data.phpipam_address.addresses.*.ip_address]
}
```

## Argument Reference

The data source takes the following parameters:

- `subnet_id` (Required) - The ID of the subnet that the address resides in. This is
   required to search on the `description` or `hostname` fields.

One of the following fields is required alongside `subnet_id`:

- `description` - The description of the IP address. `subnet_id` is required
   when using this field.
- `hostname` - The host name of the IP address. `subnet_id` is required when
   using this field.
- `custom_field_filter` - A map of custom fields to search for. The filter
   values are regular expressions that follow the RE2 syntax for which you can
   find documentation [here](https://github.com/google/re2/wiki/Syntax). All
   fields need to match for the match to succeed.

⚠️  **NOTE:** An empty or unspecified `custom_field_filter` value is the
equivalent to a regular expression that matches everything, and hence will
return **all** addresses that contain the referenced custom field key!
Custom fileds must contain mandatory prefix `custom_`.

## Attribute Reference

The following attributes are exported:

- `address_ids` - A list of discovered IP address IDs.
