package influxdb

import (
	"context"
	"fmt"
	"time"

	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"switchbot-influxdb/config"
	"switchbot-influxdb/switchbot"
)

func StoreDevices(config config.InfluxDBConfig, devices []switchbot.Device) error {
	client := influxdb2.NewClient(config.URL, config.Token)
	defer client.Close()

	writeAPI := client.WriteAPIBlocking(config.Org, config.Bucket)
	for _, device := range devices {
		p := influxdb2.NewPoint("device",
			map[string]string{"device_id": device.Id, "device_type": device.Type},
			map[string]interface{}{"name": device.Name},
			time.Now())
		if err := writeAPI.WritePoint(context.Background(), p); err != nil {
			return err
		} else {
			fmt.Printf("Successfully wrote point for device: %s\n", device.Name)
		}
	}
	return nil
}

func StoreMeterData(config config.InfluxDBConfig, deviceID string, status map[string]interface{}) error {
	client := influxdb2.NewClient(config.URL, config.Token)
	defer client.Close()

	writeAPI := client.WriteAPIBlocking(config.Org, config.Bucket)

	fields := map[string]interface{}{
		"temperature": status["temperature"],
		"humidity":    status["humidity"],
		"battery":     status["battery"],
	}

	p := influxdb2.NewPoint("meter",
		map[string]string{"device_id": deviceID},
		fields,
		time.Now())

	if err := writeAPI.WritePoint(context.Background(), p); err != nil {
		return err
	}
	fmt.Printf("Successfully wrote meter data for device ID: %s\n", deviceID)
	return nil
}
