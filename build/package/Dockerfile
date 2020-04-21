FROM alpine:3.11

ENV WEBPATH "/"
ENV WEBLOG "off"
ENV UPDATEKEY "changeme"

COPY dashgoat /app/dashgoat
COPY web/ /web/
#RUN ls -la /web/*

EXPOSE 1323

CMD ["sh", "-c", "/app/dashgoat -webpath ${WEBPATH}  -weblog ${WEBLOG} -updatekey ${UPDATEKEY}"]