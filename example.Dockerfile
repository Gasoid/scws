FROM golang:1.16.2-stretch AS builder
WORKDIR /code
ADD go.mod /code/
#ADD go.sum /code/
RUN go mod tidy
ADD . /code/
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -o /code/scws .


FROM node:16-stretch AS demo
WORKDIR /code
RUN git clone https://github.com/ahfarmer/calculator.git
RUN cd calculator && npm install && npm run build


FROM alpine:3.6
WORKDIR /root/
RUN apk --no-cache --update add bash curl less jq openssl
COPY --from=builder /code/scws /root/
COPY --from=demo /code/calculator/build /www/calculator
CMD /root/scws
