package main

import (
	"encoding/json"
	"io/ioutil"
)

type confProperty struct {
	Name       string  `json:"Name"`
	Max        float64 `json:"max"`
	Min        float64 `json:"min"`
	ColorsFile string  `json:"colorsFile"`
}

type config struct {
	CoordsFileName string         `json:"CoordFileName"`
	URL            string         `json:"URL"`
	CsvDelimiter   string         `json:"CsvDelimiter"`
	Workers        int            `json:"Workers"`
	Step           int            `json:"Step"`
	GridSize       float64        `json:"gridSize"`
	GridWSLatLng   []float64      `json:"gridWSLatLng"`
	GridENLatLng   []float64      `json:"gridENLatLng"`
	OutFormats     []string       `json:"outFormats"`
	OutFileName    string         `json:"outFileName"`
	Properties     []confProperty `json:"properties"`
}

func readConf(confName string) (conf *config, err error) {
	var (
		buff []byte
	)
	if buff, err = ioutil.ReadFile(confName); err != nil {
		return
	}
	conf = &config{}
	err = json.Unmarshal(buff, conf)
	return
}
