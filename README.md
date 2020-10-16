# hyper-cas

[![Build Status](https://travis-ci.org/vtex/hyper-cas.svg?branch=master)](https://travis-ci.org/vtex/hyper-cas) [![Coverage Status](https://coveralls.io/repos/github/vtex/hyper-cas/badge.svg?branch=master)](https://coveralls.io/github/vtex/hyper-cas?branch=master) [![Docker](https://img.shields.io/docker/cloud/build/vtexcom/hyper-cas?label=Docker&style=flat)](https://hub.docker.com/r/vtexcom/hyper-cas/builds)

hyper-cas is a Content-Addressable Storage aimed at JAMStack websites.

## Usage - hyper-cas API

TODO.

## Usage - Synchronizing a website to the CAS

TODO.

## Dev

In order to run tests `make test`.

For getting the API up and running, `make serve`.

For serving the websites published to the API `make route`. This requires docker and will download and run an Nginx container.

For running the API with Docker, `make docker-serve`.
