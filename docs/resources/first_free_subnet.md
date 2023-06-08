# phpipam_first_free_subnet

The `phpipam_first_free_subnet` resource allow you to create new subnet automatically
in defined network without execution Terraform data instruction. You can use it
to create several subnets automatically with nested subnet creation allowed. This resource
support the same arguments as [`phpipam_subnet`](./subnet.md). An example usage is below.

```hcl
// Look up the subnet
data "phpipam_subnet" "subnet" {
  subnet_address = "10.10.2.0"
  subnet_mask    = 22
}

resource "phpipam_first_free_subnet" "new_subnet" {
  parent_subnet_id   = data.phpipam_subnet.subnet.subnet_id
  subnet_mask = 24
  description = "Managed by Terraform"
}

resource "phpipam_first_free_subnet" "Child_subnet_of_new_subnet" {
  parent_subnet_id   = phpipam_first_free_subnet.new_subnet.subnet_id
  subnet_mask = 25
  description = "Managed by Terraform"
}

resource "phpipam_first_free_subnet" "Child_subnet_of_new_subnet" {
  parent_subnet_id   = phpipam_first_free_subnet.new_subnet.subnet_id
  subnet_mask = 25
  description = "Managed by Terraform"
}
```
