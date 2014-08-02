# Contributors App

Display GitHub contributors for a specific repo. This repo represents a sample Go application that can run in a scratch Docker image. The total size of the Docker image including the contributors binary is less than 6MB. 

## Deploying with Docker

First build the binary.

```
CGO_ENABLED=0 GOOS=linux go build -a -ldflags '-w' .
``` 

Build the docker image.

```
docker build -t kelseyhightower/contributors .
```

## Run

```
docker run -d -P kelseyhightower/contributors
```
