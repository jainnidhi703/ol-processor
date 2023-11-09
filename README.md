# OpenLineage Event Processor
A `backend` to process OpenLineage Airflow events.

## Requirements

* [GoLang](https://go.dev/)
* [OpenLineage Proxy](https://github.com/OpenLineage/OpenLineage/tree/main/proxy/backend#openlineage-proxy-backend)
* [Airflow with OpenLineage Installed - Sample Docker with DAGS](https://github.com/jainnidhi703/airflow-ol)

## Running the Web Server

To build the entire project run:

Local:
```bash
$ go run cmd/ol-processor/main.go 
```

Docker:
```bash
$ docker-compose up --build
```
> **Note:** Add the configuration settings in `proxy.yml` of OpenLineage Proxy Server. Add respective App ports for Open Lineage and HTTP URL to the Web Applications URL.
```bash
server:
  applicationConnectors:
    - type: http
      port: ${OPENLINEAGE_PROXY_PORT:-4433}
  adminConnectors:
    - type: http
      port: ${OPENLINEAGE_PROXY_ADMIN_PORT:-4434}

logging:
  level: ${LOG_LEVEL:-INFO}
  appenders:
    - type: console

proxy:
  source: openLineageProxyBackend
  streams:
    - type: Console
    - type: Http
      url: http://localhost:3000/api/v1/lineage
```

## Description

Open Lineage will push all the events to `http://localhost:3000/api/v1/lineage`.
The consumed events are aggregated and cached based on the Airflow DAG name, and can be queried via the API `http://localhost:3000/api/get/graph/:dag_name`

The above requests show the Lineage graph as follows:



----


