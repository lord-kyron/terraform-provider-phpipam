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
  section_id = 3

  custom_field_filter = {
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

To search across all sections, use `search` instead of, or together with,
`section_id`:

```hcl
data "phpipam_subnets" "customer_subnets" {
  search = "Customer 1"
}
```

## Argument Reference

The data source takes the following parameters. At least one of `section_id`
or `search` is required.

- `section_id` - The ID of the section of the subnet. When used with `search`,
  limits the search results to this section.
- `search` - A free-text search across all sections, using the phpIPAM search
  controller. It matches the same fields as the search box in the phpIPAM web
  interface, such as subnet addresses, descriptions and custom fields. Must not
  contain `/`, `?` or `#`; to find a subnet by CIDR, search for its network
  address. Requires phpIPAM 1.6 or higher.

When only `section_id` is set, one of the following below parameters is
required. When `search` is set, they are optional and narrow down the search
results:

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

⚠️  **NOTE:** `search` currently fails if any IPv6 subnet matches the search
term, due to a bug in the phpIPAM SDK that will be fixed in a future release.
Until then, use `search` only with terms that match IPv4 subnets.

## Attribute Reference

The following attributes are exported:

- `subnet_ids` - A list of subnet IDs that match the given criteria.
