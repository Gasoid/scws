FROM golang:1.16.2-stretch AS builder
WORKDIR /code
ADD go.mod /code/
#ADD go.sum /code/
RUN go mod download
ADD . /code/
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -o /code/scws .
RUN chmod a+x /code/scws


FROM alpine:3.6
WORKDIR /root/
RUN apk --no-cache --update add bash curl less jq openssl
COPY --from=builder /code/scws /usr/local/bin/scws
CMD /usr/local/bin/scws
