# Build
FROM golang:1.23 as builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .
ARG TARGETOS
ARG TARGETARCH
RUN CGO_ENABLED=0 GOOS=${TARGETOS} GOARCH=${TARGETARCH} go build -ldflags="-s -w" -o go-secrets .

# Runtime 
FROM scratch

WORKDIR /app

COPY --from=builder /app/go-secrets .

ENTRYPOINT ["/app/go-secrets"]
