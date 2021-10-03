FROM node:16-stretch AS demo
WORKDIR /code
RUN git clone https://github.com/ahfarmer/calculator.git
RUN cd calculator && npm install && npm run build


FROM ghcr.io/gasoid/scws:latest
WORKDIR /root/
RUN apk --no-cache --update add bash curl less jq openssl
COPY --from=demo /code/calculator/build /www/calculator
CMD SCWS_INDEXHTML="calculator/index.html" /root/scws
