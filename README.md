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

## Resources

### File (sftp_file)

#### Arguments

* `contents` - (Optional) The file contents. Defaults to an empty string. Leave blank for file downloads (see `download` argument).
* `destroy_local_file` - (Optional) Whether to destroy the local file when the resource is destroyed. Defaults to `false`.
* `destroy_remote_file` - (Optional) Whether to destroy the remote file when the resource is destroyed. Defaults to `false`.
* `download` - (Optional) Whether to download the specified file. Defaults to `true`.
* `host` - (Required) The remote host.
* `host_key` - (Optional) The remote host's key. Defaults to an empty string.
* `local_file_path` - (Optional) The absolute path to a local file. The remote file contents will be written to this file, if `download` is set to `true`. Otherwise, the remote file contents will be retrieved from this file. Leave blank in order to use the `contents` argument.
* `password` - (Optional) The password for the remote host. Defaults to an empty string (use `private_key` as an alternative).
* `port` - (Optional) The port number for the remote host.
* `private_key` - (Optional) The private key for the remote host. Defaults to an empty string (use `password` as an alternative).
* `remote_file_path` - (Required) The absolute path to a remote file.
* `timeout` - (Optional) The connect timeout. Defaults to `5m` (5 minutes).
* `triggers` - (Optional) The key-value map to use as triggers.
* `user` - (Required) The username for the remote host.

#### Attributes

* `contents` - The file contents, if a file is being downloaded and if the `local_file_path` argument is undefined.
* `last_modified` - The last modified timestamp of the file.
* `size` - The size in bytes of the file.

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
