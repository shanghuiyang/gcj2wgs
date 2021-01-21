package main

import (
	"errors"
	"fmt"
	"math"
	"os"
)

const (
	// Krasovsky 1940
	a  = 6378245.0
	ee = 0.00669342162296594323

	delta     = 0.0001
	threshold = 0.000000001
	maxloop   = 1000
)

func main() {
	if len(os.Args) != 2 {
		fmt.Println("error: invalid args")
		printUsage()
		os.Exit(1)
	}
	input := os.Args[1]

	var lat, lon float64
	if n, err := fmt.Sscanf(input, "%f,%f", &lat, &lon); n != 2 || err != nil {
		fmt.Println("error: invalid input")
		printUsage()
		os.Exit(1)
	}

	if lat < -90 || lat > 90 {
		fmt.Println("error: invalid lat, expect range: [-90, 90]")
		os.Exit(1)
	}
	if lon < -180 || lon > 180 {
		fmt.Println("error: invalid lon, expect range: [-180, 180]")
		os.Exit(1)
	}

	wgslat, wgslon, err := gcj2wgs(lat, lon)
	if err != nil {
		fmt.Printf("error: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("%.6f,%.6f\n", wgslat, wgslon)
	os.Exit(0)
}

func wgs2gcj(lat, lon float64) (gcjLat, gcjLon float64, err error) {
	if outOfChina(lat, lon) {
		return 0, 0, errors.New("out of China")
	}
	dLat := transformLat(lat-35.0, lon-105.0)
	dLon := transformLon(lat-35.0, lon-105.0)
	rLat := lat / 180.0 * math.Pi
	magic := math.Sin(rLat)
	magic = 1 - ee*magic*magic
	sqrtMagic := math.Sqrt(magic)
	dLat = (dLat * 180.0) / ((a * (1 - ee)) / (magic * sqrtMagic) * math.Pi)
	dLon = (dLon * 180.0) / (a / sqrtMagic * math.Cos(rLat) * math.Pi)
	gcjLat = lat + dLat
	gcjLon = lon + dLon
	return
}

func gcj2wgs(lat, lon float64) (wgsLat, wgsLon float64, err error) {
	if outOfChina(lat, lon) {
		return 0, 0, errors.New("out of China")
	}

	mLat := lat - delta
	mLon := lon - delta
	pLat := lat + delta
	pLon := lon + delta

	for i := 0; i < maxloop; i++ {
		wgsLat = (mLat + pLat) / 2
		wgsLon = (mLon + pLon) / 2
		gcjlat, gcjlon, err := wgs2gcj(wgsLat, wgsLon)
		if err != nil {
			return 0, 0, err
		}
		dLat := gcjlat - lat
		dLon := gcjlon - lon
		if (math.Abs(dLat) < threshold) && (math.Abs(dLon) < threshold) {
			break
		}
		if dLat > 0 {
			pLat = wgsLat
		} else {
			mLat = wgsLat
		}

		if dLon > 0 {
			pLon = wgsLon
		} else {
			mLon = wgsLon
		}
	}
	return
}

func printUsage() {
	fmt.Println("usage: $gcj2wgs lat,lon")
	fmt.Println("example: $gcj2wgs 39.1,106.1")
}

func outOfChina(lat, lon float64) bool {
	if lat < 0.8293 || lat > 55.8271 {
		return true
	}
	if lon < 72.004 || lon > 137.8347 {
		return true
	}
	return false
}

func transformLat(lat, lon float64) float64 {
	ret := -100.0 + 2.0*lon + 3.0*lat + 0.2*lat*lat + 0.1*lon*lat + 0.2*math.Sqrt(math.Abs(lon))
	ret += (20.0*math.Sin(6.0*lon*math.Pi) + 20.0*math.Sin(2.0*lon*math.Pi)) * 2.0 / 3.0
	ret += (20.0*math.Sin(lat*math.Pi) + 40.0*math.Sin(lat/3.0*math.Pi)) * 2.0 / 3.0
	ret += (160.0*math.Sin(lat/12.0*math.Pi) + 320*math.Sin(lat*math.Pi/30.0)) * 2.0 / 3.0
	return ret
}

func transformLon(lat, lon float64) float64 {
	ret := 300.0 + lon + 2.0*lat + 0.1*lon*lon + 0.1*lon*lat + 0.1*math.Sqrt(math.Abs(lon))
	ret += (20.0*math.Sin(6.0*lon*math.Pi) + 20.0*math.Sin(2.0*lon*math.Pi)) * 2.0 / 3.0
	ret += (20.0*math.Sin(lon*math.Pi) + 40.0*math.Sin(lon/3.0*math.Pi)) * 2.0 / 3.0
	ret += (150.0*math.Sin(lon/12.0*math.Pi) + 300.0*math.Sin(lon/30.0*math.Pi)) * 2.0 / 3.0
	return ret
}
