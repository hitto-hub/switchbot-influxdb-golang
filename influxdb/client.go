package influxdb

import (
    "context"
    "fmt"
    "github.com/influxdata/influxdb-client-go/v2"
    "switchbot-influxdb/config"
    "switchbot-influxdb/switchbot"
    "time"
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
            fmt.Printf("Error writing point: %v\n", err)
            return err
        } else {
            fmt.Printf("Successfully wrote point for device: %s\n", device.Name)
        }
    }
    return nil
}
