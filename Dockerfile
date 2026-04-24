FROM golang:1.20 as builder
WORKDIR /src

# Install build deps for tree-sitter (gcc etc) - present in official image
RUN apt-get update && apt-get install -y build-essential ca-certificates git \
    && rm -rf /var/lib/apt/lists/*

COPY go.mod go.sum ./
RUN go mod download

COPY . .

# Build with the treesitter tag so tree-sitter adapter is included
RUN CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build -tags treesitter -o /out/codemind ./

FROM debian:bookworm-slim
COPY --from=builder /out/codemind /usr/local/bin/codemind
ENTRYPOINT ["/usr/local/bin/codemind"]
