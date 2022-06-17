package dbconfig

import (
	"errors"
	"fmt"
	"github.com/arangodb/go-driver"
	"github.com/golang/glog"
	"strings"
)

type Index struct {
	Field        string
	Name         string
	Username     string
	Unique       bool
	Sparse       bool
	InBackground bool
}

type TtlIndex struct {
	Field        string
	Name         string
	Username     string
	ExpiresAfter int
	InBackground bool
}

type CompositeIndex struct {
	Fields       []string
	Name         string
	Username     string
	Unique       bool
	Sparse       bool
	InBackground bool
}

type Collection struct {
	Name             string
	Indexes          []Index
	CompositeIndexes []CompositeIndex
	TtlIndex         []TtlIndex
}

func (ttlIdx *TtlIndex) Create(client driver.Collection) error {
	exists, err := client.IndexExists(nil, ttlIdx.Name)
	if err != nil {
		return err
	}
	if !exists {
		_, _, err := client.EnsureTTLIndex(nil, ttlIdx.Field, ttlIdx.ExpiresAfter, &driver.EnsureTTLIndexOptions{
			InBackground: ttlIdx.InBackground,
			Name:         ttlIdx.Username,
			Estimates:    nil,
		})
		if err != nil {
			return err
		}
	}
	return nil
}

func (cIdx *CompositeIndex) Create(client driver.Collection) error {
	exists, err := client.IndexExists(nil, cIdx.Name)
	if err != nil {
		return err
	}
	if !exists {
		if len(cIdx.Fields) < 2 {
			return errors.New("composite index must at least contain two fields")
		}
		_, _, err := client.EnsurePersistentIndex(nil, cIdx.Fields, &driver.EnsurePersistentIndexOptions{
			Unique:       cIdx.Unique,
			Sparse:       cIdx.Sparse,
			InBackground: cIdx.InBackground,
			Name:         cIdx.Username,
			Estimates:    nil,
		})
		if err != nil {
			return err
		}
	}
	return nil
}

func (idx *Index) Create(client driver.Collection) error {
	exists, err := client.IndexExists(nil, idx.Name)
	if err != nil {
		return err
	}
	if !exists {
		_, _, err := client.EnsurePersistentIndex(nil, []string{idx.Field}, &driver.EnsurePersistentIndexOptions{
			Unique:       idx.Unique,
			Sparse:       idx.Sparse,
			InBackground: idx.InBackground,
			Name:         idx.Username,
			Estimates:    nil,
		})
		if err != nil {
			return err
		}
	}
	return nil
}

func (col *Collection) Create(client driver.Database) error {
	exists, err := client.CollectionExists(nil, col.Name)
	if err != nil {
		return err
	}
	var collection driver.Collection
	if !exists {
		collection, err = client.CreateCollection(nil, col.Name, nil)
		if err != nil {
			return err
		}
	} else {
		collection, err = client.Collection(nil, col.Name)
		if err != nil {
			return err
		}
	}
	for _, idx := range col.Indexes {
		glog.Info(fmt.Sprintf("create index for: %s", idx.Field))
		err = idx.Create(collection)
		if err != nil {
			return err
		}
	}
	for _, cIdx := range col.CompositeIndexes {
		glog.Info(fmt.Sprintf("create composit index for: %s", strings.Join(cIdx.Fields, ", ")))
		err = cIdx.Create(collection)
		if err != nil {
			return err
		}
	}
	for _, ttlIdx := range col.TtlIndex {
		glog.Info(fmt.Sprintf("create TTL index for: %s", ttlIdx.Field))
		err = ttlIdx.Create(collection)
		if err != nil {
			return err
		}
	}
	return nil
}
