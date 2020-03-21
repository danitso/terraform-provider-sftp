/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/. */

package sftptf

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"time"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
)

const (
	mkDataSourceRemoteFileAllowMissing = "allow_missing"
	mkDataSourceRemoteFileContents     = "contents"
	mkDataSourceRemoteFileHost         = "host"
	mkDataSourceRemoteFileHostKey      = "host_key"
	mkDataSourceRemoteFileLastModified = "last_modified"
	mkDataSourceRemoteFilePassword     = "password"
	mkDataSourceRemoteFilePath         = "path"
	mkDataSourceRemoteFilePort         = "port"
	mkDataSourceRemoteFilePrivateKey   = "private_key"
	mkDataSourceRemoteFileSize         = "size"
	mkDataSourceRemoteFileTimeout      = "timeout"
	mkDataSourceRemoteFileTriggers     = "triggers"
	mkDataSourceRemoteFileUser         = "user"
)

// dataSourceRemoteFile retrieves information about a remote file.
func dataSourceRemoteFile() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			mkDataSourceRemoteFileAllowMissing: {
				Type:        schema.TypeBool,
				Description: "Whether to ignore that the file is missing",
				Optional:    true,
				Default:     false,
				ForceNew:    true,
			},
			mkDataSourceRemoteFileContents: {
				Type:        schema.TypeString,
				Description: "The file contents",
				Computed:    true,
				ForceNew:    true,
			},
			mkDataSourceRemoteFileHost: {
				Type:        schema.TypeString,
				Description: "The hostname",
				Required:    true,
				ForceNew:    true,
			},
			mkDataSourceRemoteFileHostKey: {
				Type:        schema.TypeString,
				Description: "The host key",
				Optional:    true,
				Default:     "",
				ForceNew:    true,
			},
			mkDataSourceRemoteFileLastModified: {
				Type:        schema.TypeString,
				Description: "The last modified timestamp",
				Computed:    true,
				ForceNew:    true,
			},
			mkDataSourceRemoteFilePassword: {
				Type:        schema.TypeString,
				Description: "The password",
				Optional:    true,
				Default:     "",
				ForceNew:    true,
			},
			mkDataSourceRemoteFilePath: {
				Type:        schema.TypeString,
				Description: "The file path",
				Required:    true,
				ForceNew:    true,
			},
			mkDataSourceRemoteFilePort: {
				Type:        schema.TypeInt,
				Description: "The port number",
				Optional:    true,
				Default:     22,
				ForceNew:    true,
			},
			mkDataSourceRemoteFilePrivateKey: {
				Type:        schema.TypeString,
				Description: "The private key",
				Optional:    true,
				Default:     "",
				ForceNew:    true,
			},
			mkDataSourceRemoteFileSize: {
				Type:        schema.TypeInt,
				Description: "The file size (in bytes)",
				Computed:    true,
				ForceNew:    true,
			},
			mkDataSourceRemoteFileTimeout: {
				Type:        schema.TypeString,
				Description: "The connect timeout",
				Optional:    true,
				Default:     "5m",
				ForceNew:    true,
			},
			mkDataSourceRemoteFileTriggers: {
				Type:        schema.TypeMap,
				Description: "The triggers",
				Optional:    true,
				ForceNew:    true,
			},
			mkDataSourceRemoteFileUser: {
				Type:        schema.TypeString,
				Description: "The username",
				Optional:    true,
				Default:     "",
				ForceNew:    true,
			},
		},

		Read: dataSourceRemoteFileRead,
	}
}

