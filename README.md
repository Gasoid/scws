# Static Content Web Server
The main purpose of the project is to develop static server that can be used with modern javascript frameworks (React, vue.js and so on)

You don't need to use nginx in order to serve your react project.

## Features:
- Settings, env var exporting
- Jaeger support
- Prometeus metrics
- Vault support

## Settings
SCWS allows you to export all env variables that begin with SCWS_SETTINGS_VAR_ (it is configurable).
You can use this feature to provide to your javascript application a config.

E.g. we can create following Dockerfile file:
```dockerfile
FROM node:16-stretch AS demo
WORKDIR /code
RUN git clone https://github.com/Gasoid/test-client.git
ENV REACT_APP_SETTINGS_API_URL="http://127.0.0.1:8080/_/settings"
RUN cd test-client && npm install && npm run build


FROM ghcr.io/gasoid/scws:latest
WORKDIR /root/
RUN apk --no-cache --update add bash curl less jq openssl
COPY --from=demo /code/test-client/build/ /www/
CMD SCWS_SETTINGS_VAR_WEBSITE="example.com" SCWS_SETTINGS_VAR_GOOGLE_MAP_IS_ENABLED="false" scws
```
after building and we can run it:
```bash
docker build -t scws:example -f example.Dockerfile
docker run --rm -p 8080:8080  scws:example
```
if we run `curl 127.0.0.1:8080/_/settings`, we will see folloing

```json
{"GOOGLE_MAP_IS_ENABLED":"false","WEBSITE":"example.com"}
```
By default settings path is `/_/settings`. Prefix is `SCWS_SETTINGS_VAR_`


## Prometheus metrics
path is `/_/metrics`
- `scws_response_size_bytes` histogram
- `scws_requests_total` counter

## Health check
Scws checks whether index.html exists or not and accordingly responds 200 or 500 code.
health check path is `/_/health`


## Storage types:
- local filesystem
- aws s3


## Variables

| Variable Name  | Default value | Description |
| ------------- | ------------- | ------------- |
| SCWS_INDEX_HTML | "/index.html" | index file |
| SCWS_STORAGE | "filesystem" |storage type: filesystem, s3 |
| SCWS_PORT | "8080" | port |
| SCWS_FS_ROOT | "/www/" | root path for filesystem |
| SCWS_SETTINGS_PREFIX | "SCWS_SETTINGS_VAR_" | prefix for env variables, which will be exposed for client, you can get it from /_/settings as json. e.g. SCWS_SETTINGS_VAR_WEBSITE="mycoolwebsite" |
| SCWS_SETTINGS_PATH | "/_/settings" | path for settings |
| SCWS_VAULT_ADDRESS | "" | vault address, e.g. http://vault:8200/ |
| SCWS_VAULT_PATHS | "" | list of paths, e.g. "secrets/aws/scws,secrets/aws/scws2" |
| SCWS_VAULT_TOKEN | "" | vault token |
| SCWS_S3_BUCKET | "" | s3 bucket where content is |
| SCWS_S3_PREFIX | "" | s3 prefix where content is |
| SCWS_S3_AWS_ACCESS_KEY_ID | "" | please set up SCWS_S3_AWS_ACCESS_KEY_ID, SCWS_S3_AWS_SECRET_ACCESS_KEY and AWS_REGION if storage type is "s3" |
| SCWS_S3_AWS_SECRET_ACCESS_KEY | "" |  |
| SCWS_S3_AWS_REGION | "" | REGION |
| JAEGER_AGENT_HOST | "localhost" | jaeger host |
| JAEGER_AGENT_PORT | "6831" | jaeger port |
| JAEGER_TAGS | "" | jaeger tags |
| JAEGER_SERVICE_NAME | "" | jaeger service name |

**Jaeger lib has more variables.** Please check its github readme https://github.com/jaegertracing/jaeger-client-go


## Useful URLS
- /_/health
- /_/metrics
- /_/settings


## Docker image
```bash
docker pull ghcr.io/gasoid/scws:latest
```

## Docker example

```dockerfile
FROM node:16-stretch AS demo
WORKDIR /code
RUN git clone https://github.com/Gasoid/test-client.git
ENV REACT_APP_SETTINGS_API_URL="http://127.0.0.1:8080/_/settings"
RUN cd test-client && npm install && npm run build


FROM ghcr.io/gasoid/scws:latest
WORKDIR /www/
RUN apk --no-cache --update add bash curl less jq openssl
COPY --from=demo /code/test-client/build/ /www/
CMD SCWS_SETTINGS_VAR_WEBSITE="example.com" SCWS_SETTINGS_VAR_GOOGLE_MAP_IS_ENABLED="false" scws
```

## K8s/helm example for settings
values.yml:
```yaml
# ...
react_variables:
  IS_MAP_ENABLED: true
  TITLE: "example"
# ...
```

deployment.yml:
```yaml
#...
      containers:
        - name: {{ .Chart.Name }}
          securityContext:
            {{- toYaml .Values.securityContext | nindent 12 }}
          {{- if .Values.image.tag }}
          image: "{{ .Values.image.repository }}:{{ .Values.image.tag | default .Chart.AppVersion }}"
          {{- end}}
          imagePullPolicy: {{ .Values.image.pullPolicy }}
          env:
          {{- range $key, $val := .Values.react_variables }}
          - name:  SCWS_SETTINGS_VAR_{{ $key }}
            value: {{ $val | quote }}
          {{- end }}
#...
```

## License
This program is published under the terms of the MIT License. Please check the LICENSE file for more details.
