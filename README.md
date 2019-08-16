# GoZix Zap

## Dependencies

* [viper](https://github.com/gozix/viper)

## Configuration example

time_encoder config can be ("iso8601", "millis", "nanos") and only for "console" and "json" encoding

```json
{
  "zap": {
    "cores": [{
      "level": "debug",
      "encoding": "console",
      "time_encoder": "iso8601"
    }, {
      "level": "debug",
      "encoding": "json",
      "time_encoder": "millis"
    }, {
      "addr": "127.0.0.1:12001",
      "level": "debug",
      "encoding": "gelf"
    }],
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
