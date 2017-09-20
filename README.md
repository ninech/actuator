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

## Openshift

### Deploy Actuator

Actuator can easily be deployed with the provided Openshift template.

    oc create -f https://raw.githubusercontent.com/ninech/actuator/master/template.yml
    oc new-app actuator -p ACTUATOR_DOMAIN=actuator.example.com -p GITHUB_ACCESS_TOKEN=<secret token from github>
    oc policy add-role-to-user edit -z actuator

After that you can change Actuator's configuration. Every change to the config needs a new deployment.

    oc edit cm actuator
    oc rollout latest actuator

### Template Parameters

The following parameters can be used in a template. They get automatically filled.

| Parameter          | Description     |
| :----------------- | :-------------- |
| `BRANCH_NAME`      | The name of the pull request branch. This is provided by the Github webhook event. |

## API

| Endpoint           | Description     |
| :----------------- | :-------------- |
| `GET /v1/health`   | Returns status code `200` and a JSON formatted message if the server is running. |
| `POST /v1/event`   | Post a [Github Pull Request event](https://developer.github.com/v3/activity/events/types/#pullrequestevent) and trigger the appropriate action. |

## Development

### Call the hook

There is an example event payload in `examples/pull-request-event.json`. You can use that to test against the event handler API.

    curl -vX POST http://localhost:8080/v1/event \
         -d @examples/pull-request-event.json \
         --header "Content-Type: application/json"

    # or using HTTPie
    http POST localhost:8080/v1/event @examples/pull-request-event.json

For the above noted commands to work you will need to provide a signed secret. To make that easier during development there is a wrapper script to calculate the signature from the secret and call the hook:

    $ examples/send-event.rb

### Test template

First in your Openshift project import the sample template:

    oc create -f examples/test-template.yml

Whenever this template is applied there will be a new ConfigMap with a partly random name `actuator-test-*` in your project. These ConfigMaps can be deleted with the following command:

    oc delete cm -l "actuator.nine.ch/create-reason"
