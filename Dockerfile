FROM golang:1.23.7 AS build
# Set the working directory
WORKDIR /build

RUN go env -w GOMODCACHE=/root/.cache/go-build

# Install dependencies
COPY go.mod go.sum ./
RUN --mount=type=cache,target=/root/.cache/go-build go mod download

COPY . ./
# Build the server
RUN --mount=type=cache,target=/root/.cache/go-build CGO_ENABLED=0 go build -trimpath -ldflags '-s -w' \
    -o mcp-corapanel main.go

# Make a stage to run the app
FROM alpine:3.21.3
# Set the working directory
WORKDIR /server
# Copy the binary from the build stage
COPY --from=build /build/mcp-corapanel .
# Command to run the server
CMD ["./mcp-corapanel", "stdio"]
