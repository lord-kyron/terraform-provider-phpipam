# phpipam_section

The `phpipam_section` resource manages a PHPIPAM section - a top-level category
that subnets and IP addresses are entered into. Use this resource if you want to
manage a section entirely from Terraform. If you just need to get information on
a section use the [`phpipam_section` data source](../data-sources/section.md)
instead.

**Example:**

```hcl
// Create a section
resource "phpipam_section" "section" {
  name        = "tf-test"
  description = "Terraform test section"
}
```

## Argument Reference

The resource takes the following parameters:

- `name` (Required) - The name of the section.
- `description` (Optional) - The section's description.
- `master_section_id` (Optional) - The ID of the parent section in the PHPIPAM
  database.
- `strict_mode` (Optional) - `true` if this subnet is set up to check that IP
  addresses are valid for the subnets they are in.
- `subnet_ordering` (Optional) - How subnets in this section are ordered.
- `display_order` (Optional) - The section's display order number.
- `show_vlan_in_subnet_listing` (Optional) - `true` if VLANs are being shown in
  the subnet listing for this section.
- `show_vrf_in_subnet_listing` (Optional) - `true` if VRFs are being shown in
  the subnet listing for this section.
- `show_supernet_only` (Optional) - `true` if supernets are only being shown in
  the subnet listing.
- `dns_resolver_id` (Optional) - The ID of the DNS resolver to use in the
  section.

## Attribute Reference

The following attributes are exported:

- `section_id` - The ID of the section in the PHPIPAM database.
- `edit_date` - The date this resource was last edited.
