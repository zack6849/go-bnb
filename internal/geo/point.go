package geo

import (
	"fmt"
	"sort"

	"github.com/mmcloughlin/geohash"
)

type Point struct {
	Latitude    float64
	Longitude   float64
	BoundingBox *geohash.Box
}

func NewPoint(lat float64, lng float64) *Point {
	return &Point{
		Latitude:  lat,
		Longitude: lng,
	}
}

func (p *Point) ToString() string {
	return fmt.Sprintf("(%f, %f)", p.Latitude, p.Longitude)
}

func (p *Point) Hash() string {
	return p.HashWithPrecision(HashLevelCity) //default to city-level hashes
}

func (p *Point) HashWithPrecision(level HashLevel) string {
	return geohash.EncodeWithPrecision(p.Latitude, p.Longitude, uint(level))
}

func (p *Point) GetBoundingBox(level HashLevel) *geohash.Box {
	if p.BoundingBox != nil {
		return p.BoundingBox
	}
	hashed := geohash.EncodeWithPrecision(p.Latitude, p.Longitude, uint(level))
	box := geohash.BoundingBox(hashed)
	p.BoundingBox = &box
	return &box
}

func (p *Point) GetBoxEdges(level HashLevel) []*Point {
	box := p.GetBoundingBox(level)
	topLeft := NewPoint(box.MaxLat, box.MinLng)
	topRight := NewPoint(box.MaxLat, box.MaxLng)
	bottomLeft := NewPoint(box.MinLat, box.MinLng)
	bottomRight := NewPoint(box.MinLat, box.MaxLng)
	return []*Point{topLeft, topRight, bottomLeft, bottomRight}
}

func (p *Point) GetBoxEdgeMidpoints(level HashLevel) []*Point {
	box := p.GetBoundingBox(level)
	centerLat, centerLng := box.Center()
	northMidpoint := NewPoint(box.MaxLat, centerLng)
	southMidpoint := NewPoint(box.MinLat, centerLng)
	eastMidpoint := NewPoint(centerLat, box.MinLng)
	westMidpoint := NewPoint(
		centerLat,
		box.MaxLng,
	)
	return []*Point{northMidpoint, southMidpoint, eastMidpoint, westMidpoint}
}

type SubPoint struct {
	Point *Point
	Score float64
}

func (p *Point) SubdivideIntoGrid(gridSize int, level HashLevel) []SubPoint {
	box := p.GetBoundingBox(level)
	candidates := make([]SubPoint, 0, gridSize*gridSize)
	latStep := (box.MaxLat - box.MinLat) / float64(gridSize+1)
	lngStep := (box.MaxLng - box.MinLng) / float64(gridSize+1)

	for i := range gridSize {
		lat := box.MinLat + float64(i)*latStep
		for j := range gridSize {
			lng := box.MinLng + float64(j)*lngStep
			dLat := lat - p.Latitude
			dLng := lng - p.Longitude
			candidates = append(candidates, SubPoint{
				Point: NewPoint(lat, lng),
				Score: dLat*dLat + dLng*dLng,
			})
		}
	}
	sort.Slice(candidates, func(i, j int) bool {
		return candidates[i].Score < candidates[j].Score
	})
	return candidates
}
