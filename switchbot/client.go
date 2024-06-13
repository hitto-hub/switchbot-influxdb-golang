package switchbot

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"
)

const apiUrl = "https://api.switch-bot.com/v1.1/devices"

type Device struct {
	Id   string `json:"deviceId"`
	Name string `json:"deviceName"`
	Type string `json:"deviceType"`
}

type DevicesResponse struct {
	StatusCode int             `json:"statusCode"`
	Message    string          `json:"message"`
	Body       json.RawMessage `json:"body"`
}

type DeviceList struct {
	Devices []Device `json:"deviceList"`
}

func FetchDevices(apiToken, apiSecret string) ([]Device, error) {
	t := strconv.FormatInt(time.Now().Unix()*1000, 10)
	nonce := "requestID"
	stringToSign := fmt.Sprintf("%s%s%s", apiToken, t, nonce)
	mac := hmac.New(sha256.New, []byte(apiSecret))
	mac.Write([]byte(stringToSign))
	sign := base64.StdEncoding.EncodeToString(mac.Sum(nil))

	client := &http.Client{}
	req, err := http.NewRequest("GET", apiUrl, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+apiToken)
	req.Header.Set("T", t)
	req.Header.Set("Sign", sign)
	req.Header.Set("Nonce", nonce)

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var devicesResponse DevicesResponse
	if err := json.Unmarshal(body, &devicesResponse); err != nil {
		return nil, err
	}

	var deviceList DeviceList
	if err := json.Unmarshal(devicesResponse.Body, &deviceList); err != nil {
		return nil, err
	}

	return deviceList.Devices, nil
}

func FetchDeviceStatus(apiToken, apiSecret, deviceId string) (map[string]interface{}, error) {
	t := strconv.FormatInt(time.Now().Unix()*1000, 10)
	nonce := "requestID"
	stringToSign := fmt.Sprintf("%s%s%s", apiToken, t, nonce)
	mac := hmac.New(sha256.New, []byte(apiSecret))
	mac.Write([]byte(stringToSign))
	sign := base64.StdEncoding.EncodeToString(mac.Sum(nil))

	url := fmt.Sprintf("https://api.switch-bot.com/v1.1/devices/%s/status", deviceId)

	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+apiToken)
	req.Header.Set("T", t)
	req.Header.Set("Sign", sign)
	req.Header.Set("Nonce", nonce)

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var statusResponse map[string]interface{}
	if err := json.Unmarshal(body, &statusResponse); err != nil {
		return nil, err
	}

	return statusResponse, nil
}
