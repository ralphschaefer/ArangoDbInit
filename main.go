package main

import (
	"flag"
	"fmt"
//	"os"
	"github.com/golang/glog"
//	"gopkg.in/yaml.v3"
//	"github.com/arangodb/go-driver/http"
//	driver "github.com/arangodb/go-driver"
)

var Config string

func init() {
	flag.Set("logtostderr", "true")
	flag.StringVar(&Config, "config", "./temp/arangoInit.yaml", "file for config")
	flag.Parse()
}

func main() {
	glog.Info("run arangoInit")
	glog.Info(fmt.Sprintf("config file: '%s'", Config))
	fmt.Println("Start")
}

