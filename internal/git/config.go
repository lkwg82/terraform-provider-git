package git

import (
	"crypto/tls"
	"fmt"
	"net/http"

	"github.com/go-git/go-git/v5/plumbing/transport"
	"github.com/go-git/go-git/v5/plumbing/transport/client"
	githttp "github.com/go-git/go-git/v5/plumbing/transport/http"
)

type Config struct {
	Username              string
	Password              string
	PrivateKey            string
	PrivateKeyFile        string
	InsecureIgnoreHostKey bool
	InsecureSkipTLSVerify bool
}

type Meta struct {
	Auth transport.AuthMethod
}

func (c *Config) Client() (interface{}, error) {
	var meta Meta

	if c.PrivateKey != "" {
		auth, err := getSSHKey(c.PrivateKey, c.InsecureIgnoreHostKey)
		if err != nil {
			return nil, fmt.Errorf("unable to get ssh key: %s", err)
		}

		meta.Auth = auth
	}

	if c.PrivateKeyFile != "" {
		auth, err := getSSHKeyFromFile(c.PrivateKeyFile, c.InsecureIgnoreHostKey)
		if err != nil {
			return nil, fmt.Errorf("unable to get ssh key from file: %s", err)
		}

		meta.Auth = auth
	}

	if c.Username != "" && c.Password != "" {
		meta.Auth = &githttp.BasicAuth{
			Username: c.Username,
			Password: c.Password,
		}
	}

	httpsClient := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: c.InsecureSkipTLSVerify},
		},
	}

	client.InstallProtocol("https", githttp.NewClient(httpsClient))

	return &meta, nil
}
