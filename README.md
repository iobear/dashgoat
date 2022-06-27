# dashgoat

Dashgoat - A simple dashboard, easy to deploy.

![Alt dashgoat](doc/dashgoat.png?raw=true "DashGoat")

## Features

 * Easy to use
 * Configuration management frendly
 * Non hierarchical cluster option
 * Lightweight
 * HTTP(s) only, no special ports/protocols

## Golang Hello world

```go build  ./cmd/dashgoat```

```./dashgoat -updatekey my-precious!```

curl API example;

```bash
curl --request POST \
  --url http://127.0.0.1:1323/update \
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
```http://127.0.0.1:1323```

Update status to error;

```bash
curl --request POST \
  --url http://127.0.0.1:1323/update \
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

### Lost heartbeat

If you expect regular updates from a service, and you want to keep track of the service updates, you can use the ```nextupdatesec``` parameter, this will warn you if dashGoat is missing updates within the seconds defined.

```bash
curl --request POST \
  --url http://127.0.0.1:1323/update \
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


## Docker/Podman Hello world

```docker run -e UPDATEKEY=my-precious! -p 1323:1323 --rm --name=dashgoat analogbear/dashgoat```


## Buddy system

dashGoat can have a buddy, to share state or just to gossip to.

 1. When defining a Buddy, dashGoat will at start-up ask for a full list of service states from its buddy.
 2. When receiving an update, dashGoat will forward the update to its buddies.
 3. If dashGoats buddy is down, it will spool the updates, and tell buddy later when its back.


### Buddy, hello world

To run it on your local machine, you can expose two different ports, for two different instances.
Start your first instance:
```bash
./dashgoat -updatekey my-precious! -buddyurl http://localhost:2001
```
Have a look at your browser again:
```http://127.0.0.1:1323```

There should be something about "My buddy is down" in the dashboard.
Start your second instance:

```bash
./dashgoat -updatekey my-precious! -buddyurl http://localhost:1323 -ipport :2001
```
Your first dashboard should be happy now. If you check your new dashboard at ```http://localhost:2001```, it should say "Waiting for first update".

Now try doing the same updates as before, and you should see both dashGoat instances update, on both port 2001, and 1323.

If you want more buddies, you can define them in a list, in the dashgoat.yaml file, instead of using the -buddyurl parameter.

## Full Api
For a full API feature list, go to the doc folder and import the ```dashGoat.postman_collection.json``` file to Postman, Insomnia or Paw. Or read the json file :-)

## TODO

 * Better auth
 * Automatic cleanup option
 * lots more..
