# RISU (Rapid Image Supplying Unit)
Risu is a build tool for docker image with original cache mechanism.

![](https://cloud.githubusercontent.com/assets/261700/9082260/c51e910c-3b9d-11e5-9202-f0ab05207ac6.png)

## Quick Start
First, run risu server as docker container.

```bash
$ docker run \
    --name risu \
    -e GITHUB_ACCESS_TOKEN=XXXXXXXXXXXXXXXXXXXXXX \
    -e DOCKER_AUTH_USER_NAME=your_name \
    -e DOCKER_AUTH_USER_PASSWORD=your_password \
    -e DOCKER_AUTH_USER_EMAIL=your_email \
    -p 8080:8080 \
    -v /var/run/docker.sock:/var/run/docker.sock \
    quay.io/wantedly/risu:latest
```

Second, trigger a new build via risu API.

```bash
$ curl -n -X POST https://<your-risu-server>.com/builds \
  -H "Content-Type: application/json" \
 \
  -d '{
  "source_repo": "wantedly/risu",
  "source_branch": "master",
  "image_name": "quay.io/wantedly/risu:latest",
  "dockerfile": "Dockerfile.dev",
  "cache_directories": [
    {
      "source": "vendor/bundle",
      "container": "/app/vendor/bundle"
    },
    {
      "source": "vendor/assets",
      "container": "/app/vendor/assets"
    }
  ]
}'
```

Then, risu server build docker image with original cache mechanism and push it to docker registry.

That's it!

## Documentation
### HTTP API

* [Build](https://github.com/wantedly/risu/blob/master/docs/api-v1-alpha.md#build)
 * [Create](https://github.com/wantedly/risu/blob/master/docs/api-v1-alpha.md#build-create)
 * [Info](https://github.com/wantedly/risu/blob/master/docs/api-v1-alpha.md#build-info)
 * [List](https://github.com/wantedly/risu/blob/master/docs/api-v1-alpha.md#build-list)

## Requirements

* `GITHUB_ACCESS_TOKEN`
* `DOCKER_AUTH_USER_NAME`
* `DOCKER_AUTH_USER_PASSWORD`
* `DOCKER_AUTH_USER_EMAIL`


## How It Works
TBD

## Registry Backend
### localfs Backend Registry

```bash
$ docker run \
    --name risu \
    -e GITHUB_ACCESS_TOKEN=XXXXXXXXXXXXXXXXXXXXXX \
    -e DOCKER_AUTH_USER_NAME=your_name \
    -e DOCKER_AUTH_USER_PASSWORD=your_password \
    -e DOCKER_AUTH_USER_EMAIL=your_email \
    -p 8080:8080 \
    -v /var/run/docker.sock:/var/run/docker.sock \
    quay.io/wantedly/risu:latest
```

### etcd Backend Registry

```bash
$ docker run \
    --name risu \
    -e GITHUB_ACCESS_TOKEN=XXXXXXXXXXXXXXXXXXXXXX \
    -e DOCKER_AUTH_USER_NAME=your_name \
    -e DOCKER_AUTH_USER_PASSWORD=your_password \
    -e DOCKER_AUTH_USER_EMAIL=your_email \
    -e REGISTRY_BACKEND=etcd \
    -e REGISTRY_ENDPOINT=http://172.17.8.101:4001 \
    -p 8080:8080 \
    -v /var/run/docker.sock:/var/run/docker.sock \
    quay.io/wantedly/risu:latest
```

## Cache Backend
### localfs

```bash
$ docker run \
    --name risu \
    -e GITHUB_ACCESS_TOKEN=XXXXXXXXXXXXXXXXXXXXXX \
    -e DOCKER_AUTH_USER_NAME=your_name \
    -e DOCKER_AUTH_USER_PASSWORD=your_password \
    -e DOCKER_AUTH_USER_EMAIL=your_email \
    -p 8080:8080 \
    -v /var/run/docker.sock:/var/run/docker.sock \
    quay.io/wantedly/risu:latest
```

### S3

```bash
$ docker run \
    --name risu \
    -e GITHUB_ACCESS_TOKEN=XXXXXXXXXXXXXXXXXXXXXX \
    -e DOCKER_AUTH_USER_NAME=your_name \
    -e DOCKER_AUTH_USER_PASSWORD=your_password \
    -e DOCKER_AUTH_USER_EMAIL=your_email \
    -e CACHE_BACKEND=s3 \
    -e AWS_ACCESS_KEY_ID=XXXXXXXXXXXXXXXXXXXX \
    -e AWS_SECRET_ACCESS_KEY=xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx \
    -e AWS_REGION=xx-yyyy-0 \
    -e RISU_CACHE_BUCKET=xxxx \
    -p 8080:8080 \
    -v /var/run/docker.sock:/var/run/docker.sock \
    quay.io/wantedly/risu:latest
```

## Contribution

1. Fork it ( http://github.com/wantedly/risu )
2. Create your feature branch (git checkout -b my-new-feature)
3. Commit your changes (git commit -am 'Add some feature')
4. Push to the branch (git push origin my-new-feature)
5. Create new Pull Request
