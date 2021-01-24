FROM golang:1.14.13-buster AS builder

WORKDIR /go/src/app

# Create folder structure
RUN mkdir -p bin/
RUN mkdir -p logs/

# Install dependencies
COPY ./go.mod ./go.mod
COPY ./go.sum ./go.sum
RUN go mod download

# Copy source code
COPY ./cmd ./cmd
COPY ./config ./config
COPY ./internal ./internal
COPY ./pkg ./pkg

# Build binary
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o bin/yorpoll ./cmd/yorpoll/main.go && \
    chmod +x bin/yorpoll

FROM alpine:3.13.0 AS app
WORKDIR /app
COPY --from=builder /go/src/app/bin/yorpoll ./yorpoll
COPY ./scripts /app/scripts
COPY ./swaggerui /app/swaggerui
CMD ["./yorpoll"]
