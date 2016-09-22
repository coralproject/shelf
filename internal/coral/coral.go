// Package coral interact with the coral ecosystem endpoints.
package coral

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"

	"github.com/ardanlabs/kit/log"
	"github.com/ardanlabs/kit/web/app"
)

// Response encapsulates a http response.
type Response struct {
	Status     string
	Header     http.Header
	Body       []byte
	StatusCode int
}

// DoRequest will do a request to the web service.
func DoRequest(c *app.Context, method string, urlStr string, payload io.Reader) (*Response, error) {

	var err error
	request, err := http.NewRequest(method, urlStr, payload)
	if err != nil {
		log.Error(c, "coral.doRequest", err, "New http request.")
		return nil, err
	}
	request.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	var response *http.Response

	response, err = client.Do(request)
	if err != nil {
		log.Error(c, "coral.doRequest", err, "Sending request number.")
		return nil, err
	}
	defer response.Body.Close()

	if response.StatusCode != 200 {
		err = fmt.Errorf("Not succesful status code: %s.", response.Status)
		return nil, err
	}
	resBody, _ := ioutil.ReadAll(response.Body)

	return &Response{
		response.Status,
		response.Header,
		resBody,
		response.StatusCode,
	}, nil

}
