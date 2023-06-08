# phpipam_first_free_address

The `phpipam_first_free_address` data source allows you to get the next
available IP address in a specific subnet in PHPIPAM. Using this resource allows
you to automatically allocate an IP address that can be used as an IP address in
resources such as `vsphere_virtual_machine`, or other virtual machine-like
resources that require static IP addresses.

Note that not having any addresses available will cause the Terraform run to
fail. Conversely, marking a subnet as unavailable or used will not prevent this
data source from returning an IP address, so be aware of this while using this
resource.

## Argument Reference

The data source takes the following parameter:

- `subnet_id` (Required) - The ID of the subnet that the address resides in.

**Example:**

```hcl
// Look up the subnet
data "phpipam_subnet" "subnet" {
  subnet_address = "10.10.2.0"
  subnet_mask    = 24
}

// Get the first available address
data "phpipam_first_free_address" "next_address" {
  subnet_id = data.phpipam_subnet.subnet.subnet_id
}

// Reserve the address. Note that we use ignore_changes here to ensure that we
// don't end up re-allocating this address on future Terraform runs.
resource "phpipam_address" "newip" {
  subnet_id   = data.phpipam_subnet.subnet.subnet_id
  ip_address  = data.phpipam_first_free_address.next_address.ip_address
  hostname    = "tf-test-host.example.internal"
  description = "Managed by Terraform"

  lifecycle {
    ignore_changes = [
      subnet_id,
      ip_address,
    ]
  }
}

// Supply the IP address to an instance. Note that we are also ignoring
// network_interface here to ensure the IP address does not get re-calculated.
resource "vsphere_virtual_machine" "web" {
  name   = "terraform-web"
  vcpu   = 2
  memory = 4096

  network_interface {
    label        = "VM Network"
    ipv4_address = data.phpipam_first_free_address.next_address.ip_address
  }

  disk {
    template = "centos-7"
  }

  ignore_changes = [
    network_interface,
  ]
}
```
