FROM alpine:latest

RUN mkdir /app /certificates

COPY ./bin/oauth /app

COPY ./certificates /certificates

CMD ["/app/oauth"]