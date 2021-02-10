[![Build Status](https://api.travis-ci.com/danitso/terraform-provider-sftp.svg?branch=master)](https://travis-ci.com/danitso/terraform-provider-sftp)
[![Go Report Card](https://goreportcard.com/badge/github.com/danitso/terraform-provider-sftp)](https://goreportcard.com/report/github.com/danitso/terraform-provider-sftp)
[![GoDoc](https://godoc.org/github.com/danitso/terraform-provider-sftp?status.svg)](http://godoc.org/github.com/danitso/terraform-provider-sftp)

# Terraform Provider for SFTP
A Terraform Provider which adds additional SFTP functionality.

## Requirements

- [Terraform](https://www.terraform.io/downloads.html) 0.13+
- [Go](https://golang.org/doc/install) 1.15+ (to build the provider plugin)
- [GoReleaser](https://goreleaser.com/install/) 0.155+ (to build the provider plugin)

## Building the Provider
- Clone the repository to `$GOPATH/src/github.com/danitso/terraform-provider-sftp`:

    ```sh
    $ mkdir -p "${GOPATH}/src/github.com/danitso"
    $ cd "${GOPATH}/src/github.com/danitso"
    $ git clone git@github.com:danitso/terraform-provider-sftp
    ```

- Enter the provider directory and build it:

    ```sh
    $ cd "${GOPATH}/src/github.com/danitso/terraform-provider-sftp"
    $ make build
    ```

## Using the Provider
You can find the latest release in the [Terraform Registry](https://registry.terraform.io/providers/danitso/sftp/latest).

## Testing the Provider
In order to test the provider, you can simply run `make test`.

```sh
$ make test
```

Tests are limited to regression tests, ensuring backwards compability.
