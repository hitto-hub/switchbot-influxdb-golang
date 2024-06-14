package main

import (
    "log"
    "time"

    "switchbot-influxdb/config"
    "switchbot-influxdb/influxdb"
    "switchbot-influxdb/switchbot"
)

func main() {
    config := config.LoadConfig()

    devices, err := switchbot.FetchDevices(config.SwitchBotAPIToken, config.SwitchBotSecret)
    if err != nil {
        log.Fatalf("Error fetching devices: %v\n", err)
    }

    // デバイスリストの更新 (1日1回)
    go func() {
        for {
            devices, err = switchbot.FetchDevices(config.SwitchBotAPIToken, config.SwitchBotSecret)
            if err != nil {
                log.Printf("Error updating device list: %v", err)
            } else {
                log.Println("Device list updated successfully")
            }
            time.Sleep(24 * time.Hour)
        }
    }()

    // 温湿度データの定期取得 (1分に1回)
    ticker := time.NewTicker(1 * time.Minute)
    defer ticker.Stop()

    for range ticker.C {
        for _, device := range devices {
            if device.Type == "Meter" {
                status, err := switchbot.FetchDeviceStatus(config.SwitchBotAPIToken, config.SwitchBotSecret, device.Id)
                if err != nil {
                    log.Printf("Error fetching device status: %v\n", err)
                    continue
                }

                body, ok := status["body"].(map[string]interface{})
                if !ok {
                    log.Printf("Error: device status body is nil or not in expected format for device ID: %s\n", device.Id)
                    continue
                }

                err = influxdb.StoreMeterData(config.InfluxDBConfig, device.Id, body)
                if err != nil {
                    log.Printf("Error storing meter data in InfluxDB: %v\n", err)
                    continue
                }

                log.Printf("Successfully stored meter data for device ID: %s\n", device.Id)
            }
        }
    }
}
