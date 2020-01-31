# Xin

[![Go Report Card](https://goreportcard.com/badge/github.com/xieqiaoyu/xin)](https://goreportcard.com/report/github.com/xieqiaoyu/xin)
[![GoDoc](https://godoc.org/github.com/xieqiaoyu/xin?status.svg)](https://godoc.org/github.com/xieqiaoyu/xin)


Xin is a framework focus on building configurable server service easily

Xin is  based on many other fantastic repo,thanks for their author's work!

## At a glance

#### HTTP service

assume the following codes in example.go

```go
package main

import (
    "github.com/gin-gonic/gin"
    "github.com/spf13/cobra"
    "github.com/xieqiaoyu/xin"
    xhttp "github.com/xieqiaoyu/xin/http"
)

var configString = `
# app envirment
env = "dev"
#env = "release"

# http server setting
[http]
    listen=":8080"
`

//HttpDemoService Demo http service
type HttpDemoService struct{}

//RegisterRouter xhttp.ServerInterface implement
func (s *HttpDemoService) RegisterRouter(e *gin.Engine) {
    e.Use(gin.Logger(), gin.Recovery())

    e.GET("/ping", func(c *gin.Context) {
        c.JSON(200, gin.H{
            "message": "pong",
        })
    })
}

//InitializeHTTPServer define an instance of  xhttp.InitializeServerFunc
func InitializeHTTPServer() (xhttp.ServerInterface, error) {
    configLoader := xin.NewStringConfigLoader(configString, "toml")
    config := xin.NewConfig(configLoader, nil)
    //init toml config from a string
    config.Init()

    env := xin.NewEnvSetting(config)
    return xhttp.NewServer(env, config, &HttpDemoService{}), nil
}

func main() {
    httpCmd := xhttp.NewHTTPCmd(InitializeHTTPServer)

    rootCmd := &cobra.Command{}
    rootCmd.AddCommand(httpCmd)
    rootCmd.Execute()
}
```

```bash
# run the demo and visit 0.0.0.0:8080/ping on browser
$ go run example.go http
```



## Getting Start

#### Installation

xin require go 1.13+

```bash
$ go get -u github.com/xieqiaoyu/xin
```

### Command line

xin use [cobra](https://github.com/spf13/cobra) to implement command line entry

xin provide some out-of-the-box subcommand generator such as

`github.com/xieqiaoyu/xin/http.NewHTTPCmd`

Add  subcommands into a cobra Command makes the application flexible

```go
import (
    "github.com/spf13/cobra"
    "github.com/xieqiaoyu/xin"
    xhttp "github.com/xieqiaoyu/xin/http"
)

...

func main() {
    httpCmd := xhttp.NewHTTPCmd(InitializeHTTPServer)
    versionCmd := xin.NewVersionCmd("v0.1.0-dev")


    rootCmd := &cobra.Command{}
    rootCmd.AddCommand(httpCmd)
    rootCmd.AddCommand(versionCmd)
    rootCmd.Execute()
}
```



```
# call version subcommand to get application version
$ go run example.go version
v0.1.0-dev
```



### Use config

### Define http server

### Use database
