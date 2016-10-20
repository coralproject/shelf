package service

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/ardanlabs/kit/web"
	"github.com/coralproject/shelf/internal/sponge/item"
)

var (
	// ErrItemNotFound is when the item is not found.
	ErrItemNotFound = errors.New("Item not found")

	// ErrNotAnItem is returned when the interface{} is not an Item{}
	ErrNotAnItem = errors.New("Not an item")

	// ErrServiceNotSet is returned when the URL for the service is not setup
	ErrServiceNotSet = errors.New("Service Sponged not found")
)

// GetItemByID returns an item based on its item_id
func GetItemByID(c *web.Context, itemID string) (*item.Item, error) {
	spongedURL, ok := c.Web.Ctx["spongedURL"].(string)
	if !ok {
		return nil, ErrServiceNotSet
	}

	// Get the item by ID
	url := spongedURL + "/v1/item/" + itemID

	resp, err := requestSponge(c, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var items []item.Item
	if err := json.NewDecoder(resp.Body).Decode(&items); err != nil {
		return nil, err
	}

	if len(items) == 0 {
		return nil, ErrItemNotFound
	}

	// We are only retrieving one item.
	itm := &items[0]

	return itm, nil
}

// UpsertItem upserts an item
func UpsertItem(c *web.Context, itm *item.Item) error {
	spongedURL, ok := c.Web.Ctx["spongedURL"].(string)
	if !ok {
		return ErrServiceNotSet
	}

	// Upsert the target with the new actions.
	url := spongedURL + "/v1/item"

	// Send the item into Sponge.
	payload, err := json.Marshal(&itm)
	if err != nil {
		return err
	}

	resp, err := requestSponge(c, http.MethodPost, url, payload)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return nil
}

/*=======================================================================================*/

func requestSponge(c *web.Context, verb string, url string, payload []byte) (*http.Response, error) {
	req, err := http.NewRequest(verb, url, bytes.NewBuffer(payload))
	if err != nil {
		return nil, err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	return resp, nil
}
