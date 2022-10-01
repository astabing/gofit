package main

import (
	"flag"
	"log"
	"time"

	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
)

func main() {
	var forSet, forDate, forTime string

	flag.StringVar(&forSet, "s", "all", `The set for which data is to be collected. Supported: all,heart, Default: all`)
	flag.StringVar(&forDate, "d", "", `The date for which data is to be collected. Expected argument format: YYYY-MM-DD.
Default: current date. If any past date is provided [t] argument will be internally ignored`)
	flag.StringVar(&forTime, "t", "", `The time interval between current time and [t] for which data should be collected. 
Expected argument format: 1h|60m. Maximal interval value is 24h. Default: 00:00`)
	flag.Parse()

	forDate, forTime = parseArgs(forDate, forTime)

	// Init influxdb client environment and open WriteAPI
	InfluxClient, err := idb.newInfluxdb()
	if err != nil {
		log.Fatal(err)
	}
	InfluxWrite := InfluxClient.WriteAPI(idb.org, idb.bucket)

	defer func() {
		// Force all unwritten data to be sent
		InfluxWrite.Flush()
		// Ensures background processes finishes
		InfluxClient.Close()
	}()

	// Init Fitbit API
	fit, err := newFitbit()
	if err != nil {
		log.Fatal(err)
	}

	if forSet == "heart" || forSet == "all" {
		// Intraday heart time series
		heartTimeSeries, err := fit.getActivity("heart", forDate, forTime)
		if err != nil {
			log.Fatal(err)
		}
		for _, v := range *heartTimeSeries {
			p := influxdb2.NewPoint(
				"heartrate",
				map[string]string{"type": v.tp},
				map[string]interface{}{
					"rate": v.value,
				},
				v.time)
			// write asynchronously
			InfluxWrite.WritePoint(p)
		}
	}
	if forSet == "steps" || forSet == "all" {
		stepsTimeSeries, err := fit.getActivity("steps", forDate, forTime)
		if err != nil {
			log.Fatal(err)
		}
		for _, v := range *stepsTimeSeries {
			p := influxdb2.NewPoint(
				"activities",
				map[string]string{"activity": "steps"},
				map[string]interface{}{
					"steps": v.value,
				},
				v.time)
			// write asynchronously
			InfluxWrite.WritePoint(p)
		}
	}
	if forSet == "distance" || forSet == "all" {
		stepsTimeSeries, err := fit.getActivity("distance", forDate, forTime)
		if err != nil {
			log.Fatal(err)
		}
		for _, v := range *stepsTimeSeries {
			p := influxdb2.NewPoint(
				"activities",
				map[string]string{"activity": "distance"},
				map[string]interface{}{
					"distance": v.valueF,
				},
				v.time)
			// write asynchronously
			InfluxWrite.WritePoint(p)
		}
	}
	if forSet == "floors" || forSet == "all" {
		stepsTimeSeries, err := fit.getActivity("floors", forDate, forTime)
		if err != nil {
			log.Fatal(err)
		}
		for _, v := range *stepsTimeSeries {
			p := influxdb2.NewPoint(
				"activities",
				map[string]string{"activity": "floors"},
				map[string]interface{}{
					"floors": v.value,
				},
				v.time)
			// write asynchronously
			InfluxWrite.WritePoint(p)
		}
	}
}

func parseArgs(d, t string) (string, string) {
	ct := time.Now()
	if d != "" {
		dt, err := time.ParseInLocation(dateFormat, d, time.Local)
		if err != nil {
			log.Fatal(err)
		}
		d = dt.Format(dateFormat)
	} else {
		// Empty date argument provided, set to current date
		d = ct.Format(dateFormat)
	}

	if t != "" {
		td, err := time.ParseDuration("-" + t)
		if err != nil {
			log.Fatal(err)
		}
		// Only consider current day time
		// Set 00:00 if provided time duration goes to previous day
		if isToday(ct.Add(td)) {
			t = ct.Add(td).Format(timehmFormat)
		} else {
			t = "00:00"
		}
	} else {
		t = "00:00"
	}

	return d, t
}
