FROM alpine:3.9

ADD . .
COPY ./ ./

RUN apk update

EXPOSE 8000
EXPOSE 3000

WORKDIR /cmd

CMD ["./service"]
