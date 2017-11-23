package main

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

//
// weatherUnmarshalError
//
type weatherUnmarshalError struct {
	msg string
	err error
}

func newWeatherUnmarshalError(msg string, err error) weatherUnmarshalError {
	return weatherUnmarshalError{msg, err}
}

func (w weatherUnmarshalError) Error() string {
	return fmt.Sprintf("WeatherUnmarshalError: %s\n%v", w.msg, w.err)
}

//
// response
//
type response struct {
	polygonCoordinates         []*latLng
	polygonWkt                 string
	minTemp, maxTemp           float64
	time, coord                string
	maxTempColor, minTempColor string
}

func (r *response) rowHeader() []string {
	return []string{"Polygon", "Time", "Coordinate", "minTemp", "minTempColor", "maxTemp", "maxTempColor"}
}
func (r *response) row() []string {
	return []string{fmt.Sprintf(`"%s"`, r.polygonWkt),
		r.time, r.coord, fmt.Sprintf("%f", r.minTemp),
		r.minTempColor, fmt.Sprintf("%f", r.maxTemp),
		r.maxTempColor}
}

func (r *response) polygonCoords() [][]float64 {
	var res [][]float64
	for _, coord := range r.polygonCoordinates {
		res = append(res, coord.coord())
	}
	res = append(res, r.polygonCoordinates[0].coord())
	return res
}

func (r *response) properties() (props map[string]interface{}) {
	props = make(map[string]interface{})
	props["maxTemp"] = r.maxTemp
	props["maxTempColor"] = r.maxTempColor
	props["minTemp"] = r.minTemp
	props["minTempColor"] = r.minTempColor
	return
}

//
// functions
//
func readWeatherData(url string) (body []byte, err error) {
	var (
		resp *http.Response
		req  *http.Request
	)
	client := &http.Client{}
	if req, err = http.NewRequest("GET", url, nil); err != nil {
		return
	}
	// req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 6.1; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/61.0.3163.100 Safari/537.36")
	if resp, err = client.Do(req); err != nil {
		return
	}
	defer resp.Body.Close()
	body, err = ioutil.ReadAll(resp.Body)
	return
}

func sendReq(templateURL string, polygons []*polygon) (responses []*response, err error) {
	var (
		xmlResponse []byte
	)
	latLngs := make([]string, len(polygons), len(polygons))
	for i, pol := range polygons {
		latLngs[i] = pol.center.String()
	}
	reqURL := fmt.Sprintf(templateURL, strings.Join(latLngs, "%20"))
	if xmlResponse, err = readWeatherData(reqURL); err != nil {
		fmt.Println(err)
		return
	}
	forecast := newForecast()
	if err = xml.Unmarshal(xmlResponse, forecast); err != nil {
		log.Printf("Error for url: %s\n\n", reqURL)
		err = newWeatherUnmarshalError(reqURL, err)
		return
	}
	for i, feature := range forecast.FeatureMembers {
		if feature.Result.MaxTemp > 900 {
			continue
		}
		featureCoords := feature.Coordinates()
		for _, polygon := range polygons {
			if polygon.centerStr == featureCoords {
				res := &response{polygonWkt: polygons[i].wkt(),
					minTemp:            feature.Result.MinTemp,
					maxTemp:            feature.Result.MaxTemp,
					time:               feature.Result.Time,
					coord:              feature.Result.Position.Point.Coordinates,
					polygonCoordinates: polygons[i].coords}
				responses = append(responses, res)
				break
			}
		}
	}
	return
}
