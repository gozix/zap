# GoZix Zap

## Dependencies

* [viper](https://github.com/gozix/viper)

## Configuration example

```json
{
  "zap": {
    "cores": {
      "console": {
        "type": "stream",
        "level": "debug",
        "encoding": "console"
      },
      "json": {
        "type": "stream",
        "level": "debug",
        "encoding": "json"
      }
    },
    "caller": true,
    "fields": [{
      "key": "team",
      "value": "any team name"
    }, {
      "key": "service",
      "value": "any service name"
    }],
    "stacktrace": "error",
    "message_key": "message",
    "development": true
  }
}
```

## Built-in Tags

| Symbol                          | Value              | Description     | 
| ------------------------------- | ------------------ | ----------------|
| [core.TagFactory](core/core.go) | zap.core.factory   | Add an factory  |

## Cores

- [gelf](https://github.com/gozix/zap-gelf)
- [stream](core/stream/stream.go)
