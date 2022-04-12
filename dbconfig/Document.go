package dbconfig

import (
	"encoding/json"
	"fmt"
	"github.com/arangodb/go-driver"
	"github.com/golang/glog"
)

type Document struct {
	Key   string
	Value map[string]interface{}
}

func NewDocument(key string, value string) (*Document, error) {
	var d Document
	err := json.Unmarshal([]byte(value), &d.Value)
	if err != nil {
		return nil, err
	}
	d.Key = key
	return &d, nil
}

func (doc *Document) Create(client driver.Collection) error {
	glog.Info(fmt.Sprintf("check document in '%s' key '%s'", client.Name(), doc.Key))
	exists, err := client.DocumentExists(nil, doc.Key)
	if err != nil {
		return err
	}
	if !exists {
		var v map[string]interface{}
		v = doc.Value
		v["_key"] = doc.Key
		_, err = client.CreateDocument(nil, v)
		if err != nil {
			return err
		}
	}
	return nil
}
