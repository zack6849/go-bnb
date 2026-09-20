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

type ReverseGeocodingResult struct {
	datasource    Datasource
	Name          string                 `json:"name"`    //albany
	Country       string                 `json:"country"` //united states
	CountryCode   string                 `json:"country_code"`
	City          string                 `json:"city"`
	State         string                 `json:"iso3166_2"` //equivalent to a state, eg: US-NY
	Latitude      float64                `json:"lat"`
	Longitude     float64                `json:"lon"`
	ResultLevel   ResultSpecificityLevel `json:"result_type"`
	FormattedName string                 `json:"formatted"`
	AddressLine1  string                 `json:"address_line1"`
	AddressLine2  string                 `json:"address_line2"`
	Timezone      Timezone
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
