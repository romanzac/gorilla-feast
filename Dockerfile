FROM alpine:3.16.4

RUN apk update
RUN apk add ca-certificates

COPY gorilla-feast-linux /app/
COPY ./config/gorilla-feast.yaml /app/gorilla-feast.yaml
COPY ./keys /app/
COPY ./scripts/start_gorilla_feast.sh /app/
COPY ./scripts/start_gorilla_feast_wsclient.sh /app/

RUN chmod +x /app/gorilla-feast-linux

ENV PORT 4439
EXPOSE 4439

WORKDIR /app

ENTRYPOINT ["./start_gorilla_feast.sh"]





