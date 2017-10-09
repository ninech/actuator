[![Build Status](https://travis-ci.org/ninech/actuator.svg?branch=master)](https://travis-ci.org/ninech/actuator)

# Actuator

More documentation can be found at https://ninech.github.io/actuator.

> **actuator** (noun  ac·tu·a·tor \ˈak-chə-ˌwā-tər, -shə-\ ) a mechanical device for moving or controlling something

Actuator keeps an eye on your Github pull-requests and automatically spawns staging environments in Openshift.

It's in a very early stage of development. Some people would say it's just a prototype. So please come back later if you're interested in using Actuator.


## Development

### Call the hook

There is an example event payload in `examples/pull-request-event-opened.json`. You can use that to test against the event handler API.

Because the event handler validates the HTTP request, there is a wrapper script to calculate the signature from the secret and call the hook:

    examples/send-event.rb -f examples/pull-request-event-opened.json -s supersecret

    examples/send-event.rb -f examples/pull-request-event-closed.json -s supersecret

### Test template

Import the sample template into your Openshift project:

    oc create -f examples/test-template.yml

Whenever this template is applied there will be a new ConfigMap with a partly random name `actuator-test-*` in your project. These ConfigMaps can be deleted with the following command:

    oc delete cm -l "actuator.nine.ch/create-reason"

## About

This tool is currently maintained and funded by [nine](https://nine.ch).

[![logo of the company 'nine'](https://logo.apps.at-nine.ch/Dmqied_eSaoBMQwk3vVgn4UIgDo=/trim/500x0/logo_claim.png)](https://www.nine.ch)
