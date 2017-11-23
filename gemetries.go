package main

import (
	"fmt"
	"math"
	"strings"
)

//
// latLng
//
type latLng struct {
	lat, lng float64
}

func newLatLng(lat, lng float64) *latLng {
	return &latLng{lat, lng}
}

func (c *latLng) String() string {
	return fmt.Sprintf("%.5f,%.5f", c.lat, c.lng)
}
func (c *latLng) StringRev() string {
	return fmt.Sprintf("%.5f,%.5f", c.lng, c.lat)
}

func (c *latLng) wkt() string {
	return fmt.Sprintf("%f %f", c.lng, c.lat)
}
func (c *latLng) coord() []float64 {
	return []float64{c.lat, c.lng}
}

func (c *latLng) toBounds(sizeInMeters float64) *latLngBounds {
	var (
		latAccuracy, lngAccuracy float64
	)
	latAccuracy = (180 * sizeInMeters) / 40075017
	lngAccuracy = latAccuracy / math.Cos((math.Pi/180)*c.lat)
	return newLatLngBounds(newLatLng(c.lat-latAccuracy, c.lng-lngAccuracy), newLatLng(c.lat+latAccuracy, c.lng+lngAccuracy))

}

//
// latLngBounds
//
type latLngBounds struct {
	southWest, northEast *latLng
}

func newLatLngBounds(southWest, northEast *latLng) *latLngBounds {
	return &latLngBounds{southWest, northEast}
}

func (b *latLngBounds) getCenter() *latLng {
	return newLatLng((b.southWest.lat+b.northEast.lat)/2, (b.southWest.lng+b.northEast.lng)/2)
}

func (b *latLngBounds) getSouthWest() *latLng {
	return b.southWest
}
func (b *latLngBounds) getNorthEast() *latLng {
	return b.northEast
}
func (b *latLngBounds) getNorthWest() *latLng {
	return newLatLng(b.getNorth(), b.getWest())
}

func (b *latLngBounds) getSouthEast() *latLng {
	return newLatLng(b.getSouth(), b.getEast())
}
func (b *latLngBounds) getWest() float64 {
	return b.southWest.lng
}
func (b *latLngBounds) getSouth() float64 {
	return b.southWest.lat
}
func (b *latLngBounds) getEast() float64 {
	return b.northEast.lng
}
func (b *latLngBounds) getNorth() float64 {
	return b.northEast.lat
}

//
// polygon
//
type polygon struct {
	center    *latLng
	coords    []*latLng
	centerStr string
}

func newpolygon(coords []*latLng) *polygon {
	centerLat := (coords[0].lat + coords[2].lat) / 2
	centerLng := (coords[0].lng + coords[2].lng) / 2
	center := newLatLng(centerLat, centerLng)
	return &polygon{coords: coords,
		center:    center,
		centerStr: center.StringRev(),
	}
}

func (p *polygon) wkt() string {
	wkt := make([]string, 5, 5)
	for i, c := range p.coords {
		wkt[i] = c.wkt()
	}
	wkt[4] = p.coords[0].wkt()
	return fmt.Sprintf("POLYGON((%s))", strings.Join(wkt, ", "))
}

//
// Functions
//
func calculatePolygons(coords []*latLng) (polygons []*polygon, err error) {
	var prevCoord *latLng
	subgrid := make([]*latLng, 0, 0)
	grid := make([][]*latLng, 0, 0)

	for _, coordinate := range coords {
		if prevCoord != nil && (prevCoord.lat-coordinate.lat) > 10 {
			grid = append(grid, subgrid)
			subgrid = make([]*latLng, 0, 0)
		}
		subgrid = append(subgrid, coordinate)
		prevCoord = coordinate
	}
	for i := 0; i < len(grid)-1; i++ {
		subgrid1 := grid[i]
		subgrid2 := grid[i+1]
		subgrid1Len := len(subgrid1) - 1
		subgrid2Len := len(subgrid2) - 1
		for j := 0; j < len(subgrid1)-1; j++ {
			if j+1 > subgrid1Len || j+1 > subgrid2Len {
				continue
			}
			poly := []*latLng{subgrid1[j], subgrid2[j], subgrid2[j+1], subgrid1[j+1]}
			polygons = append(polygons, newpolygon(poly))
		}
	}

	return
}
