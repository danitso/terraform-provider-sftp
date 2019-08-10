package main

import (
	"github.com/hashicorp/terraform/helper/schema"
)

// Provider returns the object for this provider.
func Provider() *schema.Provider {
	return &schema.Provider{
		ConfigureFunc: providerConfigure,
		DataSourcesMap: map[string]*schema.Resource{
			"sftp_remote_file": dataSourceRemoteFile(),
		},
		ResourcesMap: map[string]*schema.Resource{},
		Schema:       map[string]*schema.Schema{},
	}
}

// providerConfigure() configures the provider before processing any IronMQ resources.
func providerConfigure(d *schema.ResourceData) (interface{}, error) {
	return nil, nil
}
