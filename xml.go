package main

import (
	"encoding/xml"
	"fmt"
)

type point struct {
	XMLName     xml.Name `xml:"Point"`
	Coordinates string   `xml:"coordinates"`
}

type position struct {
	XMLName xml.Name `xml:"position"`
	Point   point
}

type result struct {
	XMLName  xml.Name `xml:"Forecast_Gml2Point"`
	Position position `xml:"position"`
	Time     string   `xml:"validTime"`
	MaxTemp  float64  `xml:"maximumTemperature"`
	MinTemp  float64  `xml:"minimumTemperature"`
}
type featureMember struct {
	XMLName xml.Name `xml:"featureMember"`
	Result  result   `xml:"Forecast_Gml2Point"`
}

func (f *featureMember) Coordinates() string {
	return f.Result.Position.Point.Coordinates
}

type box struct {
	XMLName     xml.Name `xml:"Box"`
	Coordinates string   `xml:"coordinates"`
}
type boundedBy struct {
	XMLName xml.Name `xml:"boundedBy"`
	Box     box      `xml:"Box"`
}

type forecast struct {
	Name           xml.Name        `xml:"NdfdForecastCollection"`
	Where          string          `xml:"app,attr"`
	BoundedBy      boundedBy       `xml:"boundedBy"`
	FeatureMembers []featureMember `xml:"featureMember"`
}

func newForecast() *forecast {
	return &forecast{}
}

func readXML(buff []byte) {
	body := forecast{}
	err := xml.Unmarshal(buff, &body)
	if err != nil {
		fmt.Println(err)
	}
	for _, fm := range body.FeatureMembers {
		fmt.Println(fm.Result.Position.Point.Coordinates)
		fmt.Println(fm.Result.Time)
		fmt.Println(fm.Result.MaxTemp)
		fmt.Println(fm.Result.MinTemp)
		fmt.Println("-")
	}
}
