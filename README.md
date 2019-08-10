# Terraform Provider for SFTP
A Terraform Provider which adds SFTP functionality.

# Requirements

- [Terraform](https://www.terraform.io/downloads.html) 0.11+
- [Go](https://golang.org/doc/install) 1.12 (to build the provider plugin)

# Building the Provider
Clone repository to: `$GOPATH/src/github.com/danitso/terraform-provider-sftp`

```sh
$ mkdir -p $GOPATH/src/github.com/danitso; cd $GOPATH/src/github.com/danitso
$ git clone git@github.com:danitso/terraform-provider-sftp
```

Enter the provider directory, initialize and build the provider

```sh
$ cd $GOPATH/src/github.com/danitso/terraform-provider-sftp
$ make init
$ make build
```

# Using the Provider
If you're building the provider, follow the instructions to [install it as a plugin.](https://www.terraform.io/docs/plugins/basics.html#installing-plugins) After placing it into your plugins directory,  run `terraform init` to initialize it.

## Data Sources

### Remote File (sftp_remote_file)

#### Arguments

* `host` - (Required) The remote host.
* `host_key` - (Optional) The remote host's key. Defaults to an empty string.
* `password` - (Optional) The password for the remote host. Defaults to an empty string (use `private_key` for key based authorization).
* `path` - (Required) The absolute path to the file.
* `port` - (Optional) The port number for the remote host.
* `private_key` - (Optional) The private key for the remote host. Defaults to an empty string (use `password` for regular password authorization).
* `timeout` - (Optional) The connect timeout. Defaults to `5m` (5 minutes).
* `user` - (Required) The username for the remote host.

#### Attributes

* `contents` - The file contents.
* `last_modified` - The last modified timestamp of the file.
* `size` - The file size (in bytes).

# Developing the Provider
If you wish to work on the provider, you'll first need [Go](http://www.golang.org) installed on your machine (version 1.12+ is *required*).
You'll also need to correctly setup a [GOPATH](http://golang.org/doc/code.html#GOPATH), as well as adding `$GOPATH/bin` to your `$PATH`.

To compile the provider, run `make build`. This will build the provider and put the provider binary in the `$GOPATH/bin` directory.

```sh
$ make build
...
$ $GOPATH/bin/terraform-provider-sftp
...
```

If you wish to contribute to the provider, the following requirements must be met,

* All tests must pass using `make test`
* The Go code must be formatted using Gofmt
* Dependencies are installed by `make init`

# Testing the Provider
In order to test the provider, you can simply run `make test`.

```sh
$ make test
```

Tests are limited to regression tests, ensuring backwards compability.
