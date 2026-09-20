package listing_search

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/gin-gonic/gin"
)

func GetGinTestContext() *gin.Context {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	return c
}

func SendRequest(values url.Values) (ListingSearchParameters, error) {
	c := GetGinTestContext()
	route := ListingApiUrl
	params := values.Encode()
	if len(params) > 0 {
		route += "?" + values.Encode()
	}
	c.Request, _ = http.NewRequest("GET", route, nil)
	return GetRequestsFromParams(c)
}

const ListingApiUrl = "/api/listings"
const SearchLng = -73.7782
const SearchLat = 42.6657
const SearchDistanceMeters = 30_000 //30km

func GetDefaultRequestParams() url.Values {
	v := url.Values{}
	v.Set("latitude", fmt.Sprintf("%f", SearchLat))
	v.Set("longitude", fmt.Sprintf("%f", SearchLng))
	v.Set("distance", fmt.Sprintf("%d", SearchDistanceMeters))
	return v
}

func TestConvertsQueryParametersToSearchParameters(t *testing.T) {
	v := GetDefaultRequestParams()
	params, err := SendRequest(v)
	assert.Equal(t, err, nil)
	assert.Equal(t, params.MaxDistance, SearchDistanceMeters)
	assert.InDelta(t, params.SearchNear.Latitude, SearchLat, 1)
	assert.InDelta(t, params.SearchNear.Longitude, SearchLng, 1)
}

// TODO: more testing
func TestReturnsErrorWhenConvertingFields(t *testing.T) {
	v := GetDefaultRequestParams()
	v.Set("distance", "foo")
	_, err := SendRequest(v)
	assert.ErrorContains(t, err, "failed to parse distance")
	assert.ErrorContains(t, err, "foo")
}

func TestRequiresLocation(t *testing.T) {
	v := GetDefaultRequestParams()
	v.Del("latitude")
	v.Del("longitude")
	_, err := SendRequest(v)
	assert.ErrorContains(t, err, "latitude and longitude are required to use this endpoint")
}
