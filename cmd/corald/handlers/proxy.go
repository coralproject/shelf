package handlers

import (
	"github.com/ardanlabs/kit/web"
	"github.com/coralproject/shelf/cmd/corald/service"
)

// Proxy will setup a direct proxy inbetween this service and the destination
// service using the rewrite function if specified. If the rewrite function is
// not specified, the path on the target will be set to the target path
// concatenated with the request path.
func Proxy(targetURL string, rewrite func(*web.Context) string) web.Handler {

	f := func(c *web.Context) error {

		// If specified, the rewrite will rewrite the request path.
		var targetPath string
		if rewrite != nil {
			targetPath = rewrite(c)
		}

		// Perform the actual proxy of the service request. The only error that
		// can be returned by this service is as a result of the targetURL not
		// being valid. All other errors cannot be interpreted and instead will
		// be forwarded to the requester. We will also rewrite the path as passed
		// by the rewrite function
		return c.Proxy(targetURL, service.RewritePath(c, targetPath))
	}

	return f
}
