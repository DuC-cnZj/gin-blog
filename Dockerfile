FROM golang:1.14 AS builder
COPY . /app
WORKDIR /app
RUN groupadd -r appuser && useradd --no-log-init -r -g appuser appuser
RUN go env -w GOPROXY=https://goproxy.cn,direct && \
 go mod download

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o /bin/app main.go

FROM scratch
COPY --from=builder /etc/passwd /etc/passwd
COPY --from=builder /bin/app /bin/app
USER appuser
ENTRYPOINT [ "/bin/app" , "serve"]
