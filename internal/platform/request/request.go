package request

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	stdpath "path"

	"github.com/coralproject/shelf/internal/platform/auth"
)

// Client contains the necessary pieces to perform requests down to a service
// layer that has platform authentication enabled.
type Client struct {
	BaseURL string
	Signer  auth.Signer
}

// New creates a new request sourced from the client. If a signer is present on
// the client, requests will automatically be signed.
func (c *Client) New(context interface{}, verb, path string, body io.Reader) (*http.Request, error) {

	// Parse the base url passed in.
	uri, err := url.Parse(c.BaseURL)
	if err != nil {
		return nil, err
	}

	// Join the paths from the base url and the request path.
	uri.Path = stdpath.Join(uri.Path, path)

	req, err := http.NewRequest(verb, uri.String(), body)
	if err != nil {
		return nil, err
	}

	// If the authSigner is defined, it means that we can now sign the request
	// with the authentication token.
	if c.Signer != nil {

		// Sign the actual request without any extra claims added.
		if err := auth.SignRequest(context, c.Signer, nil, req); err != nil {
			return nil, err
		}
	}

	return req, nil
}

// Do executes the http request on the default http client and returns the bytes
// in the event that the response code was < 400.
func (c *Client) Do(req *http.Request) ([]byte, error) {

	// Perform the request with the default client.
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// If the error returned by the endpoint is a non ok return, then we should
	// return this as an error.
	if resp.StatusCode >= http.StatusBadRequest {
		return nil, fmt.Errorf("Status[%d]", resp.StatusCode)
	}

	// Read the response into
	return ioutil.ReadAll(resp.Body)
}
