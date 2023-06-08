# phpipam_subnets

The `phpipam_subnets` data source allows you to search for subnets, much in the
same way as you can in the single-form [`phpipam_subnet`](./subnet.md) data
source.  However, multiple subnets are returned from this data source as a
single list of subnet IDs as they are found in the PHPIPAM database. You can
then use the single-form [`phpipam_subnet`](./subnet.md) data source to extract
the subnet data for each matched network in the database.

**Example:**

⚠️  **NOTE:** The below example requires Terraform v0.12.0 or later!

```hcl
data "phpipam_subnets" "subnet_search" {
  subnet_id = 3

  custom_field_filter {
    custom_CustomTestSubnets = ".*terraform.*"
  }
}

data "phpipam_subnet" "subnets" {
  count     = length(data.phpipam_subnets.subnet_search.subnet_ids)
  subnet_id = element(data.phpipam_subnets.subnet_search.subnet_ids, count.index)
}

output "subnet_addresses" {
  value = [data.phpipam_subnet.subnets.*.ip_address]
}

output "subnet_cidrs" {
  value = [formatlist("%s/%d", data.phpipam_subnet.subnets.*.subnet_address, data.phpipam_subnet.subnets.*.subnet_mask)]
}
```

## Argument Reference

The data source takes the following parameters:

- `section_id` (Required) - The ID of the section of the subnet.

One of the following below parameters is required:

- `description` - The subnet's description.
- `description_match` - A regular expression to match against when searching
  for a subnet.
- `custom_field_filter` - A map of custom fields to search for. The filter
  values are regular expressions. All fields need to match for the match to
  succeed.

You can find documentation for the regular expression syntax used with the
`description_match` and `custom_field_filter` attributes
[here](https://github.com/google/re2/wiki/Syntax).

⚠️  **NOTE:** An empty or unspecified `custom_field_filter` value is the
equivalent to a regular expression that matches everything, and hence will
return **all** subnets that contain the referenced custom field key!
Custom fileds must contain mandatory prefix `custom_`.

## Attribute Reference

The following attributes are exported:

- `subnet_ids` - A list of subnet IDs that match the given criteria.
