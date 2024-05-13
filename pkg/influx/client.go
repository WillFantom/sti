package influx

import (
	"context"
	"fmt"
	"time"

	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"github.com/willfantom/influx-speedtest/pkg/tester"
)

func WriteData(URL, org, bucket, token string, testName string, result *tester.Result) error {
	client := influxdb2.NewClientWithOptions(URL, token, influxdb2.DefaultOptions().SetBatchSize(20).AddDefaultTag("sti", "true"))
	alive, err := client.Ping(context.Background())
	if err != nil {
		return fmt.Errorf("failed to ping influx instance: %w", err)
	}
	if !alive {
		return fmt.Errorf("influx instance is not alive")
	}
	writeAPI := client.WriteAPIBlocking(org, bucket)
	return writeAPI.WritePoint(context.Background(), influxdb2.NewPoint(testName, result.Labels, result.Data, time.Now()))
}
