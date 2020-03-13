# Terraform PHPIPAM provider - version 1.0
# ATTENTION!!! 
This repository is based on the original work of github user paybyphone.
However, the version of the provider in this repo is updated and revised to support working with Terraform 12.x+
The build here is currently based on the paybyphone original repo + hashicorp original terraform repo and was build around Terraform version 0.12.23
All credit should go to https://github.com/paybyphone/terraform-provider-phpipam - I've just modernized his work!
 
# Terraform Provider Plugin for PHPIPAM

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

## Installing

See the [Plugin Basics][5] page of the Terraform docs to see how to plunk this
into your config. Check the [releases page][6] of this repo to get releases for
Linux, OS X, and Windows.

[5]: https://www.terraform.io/docs/plugins/basics.html
[6]: https://github.com/lord-kyron/terraform-provider-phpipam-0.3.1/releases

Examle for CentOS 7:
```
yum install golang git
mkdir -p $HOME/development/terraform-providers/
cd $HOME/development/terraform-providers/
git clone https://github.com/lord-kyron/terraform-provider-phpipam
# In some cases need execute go install twice
go install
go build
cp terraform-provider-phpipam ~/.terraform.d/plugins/
```


## Usage

After installation, to use the plugin, simply use any of its resources or data
sources (such as `phpipam_subnet` or `phpipam_address` in a Terraform
configuration.

Credentials can be supplied via configuration variables to the `phpipam`
provider instance, or via environment variables. These are documented in the
next section.

You can see the following example below for a simple usage example that reserves
the first available IP address in a subnet. This address could then be passed
along to the configuration for a VM, say, for example, a
[`vsphere_virtual_machine`][7] resource.

[7]: https://www.terraform.io/docs/providers/vsphere/r/virtual_machine.html

```
provider "phpipam" {
  app_id   = "test"
  endpoint = "https://phpipam.example.com/api"
  password = "PHPIPAM_PASSWORD"
  username = "Admin"
}

data "phpipam_subnet" "subnet" {
  subnet_address = "10.10.2.0"
  subnet_mask    = 24
}

data "phpipam_first_free_address" "next_address" {
  subnet_id = data.phpipam_subnet.subnet.subnet_id
}

resource "phpipam_address" {
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

### Plugin Options

The options for the plugin are as follows: 

 * `app_id` - The API application ID, configured in the PHPIPAM API panel. This
   application ID should have read/write access if you are planning to use the
   resources, but read-only access should be sufficient if you are only using
   the data sources. Can also be supplied by the `PHPIPAM_APP_ID` environment
   variable.
 * `endpoint` - The full URL to the PHPIPAM API endpoint, such as
   `https://phpipam.example.com/api`. Can also be supplied by the
   `PHPIPAM_ENDPOINT_ADDR` environment variable.
 * `password` - The password to access the PHPIPAM API with. Can also be
   supplied via `PHPIPAM_PASSWORD` to prevent plain text password storage in
   config.
 * `username` - The user name to access the PHPIPAM API with. Can also be
   supplied via the `PHPIPAM_USER_NAME` variable.

### Data Sources

The following data sources are supplied by this plugin:

#### The `phpipam_address` Data Source

The `phpipam_address` data source allows one to get information about a specific
IP address within PHPIPAM. Use this address to get general information about a
specific IP address such as its host name, description and more.

Lookups for IP addresses can only happen at this time via its entry in the
database, or the IP address itself. Future versions of this resource, when such
features become generally available in the PHPIPAM API, will allow lookup based
on host name, allowing for better ability for this resource to discover IP
addresses that have been pre-assigned for a specific resource.

**Example:**

```
data "phpipam_address" "address" {
  ip_address = "10.10.1.1"
}

output "address_description" {
  value = data.phpipam_address.address.description
}
```

**Example With `description`:**

```
data "phpipam_address" "address" {
  subnet_id         = 3
  description_match = "Customer 1"
}

output "address_description" {
  value = data.phpipam_address.address.description
}
```

**Example With `custom_field_filter`:**

```
data "phpipam_address" "address" {
  subnet_id = 3

  custom_field_filter {
    CustomTestAddresses = ".*terraform.*"
  }
}

output "address_description" {
  value = data.phpipam_address.address.description
}
```

##### Argument Reference

The data source takes the following parameters:

 * `address_id` - The ID of the IP address in the PHPIPAM database.
 * `ip_address` - The actual IP address in PHPIPAM.
 * `subnet_id` - The ID of the subnet that the address resides in. This is
   required to search on the `description` or `hostname` fields.
 * `description` - The description of the IP address. `subnet_id` is required
   when using this field.
 * `hostname` - The host name of the IP address. `subnet_id` is required when
   using this field.
 * `custom_field_filter` - A map of custom fields to search for. The filter
   values are regular expressions that follow the RE2 syntax for which you can
   find documentation [here](https://github.com/google/re2/wiki/Syntax). All
   fields need to match for the match to succeed. 

⚠️  **NOTE:** `description`, `hostname`, and `custom_field_filter` fields return
the first match found without any warnings. If you are looking to return
multiple addresses, combine this data source with the `phpipam_addresses` data
source.

⚠️  **NOTE:** An empty or unspecified `custom_field_filter` value is the
equivalent to a regular expression that matches everything, and hence will
return the first address it sees in the subnet.

Arguments are processed in the following order of precedence:

 * `address_id`
 * `ip_address`
 * `subnet_id`, and either one of `description`, `hostname`, or
   `custom_field_filter`

##### Attribute Reference

The following attributes are exported:

 * `address_id` - The ID of the IP address in the PHPIPAM database.
 * `ip_address` - the IP address.
 * `subnet_id` - The database ID of the subnet this IP address belongs to.
 * `is_gateway` - `true` if this IP address has been designated as a gateway.
 * `description` - The description provided to this IP address.
 * `hostname` - The hostname supplied to this IP address.
 * `owner` - The owner name provided to this IP address.
 * `mac_address` - The MAC address provided to this IP address.
 * `state_tag_id` - The tag ID in the database for the IP address' specific
   state. **NOTE:** This is currently represented as an integer but may change
   to the specific string representation at a later time.
 * `skip_ptr_record` - `true` if PTR records are not being created for this IP
   address.
 * `ptr_record_id` - The ID of the associated PTR record in the PHPIPAM
   database.
 * `device_id` - The ID of the associated device in the PHPIPAM database.
 * `switch_port_label` - A string port label that is associated with this
   address.
 * `note` - The note supplied to this IP address.
 * `last_seen` - The last time this IP address answered ping probes.
 * `exclude_ping` - `true` if this address is excluded from ping probes.
 * `edit_date` - The last time this resource was modified.
 * `custom_fields` - A key/value map of custom fields for this address.

##### The `phpipam_addresses` Data Source

The `phpipam_addresses` data source allows you to search for IP addresses, much
in the same way as you can in the single-form `phpipam_address` data source.
However, multiple addresses are returned from this data source as a single list
of address IDs as they are found in the PHPIPAM database. You can then use the
single-form `phpipam_address` data source to extract the IP data for each
matched address in the database.

**Example:**

⚠️  **NOTE:** The below example requires Terraform v0.12.0 or later!

```
data "phpipam_addresses" "address_search" {
  subnet_id = 3

  custom_field_filter {
    CustomTestAddresses = ".*terraform.*"
  }
}

data "phpipam_address" "addresses" {
  count      = length(data.phpipam_addresses.address_search.address_ids)
  address_id = element(data.phpipam_addresses.address_search.address_ids, count.index)
}

output "ip_addresses" {
  value = [data.phpipam_address.addresses.*.ip_address]
}
```

##### Argument Reference

The data source takes the following parameters:

 * `subnet_id` (Required) - The ID of the subnet that the address resides in. This is
   required to search on the `description` or `hostname` fields.

One of the following fields is required alongside `subnet_id`:

 * `description` - The description of the IP address. `subnet_id` is required
   when using this field.
 * `hostname` - The host name of the IP address. `subnet_id` is required when
   using this field.
 * `custom_field_filter` - A map of custom fields to search for. The filter
   values are regular expressions that follow the RE2 syntax for which you can
   find documentation [here](https://github.com/google/re2/wiki/Syntax). All
   fields need to match for the match to succeed. 

⚠️  **NOTE:** An empty or unspecified `custom_field_filter` value is the
equivalent to a regular expression that matches everything, and hence will
return **all** addresses that contain the referenced custom field key!

##### Attribute Reference

The following attributes are exported:

 * `address_ids` - A list of discovered IP address IDs.

##### The `phpipam_first_free_address` Data Source

The `phpipam_first_free_address` data source allows you to get the next
available IP address in a specific subnet in PHPIPAM. Using this resource allows
you to automatically allocate an IP address that can be used as an IP address in
resources such as `vsphere_virtual_machine`, or other virtual machine-like
resources that require static IP addresses.

Note that not having any addresses available will cause the Terraform run to
fail. Conversely, marking a subnet as unavailable or used will not prevent this
data source from returning an IP address, so be aware of this while using this
resource.

**Example:**

```
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
resource "phpipam_address" {
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

##### Argument Reference

The data source takes the following parameters:

 * `subnet_id` - The ID of the subnet to look up the address in.

##### Attribute Reference

The following attributes are exported:

 * `ip_address` - The next available IP address.

#### The `phpipam_section` Data Source

The `phpipam_section` data source allows one to look up a specific section,
either by database ID or name. This data can then be used to manage other parts
of PHPIPAM, such as in the event that the section name is known but not its ID,
which is required for managing subnets.

**Example:**

```
data "phpipam_section" "section" {
  name = "Customers"
}

resource "phpipam_subnet" "subnet" {
  section_id = data.phpipam_section.section.section_id
  subnet_address = "10.10.3.0"
  subnet_mask = 24
}
```

##### Argument Reference

The data source takes the following parameters:

 * `section_id` - The ID of the section to look up.
 * `name` - The name of the section to look up.

One of `section_id` or `name` must be supplied. If both are supplied,
`section_id` is used.

##### Attribute Reference

The following attributes are exported:

 * `section_id` - The ID of the section in the PHPIPAM database.
 * `name` - The name of the section.
 * `description` - The section's description.
 * `master_section_id` - The ID of the parent section in the PHPIPAM database.
 * `permissions` - A JSON representation of permissions for this section.
 * `strict_mode` - `true` if this subnet is set up to check that IP addresses
   are valid for the subnets they are in.
 * `subnet_ordering` - How subnets in this section are ordered.
 * `display_order` - The section's display order number.
 * `edit_date` - The date this resource was last edited.
 * `show_vlan_in_subnet_listing` - `true` if VLANs are being shown in the subnet
   listing for this section.
 * `show_vrf_in_subnet_listing` - `true` if VRFs are being shown in the subnet
   listing for this section.
 * `show_supernet_only` - `true` if supernets are only being shown in the subnet
   listing.
 * `dns_resolver_id` - The ID of the DNS resolver to use in the section.

#### The `phpipam_subnet` Data Source

The `phpipam_subnet` data source gets information on a subnet such as its ID
(required for creating addresses), description, and more.

**Example:**

```
// Look up the subnet
data "phpipam_subnet" "subnet" {
  subnet_address = "10.10.2.0"
  subnet_mask    = 24
}

// Reserve the address.
resource "phpipam_address" {
  subnet_id   = data.phpipam_subnet.subnet.subnet_id
  ip_address  = "10.10.2.10"
  hostname    = "tf-test-host.example.internal"
  description = "Managed by Terraform"
}
```

**Example with `description_match`:**

```
// Look up the subnet (matching on either case of "customer")
data "phpipam_subnet" "subnet" {
  section_id        = 1
  description_match = "[Cc]ustomer 2"
}

// Get the first available address
data "phpipam_first_free_address" "next_address" {
  subnet_id = data.phpipam_subnet.subnet.subnet_id
}

// Reserve the address. Note that we use ignore_changes here to ensure that we
// don't end up re-allocating this address on future Terraform runs.
resource "phpipam_address" {
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

**Example With `custom_field_filter`:**

```
// Look up the subnet
data "phpipam_subnet" "subnet" {
  section_id = 1

  custom_field_filter = {
    CustomTestSubnets = ".*terraform.*"
  }
}

// Get the first available address
data "phpipam_first_free_address" "next_address" {
  subnet_id = data.phpipam_subnet.subnet.subnet_id
}

// Reserve the address. Note that we use ignore_changes here to ensure that we
// don't end up re-allocating this address on future Terraform runs.
resource "phpipam_address" {
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

**Example how to get gateway IP address by subnet_id:**

```
// Get ID of phpIPAM section
data "phpipam_section" "section" {
  name = "Subnet Section"
}

// Look up the subnet parameters
data "phpipam_subnet" "subnet" {
  section_id = data.phpipam_section.section.id
  // prod_mgmt - this is subnet desctiption (subnet name)
  description_match = "prod_mgmt"
}

#// Determine Gateway IP by gateway_id
data "phpipam_address" "gateway" {
  address_id    = data.phpipam_subnet.subnet.gateway_id
}
```

##### Argument Reference

The data source takes the following parameters:

 * `section_id` - The ID of the section of the subnet. Required if you are
   looking up a subnet using the `description` or `description_match` arguments.
 * `subnet_id` - The ID of the subnet to look up.
 * `subnet_address` - The network address of the subnet to look up.
 * `subnet_mask` - The subnet mask, in bits, of the subnet to look up.
 * `description` - The subnet's description. `section_id` is required if you
   want to use this option.
 * `description_match` - A regular expression to match against when searching
   for a subnet. `section_id` is required if you want to use this option.
 * `custom_field_filter` - A map of custom fields to search for. The filter
   values are regular expressions that follow the RE2 syntax for which you can
   find documentation [here](https://github.com/google/re2/wiki/Syntax). All
   fields need to match for the match to succeed. 

⚠️  **NOTE:** Searches with the `description`, `description_match` and
`custom_field_filter` fields return the first match found without any warnings.
Conversely, the resource fails if it somehow finds multiple results on a CIDR
(subnet and mask) search - this is to assert that you are getting the subnet you
requested. If you want to return multiple results, combine this data source with
the `phpipam_subnets` data source.

⚠️  **NOTE:** An empty or unspecified `custom_field_filter` value is the
equivalent to a regular expression that matches everything, and hence will
return the first subnetit sees in the section.

Arguments are processed in the following order of precedence:

 * `subnet_id`
 * `subnet_address` and `subnet_mask`
 * `section_id`, and either one of `description`, `description_match`, or
   `custom_field_filter`

##### Attribute Reference

The following attributes are exported:

 * `subnet_id` - The ID of the subnet in the PHPIPAM database.
 * `subnet_address` - The network address of the subnet.
 * `subnet_mask` - The subnet mask, in bits.
 * `description` - The description set for the subnet.
 * `gateway` - Key map of values: ip_addr, id. ip_addr - this is gateway IP address
    (like 192.168.1.254). id - gateway ip ID
 * `section_id` - The ID of the section for this address in the PHPIPAM
   database.
 * `linked_subnet_id` - The ID of the linked subnet in the PHPIPAM database.
 * `vlan_id` - The ID of the VLAN for this subnet in the PHPIPAM database.
 * `vrf_id` - The ID of the VRF for this subnet in the PHPIPAM database.
 * `master_subnet_id` - The ID of the parent subnet for this subnet in the
   PHPIPAM database.
 * `nameserver_id` - The ID of the nameserver used to assign PTR records for
   this subnet.
 * `show_name` - `true` if the subnet name is are shown in the section, instead
   of the network address.
 * `permissions` - A JSON representation of the permissions associated with this
   subnet.
 * `create_ptr_records` - `true` if PTR records are created for addresses in
   this subnet.
 * `display_hostnames` - `true` if hostnames are displayed instead of IP
   addresses in the address listing for this subnet.
 * `allow_ip_requests` - `true` if the subnet allows IP requests in PHPIPAM.
 * `scan_agent_id` - The ID of the ping scan agent that is used for this subnet.
 * `include_in_ping` - `true` if this subnet is included in ping probes.
 * `host_discovery_enabled` - `true` if this subnet is included in new host
   scans.
 * `is_folder` - `true` if this subnet is a folder and not an actual subnet.
 * `is_full` - `true` if the subnet has been marked as full.
 * `state_tag_id` - The ID of the state tag for this subnet. This may become an
   actual string representation of this at a later time (example: `Used`).
 * `utilization_threshold` - The subnet's utilization threshold.
 * `location_id` - The ID of the location for this subnet.
 * `edit_date` - The date this resource was last updated.
 * `custom_fields` - A key/value map of custom fields for this subnet.
 * `gateway_id` - The ID of gateway IP address fot this subnet

##### The `phpipam_subnets` Data Source

The `phpipam_subnets` data source allows you to search for subnets, much in the
same way as you can in the single-form `phpipam_subnet` data source.  However,
multiple subnets are returned from this data source as a single list of subnet
IDs as they are found in the PHPIPAM database. You can then use the single-form
`phpipam_subnet` data source to extract the subnet data for each matched network
in the database.

**Example:**

⚠️  **NOTE:** The below example requires Terraform v0.12.0 or later!

```
data "phpipam_subnets" "subnet_search" {
  subnet_id = 3

  custom_field_filter {
    CustomTestSubnets = ".*terraform.*"
  }
}

data "phpipam_subnet" "subnets" {
  count      = length(data.phpipam_subnets.subnet_search.subnet_ids)
  address_id = element(data.phpipam_subnets.subnet_search.subnet_ids, count.index)
}

output "subnet_addresses" {
  value = [data.phpipam_subnet.subnets.*.ip_address]
}

output "subnet_cidrs" {
  value = [formatlist("%s/%d", data.phpipam_subnet.subnets.*.subnet_address, data.phpipam_subnet.subnets.*.subnet_mask)]
}
```

##### Argument Reference

The data source takes the following parameters:

 * `section_id` (Required) - The ID of the section of the subnet.

One of the following below parameters is required:

 * `description` - The subnet's description.
 * `description_match` - A regular expression to match against when searching
   for a subnet.
 * `custom_field_filter` - A map of custom fields to search for. The filter
   values are regular expressions. All fields need to match for the match to
   succeed. 

You can find documentation for the regular expression syntax used with the
`description_match` and `custom_field_filter` attributes
[here](https://github.com/google/re2/wiki/Syntax).

⚠️  **NOTE:** An empty or unspecified `custom_field_filter` value is the
equivalent to a regular expression that matches everything, and hence will
return **all** subnets that contain the referenced custom field key!

##### Attribute Reference

The following attributes are exported:

 * `subnet_ids` - A list of subnet IDs that match the given criteria.

#### The `phpipam_vlan` Data Source

The `phpipam_vlan` data source allows one to look up a VLAN in the PHPIPAM
database. This can then be used to assign a VLAN to a subnet in the
`phpipam_subnet` resource. It can also be used to gather other information on
the VLAN.

**Example:**

```
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

##### Argument Reference

The data source takes the following parameters:

 * `vlan_id` - The ID of the VLAN to look up. **NOTE:** this is the database ID,
   not the VLAN number - if you need this, use the `number` parameter.
 * `number` - The number of the VLAN to look up.

One of `vlan_id` or `number` must be supplied. If both are supplied,
`vlan_id` is used.

##### Attribute Reference

The following attributes are exported:

 * `vlan_id` - The ID of the VLAN to look up. **NOTE:** this is the database ID,
   not the VLAN number - if you need this, use the `number` parameter.
 * `number` - The number of the VLAN (the actual VLAN ID on your switch).
 * `l2_domain_id` - The layer 2 domain ID in the PHPIPAM database.
 * `name` - The name/label of the VLAN.
 * `description` - The description supplied to the VLAN.
 * `edit_date` - The date this resource was last updated.
 * `custom_fields` - A key/value map of custom fields for this VLAN.

### Resources

The following resources are supplied by this plugin:

#### The `phpipam_address` Resource

The `phpipam_address` resource manages an IP address in PHPIPAM. You can use it
to create IP address reservations for IP addresses that have been created by
other Terraform resources, or supplied by the `phpipam_first_free_address` data
source. An example usage is below. 

⚠️  **NOTE:** If you are using the `phpipam_first_free_address` to get the first
free IP address in a specific subnet, make sure you set `subnet_id` and
`ip_address` as ignored attributes with the `ignore_changes` lifecycle
attribute. This will prevent Terraform from perpetually deleting and
re-allocating the address when it sees a different available IP address in the
`phpipam_first_free_address` data source.

**Example:**

```
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
resource "phpipam_address" {
  subnet_id   = data.phpipam_subnet.subnet.subnet_id
  ip_address  = data.phpipam_first_free_address.next_address.ip_address
  hostname    = "tf-test-host.example.internal"
  description = "Managed by Terraform"

  custom_fields = {
    CustomTestAddresses = "terraform-test"
  }

  lifecycle {
    ignore_changes = [
      subnet_id,
      ip_address,
    ]
  }
}
```

##### Argument Reference

The resource takes the following parameters:

 * `subnet_id` (Required) - The database ID of the subnet this IP address
   belongs to.
 * `ip_address` (Required) - The IP address to reserve.
 * `is_gateway` (Optional) - `true` if this IP address has been designated as a
   gateway.
 * `description` (Optional) - The description provided to this IP address.
 * `hostname` (Optional) - The hostname supplied to this IP address.
 * `owner` (Optional) - The owner name provided to this IP address.
 * `mac_address` (Optional) - The MAC address provided to this IP address.
 * `state_tag_id` (Optional) - The tag ID in the database for the IP address'
   specific state. **NOTE:** This is currently represented as an integer but may
   change to the specific string representation at a later time.
 * `skip_ptr_record` (Optional) - `true` if PTR records are not being created
   for this IP address.
 * `ptr_record_id` (Optional) - The ID of the associated PTR record in the
   PHPIPAM database.
 * `device_id` (Optional) - The ID of the associated device in the PHPIPAM
   database.
 * `switch_port_label` (Optional) - A string port label that is associated with
   this address.
 * `note` (Optional) - The note supplied to this IP address.
 * `exclude_ping` (Optional) - `true` if this address is excluded from ping
   probes.
 * `remove_dns_on_delete` (Optional) - Removes DNS records created by PHPIPAM
   when the address is deleted from Terraform. Defaults to `true`.
 * `custom_fields` (Optional) -  A key/value map of custom fields for this address.

⚠️  **NOTE on custom fields:** PHPIPAM installations with custom fields must have
all fields set to optional when using this plugin. For more info see
[here](https://github.com/phpipam/phpipam/issues/1073). Further to this, either
ensure that your fields also do not have default values, or ensure the default
is set in your TF configuration. Diff loops may happen otherwise!

##### Attribute Reference

The following attributes are exported:

 * `address_id` - The ID of the IP address in the PHPIPAM database.
 * `last_seen` - The last time this IP address answered ping probes.
 * `edit_date` - The last time this resource was modified.

#### The `phpipam_first_free_address` Resource - Dynamic IPs creation

The `phpipam_first_free_address` resource allow to create automatically 
new IP in defined network without execution Terraform data instruction. You can use it
to create several IP addresses automatically. This resource support the same arguments as
phpipam_address.  An example usage is below.

⚠️  **NOTE:** This is experimental new feature. You can use Terraform count
instruction. But be carefull, phpIPAM currently has a bug https://github.com/phpipam/phpipam/issues/2960
Use resource with count option only with limitted terraform threads count: `**terraform apply -parallelism=1**`
.

**Example create IPs in loop with `count`:**

```
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
```
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
```
 terraform plan
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

#### The `phpipam_section` Resource

The `phpipam_section` resource manages a PHPIPAM section - a top-level category
that subnets and IP addresses are entered into. Use this resource if you want to
manage a section entirely from Terraform. If you just need to get information on
a section use the `phpipam_section` data source instead.

**Example:**

```
// Create a section
resource "phpipam_section" "section" {
  name        = "tf-test"
  description = "Terraform test section"
}
```

##### Argument Reference

The resource takes the following parameters:

 * `name` (Required) - The name of the section.
 * `description` (Optional) - The section's description.
 * `master_section_id` (Optional) - The ID of the parent section in the PHPIPAM
   database.
 * `strict_mode` (Optional) - `true` if this subnet is set up to check that IP
   addresses are valid for the subnets they are in.
 * `subnet_ordering` (Optional) - How subnets in this section are ordered.
 * `display_order` (Optional) - The section's display order number.
 * `show_vlan_in_subnet_listing` (Optional) - `true` if VLANs are being shown in
   the subnet listing for this section.
 * `show_vrf_in_subnet_listing` (Optional) - `true` if VRFs are being shown in
   the subnet listing for this section.
 * `show_supernet_only` (Optional) - `true` if supernets are only being shown in
   the subnet listing.
 * `dns_resolver_id` (Optional) - The ID of the DNS resolver to use in the
   section.

##### Attribute Reference

The following attributes are exported:

 * `section_id` - The ID of the section in the PHPIPAM database.
 * `edit_date` - The date this resource was last edited.

#### The `phpipam_subnet` Resource

The `phpipam_subnet` resource can be used to create and manage a subnet in
PHPIPAM. Use it to manage details on subnets you create in Terraform for other
things, such as storing the IDs of the subnets you create for AWS, or for a full
top-down management of subnets and IP addresses in Terraform. If you just need
to get information on a subnet, use the `phpipam_subnet` data source instead.

**Example:**

```
data "phpipam_section" "section" {
  name = "Customers"
}

resource "phpipam_subnet" "subnet" {
  section_id     = data.phpipam_section.section.section_id
  subnet_address = "10.10.3.0"
  subnet_mask    = 24

  custom_fields = {
    CustomTestSubnets = "terraform-test"
  }
}
```

##### Argument Reference

The resource takes the following parameters:

 * `subnet_address` (Required) - The network address of the subnet.
 * `subnet_mask` (Required) - The subnet mask, in bits.
 * `description` (Optional) - The description set for the subnet.
 * `section_id` (Optional) - The ID of the section for this address in the
   PHPIPAM database.
 * `linked_subnet_id` (Optional) - The ID of the linked subnet in the PHPIPAM
   database.
 * `vlan_id` (Optional) - The ID of the VLAN for this subnet in the PHPIPAM
   database.
 * `vrf_id` (Optional) - The ID of the VRF for this subnet in the PHPIPAM
   database.
 * `master_subnet_id` (Optional) - The ID of the parent subnet for this subnet
   in the PHPIPAM database.
 * `nameserver_id` (Optional) - The ID of the nameserver used to assign PTR
   records for this subnet.
 * `show_name` (Optional) - `true` if the subnet name is are shown in the
   section, instead of the network address.
 * `create_ptr_records` (Optional) - `true` if PTR records are created for
   addresses in this subnet.
 * `display_hostnames` (Optional) - `true` if hostnames are displayed instead of
   IP addresses in the address listing for this subnet.
 * `allow_ip_requests` (Optional) - `true` if the subnet allows IP requests in
   PHPIPAM.
 * `scan_agent_id` (Optional) - The ID of the ping scan agent that is used for
   this subnet.
 * `include_in_ping` (Optional) - `true` if this subnet is included in ping
   probes.
 * `host_discovery_enabled` (Optional) - `true` if this subnet is included in
   new host scans.
 * `is_folder` (Optional) - `true` if this subnet is a folder and not an actual
   subnet.
 * `is_full` (Optional) - `true` if the subnet has been marked as full.
 * `state_tag_id` (Optional) - The ID of the state tag for this subnet. This may
   become an actual string representation of this at a later time (example:
   `Used`).
 * `utilization_threshold` (Optional) - The subnet's utilization threshold.
 * `location_id` (Optional) - The ID of the location for this subnet.
 * `custom_fields` (Optional) -  A key/value map of custom fields for this
   subnet.

⚠️  **NOTE on custom fields:** PHPIPAM installations with custom fields must have
all fields set to optional when using this plugin. For more info see
[here](https://github.com/phpipam/phpipam/issues/1073). Further to this, either
ensure that your fields also do not have default values, or ensure the default
is set in your TF configuration. Diff loops may happen otherwise!

##### Attribute Reference

The following attributes are exported:

 * `subnet_id` - The ID of the subnet in the PHPIPAM database.
 * `permissions` - A JSON representation of the permissions associated with this
   subnet.
 * `edit_date` - The date this resource was last updated.

#### The `phpipam_vlan` Resource

The `phpipam_vlan` resource can be used to manage a VLAN on PHPIPAM. Use it to
set up a VLAN through Terraform, or update details such as its name or
description. If you are just looking for information on a VLAN, use the
`phpipam_vlan` data source instead.

**Example:**

```
resource "phpipam_vlan" "vlan" {
  name        = "tf-test"
  number      = 1000
  description = "Managed by Terraform"

  custom_fields = {
    CustomTestVLANs = "terraform-test"
  }
}
```

##### Argument Reference

The resource takes the following parameters:

 * `name` (Required) - The name/label of the VLAN.
 * `number` (Required) - The number of the VLAN (the actual VLAN ID on your switch).
 * `l2_domain_id` (Optional) - The layer 2 domain ID in the PHPIPAM database.
 * `description` (Optional) - The description supplied to the VLAN.
 * `edit_date` (Optional) - The date this resource was last updated.
 * `custom_fields` (Optional) -  A key/value map of custom fields for this
   VLAN.

⚠️  **NOTE on custom fields:** PHPIPAM installations with custom fields must have
all fields set to optional when using this plugin. For more info see
[here](https://github.com/phpipam/phpipam/issues/1073). Further to this, either
ensure that your fields also do not have default values, or ensure the default
is set in your TF configuration. Diff loops may happen otherwise!

##### Attribute Reference

The following attributes are exported:

 * `vlan_id` - The ID of the VLAN to look up. **NOTE:** this is the database ID,
   not the VLAN number - if you need this, use the `number` parameter.
 * `edit_date` - The date this resource was last updated.

## LICENSE

```
Copyright 2017 PayByPhone Technologies, Inc.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
```
# terraform-provider-phpipam-0.3.1
