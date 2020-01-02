FROM alpine:3.11 AS qemu
RUN apk --update add --no-cache curl
RUN cd /tmp && \
    curl -L https://github.com/balena-io/qemu/releases/download/v3.0.0%2Bresin/qemu-3.0.0+resin-arm.tar.gz | tar zxvf - -C . && mv qemu-3.0.0+resin-arm/qemu-arm-static .

FROM arm32v6/golang:1.13-alpine AS builder
COPY --from=qemu /tmp/qemu-arm-static /usr/bin/
RUN apk --update add --no-cache git
RUN export GOBIN=$HOME/work/bin
WORKDIR /go/src/app
ADD src/ .
RUN go get -d -v ./...
RUN CGO_ENABLED=0 go build -o main .

FROM arm32v6/alpine:3.11
ARG BUILD_DATE
ARG VCS_REF
ARG VERSION
LABEL org.label-schema.build-date=$BUILD_DATE \
        org.label-schema.name="JWT Auth Proxy" \
        org.label-schema.description="A lightweight authentication proxy written in Go designed for use in Docker/Kubernetes environments." \
        org.label-schema.vcs-ref=$VCS_REF \
        org.label-schema.vcs-url="https://github.com/virtualzone/jwt-auth-proxy" \
        org.label-schema.version=$VERSION \
        org.label-schema.schema-version="1.0"
COPY --from=qemu /tmp/qemu-arm-static /usr/bin/
RUN adduser -S -D -H -h /app appuser
USER appuser
COPY --from=builder /go/src/app/main /app/
ADD res/ /app/res/
WORKDIR /app
CMD ["./main"]