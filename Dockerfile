FROM scratch

COPY build/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY build/zoneinfo-berlin /usr/share/zoneinfo/Europe/Berlin
COPY build/gome /gome
COPY web /web

CMD ["/gome"]
