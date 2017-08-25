[![Build Status](https://travis-ci.org/ninech/frau-schultz.svg?branch=master)](https://travis-ci.org/ninech/frau-schultz)

# Frau Schultz

Frau Schultz keeps an eye on your Github pull-requests and automatically spawns staging environments in Openshift.

It's in a very early stage of development. Some people would say it's just a prototype. So please come back later if you're interested in using Frau Schultz.

```sh
$ frau-schultz
Hallo, ich bin Frau Schultz!
[GIN] 2017/08/25 - 09:04:54 | 200 | 506.684µs | ::1 | GET /v1/health
```

## API

| Endpoint           | Description     |
| :----------------- | :-------------- |
| `GET /v1/health`   | Returns status code `200` and a JSON formatted message if the server is running. |
