# phpipam_vlan

The `phpipam_vlan` data source allows one to look up a VLAN in the PHPIPAM
database. This can then be used to assign a VLAN to a subnet in the
[`phpipam_subnet` resource](../resources/subnet.md). It can also be used
to gather other information on the VLAN.

**Example:**

```hcl
data "phpipam_section" "section" {
  name = "Customers"
}

data "phpipam_vlan" "vlan" {
  number = 1000
}

resource "phpipam_subnet" "subnet" {
  section_id     = data.phpipam_section.section.section_id
  subnet_address = "10.10.3.0"
  subnet_mask    = 24
  vlan_id        = data.phpipam_vlan.vlan.vlan_id
}
```

## Argument Reference

The data source takes the following parameters:

- `vlan_id` - The ID of the VLAN to look up. **NOTE:** this is the database ID,
   not the VLAN number - if you need this, use the `number` parameter.
- `number` - The number of the VLAN to look up.

One of `vlan_id` or `number` must be supplied. If both are supplied,
`vlan_id` is used.

## Attribute Reference

The following attributes are exported:

- `vlan_id` - The ID of the VLAN to look up. **NOTE:** this is the database ID,
  not the VLAN number - if you need this, use the `number` parameter.
- `number` - The number of the VLAN (the actual VLAN ID on your switch).
- `l2_domain_id` - The layer 2 domain ID in the PHPIPAM database.
- `name` - The name/label of the VLAN.
- `description` - The description supplied to the VLAN.
- `edit_date` - The date this resource was last updated.
- `custom_fields` - A key/value map of custom fields for this VLAN.
