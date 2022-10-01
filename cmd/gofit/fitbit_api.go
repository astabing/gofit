package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
)

const (
	fitHost               = "https://api.fitbit.com"
	fitHeartPath          = "/1/user/-/activities/heart/date/today/today/1sec/time/%s/%s.json"
	fitHeartDatePath      = "/1/user/-/activities/heart/date/%s/1d/1sec.json"
	fitActivitiesPath     = "/1/user/-/activities/%s/date/today/today/15min/time/%s/%s.json"
	fitActivitiesDatePath = "/1/user/-/activities/%s/date/%s/1d.json"
	dateFormat            = "2006-01-02"
	timehmFormat          = "15:04"
	timehmsFormat         = "15:04:05"
)

var errInvalidActivity = errors.New("Invalid activity type")

type fit struct {
	Client *http.Client
	Data   timeSeriesWriter
	Host   string
	// /1/user/-/activities/heart/date/{date}/{end-date}/{detail-level}/time/{start-time}/{end-time}.json
	HeartPath string
	// /1/user/-/activities/heart/date/{date}/{end-date}/{detail-level}.json
	HeartDatePath string
	// /1/user/-/activities/steps/date/{date}/{end-date}/{detail-level}/time/{start-time}/{end-time}.json
	ActivitiesPath string
	// /1/user/-/activities/steps/date/{date}/{period}.json
	ActivitiesDatePath string
}

type timeSeries struct {
	tp     string
	time   time.Time
	value  int
	valueF float64
}

type timeSeriesWriter interface {
	getTimeSeries(date string, body *[]byte) (*[]timeSeries, error)
}

// newFitbit creates Fitbit API client
func newFitbit() (*fit, error) {
	var data timeSeriesWriter
	client, err := GetOauthClient()
	if err != nil {
		return nil, err
	}

	return &fit{
		Client:             client,
		Data:               data,
		Host:               fitHost,
		HeartPath:          fitHeartPath,
		HeartDatePath:      fitHeartDatePath,
		ActivitiesPath:     fitActivitiesPath,
		ActivitiesDatePath: fitActivitiesDatePath,
	}, nil

}

// Returns true if provided timestamp is today
func isToday(d time.Time) bool {
	ct := time.Now()
	return d.Year() == ct.Year() && d.Month() == ct.Month() && d.Day() == ct.Day()
}

// Returns time.Time from concatenated date and time string arguments
func getTimestamp(d, t string) time.Time {
	var tt time.Time
	dd, err := time.ParseInLocation(dateFormat, d, time.Local)
	if err != nil {
		log.Fatalf("Date parse error: %v\n", err)
	}

	if len(t) == 5 {
		// "00:00" time string
		tt, err = time.ParseInLocation(timehmFormat, t, time.Local)
	} else {
		// "00:00:00" time string
		tt, err = time.ParseInLocation(timehmsFormat, t, time.Local)
	}
	if err != nil {
		log.Fatalf("Time parse error: %v\n", err)
	}

	return time.Date(dd.Year(), dd.Month(), dd.Day(), tt.Hour(), tt.Minute(), tt.Second(), 0, time.Local)
}

// Formats and returns API URL based on provided activity type and datetime
func (fit *fit) setAPIURL(tp, d, t string) (string, error) {
	var apiURL string
	var today = isToday(getTimestamp(d, t))

	if today && t != "00:00" {
		switch tp {
		case "heart":
			apiURL = fmt.Sprintf(fit.Host+fit.HeartPath, t, time.Now().Format(timehmFormat))
			fit.Data = &heartRateIntradayTimeSeries{}
		case "steps":
			apiURL = fmt.Sprintf(fit.Host+fit.ActivitiesPath, tp, t, time.Now().Format(timehmFormat))
			fit.Data = &activityStepsTimeSeries{}
		case "distance":
			apiURL = fmt.Sprintf(fit.Host+fit.ActivitiesPath, tp, t, time.Now().Format(timehmFormat))
			fit.Data = &activityDistanceTimeSeries{}
		case "floors":
			apiURL = fmt.Sprintf(fit.Host+fit.ActivitiesPath, tp, t, time.Now().Format(timehmFormat))
			fit.Data = &activityFloorsTimeSeries{}
		default:
			return "", errInvalidActivity
		}
	} else {
		switch tp {
		case "heart":
			apiURL = fmt.Sprintf(fit.Host+fit.HeartDatePath, d)
			fit.Data = &heartRateTimeSeries{}
		case "steps":
			apiURL = fmt.Sprintf(fit.Host+fit.ActivitiesDatePath, tp, d)
			fit.Data = &activityStepsTimeSeries{}
		case "distance":
			apiURL = fmt.Sprintf(fit.Host+fit.ActivitiesDatePath, tp, d)
			fit.Data = &activityDistanceTimeSeries{}
		case "floors":
			apiURL = fmt.Sprintf(fit.Host+fit.ActivitiesDatePath, tp, d)
			fit.Data = &activityFloorsTimeSeries{}
		default:
			return "", errInvalidActivity
		}
	}

	return apiURL, nil
}

