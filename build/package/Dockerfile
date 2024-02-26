FROM golang:alpine AS build-env
RUN apk add build-base git
WORKDIR /go/src/dashgoatbuild
COPY cmd/dashgoat/* /go/src/dashgoatbuild/
COPY go.* /go/src/dashgoatbuild/
COPY web/ /go/src/dashgoatbuild/web/
COPY common/ /go/src/dashgoatbuild/
#COPY cmd/dashgoat/dashgoat.yaml /go/src/dashgoatbuild/
RUN go get -d -v github.com/labstack/echo/...
RUN go get github.com/iobear/dashgoat/common

RUN ls -lah
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o dashgoat .

#FROM linuxcontainers/debian-slim
#RUN apt-get update && apt-get install -y procps iputils-ping bind9-host --no-install-recommends && rm -rf /var/lib/apt/lists/*
FROM alpine:latest

COPY --from=0 /go/src/dashgoatbuild/dashgoat /app/dashgoat
#COPY --from=0 /go/src/dashgoatbuild/dashgoat.yaml /dashgoat.yaml

EXPOSE 2000

CMD ["sh", "-c", "/app/dashgoat"]

#docker build . --progress=plain -t analogbear/dashgoat:test -f build/package/Dockerfile
