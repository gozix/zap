# GoZix Zap

[documentation-img]: https://img.shields.io/badge/godoc-reference-blue.svg?color=24B898&style=for-the-badge&logo=go&logoColor=ffffff
[documentation-url]: https://pkg.go.dev/github.com/gozix/zap/v3
[license-img]: https://img.shields.io/github/license/gozix/zap.svg?style=for-the-badge
[license-url]: https://github.com/gozix/zap/blob/master/LICENSE
[release-img]: https://img.shields.io/github/tag/gozix/zap.svg?label=release&color=24B898&logo=github&style=for-the-badge
[release-url]: https://github.com/gozix/zap/releases/latest
[build-status-img]: https://img.shields.io/github/actions/workflow/status/gozix/zap/go.yml?logo=github&style=for-the-badge
[build-status-url]: https://github.com/gozix/zap/actions
[go-report-img]: https://img.shields.io/badge/go%20report-A%2B-green?style=for-the-badge
[go-report-url]: https://goreportcard.com/report/github.com/gozix/zap
[code-coverage-img]: https://img.shields.io/codecov/c/github/gozix/zap.svg?style=for-the-badge&logo=codecov
[code-coverage-url]: https://codecov.io/gh/gozix/zap

[![License][license-img]][license-url]
[![Documentation][documentation-img]][documentation-url]

[![Release][release-img]][release-url]
[![Build Status][build-status-img]][build-status-url]
[![Go Report Card][go-report-img]][go-report-url]
[![Code Coverage][code-coverage-img]][code-coverage-url]

The bundle provide a Zap integration to GoZix application.

## Installation

```shell
go get github.com/gozix/zap/v3
```

## Dependencies

* [viper](https://github.com/gozix/zap)

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

## Built-in DI options

| Name          | Description     | 
|---------------| ----------------|
| AsCoreFactory | Add an factory  |

## Cores

- [gelf](https://github.com/gozix/zap-gelf)
- [stream](stream.go)

## Documentation

You can find documentation on [pkg.go.dev][documentation-url] and read source code if needed.

## Questions

If you have any questions, feel free to create an issue.
