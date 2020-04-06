FROM alpine:3.11

COPY dashgoat /app/dashgoat

COPY web/ /web/

RUN ls -la /web/*

EXPOSE 1323

CMD ["/app/dashgoat”]