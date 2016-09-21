package handlers

import (
	"net/http"
	"net/url"

	"github.com/ardanlabs/kit/db"
	"github.com/ardanlabs/kit/web/app"
	"github.com/coralproject/shelf/internal/xenia"
	"github.com/coralproject/shelf/internal/xenia/query"
)

// itemHandle maintains the set of handlers for the form api.
type itemHandle struct{}

// Item fronts the access to the comment service functionality.
var Item itemHandle

//==============================================================================

// List returns all the existing items in the system.
// 200 Success, 404 Not Found, 500 Internal
func (itemHandle) List(c *app.Context) error {

	set, err := query.GetByName(c.SessionID, c.Ctx["DB"].(*db.DB), c.Params["name"])
	if err != nil {
		if err == query.ErrNotFound {
			err = app.ErrNotFound
		}
		return err
	}

	var vars map[string]string
	if c.Request.URL.RawQuery != "" {
		if m, err := url.ParseQuery(c.Request.URL.RawQuery); err == nil {
			vars = make(map[string]string)
			for k, v := range m {
				vars[k] = v[0]
			}
		}
	}

	result := xenia.Exec(c.SessionID, c.Ctx["DB"].(*db.DB), set, vars)

	c.Respond(result, http.StatusOK)
	return nil

	//  {type:’comment’, content: ‘Stuff and things’, author:’userid123’}
}

// FilterByType returns all the existing items in the system.
// 200 Success, 404 Not Found, 500 Internal
func (itemHandle) FilterByType(c *app.Context) error {

	//  {type:’comment’, content: ‘Stuff and things’, author:’userid123’}
	return nil
}
