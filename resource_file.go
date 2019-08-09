package main

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"strconv"
	"time"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
)

const (
	mkResourceFileContents          = "contents"
	mkResourceFileDestroyLocalFile  = "destroy_local_file"
	mkResourceFileDestroyRemoteFile = "destroy_remote_file"
	mkResourceFileDownload          = "download"
	mkResourceFileHost              = "host"
	mkResourceFileHostKey           = "host_key"
	mkResourceFileLastModified      = "last_modified"
	mkResourceFileLocalFilePath     = "local_file_path"
	mkResourceFilePassword          = "password"
	mkResourceFilePort              = "port"
	mkResourceFilePrivateKey        = "private_key"
	mkResourceFileRemoteFilePath    = "remote_file_path"
	mkResourceFileSize              = "size"
	mkResourceFileTimeout           = "timeout"
	mkResourceFileTriggers          = "triggers"
	mkResourceFileUser              = "user"
)

// resourceFile manages a file.
func resourceFile() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			mkResourceFileContents: &schema.Schema{
				Type:        schema.TypeString,
				Description: "The file contents",
				Optional:    true,
				ForceNew:    true,
			},
			mkResourceFileDestroyLocalFile: &schema.Schema{
				Type:        schema.TypeBool,
				Description: "Whether to destroy the local file when the resource is destroyed",
				Optional:    true,
				Default:     false,
			},
			mkResourceFileDestroyRemoteFile: &schema.Schema{
				Type:        schema.TypeBool,
				Description: "Whether to destroy the remote file when the resource is destroyed",
				Optional:    true,
				Default:     false,
			},
			mkResourceFileDownload: &schema.Schema{
				Type:        schema.TypeBool,
				Description: "Whether to download the file",
				Optional:    true,
				Default:     true,
			},
			mkResourceFileHost: &schema.Schema{
				Type:        schema.TypeString,
				Description: "The hostname",
				Required:    true,
				ForceNew:    true,
			},
			mkResourceFileHostKey: &schema.Schema{
				Type:        schema.TypeString,
				Description: "The host key",
				Optional:    true,
				Default:     "",
			},
			mkResourceFileLastModified: &schema.Schema{
				Type:        schema.TypeString,
				Description: "The last modified timestamp",
				Computed:    true,
				ForceNew:    true,
			},
			mkResourceFileLocalFilePath: &schema.Schema{
				Type:        schema.TypeString,
				Description: "The local file path",
				Optional:    true,
				Default:     "",
				ForceNew:    true,
			},
			mkResourceFilePassword: &schema.Schema{
				Type:        schema.TypeString,
				Description: "The password",
				Optional:    true,
				Default:     "",
			},
			mkResourceFilePort: &schema.Schema{
				Type:        schema.TypeInt,
				Description: "The port number",
				Optional:    true,
				Default:     22,
			},
			mkResourceFilePrivateKey: &schema.Schema{
				Type:        schema.TypeString,
				Description: "The private key",
				Optional:    true,
				Default:     "",
			},
			mkResourceFileRemoteFilePath: &schema.Schema{
				Type:        schema.TypeString,
				Description: "The remote file path",
				Required:    true,
				ForceNew:    true,
			},
			mkResourceFileSize: &schema.Schema{
				Type:        schema.TypeString,
				Description: "The file size in bytes",
				Computed:    true,
				ForceNew:    true,
			},
			mkResourceFileTimeout: &schema.Schema{
				Type:        schema.TypeString,
				Description: "The connection timeout",
				Optional:    true,
				Default:     "5m",
			},
			mkResourceFileTriggers: &schema.Schema{
				Type:     schema.TypeMap,
				Optional: true,
				ForceNew: true,
			},
			mkResourceFileUser: &schema.Schema{
				Type:        schema.TypeString,
				Description: "The username",
				Optional:    true,
				Default:     "",
			},
		},

		Create: resourceFileCreate,
		Read:   resourceFileRead,
		Delete: resourceFileDelete,
	}
}

