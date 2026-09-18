package apify

import "net/url"

type ApifyClient struct {
	ApiKey string
}

type AddressLevel string

func NewClient(apiKey string) ApifyClient {
	return ApifyClient{ApiKey: apiKey}
}

func (c *ApifyClient) Geocode(query string, level AddressLevel) {
	params := c.getDefaultRequestParams()
	params.Set("text", query)
	params.Set("lang", "en")
	params.Set("limit", "1")
	params.Set("type", "city")
	params.Set("format", "json")
}

func (c *ApifyClient) getDefaultRequestParams() url.Values {
	values := url.Values{}
	values.Set("apiKey", c.ApiKey)
	return values
}
