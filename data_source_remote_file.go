package main

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"strconv"
	"time"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
)

const (
	mkDataSourceFileContents       = "contents"
	mkDataSourceFileHost           = "host"
	mkDataSourceFileHostKey        = "host_key"
	mkDataSourceFileLastModified   = "last_modified"
	mkDataSourceFilePassword       = "password"
	mkDataSourceFilePort           = "port"
	mkDataSourceFilePrivateKey     = "private_key"
	mkDataSourceFileRemoteFilePath = "remote_file_path"
	mkDataSourceFileSize           = "size"
	mkDataSourceFileTimeout        = "timeout"
	mkDataSourceFileUser           = "user"
)

// dataSourceRemoteFile retrieves information about a remote file.
func dataSourceRemoteFile() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			mkDataSourceFileContents: &schema.Schema{
				Type:        schema.TypeString,
				Description: "The file contents",
				Computed:    true,
			},
			mkDataSourceFileHost: &schema.Schema{
				Type:        schema.TypeString,
				Description: "The hostname",
				Required:    true,
				ForceNew:    true,
			},
			mkDataSourceFileHostKey: &schema.Schema{
				Type:        schema.TypeString,
				Description: "The host key",
				Optional:    true,
				Default:     "",
				ForceNew:    true,
			},
			mkDataSourceFileLastModified: &schema.Schema{
				Type:        schema.TypeString,
				Description: "The last modified timestamp",
				Computed:    true,
				ForceNew:    true,
			},
			mkDataSourceFilePassword: &schema.Schema{
				Type:        schema.TypeString,
				Description: "The password",
				Optional:    true,
				Default:     "",
				ForceNew:    true,
			},
			mkDataSourceFilePort: &schema.Schema{
				Type:        schema.TypeInt,
				Description: "The port number",
				Optional:    true,
				Default:     22,
				ForceNew:    true,
			},
			mkDataSourceFilePrivateKey: &schema.Schema{
				Type:        schema.TypeString,
				Description: "The private key",
				Optional:    true,
				Default:     "",
				ForceNew:    true,
			},
			mkDataSourceFileRemoteFilePath: &schema.Schema{
				Type:        schema.TypeString,
				Description: "The remote file path",
				Required:    true,
				ForceNew:    true,
			},
			mkDataSourceFileSize: &schema.Schema{
				Type:        schema.TypeString,
				Description: "The file size (in bytes)",
				Computed:    true,
				ForceNew:    true,
			},
			mkDataSourceFileTimeout: &schema.Schema{
				Type:        schema.TypeString,
				Description: "The connection timeout",
				Optional:    true,
				Default:     "5m",
				ForceNew:    true,
			},
			mkDataSourceFileUser: &schema.Schema{
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
	host := d.Get(mkDataSourceFileHost).(string)
	hostKey := d.Get(mkDataSourceFileHostKey).(string)
	password := d.Get(mkDataSourceFilePassword).(string)
	port := d.Get(mkDataSourceFilePort).(int)
	privateKey := d.Get(mkDataSourceFilePrivateKey).(string)
	timeout, err := time.ParseDuration(d.Get(mkDataSourceFileTimeout).(string))
	username := d.Get(mkDataSourceFileUser).(string)

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
	remoteFilePath := d.Get(mkDataSourceFileRemoteFilePath).(string)

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
		return err
	}

	remoteFile, err := sftpClient.Open(remoteFilePath)

	if err != nil {
		return err
	}

	defer remoteFile.Close()

	buffer := bytes.NewBuffer(nil)
	_, err = io.Copy(buffer, remoteFile)

	d.SetId(remoteFileInfo.Name())

	d.Set(mkDataSourceFileContents, buffer.String())
	d.Set(mkDataSourceFileLastModified, remoteFileInfo.ModTime().Format(time.RFC3339))
	d.Set(mkDataSourceFileSize, strconv.FormatInt(remoteFileInfo.Size(), 10))

	return nil
}
