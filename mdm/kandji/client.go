package kandji

import (
	"net/http"

	"github.com/johnmikee/manifester/mdm"
	"github.com/johnmikee/manifester/pkg/helpers"
	"github.com/johnmikee/manifester/pkg/logger"
	"github.com/johnmikee/manifester/pkg/requester"
)

type Client struct {
	token   string
	baseURL string
	client  *http.Client
	log     logger.Logger
}

// Setup implements mdm.Provider.
func (c *Client) Setup(config mdm.Config) {
	c.token = helpers.TokenValidator(config.Token, "Bearer")
	c.baseURL = helpers.URLShaper(config.URL, "api/v1/")
	c.client = config.Client
	c.log = logger.ChildLogger("kandji", &config.Log)
}

type offsetRange struct {
	Limit  int
	Offset int
}

func (c *Client) newRequest(method, url string, override bool, body interface{}) (*http.Request, error) {
	return requester.New(method, c.baseURL, url, override, body)
}

func (c *Client) headers(req *http.Request) {
	req.Header.Set("Authorization", c.token)
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-type", "application/json;charset=utf-8")
}

func (c *Client) do(req *http.Request, v interface{}) (*http.Response, error) {
	c.headers(req)
	return requester.Do(c.client, req, v)
}
