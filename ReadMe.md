# ArangoDb Init

Tool to initialise an ArangoDb with databases, users, collections 
and documents 

## Usage

### startup standalone
```
./arangoinit -config <CONFIGFILE> -user <USER> -pass <PASSWORD> -endpoint <ENDPOINT_URL> -timeout <DURATION> -retry <COUNT>
```
* CONFIGFILE: Path of config file
* USER: ArangoDB root user 
* PASSWORD: ArangoDB root password
* ENDPOINT_URL: ArangoDB endpoint url
* DURATION: Timeout for ArangoDB requests
* COUNT: retry counts for ArangoDB connections


### startup in docker
```
docker run -d arangodb-cloudinit:0.1.0
```
### configuration file
Example: 

```yaml
databases:
  - name: "aaaa"
    owner:
      name: "hhhh"
      password: "pass"
    collections:
      - name: "users"
        compositeindex:
          - fields: ["email", "name"]
            options: ["unique"]        
        index:
          - field: "name"
            options: ["sparse"]
        documents:
          - key: "root"
            value: '{"name": "admin", "email": "admin@admin", "humanReadableName": "Dr. Hans Wurst"}'
          - key: "test"
            value: '{"name": "test", "email": "test@user", "humanReadableName": "Max Mustermann"}'
          - key: "test2"
            value: '{"name": "test", "email": "test2@user", "humanReadableName": "Max Mustermann"}'
      - name: "dddddd"
        index: []
  - name: "bbbb"
    owner:
      name: "sdsdf"
      password: "sdfgsdfgd"
    collections:
      - name: "asasa"
        index: ["ddd", "ffff"]
      - name: "dddddd"
        index: ["ddd", "ffff"]
```
#### database item
* name: database name
* owner: database own
* collections: database collections

#### owner item
* name: username
* password: password of owner

#### collection item
* name: name of collection
* index: list of _index items_
* compositeindex: list of _composite index items_
* documents: initial list of documents for collection

#### index item
* field: name of collection field to be indexed
* options: index options. Possible option values are: 
    - Unique
    - Sparse
    - InBackground

#### composite index items
* fields: list of names of collection fields to be used as a composite index
* options: index options. Possible option values are: 
    - Unique
    - Sparse
    - InBackground

#### document item
* key: document key
* value: document value as json record string literal

## Building
* for building executable run ```make build```
* for creating a docker container ```make docker```

