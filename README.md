Marco
=====

A simple proxy for Docker containers.

![Overview](/docs/overview.png "Overview")

### Example

Running on the domain example.com

I have the containers named:
* container1 - HTTP exposed on port 80
* container2 - HTTP exposed on port 8983

The router will proxy these using the following hosts:
* container1.example.com
* container2.example.com

### How to run

Build the binary with:

```
$ go build marco.go
```

Run the binary with:

```
$ sudo ./marco
```

Run with all the args:

```
$ sudo ./marco -ports=80,8983,8080 -endpoint=tcp://localhost:2375
```

* ports - The Docker exposed ports that this proxy can use (in order).
* endpoint - Connection to the Docker daemon.

### Why?




### Demo

I will post a video very soon...