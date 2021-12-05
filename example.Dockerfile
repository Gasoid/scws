FROM node:16-stretch AS demo
WORKDIR /code
RUN git clone https://github.com/Gasoid/test-client.git
ENV REACT_APP_SETTINGS_API_URL="http://127.0.0.1:8080/_/settings"
RUN cd test-client && npm install && npm run build


#FROM ghcr.io/gasoid/scws:latest
FROM scws:check
WORKDIR /root/
RUN apk --no-cache --update add bash curl less jq openssl
COPY --from=demo /code/test-client/build/ /www/
