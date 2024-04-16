# dashGoat

dashGoat - A simple HTTP dashboard, easy to deploy.

![Alt dashgoat](doc/dashgoat.png?raw=true "DashGoat")

[CHANGELOG](CHANGELOG.md) [API](doc/API.md) [k8s](/deploy/kubernetes)

## Features

 * Easy to use
 * Configuration management friendly
 * Non hierarchical cluster option
 * Lightweight
 * HTTP(s) only, no special protocols

## Golang

`make build` or download a binary from the releases

`./dashgoat -updatekey my-precious!`

## Docker

```docker run -e UPDATEKEY=my-precious! -p 2000:2000 --rm --name=dashgoat analogbear/dashgoat```

curl API example;

```bash
curl --request POST \
  --url http://127.0.0.1:2000/update \
  --header 'content-type: application/json' \
  --data '{
	"host": "host-1",
	"service": "HTTP",
	"status": "ok",
	"message": "Hello World",
	"updatekey": "my-precious!"
}'
```

Check your browser on:
`http://127.0.0.1:2000`

Update status to error;

```bash
curl --request POST \
  --url http://127.0.0.1:2000/update \
  --header 'content-type: application/json' \
  --data '{
	"host": "host-1",
	"service": "HTTP",
	"status": "error",
	"message": "Hello World",
	"updatekey": "my-precious!"
}'
```
Check web page again.

## Watchdog - Lost heartbeat

If you expect regular updates from a service, and you want to keep track of the service updates, you can use the `nextupdatesec` parameter, this will warn you if dashGoat is missing updates within the seconds defined.

```bash
curl --request POST \
  --url http://127.0.0.1:2000/update \
  --header 'content-type: application/json' \
  --data '{
	"host": "host-1",
	"service": "rsync",
	"status": "ok",
	"message": "",
	"updatekey": "my-precious!",
	"nextupdatesec": 20
}'
```
Now you should get a warning if the update is missing for 20 seconds. This feature is not meant to be super fast (< 10 sec), this is just to keep track of "lost" agents.

If POST is not possible you can use GET like this:

```bash
curl http://127.0.0.1:2000/heartbeat/hdjsakl678dsa/router01/20/dsl,openwrt,home
```

The input is as follows
```bash
curl http://127.0.0.1:2000/heartbeat/<urnkey>/<host>/<nextupdatesec>/<tags>
```

When using HTTP GET you need to update the config with:
`urnkey: <key>` or use the enviroment variable `URNKEY=<key>`

## TTL

If you want your event to change state/disappear after a set amount of seconds, use the `ttl` parameter, like this.

```bash
curl --request POST \
  --url http://127.0.0.1:2000/update \
  --header 'content-type: application/json' \
  --data '{
	"host": "host-1",
	"service": "uptime",
	"status": "ok",
	"message": "Server has rebooted, msg gone in 10sec",
	"updatekey": "my-precious!",
	"ttl": 10
}'
```
Default behaviour (PromoteToOk)
 * if status == "ok" the event vil disappear, after ttl expire.
 * if status != "ok", state vil change to "ok", after ttl expire.
Via config there are three other modes

Remove
 * If ttl expires, serviceStatus is removed, not waiting for `TtlOkDelete`

PromoteOnce
 * If ttl expires, status moved one to the right on the list `["critical", "error", "warning", "info", "ok"]` If status becomes `ok` serviceStatus waits for  `TtlOkDelete` expires.

PromoteOneStep
 * Like promoteOnce but status keeps promoting along the list everytime ttl expires until it ends at `"ok"`

## Tags
!! TODO - Frontend is not done, no filter API !!

Tags are used to filter sevices depending on their tags, this way you can, as an example list services associated with specific customers or departments. You could also ad tags related to the service the server is running. Lets say you are running a transcoder service, the service is transcoding 4 channels, these channels can then be added as a list, like this:

```bash
curl --request POST \
  --url http://127.0.0.1:2000/update \
  --header 'content-type: application/json' \
  --data '{
	"host": "trans-1",
	"service": "nginx",
	"status": "ok",
	"message": "Smooth",
	"updatekey": "my-precious!",
	"tags": ["tr-ch1","tr-ch2","tr-ch3","tr-ch4"]
}'
```
You do that to all your transcoders, and now you can list them according to their channels, or whatever tag you use. These are also very handy when defining dependencies(DependOn)


## DependOn

DepenOn is a parameter you can add to your service updates, this will reduce the important alerts to the systems the services depend on. Like this:

```bash
curl --request POST \
  --url http://127.0.0.1:2000/update \
  --header 'content-type: application/json' \
  --data '{
	"host": "host-1",
	"service": "nginx",
	"status": "error",
	"message": "No requests",
	"updatekey": "my-precious!",
	"dependon": "loadbalancer-1"
}'
```
In this case if `loadbalancer-1` is down, all the services that has `"dependon": "loadbalancer-1"` will reduce status to `info` until its up again. If you have more that one server your service depends on then you can also use tags, the value is checked for matches with both hosts and tags.
<br /> In the above `Tags` example instead of using `dependon:trans-1`, you can use the ch2 tag `dependon:tr-ch2` and dashGoat will check if there is other services with the same tag that is up, and will only say 1/X is down. When setup correctly, this reduces events with `error` and `critical` and only show "upstream" errors.

