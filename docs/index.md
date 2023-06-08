# Terraform PHPIPAM provider

![GitHub release (latest by date)](https://img.shields.io/github/v/release/lord-kyron/terraform-provider-phpipam?color=gr&label=version&style=flat-square&logo=terraform) ![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/lord-kyron/terraform-provider-phpipam?style=flat-square&logo=go) ![GitHub](https://img.shields.io/github/license/lord-kyron/terraform-provider-phpipam?color=orange&logo=apache&style=flat-square) ![GitHub last commit](https://img.shields.io/github/last-commit/lord-kyron/terraform-provider-phpipam?style=flat-square&logo=github) ![GitHub Workflow Status](https://img.shields.io/github/actions/workflow/status/lord-kyron/terraform-provider-phpipam/go.yml?style=flat-square&logo=github) ![GitHub Release Date](https://img.shields.io/github/release-date/lord-kyron/terraform-provider-phpipam?style=flat-square&logo=github) ![GitHub release (latest SemVer)](https://img.shields.io/github/v/release/lord-kyron/terraform-provider-phpipam?color=blueviolet&style=flat-square&logo=github) ![GitHub release (latest by date)](https://img.shields.io/github/downloads/lord-kyron/terraform-provider-phpipam/latest/total?style=flat-square&color=informational&logo=github)

This repository holds a external plugin for a [Terraform][1] provider to manage
resources within [PHPIPAM][2], an open source IP address management system.

[1]: https://www.terraform.io/
[2]: https://phpipam.net/

## About PHPIPAM

[PHPIPAM][2] is an open source IP address management system written in PHP. It
has an evolving [API][3] that allows for the management and lookup of data that
has been entered into the system. Through our Go integration
[phpipam-sdk-go][4], we have been able to take this API and integrate it into
Terraform, allowing for the management and lookup of sections, VLANs, subnets,
and IP addresses, entirely within Terraform.

[3]: https://phpipam.net/api/api_documentation/
[4]: https://github.com/pavel-z1/phpipam-sdk-go

## Usage

After installation, to use the plugin, simply use any of its resources or data
sources (such as [`phpipam_subnet`](./resources/subnet.md) or
[`phpipam_address`](./data-sources/address.md)) in a Terraform configuration.

Credentials can be supplied via configuration variables to the `phpipam`
provider instance, or via environment variables. These are documented in the
next section.

You can see the following example below for a simple usage example that reserves
the first available IP address in a subnet. This address could then be passed
along to the configuration for a VM, say, for example, a
[`vsphere_virtual_machine`][7] resource.

[7]: https://www.terraform.io/docs/providers/vsphere/r/virtual_machine.html

```hcl
provider "phpipam" {
  app_id   = "test"
  endpoint = "https://phpipam.example.com/api"
  password = "PHPIPAM_PASSWORD"
  username = "Admin"
  insecure = false
}

data "phpipam_subnet" "subnet" {
  subnet_address = "10.10.2.0"
  subnet_mask    = 24
}

data "phpipam_first_free_address" "next_address" {
  subnet_id = data.phpipam_subnet.subnet.subnet_id
}

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
```

### Data Sources

- [`phpipam_address`](./data-sources/address.md)
- [`phpipam_addresses`](./data-sources/addresses.md)
- [`phpipam_first_free_address`](./data-sources/first_free_address.md)
- [`phpipam_first_free_subnet`](./data-sources/first_free_subnet.md)
- [`phpipam_section`](./data-sources/section.md)
- [`phpipam_subnet`](./data-sources/subnet.md)
- [`phpipam_subnets`](./data-sources/subnets.md)
- [`phpipam_vlan`](./data-sources/vlan.md)

### Resources

- [`phpipam_address`](./resources/address.md)
- [`phpipam_first_free_address`](./resources/first_free_address.md)
- [`phpipam_first_free_subnet`](./resources/first_free_subnet.md)
- [`phpipam_section`](./resources/section.md)
- [`phpipam_subnet`](./resources/subnet.md)
- [`phpipam_vlan`](./resources/vlan.md)

### Plugin Options

The options for the plugin are as follows:

- `app_id` - The API application ID, configured in the PHPIPAM API panel. This
   application ID should have read/write access if you are planning to use the
   resources, but read-only access should be sufficient if you are only using
   the data sources. Can also be supplied by the `PHPIPAM_APP_ID` environment
   variable.
- `endpoint` - The full URL to the PHPIPAM API endpoint, such as
  `https://phpipam.example.com/api`. Can also be supplied by the
  `PHPIPAM_ENDPOINT_ADDR` environment variable.
- `password` - The password to access the PHPIPAM API with. Can also be
  supplied via `PHPIPAM_PASSWORD` to prevent plain text password storage in
  config.
- `username` - The user name to access the PHPIPAM API with. Can also be
  supplied via the `PHPIPAM_USER_NAME` variable.
- `insecure` - Set to true to not validate the HTTPS certificate chain.
   Optional parameter, can be used only with HTTPS connections

### Resource importing

Importing all resource types are supported.

**Example:**

```hcl
resource "phpipam_subnet" "imported" {
  #parent_subnet_id = data.phpipam_subnet.gcp_cidr_pool.subnet_id
  subnet_address = "172.20.0.0"
  subnet_mask = 24
  section_id = 1
}
```

```ShellSession
$ terraform import phpipam_subnet.imported 20

$ terraform state show phpipam_subnet.imported
# phpipam_subnet.imported:
resource "phpipam_subnet" "imported" {
    allow_ip_requests      = false
    create_ptr_records     = false
    display_hostnames      = false
    gateway                = {}
    host_discovery_enabled = false
    id                     = "20"
    include_in_ping        = false
    is_folder              = false
    is_full                = false
    linked_subnet_id       = 0
    location_id            = 0
    master_subnet_id       = 8
    nameserver_id          = 0
    nameservers            = {}
    parent_subnet_id       = 8
    permissions            = jsonencode(
        {
            "2" = "2"
            "3" = "1"
        }
    )
    scan_agent_id          = 0
    section_id             = 1
    show_name              = false
    subnet_address         = "172.20.0.0"
    subnet_id              = 20
    subnet_mask            = 24
    utilization_threshold  = 0
    vlan_id                = 0
    vrf_id                 = 0
}
```

## LICENSE

> Copyright 2023 lord-kyron
>
> Licensed under the Apache License, Version 2.0 (the "License");
> you may not use this file except in compliance with the License.
> You may obtain a copy of the License at
>
> [http://www.apache.org/licenses/LICENSE-2.0](http://www.apache.org/licenses/LICENSE-2.0)
>
> Unless required by applicable law or agreed to in writing, software
> distributed under the License is distributed on an "AS IS" BASIS,
> WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
> See the License for the specific language governing permissions and
> limitations under the License.
