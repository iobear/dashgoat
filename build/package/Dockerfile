# Fetch necessary values from git
ARG VERSION=$(git tag -l --sort=-creatordate | head -n 1)
ARG COMMIT=$(git describe --always)
ARG BUILD_DATE=$(date +%Y-%m-%d)

# Build stage
FROM golang:alpine AS build-env
RUN apk add build-base git
WORKDIR /go/src/dashgoatbuild
COPY cmd/dashgoat/* /go/src/dashgoatbuild/
COPY go.* /go/src/dashgoatbuild/
COPY web/ /go/src/dashgoatbuild/web/

# Set ldflags using build arguments
ENV LDFLAGS="-X 'main.Version=${VERSION}' -X 'main.Commit=${COMMIT_SHA1}' -X 'main.BuildDate=${BUILD_DATE}'"

RUN go get -d -v github.com/labstack/echo/...
RUN go get github.com/iobear/dashgoat/common

RUN ls -lah
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="${LDFLAGS}" -a -installsuffix cgo -o dashgoat .

# Final stage
FROM alpine:latest
COPY --from=build-env /go/src/dashgoatbuild/dashgoat /app/dashgoat

EXPOSE 2000

CMD ["/app/dashgoat"]

#docker build . --progress=plain -t analogbear/dashgoat:test -f build/package/Dockerfile
