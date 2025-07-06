package common

import (
	"fmt"
	"io/fs"
	"os"
	"os/user"
	"path"
)

type (
	// String Writer which satisfies Writer interface but writes nothing
	// and returns an empty string on String() method (similar to strings.Builder)
	StringNullWriter struct{}
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

func MakeDirAllIfNotExist(dirPath string, perm fs.FileMode) error {
	dirStat, err := os.Stat(dirPath)
	if err != nil {
		if os.IsNotExist(err) {
			if err = os.MkdirAll(dirPath, perm); err != nil {
				return fmt.Errorf("creating directory '%s': %w", dirPath, err)
			}
		} else {
			return fmt.Errorf("getting info of the directory '%s': %w", dirPath, err)
		}
	} else if !dirStat.IsDir() {
		return fmt.Errorf("path (%s) exists and is not a directory: %w", dirPath, err)
	}
	return nil
}

func (nw *StringNullWriter) Write(p []byte) (int, error) {
	return len(p), nil
}

func (nw *StringNullWriter) String() string {
	return ""
}
