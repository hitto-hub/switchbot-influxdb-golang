package switchbot

import "encoding/json"

type Device struct {
    Id   string `json:"deviceId"`
    Name string `json:"deviceName"`
    Type string `json:"deviceType"`
}

type DevicesResponse struct {
    StatusCode int             `json:"statusCode"`
    Message    string          `json:"message"`
    Body       json.RawMessage `json:"body"` // RawMessageに変更
}

type DeviceList struct {
    Devices []Device `json:"deviceList"`
}
