package talk

import (
	"bytes"
	"encoding/json"
	"net/http"

	"github.com/coralproject/shelf/internal/sponge/item"
)

func getItemByID(spongedURL string, targetID string) (item.Item, error) {
	var itm item.Item

	// Get the item by ID
	url := spongedURL + "/v1/item/" + targetID

	client := &http.Client{}
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return itm, err
	}

	resp, err := client.Do(req)
	if err != nil {
		return itm, err
	}
	defer resp.Body.Close()

	var items []item.Item

	if err := json.NewDecoder(resp.Body).Decode(&items); err != nil {
		return itm, err
	}

	// We are only retrieving one item.
	itm = items[0]

	return itm, nil
}

func upsertItem(spongedURL string, target item.Item) error {

	// Upsert the target with the new actions.
	url := spongedURL + "/v1/item"

	// Send the target into Sponge.
	body, err := json.Marshal(target)
	if err != nil {
		return err
	}

	client := &http.Client{}
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(body))
	if err != nil {
		return err
	}

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return nil
}
