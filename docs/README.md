# Actuator

> **actuator** (noun  ac·tu·a·tor \ˈak-chə-ˌwā-tər, -shə-\ ) a mechanical device for moving or controlling something

Actuator makes it possible to have a staging environment of your application for every pull-request.

It runs in a container in your application's [Openshift](https://www.openshift.com/) project. You can then point a webhook in your repository (only Github at the moment) to you Actuator route. Whenever someone opens a new pull-request in your repository, Actuator applies a predefined template. When there is a route in the template it also posts the url to the new environment on your pull-request.


# Installation

Actuator can easily be deployed with the provided Openshift template.

It needs edit privileges in your project. This means you have to add the service account to the `edit` role. Note: Don't do this if your production environment is running in the same project.

```bash
oc process -f https://raw.githubusercontent.com/ninech/actuator/master/template.yml \
  -p ACTUATOR_DOMAIN=actuator.example.com \
  -p GITHUB_ACCESS_TOKEN=supersecrettoken | oc create -f -

oc policy add-role-to-user edit -z actuator
```

Replace the values for `ACTUATOR_DOMAIN` and `GITHUB_ACCESS_TOKEN`. See bellow what their meaning is.

After that you can change Actuator's configuration. Every change to the config needs a new deployment.

```bash
oc edit cm actuator
oc rollout latest actuator
```


# Configuration

Most of the configuration happens in `actuator.yml`. There is an example config in this repo. The basic values are explained here.

```yaml
---
# a secret value you share with github. it needs to be entered in the webhook form on github.
github_webhook_secret: supersecret

# the token which is used to post comments on github. you have to create a token in your github settings: https://github.com/settings/tokens.
github_access_token: alsoverysecret

# configurations for every repository you'd like to handle.
repositories:
-
  # should this repository be handled? turn to `false` to disable it temporarily.
  enabled: true
  # full name of the repo in the form of `owner/repo-name`
  fullname: ninech/actuator-demo
  # the openshift template to apply whenever actuator receives a hook. it has to be present in your openshift project.
  template: actuator-demo
```

You can also configure some values via environment variables:

| Config                         | Description     |
| :----------------------------- | :-------------- |
| `ACTUATOR_WEBHOOK_SECRET`      | The webhook secret you share with Github. |
| `ACTUATOR_GITHUB_ACCESS_TOKEN` | The Github access token to create comments on Github. |

Environment variables take precedence over configurations from the config file!

# Template Parameters

On [Openshift template accepts parameters](https://docs.openshift.com/online/dev_guide/templates.html#writing-parameters). For now Actuator uses just one parameter. It is planned to add more dynamic and configurable parameters in the future.

The following parameter can be used in a template. It is automatically provided by Actuator.

| Parameter          | Description     |
| :----------------- | :-------------- |
| `BRANCH_NAME`      | The name of the pull request's head ref. This is provided by the Github webhook event. It can be used to name the objects. (ex. BuildConfig `actuator-${BRANCH_NAME}`) |

# API

| Endpoint           | Description     |
| :----------------- | :-------------- |
| `GET /v1/health`   | Returns status code `200` and a JSON formatted message if the server is running. |
| `POST /v1/event`   | Post a [Github Pull Request event](https://developer.github.com/v3/activity/events/types/#pullrequestevent) and trigger the appropriate action. |