// resourceFileCreate initializes a file download or upload.
func resourceFileCreate(d *schema.ResourceData, m interface{}) error {
	download := d.Get(mkResourceFileDownload).(bool)
	localFilePath := d.Get(mkResourceFileLocalFilePath).(string)
	remoteFilePath := d.Get(mkResourceFileRemoteFilePath).(string)

	// Create a new SFTP client.
	sshClient, err := resourceFileCreateSSHClient(d, m)

	if err != nil {
		return err
	}

	defer sshClient.Close()

	sftpClient, err := sftp.NewClient(sshClient)

	if err != nil {
		return err
	}

	defer sftpClient.Close()

	// Download or upload the file.
	var remoteFile *sftp.File

	if download {
		remoteFile, err = sftpClient.Open(remoteFilePath)
	} else {
		remoteFile, err = sftpClient.Create(remoteFilePath)
	}

	if err != nil {
		return err
	}

	defer remoteFile.Close()

	if localFilePath != "" {
		localFile, err := os.OpenFile(localFilePath, os.O_RDWR, 0644)

		if err != nil {
			return err
		}

		defer localFile.Close()

		if download {
			_, err = io.Copy(localFile, remoteFile)
		} else {
			_, err = io.Copy(remoteFile, localFile)
		}

		if err != nil {
			return err
		}
	} else {
		if download {
			buffer := bytes.NewBuffer(nil)
			_, err = io.Copy(buffer, remoteFile)

			d.Set(mkResourceFileContents, buffer.String())
		} else {
			buffer := bytes.NewBufferString(d.Get(mkResourceFileContents).(string))
			_, err = io.Copy(remoteFile, buffer)
		}

		if err != nil {
			return err
		}
	}

	// Retrieve information about the file.
	remoteFileInfo, err := sftpClient.Lstat(remoteFilePath)

	if err != nil {
		return err
	}

	d.SetId(remoteFileInfo.Name())

	d.Set(mkResourceFileLastModified, remoteFileInfo.ModTime().Format(time.RFC3339))
	d.Set(mkResourceFileSize, strconv.FormatInt(remoteFileInfo.Size(), 10))

	return nil
}

// resourceFileCreateSSHClient creates a new SSH client.
func resourceFileCreateSSHClient(d *schema.ResourceData, m interface{}) (*ssh.Client, error) {
	host := d.Get(mkResourceFileHost).(string)
	hostKey := d.Get(mkResourceFileHostKey).(string)
	password := d.Get(mkResourceFilePassword).(string)
	port := d.Get(mkResourceFilePort).(int)
	privateKey := d.Get(mkResourceFilePrivateKey).(string)
	timeout, err := time.ParseDuration(d.Get(mkResourceFileTimeout).(string))
	username := d.Get(mkResourceFileUser).(string)

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

// resourceFileRead reads information about a remote file.
func resourceFileRead(d *schema.ResourceData, m interface{}) error {
	download := d.Get(mkResourceFileDownload).(bool)
	localFilePath := d.Get(mkResourceFileLocalFilePath).(string)
	remoteFilePath := d.Get(mkResourceFileRemoteFilePath).(string)

	// Create a new SFTP client.
	sshClient, err := resourceFileCreateSSHClient(d, m)

	if err != nil {
		return err
	}

	defer sshClient.Close()

	sftpClient, err := sftp.NewClient(sshClient)

	if err != nil {
		return err
	}

	defer sftpClient.Close()

	// Retrieve information about the files in order to determine, if the resource needs to be re-created.
	remoteFileInfo, err := sftpClient.Lstat(remoteFilePath)

	if err != nil {
		if !download {
			d.SetId("")

			return nil
		}

		return err
	}

	if localFilePath != "" {
		localFileInfo, err := os.Stat(localFilePath)

		if err != nil {
			if os.IsNotExist(err) && download {
				d.SetId("")

				return nil
			}

			return err
		}

		if localFileInfo.Size() != remoteFileInfo.Size() {
			d.SetId("")

			return nil
		}
	} else {
		contents := d.Get(mkResourceFileContents).(string)

		if int64(len(contents)) != remoteFileInfo.Size() {
			d.SetId("")

			return nil
		}
	}

	d.Set(mkResourceFileLastModified, remoteFileInfo.ModTime().Format(time.RFC3339))
	d.Set(mkResourceFileSize, strconv.FormatInt(remoteFileInfo.Size(), 10))

	return nil
}

// resourceFileDelete cleans up after a file download or upload.
func resourceFileDelete(d *schema.ResourceData, m interface{}) error {
	if d.Get(mkResourceFileDestroyLocalFile).(bool) {
		localFilePath := d.Get(mkResourceFileLocalFilePath).(string)

		if localFilePath != "" {
			if _, err := os.Stat(localFilePath); err == nil {
				err := os.Remove(localFilePath)

				if err != nil {
					return err
				}
			}
		}
	}

	if d.Get(mkResourceFileDestroyRemoteFile).(bool) {
		remoteFilePath := d.Get(mkResourceFileRemoteFilePath).(string)

		// Create a new SFTP client and delete the remote file.
		sshClient, err := resourceFileCreateSSHClient(d, m)

		if err != nil {
			return err
		}

		defer sshClient.Close()

		sftpClient, err := sftp.NewClient(sshClient)

		if err != nil {
			return err
		}

		defer sftpClient.Close()

		if _, err := sftpClient.Lstat(remoteFilePath); err == nil {
			err = sftpClient.Remove(remoteFilePath)

			if err != nil {
				return err
			}
		}
	}

	d.SetId("")

	return nil
}
