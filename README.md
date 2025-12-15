# Monitoring Lab

<p align="center"><img src="assets/go-spring-logo.png" width="300" alt="Monitoring Lab logo"></p>

This project shows a complete monitoring setup with two apps (Go and Spring Boot), MongoDB, Prometheus and Grafana.

Project structure

```
/monitoring-lab
    /go-app              # Go application
    /spring-app          # Spring Boot application
    /prometheus          # Prometheus config
    /grafana             # Grafana provisioning files
    docker-compose.yml
    README.md            # Portuguese original
    README.en.md         # This file (English B1)
```

Requirements

- Docker
- Docker Compose or Podman Compose

Start all services

Build and start all containers:

```bash
docker-compose up --build
```

Stop services:

```bash
docker-compose down
```

Remove volumes too:

```bash
docker-compose down -v
```

Access links

- Go App (health): `http://localhost:8080/ping`
- Spring App (health): `http://localhost:8081/ping`
- Prometheus: `http://localhost:9090`
- Grafana: `http://localhost:3000` (user: `admin`, pass: `admin`)

Metrics

Go App (`/metrics`):
- `mongodb_pings_total` - total ping operations
- `mongodb_inserts_total` - total insert operations
- `mongodb_operation_duration_seconds` - operation latencies

Spring App (`/actuator/prometheus`):
- `mongodb_total_operations` - total operations
- `mongodb_count_velocity` - rate of counts
- `mongodb_operation_latency` - latencies

Other exporters

- MongoDB exporter (metrics about MongoDB)
- Node exporter (system metrics: CPU, memory, disk)

Open Prometheus

1. Go to `http://localhost:9090`
2. Use the query box, for example:

```
rate(mongodb_pings_total[1m])
```

Open Grafana

1. Go to `http://localhost:3000`
2. Login with `admin/admin`
3. The Prometheus data source is already configured
4. A dashboard called "Monitoring Lab - Go vs Java" is available

Example queries for the dashboard

- Go pings per second:

```
rate(mongodb_pings_total[1m])
```

- Go inserts per second:

```
rate(mongodb_inserts_total[1m])
```

- Go average latency:

```
rate(mongodb_operation_duration_seconds_sum[1m]) / rate(mongodb_operation_duration_seconds_count[1m])
```

- Java total operations:

```
mongodb_total_operations
```

Configuration notes

- Go app: `MONGO_URI` (default `mongodb://mongo:27017`)
- Spring app: `SPRING_DATA_MONGODB_URI` (default `mongodb://mongo:27017/monitoring`)
- Prometheus config file: `prometheus/prometheus.yml` (scrape interval 5s and targets)

Services in docker-compose

1. MongoDB - port `27017`
2. Go App - port `8080`
3. Spring App - port `8081`
4. MongoDB Exporter - port `9216`
5. Node Exporter - port `9100`
6. Prometheus - port `9090`
7. Grafana - port `3000`

Quick checks

1. Application health:

```bash
curl http://localhost:8080/ping
curl http://localhost:8081/ping
```

2. Metrics:

```bash
curl http://localhost:8080/metrics
curl http://localhost:8081/actuator/prometheus
```

3. Logs while running:

```bash
docker-compose logs -f go-app
docker-compose logs -f spring-app
```

Troubleshooting

- Apps not connecting to MongoDB:

```bash
docker-compose ps mongo
docker-compose logs mongo
```

- Prometheus not scraping:

Open `http://localhost:9090/targets` and check targets are `UP`.

- Grafana not connecting to Prometheus:

Open `http://localhost:3000/connections/datasources` and check Prometheus is set up.
