# MiBand

Basic utilities to manipulate [Gadgetbrige](https://codeberg.org/Freeyourgadget/Gadgetbridge) data and import into InfluxDB.

## Build

### Raspberry Pi Server

```shell
GOOS=linux GOARCH=arm GOARM=7 go build ./cmd/server
```

## Usage

### Server

```shell
server -p <port> -db <influx_db_name> -influxEndpoint <influx_host_port> -user <auth_user> -pass <auth_pass>
```

### Importer

```shell
importer -user <auth_user> -pass <auth_pass> -host <server_endpoint>
```
