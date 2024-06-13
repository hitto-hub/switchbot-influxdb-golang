package main

import (
    "fmt"
    "switchbot-influxdb/config"
    "switchbot-influxdb/influxdb"
    "switchbot-influxdb/switchbot"
)

func main() {
    config := config.LoadConfig()

    devices, err := switchbot.FetchDevices(config.SwitchBotAPIToken, config.SwitchBotSecret)
    if err != nil {
        fmt.Printf("Error fetching devices: %v\n", err)
        return
    }

    // デバイス情報をターミナルに出力
    for _, device := range devices {
        fmt.Printf("Device ID: %s, Name: %s, Type: %s\n", device.Id, device.Name, device.Type)
    }

    // デバイス情報をInfluxDBに保存
    err = influxdb.StoreDevices(config.InfluxDBConfig, devices)
    if err != nil {
        fmt.Printf("Error storing devices in InfluxDB: %v\n", err)
        return
    }

    fmt.Println("Data successfully stored in InfluxDB")
}
