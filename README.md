# GoZix Zap

## Dependencies

* [viper](https://github.com/gozix/viper)

## Configuration example

time_encoder config can be ("iso8601", "millis", "nanos") and only for "console" and "json" encoding

```json
{
  "zap": {
    "cores": {
      "console": {
        "type": "stream",
        "level": "debug",
        "encoding": "console",
        "message_key": "message",
        "time_encoder": "iso8601"
      },
      "json": {
        "type": "stream",
        "level": "debug",
        "encoding": "json",
        "message_key": "message",
        "time_encoder": "millis"
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
