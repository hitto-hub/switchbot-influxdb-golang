package main

import (
    "log"
    "os"
    "os/signal"
    "sync"
    "syscall"
    "time"

    "switchbot-influxdb/config"
    "switchbot-influxdb/influxdb"
    "switchbot-influxdb/switchbot"
)

func main() {
    // 設定の読み込み
    cfg, err := config.LoadConfig()
    if err != nil {
        log.Fatalf("設定の読み込み中にエラーが発生しました: %v\n", err)
    }

    // デバイスの取得
    devices, err := switchbot.FetchDevices(cfg.SwitchBotAPIToken, cfg.SwitchBotSecret)
    if err != nil {
        log.Fatalf("デバイスの取得中にエラーが発生しました: %v\n", err)
    }

    var mu sync.Mutex // デバイスリストへのアクセスを保護するためのミューテックス

    // デバイスリストの更新 (1日1回)
    go func() {
        for {
            updatedDevices, err := switchbot.FetchDevices(cfg.SwitchBotAPIToken, cfg.SwitchBotSecret)
            if err != nil {
                log.Printf("デバイスリストの更新中にエラーが発生しました: %v", err)
            } else {
                log.Println("デバイスリストが正常に更新されました")
                mu.Lock()
                devices = updatedDevices // デバイスリストを更新
                mu.Unlock()
            }
            time.Sleep(24 * time.Hour)
        }
    }()

    // 温湿度データの定期取得 (1分に1回)
    ticker := time.NewTicker(1 * time.Minute)
    defer ticker.Stop()

    // 優雅なシャットダウンの処理
    stop := make(chan os.Signal, 1)
    signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

    go func() {
        for {
            select {
            case <-ticker.C:
                mu.Lock()
                if devices == nil {
                    log.Println("データを取得するためのデバイスがありません。")
                    mu.Unlock()
                    continue
                }

                for _, device := range devices {
                    if device.Type == "Meter" {
                        log.Printf("デバイスID %s のメーターデータを取得中です\n", device.Id)
                        status, err := switchbot.FetchDeviceStatus(cfg.SwitchBotAPIToken, cfg.SwitchBotSecret, device.Id)
                        if err != nil {
                            log.Printf("デバイスステータスの取得中にエラーが発生しました: %v\n", err)
                            continue
                        }

                        body, ok := status["body"].(map[string]interface{})
                        if !ok {
                            log.Printf("エラー: デバイスID %s のステータスボディがnilまたは期待された形式ではありません\n", device.Id)
                            continue
                        }

                        err = influxdb.StoreMeterData(cfg.InfluxDBConfig, device.Id, body)
                        if err != nil {
                            log.Printf("InfluxDBへのメーターデータの保存中にエラーが発生しました: %v\n", err)
                            continue
                        }

                        log.Printf("デバイスID %s のメーターデータが正常に保存されました\n", device.Id)
                    }
                }
                mu.Unlock()
            case <-stop:
                log.Println("優雅にシャットダウンしています...")
                ticker.Stop()
                return
            }
        }
    }()

    // メインのゴルーチンをシャットダウンシグナルを受け取るまでブロック
    <-stop
    log.Println("サービスが停止しました")
}
