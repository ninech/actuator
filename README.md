[![Build Status](https://travis-ci.org/ninech/actuator.svg?branch=master)](https://travis-ci.org/ninech/actuator)

# Actuator

> **actuator** (noun  ac·tu·a·tor \ˈak-chə-ˌwā-tər, -shə-\ ) a mechanical device for moving or controlling something

Actuator keeps an eye on your Github pull-requests and automatically spawns staging environments in Openshift.

It's in a very early stage of development. Some people would say it's just a prototype. So please come back later if you're interested in using Actuator.

```sh
$ actuator
Hello, I am actuator!
[GIN] 2017/08/25 - 09:04:54 | 200 | 506.684µs | ::1 | GET /v1/health
```

## API

| Endpoint           | Description     |
| :----------------- | :-------------- |
| `GET /v1/health`   | Returns status code `200` and a JSON formatted message if the server is running. |
| `POST /v1/event`   | Post a [Github Pull Request event](https://developer.github.com/v3/activity/events/types/#pullrequestevent) and trigger the appropriate action. |

## Development

There is an example event payload in `examples/pull-request-event.json`. You can use that to test against the event handler API.

    curl -vX POST http://localhost:8080/v1/event \
         -d @examples/pull-request-event.json \
         --header "Content-Type: application/json"

    # or using HTTPie
    http POST localhost:8080/v1/event @examples/pull-request-event.json
