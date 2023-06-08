# phpipam_subnet

The `phpipam_subnet` resource can be used to create and manage a subnet in
PHPIPAM. Use it to manage details on subnets you create in Terraform for other
things, such as storing the IDs of the subnets you create for AWS, or for a full
top-down management of subnets and IP addresses in Terraform. If you just need
to get information on a subnet, use the
[`phpipam_subnet` data source](../data-sources/subnet.md) instead.

**Example:**

```hcl
data "phpipam_section" "section" {
  name = "Customers"
}

resource "phpipam_subnet" "subnet" {
  section_id     = data.phpipam_section.section.section_id
  subnet_address = "10.10.3.0"
  subnet_mask    = 24

  custom_fields = {
    custom_CustomTestSubnets = "terraform-test"
  }
}
```

## Argument Reference

The resource takes the following parameters:

- `subnet_address` (Required) - The network address of the subnet.
- `subnet_mask` (Required) - The subnet mask, in bits.
- `description` (Optional) - The description set for the subnet.
- `section_id` (Optional) - The ID of the section for this address in the
   PHPIPAM database.
- `linked_subnet_id` (Optional) - The ID of the linked subnet in the PHPIPAM
   database.
- `vlan_id` (Optional) - The ID of the VLAN for this subnet in the PHPIPAM
   database.
- `vrf_id` (Optional) - The ID of the VRF for this subnet in the PHPIPAM
   database.
- `master_subnet_id` (Optional) - The ID of the parent subnet for this subnet
   in the PHPIPAM database.
- `nameserver_id` (Optional) - The ID of the nameserver used to assign PTR
   records for this subnet.
- `show_name` (Optional) - `true` if the subnet name is are shown in the
   section, instead of the network address.
- `create_ptr_records` (Optional) - `true` if PTR records are created for
   addresses in this subnet.
- `display_hostnames` (Optional) - `true` if hostnames are displayed instead of
   IP addresses in the address listing for this subnet.
- `allow_ip_requests` (Optional) - `true` if the subnet allows IP requests in
   PHPIPAM.
- `scan_agent_id` (Optional) - The ID of the ping scan agent that is used for
   this subnet.
- `include_in_ping` (Optional) - `true` if this subnet is included in ping
   probes.
- `host_discovery_enabled` (Optional) - `true` if this subnet is included in
   new host scans.
- `resolve_dns` (Optional) - `true` if enabled resolving of DNS names.
- `is_folder` (Optional) - `true` if this subnet is a folder and not an actual
   subnet.
- `is_full` (Optional) - `true` if the subnet has been marked as full.
- `state_tag_id` (Optional) - The ID of the state tag for this subnet. This may
   become an actual string representation of this at a later time (example:
   `Used`).
- `utilization_threshold` (Optional) - The subnet's utilization threshold.
- `location_id` (Optional) - The ID of the location for this subnet.
- `custom_fields` (Optional) -  A key/value map of custom fields for this
   subnet.

⚠️  **NOTE on custom fields:** PHPIPAM installations with custom fields must have
all fields set to optional when using this plugin. For more info see
[here](https://github.com/phpipam/phpipam/issues/1073). Further to this, either
ensure that your fields also do not have default values, or ensure the default
is set in your TF configuration. Diff loops may happen otherwise!
Custom fileds must contain mandatory prefix `custom_`.

## Attribute Reference

The following attributes are exported:

- `subnet_id` - The ID of the subnet in the PHPIPAM database.
- `permissions` - A JSON representation of the permissions associated with this
   subnet.
- `edit_date` - The date this resource was last updated.
