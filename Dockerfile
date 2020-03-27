FROM golang:1.13 as builder
RUN mkdir -p /go/src/github.com/lvl484
ENV GO111MODULE on
ENV CGO_ENABLED 0
WORKDIR /go/src/github.com/lvl484/user-manager
COPY go.mod ./
COPY go.sum ./
RUN go mod download
COPY . .
WORKDIR /go/src/github.com/lvl484/user-manager/cmd/umserver
RUN mkdir -p /opt/services/ && go build -o /opt/services/user-manager

FROM alpine:3.7
COPY --from=builder /opt/services/user-manager /opt/services/user-manager/user-manager
COPY config/viper.config.json /opt/services/user-manager/config/
WORKDIR /opt/services/user-manager