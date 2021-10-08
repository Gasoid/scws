# Static Content Web Server
The main purpose of the project is to develop static server that can be used with modern javascript frameworks (React, vue.js and so on)


### Features:
- Jaeger support
- Prometeus metrics


### Storage types:
- local filesystem
- aws s3


## Variables

| Variable Name  | Description |
| ------------- | ------------- |
| SCWS_INDEX_HTML | index file (default: "/index.html") |
| SCWS_STORAGE | storage type: filesystem, s3 (default: "filesystem") |
| SCWS_PORT | port (default: "8080") |
| SCWS_FS_ROOT | root path for filesystem (default: "/www/") |
| SCWS_S3_BUCKET | s3 bucket where content is |
| SCWS_S3_PREFIX | s3 prefix where content is (default: "") |
| SCWS_S3_AWS_ACCESS_KEY_ID | please set up SCWS_S3_AWS_ACCESS_KEY_ID, SCWS_S3_AWS_SECRET_ACCESS_KEY and AWS_REGION if storage type is "s3" |
| SCWS_S3_AWS_SECRET_ACCESS_KEY |  |
| SCWS_S3_AWS_REGION | REGION |
| JAEGER_AGENT_HOST | jaeger host (default: "localhost") |
| JAEGER_AGENT_PORT | jaeger port (default: "6831") |
| JAEGER_TAGS | jaeger tags (default: "") |
| JAEGER_SERVICE_NAME | jaeger service name (default: "") |

**Jaeger lib has more variables.** Please check its github readme https://github.com/jaegertracing/jaeger-client-go


## Useful URLS
- /_/health
- /_/metrics


## Docker image
```bash
docker pull ghcr.io/gasoid/scws:latest
```

## Docker example

```dockerfile
FROM node:16-stretch AS demo
WORKDIR /code
RUN git clone https://github.com/ahfarmer/calculator.git
RUN cd calculator && npm install && npm run build


FROM ghcr.io/gasoid/scws:latest
WORKDIR /root/
RUN apk --no-cache --update add bash curl less jq openssl
COPY --from=demo /code/calculator/build /www/calculator
CMD SCWS_INDEXHTML="calculator/index.html" /root/scws
```
