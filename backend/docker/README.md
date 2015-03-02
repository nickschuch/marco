Docker
=====

For containers to be discovered they need to be run with the following criteria.
* A DOMAIN environment variable
* A port to be exposed

Here is an example of running a container which meets the criteria above

```
docker run -d -m 128m --publish-all=true -e "DOMAIN=test.dev" google/golang-hello
```

Note: The flag --publish-all exposes port 8080 on this container (as per the Dockerfile).

### Demo

Coming soon...

