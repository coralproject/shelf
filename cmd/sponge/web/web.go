package web

import (
	"io"
	"log"
	"os"

	"github.com/ardanlabs/kit/cfg"
	"github.com/coralproject/shelf/internal/platform/auth"
	"github.com/coralproject/shelf/internal/platform/request"
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

// DefaultClient is the default client to perform downstream requests to.
var DefaultClient request.Client

func init() {
	if err := cfg.Init(cfg.EnvProvider{Namespace: "SPONGE"}); err != nil {
		log.Println("Unable to initialize configuration")
		os.Exit(1)
	}

	// Insert the base url for requests by this client.
	DefaultClient.BaseURL = cfg.MustURL(cfgWebHost).String()

	platformPrivateKey, err := cfg.String(cfgPlatformPrivateKey)
	if err != nil || platformPrivateKey == "" {
		if err != nil {
			log.Printf("Downstream Auth : Disabled : %s\n", err.Error())
			return
		}

		log.Printf("Downstream Auth : Disabled\n")
		return
	}

	// If the platformPrivateKey is provided, then we should generate the token
	// signing function to be used when composing requests down to the platform.
	signer, err := auth.NewSigner(platformPrivateKey)
	if err != nil {
		log.Printf("Downstream Auth : Error : %s", err.Error())
		os.Exit(1)
	}

	// Requests can now be signed with the given signer function which we will
	// save on the application wide context. In the event that a function
	// requires a call down to a downstream platform, we will include a signed
	// header using the signer function here.
	DefaultClient.Signer = signer

	log.Println("Downstream Auth : Enabled")
}

// Request provides support for executing commands against the
// web service.
func Request(cmd *cobra.Command, verb, path string, body io.Reader) (string, error) {
	req, err := DefaultClient.New("", verb, path, body)
	if err != nil {
		return "", nil
	}

	resp, err := DefaultClient.Do(req)
	if err != nil {
		return "", nil
	}

	return string(resp), err
}
