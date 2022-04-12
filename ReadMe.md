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
        index: ["name"]
        documents:
          - key: "root"
            value: '{"name": "admin", "email": "admin@admin", "humanReadableName": "Dr. Hans Wurst"}'
          - key: "test"
            value: '{"name": "test", "email": "test@user", "humanReadableName": "Max Mustermann"}'
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
* index: list of collections fields for witch an index is needed
* documents: initial list of documents for collection

#### document item
* key: document key
* value: document value as json record string literal

## Building
* for building executable run ```make build```
* for creating a docker container ```make docker```

