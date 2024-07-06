# Go Logger

A high-performance logger library for Go.

## Installation

To install the package, run:

```bash
go get github.com/yourusername/go-logger
```

## Usage

    go get github.com/yourusername/go-logger


```bash
package main

import (
    "github.com/yourusername/go-logger"
    "go.uber.org/zap"
)

func main() {
    config := logger.LoggerConfig{
        Format:   "json",
        LogType:  "file",
        Priority: "debug",
    }
    log := logger.InitLogger(config)
    defer log.Sync()

    log.Info("This is an info message")
    log.Debug("This is a debug message")
    log.Error("This is an error message")
}
```


