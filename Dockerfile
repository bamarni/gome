FROM golang:1

RUN mkdir /var/www
VOLUME /var/www

RUN mkdir -p /go/src/app
WORKDIR /go/src/app
COPY gome.go /go/src/app

RUN go-wrapper download
RUN go-wrapper install
CMD ["go-wrapper", "run"]
