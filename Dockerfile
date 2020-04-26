FROM golang:1.12-alpine AS build 
ENV GO111MODULE on
ENV CGO_ENABLED 0

RUN apk add git make openssl

WORKDIR /go/src/github.com/4sh/k8s-toleration-injection-webhook
ADD . .
RUN make test
RUN make app

FROM scratch
WORKDIR /app
COPY --from=build /go/src/github.com/4sh/k8s-toleration-injection-webhook/toleration-injection-server .
ENTRYPOINT ["/app/toleration-injection-server"]