// dataSourceRemoteFileCreateSSHClient creates a new SSH client.
func dataSourceRemoteFileCreateSSHClient(d *schema.ResourceData, m interface{}) (*ssh.Client, error) {
	host := d.Get(mkDataSourceRemoteFileHost).(string)
	hostKey := d.Get(mkDataSourceRemoteFileHostKey).(string)
	password := d.Get(mkDataSourceRemoteFilePassword).(string)
	port := d.Get(mkDataSourceRemoteFilePort).(int)
	privateKey := d.Get(mkDataSourceRemoteFilePrivateKey).(string)
	timeout, err := time.ParseDuration(d.Get(mkDataSourceRemoteFileTimeout).(string))
	username := d.Get(mkDataSourceRemoteFileUser).(string)

	if err != nil {
		return nil, err
	}

	if password == "" && privateKey == "" {
		return nil, errors.New("No password or private key has been specified")
	}

	var authMethod []ssh.AuthMethod

	if password != "" {
		authMethod = []ssh.AuthMethod{ssh.Password(password)}
	} else {
		privateKeySigner, err := ssh.ParsePrivateKey([]byte(privateKey))

		if err != nil {
			return nil, err
		}

		authMethod = []ssh.AuthMethod{ssh.PublicKeys(privateKeySigner)}
	}

	var hostKeyCallback ssh.HostKeyCallback

	if hostKey != "" {
		parsedHostKey, err := ssh.ParsePublicKey([]byte(hostKey))

		if err != nil {
			return nil, err
		}

		hostKeyCallback = ssh.FixedHostKey(parsedHostKey)
	} else {
		hostKeyCallback = ssh.InsecureIgnoreHostKey()
	}

	sshConfig := &ssh.ClientConfig{
		User:            username,
		Auth:            authMethod,
		HostKeyCallback: hostKeyCallback,
	}

	timeDelay := int64(10)
	timeMax := timeout.Seconds()
	timeStart := time.Now()
	timeElapsed := timeStart.Sub(timeStart)

	err = nil

	var client *ssh.Client

	for timeElapsed.Seconds() < timeMax {
		if int64(timeElapsed.Seconds())%timeDelay == 0 {
			client, err = ssh.Dial("tcp", fmt.Sprintf("%s:%d", host, port), sshConfig)

			if err == nil {
				break
			}

			time.Sleep(1 * time.Second)
		}

		time.Sleep(200 * time.Millisecond)

		timeElapsed = time.Now().Sub(timeStart)
	}

	if err != nil {
		return nil, err
	}

	return client, nil
}

// dataSourceRemoteFileRead reads information about a remote file.
func dataSourceRemoteFileRead(d *schema.ResourceData, m interface{}) error {
	allowMissing := d.Get(mkDataSourceRemoteFileAllowMissing).(bool)
	remoteFilePath := d.Get(mkDataSourceRemoteFilePath).(string)

	// Create a new SFTP client.
	sshClient, err := dataSourceRemoteFileCreateSSHClient(d, m)

	if err != nil {
		return err
	}

	defer sshClient.Close()

	sftpClient, err := sftp.NewClient(sshClient)

	if err != nil {
		return err
	}

	defer sftpClient.Close()

	// Retrieve the information for a remote file as well as its data.
	remoteFileInfo, err := sftpClient.Lstat(remoteFilePath)

	if err != nil {
		if allowMissing {
			d.SetId("missing")

			d.Set(mkDataSourceRemoteFileContents, "")
			d.Set(mkDataSourceRemoteFileLastModified, time.Now().Format(time.RFC3339))
			d.Set(mkDataSourceRemoteFileSize, -1)

			return nil
		}

		return err
	}

	remoteFile, err := sftpClient.Open(remoteFilePath)

	if err != nil {
		return err
	}

	defer remoteFile.Close()

	buffer := bytes.NewBuffer(nil)
	_, err = io.Copy(buffer, remoteFile)

	if err != nil {
		return err
	}

	d.SetId(remoteFileInfo.Name())

	d.Set(mkDataSourceRemoteFileContents, buffer.String())
	d.Set(mkDataSourceRemoteFileLastModified, remoteFileInfo.ModTime().Format(time.RFC3339))
	d.Set(mkDataSourceRemoteFileSize, int(remoteFileInfo.Size()))

	return nil
}
