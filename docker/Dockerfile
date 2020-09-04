FROM alpine:3.7
RUN mkdir /app
WORKDIR /app
COPY simple .
RUN ls -l
EXPOSE 8080:8080

CMD [ "./simple","-p","8080" ]