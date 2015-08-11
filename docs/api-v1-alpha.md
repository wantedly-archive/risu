## <a name="resource-build"></a>Build

A build represents an individual build job for docker image

### Attributes

| Name | Type | Description | Example |
| ------- | ------- | ------- | ------- |
| **id** | *uuid* | unique identifier of build | `"01234567-89ab-cdef-0123-456789abcdef"` |
| **source_repo** | *string* | source github source_repositry to build. It must includes Dockerfile | `"wantedly/risu"` |
| **source_revision** | *string* | git revision to use for build. (also you can use git tag/branch for this). | `"ada9ce1829fab49e605e5a563dbf91274f64e923"` |
| **name** | *string* | a source_repository name (and optionally a tag) to apply to the resulting image in case of success. | `"quay.io/wantedly/risu:latest"` |
| **dockerfile** | *string* | path within the build context to the Dockerfile | `"Dockerfile.dev"` |
| **status** | *string* | status of build. one of "failed" or "building" or "succeeded" | `"failed"` |
| **created_at** | *date-time* | when build was created | `"2015-01-01T12:00:00Z"` |
| **updated_at** | *date-time* | when build was updated | `"2015-01-01T12:00:00Z"` |

### Build Create

Create a new build.

```
POST /builds
```

#### Optional Parameters

| Name | Type | Description | Example |
| ------- | ------- | ------- | ------- |
| **source_repo** | *string* | source github source_repositry to build. It must includes Dockerfile | `"wantedly/risu"` |
| **source_revision** | *string* | git revision to use for build. (also you can use git tag/branch for this). | `"ada9ce1829fab49e605e5a563dbf91274f64e923"` |
| **name** | *string* | a source_repository name (and optionally a tag) to apply to the resulting image in case of success. | `"quay.io/wantedly/risu:latest"` |
| **dockerfile** | *string* | path within the build context to the Dockerfile | `"Dockerfile.dev"` |


#### Curl Example

```bash
$ curl -n -X POST https://<your-risu-server>.com//builds \
  -H "Content-Type: application/json" \
 \
  -d '{
  "source_repo": "wantedly/risu",
  "source_revision": "ada9ce1829fab49e605e5a563dbf91274f64e923",
  "name": "quay.io/wantedly/risu:latest",
  "dockerfile": "Dockerfile.dev"
}'
```


#### Response Example

```
HTTP/1.1 201 Created
```

```json
{
  "id": "01234567-89ab-cdef-0123-456789abcdef",
  "source_repo": "wantedly/risu",
  "source_revision": "ada9ce1829fab49e605e5a563dbf91274f64e923",
  "name": "quay.io/wantedly/risu:latest",
  "dockerfile": "Dockerfile.dev",
  "status": "failed",
  "created_at": "2015-01-01T12:00:00Z",
  "updated_at": "2015-01-01T12:00:00Z"
}
```

### Build Info

Info for existing build.

```
GET /builds/{build_id}
```


#### Curl Example

```bash
$ curl -n https://<your-risu-server>.com//builds/$BUILD_ID
```


#### Response Example

```
HTTP/1.1 200 OK
```

```json
{
  "id": "01234567-89ab-cdef-0123-456789abcdef",
  "source_repo": "wantedly/risu",
  "source_revision": "ada9ce1829fab49e605e5a563dbf91274f64e923",
  "name": "quay.io/wantedly/risu:latest",
  "dockerfile": "Dockerfile.dev",
  "status": "failed",
  "created_at": "2015-01-01T12:00:00Z",
  "updated_at": "2015-01-01T12:00:00Z"
}
```

### Build List

List existing builds.

```
GET /builds
```


#### Curl Example

```bash
$ curl -n https://<your-risu-server>.com//builds
```


#### Response Example

```
HTTP/1.1 200 OK
```

```json
[
  {
    "id": "01234567-89ab-cdef-0123-456789abcdef",
    "source_repo": "wantedly/risu",
    "source_revision": "ada9ce1829fab49e605e5a563dbf91274f64e923",
    "name": "quay.io/wantedly/risu:latest",
    "dockerfile": "Dockerfile.dev",
    "status": "failed",
    "created_at": "2015-01-01T12:00:00Z",
    "updated_at": "2015-01-01T12:00:00Z"
  }
]
```


