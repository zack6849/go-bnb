package domain

type GeocodeResult struct {
	query    string
	Location Location
	PlaceId  string
}
