Marco [![Build Status](https://travis-ci.org/nickschuch/marco.svg?branch=master)](https://travis-ci.org/nickschuch/marco)
=====

A simple proxy for Docker containers.

![Overview](/docs/overview.png "Overview")

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
$ go build
```

**Run the binary with**

```
$ sudo ./marco
```

**Run with all the args**

```
$ sudo ./marco --port=8080 --docker-ports=80,8983,8080 --docker-endpoint=tcp://localhost:2375
```

**Run inside a Docker container**

```
$ docker pull nickschuch/marco
$ docker run -d -p 0.0.0.0:80:80 -v /var/run/docker.sock:/var/run/docker.sock nickschuch/marco
```

#### Running a container

As mentioned in the example above Marco requires containers to be run with 2 options:
* A DOMAIN environment variable
* A port to be exposed

Here is an example of running a container which meets the criteria above

```
docker run -d -m 128m --publish-all=true -e "DOMAIN=test.dev" google/golang-hello
```

Note: The flag --publish-all exposes port 8080 on this container (as per the Dockerfile).

### Drivers

#### Backends

Anything that results in a list of http paths.

Can be passed with the `--backend` flag.

We currently only support Docker but looking to support something like:

* AWS EC2
* AWS ECS
* Tutum

#### Balancer

These are types of load balancers. Currently we only support "Round robin".

Can be passed with the `--balancer` flag.

### Why?

I created this proxy for 2 reasons.
* Local development with more than one container that needs to run on a single HTTP port (eg. 80)
* Load balance across multiple hosts with containers ready to serve a single site. Powered by Docker Swarm.

![Why](/docs/why.png "Why")

### Demo

<a href="http://www.youtube.com/watch?feature=player_embedded&v=2pzwmtCeSyQ
" target="_blank"><img src="http://img.youtube.com/vi/2pzwmtCeSyQ/0.jpg" 
alt="Marco - demo " width="240" height="180" border="10" /></a>

### Roadmap

* Http auth (maybe)
* Pluggable load balancers
