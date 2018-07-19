package xkcd

import (
	"io"
	"net/http"
	"os"
	"os/user"
	"path/filepath"

	"github.com/davecgh/go-spew/spew"
)

// expand returns the absolute path of a file and optionally does the
// expansion if there is a `~'.
func expand(path string) (string, error) {
	if len(path) == 0 || path[0] != '~' {
		return path, nil
	}

	usr, err := user.Current()
	if err != nil {
		return path, err
	}

	return filepath.Join(usr.HomeDir, path[1:]), nil
}

// FetchComicImg is a utility function to download a comic image and store
// it locally on disk.
func FetchComicImg(url, path string) error {
	path, err := expand(path)
	if err != nil {
		return err
	}

	path, err = filepath.Abs(path)
	if err != nil {
		return err
	}

	out, err := os.Create(path)
	if err != nil {
		return err
	}
	defer out.Close()

	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return err
	}

	return nil
}

// ToStr pretty prints a Comic struct.
func (c Comic) ToStr() string {
	var spewFmt = spew.ConfigState{Indent: "\t"}
	return spewFmt.Sdump(c)
}
