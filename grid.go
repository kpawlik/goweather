package main

import (
	"bytes"
	"io/ioutil"
	"strconv"
)

func readCoords(fileName string) (coords []*latLng, err error) {
	var (
		content []byte
	)
	if content, err = ioutil.ReadFile(fileName); err != nil {
		return
	}
	coordsArr := bytes.Split(content, []byte(" "))
	for _, byteCoords := range coordsArr {
		arr := bytes.Split(byteCoords, []byte(","))
		lat, _ := strconv.ParseFloat(string(arr[0]), 32)
		lng, _ := strconv.ParseFloat(string(arr[1]), 32)
		coords = append(coords, newLatLng(lat, lng))
	}

	return
}

//
// calcGrid calculate grid of coordinates between GridWSLatLng and GridENLatLng with size GridSize (in km)
//
func calcGrid(cfg *config) (latLngs []*latLng) {
	ll := newLatLng(cfg.GridWSLatLng[0], cfg.GridWSLatLng[1])
	// scale to meters
	dist := cfg.GridSize * 1000.0
	bbox := ll.toBounds(dist)
	startBBox := bbox
	poly := newpolygon([]*latLng{bbox.getSouthWest(), bbox.getSouthEast(), bbox.getNorthEast(), bbox.getNorthWest()})
	latLngs = append(latLngs, poly.center)
	for {
		if ll.lat > cfg.GridENLatLng[0] { //if ll.lat > 50 {
			ll = startBBox.getSouthEast()
			tmp := ll.toBounds(dist)
			ll = tmp.getNorthEast()
			startBBox = ll.toBounds(dist)
		} else {
			ll = bbox.getNorthEast()
		}
		if ll.lng > cfg.GridENLatLng[1] {
			break
		}
		bbox = ll.toBounds(dist)
		ll = bbox.getNorthWest()
		bbox = ll.toBounds(dist)
		poly = newpolygon([]*latLng{bbox.getSouthWest(), bbox.getSouthEast(), bbox.getNorthEast(), bbox.getNorthWest()})
		latLngs = append(latLngs, poly.center)
	}
	return

}
