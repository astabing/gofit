package main

// https://dev.fitbit.com/build/reference/web-api/explore/#/

// https://dev.fitbit.com/build/reference/web-api/explore/#/Heart%20Rate%20Intraday%20Time%20Series/getHeartByDateTimestampIntraday
type heartRateIntradayTimeSeries struct {
	ActivitiesHeart []struct {
		CustomHeartRateZones []interface{} `json:"customHeartRateZones"`
		DateTime             string        `json:"dateTime"`
		HeartRateZones       []struct {
			CaloriesOut float64 `json:"caloriesOut"`
			Max         int     `json:"max"`
			Min         int     `json:"min"`
			Minutes     int     `json:"minutes"`
			Name        string  `json:"name"`
		} `json:"heartRateZones"`
		Value string `json:"value"`
	} `json:"activities-heart"`
	ActivitiesHeartIntraday struct {
		Dataset []struct {
			Time  string `json:"time"`
			Value int    `json:"value"`
		} `json:"dataset"`
		DatasetInterval int    `json:"datasetInterval"`
		DatasetType     string `json:"datasetType"`
	} `json:"activities-heart-intraday"`
}

// https://dev.fitbit.com/build/reference/web-api/explore/#/Heart%20Rate%20Intraday%20Time%20Series/getHeartByDateIntraday
type heartRateTimeSeries struct {
	ActivitiesHeart []struct {
		DateTime string `json:"dateTime"`
		Value    struct {
			CustomHeartRateZones []interface{} `json:"customHeartRateZones"`
			HeartRateZones       []struct {
				CaloriesOut float64 `json:"caloriesOut"`
				Max         int     `json:"max"`
				Min         int     `json:"min"`
				Minutes     int     `json:"minutes"`
				Name        string  `json:"name"`
			} `json:"heartRateZones"`
			RestingHeartRate int `json:"restingHeartRate"`
		} `json:"value"`
	} `json:"activities-heart"`
	ActivitiesHeartIntraday struct {
		Dataset []struct {
			Time  string `json:"time"`
			Value int    `json:"value"`
		} `json:"dataset"`
		DatasetInterval int    `json:"datasetInterval"`
		DatasetType     string `json:"datasetType"`
	} `json:"activities-heart-intraday"`
}

type activityStepsTimeSeries struct {
	ActivitiesSteps []struct {
		DateTime string `json:"dateTime"`
		Value    string `json:"value"`
	} `json:"activities-steps"`
	ActivitiesStepsIntraday struct {
		Dataset []struct {
			Time  string `json:"time"`
			Value int    `json:"value"`
		} `json:"dataset"`
		DatasetInterval int    `json:"datasetInterval"`
		DatasetType     string `json:"datasetType"`
	} `json:"activities-steps-intraday"`
}

type activityDistanceTimeSeries struct {
	ActivitiesDistance []struct {
		DateTime string `json:"dateTime"`
		Value    string `json:"value"`
	} `json:"activities-distance"`
	ActivitiesDistanceIntraday struct {
		Dataset []struct {
			Time  string  `json:"time"`
			Value float64 `json:"value"`
		} `json:"dataset"`
		DatasetInterval int    `json:"datasetInterval"`
		DatasetType     string `json:"datasetType"`
	} `json:"activities-distance-intraday"`
}

type activityFloorsTimeSeries struct {
	ActivitiesFloors []struct {
		DateTime string `json:"dateTime"`
		Value    string `json:"value"`
	} `json:"activities-floors"`
	ActivitiesFloorsIntraday struct {
		Dataset []struct {
			Time  string `json:"time"`
			Value int    `json:"value"`
		} `json:"dataset"`
		DatasetInterval int    `json:"datasetInterval"`
		DatasetType     string `json:"datasetType"`
	} `json:"activities-floors-intraday"`
}
