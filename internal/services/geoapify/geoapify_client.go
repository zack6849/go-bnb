package geoapify

import (
	"context"
	"encoding/json"
	"fmt"
	"gobnb/internal/geo"
	"io"
	"net/http"
	"net/url"
	"time"

	"golang.org/x/time/rate"
)

const BaseUrl = "https://api.geoapify.com"

type Client struct {
	ApiKey      string
	RateLimiter *rate.Limiter
}

func CreateClient(apiKey string, limit rate.Limit) Client {
	return Client{
		ApiKey:      apiKey,
		RateLimiter: rate.NewLimiter(limit, 1),
	}
}

type AddressLevel string

func (c *Client) Autocomplete(query string, level ResultSpecificityLevel, ctx context.Context) (GeocodeSearchResult, error) {
	zero := GeocodeSearchResult{}
	target := GeocodeSearchResponse{}
	params := c.getDefaultRequestParams()
	params.Set("text", query)
	params.Set("lang", "en")
	params.Set("limit", "1")
	params.Set("type", string(level))
	params.Set("filter", "countrycode:us,ca") //na-bias because our dataset is na-centric
	params.Set("format", "json")
	req, err := c.CreateRequest("GET", "/v1/geocode/autocomplete", params)
	if err != nil {
		return zero, err
	}
	res, err := c.Execute(req, ctx)
	if err != nil {
		return zero, err
	}
	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return zero, err
	}
	err = json.Unmarshal(body, &target)
	if err != nil {
		return zero, err
	}
	if len(target.Results) == 0 {
		return zero, fmt.Errorf("no results returned for query")
	}
	return target.Results[0], nil
}

func (c *Client) Geocode(query string, level ResultSpecificityLevel, ctx context.Context) (GeocodeSearchResult, error) {
	zero := GeocodeSearchResult{}
	target := GeocodeSearchResponse{}
	params := c.getDefaultRequestParams()
	params.Set("text", query)
	params.Set("lang", "en")
	params.Set("limit", "1")
	params.Set("type", string(level))
	params.Set("filter", "countrycode:us,ca") //na-bias because our dataset is na-centric
	params.Set("format", "json")
	req, err := c.CreateRequest("GET", "/v1/geocode/search", params)
	if err != nil {
		return zero, err
	}
	res, err := c.Execute(req, ctx)
	if err != nil {
		return zero, err
	}
	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return zero, err
	}
	err = json.Unmarshal(body, &target)
	if err != nil {
		return zero, err
	}
	if len(target.Results) == 0 {
		return zero, fmt.Errorf("no results returned for query")
	}
	return target.Results[0], nil
}

func (c *Client) ReverseGeocodePoint(point *geo.Point, level ResultSpecificityLevel, ctx context.Context) (ReverseGeocodingResult, error) {
	zero := ReverseGeocodingResult{}
	target := ReverseGeocodingResponse{}
	params := c.getDefaultRequestParams()
	params.Set("type", string(level))
	params.Set("format", "json")
	params.Set("lat", fmt.Sprintf("%f", point.Latitude))
	params.Set("lon", fmt.Sprintf("%f", point.Longitude))
	req, err := c.CreateRequest("GET", "/v1/geocode/reverse", params)
	if err != nil {
		return zero, err
	}
	res, err := c.Execute(req, ctx)
	if err != nil {
		return zero, err
	}
	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return zero, err
	}
	err = json.Unmarshal(body, &target)
	if err != nil {
		return zero, err
	}
	if len(target.Results) == 0 {
		return zero, fmt.Errorf("no results returned for query")
	}
	return target.Results[0], nil
}

func (c *Client) getDefaultRequestParams() url.Values {
	values := url.Values{}
	values.Set("apiKey", c.ApiKey)
	return values
}

func (c *Client) GetHttpClient() http.Client {
	return http.Client{
		Timeout: time.Second * 2,
	}
}

func (c *Client) CreateRequest(method string, path string, params url.Values) (*http.Request, error) {
	formattedUrl := fmt.Sprintf("%s%s?%s", BaseUrl, path, params.Encode())
	return http.NewRequest(method, formattedUrl, nil)
}

func (c *Client) Execute(req *http.Request, ctx context.Context) (*http.Response, error) {
	client := c.GetHttpClient()
	//wait for the limiter before firing the request
	err := c.RateLimiter.Wait(ctx)
	if err != nil {
		return nil, err
	}
	return client.Do(req)
}
