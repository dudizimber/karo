# Build stage
FROM golang:1.21-alpine AS builder

WORKDIR /workspace

# Copy go mod and sum files
COPY go.mod go.mod
COPY go.sum go.sum

# Download dependencies
RUN go mod download

# Copy the source code
COPY api/ api/
COPY controllers/ controllers/
COPY webhook/ webhook/
COPY main.go main.go

# Build the binary
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -o manager main.go

# Runtime stage
FROM gcr.io/distroless/static:nonroot

WORKDIR /

COPY --from=builder /workspace/manager .

USER 65532:65532

ENTRYPOINT ["/manager"]
