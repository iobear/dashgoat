# dashgoat

Dashgoat - A simple dashboard, easy to deploy.

## Installation

Set IP and PORT in main.go, for example IP 127.0.0.1 and PORT 1323 

`./dashgoat`

## Hello world

curl API example;

`curl --request POST 
  --url http://127.0.0.1:1323/dg/update 
  --header 'content-type: application/json' 
  --data '{
	"host": "myhost01",
	"service": "HTTP",
	"status": "ok",
	"message": "Hello World"
}'`

Check your browser on:
http://127.0.0.1:1323

Update status to error;

`curl --request POST 
  --url http://127.0.0.1:1323/dg/update
  --header 'content-type: application/json' 
  --data '{
	"host": "myhost01",
	"service": "HTTP",
	"status": "error",
	"message": "Hello World"
}'`

Check web page again.