Marco
=====

A simple proxy for Docker containers.

![Overview](/docs/overview.png "Overview")

### Example

Running on the domain **example.com**

I have the containers named:
* **container1** - HTTP exposed on port 80
* **container2** - HTTP exposed on port 8983

The router will proxy these using the following hosts:
* **container1.example.com**
* **container2.example.com**

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
* Local development with more than one container that needs to run on a single HTTP port (eg. 80)
* Load balance across multiple hosts with containers ready to serve a single site. Powered by Docker Swarm.

![Why](/docs/why.png "Why")

### Demo

_I will post a video very soon..._

### Roadmap

* Test suite
  * Unit tests
  * Functional tests
* Http auth (maybe)
* Logging
* Error handling
