package main

import (
	"fmt"
	"image"
	_ "image/png"
	"io"
	"log"
	"math"
	"os"
)

func readColors(fileName string) (colors []string, err error) {
	var (
		reader io.ReadCloser
		img    image.Image
	)
	if reader, err = os.Open(fileName); err != nil {
		log.Fatal(err)
	}
	defer reader.Close()
	if img, _, err = image.Decode(reader); err != nil {
		return
	}
	bounds := img.Bounds()
	colors = make([]string, bounds.Max.X, bounds.Max.X)
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			r, g, b, _ := img.At(x, y).RGBA()
			colors[x] = fmt.Sprintf("#%02x%02x%02x", uint8(r), uint8(g), uint8(b))
		}
	}
	return
}

func matchColors(results []*response, colors []string, cfg confProperty) {
	var (
		index float64
	)
	tempRange := cfg.Max - cfg.Min
	colorsLen := float64(len(colors) - 1)
	for _, response := range results {
		index = math.Ceil(response.minTemp) + math.Abs(cfg.Min)
		index = (index * colorsLen) / tempRange
		response.minTempColor = colors[int(index)]
		index = math.Ceil(response.maxTemp) + math.Abs(cfg.Min)
		index = (index * colorsLen) / tempRange
		if index > colorsLen {
			index = colorsLen
		}
		response.maxTempColor = colors[int(index)]
	}

}
