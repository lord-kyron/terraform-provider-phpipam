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

## Data Source and Resource Documentation

Please see the [documentation directory](./docs/index.md) for how to use this
provider.

## Building

The provider's executable files are hosted on the [Terraform Repository][8] and are
ready to use without additional assembly.
See the [Plugin Basics][5] page of the Terraform docs to see how to plunk this
into your config. Check the [releases page][6] of this repo to get releases for
Linux, OS X, and Windows.

[5]: https://www.terraform.io/docs/plugins/basics.html
[6]: https://github.com/lord-kyron/terraform-provider-phpipam/releases
[8]: https://registry.terraform.io/providers/lord-kyron/phpipam/latest

How to build provider;

Example for Rocky Linux 8:

Build from repo:

```sh
sudo yum install golang git
sudo mkdir -p $HOME/development/terraform-providers/
cd $HOME/development/terraform-providers/
git clone https://github.com/lord-kyron/terraform-provider-phpipam
# In some cases need execute go install twice
go install
go build
cp terraform-provider-phpipam ~/.terraform.d/plugins/local.dev/phpipam/{version}/{os_platform}/
```

## Unit tests

Requirements:

1. Ready for usage phpIPAM instance
2. Created Custom fields for IP address, subnets, vlans objects inside phpIPAM
   should be created ext custom fields:
   Custom IP addresses fields:
     - CustomTestAddresses varchar(30)
     - CustomTestAddresses2 varchar(30)
   Custom Subnets fields:
     - CustomTestSubnets varchar(30)
     - CustomTestSubnets2 varchar(30)
   Custom VLAN fields:
     - CustomTestVLANs varchar(30)
3. Exported environment variable with phpIPAM credentials. Example:

```sh
export PHPIPAM_APP_ID="terraform"
export PHPIPAM_ENDPOINT_ADDR="http://10.10.0.1/api"
export PHPIPAM_PASSWORD="password"
export PHPIPAM_USER_NAME="Admin"
```

To start unit test exec next command:

```sh
make testacc
```

Example unit test results:

```sh
go clean -testcache; TF_ACC=1 go test -v ./plugin/providers/phpipam -run="TestAcc"
=== RUN   TestAccDataSourcePHPIPAMAddress
--- PASS: TestAccDataSourcePHPIPAMAddress (3.00s)
=== RUN   TestAccDataSourcePHPIPAMAddresses
--- PASS: TestAccDataSourcePHPIPAMAddresses (3.71s)
=== RUN   TestAccDataSourcePHPIPAMFirstFreeAddress
--- PASS: TestAccDataSourcePHPIPAMFirstFreeAddress (1.52s)
=== RUN   TestAccDataSourcePHPIPAMFirstFreeAddressNoFree
--- PASS: TestAccDataSourcePHPIPAMFirstFreeAddressNoFree (0.79s)
=== RUN   TestAccDataSourcePHPIPAMFirstFreeSubnet
--- PASS: TestAccDataSourcePHPIPAMFirstFreeSubnet (1.30s)
=== RUN   TestAccDataSourcePHPIPAMFirstFreeSubnetNoFree
--- PASS: TestAccDataSourcePHPIPAMFirstFreeSubnetNoFree (0.69s)
=== RUN   TestAccDataSourcePHPIPAML2Domain
--- PASS: TestAccDataSourcePHPIPAML2Domain (1.01s)
=== RUN   TestAccDataSourcePHPIPAMSection
--- PASS: TestAccDataSourcePHPIPAMSection (1.27s)
=== RUN   TestAccDataSourcePHPIPAMSubnet
--- PASS: TestAccDataSourcePHPIPAMSubnet (1.62s)
=== RUN   TestAccDataSourcePHPIPAMSubnet_CustomFields
--- PASS: TestAccDataSourcePHPIPAMSubnet_CustomFields (1.50s)
=== RUN   TestAccDataSourcePHPIPAMSubnets
--- PASS: TestAccDataSourcePHPIPAMSubnets (2.57s)
=== RUN   TestAccDataSourcePHPIPAMVLAN
--- PASS: TestAccDataSourcePHPIPAMVLAN (1.07s)
=== RUN   TestAccResourcePHPIPAMAddress
--- PASS: TestAccResourcePHPIPAMAddress (1.25s)
=== RUN   TestAccResourcePHPIPAMOptionalAddress
--- PASS: TestAccResourcePHPIPAMOptionalAddress (1.42s)
=== RUN   TestAccResourcePHPIPAMAddress_CustomFields
--- PASS: TestAccResourcePHPIPAMAddress_CustomFields (2.26s)
=== RUN   TestAccResourcePHPIPAML2Domain
--- PASS: TestAccResourcePHPIPAML2Domain (1.12s)
=== RUN   TestAccResourcePHPIPAMSection
--- PASS: TestAccResourcePHPIPAMSection (1.06s)
=== RUN   TestAccResourcePHPIPAMSubnet
--- PASS: TestAccResourcePHPIPAMSubnet (1.11s)
=== RUN   TestAccResourcePHPIPAMSubnet_CustomFields
--- PASS: TestAccResourcePHPIPAMSubnet_CustomFields (1.95s)
=== RUN   TestAccResourcePHPIPAMVLAN
--- PASS: TestAccResourcePHPIPAMVLAN (1.27s)
PASS
ok    github.com/lord-kyron/terraform-provider-phpipam/plugin/providers/phpipam 31.522s
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
