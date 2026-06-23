FROM golang:1.25-alpine AS builder

WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -trimpath -ldflags="-s -w" -o /out/deadliner ./cmd/deadlinerserver

FROM alpine:3.22

WORKDIR /app

RUN apk add --no-cache ca-certificates tzdata

COPY --from=builder /out/deadliner /app/bin/deadliner
COPY conf/config.json /app/conf/config.json
COPY conf/secret.example.json /app/conf/secret.example.json

EXPOSE 8080 8888

ENTRYPOINT ["/app/bin/deadliner"]
