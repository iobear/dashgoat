FROM golang:alpine AS build-env
RUN apk --no-cache add build-base git
WORKDIR /go/src/dashgoatbuild
RUN go get -d -v github.com/labstack/echo/...
RUN go get gopkg.in/validator.v2
COPY cmd/dashgoat/* /go/src/dashgoatbuild/
RUN ls -lah
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o dashgoat .

FROM alpine:latest

ENV WEBPATH "/"
ENV WEBLOG "off"
ENV UPDATEKEY "changeme"

COPY --from=0 /go/src/dashgoatbuild/dashgoat /app/dashgoat
COPY web/ /web/
#RUN ls -la /web/*

EXPOSE 1323

CMD ["sh", "-c", "/app/dashgoat -webpath ${WEBPATH}  -weblog ${WEBLOG} -updatekey ${UPDATEKEY}"]