FROM golang:1.14-alpine AS builder
COPY . /app
WORKDIR /app
RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.aliyun.com/g' /etc/apk/repositories && \
  apk add --no-cache ca-certificates tzdata
RUN go env -w GOPROXY=https://goproxy.cn,direct && \
 go mod download

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o /bin/app main.go

FROM scratch
COPY --from=builder /usr/share/zoneinfo/Asia/Shanghai /etc/localtime
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /etc/passwd /etc/passwd
COPY --from=builder /bin/app /bin/app

ENTRYPOINT [ "/bin/app" , "serve"]
