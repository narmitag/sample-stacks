mysql
sysbench
mysql-exporter
prometheus
grafana
alertmanager
sample alert logger from https://github.com/tomtom-international/alertmanager-webhook-logger

and a sample go app that queries mysql and exposes test metrics to prometheus

Prometheus url localhost:9090
Grafana localhost:3000

Build steps

```bash
docker compose build
docker compose up
```
