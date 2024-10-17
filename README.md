# Terraform Provider for Defined.net

Documentation: https://registry.terraform.io/providers/smaily/definednet/latest/docs

## Requirements

- [Terraform](https://developer.hashicorp.com/terraform/downloads)
  - HashiCorp recommends to use the two latest terraform releases (1.8.x, 1.9.x). Our test suite validates that our provider works with these versions.
  - This provider uses the [terraform plugin protocol version 6](https://developer.hashicorp.com/terraform/plugin/terraform-plugin-protocol#protocol-version-6), and should work with all tools (ie. Terraform & OpenTofu) that supports it.
- [Go](https://go.dev/doc/install) >= 1.22 (to build the provider plugin).


## Building The Provider

1. Clone the repository
1. Enter the repository directory
1. Build the provider using the Go `install` command:

```shell
go install
```

## Adding Dependencies

This provider uses [Go modules](https://github.com/golang/go/wiki/Modules).
Please see the Go documentation for the most up to date information about using Go modules.

To add a new dependency `github.com/author/dependency` to your Terraform provider:

```shell
go get github.com/author/dependency
go mod tidy
```

Then commit the changes to `go.mod` and `go.sum`.

## Developing the Provider

If you wish to work on the provider, you'll first need [Go](http://www.golang.org) installed on your machine (see [Requirements](#requirements) above).

To compile the provider, run `go install`. This will build the provider and put the provider binary in the `$GOPATH/bin` directory.

To generate or update documentation, run `make generate`.

In order to run the full suite of Acceptance tests, run `make testacc`.

*Note:* Acceptance tests create real resources, and often cost money to run.

```shell
make testacc
```

> Consult the official documentation for more information on the test framework's internals: https://developer.hashicorp.com/terraform/plugin/sdkv2/testing/acceptance-tests.
