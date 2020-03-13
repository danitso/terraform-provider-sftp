---
layout: home
title: Introduction
nav_order: 1
---

# SFTP Provider

This provider for [Terraform](https://www.terraform.io/) adds additional SFTP functionality.

Use the navigation to the left to read about the available resources.

## Example Usage

```
provider "sftp" {}
```

## Installation

You can install the latest release of the provider using either Git Bash or regular Bash:

```sh
$ export PROVIDER_PLATFORM="$([[ "$OSTYPE" =~ ^msys|cygwin$ ]] && echo "windows" || ([[ "$OSTYPE" == "darwin"* ]] && echo "darwin" || ([[ "$OSTYPE" == "linux"* ]] && echo "linux" || echo "unsupported")))"
$ export PROVIDER_VERSION="$(curl -L -s -H 'Accept: application/json' https://github.com/danitso/terraform-provider-sftp/releases/latest | sed -e 's/.*"tag_name":"\([^"]*\)".*/\1/')"
$ export PLUGINS_PATH="$([[ "$PROVIDER_PLATFORM" == "windows" ]] && cygpath -u "$APPDATA" || echo "$HOME")/terraform.d/plugins"
$ mkdir -p "$PLUGINS_PATH"
$ curl -o "${PLUGINS_PATH}/terraform-provider-sftp_v${PROVIDER_VERSION}.zip" -sL "https://github.com/danitso/terraform-provider-sftp/releases/download/${PROVIDER_VERSION}/terraform-provider-sftp_v${PROVIDER_VERSION}-custom_${PROVIDER_PLATFORM}_amd64.zip"
$ unzip -o -d "$PLUGINS_PATH" "${PLUGINS_PATH}/terraform-provider-sftp_v${PROVIDER_VERSION}.zip"
$ rm "${PLUGINS_PATH}/terraform-provider-sftp_v${PROVIDER_VERSION}.zip"
```

You can also install it manually by following the instructions to [install it as a plugin](https://www.terraform.io/docs/plugins/basics.html#installing-plugins). You can download the latest release from the [releases](https://github.com/danitso/terraform-provider-sftp/releases) page.
