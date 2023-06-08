# phpipam_first_free_subnet

The `phpipam_first_free_subnet` data source allows you to get the next
available subnet address in a specific subnet in PHPIPAM. Using this resource allows
you to automatically allocate an subnet CIDR address that can be used as an CIDR in
resources such as `aws VPC`, or other public or private cloud that require CIDR range.

Note that not having any subnet available will cause the Terraform run to
fail. Conversely, marking a subnet as unavailable or used will not prevent this
data source from returning an next available subnet address, so be aware of this while
using this resource.

**Example:**

```hcl
// Look up the subnet
data "phpipam_subnet" "subnet" {
  subnet_address = "10.10.2.0"
  subnet_mask    = 24
}

// Get the first available address
data "phpipam_first_free_subnet" "next_subnet" {
  subnet_id = data.phpipam_subnet.subnet.subnet_id
  subnet_mask = 25
}
```

## Argument Reference

The data source takes the following parameters:

- `subnet_id` - The ID of the subnet to look up the subnet in.
- `subnet_mask` - Mask that will be used to look next available subnet

## Attribute Reference

The following attributes are exported:

- `ip_address` - The next available IP address.
