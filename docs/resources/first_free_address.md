# phpipam_first_free_address

The `phpipam_first_free_address` resource allow to create automatically
new IP in defined network without execution Terraform data instruction. You can use it
to create several IP addresses automatically. This resource support the same arguments as
[`phpipam_address`](./address.md).  An example usage is below.

⚠️  **NOTE:** This is experimental new feature. You can use Terraform count
instruction. But be carefull, phpIPAM currently has a bug
[https://github.com/phpipam/phpipam/issues/2960](https://github.com/phpipam/phpipam/issues/2960)
Use resource with count option only with limited terraform threads count: `**terraform apply -parallelism=1**`
.

**Example create IPs in loop with `count`:**

```hcl
// Look up the subnet
data "phpipam_subnet" "subnet" {
  subnet_address = "10.10.2.0"
  subnet_mask    = 24
}

// Reserve the address. Note that we use ignore_changes here to ensure that we
// don't end up re-allocating this address on future Terraform runs.
resource "phpipam_first_free_address" "new_ip" {
  count = 3

  subnet_id   = data.phpipam_subnet.subnet.subnet_id
  hostname    = format("vps0%d.example.internal", count.index)
  description = "Managed by Terraform"
}

// IPs will be accessible by:
// phpipam_first_free_address.new_ip[0].ip_address
// phpipam_first_free_address.new_ip[0].ip_address
// phpipam_first_free_address.new_ip[0].ip_address
```

***Example: Create IPs in loop and use then in `for each`**

```hcl
// Look up the subnet
data "phpipam_subnet" "subnet" {
  count = 2

  // Will be defined subnet addresses: 192.168.0.0, 192.168.1.0
  subnet_address = format("192.168.%d.0", count.index)
  subnet_mask    = 24
}

resource "phpipam_first_free_address" "new_ip" {
  count = 2

  subnet_id   = element(data.phpipam_subnet.subnet.*, count.index).subnet_id
  hostname    = format("vps0%d.example.internal", count.index)
  description = "Managed by Terraform"
}

// Build oVirt VPS with two NIC in nic_configuration parameters
resource "ovirt_vm" "vm" {
    name                 = "vm-test01"
    cluster_id           = "178a56cc-2d3b-11ea-8913-00163e00bc29"
    memory               = 2028
    template_id          = "74ba668e-19c6-41a7-b5f9-b2fffadac3ff"
    cores                = 2

    initialization {
        host_name          = "vm-test01.example.com"
        dns_servers        = "8.8.8.8"

        // Dynamic declaration of several nic_configuration resources
        dynamic "nic_configuration" {
            for_each = phpipam_first_free_address.new_ip

            content {
                boot_proto = "static"
                label      = format("eth%d", nic_configuration.key)
                address    = phpipam_first_free_address.new_ip[nic_configuration.key].ip_address
                netmask    = element(data.phpipam_subnet.subnet.*, nic_configuration.key).subnet_mask
                on_boot    = true
            }
        }
    }
}
```

Result will be like this:

```text
$ terraform plan
Refreshing Terraform state in-memory prior to plan...
The refreshed state will be used to calculate this plan, but will not be
persisted to local or remote state storage.

data.phpipam_subnet.subnet[0]: Refreshing state...
data.phpipam_subnet.subnet[1]: Refreshing state...

------------------------------------------------------------------------

An execution plan has been generated and is shown below.
Resource actions are indicated with the following symbols:
  + create

Terraform will perform the following actions:

  # ovirt_vm.vm will be created
  + resource "ovirt_vm" "vm" {
      + clone             = false
      + cluster_id        = "178a56cc-2d3b-11ea-8913-00163e00bc29"
      + cores             = 2
      + high_availability = false
      + id                = (known after apply)
      + memory            = 2028
      + name              = "vm-test01"
      + sockets           = 1
      + status            = (known after apply)
      + template_id       = "74ba668e-19c6-41a7-b5f9-b2fffadac3ff"
      + threads           = 1

      + initialization {
          + dns_servers = "8.8.8.8"
          + host_name   = "vm-test01.example.com"

          + nic_configuration {
              + address    = (known after apply)
              + boot_proto = "static"
              + label      = "eth0"
              + netmask    = "24"
              + on_boot    = true
            }
          + nic_configuration {
              + address    = (known after apply)
              + boot_proto = "static"
              + label      = "eth1"
              + netmask    = "24"
              + on_boot    = true
            }
        }
    }

  # phpipam_first_free_address.new_ip[0] will be created
  + resource "phpipam_first_free_address" "new_ip" {
      + address_id        = (known after apply)
      + description       = "Managed by Terraform"
      + device_id         = (known after apply)
      + edit_date         = (known after apply)
      + exclude_ping      = (known after apply)
      + hostname          = "vps00.example.internal"
      + id                = (known after apply)
      + ip_address        = (known after apply)
      + is_gateway        = (known after apply)
      + last_seen         = (known after apply)
      + mac_address       = (known after apply)
      + note              = (known after apply)
      + owner             = (known after apply)
      + ptr_record_id     = (known after apply)
      + skip_ptr_record   = (known after apply)
      + state_tag_id      = (known after apply)
      + subnet_id         = 8
      + switch_port_label = (known after apply)
    }

  # phpipam_first_free_address.new_ip[1] will be created
  + resource "phpipam_first_free_address" "new_ip" {
      + address_id        = (known after apply)
      + description       = "Managed by Terraform"
      + device_id         = (known after apply)
      + edit_date         = (known after apply)
      + exclude_ping      = (known after apply)
      + hostname          = "vps01.example.internal"
      + id                = (known after apply)
      + ip_address        = (known after apply)
      + is_gateway        = (known after apply)
      + last_seen         = (known after apply)
      + mac_address       = (known after apply)
      + note              = (known after apply)
      + owner             = (known after apply)
      + ptr_record_id     = (known after apply)
      + skip_ptr_record   = (known after apply)
      + state_tag_id      = (known after apply)
      + subnet_id         = 28
      + switch_port_label = (known after apply)
    }

Plan: 3 to add, 0 to change, 0 to destroy.

------------------------------------------------------------------------
```
