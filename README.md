# MongoDB / gRPC project

## Configuration

Config example:
```YAML
# Application
app:
  chunkSize: 3      # CSV-file import processing chunk size
  logLevel: "info"  # application log level

# gRPC server
server:
  port: 2412        # gRPC server port

# MongoDB
mdb:
  url: "localhost"  # MongoDB URL
  port: 27017       # MongoDB port
  database: db      # MongoDB database name
```

Every config parameter can be overwritten using ENV variables. For example:
```
MDB_TUTORIAL_SERVER_PORT=1234
```

## CLI

All CLI commands have the following flags:
* `--log-level`: (optional) set logging level (default: info);
* `--config`: (optional) path to configuration file;

If config file not specified, defaults are used. Defaults can be overwritten using ENV variables.

### Server

    mdb-tutorial start --config /etc/config.yaml

Command runs a gRPC server, connects to MongoDB and starts serving gRPC services.

    mdb-tutorial client file-server $GOPATH/src/github.com/itiky/mdb-tutorial/build/resources/csv 2413

Command runs HTTP file server which might be used for testing purposes as a source of CSV files. 

Arguments:
* `args[0]`: path to directory containing files for the fileServer;
* `args[1]`: fileServer port;

### Client

    mdb-tutorial client fetch http://fileserver:2412/1.csv

Command downloads, parses and processed price entities file.

Arguments:
* `args[0]`: file path;

    
    mdb-tutorial client list client list --sort-by-price DESC --sort-by-name ASC --skip 10 --limit 100

Command requests price entities with pagination and sort parameters.

IMPORTANT: gRPC service takes into account sort priority, but CLi doesn't.

Flags:
* `--skip 10`: (optional) skip entities;
* `--limit 100`: (optional) limit entities (default 50);
* `--sort-by-name ASC`: (optional) sort entities by product name;
* `--sort-by-price DESC`: (optional) sort entities by product name;
* `--sort-by-timestamp DESC`: (optional) sort entities by import timestamp;

## Test environment

Setup:
1. `cd ${PROJECT_DIR}`
2. `make`
3. `cd ${PROJECT_DIR}/build`
4. `docker-compose up -d`

Those steps will start:
* 2 gRPC servers (ports 2420 and 2421);
* file server with mock CSV files (port 2422;
* MongoDB database (port 27017);
* Nginx as a gRPC requests balancer;

### Example commands

As Docker compose uses TLS for `nginx` CN, we should do client requests withint Docker network.

    # Switch to Docker compose directory
    cd ${PROJECT_DIR}/build
    
    # Request price entries
    docker-compose run --rm client /mdb-tutorial client list --sort-by-price --limit 100
    
    # Fetch CSV files
    docker-compose run --rm client /mdb-tutorial client fetch http://fileserver:2412/1.csv
    docker-compose run --rm client /mdb-tutorial client fetch http://fileserver:2412/2.csv
    docker-compose run --rm client /mdb-tutorial client fetch http://fileserver:2412/partially_invalid.csv
    docker-compose run --rm client /mdb-tutorial client fetch http://fileserver:2412/invalid.csv
    
    # Check logs
    docker-compose logs -f server_1
    docker-compose logs -f server_2

## Implementation details

**CSV-file Fetch operation**

* temporary file is created to reduce RAM usage for large files;
* temporary file is parsed and processed in chunks to reduce RAM usage;
* import data is "map-reduced" in parallel to optimize DB IO operations;
* importing the same CSV-file will produce duplicates (but with different timestamps);

## TODO

- [ ] avoid import duplicates: for example use fileName as a uniqueness factor;
- [ ] raise the test coverage;
- [X] add TLS support to gRPC server/client;
- [X] configure Nginx as a gRPC request balancer;
- [ ] add MongoDB security configuration;
- [ ] optimize storage layer performance (indices?);
- [ ] optimize Docker image size;
