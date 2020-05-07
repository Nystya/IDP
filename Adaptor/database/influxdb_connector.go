package database

import (
	_ "github.com/influxdata/influxdb1-client"
	influx "github.com/influxdata/influxdb1-client/v2"
	"golang.org/x/net/context"
	"time"
)

type InfluxDB interface {
	AddMeasurements(ctx context.Context, measurements []string, fields []map[string]interface{}, tags map[string]string, ts time.Time) error
}

type InfluxDBConnector struct {
	DBAdress string
	DBName   string
	Username string
	Password string
}

func NewInfluxDBConnector(dbaddress, dbname, username, password string) *InfluxDBConnector {
	return &InfluxDBConnector{
		DBAdress: dbaddress,
		DBName:   dbname,
		Username: username,
		Password: password,
	}
}

func (ifx *InfluxDBConnector) AddMeasurements(ctx context.Context, measurements []string, fields []map[string]interface{}, tags map[string]string, ts time.Time) error {
	c, err := influx.NewHTTPClient(influx.HTTPConfig{
		Addr:     ifx.DBAdress,
		Username: ifx.Username,
		Password: ifx.Password,
	})
	if err != nil {
		return err
	}
	defer c.Close()

	batchPoints, err := influx.NewBatchPoints(influx.BatchPointsConfig{Database: ifx.DBName})
	if err != nil {
		return err
	}

	for idx, measurementName := range measurements {
		point, err := influx.NewPoint(measurementName, tags, fields[idx], ts)
		if err != nil {
			return err
		}

		batchPoints.AddPoint(point)
	}

	if err = c.Write(batchPoints); err != nil {
		return err
	}

	return nil
}
