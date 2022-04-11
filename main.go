package main

import (
	"context"
	"flag"
	"fmt"
	driver "github.com/arangodb/go-driver"
	"github.com/arangodb/go-driver/http"
	"github.com/fatih/structs"
	"github.com/golang/glog"
	"gopkg.in/yaml.v3"
	"main/dbconfig"
	"os"
	"time"
)

type DbConnectionValues struct {
	SysUser  string
	SysPass  string
	Endpoint string
	Timeout  time.Duration
	Retry    int
}

type CMDLineArguments struct {
	config             string
	dbConnectionValues DbConnectionValues
}

type DBDef struct {
	Name  string `yaml:",omitempty"`
	Owner struct {
		Name     string `yaml:",omitempty"`
		Password string `yaml:",omitempty"`
	} `yaml:",omitempty"`
	Collections []struct {
		Name  string   `yaml:",omitempty"`
		Index []string `yaml:",omitempty,flow"`
	} `yaml:",omitempty"`
}

type ConfigDef struct {
	Databases []DBDef `yaml:",omitempty"`
}

var arguments CMDLineArguments

func (dbconValues *DbConnectionValues) connect() driver.Client {
	glog.Info("connect to db")
	count := 0
	for count < dbconValues.Retry {
		count += 1
		glog.Info(fmt.Sprintf("retry count: %d", count))
		success := true
		con, err := http.NewConnection(http.ConnectionConfig{
			Endpoints: []string{dbconValues.Endpoint},
		})
		if err != nil {
			glog.Error(err.Error())
			success = false
		}
		client, err := driver.NewClient(driver.ClientConfig{
			Connection:     con,
			Authentication: driver.BasicAuthentication(dbconValues.SysUser, dbconValues.SysPass),
		})
		if err != nil {
			glog.Error(err.Error())
			success = false
		}

		ctx, _ := context.WithTimeout(context.Background(), dbconValues.Timeout)
		databases, err := client.Databases(ctx)
		if err != nil {
			glog.Error(err.Error())
			success = false
		}

		if success {
			for _, database := range databases {
				glog.Info(fmt.Sprintf("database found -> '%s' ", database.Name()))
			}
			glog.Info("connect ok")
			return client
		}
		time.Sleep(dbconValues.Timeout)
	}

	glog.Fatal("arango connect finally failed")
	return nil
}

func init() {
	flag.Set("logtostderr", "true")
	flag.StringVar(&arguments.config, "config", "arangoInit.yaml", "file for config")
	flag.StringVar(&arguments.dbConnectionValues.SysUser, "user", "root", "arango sys user")
	flag.StringVar(&arguments.dbConnectionValues.SysPass, "pass", "root", "arango sys pass")
	flag.StringVar(&arguments.dbConnectionValues.Endpoint, "endpoint", "http://localhost:8529", "arango endpoint")
	flag.DurationVar(&arguments.dbConnectionValues.Timeout, "timeout", 5*time.Second, "request timeout")
	flag.IntVar(&arguments.dbConnectionValues.Retry, "retry", 10, "retry count")
	flag.Parse()

}

func iterateDbs(database DBDef, client driver.Client) {
	if structs.HasZero(database) {
		glog.Fatal(fmt.Sprintf("missing fields in config file for database '%s'", database.Name))
	}
	glog.Info(fmt.Sprintf("check database: '%s'", database.Name))
	db := dbconfig.Database{
		Name: database.Name,
		User: database.Owner.Name,
		Pass: database.Owner.Password,
	}
	err := db.Create(client)
	if err != nil {
		glog.Fatal(err)
	}
	dbcon, err := client.Database(nil, database.Name)
	if err != nil {
		glog.Fatal(err)
	}
	for _, col := range database.Collections {
		glog.Info(fmt.Sprintf("check collection: '%s/%s'", database.Name, col.Name))
		idx := make([]dbconfig.Index, len(col.Index))
		for i, v := range col.Index {
			idx[i].Field = v
			idx[i].Name = fmt.Sprintf("Index_%s", v)
		}
		collection := dbconfig.Collection{
			Name:  col.Name,
			Index: idx,
		}
		err = collection.Create(dbcon)
		if err != nil {
			glog.Fatal(err)
		}
	}
}

func main() {
	glog.Info("start arangoInit")
	glog.Info(fmt.Sprintf("config file: '%s'", arguments.config))
	content, err := os.ReadFile(arguments.config)
	if err != nil {
		glog.Fatal(err.Error())
	}
	config := ConfigDef{}
	yaml.Unmarshal(content, &config)
	if err != nil {
		glog.Fatal(err.Error())
	}

	client := arguments.dbConnectionValues.connect()

	for _, database := range config.Databases {
		iterateDbs(database, client)
	}

	glog.Info("stop arangoInit")
}
