# GoZix Zap

## Dependencies

* [viper](https://github.com/gozix/viper)

## Configuration example

```json
{
  "logger": {
    "cores": [{
      "level": "debug",
      "encoding": "console"
    }, {
      "level": "debug",
      "encoding": "json"
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
    "development": true
  }
}
```