## Alertmanager

You can forward the alerts from your diffrent alertmanagers to display on a dashgoat central screen or share with other systems. You need to set the `urnkey` in your dashgoat config for this to work.

In alertmanager config, you need to add a webhook reciever:
```yaml
    receivers:
      - name: Dashgoat
        webhook_configs:
        - send_resolved: true
          url: https://<dashgoat host>/alertmanager/<urnkey>
    route:
      - receiver: Dashgoat
        repeat_interval: 5m
        matchers:
          - alertname != Watchdog
```
## Prometheus

`/metrics` is exposed by default, you can enable to show a host-service timeline in dashGoat by providing a url for your prometheus instance.
Either by using a enviroment variable,
 `PROMETHEUSURL=http://localhost:9090 ./dashgoat`
Or adding a line to the `dashgoat.yaml`
 `prometheusurl: http://localhost:9090`

To show the timeline in dashGoat select the change time and the timeline will appear, this feature needs some rework in the UI.

## PagerDuty

You can forward events to PagerDutys diffrent tecnical services depending on hosts + services, tags and severity, there is also a catch all.
<br />Minimum config, catch all:
```yaml
pagerdutyconfig:
  pagerdutyservicemaps:
    - hostservice:
      tag:
      eapikey: 12345acc39464b01d0105f1234567890
```
All options looks like this, first matching the host `host-1` and service `cache` to one technical service key, second matching all services with tag `customer23` to another tecnical service key. Both only being forwarded to PageDuty if severity is `error` or higher.
```yaml
pagerdutyconfig:
  url: https://events.pagerduty.com/v2/enqueue
  timeout: 10s
  triggerlevel: error
  pagerdutymode: push
  pagerdutyservicemaps:
    - hostservice: host-1cache
      tag:
      eapikey: 12345acc39464b01d0105f1234567890
    - hostservice:
      tag: customer23
      eapikey: ffff5acc39464b01d0105f123456ffff
```

## Buddy system

dashGoat can have a buddy to share state, or just gossip to.

 1. When defining a Buddy, dashGoat will at start-up ask for a full list of service states from its buddy.
 2. When receiving an update, dashGoat will forward the update to its buddies.
 3. If dashGoats buddy is down, it will spool the updates, and tell buddy later when it's back.


### Buddy, hello world

To run it on your local machine, you can expose two different ports, for two different instances.
Start your first instance:
```bash
./dashgoat -updatekey my-precious! -buddyurl http://localhost:2001
```
Have a look at your browser again:
`http://127.0.0.1:2000`

There should be something about "My buddy is down" in the dashboard.
Start your second instance:

```bash
./dashgoat -updatekey my-precious! -buddyurl http://localhost:2000 -ipport :2001
```
Your first dashboard should be happy now. If you check your new dashboard at `http://localhost:2001`, it should say "Waiting for first update".

Now try doing the same updates as before, and you should see both dashGoat instances update, on both port 2001, and 2000.

If you want more buddies, you can define them in a list, in the dashgoat.yaml file, instead of using the -buddyurl parameter.

### Docker, Buddy Hello world

So for docker you can't use localhost, as every Docker container has it own .. So to compensate, use the IP on your network card instead.

first instance:

```docker run  -e BUDDYURL=http://<local-nic-ip>:2001 -p <local-nic-ip>:2000:2000 --rm --name=dashgoat analogbear/dashgoat```

Second instance:

```docker run  -e BUDDYURL=http://<local-nic-ip>:2000 -p <local-nic-ip>:2001:2000 --rm --name=dashgoat2 analogbear/dashgoat```

### k8s Buddies

For Kubernetes you can use the files from the `deploy` folder as inspiration, if you add the headless service "dashgoat-headless-svc", dashGoat will find its buddies via DNS. You can point to another DNS-name via the `nsconfig` option.

## Full Api
For a full API feature list, go to the doc folder and import the `dashGoat.postman_collection.json` file to Postman, Insomnia or Paw. Or read the [API](doc/API.md) file.

## Docker build

If you want to build your own Docker container, you can use the Dockerfile, with the included GO build environment.

```bash
docker build -f build/package/Dockerfile -t myDashgoat
```

To include the config file:
 1. comment-in the two copy commands in the Dockerfile
 2. edit dashgoat.yaml
 3. copy yaml to cmd/dashgoat/
 4. run "docker build..."


## TODO

 * Delete event on dashboard via mouse
 * MS teams support
 * API tests (in progress)
 * Configuration tests 
 * Save state
 * Tags filter view / filter API
 * Users +gravatar?
 * Ack event
 * Auth on delete
 * dashGoat client
 * clean up main.go