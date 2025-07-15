package command

import (
	"fmt"

	"github.com/aigic8/corn/lib/config"
	"github.com/melbahja/goph"
)

type RemoteClient = goph.Client

func LoginToRemote(remote *config.Remote) (*RemoteClient, error) {
	var auth goph.Auth
	var err error
	if remote.Auth.KeyAuth != nil {
		auth, err = goph.Key(remote.Auth.KeyAuth.KeyPath, remote.Auth.KeyAuth.Passphrase)
		if err != nil {
			return nil, fmt.Errorf("getting key auth: %w", err)
		}
	} else if remote.Auth.PasswordAuth != nil {
		auth = goph.Password(remote.Auth.PasswordAuth.Password)
	} else {
		return nil, fmt.Errorf("no auth method provided for remote")
	}

	client, err := newGothClient(remote.Username, remote.Address, remote.Port, auth)
	if err != nil {
		return nil, fmt.Errorf("getting ssh client: %w", err)
	}
	return client, nil
}

// This is a replacement for `goph.NewClient`. But gives the option to choose the port
func newGothClient(user, addr string, port uint, auth goph.Auth) (*goph.Client, error) {
	callback, err := goph.DefaultKnownHosts()

	if err != nil {
		return nil, err
	}

	return goph.NewConn(&goph.Config{
		User:     user,
		Addr:     addr,
		Port:     port,
		Auth:     auth,
		Timeout:  goph.DefaultTimeout,
		Callback: callback,
	})
}
