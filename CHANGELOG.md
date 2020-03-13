## 1.0
New stuff:
 * Added support for terraform multiple resource deployments (count) [Thanks to new repository contributor @pavel-z1]
 * In the process @pavel-z1 forked original phpipam-sdk-go repo and commited some changes. All links leading to the original repo were changed to lead to the updated one.
 * In the process @pavel-z1 found a bug in phpipam itself and reported it. The bug was fixed and backported for version 1.3 and 1.4, BUT keep in mind, that if you want to use the module with "count", you must apply this fix: https://github.com/phpipam/phpipam/commit/b634cb9e4e7df655f219d57e50b813733fd45afc
 otherwise you will have to run terraform apply always with "parallelism=1" or phpipam prior to version 1.5 will not be able to handle the request. More info can be found in the bug report here:
 https://github.com/phpipam/phpipam/issues/2960
 * Build around latest Terraform version 0.12.23 (still API version 5)
 * Build with latest golang v 1.14
 
## 0.3.1
New stuff:
 * Fixed paths across the whole code
 * Build around Terraform version 0.12.16 (latest release, but still API version 5)
 * Build with golang v 1.12.13

## 0.3.0
New version supporting new Terraform API version 5
 * The whole provider was re-build around teeraform source code for version 0.12.x+
   which add support (and makes possible to run the provider) on Terraform 0.12.x+

## 0.2.0

New version with some breaking changes regarding custom fields:

 * Custom field searches are not done via the `custom_field_filter` attribute in
   both the `phpipam_address` and `phpipam_subnet` data sources. This parameter
   is a map that takes custom field names, and regular expressions to search for
   against. A field does not match if any of the search criteria keys are
   missing or do not match.

Also have added two new data sources:

 * `phpipam_address` will search addresses for a `description`  or `hostname`
   exact match or a `custom_field_filter` match, much like the singular-form
   `phpipam_address` data source. A list of IP address IDs are returned, which
   can then be used to look up addresses with the `phpipam_address` data source.
   This will work better in Terraform v0.9.x, or higher, which has support for a
   computed `count` in data sources now.
 * `phpipam_subnets` has been added in the same way. This one can search on
   `description`, `description_match`, and `custom_field_filter`, in the same
   way the singular-form `phpipam_address` data source can.

## 0.1.2

Added custom field support - this plugin now supports custom fields in
addresses, subnets, and VLANs, as long as those fields are optional. Data source
searching supports addresses and subnets only, due to limitations in VLAN
searching capabilities.

## 0.1.1

Bumping release so that I have a consistent snapshot, and also so that I can
correct some tests on the compat branch.

## 0.1.0

First release!
