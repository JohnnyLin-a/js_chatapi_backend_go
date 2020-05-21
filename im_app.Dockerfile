FROM golang:alpine as builder

ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64 
    # GOARCH=arm

WORKDIR /app

COPY ./cmd/ ./cmd/
COPY ./pkg/ ./pkg/

# Get dependencies
COPY go.mod ./go.sum ./
RUN go mod download

RUN go build ./cmd/chatapisrv/main.go
RUN go build ./cmd/dbsetup/dbsetup.go

FROM alpine:latest
RUN apk --no-cache add ca-certificates

WORKDIR /root/

COPY --from=builder /app/main .
COPY --from=builder /app/dbsetup .
COPY ./.env ./
COPY ./static/ ./static/

ENTRYPOINT ["./main"]

EXPOSE 80 443