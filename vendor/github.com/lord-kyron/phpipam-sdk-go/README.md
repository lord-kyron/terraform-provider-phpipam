[![GoDoc](https://godoc.org/github.com/pavel-z1/phpipam-sdk-go?status.svg)](https://godoc.org/github.com/pavel-z1/phpipam-sdk-go)

# phpipam-sdk-go - Partial SDK for PHPIPAM

`phpipam-sdk-go` is a partial SDK for the [PHPIPAM][1] API.

[1]: https://phpipam.net/api/api_documentation/

This is a WIP and this README along with the rest of the code will develop until
it reaches an acceptable level of maturity that it can be used with some CLI
tools that we are developing to work with PHPIPAM, and possibly a Terraform
provider to help insert data gathered from AWS and beyond.

## Reference

See the [GoDoc][2] for the SDK usage details.

[2]: https://godoc.org/github.com/pavel-z1/phpipam-sdk-go

## A Note on Custom Fields

The controllers in this SDK can access custom fields in one of two ways: using
the embedded `CustomFields` map in each controller's data type, or using the
`Get` and `Update` methods in each controller designed to work with custom
fields. Which one you use depends on if you are using the **Nested custom
fields** feature in PHPIPAM (requires 1.3 or higher). Nested custom fields
require that you use the `CustomFields` map, non-nested require the use of the
aforementioned functions.

Note that when you are using un-nested custom fields, you cannot use required
fields - this is due to the fact that entries get added ahead of time without
custom fields as there is no easy way to predict the shape of the data necessary
to send to PHPIPAM in the initial creation request. If you require required
fields, enable the nested functionality - otherwise, ensure that your fields are
not required and choose sane defaults if it's absolutely necessary for data to
be present.


## License

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
