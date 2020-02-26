FROM amd64/golang:1.13-alpine AS builder
RUN apk --update add --no-cache git
RUN export GOBIN=$HOME/work/bin
WORKDIR /go/src/app
ADD src/ .
RUN go get -d -v ./...
RUN CGO_ENABLED=0 go build -o main .

FROM amd64/alpine:3.11
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
RUN adduser -S -D -H -h /app appuser
COPY --from=builder /go/src/app/main /app/
ADD res/ /app/res/
RUN mkdir /app/certs
RUN chown -R appuser /app
USER appuser
VOLUME /app/certs
WORKDIR /app
CMD ["./main"]