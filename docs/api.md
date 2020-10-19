# Using hyper-cas API

hyper-cas has a restful API that can be used to interact with the storage. Running the API is as simple as:

```
$ hyper-cas serve --config hyper-cas.yaml
```

For more information on the configuration file, check the [Configuration](config.md) page.

# Routes & Payloads

## Healthcheck

hyper-cas comes bundled with a healthcheck route so it is easy to understand whether the API is up and running.

### Request

- Method: `GET`
- URL: `/healthcheck`
- Body: `none`

### Response

```
 $ curl -vvv http://localhost:2485/healthcheck
 *   Trying 127.0.0.1...
 * TCP_NODELAY set
 * Connected to localhost (127.0.0.1) port 2485 (#0)
 > GET /healthcheck HTTP/1.1
 > Host: localhost:2485
 > User-Agent: curl/7.58.0
 > Accept: */*
 >
 < HTTP/1.1 200 OK
 < Server: fasthttp
 < Date: Mon, 19 Oct 2020 21:35:46 GMT
 < Content-Type: text/plain; charset=utf-8
 < Content-Length: 2
 <
 * Connection #0 to host localhost left intact
 OK
```

## File Storage

These are APIs meant to handle files stored in the CAS. You can either store (`PUT`), retrieve a file (`GET`) or verify if a hash is in the CAS (`HEAD`).

### Storing a file

#### Request

- Method: `PUT`
- URL: `/file`
- Body: `the contents of the file to be stored`

#### Response

```
$ curl -vvv -XPUT --data "test1" http://localhost:2485/file
*   Trying 127.0.0.1...
* TCP_NODELAY set
* Connected to localhost (127.0.0.1) port 2485 (#0)
> PUT /file HTTP/1.1
> Host: localhost:2485
> User-Agent: curl/7.58.0
> Accept: */*
> Content-Length: 5
> Content-Type: application/x-www-form-urlencoded
>
* upload completely sent off: 5 out of 5 bytes
< HTTP/1.1 200 OK
< Server: fasthttp
< Date: Mon, 19 Oct 2020 21:39:49 GMT
< Content-Type: text/plain; charset=utf-8
< Content-Length: 40
<
* Connection #0 to host localhost left intact
b444ac06613fc8d63795be9ad0beaf55011936ac
```

The response is the SHA1 hash of the contents of the file.

### Retrieving a file

> **⚠ WARNING: This API is just for DEBUG purposes.**  
> The retrieval of the files is synchronous and not cached. For production, please use NGINX or another production-grade webserver.

#### Request

- Method: `GET`
- URL: `/file/{hash}`
    - `hash`: the SHA1 hash of the file you want to retrieve from the CAS
- Body: `none`

#### Response

```
$ curl -vvv "http://localhost:2485/file/b444ac06613fc8d63795be9ad0beaf55011936ac"
*   Trying 127.0.0.1...
* TCP_NODELAY set
* Connected to localhost (127.0.0.1) port 2485 (#0)
> GET /file/b444ac06613fc8d63795be9ad0beaf55011936ac HTTP/1.1
> Host: localhost:2485
> User-Agent: curl/7.58.0
> Accept: */*
>
< HTTP/1.1 200 OK
< Server: fasthttp
< Date: Mon, 19 Oct 2020 21:47:19 GMT
< Content-Type: text/plain; charset=utf-8
< Content-Length: 5
<
* Connection #0 to host localhost left intact
test1
```

### Verifying if file already in CAS

By doing a `HEAD` request for the file hash you can verify if a file is in the CAS and if it is there's no need to upload it again.

#### Request

- Method: `HEAD`
- URL: `/file/{hash}`
    - `hash`: the SHA1 hash of the file you want to verify
- Body: `none`

#### Response

```
$ curl -vvv -I "http://localhost:2485/file/b444ac06613fc8d63795be9ad0beaf55011936ac"
*   Trying 127.0.0.1...
* TCP_NODELAY set
* Connected to localhost (127.0.0.1) port 2485 (#0)
> HEAD /file/b444ac06613fc8d63795be9ad0beaf55011936ac HTTP/1.1
> Host: localhost:2485
> User-Agent: curl/7.58.0
> Accept: */*
>
< HTTP/1.1 200 OK
HTTP/1.1 200 OK
< Server: fasthttp
Server: fasthttp
< Date: Mon, 19 Oct 2020 21:49:32 GMT
Date: Mon, 19 Oct 2020 21:49:32 GMT

<
* Connection #0 to host localhost left intact
```

If the file exists, you'll get `200` status code, `404` otherwise.

## Distribution Storage

