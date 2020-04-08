# dashgoat

Dashgoat - A simple dashboard, easy to deploy.

## Hello world

```go build  ./cmd/dashgoat```

```./dashgoat -updatekey my-precious!```

curl API example;

```
curl --request POST \
  --url http://127.0.0.1:1323/update \
  --header 'content-type: application/json' \
  --data '{
	"host": "myhost01",
	"service": "HTTP",
	"status": "ok",
	"message": "Hello World",
	"updatekey": "my-precious!"
}'
```

Check your browser on:
http://127.0.0.1:1323

Update status to error;

```
curl --request POST \
  --url http://127.0.0.1:1323/update \
  --header 'content-type: application/json' \
  --data '{
	"host": "myhost01",
	"service": "HTTP",
	"status": "error",
	"message": "Hello World",
	"updatekey": "my-precious!"
}'
```
Check web page again.