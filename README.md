# go-pushy
This package is a golang client for [pushy.me](https://pushy.me/).  

## Supported APIs

* [Send notifications](https://pushy.me/docs/api/send-notifications)

## Usage

```sh
go get github.com/healthimation/go-pushy/pushy
```

```golang
import (
    "log"
    "time"

    "github.com/healthimation/go-pushy/pushy"
)

func main() {
    timeout := 5 * time.Second
    client := pushy.NewClient("my pushy api key", timeout)

    // push to devices
    tokens := []string{"device_token1", "device_token2"}
    data := map[string]string{"msg": "hello world!"}
    err := client.PushToDevices(context.Background(), tokens, data, nil) 
    if err != nil {
        log.Printf("Error pushing to devices: %s", err.Error())
    }

    // push to a topic
    topic := "current_events"
    err := client.PushToTopic(context.Background(), topic, data, nil)
    if err != nil {
        log.Printf("Error pushing to topic: %s", err.Error())
    }
}
```
