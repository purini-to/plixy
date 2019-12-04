###############################
# Builder container
###############################
FROM golang:1 as builder

ENV GOOS=linux
ENV GOARCH=amd64
ENV TZ=Asia/Tokyo
ENV LANG=ja_JP.UTF-8
ENV LANGUAGE=ja_JP.UTF-8
ENV LC_ALL=ja_JP.UTF-8

ENV APP_DIR=/var/lib/plixy
RUN mkdir -p $APP_DIR

WORKDIR $APP_DIR

RUN groupadd -r plixy && useradd --no-log-init -r -g plixy plixy

ADD go.mod go.sum ./

RUN go mod download

ADD . .

RUN CGO_ENABLED=0 go build -ldflags '-d -w -s' -o plixy

###############################
# Run container
###############################
FROM scratch

ENV TZ=Asia/Tokyo
ENV LANG=ja_JP.UTF-8
ENV LANGUAGE=ja_JP.UTF-8
ENV LC_ALL=ja_JP.UTF-8

ENV APP_DIR /bin

WORKDIR $APP_DIR

COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo
COPY --from=builder /etc/passwd /etc/passwd
COPY --from=builder /usr/share/zoneinfo/Asia/Tokyo /usr/share/zoneinfo

COPY --from=builder /var/lib/plixy/plixy.yaml $APP_DIR/plixy.yaml
COPY --from=builder /var/lib/plixy/plixy $APP_DIR/plixy

USER plixy

EXPOSE 8080

ENTRYPOINT ["/bin/plixy"]
