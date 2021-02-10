terraform {
  required_providers {
    sftp = {
      source  = "danitso/sftp"
    }
    tls = {
      source  = "hashicorp/tls"
      version = "3.0.0"
    }
  }
}
