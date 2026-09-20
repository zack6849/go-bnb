package geoapify

type Datasource struct {
	Source      string `json:"sourcename"`
	Attribution string `json:"attribution"`
	License     string `json:"license"`
	Url         string `json:"url"`
}

type Timezone struct {
	Name      string `json:"name"`
	UTCOffset string `json:"offset_STD"` //this is the "standard time" offset
}

type ReverseGeocodingResponse struct {
	Results []ReverseGeocodingResult `json:"results"`
}

type GeocodeSearchResponse struct {
	Results []GeocodeSearchResult `json:"results"`
}

// BasicResult things that should be in basically any response with a success
// e.g.: lat, lng, country, and a formatted name
type BasicResult struct {
	ResultLevel   ResultSpecificityLevel `json:"result_type"`
	datasource    Datasource
	FormattedName string  `json:"formatted"`
	Latitude      float64 `json:"lat"`
	Longitude     float64 `json:"lon"`
	Country       string  `json:"country"` //united states
}

type CityLevelResult struct {
	BasicResult
	Name         string `json:"name"` //albany
	CountryCode  string `json:"country_code"`
	City         string `json:"city"`
	State        string `json:"iso3166_2"`     //equivalent to a state, eg: US-NY
	AddressLine1 string `json:"address_line1"` //I know it seems odd to have an address for a city
	AddressLine2 string `json:"address_line2"` //but they genuinely do return one on a city level result
	Timezone     Timezone
}

type ReverseGeocodingResult struct {
	CityLevelResult
}

type GeocodeSearchResult struct {
	Datasource Datasource
	CityLevelResult
	ResultLevel ResultSpecificityLevel `json:"result_type"`
}

type ResultSpecificityLevel string

const (
	ResultLevelAmenity  ResultSpecificityLevel = "amenity"
	ResultLevelBuilding ResultSpecificityLevel = "building"
	ResultLevelStreet   ResultSpecificityLevel = "street"
	ResultLevelSuburb   ResultSpecificityLevel = "suburb"
	ResultLevelDistrict ResultSpecificityLevel = "district"
	ResultLevelPostCode ResultSpecificityLevel = "postcode"
	ResultLevelCity     ResultSpecificityLevel = "city"
	ResultLevelCounty   ResultSpecificityLevel = "county"
	ResultLevelState    ResultSpecificityLevel = "state"
	ResultLevelCountry  ResultSpecificityLevel = "country"
)

var ResultLevels = []ResultSpecificityLevel{
	ResultLevelAmenity,
	ResultLevelBuilding,
	ResultLevelStreet,
	ResultLevelSuburb,
	ResultLevelDistrict,
	ResultLevelPostCode,
	ResultLevelCity,
	ResultLevelCounty,
	ResultLevelState,
	ResultLevelCountry,
}
