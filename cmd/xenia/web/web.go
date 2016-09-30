package web

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/ardanlabs/kit/cfg"
	"github.com/coralproject/shelf/internal/platform/auth"
	"github.com/spf13/cobra"
)

const (
	// cfgWebHost is the url to the web host which includes the hostname, port,
	// and scheme.
	cfgWebHost = "WEB_HOST"

	// cfgPlatformPrivateKey is the private key used to sign new requests to the
	// downstream service layer.
	cfgPlatformPrivateKey = "PLATFORM_PRIVATE_KEY"
)

var authSigner auth.Signer

func init() {
	if err := cfg.Init(cfg.EnvProvider{Namespace: "XENIA"}); err != nil {
		log.Println("Unable to initialize configuration")
		os.Exit(1)
	}

	platformPrivateKey, err := cfg.String(cfgPlatformPrivateKey)
	if err != nil {
		log.Printf("Downstream Auth : Disabled : %s\n", err.Error())
		return
	}

	// If the platformPrivateKey is provided, then we should generate the token
	// signing function to be used when composing requests down to the platform.
	if platformPrivateKey != "" {
		signer, err := auth.NewSigner(platformPrivateKey)
		if err != nil {
			log.Printf("Downstream Auth : Error : %s", err.Error())
			os.Exit(1)
		}

		// Requests can now be signed with the given signer function which we will
		// save on the application wide context. In the event that a function
		// requires a call down to a downstream platform, we will include a signed
		// header using the signer function here.
		authSigner = signer
		log.Println("Downstream Auth : Enabled")
	} else {
		log.Printf("Downstream Auth : Disabled : %s\n", err.Error())
	}
}

// Request provides support for executing commands against the
// web service.
func Request(cmd *cobra.Command, verb, path string, post io.Reader) (string, error) {
	url, err := cfg.URL(cfgWebHost)
	if err != nil {
		return "", err
	}

	// We're using the url from the environment to load in the details of the web
	// host, but we will include the path passed in.
	url.Path = path

	cmd.Printf("%s : %s\n", verb, url.String())
	req, err := http.NewRequest(verb, url.String(), post)
	if err != nil {
		return "", err
	}

	// If the authSigner is defined, it means that we can now sign the request
	// with the authentication token.
	if authSigner != nil {

		// Perform the signing here without any additional claims attached.
		if err := auth.SignRequest("", authSigner, nil, req); err != nil {
			return "", err
		}
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	// If the error returned by the endpoint is a non ok return, then we should
	// return this as an error.
	if resp.StatusCode >= http.StatusBadRequest {
		return "", fmt.Errorf("Status : %d", resp.StatusCode)
	}

	contents, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(contents), nil
}
