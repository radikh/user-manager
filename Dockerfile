FROM golang:1.13 as modules
ADD go.mod go.sum /m/
RUN cd /m && go mod download && mkdir -p /go/pkg/

FROM golang:1.13 as builder
ENV GO111MODULE on
ENV CGO_ENABLED 0
COPY --from=modules /go/pkg/ /go/pkg
RUN mkdir -p /opt/resource/ 
WORKDIR /opt/resource/
COPY . .
WORKDIR /opt/resource/cmd/umserver
RUN go build -o /opt/services/user-manager/cmd/umserver .

FROM alpine:3.7
COPY --from=builder /opt/services/user-manager/cmd/umserver /opt/services/user-manager
# Usage template for email sending
COPY server/mail/mail_template server/mail/mail_template
RUN chmod +x /opt/services/user-manager
CMD /opt/services/user-manager
