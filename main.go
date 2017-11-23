package main

import (
	"bufio"
	"encoding/csv"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/kpawlik/geojson"
)

type writeOutFunc func([]*response, *config) error

var (
	confFile    string
	writeOutMap = map[string]writeOutFunc{"json": writeGeoJSON, "csv": writeCSV}
)

func init() {
	flag.StringVar(&confFile, "c", "conf.json", "Config file")
	flag.Parse()

}

func main() {
	startTime := time.Now()
	cfg, err := readConf(confFile)
	if err != nil {
		log.Fatal(err)
	}
	coords := calcGrid(cfg)

	if cfg.Workers == 0 {
		cfg.Workers = 1
	}
	concurrent(cfg, coords)
	log.Println("Total: ", time.Now().Sub(startTime))
}

func splitPolygons(polygons []*polygon, step int) (chunks [][]*polygon) {
	start := 0
	next := 0

	stop := len(polygons) - 1
	for {
		next += step
		if next >= stop {
			break
		}
		if (start + step) > len(polygons) {
			next = len(polygons) - 1
		}
		chunk := polygons[start:next]
		chunks = append(chunks, chunk)
		start = next
	}
	return
}

func concurrent(cfg *config, coords []*latLng) {
	polygons, err := calculatePolygons(coords)
	if err != nil {
		log.Panic(err)
	}
	chunks := splitPolygons(polygons, cfg.Step)
	results := make([]*response, 0, 0)
	in := make(chan []*polygon, len(chunks))
	out := make(chan []*response, cfg.Workers)
	for i := 0; i < cfg.Workers; i++ {
		go worker(cfg, in, out)
	}
	for _, chunk := range chunks {
		in <- chunk
	}
	for i := 0; i < len(chunks); i++ {
		chunk := <-out
		log.Println(i, len(chunks))
		results = append(results, chunk...)
	}
	for _, prop := range cfg.Properties {
		switch prop.Name {
		case "maxt", "mint":
			colors, err := readColors(prop.ColorsFile)
			if err != nil {
				log.Println(err)
			}
			matchColors(results, colors, prop)
		}
	}
	writeOut(results, cfg)

}

func worker(cfg *config, in chan []*polygon, out chan []*response) {
	var reqProps []string
	for _, props := range cfg.Properties {
		reqProps = append(reqProps, props.Name)
	}
	url := cfg.URL + strings.Join(reqProps, ",")
	for {
		chunk := <-in
		chunkResults, err := sendReq(url, chunk)
		if err != nil {
			if err, ok := err.(weatherUnmarshalError); ok {
				log.Println(err.msg)
				in <- chunk
			} else {
				log.Println(err)
			}
		} else {
			out <- chunkResults
		}
	}
}

func writeCSV(results []*response, cfg *config) (err error) {
	var f *os.File
	if f, err = os.Create(fmt.Sprintf("%s.%s", cfg.OutFileName, "csv")); err != nil {
		return
	}
	defer f.Close()
	csvWriter := csv.NewWriter(f)
	csvWriter.Comma = rune(cfg.CsvDelimiter[0])
	if len(results) == 0 {
		return
	}
	csvWriter.Write(results[0].rowHeader())
	for _, res := range results {
		if err = csvWriter.Write(res.row()); err != nil {
			return
		}
	}
	csvWriter.Flush()
	return
}

func writeGeoJSON(results []*response, cfg *config) (err error) {
	var f *os.File
	if f, err = os.Create(fmt.Sprintf("%s.%s", cfg.OutFileName, "json")); err != nil {
		return
	}
	fc := geojson.NewFeatureCollection(nil)
	for _, result := range results {
		coords := geojson.Coordinates{}
		for _, polygonCoord := range result.polygonCoords() {
			gsCoord := geojson.Coordinate{geojson.CoordType(polygonCoord[1]), geojson.CoordType(polygonCoord[0])}
			coords = append(coords, gsCoord)
		}
		p := geojson.NewPolygon(nil)
		p.AddCoordinates(coords)

		f := geojson.NewFeature(p, result.properties(), nil)
		fc.AddFeatures(f)
	}
	w := bufio.NewWriter(f)
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	if err = enc.Encode(fc); err != nil {
		return
	}
	return
}

func writeOut(results []*response, cfg *config) {
	for _, format := range cfg.OutFormats {
		if writeFunc, ok := writeOutMap[format]; ok {
			writeFunc(results, cfg)
		}
	}
}
