# switchbot-influxdb-golang

## graphana

```flex
from(bucket: "データベース名")
  |> range(start: -1h)
  |> filter(fn: (r) => r._measurement == "meter" and r._field == "temperature")
```

```flex
from(bucket: "データベース名")
  |> range(start: -1h)
  |> filter(fn: (r) => r._measurement == "meter" and r._field == "humidity")
```