These are APIs meant to manage distributions. Distributions in hyper-cas are [Merkle Trees](https://en.wikipedia.org/wiki/Merkle_tree) of files in specific paths. This means that if a file content changes, or their path changes, we get a new distribution tree.

Distributions are the base for how hyper-cas allows many versions of a given website (or repository for that matter) to exist at the same time, while not wasting useful resources with repeated files.

They are basically a map of where in the CAS storage all files for a given version of the website are.

Whenever a distribution is created, hyper-cas will create a mapping of symbolic links in the filesystem pointing each file to a hash content in the CAS. This is the magic that allows us to have infinite version in the system.

### Storing a distribution

#### Request

- Method: `PUT`
- URL: `/distro`
- Body: `one line for each file in the distribution, with the format of PATH:HASH`

```
multiline.txt
where/my/file/is:b444ac06613fc8d63795be9ad0beaf55011936ac
where/my/other/file/is:b444ac06613fc8d63795be9ad0beaf55011936ac
some/other/file:b444ac06613fc8d63795be9ad0beaf55011936ac
```

Notice that it is fine to have the same file in multiple locations.

#### Response

```
curl -vvv -XPUT --data-binary "@multiline.txt" -H "Expect:" -H'Content-Type: text/html; charset=utf-8' http://localhost:2485/distro
*   Trying 127.0.0.1...
* TCP_NODELAY set
* Connected to localhost (127.0.0.1) port 2485 (#0)
> PUT /distro HTTP/1.1
> Host: localhost:2485
> User-Agent: curl/7.58.0
> Accept: */*
> Content-Type: text/html; charset=utf-8
> Content-Length: 179
>
* upload completely sent off: 179 out of 179 bytes
< HTTP/1.1 200 OK
< Server: fasthttp
< Date: Mon, 19 Oct 2020 21:58:49 GMT
< Content-Type: text/plain; charset=utf-8
< Content-Length: 40
<
* Connection #0 to host localhost left intact
768706dd535495cd5e64b94c5a603244b21237d3
```

The response is the SHA1 hash of the contents of the distribution tree.

### Retrieving a Distribution

> **⚠ WARNING: This API is just for DEBUG purposes.**  
> The retrieval of the distribution is synchronous and not cached. For production, please use NGINX or another production-grade webserver.

### Request

- Method: `GET`
- URL: `/distro/{hash}`
    - `hash`: the SHA1 hash of the distribution tree you want to get
- Body: `none`

### Response

```
$ curl -vvv "http://localhost:2485/distro/768706dd535495cd5e64b94c5a603244b21237d3"
*   Trying 127.0.0.1...
* TCP_NODELAY set
* Connected to localhost (127.0.0.1) port 2485 (#0)
> GET /distro/768706dd535495cd5e64b94c5a603244b21237d3 HTTP/1.1
> Host: localhost:2485
> User-Agent: curl/7.58.0
> Accept: */*
>
< HTTP/1.1 200 OK
< Server: fasthttp
< Date: Mon, 19 Oct 2020 22:02:22 GMT
< Content-Type: text/plain; charset=utf-8
< Content-Length: 186
<
* Connection #0 to host localhost left intact
["where/my/file/is:b444ac06613fc8d63795be9ad0beaf55011936ac","where/my/other/file/is:b444ac06613fc8d63795be9ad0beaf55011936ac","some/other/file:b444ac06613fc8d63795be9ad0beaf55011936ac"]
```

### Verifying if distribution already in CAS

By doing a `HEAD` request for the distribution hash you can verify if it is in the CAS and if it is there's no need to upload it again.

### Request

- Method: `HEAD`
- URL: `/distro/{hash}`
    - `hash`: the SHA1 hash of the distribution tree you want to verify
- Body: `none`

### Response

```
$ curl -vvv -I "http://localhost:2485/distro/768706dd535495cd5e64b94c5a603244b21237d3"
*   Trying 127.0.0.1...
* TCP_NODELAY set
* Connected to localhost (127.0.0.1) port 2485 (#0)
> HEAD /distro/768706dd535495cd5e64b94c5a603244b21237d3 HTTP/1.1
> Host: localhost:2485
> User-Agent: curl/7.58.0
> Accept: */*
>
< HTTP/1.1 200 OK
HTTP/1.1 200 OK
< Server: fasthttp
Server: fasthttp
< Date: Mon, 19 Oct 2020 22:02:27 GMT
Date: Mon, 19 Oct 2020 22:02:27 GMT

<
* Connection #0 to host localhost left intact
```

If the file exists, you'll get `200` status code, `404` otherwise.

## Label Storage

Labels are pointers to distributions in hyper-cas. They are your gateway into files, since it is very simple to point a label to another distribution (rollback or roll forward).

Whenever a label is created or updated, hyper-cas will generate the according nginx configuration file, mapping to a given distribution root path.

TODO: Write the rest of the API.
