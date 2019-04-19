# MiBand

Basic utilities to manipulate Gadgetbrige data and import into InfluxBD.

## Usage

### Server

```shell
server -p <port> -db <influx_db_name> -influxEndpoint <influx_host_port> -user <auth_user> -pass <auth_pass>
```

### Importer

```shell
importer -user <auth_user> -pass <auth_pass> -host <server_endpoint>
```
