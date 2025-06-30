package common

import (
	"fmt"
	"os/user"
	"path"
)

func FromHome(paths ...string) ([]string, error) {
	usr, err := user.Current()
	if err != nil {
		return nil, fmt.Errorf("getting the current user: %w", err)
	}

	res := make([]string, 0, len(paths))

	for _, item := range paths {
		res = append(res, path.Join(usr.HomeDir, item))
	}

	return res, nil
}
