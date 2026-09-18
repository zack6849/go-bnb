package geocode

import (
	"fmt"
	"gobnb/internal/configuration"
	"gobnb/internal/domain"
	"net/http"
	"net/url"
	"strings"
	"time"
)

func NormalizeLocationString(query string) string {
	return strings.ToLower(strings.Trim(query, " !'."))
}

func SearchCities(query string) {

}

func getRequestParameters(query string, apiKey string) url.Values {
	values := url.Values{}
	values.Set("apiKey", apiKey)
	values.Set("text", query)
	return values
}

func fetchGeocodingResult(query string) (domain.GeocodeResult, error) {
	apiKey := configuration.GetString("GEOAPIFY_API_KEY", "")
	if len(apiKey) == 0 {
		return domain.GeocodeResult{}, fmt.Errorf("geocoding disabled, no geoapify api key provided, (no GEOAPIFY_API_KEY)")
	}
	values := getRequestParameters(query, apiKey)
	url := "https://api.geoapify.com/v1/geocode/search?" + values.Encode()
	client := http.Client{
		Timeout: time.Second * 2,
	}
	resp, err := client.Get(url)

	if err != nil {
		return domain.GeocodeResult{}, err
	}
}
