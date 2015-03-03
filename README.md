Marco [![Build Status](https://travis-ci.org/nickschuch/marco.svg?branch=master)](https://travis-ci.org/nickschuch/marco)
=====

An API driver load balancer for modern day application infrastructure.

![Overview](/docs/overview.png "Overview")

### Balancers and backends

**Balancer**

* First - Will always return the first URL endpoint.
* Round - Will round robin through the list of URL endpoints.

**Drivers**

* [Docker](/backend/docker/README.md) - The Docker daemon and by extension, the Docker Swarm project.
* [Tutum](/backend/tutum/README.md) - https://www.tutum.co
* AWS ECS - Coming soon...

Note: As more providers go into Docker Swarm we won't need to have so many drivers!

### Example

I have the containers:
* **container1**
  * Port 80 exposed
  * Environment variable DOMAIN set to _www.example.com_
* **container2**
  * Port 80 exposed
  * Environment variable DOMAIN set to _www.example.com_
* **container3**
  * Port 8983 exposed
  * Environment variable DOMAIN set to _www.foobar.com_
* **container4**
  * Port 80 exposed
* **container5**
  * Environment variable DOMAIN set to _www.baz.com_

This proxy will setup the following routes:
* A "random" load balanced connection between _container1_ and _container2_ on the domain www.example.com
* A proxy connection to _container3_ on the domain www.foobar.com
* No proxy connection will be setup for _container4_ given it doesn't have a DOMAIN environment variable set.
* No proxy connection will be setup for _container5_ given it doesn't have a Port exposedt.

### Usage

#### Running Marco

**Build the binary with**

```
$ make build
```

**Run the binary with**

```
$ sudo ./marco
```

**Run with on a different port**

```
$ sudo ./marco --port=8080
```

**Run inside a Docker container**

```
$ docker pull nickschuch/marco
$ docker run -d -p 0.0.0.0:80:80 -v /var/run/docker.sock:/var/run/docker.sock nickschuch/marco
```

Please see the CLI for more configuration.

### Drivers

#### Backends

Anything that results in a list of http paths.

Can be passed with the `--backend` flag.

#### Balancer

These are types of load balancers. Currently we support "Round robin" and "First" balancers.

Can be passed with the `--balancer` flag.

### Why?

I created this proxy for 2 reasons.
* Local development with more than one container that needs to run on a single HTTP port (eg. 80)
* Load balance across multiple hosts with containers ready to serve a single site. Powered by Docker Swarm.

![Why](/docs/why.png "Why")
