# syntax=docker/dockerfile:1

FROM golang:alpine as builder

# Add tools
RUN apk add git

WORKDIR /build
RUN mkdir ./target

# Download necessary Go modules
COPY go.mod go.sum ./
RUN go mod download

# Copy src
COPY . .

# Build with compile-time parameters
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags '-extldflags "-static"' -o ./target ./...

# Get clean app bin
FROM scratch
WORKDIR /app
COPY --from=builder /build/target/ ./

EXPOSE 8080
EXPOSE 8443

CMD ["./app"]