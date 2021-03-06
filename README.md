<p align="center">
  <h2 align="center">ICON Logs Service</h2>
</p>

[![loopchain](https://img.shields.io/badge/ICON-API-blue?logoColor=white&logo=icon&labelColor=31B8BB)](https://shields.io) [![GitHub Release](https://img.shields.io/github/release/sudoblockio/icon-logs.svg?style=flat)]() ![](https://github.com/sudoblockio/icon-logs/workflows/push-main/badge.svg?branch=main) [![codecov](https://codecov.io/gh/sudoblockio/icon-logs/branch/main/graph/badge.svg)](https://codecov.io/gh/sudoblockio/icon-logs) ![Uptime](https://img.shields.io/endpoint?url=https%3A%2F%2Fraw.githubusercontent.com%2Fsudoblockio%2Ficon-status-page%2Fmaster%2Fapi%2Fdev-logs-service%2Fuptime.json) ![](https://img.shields.io/github/license/sudoblockio/icon-logs)

[Live API Docs](https://explorer.icon.geometry-dev.net/api/v1/logs/docs/)

Off chain indexer for the ICON Blockchain serving the **logs** context of the [icon-explorer](https://github.com/sudoblockio/icon-explorer). Service is broken up into API and worker components that are run as individual docker containers. It depends on data coming in from [icon-etl](https://github.com/sudoblockio/icon-etl) over a Kafka message queue. For websockets and various buffers it uses redis with persistence on a postgres database. 

### Deployment 

Service can be run in the following ways:

1. Independently from this repo with docker compose:
```bash
docker-compose -f docker-compose.db.yml -f docker-compose.yml up -d
# Or alternatively 
make up 
```   

2. With the whole stack from the main [icon-explorer](https://github.com/sudoblockio/icon-explorer) repo. 

Run `make help` for more options. 

### Development 

For local development, you will want to run the `docker-compose.db.yml` as you develop. To run the tests, 

```bash
make test 
```

### License 

Apache 2.0

