## Deployment

1.
```
mv env.example .env
```

2. Set `.env` variables according to local environment.
Influxdb variables described on official Influxdb docker image page: https://hub.docker.com/_/influxdb/.
Grafana - https://grafana.com/docs/grafana/v9.0/setup-grafana/installation/docker/

3. Grafana runs under `472` uid, thus Grafana data dir ownership should be set:
```
chown 472 grafana/data
```

4. 
```
docker-compose up -d
```

5. Setup `gofit` environment:
```
cp -a gofit/env_gofit.example gofit/.env_gofit
```

6. Run `gofit` to collect current day data:
```
./gofit/run.sh
```
 
