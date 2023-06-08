# phpipam_section

The `phpipam_section` data source allows one to look up a specific section,
either by database ID or name. This data can then be used to manage other parts
of PHPIPAM, such as in the event that the section name is known but not its ID,
which is required for managing subnets.

**Example:**

```hcl
data "phpipam_section" "section" {
  name = "Customers"
}

resource "phpipam_subnet" "subnet" {
  section_id = data.phpipam_section.section.section_id
  subnet_address = "10.10.3.0"
  subnet_mask = 24
}
```

## Argument Reference

The data source takes the following parameters:

- `section_id` - The ID of the section to look up.
- `name` - The name of the section to look up.

One of `section_id` or `name` must be supplied. If both are supplied,
`section_id` is used.

## Attribute Reference

The following attributes are exported:

- `section_id` - The ID of the section in the PHPIPAM database.
- `name` - The name of the section.
- `description` - The section's description.
- `master_section_id` - The ID of the parent section in the PHPIPAM database.
- `permissions` - A JSON representation of permissions for this section.
- `strict_mode` - `true` if this subnet is set up to check that IP addresses
   are valid for the subnets they are in.
- `subnet_ordering` - How subnets in this section are ordered.
- `display_order` - The section's display order number.
- `edit_date` - The date this resource was last edited.
- `show_vlan_in_subnet_listing` - `true` if VLANs are being shown in the subnet
   listing for this section.
- `show_vrf_in_subnet_listing` - `true` if VRFs are being shown in the subnet
   listing for this section.
- `show_supernet_only` - `true` if supernets are only being shown in the subnet
   listing.
- `dns_resolver_id` - The ID of the DNS resolver to use in the section.
