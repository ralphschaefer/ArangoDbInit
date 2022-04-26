package dbconfig

import "github.com/arangodb/go-driver"

type Index struct {
	Field        string
	Name         string
	Username     string
	Unique       bool
	Sparse       bool
	InBackground bool
}

type Collection struct {
	Name  string
	Index []Index
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
	for _, idx := range col.Index {
		err = idx.Create(collection)
		if err != nil {
			return err
		}
	}
	return nil
}
