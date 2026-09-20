package domain

import (
	"database/sql/driver"
	"fmt"
	"gobnb/internal/geo"
	"strconv"

	"github.com/twpayne/go-geom"
	"github.com/twpayne/go-geom/encoding/ewkbhex"
)

// World Geodetic Sys 1984, used by GPS
const wgs84SRID = 4326

type Location struct {
	Latitude  float64
	Longitude float64
}

func (l Location) ParseLocation(latStr string, lngStr string) (Location, error) {
	zero := Location{}
	lat, err := strconv.ParseFloat(latStr, 64)
	if err != nil {
		return zero, err
	}
	lng, err := strconv.ParseFloat(lngStr, 64)
	return Location{
		Latitude:  lat,
		Longitude: lng,
	}, nil
}

// GormDataType tells GORM this field maps to a single column, not a relation.
func (Location) GormDataType() string {
	return "geography"
}

// Value encodes the value as an EWKB hex value on store for PostGIS
func (l Location) Value() (driver.Value, error) {
	// PostGIS orders coordinates x, y — longitude first.
	point := geom.NewPointFlat(geom.XY, []float64{l.Longitude, l.Latitude}).SetSRID(wgs84SRID)
	return ewkbhex.Encode(point, ewkbhex.NDR)
}

func (l Location) ToPoint() *geo.Point {
	return &geo.Point{
		Latitude:  l.Latitude,
		Longitude: l.Longitude,
	}
}

// Scan decode hex-encoded EWKB that PostGIS returns for geography columns.
func (l *Location) Scan(src any) error {
	if src == nil {
		*l = Location{}
		return nil
	}

	var raw string
	switch v := src.(type) {
	//if it's a string, parse it
	case string:
		raw = v
	//if it's a byte[] parse that as a string
	case []byte:
		raw = string(v)
	default:
		return fmt.Errorf("location: cannot scan %T", src)
	}
	//decode our raw value with ewkbhex
	decoded, err := ewkbhex.Decode(raw)
	if err != nil {
		//error decoding
		return fmt.Errorf("location: failure decoding %q: %w", raw, err)
	}
	//convert the decoded data into a geom point
	point, ok := decoded.(*geom.Point)
	if !ok {
		//not ok, failed to cast to a geom point, was something else
		return fmt.Errorf("location: expected a point, got %T", decoded)
	}
	//if we get an empty point back. return nil values
	if point.Empty() {
		*l = Location{}
		return nil
	}
	//otherwise set the location's lat & lng
	l.Longitude = point.X()
	l.Latitude = point.Y()
	return nil
}
