Marco
=====

An API driver load balancer for modern day application infrastructure.

![Diagram](/docs/diagram.png "Diagram")

### Balancers and backends

**Balancer**

* Round - We are using https://github.com/mailgun/oxy under the hood for Round Robin balancing.

**Drivers**

* AWS ECS - https://github.com/nickschuch/marco-ecs

### Usage

#### Binary

```bash
$ marco
INFO[0000] Balancing connections on port 80             
INFO[0000] Receiving backend data on port 81 
```

#### Docker

```bash
$ docker run -d --name=marco -p 0.0.0.0:80:80 nickschuch/marco
INFO[0000] Balancing connections on port 80             
INFO[0000] Receiving backend data on port 81 
```

### Libraries

* https://github.com/nickschuch/marco-lib - Used for backend services to leverage for consistent code when pushing to Marco.
