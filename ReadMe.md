# Xin

[![Go Report Card](https://goreportcard.com/badge/github.com/xieqiaoyu/xin)](https://goreportcard.com/report/github.com/xieqiaoyu/xin)
[![GoDoc](https://godoc.org/github.com/xieqiaoyu/xin?status.svg)](https://godoc.org/github.com/xieqiaoyu/xin)


Xin is a framework focus on building configurable server service easily

It is based on many other fantastic repo,thanks for their author's work!

## At a glance

#### HTTP server service

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

xin use [cobra](https://github.com/spf13/cobra) to implement app command line entry

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

read [this doc](https://github.com/spf13/cobra#overview) for more usage about cobra

```
# call version subcommand to get application version
$ go run example.go version
v0.1.0-dev
```



### Using config

Config is something could not be decided while coding or something should be assigned later.

A config can be in many form (file、string and so on)  and  have many type (`json`、`yaml`、`toml`).

Xin provide a config struct base on [viper](https://github.com/spf13/viper) . 

```go
type Config struct{
  ...
}
func NewConfig(configloader ConfigLoader, configVerifier ConfigVerifier) *Config
```

```go
//  create a new xin config
config := xin.NewConfig(configLoader,nil)
```

A config need a `ConfigLoader` and an optional `ConfigVerifier`

A `ConfigLoader` define how to load config

```
type ConfigLoader interface {
    LoadConfig(vc *viper.Viper) error
}
```

You can use your own config loader if you need . For convenience, xin provide several config loader

```go
// load config from file
fileConfigLoader := xin.NewFileConfigLoader("another_config.toml","toml")
fileConfigLoader := xin.NewFileConfigLoader("another_config.json","json")

// load config from string
stringConfigLoader := xin.NewStringConfigLoader(configString, "toml")
```

Config should call `Init`before realy use , or there will be a panic

```
config := xin.NewConfig(configLoader,nil)
err := config.Init()
```

use the  [viper](https://github.com/spf13/viper) instance to get config setting, read [this doc](https://github.com/spf13/viper#getting-values-from-viper)  for more usage about viper

```
v := config.Viper()
httpListen := v.GetString("http.listen")
env := v.GetString("env")
```



### HTTP server

The concept of http server in xin is  a generator of `net/http.Server`.

Watch these definition in `xin/http` :

```go
//package xin/http
import (
  "net/http"
  "github.com/spf13/cobra"
)

//ServerInterface a server can provide http server
type ServerInterface interface {
    // provide the http server service
    GetHTTPServer() *http.Server
}

//InitializeServerFunc an init http Server function gives the posibility for dependence inject
type InitializeServerFunc func() (ServerInterface, error)

//NewHTTPCmd Get a cobra command start http server
func NewHTTPCmd(getServer InitializeServerFunc) *cobra.Command
```

You can get a cobra command which can start http server by  calling `NewHTTPCmd`

You need to provide a function to tell xin how to get the http server , **beware**  the server we are talking about is a `ServerInterface`

#### Using xin http Server

Xin has an  `ServerInterface` implementation : `xin/http.Server` , it use [gin](https://github.com/gin-gonic/gin) as a Low-level implementation

```go
//package github.com/xieqiaoyu/xin/http
import (
  "github.com/gin-gonic/gin"
  "github.com/xieqiaoyu/xin"
  "net/http"
)

//Service http service interface
type Service interface {
    // register route and middleware into gin engine
    RegisterRouter(*gin.Engine)
}
//ServerConfig config provide HTTP server setting
type ServerConfig interface {
    HTTPListen() string
}

//NewServer Create a new HTTP server
func NewServer(env xin.Envirment, config ServerConfig, service Service) *Server
```

Interface `Service`  register router and middleware into gin engine, you should implement your own service to complete the app. 

env and config are all interface too , you can implement your own ,or you can use `xin.EnvSetting` and `xin.Config` directly

see [demo](#http-server-service) as an example usage

##### middlewares and tools for RESTful api

gin has already provide a good RESTful api develop experience , still we have to  deal  some annoying things

xin provide some handy middlewares and tools to make life easier

###### api response status



###### api authroization 



### Database

