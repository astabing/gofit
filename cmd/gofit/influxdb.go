package main

import (
	"errors"
	"os"

	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
)

type influxdb struct {
	host   string
	org    string
	bucket string
	token  string
}

var errInfluxdbUndefEnv = errors.New("InfluxDB env variables undefined")
var idb = &influxdb{}

func (idb *influxdb) newInfluxdb() (influxdb2.Client, error) {
	idb.host = os.Getenv("INFLUXDB_HOST")
	idb.org = os.Getenv("INFLUXDB_ORG")
	idb.bucket = os.Getenv("INFLUXDB_BUCKET")
	idb.token = os.Getenv("INFLUXDB_TOKEN")

	if idb.host == "" || idb.org == "" || idb.bucket == "" || idb.token == "" {
		return nil, errInfluxdbUndefEnv
	}

	// Return non-blocking write client
	return influxdb2.NewClientWithOptions(idb.host, idb.token,
		influxdb2.DefaultOptions().SetBatchSize(120)), nil

}
