package okta

import (
	"net/http"
	"strings"

	"github.com/johnmikee/manifester/pkg/helpers"
	"github.com/johnmikee/manifester/pkg/logger"
	"github.com/johnmikee/manifester/pkg/requester"
)

// Client represents the Okta client.
type Client struct {
	token      string
	baseURL    string
	domain     string
	blockGroup string
	client     *http.Client
	log        logger.Logger
}

// Config represents the configuration for the Okta client.
type Config struct {
	Domain string         `json:"domain,omitempty"`
	URL    string         `json:"url,omitempty"`
	Token  string         `json:"token,omitempty"`
	Client *http.Client   `json:"client,omitempty"`
	Log    *logger.Logger `json:"log,omitempty"`
}

// New returns a pointer with the Client after validating the arguments passed.
func New(c *Config) *Client {
	return &Client{
		token:   helpers.TokenValidator(c.Token, "SSWS"),
		baseURL: helpers.URLShaper(c.URL, "api/v1/"),
		domain:  c.Domain,
		client:  c.Client,
		log:     logger.ChildLogger("okta", c.Log),
	}
}

func (o *Client) newRequest(method, url string, override bool, body interface{}) (*http.Request, error) {
	return requester.New(method, o.baseURL, url, override, body)
}

func (o *Client) headers(req *http.Request) *http.Request {
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", o.token)

	return req
}

func (o *Client) do(req *http.Request, v interface{}) (*http.Response, error) {
	o.headers(req)
	return requester.Do(o.client, req, v)
}

func linkSorter(l []string) string {
	var link string
	for _, i := range l {
		res := strings.Split(i, ";")
		check := strings.Split(res[1], "=")
		if strings.Contains(check[1], "next") {
			link = res[0]
		}
	}
	link = strings.ReplaceAll(link, "<", "")
	link = strings.ReplaceAll(link, ">", "")

	return link
}
