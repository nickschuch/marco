Marco
=====

An API driver load balancer for modern day application infrastructure.

### Balancers and backends

**Balancer**

* Round - We are using https://github.com/mailgun/oxy under the hood for Round Robin balancing.

**Drivers**

* AWS ECS - https://github.com/nickschuch/marco-ecs

### Usage

#### Binary

```bash
$ marco --port=80 \
        --receive=81
```

### Docker

```bash
$ docker run -d --name=marco -p 0.0.0.0:80:80 nickschuch/marco
```
