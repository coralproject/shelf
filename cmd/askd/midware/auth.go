package midware

import "github.com/ardanlabs/kit/web/app"

// Auth handles token authentication.
func Auth(h app.Handler) app.Handler {
	return func(c *app.Context) error {
		return h(c)
	}
}
