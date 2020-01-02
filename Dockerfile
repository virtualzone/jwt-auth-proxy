FROM golang:1.13-alpine AS builder
RUN apk --update add --no-cache git
RUN export GOBIN=$HOME/work/bin
WORKDIR /go/src/app
ADD src/ .
RUN go get -d -v ./...
RUN CGO_ENABLED=0 go build -o main .

FROM alpine:3.11
RUN adduser -S -D -H -h /app appuser
USER appuser
COPY --from=builder /go/src/app/main /app/
ADD res/ /app/res/
WORKDIR /app
CMD ["./main"]