package update

import (
	"net/http"

	up "github.com/inconshreveable/go-update"
)

func doUpdate(url string) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	err = up.Apply(resp.Body, up.Options{})
	if err != nil {
		// error handling
	}
	return err
}
