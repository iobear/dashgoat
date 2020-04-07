FROM alpine:3.11

ENV webpath "/"
ENV weblog = "off"

COPY dashgoat /app/dashgoat
COPY web/ /web/
#RUN ls -la /web/*

EXPOSE 1323

CMD ["sh", "-c", "/app/dashgoat -webpath ${webpath}  -weblog ${weblog}"]