func (fit *fit) getActivity(tp, d, t string) (*[]timeSeries, error) {
	apiURL, err := fit.setAPIURL(tp, d, t)
	if err != nil {
		return nil, err
	}

	res, err := fit.Client.Get(apiURL)
	if err != nil {
		return nil, err
	}

	body, err := io.ReadAll(res.Body)
	defer res.Body.Close()
	if res.StatusCode > 299 {
		return nil, fmt.Errorf("Response failed with status code: %d and\nbody: %s", res.StatusCode, body)
	}
	if err != nil {
		return nil, err
	}

	return fit.Data.getTimeSeries(d, &body)
}

func (data *heartRateTimeSeries) getTimeSeries(date string, body *[]byte) (*[]timeSeries, error) {
	var tp = "heart"
	var t time.Time
	if err := json.Unmarshal(*body, &data); err != nil {
		return nil, err
	}
	ts := make([]timeSeries, len(data.ActivitiesHeartIntraday.Dataset))
	for i, v := range data.ActivitiesHeartIntraday.Dataset {
		t = getTimestamp(date, v.Time)
		ts[i] = timeSeries{tp: tp, time: t, value: v.Value}
	}
	// Add resting heart rate
	for _, v := range data.ActivitiesHeart {
		t = getTimestamp(v.DateTime, "00:00:00")

		for i := 0; i < 8; i++ {
			// Add 8 values a day with 3h interval. Allow Grafana to use smaller than day intervals
			ts = append(ts,
				timeSeries{tp: "resting", time: t, value: v.Value.RestingHeartRate})
			t = t.Add(time.Hour * 3)
		}
	}

	return &ts, nil
}

func (data *heartRateIntradayTimeSeries) getTimeSeries(date string, body *[]byte) (*[]timeSeries, error) {
	var tp = "heart"
	var t time.Time
	if err := json.Unmarshal(*body, &data); err != nil {
		return nil, err
	}
	ts := make([]timeSeries, len(data.ActivitiesHeartIntraday.Dataset))
	for i, v := range data.ActivitiesHeartIntraday.Dataset {
		t = getTimestamp(date, v.Time)
		ts[i] = timeSeries{tp: tp, time: t, value: v.Value}
	}

	return &ts, nil
}

func (data *activityStepsTimeSeries) getTimeSeries(date string, body *[]byte) (*[]timeSeries, error) {
	var tp = "steps"
	var t time.Time
	if err := json.Unmarshal(*body, &data); err != nil {
		return nil, err
	}
	ts := make([]timeSeries, len(data.ActivitiesStepsIntraday.Dataset))
	for i, v := range data.ActivitiesStepsIntraday.Dataset {
		t = getTimestamp(date, v.Time)
		ts[i] = timeSeries{tp: tp, time: t, value: v.Value}
	}

	return &ts, nil
}

func (data *activityFloorsTimeSeries) getTimeSeries(date string, body *[]byte) (*[]timeSeries, error) {
	var tp = "floors"
	var t time.Time
	if err := json.Unmarshal(*body, &data); err != nil {
		return nil, err
	}
	ts := make([]timeSeries, len(data.ActivitiesFloorsIntraday.Dataset))
	for i, v := range data.ActivitiesFloorsIntraday.Dataset {
		t = getTimestamp(date, v.Time)
		ts[i] = timeSeries{tp: tp, time: t, value: v.Value}
	}

	return &ts, nil
}

func (data *activityDistanceTimeSeries) getTimeSeries(date string, body *[]byte) (*[]timeSeries, error) {
	var tp = "distance"
	var t time.Time
	if err := json.Unmarshal(*body, &data); err != nil {
		return nil, err
	}
	ts := make([]timeSeries, len(data.ActivitiesDistanceIntraday.Dataset))
	for i, v := range data.ActivitiesDistanceIntraday.Dataset {
		t = getTimestamp(date, v.Time)
		ts[i] = timeSeries{tp: tp, time: t, valueF: v.Value}
	}

	return &ts, nil
}
