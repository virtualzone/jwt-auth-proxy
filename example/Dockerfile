FROM node:13-alpine AS frontend_builder
RUN mkdir -p /usr/src/app
WORKDIR /usr/src/app
ADD frontend/. .
RUN npm install
RUN npm run build

FROM golang:1.13-alpine AS backend_builder
RUN apk --update add --no-cache git
RUN export GOBIN=$HOME/work/bin
WORKDIR /go/src/app
ADD backend/. .
RUN go get -d -v ./...
RUN CGO_ENABLED=0 go build -o main .

FROM alpine:3.11
RUN adduser -S -D -H -h /app appuser
USER appuser
COPY --from=frontend_builder /usr/src/app/build/. /app/static/
COPY --from=backend_builder /go/src/app/main /app/
WORKDIR /app
CMD ["./main"]