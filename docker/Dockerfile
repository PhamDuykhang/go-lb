ARG NETWORK_NAME
FROM alpine:3.7

RUN mkdir /app
WORKDIR /app
COPY klb2 .

EXPOSE 8080:8080

CMD [ "./klb2","-np","docker_kapp"]