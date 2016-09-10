package web

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"

	"github.com/ardanlabs/kit/cfg"
	"github.com/spf13/cobra"
)

const (
	cfgHost = "WEB_HOST"
	cfgAuth = "WEB_AUTH"
)

// Request provides support for executing commands against the
// web service.
func Request(cmd *cobra.Command, verb string, url string, post io.Reader) (string, error) {
	host, err := cfg.String(cfgHost)
	if err != nil {
		return "", err
	}

	url = "http://" + host + url

	cmd.Printf("%s : %s\n", verb, url)
	req, err := http.NewRequest(verb, url, post)
	if err != nil {
		return "", err
	}

	auth, err := cfg.String(cfgAuth)
	if err == nil {
		cmd.Println("Using Authentication")
		req.Header.Add("Authorization", auth)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNoContent {
		return "", fmt.Errorf("Status : %d", resp.StatusCode)
	}

	defer resp.Body.Close()

	contents, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(contents), nil
}
