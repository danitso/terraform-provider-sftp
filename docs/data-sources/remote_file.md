---
layout: page
title: sftp_remote_file
permalink: /data-sources/remote_file
nav_order: 1
parent: Data Sources
---

# Data Source: sftp_remote_file

Retrieves the contents of a remote file.

## Example Usage

```
data "sftp_remote_file" "some_configuration_file" {
  host        = "10.0.0.2"
  user        = "root"
  path        = "/etc/some/configuration.file"
  private_key = "${tls_private_key.automation_key.private_key_pem}"
}

resource "tls_private_key" "automation_key" {
  algorithm = "RSA"
  rsa_bits  = 4096
}
```

## Argument Reference

* `allow_missing` - (Optional) Whether to ignore that the file is missing. Defaults to `false`.
* `host` - (Required) The remote host.
* `host_key` - (Optional) The remote host's key. Defaults to an empty string.
* `password` - (Optional) The password for the remote host. Defaults to an empty string (use `private_key` for key based authorization).
* `path` - (Required) The absolute path to the file.
* `port` - (Optional) The port number for the remote host.
* `private_key` - (Optional) The private key for the remote host. Defaults to an empty string (use `password` for regular password authorization).
* `timeout` - (Optional) The connect timeout. Defaults to `5m` (5 minutes).
* `triggers` - (Optional) The triggers.
* `user` - (Required) The username for the remote host.

## Attribute Reference

* `contents` - The file contents.
* `last_modified` - The last modified timestamp of the file.
* `size` - The file size (in bytes).
