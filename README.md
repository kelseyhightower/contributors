# Contributors App

Display GitHub contributors for a specific repo.

The Contributors App is designed to run in the scratch Docker image. The total size of the Docker image including the contributors binary is less than 6MB.

- Avoid "x509: failed to load system roots and no roots provided" by bundling root certificates.
- Avoid dynamic linking by using the pure Go net package (-tags netgo)
- Avoid dynamic linking by disabling cgo (CGO_ENABLED=0)
- Reduce binary size by omitting dwarf information (-ldflags '-w')


## Build

### Binary

The following command will produce a statically linked Go binary without debugging (dwarf) information.

```
CGO_ENABLED=0 GOOS=linux go build -a -tags netgo -ldflags '-w' .
``` 

### Docker Image

```
docker build -t kelseyhightower/contributors .
```

## Run

```
docker run -d -P kelseyhightower/contributors
```

## Testing with curl

```
curl --data "repo=docker&owner=docker" http://localhost:49153
```

Replace port 49153 with port from docker ps
