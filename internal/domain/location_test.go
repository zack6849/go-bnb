package domain

import (
	"strings"
	"testing"
)

// NYC lat & lng captured from PostGIS
const NycLat = 40.7128
const NycLng = -74.006

// SELECT ('SRID=4326;POINT(-74.006 40.7128)'::geography)::text;
const postGISPoint = "0101000020E6100000AAF1D24D628052C05E4BC8073D5B4440"

func TestLocationValueMatchesPostGIS(t *testing.T) {
	v, err := Location{Latitude: NycLat, Longitude: NycLng}.Value()
	if err != nil {
		t.Fatalf("Value: %v", err)
	}
	got, ok := v.(string)
	if !ok {
		t.Fatalf("Value returned %T, want string", v)
	}
	if !strings.EqualFold(got, postGISPoint) {
		t.Errorf("Value = %s, want %s", got, postGISPoint)
	}
}

// scanning from various forms, strings, bytes
func TestLocationScan(t *testing.T) {
	for name, src := range map[string]any{
		"string": postGISPoint,
		"bytes":  []byte(postGISPoint),
	} {
		t.Run(name, func(t *testing.T) {
			var l Location
			//try to scan a location from the given source
			if err := l.Scan(src); err != nil {
				t.Fatalf("Scan: %v", err)
			}
			//if the lat and lng don't match, we have an error!
			if l.Latitude != NycLat || l.Longitude != NycLng {
				t.Errorf("Scan = %+v, want {%f, %f}", l, NycLat, NycLng)
			}
		})
	}
}

func TestLocationScanNil(t *testing.T) {
	//pass a location with values
	l := Location{Latitude: 1, Longitude: 2}
	//scan a nil column
	if err := l.Scan(nil); err != nil {
		t.Fatalf("Scan: %v", err)
	}
	if (l != Location{}) {
		//should return a zero value for a nil location
		t.Errorf("Scan(nil) = %+v, want zero value", l)
	}
}

func TestLocationScanRejectsGarbage(t *testing.T) {
	var l Location
	if err := l.Scan("not hex at all"); err == nil {
		t.Error("Scan accepted a non-hex string")
	}
	if err := l.Scan(42); err == nil {
		t.Error("Scan accepted an int")
	}
}
