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
* A "random" load balanced connection between _container1_ and _container2_
* A proxy connection to _container3_
* No proxy connection will be setup for _container4_ given it doesn't have a DOMAIN environment variable set.
* No proxy connection will be setup for _container5_ given it doesn't have a Port exposedt.

### How to run

**Build the binary with**

```
$ go build marco.go
```

**Run the binary with**

```
$ sudo ./marco
```

**Run with all the args**

```
$ sudo ./marco -bind=8080 -ports=80,8983,8080 -endpoint=tcp://localhost:2375
```

* bind - Server traffic through the following port.
* ports - The Docker exposed ports that this proxy can use (in order).
* endpoint - Connection to the Docker daemon.

**Run inside a Docker container**

```
$ docker pull nickschuch/marco
$ docker run -d -p 0.0.0.0:80:80 -v /var/run/docker.sock:/var/run/docker.sock nickschuch/marco
```

### Why?

I created this proxy for 2 reasons.
* To be able to implement this on a local development host with multiple containers.
* So loadbalancers don't have to be on the same network as the other Docker containers.

![Why](/docs/why.png "Why")

### Demo

_I will post a video very soon..._

### Roadmap

* Http auth (maybe)
* Logging
* Error handling
