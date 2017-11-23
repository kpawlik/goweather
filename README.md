# goweather

Program to download data from REST service weather.gov based on conf file. Program will split bounding box (gridWSLatLng, gridENLatLng) to grid with size = "gridSize", and request service for data at these points.


## Installation

```
go get github.com/kpawlik/goweather

go install github.com/kpawlik/goweather
```

## Run

goweather.exe -c PATH_TO_CONFIG_FILE

## Config file

JSON file default (conf.json)

### Example 
```
{
    "url": "https://graphical.weather.gov/xml/sample_products/browser_interface/ndfdXMLclient.php?whichClient=GmlLatLonList&gmlListLatLon=%s&featureType=Forecast_Gml2Point&compType=Between&propertyName=",
    "csvDelimiter": ";",
    "workers": 4,
    "step": 50,
    "gridSize": 100,
    "gridWSLatLng": [
        21,
        -127
    ],
    "gridENLatLng": [
        50,
        -60
    ],
    "outFormats": [
        "json",
        "csv"
    ],
    "outFileName": "tmp/weather_out",
    "properties": [
        {
            "name": "mint",
            "max": 100,
            "min": -100,
            "colorsFile": "colors.png"
        },
        {
            "name": "maxt",
            "max": 100,
            "min": -100,
            "colorsFile": "colors.png"
        }
    ]
}
```

* url - address of weather service. Default https://graphical.weather.gov/xml/sample_products/browser_interface/ndfdXMLclient.php?whichClient=GmlLatLonList&gmlListLatLon=%s&featureType=Forecast_Gml2Point&compType=Between&propertyName=
* csvDelimiter - delimiter for CSV out
* workers - number of concurrent workers
* step - number of coordinates which will be send to REST service in one request
* gridSize - size of grid in kilometers
* gridWSLatLng, gridENLatLng - bounds to request of data
* outFileName - name/path to output file name
* properties - max and min temperature properties.
    * name - name of the parameter to download (https://graphical.weather.gov/xml/rest.php#what)
    * min,max - MIN and MAX values to calculate color value for temp.
    * colorsFile - name of file with colors