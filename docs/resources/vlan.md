# phpipam_vlan

The `phpipam_vlan` resource can be used to manage a VLAN on PHPIPAM. Use it to
set up a VLAN through Terraform, or update details such as its name or
description. If you are just looking for information on a VLAN, use the
[`phpipam_vlan` data source](../data-sources/vlan.md) instead.

**Example:**

```hcl
resource "phpipam_vlan" "vlan" {
  name        = "tf-test"
  number      = 1000
  description = "Managed by Terraform"

  custom_fields = {
    custom_CustomTestVLANs = "terraform-test"
  }
}
```

## Argument Reference

The resource takes the following parameters:

- `name` (Required) - The name/label of the VLAN.
- `number` (Required) - The number of the VLAN (the actual VLAN ID on your switch).
- `l2_domain_id` (Optional) - The layer 2 domain ID in the PHPIPAM database.
- `description` (Optional) - The description supplied to the VLAN.
- `edit_date` (Optional) - The date this resource was last updated.
- `custom_fields` (Optional) -  A key/value map of custom fields for this
   VLAN.

⚠️ **NOTE on custom fields:** PHPIPAM installations with custom fields must have
all fields set to optional when using this plugin. For more info see
[here](https://github.com/phpipam/phpipam/issues/1073). Further to this, either
ensure that your fields also do not have default values, or ensure the default
is set in your TF configuration. Diff loops may happen otherwise!

## Attribute Reference

The following attributes are exported:

- `vlan_id` - The ID of the VLAN to look up. **NOTE:** this is the database ID,
   not the VLAN number - if you need this, use the `number` parameter.
- `edit_date` - The date this resource was last updated.
