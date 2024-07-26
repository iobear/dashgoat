# Run this file Makefile command "make docker"
# Build stage
FROM golang:alpine AS build-env
RUN apk add build-base git
WORKDIR /go/src/dashgoatbuild

COPY cmd/dashgoat/* ./
COPY go.* ./
COPY web/ ./

# Set ldflags using build arguments
ARG VERSION
ARG COMMIT
ARG BUILD_DATE
ENV VERSION=$VERSION
ENV COMMIT=$COMMIT
ENV BUILD_DATE=$BUILD_DATE

ARG LDFLAGS="-X 'main.Version=${VERSION}' -X 'main.Commit=${COMMIT}' -X 'main.BuildDate=${BUILD_DATE}'"

RUN go get -d -v github.com/labstack/echo/...
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="${LDFLAGS}" -a -installsuffix cgo -o dashgoat .

# Final stage
FROM alpine:latest
COPY --from=build-env /go/src/dashgoatbuild/dashgoat /app/dashgoat

EXPOSE 2000
CMD ["/app/dashgoat"]

#docker build . --progress=plain -t analogbear/dashgoat:test -f build/package/Dockerfile
